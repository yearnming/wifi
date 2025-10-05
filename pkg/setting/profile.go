package setting

import "os"

// profile config
var (
	// to do 操作系统临时目录
	ProfileXMLPath = os.TempDir() + "temp.xml"
	ProfileTmpl    = `<WLANProfile xmlns="http://www.microsoft.com/networking/WLAN/profile/v1">
<name>%s</name>
<SSIDConfig>
	<SSID>
		<hex>%s</hex>
		<name>%s</name>
	</SSID>
</SSIDConfig>
<connectionType>ESS</connectionType>
<connectionMode>%s</connectionMode>
<MSM>
	<security>
		<authEncryption>
			<authentication>WPA2PSK</authentication>
			<encryption>AES</encryption>
			<useOneX>false</useOneX>
		</authEncryption>
		<sharedKey>
			<keyType>passPhrase</keyType>
			<protected>false</protected>
			<keyMaterial>%s</keyMaterial>
		</sharedKey>
	</security>
</MSM>
</WLANProfile>`
)
