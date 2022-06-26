package collector

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"gitlab.com/voipbin/voip/asterisk-exporter.git/models/channel"
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

// // channelCollects collects the channel info
// example: asterisk -rx "core show channels concise"
// PJSIP/test-00000000!demo!1234!50005!Up!Dial!Console/dsp,20!test!!!3!13!df20e970-d58e-476c-8738-24f345bfd190!1616987975.0
// Console/dsp!default!!1!Up!AppDial!(Outgoing Line)!1234!!!3!13!df20e970-d58e-476c-8738-24f345bfd190!1616987975.1
func (h *collector) channelCollects() error {
	log := logrus.WithFields(logrus.Fields{
		"func": "channelCollects",
	})

	res, err := exec.Command("asterisk", "-rx", `core show channels concise`).Output()
	if err != nil {
		log.Errorf("Could not execute the asterisk command. err: %v", err)
		return err
	}

	channels := h.channelParser(string(res))

	channelTech := map[string]int{}    // channel tech: count
	channelContext := map[string]int{} // channel context: count

	for _, channel := range channels {
		// tech
		tech := getTech(channel.Name)
		channelTech[tech]++

		// channel context
		channelContext[channel.Context]++

		promChannelDuration.WithLabelValues(tech, channel.Context).Observe(float64(channel.CallDuration))
	}

	// set metrics
	for k, v := range channelTech {
		promCurrentChannelTech.WithLabelValues(k).Set(float64(v))
	}
	for k, v := range channelContext {
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

// asterisk -rx "core show channels concise"
// PJSIP/call-in-0000060e!call-in!+821100000005!8!Up!Stasis!voipbin,context=call-in,domain=pstn.voipbin.net,source=211.200.20.28!test11!!!3!47!5b20bec1-14a3-4b71-b3e1-7ed4d290b244!asterisk-call-6f58d55c9c-mdqhx-1656228434.2944
func (h *collector) channelParser(data string) []*channel.Channel {

	log := logrus.WithFields(logrus.Fields{
		"func": "channelParser",
	})

	res := []*channel.Channel{}

	datalines := strings.Split(string(data), "\n")
	for _, dataline := range datalines {
		items := strings.Split(dataline, "!")
		if len(items) < 2 {
			continue
		}

		// parse channel duration
		duration, err := strconv.Atoi(items[idxChannelCallDuration])
		if err != nil {
			log.Errorf("Could not parse the channel duration. err: %v", err)
			duration = 0
		}

		tmp := &channel.Channel{
			Name:            items[idxChannelName],
			Context:         items[idxChannelContext],
			Extension:       items[idxChannelExtension],
			Priority:        items[idxChannelPriority],
			State:           items[idxChannelState],
			Application:     items[idxChannelApplication],
			ApplicationData: items[idxChannelApplicationData],
			CallerNumber:    items[idxChannelCallerNumber],
			AccountCode:     items[idxChannelAccountCode],
			Account:         items[idxChannelAccount],

			AMAFlag:      items[idxChannelAMAFlag],
			CallDuration: duration,
			BridgeID:     items[idxChannelBridgeID],
			UniqueID:     items[idxChannelUniqueID],
		}

		res = append(res, tmp)
	}

	return res
}
