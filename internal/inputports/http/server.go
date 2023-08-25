package http

import (
	"context"
	"log"

	rcp "receipt-processor-challenge/internal/domain/receipt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	processPath string = "/process"
	pointsPath  string = "/:id/points"
)

type ReceiptAPI interface {
	SavePoints(ctx context.Context, r rcp.Receipt) (uuid.UUID, error)
	GetPoints(ctx context.Context, id uuid.UUID) (*rcp.Points, error)
}

type Server struct {
	receiptApp ReceiptAPI
	router     *echo.Echo
}

func NewServer(ctx context.Context, app ReceiptAPI) *Server {
	return &Server{
		receiptApp: app,
		router:     echo.New(),
	}
}

func (s *Server) routes() {
	gReceipt := s.router.Group("/receipt")
	gReceipt.POST(processPath, s.saveReceiptPoints)
	gReceipt.GET(pointsPath, s.getReceiptPoints)
}

func (s *Server) Start() {
	s.routes()

	log.Fatal(s.router.Start(":8080"))
}
