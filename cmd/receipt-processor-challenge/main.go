package main

import (
	"context"
	"receipt-processor-challenge/internal/app"
	"receipt-processor-challenge/internal/app/receipt/calculator"
	"receipt-processor-challenge/internal/inputports/http"
	"receipt-processor-challenge/internal/interfaceadapters/storage/memory"
)

func main() {
	ctx := context.Background()
	repo := memory.New(ctx)
	calc := calculator.New(ctx)
	app := app.NewServices(repo, calc)

	server := http.NewServer(ctx, app)
	server.Start()
}
