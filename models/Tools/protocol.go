package Tools

import (
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/lunny/log"
	"hash/crc32"
	"wechatdll/baseinfo"
	"wechatdll/comm"
	"wechatdll/models/baseutils"
)

// PackHeader 请求数据包头
type PackHeader struct {
	ReqData        []byte
	RetCode        int32
	Signature      byte
	HeadLength     byte
	CompressType   byte
	EncodeType     byte
	ServerVersion  uint32
	Uin            uint32
	Session        []byte
	SeqId          uint32
	URLID          uint32
	SrcLen         uint32
	ZipLen         uint32
	EncodeVersion  uint32
	HeadDeviceType byte
	CheckSum       uint32
	RunState       byte
	RqtCode        uint32
	EndFlag        byte
	Data           []byte
	HybridKeyVer   byte
}

func (p PackHeader) GetRetCode() int32 {
	return p.RetCode
}

// CalcHeadCheckSum 计算HeadCheckSum值
func CalcHeadCheckSum(uin uint32, checkSumKey []byte, srcData []byte) uint32 {
	uinBytes := baseutils.Int32ToBytes(uin)
	tmpBytes := append(uinBytes, checkSumKey[0:]...)
	md5Value := baseutils.Md5Value16(tmpBytes)

	dataLen := uint32(len(srcData))
	dataLenBytes := baseutils.Int32ToBytes(dataLen)
	tmpBytes = append(dataLenBytes, checkSumKey[0:]...)
	tmpBytes = append(tmpBytes, md5Value[0:]...)
	md5Value = baseutils.Md5Value16(tmpBytes)

	tmpBytes = append([]byte{}, md5Value[0:]...)
	// 计算返回
	/*tmpSum := baseutils.Adler32(1, tmpBytes)
	return baseutils.Adler32(tmpSum, srcData)*/
	return crc32.ChecksumIEEE(tmpBytes)
}

// CreatePackHead 创建包头
func CreatePackHead(userInfo *comm.LoginData, compressType byte, urlID uint32, srcData []byte, encodeData []byte, zipLen uint32, encodeType byte, encodeVersion uint32) *PackHeader {
	retHeader := &PackHeader{}

	// Signature
	retHeader.Signature = 0xbf
	retHeader.CompressType = compressType
	retHeader.EncodeType = encodeType << 4
	retHeader.ServerVersion = uint32(0x18001621)
	retHeader.Uin = userInfo.Uin
	retHeader.Session = userInfo.Cooike
	retHeader.URLID = urlID
	retHeader.SrcLen = uint32(len(srcData))
	retHeader.ZipLen = zipLen
	retHeader.EncodeVersion = encodeVersion
	retHeader.HeadDeviceType = baseinfo.MMHeadDeviceTypeIpadUniversal
	retHeader.CheckSum = 0x00
	// 如果有压缩，则计算Sum值
	if retHeader.CompressType == baseinfo.MMPackDataTypeCompressed {
		//retHeader.CheckSum = CalcHeadCheckSum(userInfo.Uin, userInfo.CheckSumKey, srcData)
	}

	retHeader.RunState = baseinfo.MMAppRunStateNormal
	retHeader.RqtCode = baseutils.CalcMsgCrcForString_807(baseutils.Md5ValueByte(encodeData, false))
	retHeader.EndFlag = 0x00
	retHeader.Data = encodeData

	return retHeader
}

// PackHeaderSerialize 序列化PackHeader
func PackHeaderSerialize(packHeader *PackHeader, needCookie bool) []byte {
	retBytes := make([]byte, 0)
	retBytes = append(retBytes, packHeader.Signature)
	retBytes = append(retBytes, 0)
	encodeType := packHeader.EncodeType
	if needCookie {
		packHeader.EncodeType = packHeader.EncodeType + 0xf
	}
	retBytes = append(retBytes, packHeader.EncodeType)
	retBytes = append(retBytes, baseutils.Int32ToBytes(packHeader.ServerVersion)[0:]...)
	retBytes = append(retBytes, baseutils.Int32ToBytes(packHeader.Uin)[0:]...)
	if needCookie {
		retBytes = append(retBytes, packHeader.Session[0:]...)
	}
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.URLID)[0:]...)
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.SrcLen)[0:]...)
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.ZipLen)[0:]...)
	// hybrid
	if encodeType>>4 == 12 {
		retBytes = append(retBytes, []byte{byte(packHeader.HybridKeyVer)}...)
	}
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.EncodeVersion)[0:]...)
	retBytes = append(retBytes, packHeader.HeadDeviceType)
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.CheckSum)[0:]...)
	retBytes = append(retBytes, packHeader.RunState)
	retBytes = append(retBytes, baseutils.EncodeVByte32(packHeader.RqtCode)[0:]...)
	retBytes = append(retBytes, packHeader.EndFlag)
	headLen := byte(len(retBytes))
	retBytes[1] = packHeader.CompressType + headLen<<2
	//log.Println(hex.EncodeToString(retBytes))
	retBytes = append(retBytes, packHeader.Data[0:]...)
	return retBytes
}

