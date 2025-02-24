package Algorithm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

var Sea06_Data, _ = hex.DecodeString(SeaDatIpad)
var Sea01_Data, _ = hex.DecodeString(SeaDatAndroid)

//////////////////////////////////////////////// sae06加密 /////////////////////////////////////////////////////

func SaeEncrypt06(data []byte) []byte {
	sae := Sea06_Data
	in_bytes := data
	in_len := len(data)
	xor := sae[9:25]
	input_val := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	output_val := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var result []byte
	for i := 0; i < in_len/16; i++ {
		for j := 0; j < 16; j++ {
			input_val[j] = xor[j] ^ in_bytes[i*16+j]
		}
		output_val = DoEncryptInput(input_val, sae)
		xor = output_val
		result = BytesCombine1(result, output_val)
		input_val = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		output_val = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}
	return result
}

func SaeEncrypt01(data []byte) []byte {
	sae := Sea01_Data
	in_bytes := data
	in_len := len(data)
	xor := sae[9:25]
	input_val := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	output_val := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var result []byte
	for i := 0; i < int(in_len/16); i++ {
		for j := 0; j < 16; j++ {
			input_val[j] = xor[j] ^ in_bytes[i*16+j]
		}
		output_val = DoEncryptInput(input_val, sae)
		xor = output_val
		result = BytesCombine1(result, output_val)
		input_val = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		output_val = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}
	return result
}

func DoEncryptInput(input_val, sae_val []byte) []byte {
	in_p := input_val
	output := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			output[i*4+j] = in_p[j*4+i]
		}
	}
	pos := -0x24000
	sae_pos := 0x82030
	for {
		output = LeftShift(output, 4, 1)
		output = LeftShift(output, 8, 2)
		output = LeftShift(output, 12, 3)
		if pos < 0 {
			buf := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
			buf = Sub106c03420(buf, output, sae_val[0x43030+pos:])
			output = Sub106c036a8(output, buf, sae_val[sae_pos-0x3f000:])
			sae_pos = sae_pos + 0x3000
			pos = pos + 0x4000
			continue
		}
		break
	}
	result := Sub106c0397c(output, sae_val[0xbc030:])
	result_p := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < 0x10; i++ {
		result_p[i] = result[int(i/4)+int((i%4)*4)]
	}
	return result_p
}

func LeftShift(data []byte, offset, distance int) []byte {

	/*左轮转
	:param data: 需要被轮转的数据
	:param offset: 偏移量
	:param distance: 轮转距离
	:return 轮转结果*/

	for i := 0; i < distance; i++ {
		tmp := data[offset]
		data[offset] = data[1+offset]
		data[1+offset] = data[2+offset]
		data[2+offset] = data[3+offset]
		data[3+offset] = tmp
	}
	return data
}

func Sub106c03420(buf, data, sae []byte) []byte {
	//sae06字典加密
	v3 := 0
	buf_pos := 0
	sae_pos := 0
	for {
		if v3 == 4 {
			break
		}
		v4 := 0
		v5 := buf_pos
		v6 := sae_pos
		for {
			if v4 == 4 {
				break
			}
			v7 := 0
			v8 := v6
			for {
				if v7 == 0x40 {
					break
				}
				buf[v5+v7] = sae[v8+4*int(data[4*v3+v4])] //根据data序列查字典, 结果放入buf
				v8 = v8 + 1
				v7 = v7 + 0x10
			}
			v4 = v4 + 1
			v6 = v6 + 0x400
			v5 = v5 + 4
		}
		v3 = v3 + 1
		sae_pos = sae_pos + 0x1000
		buf_pos = buf_pos + 1
	}
	return buf
}

