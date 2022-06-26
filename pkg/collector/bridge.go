package collector

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"gitlab.com/voipbin/voip/asterisk-exporter.git/models/bridge"
)

// List of bridge parse item index
const (

	// Bridge-ID    Chans    Type    Technology    Duration

	idxBridgeID         = 0
	idxBridgeChannels   = 1
	idxBridgeType       = 2
	idxBridgeTechnology = 3
	idxBridgeDuration   = 4
)

// bridgeCollects collects the bridge info
// example: asterisk -rx "bridge show all"
// 008e2905-1aa7-4106-b388-3e3f5c157265     0 stasis          simple_bridge   218:44:47
func (h *collector) bridgeCollects() error {
	log := logrus.WithFields(logrus.Fields{
		"func": "bridgeCollects",
	})

	res, err := exec.Command("asterisk", "-rx", `bridge show all`).Output()
	if err != nil {
		log.Errorf("Could not execute asterisk command. err: %v", err)
		return err
	}
	log.Debugf("tmp res. res: %v", res)

	bridges := h.bridgeParser(string(res))

	// type:tech:count
	mapBridgeCount := map[string]map[string]int{}

	for _, bridge := range bridges {
		_, ok := mapBridgeCount[bridge.Type]
		if !ok {
			mapBridgeCount[bridge.Type] = map[string]int{}
		}

		mapBridgeCount[bridge.Type][bridge.Technology]++

		promBridgeDuration.WithLabelValues(bridge.Type, bridge.Technology).Observe(bridge.Duration)
	}

	// set count
	for bridgeType, tmp := range mapBridgeCount {
		for bridgeTech, bridgeCount := range tmp {
			promCurrentBridgeCount.WithLabelValues(bridgeType, bridgeTech).Set(float64(bridgeCount))
		}
	}

	return nil
}

// bridgeCollects collects the bridge info
// example: asterisk -rx "bridge show all"
// 008e2905-1aa7-4106-b388-3e3f5c157265     0 stasis          simple_bridge   218:44:47
func (h *collector) bridgeParser(data string) []bridge.Bridge {
	log := logrus.WithFields(logrus.Fields{
		"func": "bridgeParser",
	})

	res := []bridge.Bridge{}

	datalines := strings.Split(string(data), "\n")
	for i, dataline := range datalines {
		if i == 0 {
			continue
		}

		items := strings.Fields(dataline)
		if len(items) < 2 {
			continue
		}

		// chans
		chans, err := strconv.Atoi(items[idxBridgeChannels])
		if err != nil {
			log.Errorf("Could not parse the bridge channels correctly. chans: %s, err: %v", items[idxBridgeChannels], err)
			chans = 0
		}

		// duration
		duration := h.convertBridgeDuration(items[idxBridgeDuration])

		tmp := bridge.Bridge{
			ID:         items[idxBridgeID],
			Chans:      chans,
			Type:       items[idxBridgeType],
			Technology: items[idxBridgeTechnology],
			Duration:   duration.Seconds(),
		}

		res = append(res, tmp)
	}

	return res
}

// convertBridgeDuration convert the bridge's durationto the time.Duration
func (h *collector) convertBridgeDuration(t string) time.Duration {
	/// 218:44:47
	t = strings.Replace(t, ":", "h", 1) // hour
	t = strings.Replace(t, ":", "m", 1) // minute
	t = t + "s"                         // second

	res, err := time.ParseDuration(t)
	if err != nil {
		return time.Duration(0)
	}

	return res
}
