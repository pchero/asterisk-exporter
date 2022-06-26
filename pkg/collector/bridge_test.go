package collector

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"gitlab.com/voipbin/voip/asterisk-exporter.git/models/bridge"
)

func Test_convertBridgeDuration(t *testing.T) {

	tests := []struct {
		name string

		duration string

		expectRes time.Duration
	}{
		{
			"normal",

			"218:44:47",

			time.Duration(787487000000000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			h := &collector{}

			res := h.convertBridgeDuration(tt.duration)
			if res != tt.expectRes {
				t.Errorf("Wrong match.\nexpect: %v\ngot: %v", tt.expectRes, res)
			}
		})
	}
}

func Test_bridgeParser(t *testing.T) {

	tests := []struct {
		name string

		data string

		expectRes []bridge.Bridge
	}{
		{
			"simple",

			`Bridge-ID                            Chans Type            Technology      Duration
			008e2905-1aa7-4106-b388-3e3f5c157265     0 stasis          simple_bridge   220:34:45`,

			[]bridge.Bridge{
				{
					ID:         "008e2905-1aa7-4106-b388-3e3f5c157265",
					Chans:      0,
					Type:       "stasis",
					Technology: "simple_bridge",
					Duration:   794085,
				},
			},
		},
		{
			"more than 2 items",

			`Bridge-ID                            Chans Type            Technology      Duration
			008e2905-1aa7-4106-b388-3e3f5c157265     0 stasis          simple_bridge   220:34:45
			010d3e2d-282b-422c-a07d-4d440d4bb3c6     0 stasis          simple_bridge   221:46:48
			04ec2e17-6830-4db5-8f7a-5258fd73df3f     0 stasis          simple_bridge   221:42:00`,

			[]bridge.Bridge{
				{
					ID:         "008e2905-1aa7-4106-b388-3e3f5c157265",
					Chans:      0,
					Type:       "stasis",
					Technology: "simple_bridge",
					Duration:   794085,
				},
				{
					ID:         "010d3e2d-282b-422c-a07d-4d440d4bb3c6",
					Chans:      0,
					Type:       "stasis",
					Technology: "simple_bridge",
					Duration:   798408,
				},
				{
					ID:         "04ec2e17-6830-4db5-8f7a-5258fd73df3f",
					Chans:      0,
					Type:       "stasis",
					Technology: "simple_bridge",
					Duration:   798120,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			h := &collector{}

			res := h.bridgeParser(tt.data)
			if !reflect.DeepEqual(res, tt.expectRes) {
				t.Errorf("Wrong match.\nexpect: %v\ngot: %v", tt.expectRes, res)
			}
		})
	}
}
