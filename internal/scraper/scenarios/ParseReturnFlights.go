package scenarios

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/tarekseba/flight-scraper/internal/scraper/types"
	"github.com/tarekseba/flight-scraper/internal/scraper/utils"
)

type ParseReturnFlights struct {
	LiNodeID   cdp.NodeID
	returnDate time.Time
	flights    []types.Flight
}

func (s *ParseReturnFlights) Name() string {
	return "ParseReturnFlights"
}

func (s *ParseReturnFlights) Do(ctx context.Context) error {
	err := navigate_to_return_flights(ctx, s.LiNodeID)
	if err != nil {
		return err
	}

	nodeID, err := fetch_list_of_flights_ul(ctx, types.SEL_RETURN_FLIGHTS_UL)
	if err != nil {
		return err
	}

	liNodes, err := dom.QuerySelectorAll(nodeID, "li").Do(ctx)
	if err != nil {
		return err
	}

	for idx := range liNodes {
		var return_flight = types.Flight{
			DepDate: s.returnDate,
		}
		var parseFlight = ParseFlight{NodeID: liNodes[idx], Flight: &return_flight, WithPrice: true}
		err = parseFlight.Do(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", return_flight)
		s.flights = append(s.flights, return_flight)
	}
	go_back := NewGoBack(1)
	err = go_back.Do(ctx)
	if err != nil {
		return utils.AnnotateError(err)
	}

	return err
}

func navigate_to_return_flights(ctx context.Context, nodeID cdp.NodeID) error {
	clickable_node, err := dom.QuerySelector(nodeID, "div > div").Do(ctx)
	if err != nil {
		return utils.AnnotateError(err)
	}
	obj, err := dom.ResolveNode().WithNodeID(clickable_node).Do(ctx)
	if err != nil {
		return utils.AnnotateError(err)
	}
	_, _, err = runtime.CallFunctionOn("function() { this.click() }").WithObjectID(obj.ObjectID).Do(ctx)
	if err != nil {
		return utils.AnnotateError(err)
	}
	chromedp.Sleep(time.Second * 2)

	return nil
}

func fetch_list_of_flights_ul(ctx context.Context, sel string) (cdp.NodeID, error) {
	fetchListOfFlights := NewFetchListOfFlights(sel)
	err := fetchListOfFlights.Do(ctx)
	if err != nil {
		return 0, utils.AnnotateError(err)
	}
	return fetchListOfFlights.NodeID, nil
}
