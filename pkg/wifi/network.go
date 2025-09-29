package wifi

import (
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	//"wifi/

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
	// fmt.Printf("cmd: %v \n output: %v", cmd, output)
	// fmt.Printf("cmd: %v\noutput:\n%s", cmd, output)
	if err != nil {
		return
	}

	networks = parseWiFiList(string(output))
	// fmt.Printf("cmd: %v \n output: %v", cmd, networks)
	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Signal > networks[j].Signal
	})

	return
}
func parseWiFiList(output string) []Network {
	// 1. 编码转换
	if utf8.ValidString(output) == false {
		if val, err := simplifiedchinese.GB18030.NewDecoder().String(output); err == nil {
			output = val
		}
	}

	lines := strings.Split(output, "\n")
	var nets []Network
	var cur Network

	// 2. 状态机：只有“SSID 行”出现时才算一条新记录开始
	inBlock := false

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		// SSID 行：SSID 1 : xxx
		if strings.HasPrefix(line, "SSID") && strings.Contains(line, ":") {
			// 如果前面已经扫描到一半，把半成品扔掉
			if inBlock && cur.SSID != "" {
				nets = append(nets, cur)
			}
			cur = Network{}
			inBlock = true
			cur.SSID = strings.TrimSpace(line[strings.Index(line, ":")+1:])
			continue
		}

		// BSSID 行
		if strings.HasPrefix(line, "BSSID") && strings.Contains(line, ":") {
			cur.BSSID = strings.TrimSpace(line[strings.Index(line, ":")+1:])
			continue
		}

		// 信号行
		if strings.HasPrefix(line, "信号") && strings.Contains(line, ":") {
			percStr := strings.TrimRight(strings.TrimSpace(line[strings.Index(line, ":")+1:]), "%")
			if p, err := strconv.Atoi(percStr); err == nil {
				cur.Signal = p
			}
			// 一条完整记录结束
			if cur.SSID != "" {
				nets = append(nets, cur)
			}
			inBlock = false
			cur = Network{}
		}
	}

	// 3. 文件结尾可能还有一条半成品
	if inBlock && cur.SSID != "" {
		nets = append(nets, cur)
	}
	return nets
}

// func parseWiFiList(output string) []Network {
// 	val, err := simplifiedchinese.GB18030.NewDecoder().String(output)
// 	if err == nil {
// 		output = val
// 	}

// 	lines := strings.Split(output, "\n")
// 	var networks []Network
// 	var currentNetwork Network

// 	for _, line := range lines {
// 		line = strings.TrimSpace(line)

// 		if strings.HasPrefix(line, setting.SSIDText) {
// 			index := strings.Index(line, ":")
// 			if index != -1 {
// 				currentNetwork.SSID = strings.TrimSpace(line[index+1:])
// 			}
// 			continue
// 		}

// 		if strings.HasPrefix(line, setting.BSSIDText) {
// 			index := strings.Index(line, ":")
// 			if index != -1 {
// 				currentNetwork.BSSID = strings.TrimSpace(line[index+1:])
// 			}
// 			continue
// 		}

// 		if strings.HasPrefix(line, setting.SignalText) {
// 			index := strings.Index(line, ":")
// 			if index != -1 {
// 				n, _ := strconv.Atoi(strings.TrimRight(strings.TrimSpace(line[index+1:]), "%"))
// 				currentNetwork.Signal = n

// 				if currentNetwork.SSID != "" {
// 					networks = append(networks, currentNetwork)
// 				}

// 				currentNetwork = Network{}
// 			}
// 		}
// 	}

// 	return networks
// }
