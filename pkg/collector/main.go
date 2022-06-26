package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Collector interface
type Collector interface {
	Run()
}

type collector struct {
	MetricInterval int
}

// NewCollector returns a new Collector
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

	// bridge current numbers
	promCurrentBridgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "crruent_bridge_count",
			Help:      "Current number of bridges in the asterisk.",
		},
		[]string{"type", "tech"},
	)

	// bridge duration
	promBridgeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "bridge_duration",
			Help:      "A duration time of the bridge",
			Buckets: []float64{
				5, 10, 30, 60, 120, 300, 600, 1800, 3600, 7200, 14400,
			},
		},
		[]string{"type", "tech"},
	)
)

func init() {
	prometheus.MustRegister(
		promAsteriskHealthFail,
		promCurrentChannelTech,
		promCurrentChannelContext,
		promChannelDuration,

		promCurrentBridgeCount,
		promBridgeDuration,
	)
}

func (h *collector) Run() {
	logrus.Debugf("Running the collect.")

	for {
		if err := h.Collect(); err != nil {
			logrus.Errorf("Could not collect metrics correctly. err: %v", err)
			promAsteriskHealthFail.WithLabelValues().Inc()
		}

		// sleep
		time.Sleep(time.Second * time.Duration(h.MetricInterval))
	}
}

// Collect collects the asterisk's metrics and update the prometheus metric
func (h *collector) Collect() error {
	log := logrus.WithFields(logrus.Fields{
		"func": "Collect",
	})

	if err := h.channelCollects(); err != nil {
		log.Errorf("Could not get channel metrics. err: %v", err)
		return err
	}

	if err := h.bridgeCollects(); err != nil {
		log.Errorf("Could not get bridge metrics. err: %v", err)
		return err
	}

	return nil
}
