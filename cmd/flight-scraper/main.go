package main

import (
	"context"
	"os"

	"github.com/chromedp/chromedp"
	"github.com/tarekseba/flight-scraper/internal/logger"
	"github.com/tarekseba/flight-scraper/internal/scraper/scenarios"
	"github.com/tarekseba/flight-scraper/internal/scraper/types"
)


func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),       // Disable headless mode
		chromedp.Flag("disable-gpu", false),    // GPU enabled to improve rendering
		chromedp.Flag("start-maximized", true), // Start in maximized mode
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()


	query := types.Query{
		Weekdays:     map[types.Weekday]bool{types.Weekday(0): true /*, types.Weekday(2): true*/},
		StayDuration: 3,
		MonthHorizon: 0,
		Departure:    "Paris",
		Destination:  "Rome",
	}

	handleQuery := scenarios.HandleQuery{Query: query}

	err := scenarios.LogScenario(&handleQuery)(allocCtx)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		os.Exit(1)
	}
}
