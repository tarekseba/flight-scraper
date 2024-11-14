package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/tarekseba/flight-scraper/internal/logger"
	"github.com/tarekseba/flight-scraper/internal/scraper/scenarios"
)
const DAY = time.Hour * 24

const (
	G_FLIGHTS_URL        = "https://www.google.com/travel/flights"
	ACCEPT_COOKIES_BTN   = "button[aria-label='Accept all']"
	WHERE_FROM_INPUT     = "div[data-placeholder='Where from?'] input"
	WHERE_TO_INPUT       = "div[data-placeholder='Where to?'] input"
	D_DATE_INPUT         = "div[data-enable-prices] input[placeholder='Departure']"
	R_DATE_INPUT         = "div[data-enable-prices] input[placeholder='Return']"
	CONFIRM_DATE_BTN     = "button[aria-label='Done. Search for round trip flights, departing on December 3, 2024 and returning on December 7, 2024']"
	D_CITY               = "Paris"
	D_CITY_LABEL         = "Paris"
	A_CITY               = "Rome"
	A_CITY_LABEL         = "Rome, Italy"
	SEARCH_BUTTON_SCRIPT = `var btn = document.querySelector("button[aria-label='Search']");
			if (btn == null) {
				btn = document.querySelector("button[aria-label='Search for flights']")
			}
			btn.click();`
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
	logger.InfoLogger.Println("GET " + G_FLIGHTS_URL)
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(scenarios.LogScenario(scenarios.NewNavigateToPage(G_FLIGHTS_URL))),
		chromedp.ActionFunc(scenarios.LogScenario(scenarios.NewAcceptGFlightCookies())),
	)
	if err != nil {
		logger.ErrorLogger.Println("Error while performing the automation logic:", err)
		os.Exit(1)
	}

