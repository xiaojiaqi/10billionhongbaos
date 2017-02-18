package message

import (
	//        "fmt"
	"encoding/binary"
	"net"
	"sync"
	"time"
)

const (
	HeaderLength = 12 // 2 字节长度＋2 字节 类型 + 8 字节时间戳
)

//  用户登录
type LoginRequest struct {
	Userid uint32
}

func (p *LoginRequest) GetLength() int {
	return 4
}

//  用户登录回应
type LoginResponse struct {
	Userid uint32
	Result uint32
}

func (p *LoginResponse) GetLength() int {
	return 4 + 4
}

//  摇红包
type YaoHongbaoRequest struct {
	Userid uint32
}

func (p *YaoHongbaoRequest) GetLength() int {
	return 4
}

//  摇红包回应
type YaoHongbaoResponse struct {
	Userid uint32
	Result uint32 // 00 没摇dao
	// 01 摇到了
	// 02 系统过载
	Hongbaoid uint32
	Amount    uint32
}

func (p *YaoHongbaoResponse) GetLength() int {
	return 4 * 4
}

//  发红包
type FaHongbaoRequest struct {
	Hongbaoid  uint32
	Amount     uint32
	Useramount uint32
	Userid1    uint32
	Userid2    uint32
	Userid3    uint32
}

func (p *FaHongbaoRequest) GetLength() int {
	return 4 * 6
}

//  拿红包
type GetHongbaoRequest struct {
	Userid    uint32
	Hongbaoid uint32
}

func (p *GetHongbaoRequest) GetLength() int {
	return 4 * 2
}

//  拿红包回应
type GetHongbaoRespone struct {
	Hongbaoid uint32
	Amount    uint32
}

func (p *GetHongbaoRespone) GetLength() int {
	return 4 * 2
}

// 还应该有聊天什么的
// ...
const (
	LoginRequest_ID       = 0x0000
	LoginResponse_ID      = 0x0001
	YaoHongbaoRequest_ID  = 0x0002
	YaoHongbaoResponse_ID = 0x0003
	FaHongbaoRequest_ID   = 0x0004
	GetHongbaoRequest_ID  = 0x0005
	GetHongbaoRespone_ID  = 0x0006
)

type Message struct {
	Userid    uint
	Cmd       int
	TimeStamp uint64
	Conn      *net.Conn
	P1        *LoginRequest
	P2        *LoginResponse
	P3        *YaoHongbaoRequest
	P4        *YaoHongbaoResponse
	P5        *FaHongbaoRequest
	P6        *GetHongbaoRequest
	P7        *GetHongbaoRespone
	Buff      []byte
	Bufflen   int
	Pool      *MessagePool
}

func (p *Message) Init() {
	p.Buff = make([]byte, 100)
	p.Bufflen = 0
}

func (p *Message) SetCmd(i int) {
	p.Cmd = i
}

/*
 Message format
 0-1   length
 2-3   cmdtype
*/

func (p *LoginRequest) serialization(buff []byte) {
	binary.BigEndian.PutUint32(buff, p.Userid)
	//fmt.Println(p.Userid, buff)
}

func (p *LoginRequest) unserialization(buff []byte) {

	p.Userid = binary.BigEndian.Uint32(buff)
	//fmt.Println(p.Userid, buff)
}

func (p *LoginResponse) serialization(buff []byte) {
	binary.BigEndian.PutUint32(buff, p.Userid)
	binary.BigEndian.PutUint32(buff[4:], p.Result)
	//  fmt.Println(p.Userid, p.Result,buff)
}

func (p *LoginResponse) unserialization(buff []byte) {
	p.Userid = binary.BigEndian.Uint32(buff)
	p.Result = binary.BigEndian.Uint32(buff[4:])
	//    fmt.Println(p.Userid, p.Result,buff)
}

func (p *YaoHongbaoRequest) serialization(buff []byte) {

	binary.BigEndian.PutUint32(buff, p.Userid)
	//fmt.Println(p.Userid, buff)
}

