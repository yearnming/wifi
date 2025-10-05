package wifi

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/yearnming/wifi/pkg/logger"
	"github.com/yearnming/wifi/pkg/setting"
)

// 默认超时：30 秒可覆盖
var DefaultConnectTimeout = 30 * time.Second

// WIFI WIFI
type WIFI struct {
	Ssid, Password string
}

// New New
func New(ssid, password string) WIFI {
	return WIFI{
		Ssid:     ssid,
		Password: password,
	}
}

// Connect 带超时/防卡死
func (wc *WIFI) Connect() (Stat, error) {
	// 1. 整体超时控制
	ctx, cancel := context.WithTimeout(context.Background(), DefaultConnectTimeout)
	defer cancel()

	// 2. 加 profile
	if err := wc.addProfile(); err != nil {
		return "", err
	}

	// 3. 启动连接（带 ctx）
	if err := wc.connectCtx(ctx); err != nil {
		return "", err
	}

	// 4. 轮询状态（同样监听 ctx）
	tick := time.NewTicker(setting.CheckStatDuration)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("连接超时: %w", ctx.Err())
		case <-tick.C:
			stat, err := GetWIFIStat()
			if err != nil {
				return "", err
			}
			if stat == Disconnected || stat == Connected {
				return stat, nil
			}
		}
	}
}

// connectCtx 用 ctx 控制 netsh 生命周期
func (wc *WIFI) connectCtx(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "netsh", "wlan", "connect", "name="+wc.Ssid)

	// 保险丝：如果 ctx 超时但 cmd 还没返回，直接 Kill
	killer := time.AfterFunc(DefaultConnectTimeout+2*time.Second, func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})
	defer killer.Stop()

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Info(out)
		return fmt.Errorf("netsh connect: %w", err)
	}
	return nil
}

// DeleteProfile Delete Profile
func (wc *WIFI) DeleteProfile() error {
	cmd := exec.Command("netsh", "wlan", "delete", "profile", "name="+wc.Ssid)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Info(out)
	}
	return err
}

// C:\ProgramData\Microsoft\Wlansvc\Profiles\Interfaces
func (wc *WIFI) addProfile() error {
	profileXML := fmt.Sprintf(setting.ProfileTmpl, wc.Ssid, hex.EncodeToString([]byte(wc.Ssid)), wc.Ssid, Manual, wc.Password)

	// gologger.Info().Msgf("[%s]hex:[%s]", wc.Ssid, hex.EncodeToString([]byte(wc.Ssid)))
	//fmt.Println("配置文件: ", profileXML)
	profileXMLPath := setting.ProfileXMLPath

	if err := os.WriteFile(profileXMLPath, []byte(profileXML), 0644); err != nil {
		return err
	}
	defer os.Remove(profileXMLPath)
	// >>> 新增：打印即将导入的配置文件内容
	// if xmlBytes, e := os.ReadFile(profileXMLPath); e == nil {
	// 	fmt.Println("--- 待导入的 WiFi 配置文件内容 ---")
	// 	fmt.Println(string(xmlBytes))
	// 	fmt.Println("--- 配置文件内容结束 ---")
	// }
	// <<< 新增结束
	cmd := exec.Command("netsh", "wlan", "add", "profile", "filename="+profileXMLPath, "user=all")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Info(out)
	}

	// gologger.Info().Msgf("SSID: [%s],Password: [%s],Path: [%s]", wc.Ssid, wc.Password, profileXMLPath)
	return err
}
