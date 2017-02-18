package qps

import (
	"testing"
)

func Test_hash(t *testing.T) {
	var totaluser uint32
	var i uint32

	for totaluser = 0; totaluser <= 30000; totaluser++ {
		v := make([]int, int(totaluser))
		for i = 0; i < totaluser; i++ {
			index := Hash(i, totaluser)
			v[index] += 1
		}
		for i = 0; i < totaluser; i++ {
			if v[i] != 1 {
				t.Error(i, "wrong")
			}
		}
	}
}

func Test_qps2(t *testing.T) {
	//    1484553130
	//    -min 60000 -max 119999 -total 180000 -qps 3000
	sum := 0
	sum2 := 0
	i := 0
	for i = 0; i < 180000; i++ {
		b := TriggerRequest(uint32(i), 180000, 3000, 1484553130)
		if b == true {
			sum += 1
		}

	}

	for i = 0; i < 180000; i++ {
		b := TriggerRequest2(uint32(i), 180000, 3000, 1484553130)
		if b == true {
			sum2 += 1
		}

	}

	if sum != sum2 {
		t.Error(sum, sum2)
	}
}

func Test_qps(t *testing.T) {
	if false == TriggerRequest(0, 100, 10, 0) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(10, 100, 10, 0) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(11, 100, 10, 0) {
		t.Error("wrong  ")
	}
	if false == TriggerRequest(11, 100, 10, 1) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(99, 100, 10, 9) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(0, 100, 10, 10) {
		t.Error("wrong  ")
	}

	/*
		/0-32
		/33-65
		/66-98
		/99-31
		/32-64
	*/

	// 0  /0-32

	if false == TriggerRequest(0, 100, 33, 0) {
		t.Error("wrong  ")
	}
	if false == TriggerRequest(31, 100, 33, 0) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(32, 100, 33, 0) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(33, 100, 33, 0) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(64, 100, 33, 0) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(65, 100, 33, 0) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(66, 100, 33, 0) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(98, 100, 33, 0) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(99, 100, 33, 0) {
		t.Error("wrong  ")
	}

	//1 /33-65

	if true == TriggerRequest(0, 100, 33, 1) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(31, 100, 33, 1) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(32, 100, 33, 1) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(33, 100, 33, 1) {
		t.Error("wrong  ")
	}
	if false == TriggerRequest(64, 100, 33, 1) {
		t.Error("wrong  ")
	}
	if false == TriggerRequest(65, 100, 33, 1) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(66, 100, 33, 1) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(98, 100, 33, 1) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(99, 100, 33, 1) {
		t.Error("wrong  ")
	}

	//2  /66-98

	if true == TriggerRequest(0, 100, 33, 2) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(31, 100, 33, 2) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(32, 100, 33, 2) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(33, 100, 33, 2) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(64, 100, 33, 2) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(65, 100, 33, 2) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(66, 100, 33, 2) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(98, 100, 33, 2) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(99, 100, 33, 2) {
		t.Error("wrong  ")
	}

	// 3 /99-31

	if false == TriggerRequest(0, 100, 33, 3) {
		t.Error("wrong  ")
	}
	if false == TriggerRequest(31, 100, 33, 3) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(32, 100, 33, 3) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(33, 100, 33, 3) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(64, 100, 33, 3) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(65, 100, 33, 3) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(66, 100, 33, 3) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(98, 100, 33, 3) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(99, 100, 33, 3) {
		t.Error("wrong  ")
	}

	//4  /32-64
	if true == TriggerRequest(0, 100, 33, 4) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(31, 100, 33, 4) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(32, 100, 33, 4) {
		t.Error("wrong  ")
	}

	if false == TriggerRequest(33, 100, 33, 4) {
		t.Error("wrong  ")
	}
	if false == TriggerRequest(64, 100, 33, 4) {
		t.Error("wrong  ")
	}
	if true == TriggerRequest(65, 100, 33, 4) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(66, 100, 33, 4) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(98, 100, 33, 4) {
		t.Error("wrong  ")
	}

	if true == TriggerRequest(99, 100, 33, 4) {
		t.Error("wrong  ")
	}
}
