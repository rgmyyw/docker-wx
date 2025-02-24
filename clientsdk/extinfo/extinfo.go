package extinfo

import (
	"fmt"
	"github.com/gogf/guuid"
	"github.com/gogo/protobuf/proto"
	"hash/crc32"
	"math/rand"
	"strings"
	"time"
	"wechatdll/Cilent/wechat"
	"wechatdll/clientsdk/ccdata"
	"wechatdll/clientsdk/mmproto"
	"wechatdll/comm"
	"wechatdll/models/baseutils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type GetCcDataRep struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

func MakeXorKey(key int64) uint8 {
	var un int64 = int64(0xffffffed)
	xorKey := (uint8)(key*un + 7)
	return xorKey
}

func exponent(a, n uint64) uint64 {
	result := uint64(1)
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			result *= a
		}
		a *= a
	}
	return result
}

func Hex2int(hexB *[]byte) uint64 {
	var retInt uint64
	hexLen := len(*hexB)
	for k, v := range *hexB {
		retInt += b2m_map[v] * exponent(16, uint64(2*(hexLen-k-1)))
	}
	return retInt
}

func DeviceNumber(DeviceId string) int64 {
	ssss := []byte(baseutils.Md5Value(DeviceId))
	ccc := Hex2int(&ssss) >> 8
	ddd := ccc + 60000000000000000
	if ddd > 80000000000000000 {
		ddd = ddd - (80000000000000000 - ddd)
	}
	return int64(ddd)
}

var wifiPrefix = []string{"TP_", "360_", "ChinaNet-", "MERCURY_", "DL-", "VF_", "HUAW-"}

func BuildRandomWifiSsid() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	i := r.Intn(len(wifiPrefix))
	randChar := make([]byte, 6)
	for x := 0; x < 6; x++ {
		randChar[x] = byte(r.Intn(26) + 65)
	}
	return wifiPrefix[i] + string(randChar)
}

