package Algorithm

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/hkdf"
	"io"
	"wechatdll/Cilent/mm"
)

func (h *Client) encryptoIOS(Data []byte) []byte {
	ecdhkey := DoECDH415Key(h.Privkey, h.InitPubKey)
	m := sha256.New()
	m.Write(ecdhkey)
	ecdhkey = m.Sum(nil)
	mClientpubhash := sha256.New()
	mClientpubhash.Write([]byte("1"))
	mClientpubhash.Write([]byte("415"))
	mClientpubhash.Write(h.PubKey)
	mClientpubhash_digest := mClientpubhash.Sum(nil)

	mRandomEncryptKey := make([]byte, 32)
	io.ReadFull(rand.Reader, mRandomEncryptKey)
	mNonce := make([]byte, 12)
	io.ReadFull(rand.Reader, mNonce)

	mEncryptdata := AesGcmEncryptWithCompressZlib(ecdhkey[:24], mRandomEncryptKey, mNonce, mClientpubhash_digest)
	var mExternEncryptdata []byte
	if len(h.Externkey) == 0x20 {
		mExternEncryptdata = AesGcmEncryptWithCompressZlib(h.Externkey[:24], mRandomEncryptKey, mNonce, mClientpubhash_digest)
	}
	hkdfexpand_security_key := HybridHkdfExpand([]byte("security hdkf expand"), mRandomEncryptKey, mClientpubhash_digest, 56)

	mClientpubhashFinal := sha256.New()
	mClientpubhashFinal.Write([]byte("1"))
	mClientpubhashFinal.Write([]byte("415"))
	mClientpubhashFinal.Write(h.PubKey)
	mClientpubhashFinal.Write(mEncryptdata)
	mClientpubhashFinal.Write(mExternEncryptdata)
	mClientpubhashFinal_digest := mClientpubhashFinal.Sum(nil)

	mEncryptdataFinal := AesGcmEncryptWithCompressZlib(hkdfexpand_security_key[:24], Data, mNonce, mClientpubhashFinal_digest)

	h.clientHash.Write(mEncryptdataFinal)

	h.serverHash.Write(hkdfexpand_security_key[24:56])
	h.serverHash.Write(Data)

	HybridEcdhRequest := &mm.HybridEcdhRequest{
		Type: proto.Int32(1),
		SecECDHKey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(415),
			Buffer: h.PubKey,
		},
		Randomkeydata:       mEncryptdata,
		Randomkeyextenddata: mExternEncryptdata,
		Encyptdata:          mEncryptdataFinal,
	}
	reqdata, _ := proto.Marshal(HybridEcdhRequest)
	return reqdata
}

func (h *Client) decryptoIOS(Data []byte) []byte {
	HybridEcdhResponse := &mm.HybridEcdhResponse{}
	err := proto.Unmarshal(Data, HybridEcdhResponse)
	if err != nil {
		return nil
	}
	decrptecdhkey := DoECDH415Key(h.Privkey, HybridEcdhResponse.GetSecECDHKey().GetBuffer())
	m := sha256.New()
	m.Write(decrptecdhkey)
	decrptecdhkey = m.Sum(nil)
	h.serverHash.Write([]byte("415"))
	h.serverHash.Write(HybridEcdhResponse.GetSecECDHKey().GetBuffer())
	h.serverHash.Write([]byte("1"))
	mServerpubhashFinal_digest := h.serverHash.Sum(nil)

	outdata := AesGcmDecryptWithcompressZlib(decrptecdhkey[:24], HybridEcdhResponse.GetDecryptdata(), mServerpubhashFinal_digest)
	return outdata
}

func (h *Client) encryptAndroid(input []byte) []byte {
	if h.Status != HYBRID_ENC {
		return nil
	}

	priv, x, y, error := elliptic.GenerateKey(h.curve, rand.Reader)
	if error != nil {
		return nil
	}
	h.Privkey = priv
	h.PubKey = elliptic.Marshal(h.curve, x, y)

	ecdhKey := Ecdh(h.curve, h.InitPubKey, h.Privkey)

	//hash1
	h1 := sha256.New()
	h1.Write([]byte("1"))
	h1.Write([]byte("415"))
	h1.Write(h.PubKey)
	h1Sum := h1.Sum(nil)

	//Random
	random := make([]byte, 32)
	io.ReadFull(rand.Reader, random)

	nonce1 := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce1)
	gcm1 := AesGcmEncryptWithCompress(ecdhKey[0:0x18], nonce1, random, h1Sum)
	//hkdf
	salt, _ := hex.DecodeString("73656375726974792068646B6620657870616E64")
	hkdfKey := make([]byte, 56)
	hkdf.New(sha256.New, random, salt, h1Sum).Read(hkdfKey)

	//hash2
	h2 := sha256.New()
	h2.Write([]byte("1"))
	h2.Write([]byte("415"))
	h2.Write(h.PubKey)
	h2.Write(gcm1)
	h2Sum := h2.Sum(nil)

	nonce2 := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce2)
	gcm2 := AesGcmEncryptWithCompress(hkdfKey[0:0x18], nonce2, input, h2Sum)

	var nid int32 = 415
	secKey := &mm.SecKey{
		Nid: &nid,
		Key: h.PubKey,
	}

	var ver int32 = 1
	he := &mm.HybridEcdhReq{
		Version: &ver,
		SecKey:  secKey,
		Gcm1:    gcm1,
		Autokey: []byte{},
		Gcm2:    gcm2,
	}

	protoMsg, _ := proto.Marshal(he)

	// update client
	h.clientHash.Write(hkdfKey[0x18:0x38])
	h.clientHash.Write(input)

	// update server
	h.serverHash.Write(gcm2)

	h.Status = HYBRID_DEC

	return protoMsg
}

func (h *Client) decryptAndroid(input []byte) []byte {

	if h.Status != HYBRID_DEC {
		return nil
	}

	var resp mm.HybridEcdhResp
	proto.Unmarshal(input, &resp)

	h.serverHash.Write(resp.GetGcm1())
	//	hServ := h.serverHash.Sum(nil)
	//	fmt.Printf("%x\n", hServ)

	//	ecdsa.Verify(h.clientEcdsaPub, resp.GetGcm2(), hServ)
	ecKey := Ecdh(h.curve, resp.GetSecKey().GetKey(), h.Privkey)

	h.clientHash.Write([]byte("415"))
	h.clientHash.Write(resp.GetSecKey().GetKey())
	h.clientHash.Write([]byte("1"))
	hCli := h.clientHash.Sum(nil)

	plain := AesGcmDecryptWithUnCompress(ecKey[0:0x18], resp.GetGcm1(), hCli)
	return plain
}
