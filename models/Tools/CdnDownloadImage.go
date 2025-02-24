package Tools

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strconv"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

func CdnDownloadImg(Data CdnDownloadImageParam, Dns mm.GetCDNDnsResponse) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	// 建立连接
	cdnAddress := *Dns.FakeDnsInfo.ZoneIPList[0].String_ + ":443"
	conn, connErr := net.Dial("tcp", cdnAddress)
	if connErr != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", connErr.Error()),
			Data:    nil,
		}
	}
	inquiry := AskCdnReadyPack(*Dns.FakeDnsInfo.Uin, Dns.FakeDnsInfo.AuthKey.Buffer, D.ClientVersion, D.DeviceType)
	// fmt.Printf("%x\n", inquiry)
	conn.Write(inquiry)
	readyBuf := make([]byte, 0x20000)
	_, err = conn.Read(readyBuf)
	if err != nil && err != io.EOF {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", connErr.Error()),
			Data:    nil,
		}
	}

	// 请求文件
	requestBytes := RequestImagePack(*Dns.FakeDnsInfo.Uin, Dns.FakeDnsInfo.AuthKey.Buffer, D.ClientVersion, D.DeviceType, Data.FileNo, Data.FileAesKey)
	// fmt.Printf("%x\n", requestBytes)
	conn.Write(requestBytes)
	// 接收文件
	var fileBytes []byte
	recBuf := make([]byte, 0x20000)
	len, err := conn.Read(recBuf)
	// fmt.Printf("%x\n", recBuf)
	for len > 0 {
		messageBlock := recBuf[:len]
		fileBytes = append(fileBytes, messageBlock...)
		finished, fileBytes := CheckReceiveFinished(fileBytes)
		if finished {
			conn.Close()
			resMap := UnpackMessage(fileBytes)
			encryptedImage, ok := resMap["filedata"]
			if !ok {
				return models.ResponseResult{
					Code:    -8,
					Success: false,
					Message: fmt.Sprintf("异常：下载失败"),
					Data:    nil,
				}
			}

			aesKeyByte, _ := hex.DecodeString(Data.FileAesKey)
			plainImage := Algorithm.AesDecryptECB([]byte(encryptedImage), aesKeyByte)
			return models.ResponseResult{
				Code:    0,
				Success: true,
				Message: "成功",
				Data:    CdnImageBase64{
					Image: base64.StdEncoding.EncodeToString(plainImage),
			},
			}
		} else {
			len, err = conn.Read(recBuf)
		}
	}
	return models.ResponseResult{
		Code:    -8,
		Success: false,
		Message: fmt.Sprintf("请求CDN数据失败"),
		Data:    nil,
	}
}

