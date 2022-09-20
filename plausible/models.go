package plausible

import (
	"strings"
	"time"
)

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) // get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}
	*d = Date(t) // set result using the pointer
	return nil
}

type TimeseriesData struct {
	Pageviews     uint
	BounceRate    float32
	VisitDuration float32
	Visitors      uint
}

type tsDTO struct {
	Results struct {
		Pageviews struct {
			Value uint
		}
		BounceRate struct {
			Value float32
		} `json:"bounce_rate"`
		VisitDuration struct {
			Value float32
		} `json:"visit_duration"`
		Visitors struct {
			Value uint
		}
	}
}

func (d tsDTO) ToTimeseriesData() *TimeseriesData {
	return &TimeseriesData{
		Pageviews:     d.Results.Pageviews.Value,
		BounceRate:    d.Results.BounceRate.Value,
		VisitDuration: d.Results.VisitDuration.Value,
		Visitors:      d.Results.Visitors.Value,
	}
}
