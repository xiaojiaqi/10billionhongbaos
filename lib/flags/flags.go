package flags

import (
	"flag"
)

var ListPort *string

var Server *string
var MonitorServer *string

var QPS uint32
var MinUserId uint32
var MaxUserId uint32
var TotalUser uint32

var PQPS *uint
var PMinUserId *uint
var PMaxUserId *uint
var PTotalUser *uint

func Parse() {

	ListPort = flag.String("list", "1025", "list port")

	Server = flag.String("server", "10.29.101.192:1025", "server address")
	MonitorServer = flag.String("monitor", "10.29.101.3:8002", "Monitor server address")

	PQPS = flag.Uint("qps", 100, "QPS")
	PMinUserId = flag.Uint("min", 0, "MinUserId")
	PMaxUserId = flag.Uint("max", 99, "MaxUserId")
	PTotalUser = flag.Uint("total", 100, "TotalUser")
	flag.Parse()

	QPS = (uint32)(*PQPS)
	MinUserId = (uint32)(*PMinUserId)
	MaxUserId = (uint32)(*PMaxUserId)
	TotalUser = (uint32)(*PTotalUser)
}