// Pack 打包加密数据
func Pack(userInfo *comm.LoginData, src []byte, urlID uint32, encodeType byte) []byte {
	retData := make([]byte, 0)
	if encodeType == 7 || encodeType == 17 {
		//加密类型7
		encodeData := src
		if encodeType == 7 {
			encodeData = baseutils.NoCompressRsaByVer(src, userInfo.GetLoginRsaVer())
		}
		packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, urlID, src, encodeData, uint32(len(src)), 7, userInfo.GetLoginRsaVer())
		retData = PackHeaderSerialize(packHeader, false)
	} else if encodeType == 5 {
		// 加密类型5
		zipBytes := baseutils.CompressByteArray(src)
		encodeData := baseutils.AesEncrypt(zipBytes, userInfo.Sessionkey)
		packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeCompressed, urlID, src, encodeData, uint32(len(zipBytes)), encodeType, 0)
		retData = PackHeaderSerialize(packHeader, true)
	} else if encodeType == 9 {
		// 加密类型9
		encodeData := src
		packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, urlID, src, encodeData, uint32(len(src)), encodeType, userInfo.GetLoginRsaVer())
		retData = PackHeaderSerialize(packHeader, true)
	} else if encodeType == 1 {
		// 加密类型1
		encodeData := baseutils.NoCompressRsaByVer(src, userInfo.GetLoginRsaVer())
		packHeader := CreatePackHead(userInfo, baseinfo.MMPackDataTypeUnCompressed, urlID, src, encodeData, uint32(len(src)), encodeType, userInfo.GetLoginRsaVer())
		retData = PackHeaderSerialize(packHeader, true)
	}
	return retData
}

// GetRespErrorCode 当response 小于32个字节时，调用这个接口获取响应的错误码
func GetRespErrorCode(data []byte) int32 {
	tmpData := make([]byte, 0)
	tmpData = append(tmpData, data[2:6]...)

	tmpRet := binary.BigEndian.Uint32(tmpData)
	return int32(tmpRet)
}

func DecodePackHeader(respData []byte, reqData []byte) (*PackHeader, error) {
	packHeader := &PackHeader{}
	packHeader.ReqData = reqData
	packHeader.RetCode = 0
	// 如果数据长度小于等于32, 则表明请求出错
	if len(respData) <= 32 {
		packHeader.RetCode = GetRespErrorCode(respData)
		return packHeader, errors.New("DecodePackHeader err: len(respData) <= 32")
	}

	current := 0
	packHeader.Signature = respData[current]
	current++
	packHeader.HeadLength = (respData[current]) >> 2
	packHeader.CompressType = (respData[current]) & 3
	current++
	packHeader.EncodeType = respData[current] >> 4
	sessionLen := int(respData[current] & 0x0f)
	current++
	packHeader.ServerVersion = baseutils.BytesToInt32(respData[current : current+4])
	current = current + 4
	packHeader.Uin = baseutils.BytesToInt32(respData[current : current+4])
	current = current + 4
	if sessionLen > 0 {
		packHeader.Session = respData[current : current+sessionLen]
		current = current + sessionLen
	}
	retLen := uint32(0)
	packHeader.URLID, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.SrcLen, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.ZipLen, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.EncodeVersion, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.HeadDeviceType = respData[current]
	current = current + 1
	packHeader.CheckSum, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.RunState = respData[current]
	current = current + 1
	packHeader.RqtCode, retLen = baseutils.DecodeVByte32(respData, uint32(current))
	current = current + int(retLen)
	packHeader.EndFlag = respData[current]
	current = current + 1
	// // 后面还有一个字节-- 可能是7.10新版本增加的一个字节，待后面分析
	// current = current + 1
	// if current != int(packHeader.HeadLength) {
	// 	return nil, errors.New("DecodePackHeader failed current != int(packHeader.HeadLength")
	// }
	packHeader.Data = respData[packHeader.HeadLength:]
	return packHeader, nil
}

// ParseResponseData 解析相应数据
func ParseResponseData(userInfo *comm.LoginData, packHeader *PackHeader, response proto.Message) error {
	//  判断包体长度是否大于0
	if len(packHeader.Data) <= 0 {
		log.Error("ParseResponseData err: len(packHeader.Data) <= 0")
		return errors.New("ParseResponseData err: len(packHeader.Data) <= 0")
	}
	var decptBody []byte
	var err error
	if packHeader.EncodeType == 12 {
		decptBody = packHeader.Data
	} else if packHeader.EncodeType == 5 {
		// 解密
		decptBody, err = baseutils.AesDecrypt(packHeader.Data, userInfo.Sessionkey)
		if err != nil {
			return err
		}
		// 判断是否有压缩
		if packHeader.CompressType == baseinfo.MMPackDataTypeCompressed {
			if decptBody != nil {
				//log.Println(hex.EncodeToString(decptBody))
				decptBody, err = baseutils.UnzipByteArray(decptBody)
				if err != nil {
					log.Error("ParseResponseData err:", err.Error(), packHeader.URLID)
					return err
				}
			} else {
				return errors.New("decptBody err: len(decptBody) == nil")
			}
		}
	} else {
		// 解密
		decptBody, err = baseutils.AesDecrypt(packHeader.Data, userInfo.Sessionkey)
		//log.Println(hex.EncodeToString(decptBody))
		if err != nil {
			//log.Error("ParseResponseData err:", err.Error())
			return err
		}
		// 判断是否有压缩
		if packHeader.CompressType == baseinfo.MMPackDataTypeCompressed {
			if decptBody != nil {
				decptBody, err = baseutils.UnzipByteArray(decptBody)
				if err != nil {
					//log.Println(hex.EncodeToString(packHeader.Data))
					log.Error("ParseResponseData err:", err.Error(), packHeader.URLID)
					return err
				}
			} else {
				return errors.New("decptBody err: len(decptBody) == nil")
			}
		}
	}

	// 更新UserInfo
	if len(packHeader.Session) > 6 {
		userInfo.Cooike = packHeader.Session
	}

	if packHeader.Uin != 0 {
		userInfo.Uin = packHeader.Uin
	}
	//log.Println(hex.EncodeToString(decptBody))
	// 解包ProtoBuf
	err = proto.Unmarshal(decptBody, response)
	if err != nil {
		log.Error("ParseResponseData err:", err.Error())
		return err
	}

	return nil
}
