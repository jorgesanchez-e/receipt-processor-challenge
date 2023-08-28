package http

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	rcp "receipt-processor-challenge/internal/domain/receipt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type receiptAPIMock struct {
	mock.Mock
}

func (rcpMock *receiptAPIMock) SavePoints(ctx context.Context, r rcp.Receipt) (uuid.UUID, error) {
	args := rcpMock.Called(ctx, r)

	if id, ok := args.Get(0).(uuid.UUID); ok {
		return id, args.Error(1)
	}

	return uuid.Nil, args.Error(1)
}

func (rcpMock *receiptAPIMock) GetPoints(ctx context.Context, id uuid.UUID) (*rcp.Points, error) {
	args := rcpMock.Called(ctx, id)

	if pnts, ok := args.Get(0).(*rcp.Points); ok {
		return pnts, args.Error(1)
	}

	return nil, args.Error(1)
}

func purchaseDate(t *testing.T, date string) time.Time {
	tdate, err := time.Parse(rcp.DatePurchaseFormat, date)
	if err != nil {
		t.Fatal(err)
	}

	return tdate
}

func purchaseTime(t *testing.T, date string) time.Time {
	tdate, err := time.Parse(rcp.TimePurchaseFormat, date)
	if err != nil {
		t.Fatal(err)
	}

	return tdate
}

