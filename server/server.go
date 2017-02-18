package main

import (
	"fmt"
	"github.com/10billionhongbaos/lib/flags"
	. "github.com/10billionhongbaos/lib/mailbox"
	. "github.com/10billionhongbaos/lib/message"
	"github.com/10billionhongbaos/lib/monitor"
	"github.com/10billionhongbaos/lib/parse"
	"github.com/10billionhongbaos/lib/utils"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	//    "runtime/debug"
	"strconv"
	"time"
)

const (
	MaxUser                 = 1200000
	MailBoxNum              = 256
	MailBoxUserNum          = MaxUser / MailBoxNum
	MailBoxMessageLength    = MailBoxUserNum
	HongBaoBoxMessageLength = 120000
	MaxQPS                  = 65000
	HongBaoMailBoxNum       = 5
	GCCycle                 = 5
	MessagePoolNum          = 1024
	MessagePoolUser         = MaxUser / MessagePoolNum
	MessagePoolMessageNum   = MaxUser / MessagePoolNum * 3
)

var gMailbox []MessageProcess
var gHongbaoMailbox []MailBox
var gPool []MessagePool

func getOnePool(num int) *MessagePool {
	for {
		p := gPool[num%MessagePoolNum].GetPool()
		if p == nil {
			num += 1
		} else {
			return p
		}
	}
}

func InitProcess() {
	gPool = make([]MessagePool, (MessagePoolNum + 2))
	for i := 0; i < (MessagePoolNum + 2); i++ {
		num := MessagePoolMessageNum
		if i >= MessagePoolNum {
			num = MessagePoolMessageNum * 4 * 15
		}
		gPool[i].Init(num, MessagePoolUser)
		for k := 0; k < num; k++ {
			p := &Message{}
			p.Init()
			p.Pool = &gPool[i]
			gPool[i].Push(p)
		}
	}

	gHongbaoMailbox = make([]MailBox, HongBaoMailBoxNum)
	for i := 0; i < HongBaoMailBoxNum; i++ {
		gHongbaoMailbox[i].Init(120000)
	}
	gMailbox = make([]MessageProcess, MailBoxNum)
	for i := 0; i < MailBoxNum; i++ {
		gMailbox[i].Init(MailBoxMessageLength, HongBaoBoxMessageLength)
	}
	for i := 0; i < MailBoxNum; i++ {
		go ProcessRequest(i)
	}

}

func ProcessRequest(index int) {
	for {
		p1, _ := gMailbox[index].PopInput()

		switch p1.Cmd {
		case YaoHongbaoRequest_ID:
			{
				p2, b2 := gMailbox[index].PopYao()
				if b2 == true && p2 != nil {
					// send this request
					b2 := parse.WriteMessage(p1.Conn, p2)
					if b2 != nil {
						(*(p1.Conn)).Close()
					}

					monitor.GMonitor.Add(monitor.SendHongbao, 1, time.Now().Unix())

				} else {
					a := p1.Pool.Pop()
					a.Cmd = YaoHongbaoResponse_ID
					p := YaoHongbaoResponse{}
					p.Userid = p1.P3.Userid
					p.Result = monitor.NoHongbao
					p.Hongbaoid = 0
					a.P4 = &p

					b2 := parse.WriteMessage(p1.Conn, a)
					if b2 != nil {
						(*(p1.Conn)).Close()
					}
				}
			}
		case GetHongbaoRequest_ID:
			{
				index := p1.P6.Hongbaoid % HongBaoMailBoxNum
				gHongbaoMailbox[index].Push(p1)
			}
		}
		p1.Rollback()

	}
}

func main() {

	//建立socket，监听端口
	flags.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	//fmt.Println(debug.SetGCPercent(-1) )

	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	//go utils.GC(GCCycle)
	netListen, err := net.Listen("tcp", "0.0.0.0:"+*flags.ListPort)
	utils.CheckErr(err)
	defer netListen.Close()
	//
	YaoHongbaoId = 0x100000
	FaHongbaoId = 0x200000
	InitProcess()
	go monitor.Show(*(flags.MonitorServer), 1)
	go httpServer()
	for i := 0; i < HongBaoMailBoxNum; i++ {
		go ProcessGetHongbao(i)
	}
	num := 0
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		num += 1
		ppool := getOnePool(num)

		//	Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(&conn, ppool)
	}
}

