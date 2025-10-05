package main

import (
	"os"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/yearnming/wifi/pkg/util"
	"github.com/yearnming/wifi/pkg/wifi"
)

var failDB *util.DB // 全局单例

func main() {
	// 测试提交
	//gologger.DefaultLogger.SetTimestamp(true, levels.LevelInfo)
	//gologger.DefaultLogger.SetTimestamp(true, levels.LevelError)
	var start = time.Now()
	Opt := wifi.ParseOptions()

	// // ===== 1. 初始化失败库 =====
	// failDB = util.New("filedb.json") // 可改成你喜欢的路径
	// abs, _ := filepath.Abs("filedb.json")
	// gologger.Info().Msgf("filedb 绝对路径: %s", abs)
	// if err := failDB.Load(); err != nil {
	// 	gologger.Fatal().Msgf("加载失败库失败: %v", err)
	// }
	// defer func() {
	// 	if err := failDB.Save(); err != nil {
	// 		gologger.Error().Msgf("保存失败库失败: %v", err)
	// 	}
	// }()

	// // ===== 捕获 Ctrl+C =====
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-c
	// 	gologger.Info().Msgf("收到中断信号，正在保存失败库...")
	// 	if err := failDB.Save(); err != nil {
	// 		gologger.Error().Msgf("信号保存失败: %v", err)
	// 	}
	// 	os.Exit(0) // 手动退出，确保 defer 不再跑
	// }()

	// // ===== 2. 把 failDB 塞进 wifi 包 =====
	// wifi.SetFailDB(failDB) // 见下一步

	//filename := Dict.Dict
	switch Opt.Mode {
	case "crack":
		wifi.RunCrack(start, Opt) // 单 Wi-Fi → 你现在的逐行逻辑搬进去
	case "multi":
		wifi.RunMulti(start, Opt) // 多 Wi-Fi → 循环调用 runCrack
	case "auto":
		wifi.RunAuto(start, Opt) // 自动排序 → 循环调用 runCrack
	case "verify":
		wifi.TrySavedPasswords()
	case "show-saved":
		wifi.ShowSaved() // 仅打印系统已保存密码
	default:
		gologger.Info().Msgf("不正确的mode!")
		os.Exit(0)
	}
	// // get wifi list
	// networks, err := wifi.GetWIFINetworks()
	// if err != nil {
	// 	gologger.Info().Msgf("Get WIFI Networks: %s", err)
	// 	return
	// }

	// // print wifi list
	// gologger.Info().Msgf("WiFi 列表:")
	// //fmt.Printf("%v", networks)
	// for i, n := range networks {
	// 	fmt.Printf("index: %d, 信号: %d%%, BSSID: %s, SSID: %s\n\n", i, n.Signal, n.BSSID, n.SSID)
	// }

	// // choose wifi
	// gologger.Info().Msgf("选择 WiFi index: ")
	// var selected int
	// fmt.Scanln(&selected)
	// if selected < 0 || selected >= len(networks) {
	// 	gologger.Error().Msgf("选择无效")
	// 	return
	// }
	// selection := networks[selected]
	// gologger.Info().Msgf("你的选择: %v", selection.SSID)

	// var counter int
	// for _, asd := range Opt.WifiDict {
	// 	counter++
	// 	now := time.Now()
	// 	fmt.Println("-------------------------- 第", counter, "次尝试--------------------------")

	// 	gologger.Info().Msgf("测试 password: %s", asd)
	// 	wc := wifi.New(selection.SSID, asd)

	// 	// 1. 带超时连接
	// 	wifi.DefaultConnectTimeout = 4 * time.Second
	// 	stat, err := wc.Connect()
	// 	if err != nil {
	// 		gologger.Error().Msgf("连接 WiFi 失败: %s,状态：%s", err, stat)
	// 		// 2. 继续跑下一个密码，而不是直接退出
	// 		continue
	// 	}

	// 	if stat == wifi.Connected {
	// 		gologger.Info().TimeStamp().Msgf("连接状态：【%s】 成功连接 WiFi (%s) : 密码(%s)", stat, selection.SSID, asd)
	// 		if err := util.WriteToFile(setting.SuccessPwdSavePath,
	// 			fmt.Sprintf("WiFi: %s, password: %s\n", wc.Ssid, asd)); err != nil {
	// 			gologger.Error().Msgf("将密码写入文件 ERR: %s ", err)
	// 		}
	// 		return // 成功就可退出了
	// 	}

	// 	gologger.Error().TimeStamp().Msgf("连接 WiFi 失败")
	// 	if err := wc.DeleteProfile(); err != nil {
	// 		gologger.Error().Msgf("删除配置文件失败: %s", err)
	// 		continue
	// 	}
	// 	gologger.Info().Msgf("删除配置文件成功")
	// 	gologger.Info().Msgf("总计时间: %s, 当前花费时间: %s",
	// 		time.Since(start).Truncate(time.Second).String(),
	// 		time.Since(now).Truncate(time.Second).String(),
	// 	)
	// }

}

