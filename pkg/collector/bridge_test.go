package collector

import (
	"reflect"
	"testing"
	"time"

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

			`Bridge-ID                            Name                                 Chans Type            Technology      Duration
			008e2905-1aa7-4106-b388-3e3f5c157265 reference_type=call,reference_id=5cef1dbc-13d8-11ec-9199-f3e0e965a469     0 stasis          simple_bridge   220:34:45`,

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

			`Bridge-ID                            Name                                 Chans Type            Technology      Duration
			008e2905-1aa7-4106-b388-3e3f5c157265 reference_type=call,reference_id=5cef1dbc-13d8-11ec-9199-f3e0e965a469     0 stasis          simple_bridge   220:34:45
			010d3e2d-282b-422c-a07d-4d440d4bb3c6 reference_type=confbridge,reference_id=8f537474-13d8-11ec-9193-7b377238c934     0 stasis          simple_bridge   221:46:48
			04ec2e17-6830-4db5-8f7a-5258fd73df3f reference_type=call-snoop,reference_id=7f6dbc1a-02fb-11ec-897b-ef9b30e25c57     0 stasis          simple_bridge   221:42:00`,

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
			`Bridge-ID                            Name                                 Chans Type            Technology      Duration`,
			0,
		},
		{
			"too few fields",
			`Bridge-ID                            Name                                 Chans Type            Technology      Duration
008e2905-1aa7-4106-b388-3e3f5c157265 reference_type=call,reference_id=5cef1dbc-13d8-11ec-9199-f3e0e965a469     0 stasis`,
			0,
		},
		{
			"5 fields instead of 6",
			`Bridge-ID                            Name                                 Chans Type            Technology      Duration
008e2905-1aa7-4106-b388-3e3f5c157265 reference_type=call,reference_id=5cef1dbc-13d8-11ec-9199-f3e0e965a469     0 stasis          simple_bridge`,
			0,
		},
		{
			"mixed valid and malformed lines",
			`Bridge-ID                            Name                                 Chans Type            Technology      Duration
008e2905-1aa7-4106-b388-3e3f5c157265 reference_type=call,reference_id=5cef1dbc-13d8-11ec-9199-f3e0e965a469     0 stasis          simple_bridge   220:34:45
incomplete line
010d3e2d-282b-422c-a07d-4d440d4bb3c6 reference_type=confbridge,reference_id=8f537474-13d8-11ec-9193-7b377238c934     0 stasis          simple_bridge   221:46:48`,
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

	data := `Bridge-ID                            Name                                 Chans Type            Technology      Duration
008e2905-1aa7-4106-b388-3e3f5c157265 reference_type=call,reference_id=5cef1dbc-13d8-11ec-9199-f3e0e965a469     0 stasis          simple_bridge   220:34:45
incomplete line
010d3e2d-282b-422c-a07d-4d440d4bb3c6 reference_type=confbridge,reference_id=8f537474-13d8-11ec-9193-7b377238c934     0 stasis          simple_bridge   221:46:48`

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

// Test_bridgeParser_voipbinReferenceTaggedBridge is a regression test for VOIP-1468.
//
// Before this fix, a VoIPBin-tagged bridge Name (e.g.
// "reference_type=confbridge,reference_id=<uuid>", set by bin-call-manager on every
// call/confbridge bridge) was mistaken for the Chans column because the parser's index
// table didn't account for Asterisk's Name column. Since that Name string never
// converts to an integer, strconv.Atoi failed on every single VoIPBin bridge and
// silently reported Chans as 0, while also flooding logs with an ERROR per bridge.
func Test_bridgeParser_voipbinReferenceTaggedBridge(t *testing.T) {

	h := &collector{}

	data := `Bridge-ID                            Name                                 Chans Type            Technology      Duration
ff63ebcc-e6d8-430d-bb84-6633ebdcb2c6 reference_type=confbridge,reference_id=b72c115f-fee9-4cfe-86bf-99ec4b0fcf4c     2 stasis          simple_bridge   126:38:06`

	res := h.bridgeParser(data)
	if len(res) != 1 {
		t.Fatalf("Expected 1 parsed bridge, got %d", len(res))
	}
	if res[0].ID != "ff63ebcc-e6d8-430d-bb84-6633ebdcb2c6" {
		t.Errorf("Expected bridge ID ff63ebcc-e6d8-430d-bb84-6633ebdcb2c6, got %s", res[0].ID)
	}
	if res[0].Chans != 2 {
		t.Errorf("Expected Chans 2 (not misparsed from the Name column), got %d", res[0].Chans)
	}
	if res[0].Type != "stasis" {
		t.Errorf("Expected Type stasis, got %s", res[0].Type)
	}
	if res[0].Technology != "simple_bridge" {
		t.Errorf("Expected Technology simple_bridge, got %s", res[0].Technology)
	}
}
