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
