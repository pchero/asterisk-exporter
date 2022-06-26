package collector

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"gitlab.com/voipbin/voip/asterisk-exporter.git/models/channel"
)

func Test_channelParser(t *testing.T) {

	tests := []struct {
		name string

		data string

		expectRes []*channel.Channel
	}{
		{
			"simple",

			`PJSIP/call-in-0000060e!call-in!+821100000005!8!Up!Stasis!voipbin,context=call-in,domain=pstn.voipbin.net,source=211.200.20.28!test11!!!3!47!5b20bec1-14a3-4b71-b3e1-7ed4d290b244!asterisk-call-6f58d55c9c-mdqhx-1656228434.2944`,

			[]*channel.Channel{
				{
					Name:            "PJSIP/call-in-0000060e",
					Context:         "call-in",
					Extension:       "+821100000005",
					Priority:        "8",
					State:           "Up",
					Application:     "Stasis",
					ApplicationData: "voipbin,context=call-in,domain=pstn.voipbin.net,source=211.200.20.28",
					CallerNumber:    "test11",
					AccountCode:     "",
					Account:         "",
					AMAFlag:         "3",
					CallDuration:    47,
					BridgeID:        "5b20bec1-14a3-4b71-b3e1-7ed4d290b244",
					UniqueID:        "asterisk-call-6f58d55c9c-mdqhx-1656228434.2944",
				},
			},
		},
		{
			"more than 2 items",

			`PJSIP/call-in-00000610!call-in!+821100000001!8!Up!Stasis!voipbin,context=call-in,domain=pstn.voipbin.net,source=211.200.20.28!test11!!!3!9!5282069a-2abb-49e7-a8fe-0593f2ba6e2a!asterisk-call-6f58d55c9c-mdqhx-1656228956.2947
PJSIP/call-in-0000060f!call-in!+821100000005!8!Up!Stasis!voipbin,context=call-in,domain=pstn.voipbin.net,source=211.200.20.28!test11!!!3!19!0e4e297f-6822-4c55-ad41-c8c8a2519590!asterisk-call-6f58d55c9c-mdqhx-1656228946.2946`,

			[]*channel.Channel{
				{
					Name:            "PJSIP/call-in-00000610",
					Context:         "call-in",
					Extension:       "+821100000001",
					Priority:        "8",
					State:           "Up",
					Application:     "Stasis",
					ApplicationData: "voipbin,context=call-in,domain=pstn.voipbin.net,source=211.200.20.28",
					CallerNumber:    "test11",
					AccountCode:     "",
					Account:         "",
					AMAFlag:         "3",
					CallDuration:    9,
					BridgeID:        "5282069a-2abb-49e7-a8fe-0593f2ba6e2a",
					UniqueID:        "asterisk-call-6f58d55c9c-mdqhx-1656228956.2947",
				},
				{
					Name:            "PJSIP/call-in-0000060f",
					Context:         "call-in",
					Extension:       "+821100000005",
					Priority:        "8",
					State:           "Up",
					Application:     "Stasis",
					ApplicationData: "voipbin,context=call-in,domain=pstn.voipbin.net,source=211.200.20.28",
					CallerNumber:    "test11",
					AccountCode:     "",
					Account:         "",
					AMAFlag:         "3",
					CallDuration:    19,
					BridgeID:        "0e4e297f-6822-4c55-ad41-c8c8a2519590",
					UniqueID:        "asterisk-call-6f58d55c9c-mdqhx-1656228946.2946",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			defer mc.Finish()

			h := &collector{}

			res := h.channelParser(tt.data)
			if !reflect.DeepEqual(res, tt.expectRes) {
				t.Errorf("Wrong match.\nexpect: %v\ngot: %v", tt.expectRes[0], res[0])
			}
		})
	}
}
