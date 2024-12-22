package scenarios

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/tarekseba/flight-scraper/internal/logger"
	"github.com/tarekseba/flight-scraper/internal/scraper/types"
	"github.com/tarekseba/flight-scraper/internal/scraper/utils"
)

type HandleQuery struct {
	Query types.Query
}

func (s *HandleQuery) Name() string {
	return fmt.Sprintf("HandleQuery %+v", s.Query)
}

func (s *HandleQuery) Do(allocCtx context.Context) error {
	requests := s.Query.IntoRequests()

	results := make([]types.RequestResult, 0, 10)
	for idx := range requests {
		// new context per request
		ctx, cancel := context.WithCancel(allocCtx)
		timeoutCtx, _ := context.WithTimeout(ctx, time.Second*40)
		ctx, _ = chromedp.NewContext(
			timeoutCtx,
		)
		err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			// Navigate to page
			LogScenario(NewNavigateToPage(types.G_FLIGHTS_URL))(ctx)
			LogScenario(NewAcceptGFlightCookies())(ctx)
			chromedp.WaitVisible("body", chromedp.ByQuery).Do(ctx)

			err := LogScenario(NewFillAndConfirmTripInfos(requests[idx]))(ctx)
			if err != nil {
				return utils.AnnotateError(err)
			}

			fetchListOfFlights := NewFetchListOfFlights(types.SEL_OUTBOUND_FLIGHTS_UL)
			err = LogScenario(&fetchListOfFlights)(ctx)
			if err != nil {
				return utils.AnnotateError(err)
			}

			ulNodeID := fetchListOfFlights.NodeID
			parseFlightCombos := NewParseFlightCombos(ulNodeID, requests[idx])

			err = LogScenario(&parseFlightCombos)(ctx)
			if err != nil {
				return utils.AnnotateError(err)
			}

			results = append(results, parseFlightCombos.RequestRes)
			return nil
		}))
		if err != nil {
			logger.ErrorLogger.Println(utils.AnnotateError(err))
		}
		cancel()
	}

	return nil
}