func AskCdnReadyPack(uin uint32, authKey []byte, clientVersion int, deviceType string) []byte {
	bodyBytes := new(bytes.Buffer)
	// ver
	binary.Write(bodyBytes, binary.BigEndian, int32(len("ver")))
	bodyBytes.Write([]byte("ver"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("1"))
	// uin
	binary.Write(bodyBytes, binary.BigEndian, int32(len("weixinnum")))
	bodyBytes.Write([]byte("weixinnum"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(strconv.FormatUint(uint64(uin), 10))))
	bodyBytes.Write([]byte(strconv.FormatUint(uint64(uin), 10)))
	// seq
	binary.Write(bodyBytes, binary.BigEndian, int32(len("seq")))
	bodyBytes.Write([]byte("seq"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("1"))
	// clientVersion
	binary.Write(bodyBytes, binary.BigEndian, int32(len("clientversion")))
	bodyBytes.Write([]byte("clientversion"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(strconv.Itoa(clientVersion))))
	bodyBytes.Write([]byte(strconv.Itoa(clientVersion)))
	// clientostype
	binary.Write(bodyBytes, binary.BigEndian, int32(len("clientostype")))
	bodyBytes.Write([]byte("clientostype"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(deviceType)))
	bodyBytes.Write([]byte(deviceType))
	// authKey
	binary.Write(bodyBytes, binary.BigEndian, int32(len("authkey")))
	bodyBytes.Write([]byte("authkey"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(authKey)))
	bodyBytes.Write(authKey)
	// nettype
	binary.Write(bodyBytes, binary.BigEndian, int32(len("nettype")))
	bodyBytes.Write([]byte("nettype"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("1"))
	// acceptdupack
	binary.Write(bodyBytes, binary.BigEndian, int32(len("acceptdupack")))
	bodyBytes.Write([]byte("acceptdupack"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("1"))
	// sendHead := []byte{0x27, 0x14, 0x10, 0xa4, 0x65, 0x9a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	zHeader := new(bytes.Buffer)
	zHeader.Write([]byte{0xab})
	binary.Write(zHeader, binary.BigEndian, int32(bodyBytes.Len() + 25))
	zHeader.Write([]byte{0x27, 0x14})
	binary.Write(zHeader, binary.BigEndian, int32(uin))
	zHeader.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	binary.Write(zHeader, binary.BigEndian, int32(bodyBytes.Len()))
	zHeader.Write(bodyBytes.Bytes())
	return zHeader.Bytes()
}

func RequestImagePack(uin uint32, authKey []byte, clientVersion int, deviceType string, fileId string, aesKey string) []byte {
	bodyBytes := new(bytes.Buffer)
	// ver
	binary.Write(bodyBytes, binary.BigEndian, int32(len("ver")))
	bodyBytes.Write([]byte("ver"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("1"))
	// uin
	binary.Write(bodyBytes, binary.BigEndian, int32(len("weixinnum")))
	bodyBytes.Write([]byte("weixinnum"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(strconv.FormatUint(uint64(uin), 10))))
	bodyBytes.Write([]byte(strconv.FormatUint(uint64(uin), 10)))
	// seq
	binary.Write(bodyBytes, binary.BigEndian, int32(len("seq")))
	bodyBytes.Write([]byte("seq"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("2"))
	// clientVersion
	binary.Write(bodyBytes, binary.BigEndian, int32(len("clientversion")))
	bodyBytes.Write([]byte("clientversion"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(strconv.Itoa(clientVersion))))
	bodyBytes.Write([]byte(strconv.Itoa(clientVersion)))
	// clientostype
	binary.Write(bodyBytes, binary.BigEndian, int32(len("clientostype")))
	bodyBytes.Write([]byte("clientostype"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(deviceType)))
	bodyBytes.Write([]byte(deviceType))
	// authKey
	binary.Write(bodyBytes, binary.BigEndian, int32(len("authkey")))
	bodyBytes.Write([]byte("authkey"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(authKey)))
	bodyBytes.Write(authKey)
	// nettype
	binary.Write(bodyBytes, binary.BigEndian, int32(len("nettype")))
	bodyBytes.Write([]byte("nettype"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("1"))
	// acceptdupack
	binary.Write(bodyBytes, binary.BigEndian, int32(len("acceptdupack")))
	bodyBytes.Write([]byte("acceptdupack"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("1"))
	// rsaver
	binary.Write(bodyBytes, binary.BigEndian, int32(len("rsaver")))
	bodyBytes.Write([]byte("rsaver"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("1"))
	// filetype
	binary.Write(bodyBytes, binary.BigEndian, int32(len("filetype")))
	bodyBytes.Write([]byte("filetype"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("2")))
	bodyBytes.Write([]byte("2"))
	// wxchattype
	binary.Write(bodyBytes, binary.BigEndian, int32(len("wxchattype")))
	bodyBytes.Write([]byte("wxchattype"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("0")))
	bodyBytes.Write([]byte("0"))
	// fileid
	binary.Write(bodyBytes, binary.BigEndian, int32(len("fileid")))
	bodyBytes.Write([]byte("fileid"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(fileId)))
	bodyBytes.Write([]byte(fileId))
	// lastretcode
	binary.Write(bodyBytes, binary.BigEndian, int32(len("lastretcode")))
	bodyBytes.Write([]byte("lastretcode"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("0")))
	bodyBytes.Write([]byte("0"))
	// ipseq
	binary.Write(bodyBytes, binary.BigEndian, int32(len("ipseq")))
	bodyBytes.Write([]byte("ipseq"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("0")))
	bodyBytes.Write([]byte("0"))
	// wxmsgflag
	binary.Write(bodyBytes, binary.BigEndian, int32(len("wxmsgflag")))
	bodyBytes.Write([]byte("wxmsgflag"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("")))
	bodyBytes.Write([]byte(""))
	// wxautostart
	binary.Write(bodyBytes, binary.BigEndian, int32(len("wxautostart")))
	bodyBytes.Write([]byte("wxautostart"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("0")))
	bodyBytes.Write([]byte("0"))
	// downpicformat
	binary.Write(bodyBytes, binary.BigEndian, int32(len("downpicformat")))
	bodyBytes.Write([]byte("downpicformat"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("1")))
	bodyBytes.Write([]byte("1"))
	// offset
	binary.Write(bodyBytes, binary.BigEndian, int32(len("offset")))
	bodyBytes.Write([]byte("offset"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("0")))
	bodyBytes.Write([]byte("0"))
	// largesvideo
	binary.Write(bodyBytes, binary.BigEndian, int32(len("largesvideo")))
	bodyBytes.Write([]byte("largesvideo"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("0")))
	bodyBytes.Write([]byte("0"))
	// sourceflag
	binary.Write(bodyBytes, binary.BigEndian, int32(len("sourceflag")))
	bodyBytes.Write([]byte("sourceflag"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len("0")))
	bodyBytes.Write([]byte("0"))
	// rsavalue
	rsaKey := "BFEDFFB5EA28509F9C89ED83FA7FDDA8881435D444E984D53A98AD8E9410F1145EDD537890E10456190B22E6E5006455EFC6C12E41FDA985F38FBBC7213ECB810E3053D4B8D74FFBC70B4600ABD728202322AFCE1406046631261BD5EE3D44721082FEAB74340D73645DC0D02A293B962B9D47E4A64100BD7524DE00D9D3B5C1"
	aesKeyByte, _ := hex.DecodeString(aesKey)
	rsaValue := Algorithm.RSAEncrypt(aesKeyByte, rsaKey)
	binary.Write(bodyBytes, binary.BigEndian, int32(len("rsavalue")))
	bodyBytes.Write([]byte("rsavalue"))
	binary.Write(bodyBytes, binary.BigEndian, int32(len(rsaValue)))
	bodyBytes.Write(rsaValue)
	// sendHead := []byte{0x27, 0x14, 0x10, 0xa4, 0x65, 0x9a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	zHeader := new(bytes.Buffer)
	zHeader.Write([]byte{0xab})
	binary.Write(zHeader, binary.BigEndian, int32(bodyBytes.Len() + 25))
	zHeader.Write([]byte{0x4e, 0x20})
	binary.Write(zHeader, binary.BigEndian, int32(uin))
	zHeader.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	binary.Write(zHeader, binary.BigEndian, int32(bodyBytes.Len()))
	zHeader.Write(bodyBytes.Bytes())
	return zHeader.Bytes()
}

func CheckReceiveFinished(fileBytes []byte) (bool,[]byte) {
	for len(fileBytes) > 2 {
		headBytes := fileBytes[0:2]
		if headBytes[0] != 171 {
			fileBytes = fileBytes[2:]
		} else {
			var length int32
			readerHeader := bytes.NewReader(fileBytes[1:5])
			binary.Read(readerHeader, binary.BigEndian, &length)
			if length - 25 < int32(len(fileBytes)) {
				if length == 135 {
					continue
				}
				return true, fileBytes
			} else {
				return false, fileBytes
			}
		}
	}
	return false, fileBytes
}

func UnpackMessage(message []byte) map[string]string {
	var length int32
	readerHeader := bytes.NewReader(message[1:5])
	binary.Read(readerHeader, binary.BigEndian, &length)
	bodyBytes := message[25:length]
	res := make(map[string]string)
	index := 0
	for index < len(bodyBytes) {
		var keyLength int32
		keyLengthBytesReader := bytes.NewReader(bodyBytes[index : index + 4])
		index += 4
		binary.Read(keyLengthBytesReader, binary.BigEndian, &keyLength)
		key := string(bodyBytes[index : index + int(keyLength)])
		index += int(keyLength)
		var valueLength int32
		valueLengthBytesReader := bytes.NewReader(bodyBytes[index : index + 4])
		index += 4
		binary.Read(valueLengthBytesReader, binary.BigEndian, &valueLength)
		value := string(bodyBytes[index : index + int(valueLength)])
		res[key] = value
		index += int(valueLength)
	}
	return res
}