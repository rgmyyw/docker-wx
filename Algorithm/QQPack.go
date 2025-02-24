package Algorithm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
	"wechatdll/lib"
)

// 微信appid
var WechatAppId uint32 = 0x1f1d5a7a
var WechatSerialStr string = "18c867f0717aa67b2ab7347505ba07ed"
var WechatPackage string = "com.tencent.mm"

func QQPackGroup1(clientVersion uint32, qqAccount uint32) []byte {
	bodyBuffer := new(bytes.Buffer)
	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x1))         	// 固定1
	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x0))         	// 固定0
	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x600))       	// 固定0
	binary.Write(bodyBuffer, binary.BigEndian, uint32(WechatAppId)) 	// 微信appid
	binary.Write(bodyBuffer, binary.BigEndian, uint32(clientVersion)) 	// ClientVersion
	binary.Write(bodyBuffer, binary.BigEndian, uint32(qqAccount)) 		// QQ号
	binary.Write(bodyBuffer, binary.BigEndian, uint32(0x0)) 			// QQ号
	bodyBytes := bodyBuffer.Bytes()
	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x18))         	// 固定18
	binary.Write(retBuffer, binary.BigEndian, uint16(len(bodyBytes)))   // 包长
	retBuffer.Write(bodyBytes)
	return retBuffer.Bytes()
}

func QQPackGroup2(qqAccount uint32) []byte {
	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x1))         				// 固定1
	binary.Write(retBuffer, binary.BigEndian, uint16(0x14))         			// 固定14
	binary.Write(retBuffer, binary.BigEndian, uint16(0x1))         				// 固定1
	randNum := rand.Float64() * 2.147483647E9
	binary.Write(retBuffer, binary.BigEndian, uint32(randNum))         			// 随机数
	binary.Write(retBuffer, binary.BigEndian, uint32(qqAccount)) 				// QQ号
	binary.Write(retBuffer, binary.BigEndian, uint32(time.Now().Unix())) 		// 时间戳
	binary.Write(retBuffer, binary.BigEndian, uint32(0x0)) 						// 原arr
	binary.Write(retBuffer, binary.BigEndian, uint16(0x0)) 						// 填充
	return retBuffer.Bytes()
}

func QQPackGroup3(clientVersion uint32, qqAccount uint32, timeSpan uint32, passwordMd5 []byte, encryptKey []byte, deviceIdStr string) []byte {
	plainBuffer := new(bytes.Buffer)
	binary.Write(plainBuffer, binary.BigEndian, uint16(0x2))         					// 包头
	randNum := rand.Float64() * 2.147483647E9
	binary.Write(plainBuffer, binary.BigEndian, uint32(randNum))         				// 随机数
	binary.Write(plainBuffer, binary.BigEndian, uint32(0x5))         					// 固定5
	binary.Write(plainBuffer, binary.BigEndian, uint32(WechatAppId)) 					// 微信appid
	binary.Write(plainBuffer, binary.BigEndian, uint32(clientVersion)) 					// ClientVersion
	binary.Write(plainBuffer, binary.BigEndian, uint32(0x0))         					// 固定0
	binary.Write(plainBuffer, binary.BigEndian, uint32(qqAccount)) 						// QQ号
	binary.Write(plainBuffer, binary.BigEndian, uint32(timeSpan)) 						// 时间戳
	binary.Write(plainBuffer, binary.BigEndian, uint32(0x0)) 							// 固定0
	binary.Write(plainBuffer, binary.BigEndian, uint8(0x1)) 							// 固定1
	plainBuffer.Write(passwordMd5)														// 密码MD5
	plainBuffer.Write(encryptKey)														// 密钥
	binary.Write(plainBuffer, binary.BigEndian, uint32(0x0)) 							// 固定0
	binary.Write(plainBuffer, binary.BigEndian, uint8(0x1)) 							// 固定1
	plainBuffer.Write([]byte(deviceIdStr))												// deviceId的字符串编码
	binary.Write(plainBuffer, binary.BigEndian, uint16(0x0)) 							// 固定0
	plainBytes := plainBuffer.Bytes()

	keyBuffer := new(bytes.Buffer)
	keyBuffer.Write(passwordMd5)														// 密码MD5
	binary.Write(keyBuffer, binary.BigEndian, uint32(0x0)) 								// 固定0
	binary.Write(keyBuffer, binary.BigEndian, uint32(qqAccount)) 						// QQ号
	keyBytes := keyBuffer.Bytes()
	keyBytesMd5 := lib.Md5Hash(keyBytes)												// 组包3密钥

	qqCryptor := GetQQCryptor()
	encBytes := qqCryptor.Encrypt(plainBytes, keyBytesMd5)								// 加密

	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x106))         					// 固定0x147
	binary.Write(retBuffer, binary.BigEndian, uint16(len(encBytes)))   					// 包长
	retBuffer.Write(encBytes)

	return retBuffer.Bytes()
}