//处理连接
func handleConnection(conn *net.Conn, ppool *MessagePool) {
	defer (*conn).Close()
	defer ppool.Reduce()
	// 先拿到login Request
	p, result := parse.ReadMessage(conn, ppool)
	if !result {
		return
	}

	tnow := time.Now()
	monitor.GMonitor.Add(monitor.Request, 1, tnow.Unix())
	monitor.GMonitor.Add(p.Cmd, 1, tnow.Unix())
	if p.Cmd != LoginRequest_ID {
		fmt.Print("wrong request")
		return

	}

	index, _ := utils.GetMailBoxAddress(MailBoxNum, int(p.P1.Userid))
	gMailbox[index].SetConn(conn, MailBoxNum, int(p.P1.Userid))
	defer gMailbox[index].SetConn(nil, MailBoxNum, int(p.P1.Userid))

	userid := p.P1.Userid
	{
		a := ppool.Pop()
		a.Cmd = LoginResponse_ID
		p := LoginResponse{}
		p.Userid = userid
		p.Result = 0
		a.P2 = &p

		b2 := parse.WriteMessage(conn, a)
		if b2 != nil {
			(*conn).Close()
		}
	}

	for {
		p, result := parse.ReadMessage(conn, ppool)
		if !result {
			break
		}
		tnow := time.Now()
		if (uint64(tnow.Unix()) > p.TimeStamp) && (uint64(tnow.Unix())-p.TimeStamp > 5) {
			//fmt.Println("delay message", uint64(tnow.Unix()), p.TimeStamp, p.Cmd)
			p.TimeStamp = uint64(tnow.Unix())
		}
		qps := monitor.GMonitor.Add(monitor.Request, 1, tnow.Unix())
		monitor.GMonitor.Add(p.Cmd, 1, tnow.Unix())
		p.Conn = conn

		if qps > MaxQPS && p.Cmd == YaoHongbaoRequest_ID { // too many yao request

			{
				a := ppool.Pop()
				a.Cmd = YaoHongbaoResponse_ID
				p := YaoHongbaoResponse{}
				p.Userid = userid
				p.Result = monitor.OverLoad
				p.Hongbaoid = 0
				a.P4 = &p

				b2 := parse.WriteMessage(conn, a)
				if b2 != nil {
					(*conn).Close()
				}

			}
			p.Rollback()

			continue
		}
		gMailbox[index].PushInput(p)

	}

}

func getRequestValue(req *http.Request, key string) int {

	filelist, found1 := req.Form[key]

	if !(found1) {
		fmt.Println("no qps")
		return 0
	}

	value, err := strconv.Atoi(filelist[0])
	if err != nil {
		return 0
	}
	return value

}

func fahongbao(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	fmt.Println(req.Form)

	sum := 0
	qps := 0
	minuserid := 0
	maxuserid := 0

	qps = getRequestValue(req, "qps")
	sum = getRequestValue(req, "sum")
	minuserid = getRequestValue(req, "min")
	maxuserid = getRequestValue(req, "max")
	if minuserid < 0 {
		minuserid = -minuserid
	}
	if maxuserid < 0 {
		maxuserid = -maxuserid
	}
	if maxuserid < minuserid {
		t := minuserid
		minuserid = maxuserid
		maxuserid = t
	}
	hongbaoid := FaHongbaoId
	FaHongbaoId += int64(sum)
	fmt.Println(sum, qps, hongbaoid, minuserid, maxuserid)
	go FaHongbao(sum, qps, hongbaoid, minuserid, maxuserid)

}

func yaohongbao(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	fmt.Println(req.Form)
	sum := 0
	qps := 0

	qps = getRequestValue(req, "qps")
	sum = getRequestValue(req, "sum")
	fmt.Println("qps:", qps, "sum", sum)
	hongbaoid := YaoHongbaoId
	YaoHongbaoId += int64(sum)
	go YaoHongbao(sum, qps, hongbaoid)

}

var YaoHongbaoId int64
var YaoHongbaoMoney int64

func YaoHongbao(sum int, qps int, hongbaoid int64) {
	x := 0
	rand.Seed(time.Now().UnixNano() % 1e9)

	utils.CorrectSleepOneSecond()

	tnow := time.Now()
	for i := 0; i < sum; i++ {
		x += 1
		{
			a := gPool[MessagePoolNum].Pop()
			a.Cmd = YaoHongbaoResponse_ID
			p := YaoHongbaoResponse{}
			p.Userid = 0
			p.Result = monitor.GetHongbao
			p.Hongbaoid = uint32(hongbaoid)
			money := uint32(rand.Intn(300) + 10)
			p.Amount = money
			a.P4 = &p
			gMailbox[i%MailBoxNum].PushYao(a)

			monitor.GMonitor.Add(monitor.GetHongbao, 1, tnow.Unix())
			monitor.GMonitor.Add2(monitor.CountYaoHongbaoReq, 1)
			monitor.GMonitor.Add2(monitor.SumYaoHongbao, int64(money))

		}
		if x >= qps {
			utils.CorrectSleepOneSecond2(tnow)
			tnow = time.Now()
			x = 0
		}
	}

}

var FaHongbaoId int64
var FaHongbaoMoney int64

