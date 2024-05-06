package main

import (
	"bufio"
	"fmt"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/yearnming/wifi/pkg/setting"
	"github.com/yearnming/wifi/pkg/util"
	"github.com/yearnming/wifi/pkg/wifi"
	"log"
	"os"
	"sync"
	"testing"
)

func NewWifi(ssid, pwd string) {
	wc := wifi.New(ssid, pwd)

	stat, err := wc.Connect()
	if stat == wifi.Connected {
		log.Println("Connect WiFi success:", ssid, pwd)
		err = util.WriteToFile(setting.SuccessPwdSavePath, fmt.Sprintf("WiFi: %s, password: %s\n", wc.Ssid, pwd))
		if err != nil {
			log.Println("write password to file err", err)
		}
		return
	}
	log.Println("Connect WiFi failed")
	err = wc.DeleteProfile()
	if err != nil {
		log.Println("Delete profile failed:", err)
		return
	}
	log.Println("Delete profile success")

	if err != nil {
		log.Println("Connect WiFi failed:", err)
		return
	}
}

func TestName(t *testing.T) {

	ssid := "yearn"

	//pwd := "qpqp1010.."
	//dict := []string{"88887777", "qpqp1010..", "12345678"}

	filename := "wifi.txt"
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建一个扫描器来读取文件
	scanner := bufio.NewScanner(file)

	// 读取所有行到一个切片中
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 检查读取过程中是否有错误发生
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	//var now = time.Now()
	//NewWifi(ssid, pwd)

	// 使用WaitGroup等待所有goroutine完成
	var wg sync.WaitGroup
	for _, line := range lines {
		wg.Add(1) // 增加WaitGroup的计数
		go func(line string) {
			defer wg.Done() // 减少WaitGroup的计数
			NewWifi(ssid, line)
		}(line)
	}

	// 等待所有goroutine完成
	wg.Wait()
}

func TestTime(t *testing.T) {
	// 设置日志级别为Debug，这样所有级别的日志都会输出
	//gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)

	// 启用时间戳，并设置最小日志级别为Info，这样Info及以上级别的日志都会带有时间戳
	gologger.DefaultLogger.SetTimestamp(true, levels.LevelInfo)
	gologger.DefaultLogger.SetTimestamp(true, levels.LevelError)
	// 打印一条日志信息
	gologger.Info().Msg("asdasdasd")
	// 打印一条带有时间戳的日志信息
	gologger.Error().TimeStamp().Msg("This is an info level log with timestamp")
}
