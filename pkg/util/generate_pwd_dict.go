package util

import (
	"fmt"
	"github.com/yearnming/wifi/pkg/config"
	"log"
	"os"
)

var f *os.File

func generatePwdDict() {
	var err error
	f, err = os.OpenFile("./pwd_dict.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()
	generatePwd([]byte{})
	final()
}

var written int
var writeOnce = 10000
var msgs = make([]byte, 0, writeOnce)

func writeToFile(msg []byte) {
	msg = append(msg, '\n')
	msgs = append(msgs, msg...)

	written++
	if written%writeOnce == 0 {
		_, err := f.Write(msgs)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println(written)
		msgs = make([]byte, 0, writeOnce)
	}
}

func final() {
	_, err := f.Write(msgs)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(written)
}

func generatePwd(b []byte) {
	if len(b) >= config.PwdMinLen {
		writeToFile(b)
		if len(b) >= config.PwdMaxLen {
			return
		}
	}

	for _, v := range config.PwdCharDict {
		generatePwd(append(b, v))
	}
}
