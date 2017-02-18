package message2

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"
	"time"
)

// 使用protobuff 应该会更快 更好
// 但是fakechat 项目中,性能并不能让我满意，所以我简单的实现一下。
// 如果修改也非常简单

func testEq(a, b []byte) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func CheckEq(a, b []byte, t *testing.T) {
	if !testEq(a, b) {
		t.Error(a)
		t.Error(b)
		t.Error("Wrong")
	}
}

func Test_monitor(t *testing.T) {

	t1 := time.Now()
	var i uint32
	for i = 0; i < 20*10*1000; i++ {
		{
			a := &MsgItem{}
			id := int32(i)
			a.Id = &id
			R := int32(i + 1)
			a.R = &R

			a_serialization, err := proto.Marshal(a)
			if err != nil {

				t.Error("wrong  ")
			}
			b := &MsgItem{}
			err = proto.Unmarshal(a_serialization, b)
			if err != nil {
				t.Error("wrong  ")
			}
			if (*b.Id != *a.Id) || (*b.R != *a.R) {
				t.Error("wrong  ")
			}
		}

	}
	t2 := time.Now()

	fmt.Println((t2.Nanosecond() - t1.Nanosecond()) / 1e6)

}
