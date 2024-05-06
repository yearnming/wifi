package pwd

// Producer Producer
type Producer struct {
	pwdChan  chan string
	min, max int
	dict     []byte
}

// NewProducer NewProducer
func NewProducer(min, max int, dict []byte) chan string {
	p := Producer{make(chan string, 3), min, max, dict}
	go func() {
		p.generatePwd([]byte{})
		close(p.pwdChan)
	}()
	return p.pwdChan
}

func (p Producer) generatePwd(b []byte) {
	if len(b) >= p.min {
		p.pwdChan <- string(b)
		if len(b) >= p.max {
			return
		}

	}
	for _, v := range p.dict {
		p.generatePwd(append(b, v))
	}
}