func FaOneHongbao(hongbaoid int64, minuserid int, maxuserid int) {
	users := [3]int{}
	for i := 0; i < 3; i++ {
		users[i] = rand.Intn(maxuserid-minuserid) + minuserid
	}
	Amount := rand.Intn(300) + 10
	for i := 3; i >= 0; i-- {

		a := gPool[MessagePoolNum+1].Pop()
		a.Cmd = FaHongbaoRequest_ID
		p := &FaHongbaoRequest{}
		p.Hongbaoid = uint32(hongbaoid)
		p.Amount = uint32(Amount)
		p.Useramount = 3
		p.Userid1 = uint32(users[0])
		p.Userid2 = uint32(users[1])
		p.Userid3 = uint32(users[2])
		a.P5 = p
		if i < 3 {
			index, _ := utils.GetMailBoxAddress(MailBoxNum, users[i])
			gMailbox[index].PushYao(a)
		} else {
			gindex := hongbaoid % HongBaoMailBoxNum
			gHongbaoMailbox[gindex].Push(a)
			monitor.GMonitor.Add2(monitor.CountFaHongbaoReq, 1)
			monitor.GMonitor.Add2(monitor.SumFaHongbao, int64(Amount))
		}

	}
}

func FaHongbao(sum int, qps int, hongbaoid int64, minuserid int, maxuserid int) {
	x := 0
	rand.Seed(time.Now().UnixNano() % 1e9)
	utils.CorrectSleepOneSecond()

	tnow := time.Now()
	for i := 0; i < sum; i++ {
		x += 1
		{
			FaOneHongbao(hongbaoid, minuserid, maxuserid)
			hongbaoid += 1
		}
		if x >= qps {
			utils.CorrectSleepOneSecond2(tnow)
			tnow = time.Now()
			x = 0
		}
	}

}

func httpServer() {

	http.HandleFunc("/fahongbao", fahongbao)
	http.HandleFunc("/yaohongbao", yaohongbao)
	//服务器要监听的主机地址和端口号
	err := http.ListenAndServe("0.0.0.0:8989", nil)
	utils.CheckErr(err)
}

type HongbaoRecord struct {
	p             *FaHongbaoRequest
	Amount        int
	Useramount    int
	GotUseramount int
	GotAmount     int
}

func SplitHongbao(Amount int, usernum int) int {
	if usernum == 1 {
		return Amount
	}
	if Amount <= usernum {
		return 1
	}
	if Amount/usernum < 1 {
		return 1
	}
	return rand.Intn(Amount / usernum)
}

func ProcessOneHongbao(p *HongbaoRecord) (bool, int) {

	num := SplitHongbao(p.Amount-p.GotAmount, p.Useramount-p.GotUseramount)
	p.GotAmount += num
	p.GotUseramount += 1

	if p.GotUseramount == p.Useramount {
		return true, num
	}
	return false, num

}

func ProcessGetHongbao(boxindex int) {

	hongbaoMap := make(map[uint32]*HongbaoRecord)

	for {
		mess, _ := gHongbaoMailbox[boxindex].BlockingPop()
		switch mess.Cmd {
		case FaHongbaoRequest_ID:
			{
				hongbaoid := mess.P5.Hongbaoid
				_, exists := hongbaoMap[hongbaoid]
				if !exists {
					p := &HongbaoRecord{}
					p.Amount = int(mess.P5.Amount)
					p.Useramount = 3
					hongbaoMap[hongbaoid] = p
				}
			}

		case GetHongbaoRequest_ID:
			{
				hongbaoid := mess.P6.Hongbaoid
				p, exists := hongbaoMap[hongbaoid]
				finished := false
				num := 0

				tnow2 := time.Now()
				if (uint64(tnow2.Unix()) > mess.TimeStamp) && (uint64(tnow2.Unix())-mess.TimeStamp > 5) {
					fmt.Println("process delay message box: ", boxindex, " ", monitor.Request, " ", uint64(tnow2.Unix()), " ", mess.TimeStamp, uint64(tnow2.Unix())-mess.TimeStamp)
				}

				if exists {
					finished, num = ProcessOneHongbao(p)
					if finished {
						delete(hongbaoMap, hongbaoid)
						if len(hongbaoMap) == 0 {
							fmt.Println("hongbaomap", len(hongbaoMap))
						}
					}
					// send a response to client

					{
						a := mess.Pool.Pop()
						a.Cmd = GetHongbaoRespone_ID
						p := GetHongbaoRespone{}
						p.Hongbaoid = hongbaoid
						p.Amount = uint32(num)
						a.P7 = &p

						b2 := parse.WriteMessage(mess.Conn, a)
						if b2 != nil {
							(*mess.Conn).Close()
						}
						monitor.GMonitor.Add2(monitor.CountGetFaHongbaoReq, 1)

						monitor.GMonitor.Add(GetHongbaoRespone_ID, 1, time.Now().Unix())

					}
				}

			}
		}
	}
}
