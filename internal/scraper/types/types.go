package types

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	SUNDAY    = "Sunday"
	MONDAY    = "Monday"
	TUESDAY   = "Tuesday"
	WEDNESDAY = "Wednesday"
	THURSDAY  = "Thursday"
	FRIDAY    = "Friday"
	SATURDAY  = "Saturday"
)

type HashSet[T comparable] map[T]bool

type Weekday time.Weekday

func (day *Weekday) Name() string {
	switch int(*day) {
	case 0:
		return SUNDAY
	case 1:
		return MONDAY
	case 2:
		return TUESDAY
	case 3:
		return WEDNESDAY
	case 4:
		return THURSDAY
	case 5:
		return FRIDAY
	case 6:
		return SATURDAY
	}
	return ""
}
func (day *Weekday) Short() string {
	var name = day.Name()
	if len(name) >= 3 {
		return name[:3]
	}
	return ""
}

func (day *Weekday) From(d interface{}) error {
	switch v := d.(type) {
	case string:
		{
			switch strings.ToLower(v) {
			case strings.ToLower(SUNDAY):
				*day = 1
				break
			case strings.ToLower(MONDAY):
				*day = 2
				break
			case strings.ToLower(TUESDAY):
				*day = 3
			case strings.ToLower(WEDNESDAY):
				*day = 4
				break
			case strings.ToLower(THURSDAY):
				*day = 5
				break
			case strings.ToLower(FRIDAY):
				*day = 6
				break
			case strings.ToLower(SATURDAY):
				*day = 7
				break
			default:
				return errors.New(fmt.Sprintf("Wrong weekday received %s", v))
			}
		}
	case int:
		{
			if v > 0 && v < 8 {
				*day = Weekday(v)
				break
			} else {
				return errors.New(fmt.Sprintf("Wrong weekday received %d", v))
			}

		}
	default:
		{
			return errors.New(fmt.Sprintf("Wrong weekday type received %d", v))
		}
	}
	return nil
}

type Pair struct {
	A string
	B string
}

type Request struct {
	Departure     string
	Destination   string
	DepartureDate time.Time
	ReturnDate    time.Time
}

type Query struct {
	Weekdays     HashSet[Weekday] `json:"weekdays"`
	StayDuration int              `json:"stay_duration"`
	MonthHorizon int              `json:"month_horizon"`
	Departure    string           `json:"departure"`
	Destination  string           `json:"destination"`
}

func (q *Query) SanitizeQuery() error {
	weekdays := HashSet[Weekday]{}
	for k, _ := range q.Weekdays {
		intVal := int(k)
		if intVal >= 0 && intVal <= 6 {
			weekdays[k] = true
		}
	}
	if len(weekdays) == 0 {
		return fmt.Errorf("Weekdays empty after sanitization")
	}
	q.StayDuration = int(math.Max(float64(q.StayDuration), 1))
	q.MonthHorizon = int(math.Min(float64(q.MonthHorizon), 2))
	if q.Departure == "" {
		return fmt.Errorf("Departure field is empty")
	}
	if q.Destination == "" {
		return fmt.Errorf("Destination field is empty")
	}
	return nil
}

func (q *Query) IntoRequests() []Request {
	var requests = make([]Request, 0)
	var currentMonth = time.Now().Month()
	var maxMonth = MaxMonth(currentMonth, q.MonthHorizon)
	if len(q.Weekdays) <= 0 {
		return requests
	}
	currentDate := time.Now()
	currentWeekday := Weekday(currentDate.Weekday())
	daysArray := make([]Weekday, len(q.Weekdays))
	i := 0
	for d, _ := range q.Weekdays {
		daysArray[i] = d
		i++
	}
	slices.Sort(daysArray)

	for currentMonth != maxMonth {
		for d := range daysArray {
			request := Request{Departure: q.Departure, Destination: q.Destination}
			key := daysArray[d]
			diff := key.Difference(&currentWeekday)
			currentDate = currentDate.Add(DAY * time.Duration(diff))
			currentMonth = currentDate.Month()
			if currentMonth == maxMonth {
				break
			}
			currentWeekday = Weekday(currentDate.Weekday())
			request.DepartureDate = currentDate
			returnDate := currentDate.Add(DAY * time.Duration(q.StayDuration))
			request.ReturnDate = returnDate
			requests = append(requests, request)
		}
	}
	return requests
}

