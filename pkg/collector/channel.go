package collector

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// List of channel concise parse item index
const (

	// channel name ! context ! extension ! priority ! state ! application ! application data ! caller number ! account code ! account ! ama flag ! call duration ! bridge id ! unique id

	idxChannelName            = 0
	idxChannelContext         = 1
	idxChannelExtension       = 2
	idxChannelPriority        = 3
	idxChannelState           = 4
	idxChannelApplication     = 5
	idxChannelApplicationData = 6
	idxChannelCallerNumber    = 7
	idxChannelAccountCode     = 8
	idxChannelAccount         = 9

	idxChannelAMAFlag      = 10
	idxChannelCallDuration = 11
	idxChannelBridgeID     = 12
	idxChannelUniqueID     = 13
)

var (
	channelTech    map[string]int = map[string]int{}
	channelContext map[string]int = map[string]int{}
)

// Collect collects the asterisk's metrics and update the prometheus metric
func (h *collector) Collect() error {

	if err := h.collectChannels(); err != nil {
		logrus.Errorf("Could not get channel metrics. err: %v", err)
		return err
	}

	return nil
}

// collectChannels collects the channel info
// example: asterisk -rx "core show channels concise"
// PJSIP/test-00000000!demo!1234!50005!Up!Dial!Console/dsp,20!test!!!3!13!df20e970-d58e-476c-8738-24f345bfd190!1616987975.0
// Console/dsp!default!!1!Up!AppDial!(Outgoing Line)!1234!!!3!13!df20e970-d58e-476c-8738-24f345bfd190!1616987975.1
func (h *collector) collectChannels() error {
	res, err := exec.Command("asterisk", "-rx", `core show channels concise`).Output()
	if err != nil {
		return err
	}

	// reset metrics
	for k := range channelTech {
		channelTech[k] = 0
	}
	for k := range channelContext {
		channelContext[k] = 0
	}

	channels := strings.Split(string(res), "\n")
	for _, channel := range channels {
		c := strings.Split(channel, "!")
		if len(c) < 2 {
			continue
		}

		// tech
		tech := getTech(c[idxChannelName])
		if _, ok := channelTech[tech]; !ok {
			channelTech[tech] = 0
		}
		channelTech[tech] = channelTech[tech] + 1

		// context
		ctx := c[idxChannelContext]
		if _, ok := channelContext[ctx]; !ok {
			channelContext[ctx] = 0
		}
		channelContext[ctx] = channelContext[ctx] + 1

		// set channel duration
		duration, err := strconv.Atoi(c[idxChannelCallDuration])
		if err != nil {
			logrus.Errorf("Could not parse the channel duration. err: %v", err)
			continue
		}
		promChannelDuration.WithLabelValues(tech, ctx).Observe(float64(duration))
	}

	// set metrics
	for k, v := range channelTech {
		promCurrentChannelTech.WithLabelValues(k).Set(float64(v))
		if v <= 0 {
			delete(channelTech, k)
		}
	}
	for k, v := range channelContext {
		promCurrentChannelContext.WithLabelValues(k).Set(float64(v))
		if v <= 0 {
			delete(channelContext, k)
		}
	}

	return nil
}

func getTech(c string) string {

	items := strings.Split(c, "/")
	if len(items) < 1 {
		return "UNKNOWN"
	}

	return items[0]
}
