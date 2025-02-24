package Algorithm

import (
	"bytes"
	"encoding/binary"
	"math/rand"
)

var Max32bitValue uint64 = 4294967295

type QQCryptor struct {
	contextStart 	int
	crypt           int
	header         	bool
	key         	[]byte
	out             []byte
	padding         int
	plain           []byte
	pos				int
	preCrypt		int
	prePlain		[]byte
}

func GetQQCryptor() *QQCryptor {
	retVal := &QQCryptor{
		contextStart: 0,
		crypt:        0,
		header:       true,
		key:          nil,
		out:          nil,
		padding:      0,
		plain:        nil,
		pos:          0,
		preCrypt:     0,
		prePlain:     nil,
	}
	return retVal
}

func getUnsignedInt(bArr []byte, i2 int, i3 int) uint64 {
	var i4 int
	var j uint64 = 0
	if i3 > 8 {
		i4 = i2 + 8
	} else {
		i4 = i2 + i3
	}
	for {
		if i2 >= i4 {
			break
		}
		j = (j << 8) | ((uint64) (bArr[i2] & 255))
		i2++
	}
	return (Max32bitValue & j) | (j >> 32)
}

func arrayCopy(source []byte, sourceStart int, dest []byte, destStart int, length int) {
	for i := 0; i < length; i++ {
		dest[destStart + i] = source[sourceStart + i]
	}
}

func (this *QQCryptor) decipher(bArr []byte, i2 int) []byte  {
	var i3 int = 16
	unsignedInt := getUnsignedInt(bArr, i2, 4)
	unsignedInt2 := getUnsignedInt(bArr, i2 + 4, 4)
	unsignedInt3 := getUnsignedInt(this.key, 0, 4)
	unsignedInt4 := getUnsignedInt(this.key, 4, 4)
	unsignedInt5 := getUnsignedInt(this.key, 8, 4)
	unsignedInt6 := getUnsignedInt(this.key, 12, 4)
	var j uint64 = uint64(0xE3779B90) & Max32bitValue		// -478700656
	var j2 uint64 = uint64(0x9E3779B9) & Max32bitValue		// -1640531527
	for {
		var i4 int = i3 - 1
		if i3 <= 0 {
			byteBuffer := new(bytes.Buffer)
			binary.Write(byteBuffer, binary.BigEndian, uint32(unsignedInt))
			binary.Write(byteBuffer, binary.BigEndian, uint32(unsignedInt2))
			return byteBuffer.Bytes()
		}
		unsignedInt2 = (unsignedInt2 - ((((unsignedInt << 4) + unsignedInt5) ^ (unsignedInt + j)) ^ ((unsignedInt >> 5) + unsignedInt6))) & Max32bitValue
		unsignedInt = (unsignedInt - ((((unsignedInt2 << 4) + unsignedInt3) ^ (unsignedInt2 + j)) ^ ((unsignedInt2 >> 5) + unsignedInt4))) & Max32bitValue
		j = (j - j2) & Max32bitValue
		i3 = i4
	}
}

func (this *QQCryptor) encipher(bArr []byte) []byte {
	i2 := 16
	unsignedInt := getUnsignedInt(bArr, 0, 4)
	unsignedInt2 := getUnsignedInt(bArr, 4, 4)
	unsignedInt3 := getUnsignedInt(this.key, 0, 4)
	unsignedInt4 := getUnsignedInt(this.key, 4, 4)
	unsignedInt5 := getUnsignedInt(this.key, 8, 4)
	unsignedInt6 := getUnsignedInt(this.key, 12, 4)
	var j uint64 = 0
	var j2 uint64 = uint64(0x9E3779B9) & Max32bitValue
	for {
		i3 := i2 - 1
		if i2 <= 0 {
			byteBuffer := new(bytes.Buffer)
			binary.Write(byteBuffer, binary.BigEndian, uint32(unsignedInt))
			binary.Write(byteBuffer, binary.BigEndian, uint32(unsignedInt2))
			retBytes := byteBuffer.Bytes()
			return retBytes
		}
		j = (j + j2) & Max32bitValue
		unsignedInt = (unsignedInt + ((((unsignedInt2 << 4) + unsignedInt3) ^ (unsignedInt2 + j)) ^ ((unsignedInt2 >> 5) + unsignedInt4))) & Max32bitValue
		unsignedInt2 = (unsignedInt2 + ((((unsignedInt << 4) + unsignedInt5) ^ (unsignedInt + j)) ^ ((unsignedInt >> 5) + unsignedInt6))) & Max32bitValue
		i2 = i3
	}
}

func (this *QQCryptor) decrypt8Bytes(bArr []byte, i2 int, i3 int) bool {
	this.pos = 0
	for{
		if this.pos >= 8 {
			break
		}
		if this.contextStart + this.pos >= i3 {
			return true
		}
		bArr2 := this.prePlain
		i4 := this.pos
		bArr2[i4] = bArr2[i4] ^ bArr[(this.crypt + i2) + this.pos]
		this.pos++
 	}
 	this.prePlain = this.decipher(this.prePlain, 0)
 	if this.prePlain == nil {
 		return false
	}
	this.contextStart += 8
	this.crypt += 8
	this.pos = 0
	return true
}

