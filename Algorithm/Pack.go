package Algorithm

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"hash/adler32"
	"hash/crc32"
	"io"
	"wechatdll/lib"
)

func (h *Client) HybridEcdhPackIosEn(Cgi, Uin uint32, Cookies, Data []byte) []byte {
	header := new(bytes.Buffer)
	header.Write([]byte{0xbf})
	header.Write([]byte{0x02}) //加密模式占坑,默认不压缩走12

	encryptdata := h.encryptoIOS(Data)
	// log.Println(hex.EncodeToString(encryptdata))

	cookielen := len(Cookies)
	header.Write([]byte{byte((12 << 4) + cookielen)})
	binary.Write(header, binary.BigEndian, int32(h.Version))
	if Uin != 0 {
		binary.Write(header, binary.BigEndian, int32(Uin))
	} else {
		header.Write([]byte{0x00, 0x00, 0x00, 0x00})
	}

	if len(Cookies) == 0xF {
		header.Write(Cookies)
	}
	//log.Infof("计算RQT数据: %s",hex.EncodeToString(encryptdata))
	//rqt := RQT(encryptdata)
	rqtx := CalcMsgCrcForData_7019(encryptdata)
	//log.Infof("原RQT: %v", rqt)
	//log.Infof("新RQT: %v", rqtx)
	header.Write(proto.EncodeVarint(uint64(Cgi)))
	header.Write(proto.EncodeVarint(uint64(len(Data))))
	header.Write(proto.EncodeVarint(uint64(len(Data))))
	header.Write([]byte{0x90, 0x4E, 0x0D, 0x00, 0xFF})
	header.Write(proto.EncodeVarint(uint64(rqtx)))
	header.Write([]byte{0x00})
	lens := len(header.Bytes())<<2 + 2
	header.Bytes()[1] = byte(lens)
	header.Write(encryptdata)
	// log.Println(hex.EncodeToString(header.Bytes()))
	return header.Bytes()
}
func GenSignature(uiCryptin uint32, salt, data []byte) uint32 {
	var b1 bytes.Buffer
	binary.Write(&b1, binary.BigEndian, uiCryptin)
	h1 := md5.New()
	h1.Write(b1.Bytes())
	h1.Write(salt)
	sum1 := h1.Sum(nil)

	dataSize := len(data)
	var b2 bytes.Buffer
	binary.Write(&b2, binary.BigEndian, uint32(dataSize))

	h2 := md5.New()
	h2.Write(b2.Bytes())
	h2.Write(salt)
	h2.Write(sum1)
	sum2 := h2.Sum(nil)

	a := adler32.New()
	a.Write(nil)
	a.Write(sum2)
	a.Write(data)
	return a.Sum32()
}
func (h *Client) HybridEcdhPackIosEn2(Cgi, Uin uint32, Cookies, Data, loginecdhkey []byte) []byte {
	header := new(bytes.Buffer)
	header.Write([]byte{0xbf})
	header.Write([]byte{0x02}) //加密模式占坑,默认不压缩走12

	encryptdata := h.encryptoIOS(Data)

	cookielen := len(Cookies)
	header.Write([]byte{byte((12 << 4) + cookielen)})
	binary.Write(header, binary.BigEndian, int32(h.Version))
	if Uin != 0 {
		binary.Write(header, binary.BigEndian, int32(Uin))
	} else {
		header.Write([]byte{0x00, 0x00, 0x00, 0x00})
	}

	if len(Cookies) == 0xF {
		header.Write(Cookies)
	}

	header.Write(proto.EncodeVarint(uint64(Cgi)))
	header.Write(proto.EncodeVarint(uint64(len(encryptdata))))
	header.Write(proto.EncodeVarint(uint64(len(encryptdata))))
	header.Write(proto.EncodeVarint(10003))
	header.Write([]byte{0x00})
	header.Write(proto.EncodeVarint(uint64(GenSignature(Uin, loginecdhkey, Data))))
	header.Write([]byte{0xff})
	header.Write(proto.EncodeVarint(uint64(CalcMsgCrcForData_7019(encryptdata))))
	header.Write([]byte{0x00})
	lens := len(header.Bytes())<<2 + 2
	header.Bytes()[1] = byte(lens)
	header.Write(encryptdata)
	return header.Bytes()
}
func (h *Client) HybridEcdhPackIosUn(Data []byte) *PacketHeader {
	var ph PacketHeader
	var body []byte
	var nCur int64
	var bfbit byte
	srcreader := bytes.NewReader(Data)
	binary.Read(srcreader, binary.BigEndian, &bfbit)
	if bfbit == byte(0xbf) {
		nCur += 1
	}
	nLenHeader := Data[nCur] >> 2
	nCur += 1
	nLenCookie := Data[nCur] & 0xf
	nCur += 1
	nCur += 4
	srcreader.Seek(nCur, io.SeekStart)
	binary.Read(srcreader, binary.BigEndian, &ph.Uin)
	nCur += 4
	cookie_temp := Data[nCur : nCur+int64(nLenCookie)]
	ph.Cookies = cookie_temp
	nCur += int64(nLenCookie)
	cgidata := Data[nCur:]
	_, nSize := proto.DecodeVarint(cgidata)
	nCur += int64(nSize)
	LenProtobufData := Data[nCur:]
	_, nLenProtobuf := proto.DecodeVarint(LenProtobufData)
	nCur += int64(nLenProtobuf)
	body = Data[nLenHeader:]
	protobufdata := h.decryptoIOS(body)
	ph.Data = protobufdata
	return &ph
}

