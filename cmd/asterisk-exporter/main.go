package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "asterisk_exporter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  102400,
	WriteBufferSize: 102400,
}

var (
	astAriAddress  = flag.String("ast_ari_address", "127.0.0.1:8088", "Asterisk's HTTP open address with port number for ARI.")
	astAriUserinfo = flag.String("ast_ari_userinfo", "asterisk:asterisk", "Asterisk's ARI user info. username:password.")
	astAriAppname  = flag.String("ast_ari_appname", "asterisk-exporter", "Asterisk's ARI application name.")

	promAddress = flag.String("prom_address", ":9200", "Prometheus listen port")
)

var (
	chanAriMsg = make(chan string, 1024000)
)

var (
	promAstChannelGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: prometheus.BuildFQName(namespace, "", "current_channels"),
		Help: "Created new channel count",
	})
)

func main() {

	flag.Parse()
	log.SetFlags(0)

	chanInterrupt := make(chan os.Signal, 1)
	signal.Notify(chanInterrupt, os.Interrupt)

	go ariHandler(chanAriMsg)

	go ariMsgHandler(chanAriMsg)

	go promHandler()

	// interrupt handler
	select {
	case <-chanInterrupt:
		log.Println("Interrupted.")

		break
	}

	fmt.Println("Finished.")

	return
}

func promHandler() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*promAddress, nil)
}

func ariHandler(chanAriMsg chan string) {
	// create url parameter
	rawQuery := fmt.Sprintf("api_key=%s&app=%s&subscribeAll=true", *astAriUserinfo, *astAriAppname)

	u := url.URL{
		Scheme:   "ws",
		Host:     *astAriAddress,
		Path:     "/ari/events",
		RawQuery: rawQuery,
	}
	log.Printf("Dial string: %s", u.String())

	for {
		// connect
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Println("Could not connect to server. err: ", err)

			// sleep for every second
			time.Sleep(1 * time.Second)
			continue
		}
		defer c.Close()

		if err := initAri(); err != nil {
			log.Println("Could not initiate ARI. err: ", err)
			continue
		}

		// receiver
		for {
			msgType, msgStr, err := c.ReadMessage()
			if err != nil {
				log.Printf("Could not read message. msgType: %d, err: %s", msgType, err)
				break
			}
			// log.Printf("Message received. type: %d, message: %s", msgType, msgStr)

			// insert msg into queue
			chanAriMsg <- string(msgStr)
		}

		// sleep 1 second for reconnect
		time.Sleep(1 * time.Second)
	}
}

func initAri() error {
	// send channels
	url := fmt.Sprintf("http://%s/ari/channels?api_key=%s", *astAriAddress, *astAriUserinfo)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var msg []interface{}
	if err := json.Unmarshal([]byte(body), &msg); err != nil {
		log.Println("Could not parse the message. err: ", err)
		return err
	}

	// get current channels
	for i := range msg {
		log.Println(msg[i])
		promAstChannelGauge.Inc()
	}

	return nil
}

func ariMsgAnalyzer(msg map[string]interface{}) {
	switch msg["type"] {
	case "ChannelCreated":
		promAstChannelGauge.Inc()

	case "ChannelDestroyed":
		promAstChannelGauge.Dec()
	}
}

func ariMsgHandler(chanMsg chan string) {
	for {

		raw := <-chanMsg
		// log.Printf("Received message: %s", msg)

		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			log.Println("Could not parse the message. err: ", err)
			continue
		}
		log.Printf("Type: %s", msg["type"])

		ariMsgAnalyzer(msg)
	}
}
