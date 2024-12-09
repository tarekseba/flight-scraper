package scenarios

import (
	"context"
	"fmt"

	"github.com/tarekseba/flight-scraper/internal/scraper/types"
	"github.com/tarekseba/flight-scraper/internal/scraper/utils"
)

type HandleQuery struct {
	Query types.Query
}

func (s *HandleQuery) Name() string {
	return fmt.Sprintf("HandleQuery %+v", s.Query)
}

func (s *HandleQuery) Do(ctx context.Context) error {
	requests := s.Query.IntoRequests()

	for idx := range requests {
		err := LogScenario(NewFillAndConfirmTripInfos(requests[idx]))(ctx)
		if err != nil {
			return utils.AnnotateError(err)
		}

		fetchListOfFlights := NewFetchListOfFlights(types.SEL_OUTBOUND_FLIGHTS_UL)
		err = LogScenario(&fetchListOfFlights)(ctx)
		if err != nil {
			return err
		}

		ulNodeID := fetchListOfFlights.NodeID
		parseFlightCombos := NewParseFlightCombos(ulNodeID, requests[idx].DepartureDate, requests[idx].ReturnDate)

		err = LogScenario(&parseFlightCombos)(ctx)
		if err != nil {
			return err
		}
		err = LogScenario(NewNavigateToPage(types.G_FLIGHTS_URL))(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
