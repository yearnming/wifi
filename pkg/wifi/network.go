package wifi

import (
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/yearnming/wifi/pkg/util"
)

// var detectedEncoding string
// var detectOnce sync.Once

// Network Network
type Network struct {
	SSID   string
	BSSID  string
	Signal int
}

// GetWIFINetworks GetWIFINetworks
func GetWIFINetworks() (networks []Network, err error) {
	// cmd := exec.Command("cmd", "/C", "chcp 65001 && netsh wlan show networks mode=Bssid")

	cmd := exec.Command("netsh", "wlan", "show", "networks", "mode=Bssid")
	output, err := cmd.Output()
	// fmt.Printf("cmd: %v \n output: %v", cmd, output)
	// fmt.Printf("cmd: %v\noutput:\n%s", cmd, output)
	if err != nil {
		return
	}

	networks = parseWiFiList(string(output))
	// fmt.Printf("cmd: %v \n output: %v", cmd, networks)
	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Signal > networks[j].Signal
	})

	return
}
func parseWiFiList(output string) []Network {

	lines := strings.Split(output, "\n")
	var nets []Network
	var cur Network

	// 2. 状态机：只有“SSID 行”出现时才算一条新记录开始
	inBlock := false

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		// SSID 行：SSID 1 : xxx
		if strings.HasPrefix(line, "SSID") && strings.Contains(line, ":") {
			// 如果前面已经扫描到一半，把半成品扔掉
			if inBlock && cur.SSID != "" {
				nets = append(nets, cur)
			}
			cur = Network{}
			inBlock = true
			// cur.SSID = strings.TrimSpace(line[strings.Index(line, ":")+1:])
			// 获取并解码 SSID
			// rawSSID := strings.TrimSpace(line[strings.Index(line, ":")+1:])
			// 新句：只转 SSID
			if raw := strings.TrimSpace(line[strings.Index(line, ":")+1:]); raw != "" {
				cur.SSID, _ = util.DetectAndConvertEncoding([]byte(raw)) // 可能 GBK → UTF-8
			}
			// cur.SSID = decodeSSID(rawSSID)
			continue
		}

		// BSSID 行
		if strings.HasPrefix(line, "BSSID") && strings.Contains(line, ":") {
			cur.BSSID = strings.TrimSpace(line[strings.Index(line, ":")+1:])
			continue
		}

		// 信号行
		// if strings.HasPrefix(line, "信号") && strings.Contains(line, ":") {
		// 	percStr := strings.TrimRight(strings.TrimSpace(line[strings.Index(line, ":")+1:]), "%")
		// 	if p, err := strconv.Atoi(percStr); err == nil {
		// 		cur.Signal = p
		// 	}
		// 	// 一条完整记录结束
		// 	if cur.SSID != "" {
		// 		nets = append(nets, cur)
		// 	}
		// 	inBlock = false
		// 	cur = Network{}
		// }
		// 使用方式
		// 匹配 信号
		if i := strings.Index(line, ":"); i != -1 &&
			strings.Contains(line[i:], "%") { // 后半段必须带 %
			// 提取数字
			val := strings.TrimSpace(line[i+1:])
			val = strings.TrimRight(val, "%")
			if p, err := strconv.Atoi(val); err == nil {
				cur.Signal = p
			}

			// 一条完整记录结束
			if cur.SSID != "" {
				nets = append(nets, cur)
			}
			inBlock = false
			cur = Network{}
		}

	}

	// 3. 文件结尾可能还有一条半成品
	if inBlock && cur.SSID != "" {
		nets = append(nets, cur)
	}
	return nets
}

// // DetectAndConvertEncoding 根据已检测的系统编码转换为 UTF-8
// func DetectAndConvertEncoding(input []byte) (string, error) {
// 	// 1. UTF-8 优先
// 	// gologger.Info().Msgf("长度：%d", len(input))
// 	if detectedEncoding == "UTF-8" || utf8.Valid(input) {
// 		// gologger.Info().Msgf("SSID：%s", string(input))
// 		if len(input)%2 == 0 && isHex(input) {
// 			if hexBytes, err := hex.DecodeString(string(input)); err == nil {
// 				// 解码后再看是不是 UTF-8，否则按 GBK
// 				if utf8.Valid(hexBytes) {
// 					// gologger.Info().Msgf("utf8解码后：%s", string(hexBytes))
// 					return string(hexBytes), nil
// 				}
// 				if gbkStr, err := simplifiedchinese.GB18030.NewDecoder().String(string(hexBytes)); err == nil {
// 					// gologger.Info().Msgf("GB18030解码后：%s", gbkStr)
// 					return gbkStr, nil
// 				}
// 			}
// 		}
// 		return string(input), nil
// 	}