func QQPackGroup4(clientVersion uint32) []byte {
	bodyBuffer := new(bytes.Buffer)
	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x1))         	// 固定1
	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x0))         	// 固定0
	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x5))       		// 固定5
	binary.Write(bodyBuffer, binary.BigEndian, uint32(WechatAppId)) 	// 微信appid
	binary.Write(bodyBuffer, binary.BigEndian, uint32(0x1)) 			// 固定1
	binary.Write(bodyBuffer, binary.BigEndian, uint32(clientVersion)) 	// ClientVersion
	binary.Write(bodyBuffer, binary.BigEndian, uint32(0x2040)) 			// 固定8256
	bodyBytes := bodyBuffer.Bytes()
	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x100))         	// 固定18
	binary.Write(retBuffer, binary.BigEndian, uint16(len(bodyBytes)))   // 包长
	retBuffer.Write(bodyBytes)
	return retBuffer.Bytes()
}

func QQPackGroup5() []byte {
	retByte, _ := hex.DecodeString("01070006000001900001")
	return retByte
}

func QQPackGroup6() []byte {
	retByte, _ := hex.DecodeString("0116000a000000057c0001040000")
	return retByte
}

func QQPackGroup7(deviceIdStr string) []byte {
	bodyBuffer := new(bytes.Buffer)
	bodyBuffer.Write([]byte(deviceIdStr))												// deviceId的字符串编码
	bodyBytes := bodyBuffer.Bytes()
	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x145))         					// 固定18
	binary.Write(retBuffer, binary.BigEndian, uint16(len(bodyBytes)))   				// 包长
	retBuffer.Write(bodyBytes)
	return retBuffer.Bytes()
}

func QQPackGroup8(ksid []byte) []byte {
	bodyBuffer := new(bytes.Buffer)
	bodyBuffer.Write(ksid)																// 16位ksid: fb7587cd635cf1bed36b868518d7e077
	bodyBytes := bodyBuffer.Bytes()
	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x108))         					// 固定18
	binary.Write(retBuffer, binary.BigEndian, uint16(len(bodyBytes)))   				// 包长
	retBuffer.Write(bodyBytes)
	return retBuffer.Bytes()
}

func QQPackGroup9(deviceIdStr string) []byte {
	bodyBuffer := new(bytes.Buffer)
	bodyBuffer.Write([]byte(deviceIdStr))												// deviceId的字符串编码
	bodyBytes := bodyBuffer.Bytes()
	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x109))         					// 固定18
	binary.Write(retBuffer, binary.BigEndian, uint16(len(bodyBytes)))   				// 包长
	retBuffer.Write(bodyBytes)
	return retBuffer.Bytes()
}

func QQPackGroup10Sub1(mSystem string, mSystemVersion string, ispName string, ispType string) []byte {
	bodyBuffer := new(bytes.Buffer)
	mSystemBytes := []byte(mSystem)														// android字符串ascii码
	binary.Write(bodyBuffer, binary.BigEndian, uint16(len(mSystemBytes)))
	bodyBuffer.Write(mSystemBytes)
	mSystemVersionByte := []byte(mSystemVersion)										// 系统版本10
	binary.Write(bodyBuffer, binary.BigEndian, uint16(len(mSystemVersionByte)))
	bodyBuffer.Write(mSystemVersionByte)
	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x2))								// 序号2
	ispNameBytes := []byte(ispName)														// "中国移动"编码
	binary.Write(bodyBuffer, binary.BigEndian, uint16(len(ispNameBytes)))
	bodyBuffer.Write(ispNameBytes)
	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x0))         					// 固定0
	ispTypeBytes := []byte(ispType)														// wifi字符串的ascii码
	binary.Write(bodyBuffer, binary.BigEndian, uint16(len(ispTypeBytes)))
	bodyBuffer.Write(ispTypeBytes)
	bodyBytes := bodyBuffer.Bytes()

	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x124))         					// 固定0x124
	binary.Write(retBuffer, binary.BigEndian, uint16(len(bodyBytes)))   				// 包长
	retBuffer.Write(bodyBytes)
	return retBuffer.Bytes()
}

