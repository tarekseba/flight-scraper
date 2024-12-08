package scenarios

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/tarekseba/flight-scraper/internal/scraper/types"
)

type AcceptGFlightCookies struct {
}

func (s *AcceptGFlightCookies) Name() string {
	return "AcceptGFlightCookies"
}

func (s *AcceptGFlightCookies) Do(ctx context.Context) error {
	err := chromedp.WaitVisible(types.SELECTOR_ACCEPT_COOKIES_BTN, chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Click(types.SELECTOR_ACCEPT_COOKIES_BTN, chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(1 * time.Second).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewAcceptGFlightCookies() *AcceptGFlightCookies {
	return &AcceptGFlightCookies{}
}
