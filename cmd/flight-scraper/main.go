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
	)
	if err != nil {
		logger.ErrorLogger.Println("Error while performing the automation logic:", err)
		os.Exit(1)
	}

	var dTime = time.Now().Add(time.Hour * 24)
	var rTime = time.Now().Add(time.Hour * 24 * 4)
	var request types.Request = types.Request{
		Departure:     "Paris",
		Destination:   "Rome",
		DepartureDate: dTime,
		ReturnDate:    rTime,
	}
	err = chromedp.Run(ctx,
		chromedp.WaitVisible("body", chromedp.ByQuery),
	)

	fetchListOfFlights := scenarios.NewFetchListOfFlights(types.SEL_OUTBOUND_FLIGHTS_UL)
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(scenarios.LogScenario(scenarios.NewFillAndConfirmTripInfos(request))),
		chromedp.ActionFunc(scenarios.LogScenario(&fetchListOfFlights)),
	)

	ulNodeID := fetchListOfFlights.NodeID

	parseFlightCombos := scenarios.NewParseFlightCombos(ulNodeID, request.DepartureDate, request.ReturnDate)

	err = chromedp.Run(ctx, chromedp.ActionFunc(scenarios.LogScenario(&parseFlightCombos)))
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		os.Exit(1)
	}
}
