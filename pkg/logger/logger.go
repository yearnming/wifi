package logger

import (
	"fmt"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// Info Info
func Info(data []byte) {
	val, err := simplifiedchinese.GB18030.NewDecoder().Bytes(data)
	if err != nil {
		fmt.Println(string(data))
		return
	}
	fmt.Println(string(val))
}
