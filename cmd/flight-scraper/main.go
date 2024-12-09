package main

import (
	"context"
	"os"
	"time"

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

	ctx, cancel := chromedp.NewContext(
		allocCtx,
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	logger.InfoLogger.Println("Starting the application")
	logger.InfoLogger.Println("GET " + types.G_FLIGHTS_URL)
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(scenarios.LogScenario(scenarios.NewNavigateToPage(types.G_FLIGHTS_URL))),
		chromedp.ActionFunc(scenarios.LogScenario(scenarios.NewAcceptGFlightCookies())),
		chromedp.WaitVisible("body", chromedp.ByQuery),
	)
	if err != nil {
		logger.ErrorLogger.Println("Error while performing the automation logic:", err)
		os.Exit(1)
	}

	query := types.Query{
		Weekdays:     map[types.Weekday]bool{types.Weekday(0): true, types.Weekday(2): true},
		StayDuration: 3,
		MonthHorizon: 2,
		Departure:    "Paris",
		Destination:  "Rome",
	}

	handleQuery := scenarios.HandleQuery{Query: query}

	err = chromedp.Run(ctx, chromedp.ActionFunc(scenarios.LogScenario(&handleQuery)))
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		os.Exit(1)
	}
}
