package lib

import (
	"encoding/hex"
	"strings"
)

func Get62Key(Key string) string {
	if len(Key) < 344 {
		return MD5ToLower(RandSeq(15))
	}
	FinIndex := strings.Index(Key,"6E756C6C5F1020")
	if  FinIndex != -1 {
		head := FinIndex + len("6E756C6C5F1020")
		K, _ := hex.DecodeString(Key[head:head+64])
		return string(K)
	}
	K, _ := hex.DecodeString(Key[134:198])
	return string(K)
}
