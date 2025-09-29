package wifi

import (
	"testing"

	"github.com/projectdiscovery/gologger"
)

func TestGetWIFINetwork(t *testing.T) {
	networks, err := GetWIFINetworks()
	if err != nil {
		// t.Fatal(err)
		gologger.Error().Msgf(err.Error())
	}

	for _, network := range networks {
		// t.Logf("SSID: %s, BSSID: %s\n", network.SSID, network.BSSID)
		gologger.Info().Msgf("SSID: %s, BSSID: %s\n", network.SSID, network.BSSID)
	}
}