func QQPackGroup10Sub2(mBrand string, deviceIdStr string) []byte {
	bodyBuffer := new(bytes.Buffer)
	binary.Write(bodyBuffer, binary.BigEndian, uint32(0x1))								// 固定1

	binary.Write(bodyBuffer, binary.BigEndian, uint32(0x0))								// 固定0

	binary.Write(bodyBuffer, binary.BigEndian, uint8(0x0))								// 固定0

	mBrandBytes := []byte(mBrand)														// 系统品牌ascii码"MI 8"
	binary.Write(bodyBuffer, binary.BigEndian, uint16(len(mBrandBytes)))
	bodyBuffer.Write(mBrandBytes)

	deviceIdStrByte := []byte(deviceIdStr)												// 大写DeviceId的Ascii码
	binary.Write(bodyBuffer, binary.BigEndian, uint16(len(deviceIdStrByte)))
	bodyBuffer.Write(deviceIdStrByte)

	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x0))         					// 固定0
	bodyBytes := bodyBuffer.Bytes()

	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x128))         					// 固定0x128
	binary.Write(retBuffer, binary.BigEndian, uint16(len(bodyBytes)))   				// 包长
	retBuffer.Write(bodyBytes)
	return retBuffer.Bytes()
}

func QQPackGroup10Sub3(wxVersionStr string) []byte {
	bodyBuffer := new(bytes.Buffer)
	binary.Write(bodyBuffer, binary.BigEndian, uint32(WechatAppId))								// appid

	mBrandBytes := []byte(wxVersionStr)															// 微信版本号Ascii "8.0.9"
	binary.Write(bodyBuffer, binary.BigEndian, uint16(len(mBrandBytes)))
	bodyBuffer.Write(mBrandBytes)

	WechatSerial, _ := hex.DecodeString(WechatSerialStr)										// 固定微信包号
	binary.Write(bodyBuffer, binary.BigEndian, uint16(len(WechatSerial)))
	bodyBuffer.Write(WechatSerial)

	bodyBytes := bodyBuffer.Bytes()

	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x147))         					// 固定0x147
	binary.Write(retBuffer, binary.BigEndian, uint16(len(bodyBytes)))   				// 包长
	retBuffer.Write(bodyBytes)
	return retBuffer.Bytes()
}

func QQPackGroup10(deviceIdStr string, mSystem string, mSystemVersion string, ispName string, ispType string, mBrand string, wxVersionStr string, encryptKey []byte) []byte {
	packGroup9 := QQPackGroup9(deviceIdStr)
	packGroup101 := QQPackGroup10Sub1(mSystem, mSystemVersion, ispName, ispType)
	packGroup102 := QQPackGroup10Sub2(mBrand, deviceIdStr)
	packGroup103 := QQPackGroup10Sub3(wxVersionStr)

	plainBuffer := new(bytes.Buffer)
	binary.Write(plainBuffer, binary.BigEndian, uint16(0x4))         					// 4个包
	plainBuffer.Write(packGroup9)
	plainBuffer.Write(packGroup101)
	plainBuffer.Write(packGroup102)
	plainBuffer.Write(packGroup103)
	plainBytes := plainBuffer.Bytes()

	qqCryptor := GetQQCryptor()
	encryptText := qqCryptor.Encrypt(plainBytes, encryptKey)

	group10Buffer := new(bytes.Buffer)
	binary.Write(group10Buffer, binary.BigEndian, uint16(0x144))         					// 固定144
	binary.Write(group10Buffer, binary.BigEndian, uint16(len(encryptText)))         		// 长度
	group10Buffer.Write(encryptText)														// 密文

	return group10Buffer.Bytes()
}

func QQPackGroup11() []byte {
	bodyBuffer := new(bytes.Buffer)
	binary.Write(bodyBuffer, binary.BigEndian, uint16(0x0))							// 固定0

	mBrandBytes := []byte(WechatPackage)											// 'com.tencent.mm'的ASCII码
	binary.Write(bodyBuffer, binary.BigEndian, uint16(len(mBrandBytes)))
	bodyBuffer.Write(mBrandBytes)

	bodyBytes := bodyBuffer.Bytes()

	retBuffer := new(bytes.Buffer)
	binary.Write(retBuffer, binary.BigEndian, uint16(0x142))         					// 固定0x142
	binary.Write(retBuffer, binary.BigEndian, uint16(len(bodyBytes)))   				// 包长
	retBuffer.Write(bodyBytes)
	return retBuffer.Bytes()
}

