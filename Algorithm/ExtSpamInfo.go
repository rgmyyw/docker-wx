package Algorithm

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"time"
	"wechatdll/Cilent/mm"
)

func Wcstf(Username string, T int64) []byte {
	curtime := uint64(T / 1e6)
	contentlen := len(Username)

	var ct []uint64
	ut := curtime
	for i := 0; i < contentlen; i++ {
		ut += uint64(rand.Intn(10000))
		ct = append(ct, ut)
	}
	ccd := &mm.Wcstf{
		StartTime: &curtime,
		CheckTime: &curtime,
		Count:     proto.Uint32(uint32(contentlen)),
		EndTime:   ct,
	}

	pb, _ := proto.Marshal(ccd)

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, _ = w.Write(pb)
	_ = w.Close()

	zt := new(ZT)
	zt.Init()
	encData := SaeEncrypt07(b.Bytes())

	Ztdata := &mm.ZTData{
		Version:   proto.String("00000007"),
		Encrypted: proto.Uint32(1),
		Data:      encData,
		TimeStamp: proto.Uint32(uint32(T)),
		Optype:    proto.Uint32(5),
		Uin:       proto.Uint32(0),
	}
	MS, _ := proto.Marshal(Ztdata)
	return MS
}

func Wcste(A, B uint64, T int64) []byte {

	curtime := uint32(T)
	curNanoTime := uint64(time.Now().UnixNano())

	ccd := &mm.Wcste{
		Checkid:   proto.String("<LoginByID>"),
		StartTime: &curtime,
		CheckTime: &curtime,
		Count1:    proto.Uint32(0),
		Count2:    proto.Uint32(1),
		Count3:    proto.Uint32(0),
		Const1:    proto.Uint64(A),
		Const2:    &curNanoTime,
		Const3:    &curNanoTime,
		Const4:    &curNanoTime,
		Const5:    &curNanoTime,
		Const6:    proto.Uint64(B),
	}

	pb, _ := proto.Marshal(ccd)

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, _ = w.Write(pb)
	_ = w.Close()

	zt := new(ZT)
	zt.Init()
	encData := SaeEncrypt07(b.Bytes())

	Ztdata := &mm.ZTData{
		Version:   proto.String("00000007"),
		Encrypted: proto.Uint32(1),
		Data:      encData,
		TimeStamp: proto.Uint32(uint32(T)),
		Optype:    proto.Uint32(5),
		Uin:       proto.Uint32(0),
	}

	MS, _ := proto.Marshal(Ztdata)
	return MS
}
func GetIOSExtSpamInfoAndDeviceToken(Wxid, Deviceid_str, DeviceName string, DeviceToken mm.TrustResponse, T int64) []byte {
	ccData := &mm.CryptoData{
		Version:     []byte("00000007"),
		Type:        proto.Uint32(1),
		EncryptData: GetiPadNewSpamData(Deviceid_str, DeviceName, DeviceToken),
		Timestamp:   proto.Uint32(uint32(T)),
		Unknown5:    proto.Uint32(5),
		Unknown6:    proto.Uint32(0),
	}
	ccDataseq, _ := proto.Marshal(ccData)
	fmt.Printf("mm:%+v\n", DeviceToken)
	fmt.Println(DeviceToken.GetTrustResponseData().GetDeviceToken())
	DeviceTokenCCD := &mm.DeviceToken{
		Version:   proto.String(""),
		Encrypted: proto.Uint32(1),
		Data: &mm.SKBuiltinStringT{
			String_: proto.String(DeviceToken.GetTrustResponseData().GetDeviceToken()),
		},
		TimeStamp: proto.Uint32(uint32(T)),
		Optype:    proto.Uint32(2),
		Uin:       proto.Uint32(0),
	}
	DeviceTokenCCDPB, _ := proto.Marshal(DeviceTokenCCD)

	Wcstf := Wcstf(Wxid, T)
	Wcste := Wcste(0, 0, T)

	WCExtInfo := &mm.WCExtInfo{
		Wcstf: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(Wcstf))),
			Buffer: Wcstf,
		},
		Wcste: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(Wcste))),
			Buffer: Wcste,
		},
		CcData: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(ccDataseq))),
			Buffer: ccDataseq,
		},
		DeviceToken: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(DeviceTokenCCDPB))),
			Buffer: DeviceTokenCCDPB,
		},
	}

	WCExtInfoseq, _ := proto.Marshal(WCExtInfo)
	fmt.Println(WCExtInfoseq)

	return WCExtInfoseq
}
