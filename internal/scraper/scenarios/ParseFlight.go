package scenarios

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/runtime"
	"github.com/tarekseba/flight-scraper/internal/scraper/types"
	"github.com/tarekseba/flight-scraper/internal/scraper/utils"
)

type ParseFlight struct {
	NodeID    cdp.NodeID
	Flight    *types.Flight
	WithPrice bool
}

func (s *ParseFlight) Name() string {
	return "ParseFlight"
}

func (s *ParseFlight) Do(ctx context.Context) error {
	departure_time, err := getTime(ctx, s.NodeID, "Departure time")
	if err != nil {
		return utils.AnnotateError(err)
	}
	s.Flight.DepTime = departure_time

	arrival_time, err := getTime(ctx, s.NodeID, "Arrival time")
	if err != nil {
		return utils.AnnotateError(err)
	}
	s.Flight.ArrTime = arrival_time

	airline, err := getAirlineCompany(ctx, s.NodeID)
	if err != nil {
		return utils.AnnotateError(err)
	}
	airline = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`\"`, r) {
			return -1
		}
		return r
	}, airline)
	s.Flight.Company = airline

	flight_duration, err := getFlightDuration(ctx, s.NodeID)
	if err != nil {
		return utils.AnnotateError(err)
	}
	s.Flight.Duration = flight_duration

	dep_airport, arr_airport, err := getAirports(ctx, s.NodeID)
	if err != nil {
		return utils.AnnotateError(err)
	}
	s.Flight.Airports = dep_airport + " - " + arr_airport

	stops, err := getStops(ctx, s.NodeID)
	if err != nil {
		return utils.AnnotateError(err)
	}
	s.Flight.Stops = stops

	if s.WithPrice {
		price, currency, err := getPrice(ctx, s.NodeID)
		if err != nil {
			return utils.AnnotateError(err)
		}
		s.Flight.Price = price
		s.Flight.Currency = currency
	}

	return nil
}

func getTime(ctx context.Context, nodeID cdp.NodeID, selector string) (time.Time, error) {
	t := time.Now()
	spanId, err := dom.QuerySelector(nodeID, fmt.Sprintf("div > div > div span[aria-label^='%s'] > span", selector)).Do(ctx)
	if err != nil {
		return t, utils.AnnotateError(err)
	}
	attrs, err := dom.GetAttributes(spanId).Do(ctx)
	if err != nil {
		return t, utils.AnnotateError(err)
	}
	var attrValue = ""
	for i := 0; i < len(attrs)-2; i += 2 {
		if attrs[i] == "aria-label" {
			attrValue = attrs[i+1]
			break
		}
	}
	t, err = parseTime(attrValue)
	if err != nil {
		return t, utils.AnnotateError(err)
	}

	return t, nil
}

func parseTime(value string) (time.Time, error) {
	if value != "" {
		timeTextSplit := strings.Split(value, " ")
		timeText := timeTextSplit[2]
		timeText = strings.ReplaceAll(timeText, "\xe2\x80\xaf", " ")
		timeText = strings.ReplaceAll(timeText, ".", "")
		timeText = strings.Trim(timeText, " ")
		timeText = strings.ReplaceAll(timeText, ".", "")
		var timeFormat = types.TIME_FORMAT_FULL
		if len(timeText) == 7 {
			timeFormat = types.TIME_FORMAT_PARTIAL
		}
		time, err := time.Parse(timeFormat, timeText)
		if err != nil {
			return time, utils.AnnotateError(err)
		}
		return time, nil
	}
	return time.Now(), nil
}

func getAirlineCompany(ctx context.Context, nodeID cdp.NodeID) (string, error) {
	var airline = ""
	node, err := dom.QuerySelector(nodeID, "div > div > div > div > div > div > span:not([aria-label])").Do(ctx)
	if err != nil {
		return airline, utils.AnnotateError(err)
	}
	object, err := dom.ResolveNode().WithNodeID(node).Do(ctx)
	if err != nil {
		return airline, utils.AnnotateError(err)
	}
	res, _, err := runtime.CallFunctionOn(types.INNER_TEXT_FUNC).WithObjectID(object.ObjectID).Do(ctx)
	if err != nil {
		return airline, utils.AnnotateError(err)
	}
	airline = string(res.Value)

	return airline, nil
}

