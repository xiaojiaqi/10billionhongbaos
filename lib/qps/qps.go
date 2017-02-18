package qps

func Hash(userid uint32, totaluser uint32) uint32 {
	var bucket uint32
	bucket = 13
	if totaluser <= bucket {
		return userid
	}
	length := totaluser / bucket
	remainder := totaluser % bucket
	u_length := userid / bucket
	u_remainder := userid % bucket
	newid := length * u_remainder
	if u_remainder > remainder {
		newid += remainder
	} else {
		newid += u_remainder
	}
	newid += u_length

	return newid
}

func TriggerRequest2(userid uint32, totaluser uint32, qps uint32, timeStamp int64) bool {

	return TriggerRequest(Hash(userid, totaluser), totaluser, qps, timeStamp)
}

func TriggerRequest(userid uint32, totaluser uint32, qps uint32, timeStamp int64) bool {
	second := uint32(timeStamp)
	if qps > totaluser {
		qps = totaluser
	}
	if totaluser%qps == 0 {
		id := totaluser / qps
		if second%id == userid%id {
			return true
		} else {
			return false
		}
	}

	index := second * qps % totaluser
	if index+qps > totaluser {
		if (userid >= index) && userid <= (totaluser-1) {
			return true
		}
		if userid >= 0 && userid <= (qps-(totaluser-index)-1) {
			return true
		}
		return false
	} else {
		if userid >= index && userid <= (index+qps-1) {
			return true
		}
		return false
	}

}
