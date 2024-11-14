package types

import (
	"context"
	"errors"
	"fmt"
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

type Weekday time.Weekday

func (day *Weekday) Name() string {
	switch int(*day) {
	case 1:
		return SUNDAY
	case 2:
		return MONDAY
	case 3:
		return TUESDAY
	case 4:
		return WEDNESDAY
	case 5:
		return THURSDAY
	case 6:
		return FRIDAY
	case 7:
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
	Weekdays     map[Weekday]bool
	StayDuration int
	MonthHorizon int
}

func (q *Query) IntoRequests() []Request {
	var requests = make([]Request, 0)
	var currentMonth = time.Now().Month()
	var maxMonth = MaxMonth(currentMonth, q.MonthHorizon)
	fmt.Println(maxMonth)
	return requests
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
