package main

import (
	"fmt"
	"github.com/10billionhongbaos/lib/flags"
	. "github.com/10billionhongbaos/lib/message"
	"github.com/10billionhongbaos/lib/monitor"
	"github.com/10billionhongbaos/lib/parse"
	"github.com/10billionhongbaos/lib/qps"
	"github.com/10billionhongbaos/lib/utils"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func readConn(conn *net.Conn, userid uint32, gPool *MessagePool) {

	defer (*conn).Close()
	for {
		p, result := parse.ReadMessage(conn, gPool)
		if !result {
			break
		}
		//fmt.Println(p.Cmd)
		tnow := time.Now().Unix()
		monitor.GMonitor.Add(p.Cmd, 1, tnow)

		monitor.GMonitor.Add(monitor.Response, 1, tnow)
		switch p.Cmd {
		case YaoHongbaoResponse_ID:
			{
				monitor.GMonitor.Add(int(p.P4.Result), 1, tnow)
				if int(p.P4.Result) == monitor.GetHongbao {
					monitor.GMonitor.Add2(monitor.CountYaoHongbaoReq, 1)
					monitor.GMonitor.Add2(monitor.SumGetYaohoHongbao, int64(p.P4.Amount))
				}
			}
		case FaHongbaoRequest_ID:
			{
				// send a  GetHongbaoRequest_ID
				// bug: 存着 竞争问题

				//        fmt.Println("got a fahonbao request", " userid",userid, p.P5.Hongbaoid, p.P5.Userid1, p.P5.Userid2,p.P5.Userid3)
				{
					a := gPool.Pop()
					a.Cmd = GetHongbaoRequest_ID
					p2 := &GetHongbaoRequest{}
					p2.Hongbaoid = p.P5.Hongbaoid
					a.P6 = p2

					parse.WriteMessage(conn, a)
					monitor.GMonitor.Add2(monitor.CountFaHongbaoReq, 1)
					monitor.GMonitor.Add(GetHongbaoRequest_ID, 1, time.Now().Unix())
				}
			}

		case GetHongbaoRespone_ID:
			{

				//        fmt.Println("got a GetHongbaoRespone_ID", userid, p.P7.Hongbaoid,  p.P7.Amount)
				monitor.GMonitor.Add2(monitor.CountGetFaHongbaoReq, 1)
				monitor.GMonitor.Add2(monitor.SumGetFaHongbao, int64(p.P7.Amount))
			}

		}
		p.Rollback()

	}

}

func sender(conn net.Conn, userid uint32, gPool *MessagePool) bool {

	{
		a := gPool.Pop()
		a.Userid = uint(userid)
		a.Cmd = YaoHongbaoRequest_ID
		p := YaoHongbaoRequest{}
		p.Userid = userid
		a.P3 = &p

		monitor.GMonitor.Add(monitor.Request, 1, time.Now().Unix())

		err := parse.WriteMessage(&conn, a)
		if err != nil {
			return false
		}

	}

	return true

}

func OneClient(userid uint32) {
	var gPool MessagePool

	gPool.Init(30, 30)
	for k := 0; k < 30; k++ {
		p := &Message{}
		p.Init()
		p.Pool = &gPool
		gPool.Push(p)
	}

	conn, err := net.Dial("tcp", *(flags.Server))
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer conn.Close()
	go readConn(&conn, userid, &gPool)

	utils.CorrectSleepOneSecond()

	// send a login request

	{
		a := gPool.Pop()
		a.Userid = uint(userid)
		a.Cmd = LoginRequest_ID
		p := LoginRequest{}
		p.Userid = userid
		a.P1 = &p

		monitor.GMonitor.Add(monitor.Request, 1, time.Now().Unix())

		err := parse.WriteMessage(&conn, a)
		if err != nil {
			return
		}
	}

	utils.CorrectSleepOneSecond()

	for {
		tnow := time.Now()
		if qps.TriggerRequest2(userid, flags.TotalUser, flags.QPS, time.Now().Unix()) {
			//	fmt.Println("time: ", time.Now().Unix(), " userid ", userid, " send ", conn.LocalAddr())

			//            time.Sleep(300* time.Microsecond)
			if !sender(conn, userid, &gPool) {
				return
			}
		}
		utils.CorrectSleepOneSecond2(tnow)
	}
}

func Setqps(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	fmt.Println(req.Form)
	filelist, found1 := req.Form["qps"]

	if !(found1) {
		fmt.Println("no id")
		w.Write([]byte("no"))
		return
	}

	nowQps, err := strconv.Atoi(filelist[0])
	if err != nil {
		return
	}

	flags.QPS = uint32(nowQps)
}

func httpServer() {

	http.HandleFunc("/qps", Setqps)
	//服务器要监听的主机地址和端口号
	err := http.ListenAndServe("0.0.0.0:9090", nil)
	utils.CheckErr(err)

}

func runOneClient(i uint32) {
	for {
		time.Sleep(time.Duration(int(i)%20) * time.Second)
		OneClient(i)
	}
}

func main() {
	flags.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	go monitor.Show(*(flags.MonitorServer), 0)
	go httpServer()
	var i uint32
	fmt.Println(flags.MinUserId, flags.MaxUserId, flags.TotalUser)
	for i = flags.MinUserId; i <= flags.MaxUserId; i++ {
		go runOneClient(i)
	}

	for {
		utils.CorrectSleepOneSecond()

	}

}