// 获取FilePathCrc
func GetFileInfo(xorKey byte) []*wechat.FileInfo {
	return []*wechat.FileInfo{
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/20278CED-292E-4C4F-930C-0D1E86B02AD1/WeChat.app/WeChat", xorKey)),
			Fileuuid: proto.String(XorEncrypt("B069D479-E08E-3557-B7C7-411AF31B4919", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/20278CED-292E-4C4F-930C-0D1E86B02AD1/WeChat.app/Frameworks/OpenSSL.framework/OpenSSL", xorKey)),
			Fileuuid: proto.String(XorEncrypt("7CED8A7F-509A-3820-9EBC-8EB3AE5B99D9", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/20278CED-292E-4C4F-930C-0D1E86B02AD1/WeChat.app/Frameworks/ProtobufLite.framework/ProtobufLite", xorKey)),
			Fileuuid: proto.String(XorEncrypt("69971FE3-4728-3F01-9137-4FD3FFC26AB4", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/20278CED-292E-4C4F-930C-0D1E86B02AD1/WeChat.app/Frameworks/marsbridgenetwork.framework/marsbridgenetwork", xorKey)),
			Fileuuid: proto.String(XorEncrypt("CFED9A03-C881-3D50-B014-732D0A09879F", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/zstd.framework/zstd", xorKey)),
			Fileuuid: proto.String(XorEncrypt("F326111B-2EBA-34E8-8830-657315905F43", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/TXLiteAVSDK_Smart_No_VOD.framework/TXLiteAVSDK_Smart_No_VOD", xorKey)),
			Fileuuid: proto.String(XorEncrypt("606505A0-69D9-3EBC-A4DB-758D54BE46AA", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/20278CED-292E-4C4F-930C-0D1E86B02AD1/WeChat.app/Frameworks/matrixreport.framework/matrixreport", xorKey)),
			Fileuuid: proto.String(XorEncrypt("1E7F06D2-DD36-31A8-AF3B-00D62054E1F9", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/20278CED-292E-4C4F-930C-0D1E86B02AD1/WeChat.app/Frameworks/andromeda.framework/andromeda", xorKey)),
			Fileuuid: proto.String(XorEncrypt("0AE6A3E2-31A4-352E-9CAC-A1011043951A", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/YTFaceProSDK.framework/YTFaceProSDK", xorKey)),
			Fileuuid: proto.String(XorEncrypt("4F58B750-3134-36A2-9524-147872E9A607", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/GPUImage.framework/GPUImage", xorKey)),
			Fileuuid: proto.String(XorEncrypt("847B0606-87CA-3896-AE16-90DF0370DC94", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/WCDB.framework/WCDB", xorKey)),
			Fileuuid: proto.String(XorEncrypt("09941A29-AF95-3C1F-A6EC-A571EC4529FC", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/MMCommon.framework/MMCommon", xorKey)),
			Fileuuid: proto.String(XorEncrypt("917922FE-F15C-3A0C-9F9B-29DAE8BC08AF", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/MultiMedia.framework/MultiMedia", xorKey)),
			Fileuuid: proto.String(XorEncrypt("7FBD2D5F-806D-38DE-A54F-AE471B3F342F", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/QBar.framework/QBar", xorKey)),
			Fileuuid: proto.String(XorEncrypt("95DE8BA9-6A55-3642-8B1F-363BB507E2CD", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/QMapKit.framework/QMapKit", xorKey)),
			Fileuuid: proto.String(XorEncrypt("231A339F-6CBC-39BC-A88B-96FF1282390F", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/3C5AC4DE-D87D-4669-B996-D839E9100456/WeChat.app/Frameworks/ConfSDK.framework/ConfSDK", xorKey)),
			Fileuuid: proto.String(XorEncrypt("BF039435-9428-3118-B245-219EADDB825A", xorKey)),
		},
		{
			Filepath: proto.String(XorEncrypt("/var/containers/Bundle/Application/20278CED-292E-4C4F-930C-0D1E86B02AD1/WeChat.app/Frameworks/mars.framework/mars", xorKey)),
			Fileuuid: proto.String(XorEncrypt("B14DA1FF-1E28-32C4-B198-672A610690A7", xorKey)),
		},
	}
}

func CheckSoftType5() uint32 {
	sec := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(999) * 1000
	v79 := uint32(sec)&0xe | 1
	key := v79

	v77 := uint32(134217728)
	n := uint32(4)

	for true {
		dwTmp := n & 3
		if dwTmp == 0 {
			v79 = (3877*v79 + 5) & 0xf
		}

		dwTmp = uint32(((int(v79) >> int(dwTmp)) & 1)) << int(n)
		v77 ^= dwTmp
		n++
		if n == 24 {
			break
		}
	}
	return v77 | key
}

// 我爱文文
func GetNewSpamData(iosVersion string, deviceType, uuid1, uuid2, deviceName string, deviceToken *wechat.TrustResp, deviceId, userName, guid2 string, userInfo *comm.LoginData) ([]byte, error) {
	dateTimeSramp := time.Now().Unix()
	timeStamp := time.Now().Unix()
	xorKey := MakeXorKey(timeStamp)
	Lang := "zh"
	Country := "CN"
	if userInfo.DeviceInfo != nil {
		Lang = userInfo.DeviceInfo.Language
		Country = Lang
	}
	Unknown106 := CheckSoftType5()
	spamDataBody := wechat.SpamDataBody{
		UnKnown1:              proto.Int32(1),
		TimeStamp:             proto.Uint32(uint32(dateTimeSramp)),
		KeyHash:               proto.Int32(int32(MakeKeyHash(xorKey))),
		Yes1:                  proto.String(XorEncrypt("yes", xorKey)),
		Yes2:                  proto.String(XorEncrypt("yes", xorKey)),
		IosVersion:            proto.String(XorEncrypt(iosVersion, xorKey)), // 14.3.0
		DeviceType:            proto.String(XorEncrypt("iPad", xorKey)),
		UnKnown2:              proto.Int32(2), //cpu核数
		IdentifierForVendor:   proto.String(XorEncrypt(uuid1, xorKey)),
		AdvertisingIdentifier: proto.String(XorEncrypt(uuid2, xorKey)),
		Carrier:               proto.String(XorEncrypt("中国移动", xorKey)),
		BatteryInfo:           proto.Int32(1),
		NetworkName:           proto.String(XorEncrypt("en0", xorKey)),
		NetType:               proto.Int32(1),
		AppBundleId:           proto.String(XorEncrypt(userInfo.DeviceInfo.BundleID, xorKey)),
		DeviceName:            proto.String(XorEncrypt("iPad Pro", xorKey)),
		UserName:              proto.String(XorEncrypt(deviceName, xorKey)),
		Unknown3:              proto.Int64(77968568554357776), //基带版本   77968568550229002
		Unknown4:              proto.Int64(77968568554357760), //基带通讯版本   77968568550228991
		Unknown5:              proto.Int32(5),                 //IsJailbreak
		Unknown6:              proto.Int32(4),
		Lang:                  proto.String(XorEncrypt(Lang, xorKey)),    //zh
		Country:               proto.String(XorEncrypt(Country, xorKey)), //CN
		Unknown7:              proto.Int32(4),
		DocumentDir:           proto.String(XorEncrypt("/var/mobile/Containers/Data/Application/857C3940-0044-43E9-AD59-6C3240DACD91/Documents", xorKey)),
		Unknown8:              proto.Int32(0),
		Unknown9:              proto.Int32(0),
		HeadMD5:               proto.String(XorEncrypt("d55ce16228afb0ea5205380af376761e", xorKey)), //XorEncrypt("d13610700984cf481b7d3f5fa2011c30", xorKey)
		AppUUID:               proto.String(XorEncrypt(uuid1, xorKey)),
		SyslogUUID:            proto.String(XorEncrypt(uuid2, xorKey)),
		AppName:               proto.String(XorEncrypt("微信", xorKey)),
		SshPath:               proto.String(XorEncrypt("/usr/bin/ssh", xorKey)),
		TempTest:              proto.String(XorEncrypt("/tmp/test.txt", xorKey)), //XorEncrypt("/tmp/test.txt", xorKey)
		Unknown12:             proto.String(XorEncrypt("yyyy/MM/dd HH:mm", xorKey)),
		IsModify:              proto.Int32(0),
		ModifyMD5:             proto.String(XorEncrypt("83d6ad7f1c5045ab8112a8411c8091f2", xorKey)),
		RqtHash:               proto.Int64(288512216272273664),
		Unknown13:             proto.Int32(0),
		Unknown14:             proto.Int32(0),
		Ssid:                  proto.String(XorEncrypt("F3:48:B2:4B:34:EB", xorKey)),
		Unknown15:             proto.Int32(0),
		Bssid:                 proto.String(XorEncrypt("F3:48:B2:4B:34:EC", xorKey)),
		IsJail:                proto.Int32(0),
		Seid:                  proto.String(XorEncrypt("D2:1B:61:E1:DE:C6", xorKey)),
		Unknown16:             proto.Int32(60),
		Unknown17:             proto.Int32(61),
		Unknown18:             proto.Int32(62),
		WifiOn:                proto.Int32(1),
		WifiName:              proto.String(XorEncrypt("Xiaomi_L52C", xorKey)),
		WifiMac:               proto.String(XorEncrypt(userInfo.DeviceInfo.DeviceMac, xorKey)),
		BluethOn:              proto.Int32(0),
		BluethName:            proto.String(XorEncrypt("iPad Pro", xorKey)),
		BluethMac:             proto.String(XorEncrypt("F3:48:B2:4B:7B:65", xorKey)),
		Unknown19:             proto.Int32(67),
		Unknown20:             proto.Int32(68),
		Unknown26:             proto.Int32(69),
		HasSim:                proto.Int32(0),
		UsbState:              proto.Int32(0),
		Unknown27:             proto.Int32(1300),
		Unknown28:             proto.Int32(73),
		Sign:                  proto.String(XorEncrypt("", xorKey)),
		PackageFlag:           proto.Uint32(0x01),
		AccessFlag:            proto.Uint32(0x03),
		Imei:                  proto.String(XorEncrypt(userInfo.DeviceInfo.Imei, xorKey)),
		DevMD5:                proto.String(XorEncrypt("0582444e4a124e1cfb350384140fb5f4", xorKey)),
		DevUser:               proto.String(XorEncrypt("iPad", xorKey)),
		DevPrefix:             proto.String(XorEncrypt("Apple", xorKey)),
		DevSerial:             proto.String(XorEncrypt(strings.ReplaceAll(guuid.New().String(), "-", "")[0:16], xorKey)),
		Unknown29:             proto.Uint32(0x1d),
		Unknown30:             proto.Uint32(0x1e),
		Unknown31:             proto.Uint32(0x1f),
		Unknown32:             proto.Uint32(0x20),
		AppNum:                proto.Uint32(0x10),
		Totcapacity:           proto.String(XorEncrypt("0x200", xorKey)),
		Avacapacity:           proto.String(XorEncrypt("0x6685f", xorKey)),
		Unknown33:             proto.Uint32(0x21),
		Unknown34:             proto.Uint32(0x28),
		Unknown35:             proto.Uint32(0x69),
		Unknown103:            proto.Int32(0),
		Unknown104:            proto.Int32(0),
		Unknown105:            proto.Int32(0),
		Unknown106:            &Unknown106,
		Unknown107:            proto.Int32(107),
		Unknown108:            proto.Int32(0),
		Unknown109:            proto.Int32(109),
		Unknown110:            proto.Int32(0),
		Unknown111:            proto.Int32(13),
		Unknown112:            proto.Uint32(0xf5),
		AppFileInfo:           GetFileInfo(xorKey),
	}

	data, err := proto.Marshal(&spamDataBody)
	if err != nil {
		return nil, err
	}

	newClientCheckData := &wechat.NewClientCheckData{
		C32CData:  proto.Int64(int64(crc32.ChecksumIEEE(data))),
		TimeStamp: proto.Int64(time.Now().Unix()),
		DataBody:  data,
	}

	ccData, err := proto.Marshal(newClientCheckData)
	if err != nil {
		return nil, err
	}

	afterCompressionCCData := baseutils.CompressByteArray(ccData)
	afterEnData, err := ccdata.EncodeZipData(afterCompressionCCData, 0x3060)
	if err != nil {
		return nil, err
	}

	//压缩数据
	//compressdata := DoZlibCompress(ccData)
	////compressdata := AE(ccddata)
	//
	//zt := new(ZT)
	//zt.Init()
	//encData := zt.Encrypt(compressdata)

	//压缩数据
	//compressdata := AE(ccData)

	// Zero: 03加密改06加密
	// zt := new(ZT)
	// zt.Init()
	// encData := zt.Encrypt(compressdata)
	//encData := Algorithm.SaeEncrypt06(compressdata)

	deviceRunningInfo := &wechat.DeviceRunningInfoNew{
		Version:     []byte("00000003"),
		Type:        proto.Uint32(1),
		EncryptData: afterEnData,
		Timestamp:   proto.Uint32(uint32(timeStamp)),
		Unknown5:    proto.Uint32(5),
		Unknown6:    proto.Uint32(0),
	}
	return proto.Marshal(deviceRunningInfo)
}

// 获取DeviceToken
func GetDeviceToken(deviceToken string) *mmproto.DeviceToken {
	curtime := uint32(time.Now().Unix())
	return &mmproto.DeviceToken{
		Version:   proto.String(""),
		Encrypted: proto.Uint32(1),
		Data: &mmproto.SKBuiltinStringt{
			String_: proto.String(deviceToken),
		},
		TimeStamp: &curtime,
		Optype:    proto.Uint32(2),
		Uin:       proto.Uint32(0),
	}
}
func GenGUId(DeviceId, Cid string) string {
	Md5Data := baseutils.Md5Value(DeviceId + Cid)
	return fmt.Sprintf("%x-%x-%x-%x-%x", Md5Data[0:8], Md5Data[2:6], Md5Data[3:7], Md5Data[1:5], Md5Data[20:32])
}

// 获取CCD
func GetCCDPbLib(iosVersion, deviceType, uuid1, uuid2, deviceName string, deviceToken *wechat.TrustResp, deviceId, userName, guid2 string, userInfo *comm.LoginData) ([]byte, error) {
	ccData1, err := GetNewSpamData(iosVersion, deviceType, uuid1, uuid2, deviceName, deviceToken, deviceId, userName, guid2, userInfo)
	if err != nil {
		return nil, err
	}
	deviceTokenObj := GetDeviceToken(deviceToken.GetTrustResponseData().GetDeviceToken())
	dt, err := proto.Marshal(deviceTokenObj)
	if err != nil {
		return nil, err
	}

	wcExtInfo := &wechat.WCExtInfoNew{

		CcData: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(ccData1))),
			Buffer: ccData1,
		},
		DeviceToken: &wechat.BufferT{
			ILen:   proto.Uint32(uint32(len(dt))),
			Buffer: dt,
		},
	}
	return proto.Marshal(wcExtInfo)
}