type QQmobileRequest struct{
	ClientVersion			uint32			// 微信版本号
	Account					uint32			// QQ号
	Password				string			// 密码
	DeviceIdStr				string			// deviceId
	Ksid					[]byte			// 16位ksid, 没有就占位 []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
	TimeSpan				uint32			// 发起时间戳
	MobileSystem			string			// "android"
	MobileSystemVersion		string			// "10"
	ISPName					string			// 运营商: "中国移动"
	ISPType					string			// "wifi"
	MobileBrand				string			// "MI 8"
	WechatVersion			string			// "8.0.9"
}

func QQPackBody(req QQmobileRequest) ([]byte, []byte) {
	packGroup1 := QQPackGroup1(req.ClientVersion, req.Account)
	packGroup2 := QQPackGroup2(req.Account)
	passwordMd5 := lib.Md5Hash([]byte(req.Password))
	encryptKey := []byte(lib.RandSeq(16))
	packGroup3 := QQPackGroup3(req.ClientVersion, req.Account, req.TimeSpan, passwordMd5, encryptKey, req.DeviceIdStr)
	packGroup6 := QQPackGroup6()
	packGroup4 := QQPackGroup4(req.ClientVersion)
	packGroup5 := QQPackGroup5()
	//packGroup8 := QQPackGroup8(req.Ksid)
	packGroup10 := QQPackGroup10(req.DeviceIdStr, req.MobileSystem, req.MobileSystemVersion, req.ISPName, req.ISPType,
		req.MobileBrand, req.WechatVersion, encryptKey)
	packGroup11 := QQPackGroup11()
	packGroup7 := QQPackGroup7(req.DeviceIdStr)


	plainBuffer := new(bytes.Buffer)
	binary.Write(plainBuffer, binary.BigEndian, uint16(0x9))         					// 固定9
	binary.Write(plainBuffer, binary.BigEndian, uint16(0x9))         					// 10个包
	plainBuffer.Write(packGroup1)
	plainBuffer.Write(packGroup2)
	plainBuffer.Write(packGroup3)
	plainBuffer.Write(packGroup6)
	plainBuffer.Write(packGroup4)
	plainBuffer.Write(packGroup5)
	//plainBuffer.Write(packGroup8)
	plainBuffer.Write(packGroup10)
	plainBuffer.Write(packGroup11)
	plainBuffer.Write(packGroup7)
	plainBytes := plainBuffer.Bytes()
	fmt.Printf("总包体: %x\n", plainBytes)
	bodyEncKey := []byte(lib.RandSeq(16))

	qqCryptor := GetQQCryptor()
	encryptByte := qqCryptor.Encrypt(plainBytes, bodyEncKey)

	bodyBuffer := new(bytes.Buffer)
	bodyBuffer.Write(bodyEncKey)
	bodyBuffer.Write(encryptByte)
	return bodyBuffer.Bytes(), bodyEncKey
}

func QQPackMessage(req QQmobileRequest) []byte {
	bodyBytes, _ := QQPackBody(req)

	headerBuffer := new(bytes.Buffer)
	binary.Write(headerBuffer, binary.BigEndian, uint16(0x1f41))								// 固定8001
	binary.Write(headerBuffer, binary.BigEndian, uint16(0x810))									// 固定810
	binary.Write(headerBuffer, binary.BigEndian, uint16(0x1))									// 固定1
	binary.Write(headerBuffer, binary.BigEndian, uint32(req.Account))							// qq号

	binary.Write(headerBuffer, binary.BigEndian, uint8(0x3))									// 固定3
	binary.Write(headerBuffer, binary.BigEndian, uint16(0x0))									// 固定0
	binary.Write(headerBuffer, binary.BigEndian, uint16(0x0))									// 固定0
	binary.Write(headerBuffer, binary.BigEndian, uint16(0x0))									// 固定0

	binary.Write(headerBuffer, binary.BigEndian, uint32(req.ClientVersion))						// 微信版本
	binary.Write(headerBuffer, binary.BigEndian, uint32(0x0))									// 固定0
	headerBytes := headerBuffer.Bytes()

	msgBuffer := new(bytes.Buffer)
	binary.Write(msgBuffer, binary.BigEndian, uint8(0x2))										// 固定2
	binary.Write(msgBuffer, binary.BigEndian, uint16(len(headerBytes) + len(bodyBytes) + 4))	// 总长度
	msgBuffer.Write(headerBytes)
	msgBuffer.Write(bodyBytes)
	binary.Write(msgBuffer, binary.BigEndian, uint8(0x3))										// 固定3
	return msgBuffer.Bytes()
}