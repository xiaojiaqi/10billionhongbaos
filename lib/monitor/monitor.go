package monitor

import (
	"encoding/binary"
	"fmt"
	. "github.com/10billionhongbaos/lib/utils"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	Request    = 10 // 发送请求总数目
	Response   = 11 // 收到响应数目
	GetHongbao = 14 // 获得红包的数目
	NoHongbao  = 15 // 没有获得红包
	AllHongbao = 16 // 已经获取了所有红包
	OverLoad   = 17 // 过载保护

	SendHongbao = 18

	Connection = 20 // 累计连接
	Discon     = 21 // 累计断开

	CountYaoHongbaoReq     = 22 //累计发放摇红包数目
	CountGetYahoHongbaoReq = 23 //累计收到摇红包数目

	SumYaoHongbao        = 24   //累计发放摇红包金额数目
	SumGetYaohoHongbao   = 25   //累计收到摇红包金额数目
	CountFaHongbaoReq    = 30   // 累计普通红包数目
	CountGetFaHongbaoReq = 31   // 累计获取普通红包数目
	SumFaHongbao         = 32   // 累计普通红包金额总数
	SumGetFaHongbao      = 33   // 累计获取普通红包金额总数
	TimeRange            = 6000 //只记录600秒内的结果
	ArrayLength          = 60   // 最多纪录项目
)

const (
	UDP_PORT = 8002
)

type Monitor struct {
	data            [ArrayLength][TimeRange]int64
	PersistenceData [ArrayLength]int64
}

var GMonitor Monitor

func (m *Monitor) Add(index int, i int64, timeStamp int64) int64 {
	v := timeStamp % TimeRange
	return atomic.AddInt64(&m.data[index][v], i)
}

func (m *Monitor) Add2(index int, i int64) int64 {
	return atomic.AddInt64(&m.PersistenceData[index], i)
}

func (m *Monitor) Get(timeStamp int64) ([ArrayLength]int64, [ArrayLength]int64) {
	data := [ArrayLength]int64{}
	pdata := [ArrayLength]int64{}
	v := timeStamp % TimeRange
	for i := 0; i < ArrayLength; i++ {
		data[i] = atomic.LoadInt64(&m.data[i][v])
		pdata[i] = atomic.LoadInt64(&m.PersistenceData[i])

	}
	return data, pdata
}

func (m *Monitor) Clean(timeStamp int64) {
	v := timeStamp % TimeRange
	for i := 0; i < ArrayLength; i++ {
		atomic.StoreInt64(&m.data[i][v], 0)
	}
}
func itoa(x int64) string {

	return strconv.Itoa(int(x))

}

func PrintData(timeStamp int64, types int, v [ArrayLength]int64, p [ArrayLength]int64) string {
	s := ""

	s += itoa(timeStamp)
	s += " " + itoa(int64(types))
	s += " Request: " + itoa(v[Request])
	s += " Response: " + itoa(v[Response])
	s += " GetHongbao: " + itoa(v[GetHongbao])
	s += " NoHongbao: " + itoa(v[NoHongbao])
	s += " AllHongbao: " + itoa(v[AllHongbao])
	s += " OverLoad: " + itoa(v[OverLoad])
	s += " yaohongbao: "
	s += itoa(p[CountYaoHongbaoReq]) + "/" + itoa(p[CountGetYahoHongbaoReq])
	s += " " + itoa(p[SumYaoHongbao]) + " " + itoa(p[SumGetYaohoHongbao])
	s += " fahongbao:"
	s += itoa(p[CountFaHongbaoReq]) + "/" + itoa(p[CountGetFaHongbaoReq])
	s += " " + itoa(p[SumFaHongbao]) + " " + itoa(p[SumGetFaHongbao])
	s += " detail "
	for i := 0; i <= 6; i++ {
		s += itoa(v[i]) + " "
	}
	return s
}

func (m *Monitor) ShowStatus(timeStamp int64, types int) ([ArrayLength]int64, [ArrayLength]int64) {
	v, p := m.Get(timeStamp)
	s := PrintData(timeStamp, types, v, p)
	fmt.Println(s)
	return v, p
}

func Show(host string, types int) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buff := make([]byte, ArrayLength*8*2+8+8)
	tlast := time.Now().Unix() - 4
	for {
		tnow := time.Now()
		t := tnow.Unix() - 3

		//  fmt.Println("tlast", tlast, "t", t)
		for ; tlast < t; tlast += 1 {
			//     fmt.Println("for tlast", tlast, "t", t)
			v, p := GMonitor.ShowStatus(tlast, types)
			var tmp uint64
			tmp = uint64(tlast)
			binary.BigEndian.PutUint64(buff, tmp)
			tmp = uint64(types)
			binary.BigEndian.PutUint64(buff[8:], tmp)
			index := 16
			for i := 0; i < ArrayLength; i++ {
				tmp = uint64(v[i])
				binary.BigEndian.PutUint64(buff[index:], tmp)
				index += 8
			}
			for i := 0; i < ArrayLength; i++ {
				tmp = uint64(p[i])
				binary.BigEndian.PutUint64(buff[index:], tmp)
				index += 8
			}

			_, err = conn.Write(buff)
			if err != nil {
				fmt.Println("发送数据失败!", err)
				//return
			}
			GMonitor.Clean(tlast)
		}
		CorrectSleepOneSecond2(tnow)
	}
}
