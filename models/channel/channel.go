package channel

// Channel defines
type Channel struct {
	Name            string
	Context         string
	Extension       string
	Priority        string
	State           string
	Application     string
	ApplicationData string
	CallerNumber    string
	AccountCode     string
	Account         string

	AMAFlag      string
	CallDuration int
	BridgeID     string
	UniqueID     string
}
