package wifi

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	logutil "github.com/projectdiscovery/utils/log"
	"github.com/yearnming/wifi/pkg/config"
	"github.com/yearnming/wifi/pkg/util"
)

type Options struct {
	Dict      string   // Wi-Fi 密码字典文件
	WifiDict  []string // 密码字典的内容
	Mode      string   // 模式：crack, multi, auto, help
	ShowSaved bool     // 是否显示系统保存的密码
	Timeout   int      // 新增：连接超时（秒）
	FailDB    string   // 失败库 JSON 路径
}

func ParseOptions() *Options {
	logutil.DisableDefaultLogger()

	options := &Options{}

	var err error
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`自动破解WiFi密码`)

	// 添加新的命令行选项
	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&options.Dict, "dict", "l", "", "Wi-Fi 密码字典文件，已内嵌默认字典，可不指定字典"),
		flagSet.BoolVarP(&options.ShowSaved, "show-saved", "s", false, "显示系统保存的 Wi-Fi 密码"),
		flagSet.IntVarP(&options.Timeout, "timeout", "t", 10, "单密码连接超时（秒）默认10s"),
		flagSet.StringVarP(&options.Mode, "mode", "m", "", "选择程序运行模式：crack, multi, auto，verify"),
		flagSet.StringVarP(&options.FailDB, "fail-db", "d", "", "失败密码库 JSON 路径（默认 ~/.wifi-cracker/filedb.json）"),
	)

	// // 增加 --mode 参数，用来选择不同的模式（破解单个、多个、自动等）
	// flagSet.CreateGroup("mode", "Mode",
	// 	flagSet.StringVarP(&options.Mode, "mode", "m", "", "选择程序运行模式：crack, multi, auto，verify"),
	// )
	// // 放在原来的 CreateGroup 任意位置，或新建一组
	// flagSet.CreateGroup("timeout", "Timeout",
	// 	flagSet.IntVarP(&options.Timeout, "timeout", "t", 4, "单密码连接超时（秒）默认4s"),
	// )
	// // 添加 --show-saved 参数，用于显示系统保存的密码
	// flagSet.CreateGroup("saved", "Saved",
	// 	flagSet.BoolVarP(&options.ShowSaved, "show-saved", "s", false, "显示系统保存的 Wi-Fi 密码"),
	// )

	// 设置帮助文本
	flagSet.SetCustomHelpText("使用示例:\n" +
		"go run main.go --mode crack --dict common.txt  // 破解单个指定 Wi-Fi\n" +
		"go run main.go --mode multi --dict common.txt  // 破解多个 Wi-Fi\n" +
		"go run main.go --mode auto --dict common.txt   // 自动破解 Wi-Fi\n" +
		"go run main.go --show-saved                   // 显示已保存的 Wi-Fi 密码")

	// 解析命令行参数
	if err := flagSet.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// 未给任何参数时主动打印帮助并退出
	if options.Mode == "" && options.ShowSaved == false {
		// 或 fmt.Println(flagSet.GetHelpText())
		flagSet.CommandLine.PrintDefaults()
		os.Exit(0)
	}

	// 校验选项有效性
	err = options.validateOptions()

	if err != nil {
		gologger.Fatal().Msgf("程序退出: %s\n", err)
	}

	return options
}

// validateOptions 验证传递的配置选项
func (options *Options) validateOptions() error {

	// ===== 1. 失败库路径：优先用用户指定的，否则默认 =====
	if options.FailDB == "" {
		options.FailDB = util.DefaultFailDBPath() // 你前面写的 ~/.wifi-cracker/filedb.json
	}
	failDB = util.New(options.FailDB)

	abs, _ := filepath.Abs(options.FailDB) // 只算一次绝对路径
	gologger.Info().Msgf("失败库绝对路径: %s", abs)

	if err := failDB.Load(); err != nil {
		gologger.Fatal().Msgf("加载失败库失败: %v", err)
	}
	defer func() {
		if err := failDB.Save(); err != nil {
			gologger.Error().Msgf("保存失败库失败: %v", err)
		}
	}()

	// ===== 2. Ctrl+C 保盘 =====
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		gologger.Info().Msgf("收到中断信号，正在保存失败库...")
		if err := failDB.Save(); err != nil {
			gologger.Error().Msgf("信号保存失败: %v", err)
		}
		os.Exit(0)
	}()

	// ===== 3. 注入 wifi 包 =====
	SetFailDB(failDB)

	if options.ShowSaved == true {
		ShowSaved()
		os.Exit(0)
		// return nil
	}
	// 在 validateOptions() 里追加
	if options.Timeout <= 0 || options.Timeout > 60 {
		return errors.New("超时时间必须在 1-60 秒之间")
	}

	DefaultConnectTimeout = time.Duration(options.Timeout) * time.Second
	gologger.Info().Msgf("当前超时时间为：%ds", options.Timeout)

	if options.Mode != "crack" && options.Mode != "multi" && options.Mode != "auto" && options.Mode != "verify" {
		return errors.New("无效的模式选项，支持的模式: crack, multi, auto, verify")
	}
	// verify 模式完全不需要字典
	if options.Mode == "verify" {
		return nil
	}

	if options.Dict == "" {
		options.WifiDict = options.LoadWifiDict(config.Dictfile)
		gologger.Info().Msgf("使用默认字典")

		//return errors.New("没有提供参数")
	} else {
		// 使用 Stat 函数检查文件
		if _, err := os.Stat(options.Dict); err == nil {
			gologger.Info().Msgf("文件存在: %s\n", options.Dict)
		} else if os.IsNotExist(err) {
			gologger.Error().Msgf("WiFi密码字典不存在: %s\n", options.Dict)
			return err
		} else {
			gologger.Error().Msgf("检查文件时发生错误: %s\n", err)
			return err
		}
		//filename := "common.txt"
		// 打开文件
		file, err := os.Open(options.Dict)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// 创建一个扫描器来读取文件
		scanner := bufio.NewScanner(file)

		// 读取所有行到一个切片中
		//var lines []string
		for scanner.Scan() {
			options.WifiDict = append(options.WifiDict, scanner.Text())
		}

		// 检查读取过程中是否有错误发生
		if err := scanner.Err(); err != nil {
			panic(err)
		}

	}

	return nil
}

// LoadWifiDict 加载内置WiFi密码字典
func (options *Options) LoadWifiDict(dictContent string) []string {
	// 根据操作系统的换行符进行分割
	Dict := strings.Split(strings.ReplaceAll(dictContent, "\r\n", "\n"), "\n")
	gologger.Info().Msgf("字典个数：%v", len(Dict))
	return Dict
}
