package scenarios

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type NavigateToPage struct {
	Url string
}

func (s *NavigateToPage) Name() string {
	return "NavigateToPageScenario"
}

func (s *NavigateToPage) Do(ctx context.Context) error {
	err := chromedp.Navigate(s.Url).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(2 * time.Second).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewNavigateToPage(url string) *NavigateToPage {
	return &NavigateToPage{
		Url: url,
	}
}