func (q *Query) StringifyWeekdays() string {
	weekdays := ""
	for k, _ := range q.Weekdays {
		if k > 6 || k < 0 {
			continue
		}
		weekdays = weekdays + strconv.Itoa(int(k))
	}
	return weekdays
}

func PlusMonths(month time.Month, m int) time.Month {
	for _ = range m {
		month = (month % 12) + 1
	}
	return month
}

func MaxMonth(month time.Month, m int) time.Month {
	return PlusMonths(month, m+1)
}

func (after *Weekday) Difference(before *Weekday) int {
	var x = int(*after) - int(*before)
	if x <= 0 {
		return 7 + x
	}
	return x
}

type Trip struct {
	From          string
	To            string
	Aeroports     string
	properties    string
	DepartureDate time.Time
	ReturnDate    time.Time
	Price         int
	companyLogo   string
}

type Scenario interface {
	Name() string
	Do(ctx context.Context) error
}

type RoundTripFlight struct {
	OutboundDepTime  time.Time
	OutboundArrTime  time.Time
	OutboundDepDate  time.Time
	OutboundAirports string
	OutboundCompany  string
	OutboundDuration string
	OutboundStops    uint
	InboundDepTime   time.Time
	InboundArrTime   time.Time
	InboundDepDate   time.Time
	InboundAirports  string
	InboundCompany   string
	InboundDuration  string
	InboundStops     uint
	Price            int
	Currency         string
}

func RoundTripFlightFromFlights(outbound, inbound Flight) RoundTripFlight {
	return RoundTripFlight{
		OutboundDepTime:  outbound.DepTime,
		OutboundArrTime:  outbound.ArrTime,
		OutboundDepDate:  outbound.DepDate,
		OutboundAirports: outbound.Airports,
		OutboundCompany:  outbound.Currency,
		OutboundDuration: "",
		OutboundStops:    0,
		InboundDepTime:   inbound.DepTime,
		InboundArrTime:   inbound.ArrTime,
		InboundDepDate:   inbound.DepDate,
		InboundAirports:  inbound.Airports,
		InboundCompany:   inbound.Currency,
		InboundDuration:  "",
		InboundStops:     0,
		Price:            inbound.Price,
		Currency:         inbound.Currency,
	}
}

type Flight struct {
	DepDate  time.Time `json:"dep_date"`
	DepTime  time.Time `json:"dep_time"`
	ArrTime  time.Time `json:"arr_time"`
	Airports string    `json:"airports"`
	Company  string    `json:"company"`
	Duration string    `json:"duration"`
	Price    int       `json:"price"`
	Currency string    `json:"currency"`
	Stops    uint      `json:"stops"`
}

func (f *Flight) ID() string {
	date := f.DepDate.Format(DATE_FORMAT)
	time := strings.ReplaceAll(f.DepTime.Format(TIME_FORMAT_FULL), " ", "")
	comp := f.Company
	comp = strings.ReplaceAll(comp, " ", "")
	if len(comp) > 4 {
		comp = comp[:4]
	}
	airports := strings.ReplaceAll(f.Airports, " ", "")
	return fmt.Sprintf("%s.%s.%s", date, time, airports)
}

type RequestResult struct {
	Req    Request             `json:"req"`
	DepIds map[string]Flight   `json:"dep_ids"`
	DepRet map[string][]Flight `json:"dep_ret"`
}

func NewRequestResult(req Request, dep_ids map[string]Flight, dep_ret map[string][]Flight) RequestResult {
	return RequestResult{
		Req:    req,
		DepIds: dep_ids,
		DepRet: dep_ret,
	}
}
