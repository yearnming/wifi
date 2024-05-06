package wifi

import (
	"os/exec"
	"sort"
	"strconv"
	"strings"
	//"wifi/
	"github.com/yearnming/wifi/pkg/setting"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// Network Network
type Network struct {
	SSID   string
	BSSID  string
	Signal int
}

// GetWIFINetworks GetWIFINetworks
func GetWIFINetworks() (networks []Network, err error) {
	cmd := exec.Command("netsh", "wlan", "show", "networks", "mode=Bssid")
	output, err := cmd.Output()
	//fmt.Printf("cmd: %v \n output: %v", cmd, output)
	if err != nil {
		return
	}

	networks = parseWiFiList(string(output))

	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Signal > networks[j].Signal
	})

	return
}

func parseWiFiList(output string) []Network {
	val, err := simplifiedchinese.GB18030.NewDecoder().String(output)
	if err == nil {
		output = val
	}

	lines := strings.Split(output, "\n")
	var networks []Network
	var currentNetwork Network

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, setting.SSIDText) {
			index := strings.Index(line, ":")
			if index != -1 {
				currentNetwork.SSID = strings.TrimSpace(line[index+1:])
			}
			continue
		}

		if strings.HasPrefix(line, setting.BSSIDText) {
			index := strings.Index(line, ":")
			if index != -1 {
				currentNetwork.BSSID = strings.TrimSpace(line[index+1:])
			}
			continue
		}

		if strings.HasPrefix(line, setting.SignalText) {
			index := strings.Index(line, ":")
			if index != -1 {
				n, _ := strconv.Atoi(strings.TrimRight(strings.TrimSpace(line[index+1:]), "%"))
				currentNetwork.Signal = n

				if currentNetwork.SSID != "" {
					networks = append(networks, currentNetwork)
				}

				currentNetwork = Network{}
			}
		}
	}

	return networks
}
