package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	rcp "receipt-processor-challenge/internal/domain/receipt"
	"receipt-processor-challenge/internal/interfaceadapters/storage/memory"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrValidator      = errors.New("validator error")
	ErrDecode         = errors.New("decode error")
)

type receipt struct {
	Retailer     string `json:"retailer" validate:"required"`
	PurchaseDate string `json:"purchaseDate" validate:"required,datetime=2006-01-02"`
	PurchaseTime string `json:"purchaseTime" validate:"required,datetime=15:04"`
	Items        []item `json:"items" validate:"required"`
	Total        string `json:"total" validate:"required"`
}

type item struct {
	ShortDescription string `json:"shortDescription" validate:"required"`
	Price            string `json:"price" validate:"required"`
}

type points struct {
	Points int `json:"points"`
}

type id struct {
	ID string `json:"id" validate:"required"`
}

func (s *Server) saveReceiptPoints(e echo.Context) (err error) {
	ctx := e.Request().Context()
	r := new(receipt)
	newID := new(id)

	defer func() {
		if err != nil {
			err = apiReceiptResponse(e, err)
		} else {
			err = apiReceiptResponse(e, newID)
		}
	}()

	bErr := e.Bind(r)
	if bErr != nil {
		err = fmt.Errorf("%s:%w", bErr.Error(), ErrDecode)
		return err
	}

	receipt, err := r.toReceiptDomain()
	if err != nil {
		return err
	}

	uuid, err := s.receiptApp.SavePoints(ctx, *receipt)
	if err != nil {
		return err
	}

	*newID = id{ID: uuid.String()}

	return nil
}

func (s *Server) getReceiptPoints(e echo.Context) (err error) {
	ctx := e.Request().Context()
	paramID := e.Param("id")
	response := new(points)

	defer func() {
		if err != nil {
			err = apiReceiptResponse(e, err)
		} else {
			err = apiReceiptResponse(e, response)
		}
	}()

	err = validate(id{ID: paramID})
	if err != nil {
		return err
	}

	paramUUID, err := uuid.Parse(paramID)
	if err != nil {
		err = fmt.Errorf("%s:%w", err.Error(), ErrDecode)
		return
	}

	p, err := s.receiptApp.GetPoints(ctx, paramUUID)
	if err != nil {
		err = fmt.Errorf("storage error:%w", err)
		return
	}

	if p == nil {
		response = nil
	} else {
		*response = points{Points: (*p).Points}
	}

	return nil
}

func validate(e interface{}) error {
	validate := validator.New()

	err := validate.Struct(e)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return ErrValidator
		}

		var vErr error
		for _, err := range err.(validator.ValidationErrors) {
			field := err.StructField()
			switch err.Tag() {
			case "required":
				vErr = fmt.Errorf("%s is required:%w", field, ErrInvalidRequest)
			case "datetime":
				vErr = fmt.Errorf("%s date/time format:%w", field, ErrInvalidRequest)
			case "number":
				vErr = fmt.Errorf("%s is not numeric:%w", field, ErrInvalidRequest)
			default:
				vErr = fmt.Errorf("%s validation error:%w", field, ErrInvalidRequest)
			}
		}

		return vErr
	}

	return nil
}

func (r receipt) toReceiptDomain() (*rcp.Receipt, error) {
	err := validate(r)
	if err != nil {
		return nil, err
	}

	items := make([]rcp.Item, len(r.Items))
	for index, item := range r.Items {
		newItem, err := item.toItemDomain()
		if err != nil {
			return nil, err
		}

		items[index] = *newItem
	}

	total, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return nil, fmt.Errorf("%s format error:%w", "Total", ErrInvalidRequest)
	}

	purchaseDate, _ := time.Parse(rcp.DatePurchaseFormat, r.PurchaseDate)
	purchasetime, _ := time.Parse(rcp.TimePurchaseFormat, r.PurchaseTime)

	return &rcp.Receipt{
		Retailer:     r.Retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchasetime,
		Items:        items,
		Total:        total,
	}, nil
}

func (i item) toItemDomain() (*rcp.Item, error) {
	err := validate(i)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(i.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("%s format error:%w", "Price", ErrInvalidRequest)
	}

	return &rcp.Item{
		ShortDescription: i.ShortDescription,
		Price:            price,
	}, nil
}

type responseErrorMsg struct {
	Msg string `json:"error"`
}

func apiReceiptResponse(e echo.Context, r interface{}) error {
	switch v := r.(type) {
	case error:
		return apiReceiptResponseError(e, v)
	case *id:
		return e.JSON(http.StatusOK, *v)
	case *points:
		if v == nil {
			return e.JSON(http.StatusNotFound, nil)
		}
		return e.JSON(http.StatusOK, *v)
	}

	return e.JSON(http.StatusInternalServerError, nil)
}

func apiReceiptResponseError(e echo.Context, err error) error {
	code := http.StatusInternalServerError
	jsonErr := responseErrorMsg{Msg: "unexpected error"}

	if errors.Is(err, ErrInvalidRequest) {
		jsonErr.Msg, _ = strings.CutSuffix(err.Error(), fmt.Sprintf(":%s", ErrInvalidRequest.Error()))
		code = http.StatusBadRequest
	}

	if errors.Is(err, ErrDecode) {
		jsonErr.Msg, _ = strings.CutSuffix(err.Error(), fmt.Sprintf(":%s", ErrDecode.Error()))
		code = http.StatusBadRequest
	}

	if errors.Is(err, memory.ErrNotFound) {
		jsonErr.Msg, _ = strings.CutPrefix(err.Error(), "storage error:")
		code = http.StatusNotFound
	}

	return e.JSON(code, jsonErr)
}
