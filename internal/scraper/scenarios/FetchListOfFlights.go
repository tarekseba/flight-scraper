package scenarios

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/tarekseba/flight-scraper/internal/scraper/utils"
)

type FetchListOfFlights struct {
	NodeID   cdp.NodeID
	selector string
}

func (s *FetchListOfFlights) Name() string {
	return "FetchListOfFlights"
}

func NewFetchListOfFlights(sel string) FetchListOfFlights {
	return FetchListOfFlights{NodeID: -1, selector: sel}
}

func (s *FetchListOfFlights) Do(ctx context.Context) error {
	var rootNodeID cdp.NodeID = -1
	var ulNodeID cdp.NodeID = -1
	err := chromedp.Sleep(time.Second * 1).Do(ctx)
	if err != nil {
		return utils.AnnotateError(err)
	}
	node, err := dom.GetDocument().Do(ctx)
	if err != nil {
		return utils.AnnotateError(err)
	}
	rootNodeID = node.NodeID

	nodeId, err := dom.QuerySelector(rootNodeID, s.selector).Do(ctx)
	if err != nil || nodeId == 0 {
		return utils.AnnotateError(err)
	}

	ulNodeID = nodeId
	s.NodeID = ulNodeID
	return nil
}
