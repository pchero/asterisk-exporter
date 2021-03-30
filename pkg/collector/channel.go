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

	chanTech := map[string]int{}
	chanContext := map[string]int{}

	channels := strings.Split(string(res), "\n")
	for _, channel := range channels {
		c := strings.Split(channel, "!")
		if len(c) < 2 {
			continue
		}

		// tech
		tech := getTech(c[idxChannelName])
		if _, ok := chanTech[tech]; !ok {
			chanTech[tech] = 0
		}
		chanTech[tech] = chanTech[tech] + 1

		// context
		ctx := c[idxChannelContext]
		if _, ok := chanContext[ctx]; !ok {
			chanContext[ctx] = 0
		}
		chanContext[ctx] = chanContext[ctx] + 1

		// set channel duration
		duration, err := strconv.Atoi(c[idxChannelCallDuration])
		if err != nil {
			logrus.Errorf("Could not parse the channel duration. err: %v", err)
			continue
		}
		promChannelDuration.WithLabelValues(tech, ctx).Observe(float64(duration))
	}

	// set metrics
	for k, v := range chanTech {
		promCurrentChannelTech.WithLabelValues(k).Set(float64(v))
	}
	for k, v := range chanContext {
		promCurrentChannelContext.WithLabelValues(k).Set(float64(v))
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