func getFlightDuration(ctx context.Context, nodeID cdp.NodeID) (string, error) {
	div, err := dom.QuerySelector(nodeID, "div > div > div div[aria-label^='Total duration']").Do(ctx)
	if err != nil {
		return "", utils.AnnotateError(err)
	}
	object, err := dom.ResolveNode().WithNodeID(div).Do(ctx)
	if err != nil {
		return "", utils.AnnotateError(err)
	}
	val, _, err := runtime.CallFunctionOn(types.INNER_TEXT_FUNC).WithObjectID(object.ObjectID).Do(ctx)
	if err != nil {
		return "", utils.AnnotateError(err)
	}
	duration := string(val.Value)
	if duration == "" {
		return "", utils.AnnotateError(errors.New("Duration was empty"))
	}
	return strings.ReplaceAll(duration, "\"", ""), utils.AnnotateError(err)
}

func getAirports(ctx context.Context, nodeID cdp.NodeID) (string, string, error) {
	elements, err := dom.QuerySelectorAll(nodeID, "div > div > div > div > div > div > span span[aria-label='']").Do(ctx)
	if err != nil {
		return "", "", utils.AnnotateError(err)
	}
	if len(elements) <= 0 {
		return "", "", utils.AnnotateError(errors.New("Number of airport names spans is nil"))
	}
	airports := [4]string{"", ""}
	for index := range elements {
		node := elements[index]
		obj, err := dom.ResolveNode().WithNodeID(node).Do(ctx)
		if err != nil {
			return "", "", utils.AnnotateError(err)
		}
		val, _, err := runtime.CallFunctionOn(types.INNER_TEXT_FUNC).WithObjectID(obj.ObjectID).Do(ctx)
		temp := string(val.Value)
		if temp == "" {
			return "", "", utils.AnnotateError(errors.New("Aeroport span text is empty"))
		}
		airports[index] = temp
	}
	return strings.ReplaceAll(airports[0], "\"", ""), strings.ReplaceAll(airports[1], "\"", ""), nil
}

func getStops(ctx context.Context, nodeID cdp.NodeID) (uint, error) {
	node, err := dom.QuerySelector(nodeID, "div > div > div span[aria-label*='stop']").Do(ctx)
	if err != nil {
		return 0, utils.AnnotateError(err)
	}
	obj, err := dom.ResolveNode().WithNodeID(node).Do(ctx)
	if err != nil {
		return 0, utils.AnnotateError(err)
	}
	val, _, err := runtime.CallFunctionOn(types.INNER_TEXT_FUNC).WithObjectID(obj.ObjectID).Do(ctx)
	if err != nil {
		return 0, utils.AnnotateError(err)
	}
	stops := string(val.Value)
	if stops == "" {
		return 0, utils.AnnotateError(errors.New("Stops span is empty"))
	}
	stops = strings.ReplaceAll(stops, "\"", "")
	if strings.HasPrefix(stops, "Non") {
		return 0, nil
	}
	regex := regexp.MustCompile(`^(\d+).*`)
	matches := regex.FindStringSubmatch(stops)
	if matches == nil {
		return 0, utils.AnnotateError(err)
	}
	stop_count, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, utils.AnnotateError(err)
	}
	return uint(stop_count), nil
}

func getPrice(ctx context.Context, nodeID cdp.NodeID) (int, string, error) {
	price_node, err := dom.QuerySelector(nodeID, "div > div > div > div > div:not([role]) > div > div > div > span[aria-label]").Do(ctx)
	if err != nil {
		return -1, "", utils.AnnotateError(err)
	}
	obj, err := dom.ResolveNode().WithNodeID(price_node).Do(ctx)
	if err != nil {
		return -1, "", utils.AnnotateError(err)
	}
	val, _, err := runtime.CallFunctionOn("function() {return this.getAttribute(\"aria-label\")}").WithObjectID(obj.ObjectID).Do(ctx)
	if err != nil {
		return -1, "", utils.AnnotateError(err)
	}
	string_price := string(val.Value)
	if string_price == "" {
		return -1, "", utils.AnnotateError(errors.New("Price is empty"))
	}

	string_price = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`\"`, r) {
			return -1
		}
		return r
	}, string_price)

	string_price_parts := strings.Split(string_price, " ")
	price, err := strconv.Atoi(strings.Trim(string_price_parts[0], " "))
	if err != nil {
		return -1, "", utils.AnnotateError(err)
	}

	return price, strings.Trim(string_price_parts[1], " "), nil
}
