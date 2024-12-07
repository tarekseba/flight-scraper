package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/tarekseba/flight-scraper/internal/scraper/utils"
)

func initCityAriaLabelMap() map[string]string {
	return map[string]string{
		"Paris":     "Paris",
		"Rome":      "Rome, Italy",
		"Milan":     "Milan, Italy",
		"Barcelona": "Barcelona, Spain",
		"Madrid":    "Madrid, Spain",
	}
}

var cityAriaLabelMap utils.Once[map[string]string] = utils.NewOnce(initCityAriaLabelMap)

func CityAriaLabelMap(city string) (string, error) {
	val := cityAriaLabelMap.Compute()[city]
	if val != "" {
		return val, nil
	}
	return "", errors.New(fmt.Sprintf("City '%s' not found in AriaLabelMap"))
}

const (
	G_FLIGHTS_URL        = "https://www.google.com/travel/flights"
	ACCEPT_COOKIES_BTN   = "button[aria-label='Accept all']"
	WHERE_FROM_INPUT     = "div[data-placeholder='Where from?'] input"
	WHERE_TO_INPUT       = "div[data-placeholder='Where to?'] input"
	D_DATE_INPUT         = "div[data-enable-prices] input[placeholder='Departure']"
	R_DATE_INPUT         = "div[data-enable-prices] input[placeholder='Return']"
	CONFIRM_DATE_BTN     = "button[aria-label='Done. Search for round trip flights, departing on December 3, 2024 and returning on December 7, 2024']"
	D_CITY               = "Paris"
	D_CITY_LABEL         = "Paris"
	A_CITY               = "Rome"
	A_CITY_LABEL         = "Rome, Italy"
	SEARCH_BUTTON_SCRIPT = `var btn = document.querySelector("button[aria-label='Search']");
			if (btn == null) {
				btn = document.querySelector("button[aria-label='Search for flights']")
			}
			btn.click();`
)

const (
	SEL_OUTBOUND_FLIGHTS_UL = "div[role='tabpanel'] > div > ul"
	SEL_RETURN_FLIGHTS_UL   = "div:has(> div > h3) > div > ul"
)

const (
	TIME_FORMAT_FULL    = "03:04 PM"
	TIME_FORMAT_PARTIAL = "3:04 PM"
	DATE_FORMAT         = "2006-01-02"
)

const (
	INNER_TEXT_FUNC = "function() {return this.innerText}"
)

const (
	DAY = time.Hour * 24
)

const (
	SELECTOR_AIRPORT         = "div > div > div > div > div > div > span span[aria-label='']:not(:has(> span))"
	SELECTOR_STOPS           = "div > div > div span[aria-label*='stop']"
	SELECTOR_PRICE           = "div > div > div > div > div:not([role]) > div > div > div > span[aria-label]"
	SELECTOR_FLIGHT_DURATION = "div > div > div div[aria-label^='Total duration']"
	SELECTOR_AIRLINE         = "div > div > div > div > div > div > span:not([aria-label])"
	SELECTOR_FLIGHT_TIME     = "div > div > div span[aria-label^='%s'] > span"
)
