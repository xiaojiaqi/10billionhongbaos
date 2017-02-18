package utils

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(-1)
	}
}

func GC(sleep int) {
	var m runtime.MemStats

	for {
		time.Sleep(time.Duration(sleep) * time.Second)
		fmt.Println("gc")
		//            runtime.GC()
		//   var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("%d,%d,%d,%d\n", m.HeapSys, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased)

	}

}

func CorrectSleepOneSecond() {
	t := time.Now()
	second := t.UnixNano()
	sleep := 1e9 - (second % 1e9)
	time.Sleep(time.Duration(sleep) * time.Nanosecond)
}

func CorrectSleepOneSecond2(t1 time.Time) {
	tnow2 := time.Now().UnixNano()

	if tnow2-t1.UnixNano() > 1e9 {
		return
	}

	time.Sleep(time.Duration(1e9-(tnow2%1e9)) * time.Nanosecond)
}

func GetMailBoxAddress(mailboxnum int, userid int) (x int, y int) {
	x = userid % mailboxnum

	y = userid / mailboxnum
	return x, y

}
