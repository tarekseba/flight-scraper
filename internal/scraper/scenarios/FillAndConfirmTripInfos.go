package scenarios

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/tarekseba/flight-scraper/internal/scraper/types"
)

type FillAndConfirmTripInfos struct {
	types.Request
}

func (s *FillAndConfirmTripInfos) Name() string {
	return "FillAndConfirmTripInfos"
}

func (s *FillAndConfirmTripInfos) Do(ctx context.Context) error {
	// Fill 'FROM' input
	err := fillFromInput(ctx, s.Departure, types.WHERE_FROM_INPUT, types.D_CITY_LABEL)
	if err != nil {
		return err
	}

	// Fill 'TO' input
	err = fillToInput(ctx, s.Destination, types.WHERE_TO_INPUT, types.A_CITY_LABEL)
	if err != nil {
		return err
	}

	err = fillDateInputsAndSearch(ctx, s.DepartureDate, s.ReturnDate)
	if err != nil {
		return err
	}
	return nil
}

func NewFillAndConfirmTripInfos(request types.Request) *FillAndConfirmTripInfos {
	return &FillAndConfirmTripInfos{Request: request}
}

func fillFromInput(ctx context.Context, city string, selector string, optionSelector string) error {
	err := chromedp.WaitVisible(selector, chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Evaluate(generateTriggerInputFieldScript(selector, city), nil).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 2).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.WaitVisible(makeOptionSelector(optionSelector), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	chromedp.Click(makeOptionSelector(optionSelector), chromedp.ByQuery)
	return nil
}

func fillToInput(ctx context.Context, city string, selector string, optionSelector string) error {
	err := chromedp.Evaluate(generateTriggerInputFieldScript(types.WHERE_TO_INPUT, types.A_CITY), nil).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.WaitVisible(makeOptionSelector(types.A_CITY_LABEL), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Click(makeOptionSelector(types.A_CITY_LABEL), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func fillDateInputsAndSearch(ctx context.Context, departure time.Time, returnD time.Time) error {
	var departureDate = formatDate(departure)
	var returnDate = formatDate(returnD)
	err := chromedp.Click(types.D_DATE_INPUT, chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 2).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.WaitVisible(generateCalendarDateSelector(departureDate), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Click(generateCalendarDateSelector(departureDate), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 2).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Click(generateCalendarDateSelector(returnDate), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 3).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Evaluate(types.SEARCH_BUTTON_SCRIPT, nil).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 2).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Evaluate(types.SEARCH_BUTTON_SCRIPT, nil).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.WaitVisible("body", chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	// chromedp.OuterHTML("body", &res6, chromedp.ByQuery).Do(c
	return nil
}

func makeOptionSelector(city string) string {
	return fmt.Sprintf(" ul > li[aria-label='%s']", city)
}

func generateTriggerInputFieldScript(selector string, text string) string {
	var format = "var input = document.querySelector(\"%s\");\n" +
		"input.value = '%s';\n" +
		"input.dispatchEvent(new Event('input', { bubbles: true }));"

	return fmt.Sprintf(format, selector, text)
}

func generateCalendarDateSelector(isoDate string) string {
	return fmt.Sprintf("div[data-iso='%s']", isoDate)
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}
