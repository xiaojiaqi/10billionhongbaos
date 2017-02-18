package monitor

import (
	"fmt"
	"testing"
	"time"
)

func Test_monitor(t *testing.T) {
	a := &Monitor{}
	timestamp := time.Now().Unix()
	var i int64
	for i = 0; i < TimeRange; i++ {
		a.Add(Request, 1+i, timestamp+i)
		a.Add(Response, 2+i, timestamp+i)
		a.Add(GetHongbao, 3+i, timestamp+i)
		a.Add(NoHongbao, 4+i, timestamp+i)
		a.Add(AllHongbao, 5+i, timestamp+i)
	}

	for i = 0; i < TimeRange; i++ {
		a.ShowStatus(timestamp+i, 0)

		v1, _ := a.Get(timestamp + i)

		fmt.Println(v1)
		a.Clean(timestamp + i)
		v1, _ = a.Get(timestamp + i)
		fmt.Println(v1)

	}
}