// 	var enc encoding.Encoding

// 	switch detectedEncoding {
// 	case "GB18030", "GBK", "GB2312":
// 		enc = simplifiedchinese.GB18030
// 	case "UTF16LE":
// 		enc = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
// 	case "UTF16BE":
// 		enc = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
// 	default:
// 		fmt.Printf("[警告] 未知编码 %s，回退到 GB18030\n", detectedEncoding)
// 		enc = simplifiedchinese.GB18030
// 	}

// 	reader := transform.NewReader(bytes.NewReader(input), enc.NewDecoder())
// 	decoded, err := io.ReadAll(reader)
// 	if err != nil {
// 		return "", fmt.Errorf("转换失败: %v", err)
// 	}
// 	if len(decoded)%2 == 0 && isHex(decoded) {
// 		if hexBytes, err := hex.DecodeString(string(decoded)); err == nil {
// 			// 解码后再看是不是 UTF-8，否则按 GBK
// 			if utf8.Valid(hexBytes) {
// 				// gologger.Info().Msgf("解码后：%s", string(hexBytes))
// 				return string(hexBytes), nil
// 			}
// 			if gbkStr, err := simplifiedchinese.GB18030.NewDecoder().String(string(hexBytes)); err == nil {
// 				// gologger.Info().Msgf("解码后：%s", gbkStr)
// 				return gbkStr, nil
// 			}
// 		}
// 	}
// 	// gologger.Info().Msgf("解码后：%s", string(decoded))

// 	return string(decoded), nil
// }

// // DetectSystemEncoding 使用 chardet 来检测系统默认编码，只检测一次
// func DetectSystemEncoding(sample []byte) {
// 	detectOnce.Do(func() {
// 		// 1. 优先尝试 UTF-8
// 		if utf8.Valid(sample) {
// 			detectedEncoding = "UTF-8"
// 			fmt.Println("[编码检测] 系统默认编码:", detectedEncoding)
// 			return
// 		}

// 		// 2. 如果不是 UTF-8，使用 chardet 检测最可能的编码
// 		detector := chardet.NewTextDetector()
// 		result, err := detector.DetectBest(sample)
// 		if err != nil {
// 			fmt.Println("[警告] 自动检测失败，使用默认 GB18030:", err)
// 			detectedEncoding = "GB18030"
// 			return
// 		}

// 		encName := strings.ToUpper(strings.ReplaceAll(result.Charset, "-", ""))
// 		fmt.Printf("[编码检测] 检测到原始编码: %s (置信度: %d%%)\n", result.Charset, result.Confidence)

// 		// 3. 合法性检查与归一化
// 		switch encName {
// 		case "UTF8":
// 			detectedEncoding = "UTF-8"
// 		case "GB18030", "GBK", "GB2312":
// 			detectedEncoding = "GB18030"
// 		case "UTF16LE":
// 			detectedEncoding = "UTF16LE"
// 		case "UTF16BE":
// 			detectedEncoding = "UTF16BE"
// 		default:
// 			fmt.Printf("[警告] 未知编码: %s，回退到 GB18030\n", encName)
// 			detectedEncoding = "GB18030"
// 		}
// 	})
// }

// // 判断是否为“真正的”十六进制字符串（必须含 A-F/a-F）
// func isHex(s []byte) bool {
// 	hasLetter := false
// 	for _, c := range s {
// 		switch {
// 		case c >= '0' && c <= '9':
// 			continue
// 		case c >= 'A' && c <= 'F', c >= 'a' && c <= 'f':
// 			hasLetter = true
// 		default:
// 			return false // 出现非法字符
// 		}
// 	}
// 	return hasLetter && len(s)%2 == 0 // 长度偶数且至少一个字母
// }
