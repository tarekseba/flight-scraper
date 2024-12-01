package scenarios

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/tarekseba/flight-scraper/internal/scraper/utils"
)

type GoBack struct {
	steps uint
}

func (s *GoBack) Name() string {
	return fmt.Sprintf("Go back %d steps", s.steps)
}

func NewGoBack(steps uint) GoBack {
	var s uint = 1
	if steps > 1 {
		s = steps
	}
	return GoBack{
		steps: s,
	}
}

func (s *GoBack) Do(ctx context.Context) error {
	for _ = range s.steps {
		if err := chromedp.NavigateBack().Do(ctx); err != nil {
			return utils.AnnotateError(err)
		}
	}
	err := chromedp.Sleep(time.Second * 2).Do(ctx)
	if err != nil {
		return utils.AnnotateError(err)
	}
	return nil
}
