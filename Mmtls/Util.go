package Mmtls

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"net/http"
	"unicode/utf8"
)

type MmtlsClient struct {
	Shakehandpubkey    []byte
	Shakehandpubkeylen int32
	Shakehandprikey    []byte
	Shakehandprikeylen int32

	Shakehandpubkey_2   []byte
	Shakehandpubkeylen2 int32
	Shakehandprikey_2   []byte
	Shakehandprikeylen2 int32

	Mserverpubhashs     []byte
	ServerSeq           int
	ClientSeq           int
	ShakehandECDHkey    []byte
	ShakehandECDHkeyLen int32

	Encrptmmtlskey  []byte
	Decryptmmtlskey []byte
	EncrptmmtlsIv   []byte
	DecryptmmtlsIv  []byte

	CurDecryptSeqIv []byte
	CurEncryptSeqIv []byte

	Decrypt_part2_hash256            []byte
	Decrypt_part3_hash256            []byte
	ShakehandECDHkeyhash             []byte
	Hkdfexpand_pskaccess_key         []byte
	Hkdfexpand_pskrefresh_key        []byte
	HkdfExpand_info_serverfinish_key []byte
	Hkdfexpand_clientfinish_key      []byte
	Hkdfexpand_secret_key            []byte

	Hkdfexpand_application_key []byte
	Encrptmmtlsapplicationkey  []byte
	Decryptmmtlsapplicationkey []byte
	EncrptmmtlsapplicationIv   []byte
	DecryptmmtlsapplicationIv  []byte

	Earlydatapart       []byte
	Newsendbufferhashs  []byte
	Encrptshortmmtlskey []byte
	Encrptshortmmtlsiv  []byte
	Decrptshortmmtlskey []byte
	Decrptshortmmtlsiv  []byte

	//http才需要
	Pskkey    string
	Pskiv     string
	MmtlsMode uint
}

//短连接mmtls模块
type HttpClientModel struct {
	mmtlsClient *MmtlsClient
	httpClient  *http.Client
	curShortip  string
	mmtlsIsInit bool
}

// 长连接mmtls模块
type TcpClientModel struct {
	MmtlsClient *MmtlsClient

}

type MmtlsPacketHeader struct {
	headerbyte      byte
	headerversion   uint16
	headerpacketLen uint16
}

func Getsha256(Data []byte) []byte {
	D := sha256.New()
	D.Write(Data)
	return D.Sum(nil)
}

func Utf8ToBytes(Data string) []byte {
	t := make([]byte, utf8.RuneCountInString(Data))
	i := 0
	for _, r := range Data {
		t[i] = byte(r)
		i++
	}
	return t
}

func GetDecryptIv(decodeIv []byte, serverSeq int) ([]byte, int) {
	last := decodeIv[8:12]
	lastInt := binary.BigEndian.Uint32(last)
	xorInt := lastInt ^ uint32(serverSeq)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, xorInt)
	ret := new(bytes.Buffer)
	ret.Write(decodeIv[:8])
	ret.Write(buf.Bytes())
	serverSeq += 1
	return ret.Bytes(), serverSeq
}

func GetEncryptIv(encodeIv []byte, clientSeq int) ([]byte, int) {
	last := encodeIv[8:12]
	lastInt := binary.BigEndian.Uint32(last)
	xorInt := lastInt ^ uint32(clientSeq)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, xorInt)
	ret := new(bytes.Buffer)
	ret.Write(encodeIv[:8])
	ret.Write(buf.Bytes())
	clientSeq += 1
	return ret.Bytes(), clientSeq
}

func GetNonce(iv []byte, seq int) ([]byte, int) {
	last := iv[8:12]
	xorInt := binary.BigEndian.Uint32(last) ^ uint32(seq)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, xorInt)
	ret := new(bytes.Buffer)
	ret.Write(iv[:8])
	ret.Write(buf.Bytes())
	return ret.Bytes(), seq + 1
}

func UInt16ToBigEndianBytes(val uint16) []byte {
	var ret []byte
	binary.BigEndian.PutUint16(ret, val)
	return ret
}

func UInt32ToBigEndianBytes(val uint32) []byte {
	var ret []byte
	binary.BigEndian.PutUint32(ret, val)
	return ret
}

func UInt64ToBigEndianBytes(val uint64) []byte {
	var ret []byte
	binary.BigEndian.PutUint64(ret, val)
	return ret
}

