package lib

import (
	"strconv"
	"time"
)

func GetClientSeqId(DeviceId string) string {
	return DeviceId + "-" + strconv.Itoa(int(time.Now().Unix()))
}
