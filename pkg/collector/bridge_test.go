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

func Test_bridgeParser_malformed(t *testing.T) {

	tests := []struct {
		name          string
		data          string
		expectedCount int
	}{
		{
			"empty input",
			``,
			0,
		},
		{
			"header only",
			`Bridge-ID                            Chans Type            Technology      Duration`,
			0,
		},
		{
			"too few fields",
			`Bridge-ID                            Chans Type            Technology      Duration
008e2905-1aa7-4106-b388-3e3f5c157265     0 stasis`,
			0,
		},
		{
			"4 fields instead of 5",
			`Bridge-ID                            Chans Type            Technology      Duration
008e2905-1aa7-4106-b388-3e3f5c157265     0 stasis          simple_bridge`,
			0,
		},
		{
			"mixed valid and malformed lines",
			`Bridge-ID                            Chans Type            Technology      Duration
008e2905-1aa7-4106-b388-3e3f5c157265     0 stasis          simple_bridge   220:34:45
incomplete line
010d3e2d-282b-422c-a07d-4d440d4bb3c6     0 stasis          simple_bridge   221:46:48`,
			2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &collector{}

			res := h.bridgeParser(tt.data)
			if len(res) != tt.expectedCount {
				t.Errorf("Expected %d parsed bridges, got %d", tt.expectedCount, len(res))
			}
		})
	}
}

func Test_bridgeParser_mixed(t *testing.T) {

	h := &collector{}

	data := `Bridge-ID                            Chans Type            Technology      Duration
008e2905-1aa7-4106-b388-3e3f5c157265     0 stasis          simple_bridge   220:34:45
incomplete line
010d3e2d-282b-422c-a07d-4d440d4bb3c6     0 stasis          simple_bridge   221:46:48`

	res := h.bridgeParser(data)
	if len(res) != 2 {
		t.Fatalf("Expected 2 parsed bridges, got %d", len(res))
	}
	if res[0].ID != "008e2905-1aa7-4106-b388-3e3f5c157265" {
		t.Errorf("Expected first bridge ID 008e2905-1aa7-4106-b388-3e3f5c157265, got %s", res[0].ID)
	}
	if res[0].Duration != 794085 {
		t.Errorf("Expected first bridge duration 794085, got %f", res[0].Duration)
	}
	if res[1].ID != "010d3e2d-282b-422c-a07d-4d440d4bb3c6" {
		t.Errorf("Expected second bridge ID 010d3e2d-282b-422c-a07d-4d440d4bb3c6, got %s", res[1].ID)
	}
	if res[1].Duration != 798408 {
		t.Errorf("Expected second bridge duration 798408, got %f", res[1].Duration)
	}
}
