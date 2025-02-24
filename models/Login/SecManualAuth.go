package Login

import (
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/Mmtls"
	"wechatdll/comm"
	"wechatdll/lib"
)

func SecManualAuth(Data comm.LoginData, mmtlshost string) (mm.UnifyAuthResponse, []byte, []byte, []byte, mm.TrustResponse, error) {
	prikey, pubkey := Algorithm.GetEcdh713Key()

	httpclient := Mmtls.GenNewHttpClient(Data.MmtlsKey, mmtlshost)

	aeskey := []byte(lib.RandSeq(16)) //获取随机密钥
	accountRequest := &mm.ManualAuthRsaReqData{
		RandomEncryKey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(aeskey))),
			Buffer: aeskey,
		},
		CliPubEcdhkey: &mm.ECDHKey{
			Nid: proto.Int32(713),
			Key: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(pubkey))),
				Buffer: pubkey,
			},
		},
		UserName: &Data.Wxid,
		Pwd:      &Data.Pwd,
		Pwd2:     &Data.Pwd,
	}
	ccData := &mm.CryptoData{
		Version:     []byte("00000003"),
		Type:        proto.Uint32(1),
		EncryptData: Algorithm.GetiPadNewSpamData(Data.Deviceid_str, Data.DeviceName,Data.DeviceToken),
		Timestamp:   proto.Uint32(uint32(time.Now().Unix())),
		Unknown5:    proto.Uint32(5),
		Unknown6:    proto.Uint32(0),
	}

	Wcstf := Algorithm.IpadWcstf(Data.Wxid)
	Wcste := Algorithm.IpadWcste(0, 0)

	ccDataseq, _ := proto.Marshal(ccData)

	DeviceTokenCCD := &mm.DeviceToken{
		Version:   proto.String(""),
		Encrypted: proto.Uint32(1),
		Data: &mm.SKBuiltinStringT{
			String_: proto.String(Data.DeviceToken.GetTrustResponseData().GetDeviceToken()),
		},
		TimeStamp: proto.Uint32(uint32(time.Now().Unix())),
		Optype:    proto.Uint32(2),
		Uin:       proto.Uint32(0),
	}
	DeviceTokenCCDPB, _ := proto.Marshal(DeviceTokenCCD)

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
			ILen:                 proto.Uint32(uint32(len(DeviceTokenCCDPB))),
			Buffer:               DeviceTokenCCDPB,
		},
	}

	WCExtInfoseq, _ := proto.Marshal(WCExtInfo)
	ClientSeqId := lib.GetClientSeqId(Data.Deviceid_str)
	Imei := Algorithm.IOSImei(Data.Deviceid_str)
	// TODO: 放到初始化上下文中生成
	SoftType := Algorithm.SoftType_iPad(Data.Deviceid_str, Algorithm.IPadOsVersion, Algorithm.IPadModel)
	uuid1, _ := Algorithm.IOSUuid(Data.Deviceid_str)

	deviceRequest := &mm.ManualAuthAesReqData{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    aeskey,
			Uin:           proto.Uint32(0),
			DeviceId:      Data.Deviceid_byte,
			ClientVersion: proto.Int32(int32(Algorithm.IPadVersion)),
			DeviceType:    []byte(Algorithm.IPadDeviceType),
			Scene:         proto.Uint32(1),
		},
		BaseReqInfo:  &mm.BaseAuthReqInfo{},
		Imei:         &Imei,
		SoftType:     &SoftType,
		BuiltinIpseq: proto.Uint32(0),
		ClientSeqId:  &ClientSeqId,
		DeviceName:   proto.String(Data.DeviceName),
		DeviceType:   proto.String("iPad"),
		Language:     proto.String("zh_CN"),
		TimeZone:     proto.String("8.0"),
		Channel:      proto.Int(0),
		TimeStamp:    proto.Uint32(uint32(time.Now().Unix())),
		DeviceBrand:  proto.String("Apple"),
		Ostype:       proto.String(Algorithm.IPadDeviceType),
		RealCountry:  proto.String("CN"),
		BundleId:     proto.String("com.tencent.xin"),
		AdSource:     &uuid1,
		IphoneVer:    proto.String(Algorithm.IPadModel),
		InputType:    proto.Uint32(2),
		ExtSpamInfo: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(WCExtInfoseq))),
			Buffer: WCExtInfoseq,
		},
	}

	accountReqData, err := proto.Marshal(accountRequest)
	log.Println("account: " + hex.EncodeToString(accountReqData))
	deviceReqData, err := proto.Marshal(deviceRequest)
	log.Println("device: " + hex.EncodeToString(deviceReqData))

	requset := &mm.SecManualLoginRequest{
		RsaReqData: accountRequest,
		AesReqData: deviceRequest,
	}
	reqdata, err := proto.Marshal(requset)
	log.Println(hex.EncodeToString(reqdata))

	hec := &Algorithm.Client{}
	hec.Init("IOS")
	hecData := hec.HybridEcdhPackIosEn(252, 0, nil, reqdata)

	recvData, err := httpclient.MMtlsPost(mmtlshost, "/cgi-bin/micromsg-bin/secmanualauth", hecData, Data.Proxy)
	ph1 := hec.HybridEcdhPackIosUn(recvData)
	loginRes := mm.UnifyAuthResponse{}
	err = proto.Unmarshal(ph1.Data, &loginRes)

	if err != nil {
		return mm.UnifyAuthResponse{}, nil, nil, nil, mm.TrustResponse{}, err
	}

	return loginRes, prikey, pubkey, ph1.Cookies, mm.TrustResponse{}, nil
}
