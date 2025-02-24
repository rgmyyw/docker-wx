package Algorithm

func XorEncodeStr(msg string, key uint8) string {
	ml := len(msg)
	pwd := ""
	for i := 0; i < ml; i++ {

		pwd += string(key ^ msg[i])

	}
	return pwd
}

func AE(msg []byte) []byte {
	ml := len(msg)
	var pwd []byte
	for i := 0; i < ml; i++ {
		S := 0XAE ^ msg[i]
		pwd = append(pwd,S)
	}
	return DoZlibCompress(pwd)
}

func XorDecodeStr(msg string, key uint8) string {
	ml := len(msg)
	pwd := ""
	for i := 0; i < ml; i++ {

		pwd += string(msg[i] ^ key)

	}
	return pwd
}