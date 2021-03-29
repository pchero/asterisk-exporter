package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type Collector interface {
	Run()
}

type collector struct {
	MetricInterval int
}

func NewCollector(interval int) Collector {
	h := &collector{
		MetricInterval: interval,
	}

	return h
}

var (
	metricsNamespace = "asterisk"

	// promAsteriskHealth = prometheus.NewCounterVec(
	// 	prometheus.CounterOpts{
	// 		Namespace: metricsNamespace,
	// 		Name:      "health",
	// 		Help:      "Asterisk health check count",
	// 	},
	// 	[]string{},
	// )

	promAsteriskHealthFail = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "health_fail",
			Help:      "Asterisk health check count",
		},
		[]string{},
	)

	promCurrentChannelTech = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "crruent_channel_tech",
			Help:      "Current number of channels(tech) in the asterisk.",
		},
		[]string{"tech"},
	)

	promCurrentChannelContext = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "crruent_channel_context",
			Help:      "Current number of channels(context) in the asterisk.",
		},
		[]string{"context"},
	)

	promChannelDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "channel_duration",
			Help:      "A duration time of the channel",
			Buckets: []float64{
				5, 10, 30, 60, 120, 300, 600, 1800, 3600,
			},
		},
		[]string{"tech", "context"},
	)
)

func init() {
	prometheus.MustRegister(
		promAsteriskHealthFail,
		promCurrentChannelTech,
		promCurrentChannelContext,
		promChannelDuration,
	)
}

func (h *collector) Run() {

	for {
		if err := h.Collect(); err != nil {
			logrus.Errorf("Could not collect metrics correctly. err: %v", err)
			promAsteriskHealthFail.WithLabelValues().Inc()
		}

		// sleep
		time.Sleep(time.Second * time.Duration(h.MetricInterval))
	}
}
