package wifi

import (
	"fmt"
	"github.com/yearnming/wifi/pkg/setting"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// Stat 状态
type Stat string

// stat list
const (
	Associating    Stat = setting.AssociatingStatText
	Authenticating Stat = setting.AuthenticatingStatText
	Disconnecting  Stat = setting.DisconnectingStatText
	Disconnected   Stat = setting.DisconnectedStatText
	Connected      Stat = setting.ConnectedStatText
)

// Stats stat map
var Stats = map[Stat]struct{}{
	Associating:    {},
	Connected:      {},
	Disconnected:   {},
	Disconnecting:  {},
	Authenticating: {},
}

// GetWIFIStat GetWIFIStat
func GetWIFIStat() (stat Stat, err error) {
	cmd := exec.Command("netsh", "wlan", "show", "interface")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	return parseWlanStat(string(output))
}

func parseWlanStat(output string) (stat Stat, err error) {
	val, err := simplifiedchinese.GB18030.NewDecoder().String(output)
	if err == nil {
		output = val
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		arr := strings.Split(line, ":")
		if strings.TrimSpace(arr[0]) == setting.StatText {
			stat = Stat(strings.TrimSpace(arr[1]))
			if _, ok := Stats[Stat(stat)]; ok {
				return stat, nil
			}
			return "", fmt.Errorf("unexpected stat: %s", output)
		}
	}

	return "", fmt.Errorf("get stat failed: %s", output)
}
