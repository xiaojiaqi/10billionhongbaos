package main

import (
	"encoding/binary"
	"fmt"
	. "github.com/10billionhongbaos/lib/monitor"
	"github.com/10billionhongbaos/lib/parse"
	. "github.com/10billionhongbaos/lib/utils"
	"net"
	"os"
	"time"
)

var client1 Monitor
var client2 Monitor

var server1 Monitor
var server2 Monitor

func show() {

	ftime := time.Now().Format("2006-01-02_15:04:05")

	f, err := os.Create("log" + ftime + ".txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := ""
	for {
		tnow := time.Now()
		t := tnow.Unix() - 30

		{
			v1, p1 := client1.Get(t)
			v2, p2 := client2.Get(t)

			for i := 0; i < ArrayLength; i++ {
				v1[i] += v2[i]
				p1[i] += p2[i]
			}
			s = PrintData(t, 0, v1, v1)
			fmt.Println(s)

			s += "\n"
			f.Write([]byte(s))
		}

		{
			v1, p1 := server1.Get(t)
			v2, p2 := server2.Get(t)

			for i := 0; i < ArrayLength; i++ {
				v1[i] += v2[i]
				p1[i] += p2[i]
			}
			s = PrintData(t, 1, v1, v1)

			fmt.Println(s)

			s += "\n"
			f.Write([]byte(s))
		}
		client1.Clean(t)
		client2.Clean(t)
		server1.Clean(t)
		server2.Clean(t)
		CorrectSleepOneSecond2(tnow)
	}
}

func PareseData(buf []byte) {

	var times uint64
	var types uint64
	var tmp uint64

	var p1 *Monitor
	var p2 *Monitor

	times = binary.BigEndian.Uint64(buf)
	types = binary.BigEndian.Uint64(buf[8:])
	//    fmt.Println(times, times)

	//	dataindex := int(times) % ArrayLength

	if types == 0 {

		p1 = &client1
		p2 = &client1
	} else {

		p1 = &server1
		p2 = &server1
	}

	index := 16

	for i := 0; i < ArrayLength; i++ {
		tmp = binary.BigEndian.Uint64(buf[index:])
		p1.Add(i, int64(tmp), int64(times))
		index += 8
	}

	for i := 0; i < ArrayLength; i++ {
		tmp = binary.BigEndian.Uint64(buf[index:])
		p2.Add(i, int64(tmp), int64(times))
		index += 8

	}
}

func main() {
	netListen, err := net.Listen("tcp", "0.0.0.0:8002")
	CheckErr(err)
	defer netListen.Close()

	go show()
	//
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		//  Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(&conn)
	}

}

func handleConnection(conn *net.Conn) {
	defer (*conn).Close()

	data := make([]byte, 65536)

	for {

		if !parse.ReadByte(*conn, data, (ArrayLength*8*2 + 8 + 8)) {
			return

		}
		PareseData(data)
	}

}