func (h *Client) HybridEcdhPackAndroidEn(cmdid, cert, uin uint32, cookie, Data []byte) []byte {
	EnData := h.encryptAndroid(Data)
	//log.Infof("计算RQT数据: %s",hex.EncodeToString(EnData))
	//rqt := RQT(EnData)
	rqtx := CalcMsgCrcForData_7019(EnData)
	//log.Infof("原RQT: %v", rqt)
	//log.Infof("新RQT: %v", rqtx)
	inputlen := len(EnData)
	pack := append([]byte{}, cookie...)
	pack = proto.EncodeVarint(uint64(cmdid))
	pack = append(pack, proto.EncodeVarint(uint64(inputlen))...)
	pack = append(pack, proto.EncodeVarint(uint64(inputlen))...)
	pack = append(pack, proto.EncodeVarint(uint64(cert))...)
	pack = append(pack, 2)
	pack = append(pack, 0)
	pack = append(pack, 0xfe)
	pack = append(pack, proto.EncodeVarint(uint64(rqtx))...)
	pack = append(pack, 0)
	headLen := len(pack) + 11
	headFlag := (12 << 12) | (len(cookie) << 8) | (headLen << 2) | 2
	var hybridpack = new(bytes.Buffer)
	hybridpack.WriteByte(0xbf)
	binary.Write(hybridpack, binary.LittleEndian, uint16(headFlag))
	binary.Write(hybridpack, binary.BigEndian, uint32(h.Version))
	binary.Write(hybridpack, binary.BigEndian, uint32(uin))
	hybridpack.Write(pack)
	hybridpack.Write(EnData)
	return hybridpack.Bytes()
}

func (h *Client) HybridEcdhPackAndroidUn(Data []byte) *PacketHeader {
	var ph PacketHeader
	readHeader := bytes.NewReader(Data)
	binary.Read(readHeader, binary.LittleEndian, &ph.PacketCryptType)
	binary.Read(readHeader, binary.LittleEndian, &ph.Flag)
	cookieLen := (ph.Flag >> 8) & 0x0f
	headerLen := (ph.Flag & 0xff) >> 2
	ph.Cookies = make([]byte, cookieLen)
	binary.Read(readHeader, binary.BigEndian, &ph.RetCode)
	binary.Read(readHeader, binary.BigEndian, &ph.UICrypt)
	binary.Read(readHeader, binary.LittleEndian, &ph.Cookies)
	ph.Data = h.decryptAndroid(Data[headerLen:])
	return &ph
}

