package mailbox

import (
	. "github.com/10billionhongbaos/lib/message"
	"testing"
)

const (
	randomid = 12345
)

func Test_mailbox(t *testing.T) {
	a := MailBox{}
	a.Init(1000)
	for i := 0; i < 1000; i++ {
		var x Message
		b := a.Push(&x)
		if b != true {
			t.Error("wrong")
		}
	}

	var x Message
	b := a.Push(&x)
	if b != false {
		t.Error("wrong")
	}

	for i := 0; i < 1000; i++ {
		_, b := a.Pop()
		if b != true {
			t.Error("wrong")
		}
	}

	_, b2 := a.Pop()
	if b2 != false {
		t.Error("wrong")
	}

}
