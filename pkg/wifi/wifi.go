package wifi

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/yearnming/wifi/pkg/logger"
	"github.com/yearnming/wifi/pkg/setting"
)

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

// Connect Connect
func (wc *WIFI) Connect() (stat Stat, err error) {
	err = wc.addProfile()
	if err != nil {
		return
	}

	err = wc.connect()
	if err != nil {
		return
	}

	for {
		<-time.After(setting.CheckStatDuration)

		stat, err := GetWIFIStat()
		if err != nil {
			return "", err
		}

		if stat == Disconnected || stat == Connected {
			return stat, nil
		}
	}
}

func (wc *WIFI) connect() error {
	cmd := exec.Command("netsh", "wlan", "connect", "name="+wc.Ssid)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Info(out)
	}
	return err
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
	profileXML := fmt.Sprintf(setting.ProfileTmpl, wc.Ssid, wc.Ssid, Manual, wc.Password)

	//fmt.Println("配置文件: ", profileXML)
	profileXMLPath := setting.ProfileXMLPath

	if err := os.WriteFile(profileXMLPath, []byte(profileXML), 0644); err != nil {
		return err
	}
	defer os.Remove(profileXMLPath)

	cmd := exec.Command("netsh", "wlan", "add", "profile", "filename="+profileXMLPath, "user=all")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Info(out)
	}

	return err
}
