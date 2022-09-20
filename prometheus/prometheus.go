package prometheus

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/riesinger/plausible-exporter/plausible"
)

type MetricsServer struct {
	pageviews     prometheus.Gauge
	visitors      prometheus.Gauge
	bounceRate    prometheus.Gauge
	visitDuration prometheus.Gauge
}

func NewServer(siteID string) *MetricsServer {
	sid := strings.ReplaceAll(siteID, "-", "_")
	sid = strings.ReplaceAll(sid, ".", "_")

	pageviews := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "plausible",
		Subsystem: sid,
		Name:      "pageviews",
		Help:      fmt.Sprintf("Number of page views for %s", siteID),
	})

	visitors := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "plausible",
		Subsystem: sid,
		Name:      "visitors",
		Help:      fmt.Sprintf("Number of visitors for %s", siteID),
	})

	bounceRate := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "plausible",
		Subsystem: sid,
		Name:      "bounce_rate",
		Help:      fmt.Sprintf("Bounce rate for %s in %%", siteID),
	})

	visitDuration := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "plausible",
		Subsystem: sid,
		Name:      "visit_duration",
		Help:      fmt.Sprintf("Average visit duration for %s in seconds", siteID),
	})

	return &MetricsServer{
		pageviews:     pageviews,
		visitors:      visitors,
		bounceRate:    bounceRate,
		visitDuration: visitDuration,
	}
}

func (srv *MetricsServer) UpdateData(data *plausible.TimeseriesData) {
	srv.pageviews.Set(float64(data.Pageviews))
	srv.visitors.Set(float64(data.Visitors))
	srv.bounceRate.Set(float64(data.BounceRate))
	srv.visitDuration.Set(float64(data.VisitDuration))
}
