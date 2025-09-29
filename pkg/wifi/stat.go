package wifi

import (
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf8"

	"github.com/yearnming/wifi/pkg/setting"

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

// func parseWlanStat(output string) (stat Stat, err error) {
// 	val, err := simplifiedchinese.GB18030.NewDecoder().String(output)
// 	if err == nil {
// 		output = val
// 	}

// 	lines := strings.Split(output, "\n")
// 	for _, line := range lines {
// 		arr := strings.Split(line, ":")
// 		if strings.TrimSpace(arr[0]) == setting.StatText {
// 			stat = Stat(strings.TrimSpace(arr[1]))
// 			if _, ok := Stats[Stat(stat)]; ok {
// 				return stat, nil
// 			}
// 			return "", fmt.Errorf("unexpected stat: %s", output)
// 		}
// 	}

// 	return "", fmt.Errorf("get stat failed: %s", output)
// }

func parseWlanStat(output string) (Stat, error) {
	// 1. 编码转换（控制台 936 → UTF-8）
	if !utf8.ValidString(output) {
		if val, e := simplifiedchinese.GB18030.NewDecoder().String(output); e == nil {
			output = val
		}
	}

	// 2. 逐行扫描
	for _, line := range strings.Split(output, "\n") {
		// 去掉首尾空格，防止左边混进空格或 BOM
		line = strings.TrimSpace(line)

		// 忽略空行
		if line == "" {
			continue
		}

		// 只认“状态”行，其它字段直接跳过
		if !strings.HasPrefix(line, setting.StatText) { // setting.StatText == "状态"
			continue
		}

		// 冒号后面那一段
		idx := strings.Index(line, ":")
		if idx == -1 {
			continue
		}
		raw := strings.TrimSpace(line[idx+1:])

		// 3. 兜底：去掉可能混进来的回车、句号、全角空格
		raw = strings.TrimRight(raw, "\r.。\x00")

		// 4. 查表
		if _, ok := Stats[Stat(raw)]; ok {
			return Stat(raw), nil
		}

		// 5. 出现未知状态，把原文打出去方便排错
		return "", fmt.Errorf("unexpected stat: %q (full line: %q)", raw, line)
	}

	return "", fmt.Errorf("stat line not found in output")
}
