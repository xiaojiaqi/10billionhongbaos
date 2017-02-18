package message

import (
	"fmt"
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
			a := Message{}
			a.Init()
			a.Userid = uint(i)
			a.Cmd = LoginRequest_ID
			p := LoginRequest{}
			p.Userid = i
			a.P1 = &p
			a_serialization := a.Serialization()
			if len(a_serialization) != (HeaderLength + a.P1.GetLength()) {
				t.Error("wrong  ")
			}
			a1 := Message{}

			a1.Init()
			a1.Unserialization(a_serialization)
			if (a1.P1.Userid != a.P1.Userid) || (a.Cmd != a1.Cmd) {
				t.Error("wrong  ")
			}
		}

		{
			a := Message{}

			a.Init()
			a.Userid = uint(i)
			a.Cmd = LoginResponse_ID
			p := LoginResponse{}
			p.Userid = i
			p.Result = i + 1
			a.P2 = &p
			a_serialization := a.Serialization()
			if len(a_serialization) != (HeaderLength + a.P2.GetLength()) {
				t.Error("wrong  ")
			}

			a1 := Message{}

			a1.Init()
			a1.Unserialization(a_serialization)
			if (a1.Cmd != a.Cmd) || (a1.P2.Userid != a.P2.Userid) || (a1.P2.Result != a.P2.Result) {
				t.Error("wrong  ", a1.Cmd, a.Cmd, a1.P2.Userid, a.P2.Userid, a1.P2.Result, a.P2.Result)

			}

		}

		{
			a := Message{}

			a.Init()
			a.Userid = uint(i)
			a.Cmd = YaoHongbaoRequest_ID
			p := YaoHongbaoRequest{}
			p.Userid = i
			a.P3 = &p
			a_serialization := a.Serialization()
			if len(a_serialization) != (HeaderLength + a.P3.GetLength()) {
				t.Error("wrong  ")
			}

			a1 := Message{}

			a1.Init()
			a1.Unserialization(a_serialization)
			if (a1.Cmd != a.Cmd) || (a1.P3.Userid != a.P3.Userid) {
				t.Error("wrong  ")
			}

		}

		{
			a := Message{}

			a.Init()

			a.Userid = uint(i)
			a.Cmd = YaoHongbaoResponse_ID
			p := YaoHongbaoResponse{}
			p.Userid = i
			p.Result = i + 1
			p.Hongbaoid = i + 3
			p.Amount = i + 1000
			a.P4 = &p
			a_serialization := a.Serialization()
			if len(a_serialization) != (HeaderLength + a.P4.GetLength()) {
				t.Error("wrong  ")
			}

			a1 := Message{}

			a1.Init()
			a1.Unserialization(a_serialization)
			if (a1.Cmd != a.Cmd) || (a1.P4.Userid != a.P4.Userid) || (a1.P4.Result != a.P4.Result) || (a1.P4.Hongbaoid != a.P4.Hongbaoid) {
				t.Error("wrong  ")
			}

		}

		{
			a := Message{}

			a.Init()
			a.Userid = uint(i)
			a.Cmd = FaHongbaoRequest_ID
			p := FaHongbaoRequest{}
			p.Hongbaoid = i
			p.Amount = i + 1
			p.Useramount = i + 2
			p.Userid1 = i + 3
			p.Userid2 = i + 4
			p.Userid3 = i + 5

			a.P5 = &p

			a_serialization := a.Serialization()
			if len(a_serialization) != (HeaderLength + a.P5.GetLength()) {
				t.Error("wrong  ")
			}

			a1 := Message{}

			a1.Init()
			a1.Unserialization(a_serialization)
			if (a1.Cmd != a.Cmd) || (a1.P5.Hongbaoid != a.P5.Hongbaoid) || (a1.P5.Amount != a.P5.Amount) || (a1.P5.Useramount != a.P5.Useramount) || (a1.P5.Userid1 != a.P5.Userid1) || (a1.P5.Userid2 != a.P5.Userid2) || (a1.P5.Userid3 != a.P5.Userid3) {
				t.Error("wrong  ")
			}

		}

		{
			a := Message{}

			a.Init()
			a.Userid = uint(i)
			a.Cmd = GetHongbaoRequest_ID
			p := GetHongbaoRequest{}
			p.Userid = i
			p.Hongbaoid = i + 3
			a.P6 = &p
			a_serialization := a.Serialization()
			if len(a_serialization) != (HeaderLength + a.P6.GetLength()) {
				t.Error("wrong  ")
			}

			a1 := Message{}

			a1.Init()
			a1.Unserialization(a_serialization)
			if (a1.Cmd != a.Cmd) || (a1.P6.Userid != a.P6.Userid) || (a1.P6.Hongbaoid != a.P6.Hongbaoid) {
				t.Error("wrong  ")
			}

		}

		{
			a := Message{}

			a.Init()
			a.Userid = uint(i)
			a.Cmd = GetHongbaoRespone_ID
			p := GetHongbaoRespone{}
			p.Amount = i + 1
			p.Hongbaoid = i + 3
			a.P7 = &p
			a_serialization := a.Serialization()
			if len(a_serialization) != (HeaderLength + a.P7.GetLength()) {
				t.Error("wrong  ")
			}

			a1 := Message{}

			a1.Init()
			a1.Unserialization(a_serialization)
			if (a1.Cmd != a.Cmd) || (a1.P7.Hongbaoid != a.P7.Hongbaoid) || (a1.P7.Amount != a.P7.Amount) {
				t.Error("wrong  ")
			}

		}

	}
	t2 := time.Now()
	fmt.Println(t2.Sub(t1))

}
