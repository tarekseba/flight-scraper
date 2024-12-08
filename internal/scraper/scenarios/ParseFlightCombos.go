package scenarios

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"

	"github.com/tarekseba/flight-scraper/internal/scraper/types"
	"github.com/tarekseba/flight-scraper/internal/scraper/utils"
)

type ParseFlightCombos struct {
	UlNodeID cdp.NodeID
	DepDate  time.Time
	RetDate  time.Time
}

func NewParseFlightCombos(ulNodeID cdp.NodeID, depDate, retDate time.Time) ParseFlightCombos {
	return ParseFlightCombos{
		UlNodeID: ulNodeID,
		DepDate:  depDate,
		RetDate:  retDate,
	}
}

func (s *ParseFlightCombos) Name() string {
	return "ParseFlightCombos"
}

func (s *ParseFlightCombos) Do(ctx context.Context) error {
	nodes, err := dom.QuerySelectorAll(s.UlNodeID, "li").Do(ctx)
	if err != nil {
		return utils.AnnotateError(err)
	}
	var ulChildren = nodes
	length := len(ulChildren)
	inbound_outbound_flights_map := make(map[types.Flight][]types.Flight)
	for i := range length {
		if i >= len(ulChildren) {
			return nil
		}
		var flight = types.Flight{
			DepDate: s.DepDate,
		}
		var parseFlight = ParseFlight{NodeID: ulChildren[i], Flight: &flight, WithPrice: true}
		err = LogScenario(&parseFlight)(ctx)
		if err != nil {
			return err
		}
		// fmt.Printf("%+v\n", flight)
		arr := insert_key_into_map(inbound_outbound_flights_map, flight)
		var parseReturnFlights = ParseReturnFlights{
			LiNodeID:   ulChildren[i],
			returnDate: s.RetDate,
			flights:    arr,
		}
		err := LogScenario(&parseReturnFlights)(ctx)
		if err != nil {
			return err
		}
		inbound_outbound_flights_map[flight] = parseReturnFlights.flights
		node, err := fetch_list_of_flights_ul(ctx, types.SEL_OUTBOUND_FLIGHTS_UL)
		if err != nil {
			return err
		}
		ulChildren, err = dom.QuerySelectorAll(node, "li").Do(ctx)
		if err != nil {
			return utils.AnnotateError(err)
		}
	}
	return nil
}

func insert_key_into_map(m map[types.Flight][]types.Flight, key types.Flight) []types.Flight {
	elt := m[key]
	if elt != nil {
		return elt
	}
	arr := make([]types.Flight, 5)
	m[key] = arr
	return arr
}
