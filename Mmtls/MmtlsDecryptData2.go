package Mmtls

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"wechatdll/Algorithm"
	"wechatdll/lib"
)

func (httpclient *HttpClientModel) MmtlsDecryptData2(Data []Separatea) []byte {
	responsr := new(bytes.Buffer)
	for i, v := range Data {
		if i == 0 {
			if v.title == "16" {
				newsendbufferhashs := new(bytes.Buffer)
				newsendbufferhashs.Write(httpclient.mmtlsClient.Newsendbufferhashs)
				newsendbufferhashs.Write(v.val)
				decrypt_serverdata_hash256 := Getsha256(newsendbufferhashs.Bytes())

				var HkdfExpand_handshake = new(bytes.Buffer)
				HkdfExpand_handshake.Write([]byte{0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x20, 0x6b, 0x65, 0x79, 0x20, 0x65, 0x78, 0x70, 0x61, 0x6e, 0x73, 0x69, 0x6f, 0x6e})
				HkdfExpand_handshake.Write(decrypt_serverdata_hash256)

				HkdfExpand_handshake_key := Algorithm.Hkdf_Expand(sha256.New, httpclient.mmtlsClient.Hkdfexpand_pskaccess_key, HkdfExpand_handshake.Bytes(), 28)
				httpclient.mmtlsClient.Decrptshortmmtlskey = HkdfExpand_handshake_key[:16]
				httpclient.mmtlsClient.Decrptshortmmtlsiv = HkdfExpand_handshake_key[16:28]
			}
		} else {
			if v.title == "17" {
				decryptserverfinishdata := NewAESGCMDecrypter(v.val, httpclient.mmtlsClient.Decrptshortmmtlskey, httpclient.mmtlsClient.Decrptshortmmtlsiv, int32(i), v.title)
				if decryptserverfinishdata != nil {
					responsr.Write(decryptserverfinishdata)
				}
			}
		}
	}
	return responsr.Bytes()
}

func NewAESGCMDecrypter(Data, key, iv []byte, seq int32, mmtlsType string) []byte {
	businessdata_aad := new(bytes.Buffer)
	businessdata_aad.Write([]byte{0x00, 0x00, 0x00, 0x00})
	businessdata_aad.Write(lib.IntToBytes(seq))
	mmtlsTypebyte, _ := hex.DecodeString(mmtlsType)
	businessdata_aad.Write(mmtlsTypebyte)
	businessdata_aad.Write([]byte{0xf1, 0x03})
	businessdata_aad.Write(lib.Int16ToBytes(int16(len(Data))))

	var xorkeyBuffer = bytes.NewReader(iv[8:12])
	var xorkeyint uint32
	binary.Read(xorkeyBuffer, binary.BigEndian, &xorkeyint)
	xorkeyint = xorkeyint ^ uint32(seq)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, xorkeyint)

	var decryptmmtlsIv_seq = new(bytes.Buffer)
	decryptmmtlsIv_seq.Write(iv[:8])
	decryptmmtlsIv_seq.Write(buf.Bytes())

	return Algorithm.NewAES_GCMDecrypter(key, Data, decryptmmtlsIv_seq.Bytes(), businessdata_aad.Bytes())
}