func Pack(src []byte, cgi int, uin uint32, sessionkey, cookies, clientsessionkey, loginecdhkey []byte, encryptType uint8, use_compress bool) []byte {
	len_proto_compressed := len(src)
	var body []byte
	if use_compress {
		if cgi == 138 {
			encryptType = 13
			mNonce := []byte(lib.RandSeq(12)) //获取随机密钥
			body = AesGcmEncryptWithCompressZlib(clientsessionkey, src, mNonce, nil)
		} else {
			body = CompressAndAes(src, sessionkey)
		}
	} else {
		if cgi == 138 {
			encryptType = 13
			mNonce := []byte(lib.RandSeq(12)) //获取随机密钥
			body = AesGcmEncryptWithCompressZlib(clientsessionkey, src, mNonce, nil)
		} else {
			body = AesEncrypt(src, sessionkey)
		}
	}

	loginecdhkeylen := int32(len(loginecdhkey))

	header := new(bytes.Buffer)
	header.Write([]byte{0xbf})
	header.Write([]byte{0x00})
	header.Write([]byte{((encryptType << 4) + 0xf)})
	binary.Write(header, binary.BigEndian, int32(IPhoneVersion))
	binary.Write(header, binary.BigEndian, int32(uin))
	header.Write(cookies)
	header.Write(proto.EncodeVarint(uint64(cgi)))

	if use_compress {
		header.Write(proto.EncodeVarint(uint64(len_proto_compressed)))
		header.Write(proto.EncodeVarint(uint64(len(body))))
	} else {
		header.Write(proto.EncodeVarint(uint64(len_proto_compressed)))
		header.Write(proto.EncodeVarint(uint64(len_proto_compressed)))
	}

	header.Write([]byte{0x00, 0x0d}) //占坑
	uinbyte := new(bytes.Buffer)
	binary.Write(uinbyte, binary.BigEndian, uin)
	m1 := md5.New()
	m1.Write(uinbyte.Bytes())
	m1.Write(loginecdhkey[:loginecdhkeylen])
	md5str := m1.Sum(nil)

	lenprotobuf := new(bytes.Buffer)
	binary.Write(lenprotobuf, binary.BigEndian, int32(len(src)))
	m2 := md5.New()
	m2.Write(lenprotobuf.Bytes())
	m2.Write(loginecdhkey[:loginecdhkeylen])
	m2.Write(md5str)

	//log.Infof("计算RQT数据: %s",hex.EncodeToString(body))
	//rqt := RQT(body)
	rqtx := CalcMsgCrcForData_7019(body)
	//log.Infof("原RQT: %v", rqt)
	//log.Infof("新RQT: %v", rqtx)

	md5str = m2.Sum(nil)
	adler32buffer := new(bytes.Buffer)
	adler32buffer.Write(md5str)
	adler32buffer.Write(src)
	//header.Write(proto.EncodeVarint(uint64(LOGIN_RSA_VER)))
	adler32 := crc32.ChecksumIEEE(adler32buffer.Bytes())
	header.Write(proto.EncodeVarint(uint64(adler32)))
	header.Write([]byte{0xFF})                     //占坑
	header.Write(proto.EncodeVarint(uint64(rqtx))) //占坑
	header.Write([]byte{0x00})                     //占坑
	if use_compress {
		lens := (len(header.Bytes()) << 2) + 1
		header.Bytes()[1] = byte(lens)
	} else {
		lens := (len(header.Bytes()) << 2) + 2
		header.Bytes()[1] = byte(lens)
	}
	header.Write(body)
	return header.Bytes()
}

func UnpackBusinessPacket(src []byte, key []byte, uin uint32, cookie *[]byte) []byte {
	var nCur int64
	var bfbit byte
	srcreader := bytes.NewReader(src)
	binary.Read(srcreader, binary.BigEndian, &bfbit)
	if bfbit == byte(0xbf) {
		nCur += 1
	}
	nLenHeader := src[nCur] >> 2
	bUseCompressed := src[nCur] & 0x3
	nCur += 1
	nLenCookie := src[nCur] & 0xf
	nCur += 1
	nCur += 4
	srcreader.Seek(nCur, io.SeekStart)
	binary.Read(srcreader, binary.BigEndian, &uin)
	nCur += 4
	cookie_temp := src[nCur : nCur+int64(nLenCookie)]
	*cookie = cookie_temp
	nCur += int64(nLenCookie)
	cgidata := src[nCur:]
	_, nSize := proto.DecodeVarint(cgidata)
	nCur += int64(nSize)
	LenProtobufData := src[nCur:]
	_, nLenProtobuf := proto.DecodeVarint(LenProtobufData)
	nCur += int64(nLenProtobuf)
	body := src[nLenHeader:]
	if bUseCompressed == 1 {
		protobufData := DecompressAndAesDecrypt(body, key)
		return protobufData
	} else {
		protobufData := AesDecrypt(body, key)
		return protobufData
	}
}

func UnpackBusinessPacketWithAesGcm(src []byte, uin uint32, cookie *[]byte, Serversessionkey []byte) []byte {
	var nCur int64
	var bfbit byte
	srcreader := bytes.NewReader(src)
	binary.Read(srcreader, binary.BigEndian, &bfbit)
	if bfbit == byte(0xbf) {
		nCur += 1
	}
	nLenHeader := src[nCur] >> 2
	nCur += 1
	nLenCookie := src[nCur] & 0xf
	nCur += 1
	nCur += 4
	srcreader.Seek(nCur, io.SeekStart)
	binary.Read(srcreader, binary.BigEndian, &uin)
	nCur += 4
	cookie_temp := src[nCur : nCur+int64(nLenCookie)]
	*cookie = cookie_temp
	nCur += int64(nLenCookie)
	cgidata := src[nCur:]
	_, nSize := proto.DecodeVarint(cgidata)
	nCur += int64(nSize)
	LenProtobufData := src[nCur:]
	_, nLenProtobuf := proto.DecodeVarint(LenProtobufData)
	nCur += int64(nLenProtobuf)
	body := src[nLenHeader:]
	protobufdata := AesGcmDecryptWithcompressZlib(Serversessionkey, body, nil)
	return protobufdata
}

