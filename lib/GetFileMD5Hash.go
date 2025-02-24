package lib

import (
	"crypto/md5"
	"encoding/hex"
)

func GetFileMD5Hash(Data []byte) string {
	hash := md5.New()
	hash.Write(Data)
	retVal := hash.Sum(nil)
	return hex.EncodeToString(retVal)
}
