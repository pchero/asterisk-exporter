package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	joonix "github.com/joonix/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"gitlab.com/voipbin/voip/asterisk-exporter.git/pkg/collector"
)

// args
var webListenAddress = flag.String("web_listen_address", ":9495", "Address to listen on for web interface and telemetry.")
var webListenPath = flag.String("web_listen_path", "/metrics", "Path under which to expose metrics.")

var asteriskMetricInterval = flag.Int("asterisk_metric_interval", 5, "Interval sec for metric getting")

// signal channels
var (
	chSigs = make(chan os.Signal, 1)
	chDone = make(chan bool, 1)
)

func main() {
	c := collector.NewCollector(*asteriskMetricInterval)
	go c.Run()

	<-chDone
}

func init() {
	// arg parse
	flag.Parse()

	// init log
	initLog()

	// init signal
	initSignal()

	// init prometheus
	initProm(*webListenPath, *webListenAddress)
}

// initLog inits log settings.
func initLog() {
	logrus.SetFormatter(joonix.NewFormatter())
	logrus.SetLevel(logrus.DebugLevel)
}

// initSignal inits sinal settings.
func initSignal() {
	signal.Notify(chSigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	go signalHandler()
}

// signalHandler catches signals and set the done
func signalHandler() {
	sig := <-chSigs
	logrus.Debugf("Received signal. sig: %v", sig)
	chDone <- true
}

// initProm inits prometheus settings
func initProm(endpoint, listen string) {
	http.Handle(endpoint, promhttp.Handler())
	go func() {
		for {
			err := http.ListenAndServe(listen, nil)
			if err != nil {
				logrus.Errorf("Could not start prometheus listener")
				time.Sleep(time.Second * 1)
				continue
			}
			break
		}
	}()
}
