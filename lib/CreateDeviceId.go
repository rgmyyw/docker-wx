package lib

import (
	"crypto/md5"
	"encoding/hex"
)

func CreateDeviceId(s string) string {
	if s == "" || s == "string" {
		s = RandSeq(15)
	}

	h := md5.New()
	h.Write([]byte(s))
	md5string := hex.EncodeToString(h.Sum(nil))
	return "49" + md5string[2:]
}