func (p *YaoHongbaoRequest) unserialization(buff []byte) {

	p.Userid = binary.BigEndian.Uint32(buff)
	//fmt.Println(p.Userid, buff)
}

func (p *YaoHongbaoResponse) serialization(buff []byte) {

	binary.BigEndian.PutUint32(buff, p.Userid)
	binary.BigEndian.PutUint32(buff[4:], p.Result)

	binary.BigEndian.PutUint32(buff[8:], p.Hongbaoid)

	binary.BigEndian.PutUint32(buff[12:], p.Amount)

	//fmt.Println(p.Userid, buff)
}

func (p *YaoHongbaoResponse) unserialization(buff []byte) {

	p.Userid = binary.BigEndian.Uint32(buff)
	p.Result = binary.BigEndian.Uint32(buff[4:])

	p.Hongbaoid = binary.BigEndian.Uint32(buff[8:])
	p.Amount = binary.BigEndian.Uint32(buff[12:])

	//fmt.Println(p.Userid, buff)
}

func (p *FaHongbaoRequest) serialization(buff []byte) {

	binary.BigEndian.PutUint32(buff, p.Hongbaoid)
	binary.BigEndian.PutUint32(buff[4:], p.Amount)
	binary.BigEndian.PutUint32(buff[8:], p.Useramount)
	binary.BigEndian.PutUint32(buff[12:], p.Userid1)
	binary.BigEndian.PutUint32(buff[16:], p.Userid2)
	binary.BigEndian.PutUint32(buff[20:], p.Userid3)

	//fmt.Println(p.Userid, buff)
}

func (p *FaHongbaoRequest) unserialization(buff []byte) {

	p.Hongbaoid = binary.BigEndian.Uint32(buff)
	p.Amount = binary.BigEndian.Uint32(buff[4:])
	p.Useramount = binary.BigEndian.Uint32(buff[8:])
	p.Userid1 = binary.BigEndian.Uint32(buff[12:])
	p.Userid2 = binary.BigEndian.Uint32(buff[16:])
	p.Userid3 = binary.BigEndian.Uint32(buff[20:])
	//fmt.Println(p.Userid, buff)
}

func (p *GetHongbaoRequest) serialization(buff []byte) {

	binary.BigEndian.PutUint32(buff, p.Userid)
	binary.BigEndian.PutUint32(buff[4:], p.Hongbaoid)
}

func (p *GetHongbaoRequest) unserialization(buff []byte) {

	p.Userid = binary.BigEndian.Uint32(buff)
	p.Hongbaoid = binary.BigEndian.Uint32(buff[4:])
}

func (p *GetHongbaoRespone) serialization(buff []byte) {

	binary.BigEndian.PutUint32(buff, p.Hongbaoid)
	binary.BigEndian.PutUint32(buff[4:], p.Amount)
}

func (p *GetHongbaoRespone) unserialization(buff []byte) {

	p.Hongbaoid = binary.BigEndian.Uint32(buff)
	p.Amount = binary.BigEndian.Uint32(buff[4:])
}

