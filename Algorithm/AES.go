package Algorithm

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

func padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

func AESEncrypt(src []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	src = padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return src
}

func AESDecrypt(src []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return src
}

func AesEncrypt(RequestSerialize []byte, key []byte) []byte {
	//根据key 生成密文
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}

	blockSize := block.BlockSize()
	RequestSerialize = PKCS5Padding(RequestSerialize, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(RequestSerialize))
	blockMode.CryptBlocks(crypted, RequestSerialize)

	return crypted
}

func AesDecrypt(body []byte, key []byte) []byte {
	// fmt.Printf("%x\n\n", body)
	// fmt.Printf("%x\n\n", key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(body))
	blockMode.CryptBlocks(origData, body)
	origData = PKCS5UnPadding(origData)
	// fmt.Printf("%x\n", origData)
	return origData
}

// =================== ECB ======================
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func CompressAndAes(RequestSerialize []byte, aeskey []byte) []byte {
	compressed := DoZlibCompress(RequestSerialize)
	return AesEncrypt(compressed, aeskey)
}

func DecompressAndAesDecrypt(body []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}

	if len(body)%aes.BlockSize != 0 {
		log.Error(fmt.Sprintf("crypto/cipher: data is not a multiple of the block size，[BodyLength：%v] [AesLength：%v]", len(body), aes.BlockSize))
		return nil
	}

	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(body))
	blockMode.CryptBlocks(origData, body)
	origData = PKCS5UnPadding(origData)
	origData = DoZlibUnCompress(origData)
	return origData
}

func AesGcmDecrypt(key, nonce, input, additional []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, input, additional)
}

func AesGcmEncrypt(key, nonce, input, additional []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	result := aesgcm.Seal(nil, nonce, input, additional)

	return result, nil
}

func AesGcmDecryptWithUnCompress(key, input, additional []byte) []byte {

	inputSize := len(input)

	nonce := make([]byte, 12)
	copy(nonce, input[inputSize-28:inputSize-16])

	tag := make([]byte, 16)
	copy(tag, input[inputSize-16:])

	cipherText := make([]byte, inputSize-28)
	copy(cipherText, input[:inputSize-28])
	cipherText = append(cipherText, tag...)

	result, _ := AesGcmDecrypt(key, nonce, cipherText, additional)

	b := bytes.NewReader(result)

	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	r.Close()

	return out.Bytes()
}

func AesGcmEncryptWithCompress(key, nonce, input, additional []byte) []byte {

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(input)
	w.Close()

	data, _ := AesGcmEncrypt(key, nonce, b.Bytes(), additional)

	encData := data[:len(data)-16]
	tag := data[len(data)-16:]

	totalData := []byte{}
	totalData = append(totalData, encData...)
	totalData = append(totalData, nonce...)
	totalData = append(totalData, tag...)
	return totalData
}

func AesGcmEncryptWithCompressZlib(key []byte, plaintext []byte, nonce []byte, aad []byte) []byte {
	compressData := DoZlibCompress(plaintext)
	//nonce := []byte(randSeq(12)) //获取随机密钥
	encrypt_data := NewAES_GCMEncrypter(key, compressData, nonce, aad)
	outdata := encrypt_data[:len(encrypt_data)-16]
	retdata := new(bytes.Buffer)
	retdata.Write(outdata)
	retdata.Write(nonce)
	retdata.Write(encrypt_data[len(encrypt_data)-16:])
	return retdata.Bytes()
}

func AesGcmDecryptWithcompressZlib(key []byte, ciphertext []byte, aad []byte) []byte {
	ciphertextinput := ciphertext[:len(ciphertext)-0x1c]
	endatanonce := ciphertext[len(ciphertext)-0x1c : len(ciphertext)-0x10]
	data := new(bytes.Buffer)
	data.Write(ciphertextinput)
	data.Write(ciphertext[len(ciphertext)-0x10 : len(ciphertext)])
	decrypt_data := NewAES_GCMDecrypter(key, data.Bytes(), endatanonce, aad)
	if len(decrypt_data) > 0 {
		return DoZlibUnCompress(decrypt_data)
	} else {
		return []byte{}
	}

}