// type Options struct {
// 	Dict      string   // Wi-Fi 密码字典文件
// 	WifiDict  []string // 密码字典的内容
// 	Mode      string   // 模式：crack, multi, auto, help
// 	ShowSaved bool     // 是否显示系统保存的密码
// }

// func ParseOptions1() *Options {
// 	logutil.DisableDefaultLogger()

// 	options := &Options{}

// 	var err error
// 	flagSet := goflags.NewFlagSet()
// 	flagSet.SetDescription(`wifi 自动尝试密码`)

// 	flagSet.CreateGroup("input", "Input",
// 		flagSet.StringVarP(&options.Dict, "dict", "l", "", "wifi密码字典"),
// 	)
// 	flagSet.SetCustomHelpText("使用示例:\ngo run main.go -l common.txt")

// 	if err := flagSet.Parse(); err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}

// 	err = options.validateOptions()
// 	if err != nil {
// 		gologger.Fatal().Msgf("程序退出: %s\n", err)
// 	}
// 	return options
// }
// func ParseOptions() *Options {
// 	logutil.DisableDefaultLogger()

// 	options := &Options{}

// 	var err error
// 	flagSet := goflags.NewFlagSet()
// 	flagSet.SetDescription(`自动破解WiFi密码`)

// 	// 添加新的命令行选项
// 	flagSet.CreateGroup("input", "Input",
// 		flagSet.StringVarP(&options.Dict, "dict", "l", "", "Wi-Fi 密码字典文件"),
// 	)

// 	// 增加 --mode 参数，用来选择不同的模式（破解单个、多个、自动等）
// 	flagSet.CreateGroup("mode", "Mode",
// 		flagSet.StringVarP(&options.Mode, "mode", "m", "help", "选择程序运行模式：crack, multi, auto, help"),
// 	)

// 	// 添加 --show-saved 参数，用于显示系统保存的密码
// 	flagSet.CreateGroup("saved", "Saved",
// 		flagSet.BoolVarP(&options.ShowSaved, "show-saved", "s", false, "显示系统保存的 Wi-Fi 密码"),
// 	)

// 	// 设置帮助文本
// 	flagSet.SetCustomHelpText("使用示例:\n" +
// 		"go run main.go --mode crack --dict common.txt  // 破解单个指定 Wi-Fi\n" +
// 		"go run main.go --mode multi --dict common.txt  // 破解多个 Wi-Fi\n" +
// 		"go run main.go --mode auto --dict common.txt   // 自动破解 Wi-Fi\n" +
// 		"go run main.go --show-saved                   // 显示已保存的 Wi-Fi 密码")

// 	// 解析命令行参数
// 	if err := flagSet.Parse(); err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}
// 	// 未给任何参数时主动打印帮助并退出
// 	if options.Mode == "help" {
// 		// 或 fmt.Println(flagSet.GetHelpText())
// 		flagSet.CommandLine.PrintDefaults()
// 		os.Exit(0)
// 	}
// 	// 校验选项有效性
// 	err = options.validateOptions()

// 	if err != nil {
// 		gologger.Fatal().Msgf("程序退出: %s\n", err)
// 	}

// 	return options
// }

// // validateOptions 验证传递的配置选项
// func (options *Options) validateOptions() error {
// 	if options.Dict == "" {
// 		options.WifiDict = options.LoadWifiDict(config.Dictfile)
// 		gologger.Info().Msgf("使用默认字典")

// 		//return errors.New("没有提供参数")
// 	} else {
// 		// 使用 Stat 函数检查文件
// 		if _, err := os.Stat(options.Dict); err == nil {
// 			gologger.Info().Msgf("文件存在: %s\n", options.Dict)
// 		} else if os.IsNotExist(err) {
// 			gologger.Error().Msgf("WiFi密码字典不存在: %s\n", options.Dict)
// 			return err
// 		} else {
// 			gologger.Error().Msgf("检查文件时发生错误: %s\n", err)
// 			return err
// 		}
// 		//filename := "common.txt"
// 		// 打开文件
// 		file, err := os.Open(options.Dict)
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer file.Close()

// 		// 创建一个扫描器来读取文件
// 		scanner := bufio.NewScanner(file)

// 		// 读取所有行到一个切片中
// 		//var lines []string
// 		for scanner.Scan() {
// 			options.WifiDict = append(options.WifiDict, scanner.Text())
// 		}

// 		// 检查读取过程中是否有错误发生
// 		if err := scanner.Err(); err != nil {
// 			panic(err)
// 		}

// 	}
// 	// 校验模式
// 	if options.Mode != "crack" && options.Mode != "multi" && options.Mode != "auto" && options.Mode != "help" {
// 		return errors.New("无效的模式选项，支持的模式: crack, multi, auto, help")
// 	}
// 	return nil
// }

// // LoadWifiDict 加载内置WiFi密码字典
// func (options *Options) LoadWifiDict(dictContent string) []string {
// 	// 根据操作系统的换行符进行分割
// 	Dict := strings.Split(strings.ReplaceAll(dictContent, "\r\n", "\n"), "\n")
// 	gologger.Info().Msgf("字典字数：%v", len(Dict))
// 	return Dict
// }
