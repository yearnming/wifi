package util

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var DetectedEncoding string
var DetectOnce sync.Once

// WriteToFile WriteToFile
func WriteToFile(filename, data string) (err error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(data)
	return
}

// DetectAndConvertEncoding 根据已检测的系统编码转换为 UTF-8
func DetectAndConvertEncoding(input []byte) (string, error) {
	// 1. UTF-8 优先
	// gologger.Info().Msgf("长度：%d", len(input))
	if DetectedEncoding == "UTF-8" || utf8.Valid(input) {
		// gologger.Info().Msgf("SSID：%s", string(input))
		if len(input)%2 == 0 && isHex(input) {
			if hexBytes, err := hex.DecodeString(string(input)); err == nil {
				// 解码后再看是不是 UTF-8，否则按 GBK
				if utf8.Valid(hexBytes) {
					// gologger.Info().Msgf("utf8解码后：%s", string(hexBytes))
					return string(hexBytes), nil
				}
				if gbkStr, err := simplifiedchinese.GB18030.NewDecoder().String(string(hexBytes)); err == nil {
					// gologger.Info().Msgf("GB18030解码后：%s", gbkStr)
					return gbkStr, nil
				}
			}
		}
		return string(input), nil
	}

	var enc encoding.Encoding

	switch DetectedEncoding {
	case "GB18030", "GBK", "GB2312":
		enc = simplifiedchinese.GB18030
	case "UTF16LE":
		enc = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case "UTF16BE":
		enc = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	default:
		// fmt.Printf("[警告] 未知编码 %s，回退到 GB18030\n", DetectedEncoding)
		enc = simplifiedchinese.GB18030
	}

	reader := transform.NewReader(bytes.NewReader(input), enc.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("转换失败: %v", err)
	}
	if len(decoded)%2 == 0 && isHex(decoded) {
		if hexBytes, err := hex.DecodeString(string(decoded)); err == nil {
			// 解码后再看是不是 UTF-8，否则按 GBK
			if utf8.Valid(hexBytes) {
				// gologger.Info().Msgf("解码后：%s", string(hexBytes))
				return string(hexBytes), nil
			}
			if gbkStr, err := simplifiedchinese.GB18030.NewDecoder().String(string(hexBytes)); err == nil {
				// gologger.Info().Msgf("解码后：%s", gbkStr)
				return gbkStr, nil
			}
		}
	}
	// gologger.Info().Msgf("解码后：%s", string(decoded))

	return string(decoded), nil
}

// DetectSystemEncoding 使用 chardet 来检测系统默认编码，只检测一次
func DetectSystemEncoding(sample []byte) {
	DetectOnce.Do(func() {
		// 1. 优先尝试 UTF-8
		if utf8.Valid(sample) {
			DetectedEncoding = "UTF-8"
			fmt.Println("[编码检测] 系统默认编码:", DetectedEncoding)
			return
		}

		// 2. 如果不是 UTF-8，使用 chardet 检测最可能的编码
		detector := chardet.NewTextDetector()
		result, err := detector.DetectBest(sample)
		if err != nil {
			fmt.Println("[警告] 自动检测失败，使用默认 GB18030:", err)
			DetectedEncoding = "GB18030"
			return
		}

		encName := strings.ToUpper(strings.ReplaceAll(result.Charset, "-", ""))
		fmt.Printf("[编码检测] 检测到原始编码: %s (置信度: %d%%)\n", result.Charset, result.Confidence)

		// 3. 合法性检查与归一化
		switch encName {
		case "UTF8":
			DetectedEncoding = "UTF-8"
		case "GB18030", "GBK", "GB2312":
			DetectedEncoding = "GB18030"
		case "UTF16LE":
			DetectedEncoding = "UTF16LE"
		case "UTF16BE":
			DetectedEncoding = "UTF16BE"
		default:
			fmt.Printf("[警告] 未知编码: %s，回退到 GB18030\n", encName)
			DetectedEncoding = "GB18030"
		}
	})
}

// 判断是否为“真正的”十六进制字符串（必须含 A-F/a-F）
func isHex(s []byte) bool {
	hasLetter := false
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			continue
		case c >= 'A' && c <= 'F', c >= 'a' && c <= 'f':
			hasLetter = true
		default:
			return false // 出现非法字符
		}
	}
	return hasLetter && len(s)%2 == 0 // 长度偶数且至少一个字母
}

// parseIndexes 把 "1,3,5" 或 "1-3,5" 转成 []int
func ParseIndexes(s string, max int) ([]int, error) {
	s = strings.ReplaceAll(s, " ", "") // 去空格
	var out []int
	parts := strings.Split(s, ",")
	for _, p := range parts {
		if strings.Contains(p, "-") {
			// 区间
			sp := strings.Split(p, "-")
			if len(sp) != 2 {
				return nil, fmt.Errorf("区间格式错误: %s", p)
			}
			start, err1 := strconv.Atoi(sp[0])
			end, err2 := strconv.Atoi(sp[1])
			if err1 != nil || err2 != nil || start > end || start < 1 || end > max {
				return nil, fmt.Errorf("区间数值无效: %s", p)
			}
			for i := start; i <= end; i++ {
				out = append(out, i)
			}
		} else {
			// 单个
			i, err := strconv.Atoi(p)
			if err != nil || i < 1 || i > max {
				return nil, fmt.Errorf("编号无效: %s", p)
			}
			out = append(out, i)
		}
	}
	return out, nil
}

func DefaultFailDBPath() string {
	// home, _ := os.UserHomeDir()
	// return filepath.Join(home, ".wifi-crifier", "filedb.json")
	return "filedb.json"
}
