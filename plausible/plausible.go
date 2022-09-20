package plausible

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	HostAPIBase *url.URL
	SiteID      string
	Token       string
}

func (clt *Client) GetTimeseriesData() (*TimeseriesData, error) {
	url := clt.HostAPIBase.JoinPath("/api/v1/stats/aggregate")
	q := url.Query()
	q.Add("site_id", clt.SiteID)
	// TODO: Do we want to be able to configure this?
	q.Add("period", "day")
	q.Add("date", time.Now().UTC().Format("2006-01-02"))
	q.Add("metrics", "visitors,pageviews,bounce_rate,visit_duration")
	url.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", clt.Token))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("plausible: request error: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return nil, fmt.Errorf("plausible: unexpected HTTP status code %d received", response.StatusCode)
	}

	var tsData tsDTO
	err = json.NewDecoder(response.Body).Decode(&tsData)
	if err != nil {
		return nil, fmt.Errorf("plausible: failed to decode response data: %w", err)
	}
	return tsData.ToTimeseriesData(), nil
}
