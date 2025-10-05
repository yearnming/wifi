package wifi

import (
	"encoding/hex"
	"fmt"
	"sync"
	"testing"

	"github.com/projectdiscovery/gologger"
	"github.com/yearnming/wifi/pkg/util"
)

func TestGetWIFINetwork(t *testing.T) {
	networks, err := GetWIFINetworks()
	if err != nil {
		t.Fatal(err)
		// gologger.Error().Msgf(err.Error())
	}

	for _, network := range networks {
		// t.Logf("SSID: %s, BSSID: %s\n", network.SSID, network.BSSID)
		gologger.Info().Msgf("信号：%d%%, SSID: %s, BSSID: %s\n", network.Signal, network.SSID, network.BSSID)
	}
}

func TestEncoding(t *testing.T) {
	// 示例数据，全部用原始字节
	gb18030Example := []byte{0xB2, 0xE2, 0xBA, 0xFE, 0xB2, 0xBB, 0xCA, 0xD0, 0xC3, 0xFB, 0xBC, 0xE4}                // 浠婂ぉ鏃╃偣涓嬬彮
	hexUTF8Example, _ := hex.DecodeString("313233E4B8ADE5BF83E69D91E4B883E5B7B732")                                 // 123中心村七巷2
	gb18030MixExample := []byte{0x66, 0x67, 0x66, 0x64, 0xCA, 0xD0, 0xB4, 0xFA, 0x6F, 0x76, 0x61, 0x20, 0x31, 0x33} // fgfd鐨刵ova 13

	tests := []struct {
		name string
		data []byte
	}{
		{"GB18030 中文", gb18030Example},
		{"UTF-8 十六进制", hexUTF8Example},
		{"GB18030 混合文本", gb18030MixExample},
	}

	for _, tt := range tests {
		fmt.Println("======================================")
		fmt.Println("测试:", tt.name)

		// 重置检测
		util.DetectOnce = sync.Once{}
		util.DetectSystemEncoding(tt.data)

		result, err := util.DetectAndConvertEncoding(tt.data)
		if err != nil {
			fmt.Println("转换失败:", err)
		} else {
			fmt.Println("转换结果:", result)
		}
	}
}