func (this *QQCryptor) encrypt8Bytes() {
	this.pos = 0
	for {
		if this.pos >= 8 {
			break
		}
		if this.header {
			bArr := this.plain
			i2 := this.pos
			bArr[i2] = bArr[i2] ^ this.prePlain[this.pos]
		} else {
			bArr2 := this.plain
			i3 := this.pos
			bArr2[i3] = bArr2[i3] ^ this.out[this.preCrypt + this.pos]
		}
		this.pos++
	}
	arrayCopy(this.encipher(this.plain), 0, this.out, this.crypt, 8)
	this.pos = 0
	for {
		if this.pos >= 8 {
			break
		}
		bArr3 := this.out
		i4 := this.crypt + this.pos
		bArr3[i4] = bArr3[i4] ^ this.prePlain[this.pos]
		this.pos++
	}
	arrayCopy(this.plain, 0, this.prePlain, 0, 8)
	this.preCrypt = this.crypt
	this.crypt += 8
	this.pos = 0
	this.header = false
}

func (this *QQCryptor) Decrypt(bArr []byte, bArr2 []byte) []byte {
	var i4 int
	i2 := 0
	i3 := len(bArr)
	this.preCrypt = 0
	this.crypt = 0
	this.key = bArr2
	bArr3 := make([]byte, i2 + 8)
	if i3 % 8 != 0 || i3 < 16 {
		return nil
	}
	this.prePlain = this.decipher(bArr, i2)
	this.pos = (int)(this.prePlain[0]) & 7
	var i5 int = (i3 - this.pos) - 10
	if i5 < 0 {
		return nil
	}
	for i6 := i2; i6 < len(bArr3); i6++ {
		bArr3[i6] = 0
	}
	this.out = make([]byte, i5)
	this.preCrypt = 0
	this.crypt = 8
	this.contextStart = 8
	this.pos++
	this.padding = 1
	bArr4 := bArr3
	for {
		if this.padding > 2 {
			break
		}
		if this.pos < 8 {
			this.pos++
			this.padding++
		}
		if this.pos == 8 {
			if !this.decrypt8Bytes(bArr, i2, i3) {
				return nil
			}
			bArr4 = bArr
		}
	}
	i7 := 0
	bArr5 := bArr4
	for {
		if i5 == 0 {
			break
		}
		if this.pos < 8 {
			this.out[i7] = (byte) (bArr5[(this.preCrypt + i2) + this.pos] ^ this.prePlain[this.pos])
			i4 = i7 + 1
			this.pos++
			i5--
		} else {
			i4 = i7
		}
		if this.pos == 8 {
			this.preCrypt = this.crypt - 8
			if !this.decrypt8Bytes(bArr, i2, i3) {
				return nil
			}
			i7 = i4
			bArr5 = bArr
		} else {
			i7 = i4
		}
	}
	this.padding = 1
	bArr6 := bArr5
	for {
		if this.padding >= 8 {
			break
		}
		if this.pos < 8 {
			if (bArr6[(this.preCrypt + i2) + this.pos] ^ this.prePlain[this.pos]) != 0 {
				return nil
			}
			this.pos++
		}
		if this.pos == 8 {
			this.preCrypt = this.crypt
			if !this.decrypt8Bytes(bArr, i2, i3) {
				return nil
			}
			bArr6 = bArr
		}
		this.padding++
	}
	return this.out
}

func (this *QQCryptor) Encrypt(bArr []byte, bArr2 []byte) []byte {
	var i4 int
	i2 := 0
	i3 := len(bArr)
	this.plain = make([]byte, 8)
	this.prePlain = make([]byte, 8)
	this.pos = 1
	this.padding = 0
	this.preCrypt = 0
	this.crypt = 0
	this.key = bArr2
	this.header = true
	this.pos = (i3 + 10) % 8
	if this.pos != 0 {
		this.pos = 8 - this.pos
	}
	this.out = make([]byte, this.pos + i3 + 10)
	this.plain[0] = (byte) ((rand.Int() & 248) | this.pos)
	for i5 := 1; i5 <= this.pos; i5++ {
		this.plain[i5] = (byte) (rand.Int() & 255)
	}

	this.pos++
	for i6 := 0; i6 < 8; i6++ {
		this.prePlain[i6] = 0
	}
	this.padding = 1
	for {
		if this.padding > 2 {
			break
		}
		if this.pos < 8 {
			bArr3 := this.plain
			i7 := this.pos
			this.pos = i7 + 1
			bArr3[i7] = (byte) (rand.Int() & 255)
			this.padding++
		}
		if this.pos == 8 {
			this.encrypt8Bytes()
			//fmt.Printf("go单一分段结果: %x\n", this.out)
		}
	}
	i8 := i2
	i9 := i3
	for {
		if i9 <= 0 {
			break
		}
		if this.pos < 8 {
			bArr4 := this.plain
			i10 := this.pos
			this.pos = i10 + 1
			i4 = i8 + 1
			bArr4[i10] = bArr[i8]
			i9--
		} else {
			i4 = i8
		}
		if this.pos == 8 {
			this.encrypt8Bytes()
			//fmt.Printf("go单一分段结果: %x\n", this.out)
		}
		i8 = i4
	}
	this.padding = 1
	for {
		if this.padding > 7 {
			break
		}
		if (this.pos < 8) {
			bArr5 := this.plain
			i11 := this.pos
			this.pos = i11 + 1
			bArr5[i11] = 0
			this.padding++
		}
		if this.pos == 8 {
			this.encrypt8Bytes()
			//fmt.Printf("go单一分段结果: %x\n", this.out)
		}
	}
	return this.out
}