func EncodePackMini05(requestBytes []byte, cmdId int, uin uint32, publicKey []byte, sessionKey []byte, cookie []byte, clientVersion int, compress bool, useRqt bool) []byte {
	// 只处理AES常规组包, encryptType := 5
	zipBytes := requestBytes
	if compress {
		zipBytes = DoZlibCompress(requestBytes)
		log.Infof("238心跳zip压缩: %x", zipBytes)
	}

	aesBytes := AesEncrypt(zipBytes, sessionKey)
	log.Infof("238心跳aes加密: %x", aesBytes)
	rqt := uint32(0)
	if useRqt {
		rqt = CalcMsgCrcForData_7019(aesBytes)
	}
	log.Infof("238心跳rqt计算: %v", rqt)
	check := MmtlsAlder(publicKey, requestBytes, uin)
	log.Infof("238心跳check计算: %v", check)
	aesWrapper := CommonRequestPack(len(requestBytes), len(zipBytes), aesBytes, uin, clientVersion, cmdId, cookie, check, rqt, compress)

	return aesWrapper
}

func CommonRequestPack(pbLen int, zipLen int, aesBytes []byte, uin uint32, clientVersion int, cmdId int, cookie []byte, check int, rqt uint32, compress bool) []byte {
	// aes组包
	byteUin := []byte{
		uint8((uin&0xff000000)>>24) & 0xff,
		uint8((uin&0x00ff0000)>>16) & 0xff,
		uint8((uin&0x0000ff00)>>8) & 0xff,
		uint8((uin & 0x000000ff) & 0xff),
	}
	packet := new(bytes.Buffer)
	packet.Write([]byte{0xbf, 0x62, 0x50})
	binary.Write(packet, binary.BigEndian, uint32(clientVersion))
	packet.Write(byteUin)
	packet.Write(cookie)
	packet.Write(proto.EncodeVarint(uint64(cmdId)))
	packet.Write(proto.EncodeVarint(uint64(pbLen)))
	packet.Write(proto.EncodeVarint(uint64(zipLen)))
	packet.Write([]byte{0x00, 0x02})
	packet.Write(proto.EncodeVarint(uint64(check)))
	packet.Write([]byte{0x02})
	packet.Write(proto.EncodeVarint(uint64(rqt)))
	packet.Write([]byte{0x00})
	if compress {
		packet.Bytes()[1] = byte((packet.Len() << 2) + 1)
	} else {
		packet.Bytes()[1] = byte((packet.Len() << 2) + 2)
	}
	packet.Bytes()[2] = byte(0x50 + len(cookie))
	packet.Write(aesBytes)
	return packet.Bytes()
}

func MmtlsAlder(publicKey []byte, data []byte, uin uint32) int {
	num := len(data)
	array2 := []byte{
		uint8(((uin & 0xff000000) >> 24) & 0xff),
		uint8(((uin & 0xff0000) >> 16) & 0xff),
		uint8(((uin & 0xff00) >> 8) & 0xff),
		uint8(uin & 0xff & 0xff),
	}
	firstHashBytes := append(array2, publicKey...)
	second := lib.Md5Hash(firstHashBytes)
	array3 := []byte{
		uint8(((num & 0xff000000) >> 24) & 0xff),
		uint8(((num & 0xff0000) >> 16) & 0xff),
		uint8(((num & 0xff00) >> 8) & 0xff),
		uint8(num & 0xff & 0xff),
	}
	secondHashBytes := append(array3, publicKey...)
	secondHashBytes = append(secondHashBytes, second...)
	second = lib.Md5Hash(secondHashBytes)
	crc := ComputeAlder32Hash(second, 1)
	crc = ComputeAlder32Hash(data, crc)
	return crc
}

func ComputeAlder32Hash(data []byte, checkSum int) int {
	if checkSum == 0 {
		checkSum = 1
	}
	s1 := checkSum & 0xFFFF
	s2 := checkSum >> 16
	bytesToRead := len(data)
	byteStart := 0
	for bytesToRead > 0 {
		n := 3800
		if bytesToRead < 3800 {
			n = bytesToRead
		}
		bytesToRead -= n
		for n > 0 {
			n--
			s1 = s1 + int(data[byteStart]&0xFF)
			byteStart++
			s2 += s1
		}
		s1 %= 65521
		s2 %= 65521
	}
	checkSum = (s2 << 16) | s1
	return checkSum
}
