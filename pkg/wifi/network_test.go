package wifi

import "testing"

func TestGetWIFINetwork(t *testing.T) {
	networks, err := GetWIFINetworks()
	if err != nil {
		t.Fatal(err)
	}

	for _, network := range networks {
		t.Logf("SSID: %s, BSSID: %s\n", network.SSID, network.BSSID)
	}
}
