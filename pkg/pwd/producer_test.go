package pwd

import (
	"fmt"
	"testing"
)

func TestNewProducer(t *testing.T) {
	pwdChan := NewProducer(1, 2, []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'})

	for n := range pwdChan {
		fmt.Println(n)
	}
}
