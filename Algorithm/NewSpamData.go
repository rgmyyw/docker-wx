package Algorithm

import (
	"github.com/golang/protobuf/proto"
	"hash/crc32"
	"time"
	"wechatdll/Cilent/mm"
)

// ipad生成ccd, 使用03加密
func GetiPadNewSpamData(Deviceid, DeviceName string, DeviceToken mm.TrustResponse) []byte {
	T := uint32(time.Now().Unix())
	timeStamp := int(T)
	xorKey := uint8((timeStamp * 0xffffffed) + 7)

	uuid1, uuid2 := IOSUuid(Deviceid)

	if len(Deviceid) < 32 {
		Dlen := 32 - len(Deviceid)
		Fill := "ff95DODUJ4EysYiogKZSmajWCUKUg9RX"
		Deviceid = Deviceid + Fill[:Dlen]
	}

	spamDataBody := &mm.SpamDataBody{
		UnKnown1:              proto.Int32(1),
		TimeStamp:             proto.Int32(int32(timeStamp)),
		KeyHash:               proto.Int32(int32(MakeKeyHash(int(xorKey)))),
		Yes1:                  proto.String(XorEncodeStr("yes", xorKey)),
		Yes2:                  proto.String(XorEncodeStr("yes", xorKey)),
		IosVersion:            proto.String(XorEncodeStr("13.5", xorKey)),
		DeviceType:            proto.String(XorEncodeStr("iPad", xorKey)),
		UnKnown2:              proto.Int32(6),
		IdentifierForVendor:   proto.String(XorEncodeStr(uuid1, xorKey)),
		AdvertisingIdentifier: proto.String(XorEncodeStr(uuid2, xorKey)),
		Carrier:               proto.String(XorEncodeStr("中国联通", xorKey)),
		BatteryInfo:           proto.Int32(1),
		NetworkName:           proto.String(XorEncodeStr("en0", xorKey)),
		NetType:               proto.Int32(0),
		AppBundleId:           proto.String(XorEncodeStr("com.tencent.xin", xorKey)),
		DeviceName:            proto.String(XorEncodeStr(DeviceName, xorKey)),
		UserName:              proto.String(XorEncodeStr("iPad11,3", xorKey)),
		Unknown3:              proto.Int64(77968568554095637),
		Unknown4:              proto.Int64(77968568554095617),
		Unknown5:              proto.Int32(5),
		Unknown6:              proto.Int32(4),
		Lang:                  proto.String(XorEncodeStr("zh", xorKey)),
		Country:               proto.String(XorEncodeStr("CN", xorKey)),
		Unknown7:              proto.Int32(4),
		DocumentDir:           proto.String(XorEncodeStr("/var/mobile/Containers/Data/Application/94E41585-A27E-4933-AF06-5ABF7C774A6F/Documents", xorKey)),
		Unknown8:              proto.Int32(0),
		Unknown9:              proto.Int32(1),
		HeadMD5:               proto.String(XorEncodeStr("901bf05e51e2cb5585760f7e0116d0ba", xorKey)),
		AppUUID:               proto.String(XorEncodeStr(uuid1, xorKey)),
		SyslogUUID:            proto.String(""),
		Unknown10:             proto.String(""),
		Unknown11:             proto.String(""),
		AppName:               proto.String(XorEncodeStr("微信", xorKey)),
		SshPath:               proto.String(""),
		TempTest:              proto.String(""),
		DevMD5:                proto.String(""),
		DevUser:               proto.String(""),
		Unknown12:             proto.String(""),
		IsModify:              proto.Int32(0),
		ModifyMD5:             proto.String(""),
		RqtHash:               proto.Int64(288530629010980929),
		Unknown43:             proto.Uint64(1586355322),
		Unknown44:             proto.Uint64(1586355519000),
		Unknown45:             proto.Uint64(0),
		Unknown46:             proto.Uint64(288530629010980929),
		Unknown47:             proto.Uint64(1),
		Unknown48:             proto.String(XorEncodeStr(Deviceid, xorKey)),
		Unknown49:             proto.String(""),
		Unknown50:             proto.String(""),
		Unknown51:             proto.String(XorEncodeStr(DeviceToken.GetTrustResponseData().GetSoftData().GetSoftConfig(), xorKey)),
		Unknown52:             proto.Uint64(0),
		Unknown53:             proto.String(""),
		Unknown54:             proto.String(XorEncodeStr(DeviceToken.GetTrustResponseData().GetDeviceToken(), xorKey)),
	}
	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/WeChat", xorKey)),
		Filepath: proto.String(XorEncodeStr("7195B97E-9078-3119-9110-8BDA959283F0", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/Library/MobileSubstrate/MobileSubstrate.dylib", xorKey)),
		Filepath: proto.String(XorEncodeStr("3134CFB2-F722-310E-A2C7-42AE4DC131AB", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/mars.framework/mars", xorKey)),
		Filepath: proto.String(XorEncodeStr("A87DAD8E-E356-3E1E-9925-D63EA1614A95", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/andromeda.framework/andromeda", xorKey)),
		Filepath: proto.String(XorEncodeStr("EB5B920E-3AE6-3534-9DA4-C32DF72E33BD", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/OpenSSL.framework/OpenSSL", xorKey)),
		Filepath: proto.String(XorEncodeStr("8FAE149B-602B-3B9D-A620-88EA75CE153F", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/ProtobufLite.framework/ProtobufLite", xorKey)),
		Filepath: proto.String(XorEncodeStr("6F0D3077-4301-3D8F-8579-E34902547580", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/marsbridgenetwork.framework/marsbridgenetwork", xorKey)),
		Filepath: proto.String(XorEncodeStr("CFED9A03-C881-3D50-B014-732D0A09879F", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/matrixreport.framework/matrixreport", xorKey)),
		Filepath: proto.String(XorEncodeStr("1E7F06D2-DD36-31A8-AF3B-00D62054E1F9", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/libswiftCore.dylib", xorKey)),
		Filepath: proto.String(XorEncodeStr("AD0CAD3B-1B51-3327-8644-8BE1FF1F0AE9", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/libswiftDispatch.dylib", xorKey)),
		Filepath: proto.String(XorEncodeStr("9FCBA8ED-D8FD-3C16-9740-5E2A31F3E959", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/libswiftFoundation.dylib", xorKey)),
		Filepath: proto.String(XorEncodeStr("9702769F-1F06-3001-AB75-5AD38E1F7D66", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/libswiftObjectiveC.dylib", xorKey)),
		Filepath: proto.String(XorEncodeStr("1180AC10-0A92-39DB-8497-2B6D4217B8EB", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/libswiftDarwin.dylib", xorKey)),
		Filepath: proto.String(XorEncodeStr("999C2967-8A06-3CD5-82D7-D156E9440A0C", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/libswiftCoreGraphics.dylib", xorKey)),
		Filepath: proto.String(XorEncodeStr("DC548EF9-00F9-3A15-B5DB-05E39D9B5C5B", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/private/var/containers/Bundle/Application/2F493AE2-C0EB-4B4E-A86C-CE9BA3C0FA14/WeChat.app/Frameworks/libswiftCoreFoundation.dylib", xorKey)),
		Filepath: proto.String(XorEncodeStr("25114AE1-4AE9-3DBC-B3DE-7F9F9A5B45D2", xorKey)),
	})

	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/Library/Frameworks/CydiaSubstrate.framework/Libraries/SubstrateLoader.dylib", xorKey)),
		Filepath: proto.String(XorEncodeStr("54645DC0-3212-31D8-8A02-2FD67A793278", xorKey)),
	})

	srcdata, _ := proto.Marshal(spamDataBody)
	
	newClientCheckData := &mm.NewClientCheckData{
		C32Cdata:  proto.Int64(int64(crc32.ChecksumIEEE([]byte(srcdata)))),
		TimeStamp: proto.Int64(time.Now().Unix()),
		Databody:  srcdata,
	}

	ccddata, _ := proto.Marshal(newClientCheckData)
	//压缩数据
	compressdata := DoZlibCompress(ccddata)
	//compressdata := AE(ccddata)

	var encData []byte
	// 03加密
	zt := new(ZT)
	zt.Init()
	encData = zt.Encrypt(compressdata)

	return encData
}
// iphone生成ccd, 使用06加密
func GetiPhoneNewSpamData(Deviceid, DeviceName string, DeviceToken mm.TrustResponse) []byte {
	timeStamp := int(time.Now().Unix())
	xorKey := uint8((timeStamp * 0xffffffed) + 7)

	uuid1, uuid2 := IOSUuid(Deviceid)

	if len(Deviceid) < 32 {
		Dlen := 32 - len(Deviceid)
		Fill := "ff95DODUJ4EysYiogKZSmajWCUKUg9RX"
		Deviceid = Deviceid + Fill[:Dlen]
	}

	spamDataBody := &mm.SpamDataBody{
		UnKnown1:              proto.Int32(1),
		TimeStamp:             proto.Int32(int32(timeStamp)),
		KeyHash:               proto.Int32(int32(MakeKeyHash(int(xorKey)))),
		Yes1:                  proto.String(XorEncodeStr("yes", xorKey)),
		Yes2:                  proto.String(XorEncodeStr("yes", xorKey)),
		IosVersion:            proto.String(XorEncodeStr(IPadOsVersion, xorKey)),
		DeviceType:            proto.String(XorEncodeStr("iPhone", xorKey)),
		UnKnown2:              proto.Int32(2),
		IdentifierForVendor:   proto.String(XorEncodeStr(uuid1, xorKey)),
		AdvertisingIdentifier: proto.String(XorEncodeStr(uuid2, xorKey)),
		Carrier:               proto.String(XorEncodeStr("中国移动", xorKey)),
		BatteryInfo:           proto.Int32(1),
		NetworkName:           proto.String(XorEncodeStr("en0", xorKey)),
		NetType:               proto.Int32(1),
		AppBundleId:           proto.String(XorEncodeStr("com.tencent.xin", xorKey)),
		DeviceName:            proto.String(XorEncodeStr(DeviceName, xorKey)),
		UserName:              proto.String(XorEncodeStr(IPadModel, xorKey)),
		Unknown3:              proto.Int64(IOSDeviceNumber(Deviceid[:29] + "FFF")),
		Unknown4:              proto.Int64(IOSDeviceNumber(Deviceid[:29] + "OOO")),
		Unknown5:              proto.Int32(1),
		Unknown6:              proto.Int32(4),
		Lang:                  proto.String(XorEncodeStr("zh", xorKey)),
		Country:               proto.String(XorEncodeStr("CN", xorKey)),
		Unknown7:              proto.Int32(4),
		DocumentDir:           proto.String(XorEncodeStr("/var/mobile/Containers/Data/Application/"+IOSGetCidUUid(Deviceid, IOSGetCid(0x10101201))+"/Documents", xorKey)),
		Unknown8:              proto.Int32(0),
		Unknown9:              proto.Int32(0),
		HeadMD5:               proto.String(XorEncodeStr(IOSGetCidMd5(Deviceid, IOSGetCid(0x0262626262626)), xorKey)),
		AppUUID:               proto.String(XorEncodeStr(uuid1, xorKey)),
		SyslogUUID:            proto.String(XorEncodeStr(uuid2, xorKey)),
		Unknown10:             proto.String(""),
		Unknown11:             proto.String(""),
		AppName:               proto.String(XorEncodeStr("微信", xorKey)),
		SshPath:               proto.String(""),
		TempTest:              proto.String(""),
		DevMD5:                proto.String(""),
		DevUser:               proto.String(""),
		Unknown12:             proto.String(""),
		IsModify:              proto.Int32(0),
		ModifyMD5:             proto.String(""),
		RqtHash:               proto.Int64(288529533794259264),
		Unknown43:			   proto.Uint64(1586355322),
		Unknown44: 			   proto.Uint64(1586355519000),
		Unknown45: 			   proto.Uint64(0),
		Unknown46: 			   proto.Uint64(288529533794259264),
		Unknown47:			   proto.Uint64(0),
		Unknown48:			   proto.String(Deviceid),
		Unknown49:			   proto.String(""),
		Unknown50:			   proto.String(""),
		Unknown51:			   proto.String(XorEncodeStr(DeviceToken.GetTrustResponseData().GetSoftData().GetSoftConfig(),xorKey)),
		Unknown52:			   proto.Uint64(0),
		Unknown53:			   proto.String(""),
		Unknown54:			   proto.String(XorEncodeStr(DeviceToken.GetTrustResponseData().GetDeviceToken(),xorKey)),

	}
	wxFile := &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/var/containers/Bundle/Application/"+IOSGetCidUUid(Deviceid, IOSGetCid(0x098521236654))+"/WeChat.app/WeChat", xorKey)),
		Filepath: proto.String(XorEncodeStr(IOSGetCidUUid(Deviceid, IOSGetCid(0x30000001)), xorKey)),
	}
	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, wxFile)

	opensslFile := &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/var/containers/Bundle/Application/"+IOSGetCidUUid(Deviceid, IOSGetCid(0x098521236654))+"/WeChat.app/Frameworks/OpenSSL.framework/OpenSSL", xorKey)),
		Filepath: proto.String(XorEncodeStr(IOSGetCidUUid(Deviceid, IOSGetCid(0x30000002)), xorKey)),
	}
	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, opensslFile)

	protoFile := &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/var/containers/Bundle/Application/"+IOSGetCidUUid(Deviceid, IOSGetCid(0x098521236654))+"/WeChat.app/Frameworks/ProtobufLite.framework/ProtobufLite", xorKey)),
		Filepath: proto.String(XorEncodeStr(IOSGetCidUUid(Deviceid, IOSGetCid(0x30000003)), xorKey)),
	}
	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, protoFile)

	marsbridgenetworkFile := &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/var/containers/Bundle/Application/"+IOSGetCidUUid(Deviceid, IOSGetCid(0x098521236654))+"/WeChat.app/Frameworks/marsbridgenetwork.framework/marsbridgenetwork", xorKey)),
		Filepath: proto.String(XorEncodeStr(IOSGetCidUUid(Deviceid, IOSGetCid(0x30000004)), xorKey)),
	}
	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, marsbridgenetworkFile)

	matrixreportFile := &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/var/containers/Bundle/Application/"+IOSGetCidUUid(Deviceid, IOSGetCid(0x098521236654))+"/WeChat.app/Frameworks/matrixreport.framework/matrixreport", xorKey)),
		Filepath: proto.String(XorEncodeStr(IOSGetCidUUid(Deviceid, IOSGetCid(0x30000005)), xorKey)),
	}
	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, matrixreportFile)

	andromedaFile := &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/var/containers/Bundle/Application/"+IOSGetCidUUid(Deviceid, IOSGetCid(0x098521236654))+"/WeChat.app/Frameworks/andromeda.framework/andromeda", xorKey)),
		Filepath: proto.String(XorEncodeStr(IOSGetCidUUid(Deviceid, IOSGetCid(0x30000006)), xorKey)),
	}
	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, andromedaFile)

	marsFile := &mm.FileInfo{
		Fileuuid: proto.String(XorEncodeStr("/var/containers/Bundle/Application/"+IOSGetCidUUid(Deviceid, IOSGetCid(0x098521236654))+"/WeChat.app/Frameworks/mars.framework/mars", xorKey)),
		Filepath: proto.String(XorEncodeStr(IOSGetCidUUid(Deviceid, IOSGetCid(0x30000007)), xorKey)),
	}
	spamDataBody.AppFileInfo = append(spamDataBody.AppFileInfo, marsFile)
	srcdata, _ := proto.Marshal(spamDataBody)

	newClientCheckData := &mm.NewClientCheckData{
		C32Cdata:  proto.Int64(int64(crc32.ChecksumIEEE([]byte(srcdata)))),
		TimeStamp: proto.Int64(int64(timeStamp)),
		Databody:  srcdata,
	}

	ccddata, _ := proto.Marshal(newClientCheckData)

	//压缩数据
	compressdata := AE(ccddata)

	// Zero: 03加密改06加密
	// zt := new(ZT)
	// zt.Init()
	// encData := zt.Encrypt(compressdata)
	encData := SaeEncrypt06(compressdata)

	return encData
}