func paddingLastByte(t *testing.T) []byte {
	t.Helper()

	b, err := hex.DecodeString("0a")
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func Test_SaveReceiptPoints(t *testing.T) {
	cases := []struct {
		name             string
		contextBuilder   func() (echo.Context, *httptest.ResponseRecorder)
		apiBuilder       func() *receiptAPIMock
		expectedResponse []byte
		expectedHTTPCode int
		expectedError    error
	}{
		{
			name: "retailer-required-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				reqText := []byte(`
								{
									"retailer": "",
									"purchaseDate": "2022-01-01",
									"purchaseTime": "13:01",
									"items": [
									  {
										"shortDescription": "Mountain Dew 12PK",
										"price": "6.49"
									  },{
										"shortDescription": "Emils Cheese Pizza",
										"price": "12.25"
									  },{
										"shortDescription": "Knorr Creamy Chicken",
										"price": "1.26"
									  },{
										"shortDescription": "Doritos Nacho Cheese",
										"price": "3.35"
									  },{
										"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
										"price": "12.00"
									  }
									],
									"total": "35.35"
								  }
								`)

				req := httptest.NewRequest(echo.POST, "http://localhost:8080/process", bytes.NewReader(reqText))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				return echo.New().NewContext(req, rec), rec
			},
			apiBuilder: func() *receiptAPIMock {
				return &receiptAPIMock{}
			},
			expectedResponse: []byte(`{"error":"Retailer is required"}`),
			expectedHTTPCode: http.StatusBadRequest,
			expectedError:    nil,
		},
		{
			name: "purchasedate-format-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				reqText := []byte(`
									{
										"retailer": "Target",
										"purchaseDate": "2022-01-45",
										"purchaseTime": "13:01",
										"items": [
										  {
											"shortDescription": "Mountain Dew 12PK",
											"price": "6.49"
										  },{
											"shortDescription": "Emils Cheese Pizza",
											"price": "12.25"
										  },{
											"shortDescription": "Knorr Creamy Chicken",
											"price": "1.26"
										  },{
											"shortDescription": "Doritos Nacho Cheese",
											"price": "3.35"
										  },{
											"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
											"price": "12.00"
										  }
										],
										"total": "35.35"
									  }
									`)

				req := httptest.NewRequest(echo.POST, "http://localhost:8080/process", bytes.NewReader(reqText))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				return echo.New().NewContext(req, rec), rec
			},
			apiBuilder: func() *receiptAPIMock {
				return &receiptAPIMock{}
			},
			expectedResponse: []byte(`{"error":"PurchaseDate date/time format"}`),
			expectedHTTPCode: http.StatusBadRequest,
			expectedError:    nil,
		},
		{
			name: "purchasetime-format-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				reqText := []byte(`
								{
									"retailer": "Target",
									"purchaseDate": "2022-01-01",
									"purchaseTime": "13::01",
									"items": [
									  {
										"shortDescription": "Mountain Dew 12PK",
										"price": "6.49"
									  },{
										"shortDescription": "Emils Cheese Pizza",
										"price": "12.25"
									  },{
										"shortDescription": "Knorr Creamy Chicken",
										"price": "1.26"
									  },{
										"shortDescription": "Doritos Nacho Cheese",
										"price": "3.35"
									  },{
										"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
										"price": "12.00"
									  }
									],
									"total": "35.35"
								  }
								`)

				req := httptest.NewRequest(echo.POST, "http://localhost:8080/process", bytes.NewReader(reqText))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				return echo.New().NewContext(req, rec), rec
			},
			apiBuilder: func() *receiptAPIMock {
				return &receiptAPIMock{}
			},
			expectedResponse: []byte(`{"error":"PurchaseTime date/time format"}`),
			expectedHTTPCode: http.StatusBadRequest,
			expectedError:    nil,
		},
		{
			name: "no-items-required-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				reqText := []byte(`
								{
									"retailer": "Target",
									"purchaseDate": "2022-01-01",
									"purchaseTime": "13:01",
									"total": "35.35"
								  }
								`)

				req := httptest.NewRequest(echo.POST, "http://localhost:8080/process", bytes.NewReader(reqText))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				return echo.New().NewContext(req, rec), rec
			},
			apiBuilder: func() *receiptAPIMock {
				return &receiptAPIMock{}
			},
			expectedResponse: []byte(`{"error":"Items is required"}`),
			expectedHTTPCode: http.StatusBadRequest,
			expectedError:    nil,
		},
		{
			name: "total-format-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				reqText := []byte(`
								{
									"retailer": "Target",
									"purchaseDate": "2022-01-01",
									"purchaseTime": "13:01",
									"items": [
									  {
										"shortDescription": "Mountain Dew 12PK",
										"price": "6.49"
									  },{
										"shortDescription": "Emils Cheese Pizza",
										"price": "12.25"
									  },{
										"shortDescription": "Knorr Creamy Chicken",
										"price": "1.26"
									  },{
										"shortDescription": "Doritos Nacho Cheese",
										"price": "3.35"
									  },{
										"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
										"price": "12.00"
									  }
									],
									"total": "35.35a"
								  }
								`)

				req := httptest.NewRequest(echo.POST, "http://localhost:8080/process", bytes.NewReader(reqText))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				return echo.New().NewContext(req, rec), rec
			},
			apiBuilder: func() *receiptAPIMock {
				return &receiptAPIMock{}
			},
			expectedResponse: []byte(`{"error":"Total format error"}`),
			expectedHTTPCode: http.StatusBadRequest,
			expectedError:    nil,
		},
		{
			name: "shortDescription-required-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				reqText := []byte(`
								{
									"retailer": "Target",
									"purchaseDate": "2022-01-01",
									"purchaseTime": "13:01",
									"items": [
									  {
										"shortDescription": "",
										"price": "6.49"
									  },{
										"shortDescription": "Emils Cheese Pizza",
										"price": "12.25"
									  },{
										"shortDescription": "Knorr Creamy Chicken",
										"price": "1.26"
									  },{
										"shortDescription": "Doritos Nacho Cheese",
										"price": "3.35"
									  },{
										"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
										"price": "12.00"
									  }
									],
									"total": "35.35"
								  }
								`)

				req := httptest.NewRequest(echo.POST, "http://localhost:8080/process", bytes.NewReader(reqText))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				return echo.New().NewContext(req, rec), rec
			},
			apiBuilder: func() *receiptAPIMock {
				return &receiptAPIMock{}
			},
			expectedResponse: []byte(`{"error":"ShortDescription is required"}`),
			expectedHTTPCode: http.StatusBadRequest,
			expectedError:    nil,
		},
		{
			name: "price-format-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				reqText := []byte(`
								{
									"retailer": "Target",
									"purchaseDate": "2022-01-01",
									"purchaseTime": "13:01",
									"items": [
									  {
										"shortDescription": "Mountain Dew 12PK",
										"price": "6.49x"
									  },{
										"shortDescription": "Emils Cheese Pizza",
										"price": "12.25"
									  },{
										"shortDescription": "Knorr Creamy Chicken",
										"price": "1.26"
									  },{
										"shortDescription": "Doritos Nacho Cheese",
										"price": "3.35"
									  },{
										"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
										"price": "12.00"
									  }
									],
									"total": "35.35"
								  }
								`)

				req := httptest.NewRequest(echo.POST, "http://localhost:8080/process", bytes.NewReader(reqText))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				return echo.New().NewContext(req, rec), rec
			},
			apiBuilder: func() *receiptAPIMock {
				return &receiptAPIMock{}
			},
			expectedResponse: []byte(`{"error":"Price format error"}`),
			expectedHTTPCode: http.StatusBadRequest,
			expectedError:    nil,
		},
		{
			name: "retailer-required-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				reqText := []byte(`
							{
								"retailer": "Test",
								"purchaseDate": "2022-01-01",
								"purchaseTime": "13:01",
								"items": [
								  {
									"shortDescription": "Mountain Dew 12PK",
									"price": "6.49"
								  }
								],
								"total": "35.35"
							  }
							`)

				req := httptest.NewRequest(echo.POST, "http://localhost:8080/process", bytes.NewReader(reqText))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				return echo.New().NewContext(req, rec), rec
			},
			apiBuilder: func() *receiptAPIMock {
				id, _ := uuid.Parse("0a25c541-2ab9-41d9-bedb-2d518df5dc43")

				apiMock := receiptAPIMock{}
				apiMock.On("SavePoints", context.Background(), rcp.Receipt{
					Retailer:     "Test",
					PurchaseDate: purchaseDate(t, "2022-01-01"),
					PurchaseTime: purchaseTime(t, "13:01"),
					Items: []rcp.Item{
						{
							ShortDescription: "Mountain Dew 12PK",
							Price:            6.49,
						},
					},
					Total: 35.35,
				}).Return(id, nil)

				return &apiMock
			},
			expectedResponse: []byte(`{"id":"0a25c541-2ab9-41d9-bedb-2d518df5dc43"}`),
			expectedHTTPCode: http.StatusOK,
			expectedError:    nil,
		},
	}

	for _, c := range cases {
		expectedError := c.expectedError
		expectedResponse := append(c.expectedResponse, paddingLastByte(t)...)
		echoContext, rec := c.contextBuilder()
		s := Server{
			receiptApp: c.apiBuilder(),
		}

		t.Run(c.name, func(t *testing.T) {
			err := s.saveReceiptPoints(echoContext)
			assert.Equal(t, expectedError, err)
			assert.Equal(t, expectedResponse, rec.Body.Bytes())
			assert.Equal(t, c.expectedHTTPCode, rec.Code)
		})
	}
}

func Test_GetReceiptPoints(t *testing.T) {
	cases := []struct {
		name             string
		contextBuilder   func() (echo.Context, *httptest.ResponseRecorder)
		apiBuilder       func() *receiptAPIMock
		expectedResponse []byte
		expectedHTTPCode int
		expectedError    error
	}{
		{
			name: "id-required-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				req := httptest.NewRequest(echo.GET, "http://localhost:8080/:id/points", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				return echo.New().NewContext(req, rec), rec
			},
			apiBuilder: func() *receiptAPIMock {
				return &receiptAPIMock{}
			},
			expectedResponse: []byte(`{"error":"ID is required"}`),
			expectedHTTPCode: http.StatusBadRequest,
			expectedError:    nil,
		},
		{
			name: "invalid-id-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				req := httptest.NewRequest(echo.GET, "http://localhost:8080/:id/points", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e := echo.New().NewContext(req, rec)
				e.SetPath("/:id/points")
				e.SetParamNames("id")
				e.SetParamValues("010101")

				return e, rec
			},
			apiBuilder: func() *receiptAPIMock {
				return &receiptAPIMock{}
			},
			expectedResponse: []byte(`{"error":"invalid UUID length: 6"}`),
			expectedHTTPCode: http.StatusBadRequest,
			expectedError:    nil,
		},
		{
			name: "store-error-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				req := httptest.NewRequest(echo.GET, "http://localhost:8080/:id/points", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e := echo.New().NewContext(req, rec)
				e.SetPath("/:id/points")
				e.SetParamNames("id")
				e.SetParamValues("0a25c541-2ab9-41d9-bedb-2d518df5dc43")

				return e, rec
			},

			apiBuilder: func() *receiptAPIMock {
				id, _ := uuid.Parse("0a25c541-2ab9-41d9-bedb-2d518df5dc43")

				apiMock := receiptAPIMock{}
				apiMock.On("GetPoints", context.Background(), id).Return(nil, errors.New("some-error"))

				return &apiMock
			},

			expectedResponse: []byte(`{"error":"unexpected error"}`),
			expectedHTTPCode: http.StatusInternalServerError,
			expectedError:    nil,
		},
		{
			name: "store-not-found-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				req := httptest.NewRequest(echo.GET, "http://localhost:8080/:id/points", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e := echo.New().NewContext(req, rec)
				e.SetPath("/:id/points")
				e.SetParamNames("id")
				e.SetParamValues("0a25c541-2ab9-41d9-bedb-2d518df5dc43")

				return e, rec
			},

			apiBuilder: func() *receiptAPIMock {
				id, _ := uuid.Parse("0a25c541-2ab9-41d9-bedb-2d518df5dc43")

				apiMock := receiptAPIMock{}
				apiMock.On("GetPoints", context.Background(), id).Return(nil, nil)

				return &apiMock
			},

			expectedResponse: []byte(`null`),
			expectedHTTPCode: http.StatusNotFound,
			expectedError:    nil,
		},
		{
			name: "success-case",
			contextBuilder: func() (echo.Context, *httptest.ResponseRecorder) {
				req := httptest.NewRequest(echo.GET, "http://localhost:8080/:id/points", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e := echo.New().NewContext(req, rec)
				e.SetPath("/:id/points")
				e.SetParamNames("id")
				e.SetParamValues("0a25c541-2ab9-41d9-bedb-2d518df5dc43")

				return e, rec
			},

			apiBuilder: func() *receiptAPIMock {
				id, _ := uuid.Parse("0a25c541-2ab9-41d9-bedb-2d518df5dc43")

				apiMock := receiptAPIMock{}
				apiMock.On("GetPoints", context.Background(), id).Return(&rcp.Points{Points: 10}, nil)

				return &apiMock
			},

			expectedResponse: []byte(`{"points":10}`),
			expectedHTTPCode: http.StatusOK,
			expectedError:    nil,
		},
	}

	for _, c := range cases {
		expectedError := c.expectedError
		expectedResponse := append(c.expectedResponse, paddingLastByte(t)...)
		echoContext, rec := c.contextBuilder()
		s := Server{
			receiptApp: c.apiBuilder(),
		}

		t.Run(c.name, func(t *testing.T) {
			err := s.getReceiptPoints(echoContext)
			assert.Equal(t, expectedError, err)
			assert.Equal(t, expectedResponse, rec.Body.Bytes())
			assert.Equal(t, c.expectedHTTPCode, rec.Code)
		})
	}

}
