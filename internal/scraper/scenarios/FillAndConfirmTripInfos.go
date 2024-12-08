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
	err := fill_from_input(ctx, s.Departure, types.SELECTOR_WHERE_FROM_INPUT, types.SELECTOR_D_CITY_LABEL)
	if err != nil {
		return err
	}

	// Fill 'TO' input
	err = fill_to_input(ctx, s.Destination, types.SELECTOR_WHERE_TO_INPUT, types.SELECTOR_A_CITY_LABEL)
	if err != nil {
		return err
	}

	err = fill_date_inputs_and_search(ctx, s.DepartureDate, s.ReturnDate)
	if err != nil {
		return err
	}
	return nil
}

func NewFillAndConfirmTripInfos(request types.Request) *FillAndConfirmTripInfos {
	return &FillAndConfirmTripInfos{Request: request}
}

func fill_from_input(ctx context.Context, city string, selector string, optionSelector string) error {
	err := chromedp.WaitVisible(selector, chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Evaluate(generate_trigger_input_field_script(selector, city), nil).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 2).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.WaitVisible(make_option_selector(optionSelector), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	chromedp.Click(make_option_selector(optionSelector), chromedp.ByQuery)
	return nil
}

func fill_to_input(ctx context.Context, city string, selector string, optionSelector string) error {
	err := chromedp.Evaluate(generate_trigger_input_field_script(types.SELECTOR_WHERE_TO_INPUT, types.SELECTOR_A_CITY), nil).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.WaitVisible(make_option_selector(types.SELECTOR_A_CITY_LABEL), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Click(make_option_selector(types.SELECTOR_A_CITY_LABEL), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func fill_date_inputs_and_search(ctx context.Context, departure time.Time, returnD time.Time) error {
	var departureDate = format_date(departure)
	var returnDate = format_date(returnD)
	err := chromedp.Click(types.SELECTOR_D_DATE_INPUT, chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 2).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.WaitVisible(generate_calendar_date_selector(departureDate), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Click(generate_calendar_date_selector(departureDate), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 2).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Click(generate_calendar_date_selector(returnDate), chromedp.ByQuery).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 3).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Evaluate(types.SELECTOR_SEARCH_BUTTON_SCRIPT, nil).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Sleep(time.Second * 2).Do(ctx)
	if err != nil {
		return err
	}
	err = chromedp.Evaluate(types.SELECTOR_SEARCH_BUTTON_SCRIPT, nil).Do(ctx)
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

func make_option_selector(city string) string {
	return fmt.Sprintf(" ul > li[aria-label='%s']", city)
}

func generate_trigger_input_field_script(selector string, text string) string {
	var format = "var input = document.querySelector(\"%s\");\n" +
		"input.value = '%s';\n" +
		"input.dispatchEvent(new Event('input', { bubbles: true }));"

	return fmt.Sprintf(format, selector, text)
}

func generate_calendar_date_selector(isoDate string) string {
	return fmt.Sprintf("div[data-iso='%s']", isoDate)
}

func format_date(t time.Time) string {
	return t.Format("2006-01-02")
}