func (p *Message) Serialization() []byte {
	len := HeaderLength
	var buff []byte
	switch p.Cmd {

	case LoginRequest_ID:
		{
			len += p.P1.GetLength()
			buff = p.Buff[0:len]
			p.Bufflen = len
			p.P1.serialization(buff[HeaderLength:])
		}

	case LoginResponse_ID:
		{
			len += p.P2.GetLength()
			//buff = make([]byte, len)
			buff = p.Buff[0:len]

			p.Bufflen = len
			p.P2.serialization(buff[HeaderLength:])

		}

	case YaoHongbaoRequest_ID:
		{

			len += p.P3.GetLength()
			//buff = make([]byte, len)
			buff = p.Buff[0:len]

			p.Bufflen = len
			p.P3.serialization(buff[HeaderLength:])
		}

	case YaoHongbaoResponse_ID:
		{
			len += p.P4.GetLength()
			//		buff = make([]byte, len)
			buff = p.Buff[0:len]

			p.Bufflen = len

			p.P4.serialization(buff[HeaderLength:])

		}

	case FaHongbaoRequest_ID:
		{

			len += p.P5.GetLength()
			//	buff = make([]byte, len)
			buff = p.Buff[0:len]

			p.Bufflen = len
			p.P5.serialization(buff[HeaderLength:])

		}

	case GetHongbaoRequest_ID:
		{
			len += p.P6.GetLength()

			//buff = make([]byte, len)
			buff = p.Buff[0:len]

			p.Bufflen = len
			p.P6.serialization(buff[HeaderLength:])

		}

	case GetHongbaoRespone_ID:
		{
			len += p.P7.GetLength()
			//buff = make([]byte, len)
			buff = p.Buff[0:len]

			p.Bufflen = len
			p.P7.serialization(buff[HeaderLength:])

		}

	}

	binary.BigEndian.PutUint16(buff[2:], uint16(p.Cmd))
	binary.BigEndian.PutUint16(buff, uint16(len))
	binary.BigEndian.PutUint64(buff[4:], uint64(time.Now().Unix()))

	return buff

}

func (p *Message) Unserialization(buff []byte) bool {
	p.Cmd = int(binary.BigEndian.Uint16(buff[2:]))
	p.TimeStamp = uint64(binary.BigEndian.Uint64(buff[4:]))
	switch p.Cmd {

	case LoginRequest_ID:
		{
			if p.P1 == nil {
				p.P1 = &LoginRequest{}
			}
			p.P1.unserialization(buff[HeaderLength:])

		}

	case LoginResponse_ID:
		{
			if p.P2 == nil {
				p.P2 = &LoginResponse{}
			}
			p.P2.unserialization(buff[HeaderLength:])

		}

	case YaoHongbaoRequest_ID:
		{
			if p.P3 == nil {
				p.P3 = &YaoHongbaoRequest{}
			}
			p.P3.unserialization(buff[HeaderLength:])
		}

	case YaoHongbaoResponse_ID:
		{
			if p.P4 == nil {
				p.P4 = &YaoHongbaoResponse{}
			}
			p.P4.unserialization(buff[HeaderLength:])
		}

	case FaHongbaoRequest_ID:
		{
			if p.P5 == nil {
				p.P5 = &FaHongbaoRequest{}
			}
			p.P5.unserialization(buff[HeaderLength:])
		}
	case GetHongbaoRequest_ID:
		{
			if p.P6 == nil {
				p.P6 = &GetHongbaoRequest{}
			}
			p.P6.unserialization(buff[HeaderLength:])
		}
	case GetHongbaoRespone_ID:
		{
			if p.P7 == nil {
				p.P7 = &GetHongbaoRespone{}
			}
			p.P7.unserialization(buff[HeaderLength:])
		}
	default:
		return false

	}
	return true
}

func (p *Message) Rollback() {
	p.Pool.Push(p)
}

type MessagePool struct {
	l       sync.Mutex
	maxuser int
	user    int
	queue   chan *Message
}

func (p *MessagePool) Init(len int, maxusernum int) {
	p.maxuser = maxusernum
	p.queue = make(chan *Message, len)

}

func (p *MessagePool) GetPool() *MessagePool {
	p.l.Lock()
	if p.user+1 < p.maxuser {
		p.user += 1
		p.l.Unlock()
		return p
	} else {
		p.l.Unlock()
		return nil
	}

}

func (p *MessagePool) Reduce() {
	p.l.Lock()
	p.user -= 1
	p.l.Unlock()

}

func (p *MessagePool) Pop() *Message {

	p2 := <-p.queue
	return p2
}

func (p *MessagePool) Push(pmessage *Message) {

	p.queue <- pmessage
}
