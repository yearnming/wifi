package scanpasswd

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"sort"

	"github.com/projectdiscovery/gologger"
)

var (
	reSSID = regexp.MustCompile(`.*\s* :\s*(.+)`)
	rePass = regexp.MustCompile(`(?:Key Content|关键内容)\s*:\s*(.+)`) //关键内容

	errExec = errors.New("wifiscan: netsh execution failed")
)

type Profile struct {
	SSID     string
	Password string
	Error    error
}

// ListProfiles 并发获取所有 Wi-Fi 配置（调试版）
func ListProfiles() ([]Profile, error) {
	// 1) 枚举 SSID
	// gologger.Info().Msg("[Step-1] 开始枚举本机保存的 WLAN 配置文件")
	cmd := exec.Command("netsh", "wlan", "show", "profiles")
	out, err := cmd.Output()
	// gologger.Info().Msgf("[Step-1] netsh 执行成功，输出 %s ", out)
	if err != nil {
		gologger.Error().Msgf("[Step-1] netsh 执行失败: %v", err)
		return nil, errExec
	}

	matches := reSSID.FindAllSubmatch(out, -1)
	gologger.Info().Msgf("[Step-1] 正则匹配到 %d 个 SSID", len(matches))
	if len(matches) == 0 {
		gologger.Error().Msg("[Step-1] 没有匹配到任何 SSID，直接返回空列表")
		return nil, nil
	}

	// // 2) 并发抓密码
	// var (
	// 	list []Profile
	// 	mu   sync.Mutex
	// 	wg   sync.WaitGroup
	// )
	// wg.Add(len(matches))
	// gologger.Info().Msgf("[Step-2] 开始并发抓取密码...")
	// for _, m := range matches {
	// 	ssid := string(bytes.TrimSpace(m[1]))
	// 	// gologger.Info().Msgf("SSID: %s", ssid)
	// 	go func(ssid string) {
	// 		defer wg.Done()
	// 		p := Profile{SSID: ssid}
	// 		p.Password, p.Error = getPassword(ssid)

	// 		mu.Lock()
	// 		list = append(list, p)
	// 		mu.Unlock()
	// 	}(ssid)
	// }
	// wg.Wait()

	// 2) 串行抓密码
	list := make([]Profile, 0, len(matches))
	for _, m := range matches {
		ssid := string(bytes.TrimSpace(m[1]))
		p := Profile{SSID: ssid}
		p.Password, p.Error = getPassword(ssid)
		list = append(list, p)
	}

	// 3) 排序 & 汇总
	sort.Slice(list, func(i, j int) bool { return list[i].SSID < list[j].SSID })
	gologger.Info().Msgf("[Step-3] 全部完成，共拿到 %d 条记录", len(list))
	return list, nil
}

// getPassword 单条提取密码（调试版）
func getPassword(ssid string) (string, error) {
	cmd := exec.Command("netsh", "wlan", "show", "profile", ssid, "key=clear")
	out, err := cmd.Output()
	if err != nil {
		gologger.Error().Msgf("获取 %s 密码时 netsh 失败: %v", ssid, err)
		return "", errExec
	}
	// 调试用：整段原样打印
	// gologger.Info().Msgf("=== netsh out for %q ===\n%s\n========================\n", ssid, out)
	if m := rePass.FindSubmatch(out); len(m) > 1 {
		pass := string(bytes.TrimSpace(m[1]))
		return pass, nil
	}

	return "", nil
}
