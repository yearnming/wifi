package util

import (
	"os"
)

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
