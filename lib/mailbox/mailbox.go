package mailbox

import (
	. "github.com/10billionhongbaos/lib/message"
	. "github.com/10billionhongbaos/lib/utils"

	"net"
)

type MailBox struct {
	mailbox chan *Message
	Net     []*net.Conn
	index   uint32
}

type MessageProcess struct {
	InputBox MailBox
	YaoBox   MailBox
}

func (p *MailBox) Init(i int) {
	p.mailbox = make(chan *Message, i)
	p.Net = make([]*net.Conn, i)
}

func (p *MailBox) Push(p1 *Message) bool {

	select {
	case p.mailbox <- p1:
		return true
	default:
		return false
	}

}

func (p *MailBox) Pop() (*Message, bool) {

	var p1 *Message
	select {
	case p1 := <-p.mailbox:
		return p1, true
	default:
		return p1, false
	}

}

func (p *MailBox) BlockingPop() (*Message, bool) {

	var p1 *Message
	p1 = <-p.mailbox
	return p1, true

}

func (p *MessageProcess) Init(i1 int, i2 int) {

	p.InputBox.Init(i1)
	p.YaoBox.Init(i2)

}

func (p *MessageProcess) SetConn(p1 *net.Conn, mailboxnum int, userid int) {
	_, y := GetMailBoxAddress(mailboxnum, userid)
	p.InputBox.Net[y] = p1
}

func (p *MessageProcess) GetConn(mailboxnum int, userid int) *net.Conn {

	_, y := GetMailBoxAddress(mailboxnum, userid)
	return p.InputBox.Net[y]
}

func (p *MessageProcess) PushInput(p1 *Message) bool {

	return p.InputBox.Push(p1)
}

func (p *MessageProcess) PushYao(p1 *Message) bool {
	return p.YaoBox.Push(p1)

}

func (p *MessageProcess) PopInput() (*Message, bool) {
	return p.InputBox.BlockingPop()
}

func (p *MessageProcess) PopYao() (*Message, bool) {
	return p.YaoBox.Pop()
}
