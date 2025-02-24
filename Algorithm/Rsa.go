package Algorithm

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"math/big"
)

func Compress_rsa(data []byte, Key string) []byte {
	strOut := new(bytes.Buffer)
	var publicKey rsa.PublicKey
	s, _ := hex.DecodeString(Key)
	publicKey.N = new(big.Int).SetBytes(s)
	rsaLen := len(Key) / 8
	if len(data) > (rsaLen - 12) {
		blockCnt := 1
		if ((len(data) / (rsaLen - 12)) + (len(data) % (rsaLen - 12))) == 0 {
			blockCnt = 0
		}

		for i := 0; i < blockCnt; i++ {
			blockSize := rsaLen - 12
			if i == blockCnt-1 {
				blockSize = len(data) - i*blockSize
			}
			temp := data[(i * (rsaLen - 12)):(i*(rsaLen-12) + blockSize)]
			encrypted, _ := rsa.EncryptPKCS1v15(rand.Reader, &publicKey, temp)
			strOut.Write(encrypted)
		}
		return strOut.Bytes()
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, &publicKey, data)
	if err != nil {
		return []byte{}
	}
	return encrypted
}

//RSAEncrypt Rsa加密
func RSAEncrypt(data []byte, Key string) []byte {
	m := Key
	M := new(big.Int)
	M.SetString(m, 16)
	pub := rsa.PublicKey{}
	pub.E = 65537
	pub.N = M
	out, _ := rsa.EncryptPKCS1v15(rand.Reader, &pub, data)
	return out
}