func Sub106c036a8(data, buf, sae []byte) []byte {
	v3 := 0
	v4 := 0
	v5 := 0
	v6 := 0x200
	v7 := 0
	for {
		if v5 == 4 {
			break
		}
		v8 := 0
		v14 := v6
		v9 := v7
		for {
			if v8 == 4 {
				break
			}
			result := buf[v3+16*v5+4*v8+3]
			v11 := v4 + 4*v5 + v8
			data[v11] = result
			v12 := v6
			v13 := 2
			for {
				if v13 == -1 {
					break
				}
				result = Sub106c033bc(buf[v9+v13], result, sae[v12:])
				data[v11] = result
				v13 = v13 - 1
				v12 = v12 - 0x100
			}
			v8 = v8 + 1
			v9 = v9 + 4
			v6 = v6 + 0x300
		}
		v5 = v5 + 1
		v7 = v7 + 0x10
		v6 = v14 + 0xc00
	}
	return data
}

func Sub106c033bc(a1 byte, a2 byte, sae []byte) byte {
	temp := 0xffffff0f
	v3 := (a1 & 0xf0) | (a2 >> 4)
	if v3&0x80 != 0 {
		v3 = sae[(0+(v3&0x7f))] >> 4
	} else {
		v3 = sae[(0+v3)] & 0xf
	}
	v4 := ((a2 & 0xf) | 0x10*a1) & 0xff
	var v5 byte
	if v4&0x80 != 0 {
		v5 = sae[(0+int(v4&0x7f)+0x80)] >> 4
	} else {
		v5 = sae[(0+int(v4&0xff)+0x80)] & 0xf
	}
	return (v5 & uint8(temp)) | 0x10*(v3&0xf)
}

func Sub106c0397c(output, sae []byte) []byte {
	v3 := 0
	v5 := 0
	result := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	v4 := 0
	for {
		if v3 == 4 {
			break
		}
		v6 := 0
		v7 := 0
		for {
			if v6 == 4 {
				break
			}
			result[v4+v6] = sae[v7+int(output[v5+v6])]
			v6 = v6 + 1
			v7 = v7 + 0x100
		}
		v3 = v3 + 1
		v5 = v5 + 4
		sae = sae[0x400:]
		v4 = v4 + 4
	}
	return result
}

func BytesCombine1(pBytes ...[]byte) []byte {
	length := len(pBytes)
	s := make([][]byte, length)
	for index := 0; index < length; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// todo
// new sae 白盒aes加密
func SaeEncrypt07(data []byte) []byte {
	//url := beego.AppConfig.String("signserver") + "/Sign"
	//method := "POST"
	//
	//payload := strings.NewReader(`{
	//     "FunctionID": 1,
	//     "FunctionData": "` + base64.StdEncoding.EncodeToString(data) + `"
	// }`)
	//
	//client := &http.Client{}
	//req, err := http.NewRequest(method, url, payload)
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return nil
	//}
	//req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7,zh-TW;q=0.6")
	//req.Header.Add("Cache-Control", "no-cache")
	//req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Pragma", "no-cache")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	//req.Header.Add("accept", "application/json")
	//req.Header.Add("content-type", "application/json")
	//
	//res, err := client.Do(req)
	//if err != nil {
	//	fmt.Println(err)
	//	return nil
	//}
	//defer res.Body.Close()
	//
	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	return nil
	//}
	//var callData RespData
	//err = json.Unmarshal(body, &callData)
	//if err != nil || callData.Code != 0 {
	//	fmt.Println(err)
	//	return nil
	//}
	////data := CalcMsgCrcForString_7019(hex.EncodeToString(h[:]))
	////fmt.Println(strconv.Atoi(callData.data))
	//sign, err := base64.StdEncoding.DecodeString(callData.Data)
	//if err != nil || callData.Code != 0 {
	//	fmt.Println(err)
	//	return nil
	//}
	//fmt.Println(len(sign))
	//return sign
	var key, iv []byte

	key, _ = hex.DecodeString("24E545FC309F1CC92B0223FAFA8C84F4")
	iv, _ = hex.DecodeString("6d1f24b8e29268d6efbf55cafb27d3bf")
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	src := PKCS7Padding(data, block.BlockSize())
	origData := make([]byte, len(src))
	blockMode.CryptBlocks(origData, src)
	return origData
}
