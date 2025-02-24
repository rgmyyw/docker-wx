package Login

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/Mmtls"
	"wechatdll/models"

	"github.com/golang/protobuf/proto"
)

func IPadGetDeviceToken(DeviceId, DeviceType, DeviceName, OsVersion string, Version int32, httpclient Mmtls.HttpClientModel, proxy models.ProxyInfo, domain string) (mm.TrustResponse, error) {
	uuid1, uuid2 := Algorithm.IOSUuid(DeviceId)
	installTime := strconv.FormatInt(time.Now().Add(-2592234*time.Second).Unix(), 10)       // 一个月前(偏移234秒)秒级时间戳
	kernBootTime := strconv.FormatInt(time.Now().Add(-2592230*time.Second).Unix()*1000, 10) // 一个月前(偏移230秒)毫秒秒级时间戳
	td := &mm.TrustReq{
		Td: &mm.TrustData{
			Tdi: []*mm.TrustDeviceInfo{
				{Key: proto.String("deviceid"), Val: proto.String(DeviceId)},
				{Key: proto.String("sdi"), Val: proto.String(Algorithm.IOSGetCidMd5(DeviceId, Algorithm.IOSGetCid(0x0262626262626)))},
				{Key: proto.String("idfv"), Val: proto.String(uuid1)},
				{Key: proto.String("idfa"), Val: proto.String(uuid2)},
				{Key: proto.String("device_model"), Val: proto.String(DeviceType)},
				{Key: proto.String("os_version"), Val: proto.String(OsVersion)},
				{Key: proto.String("core_count"), Val: proto.String("6")},
				{Key: proto.String("carrier_name"), Val: proto.String("")},
				{Key: proto.String("is_jailbreak"), Val: proto.String("0")},
				{Key: proto.String("device_name"), Val: proto.String(DeviceName)},
				{Key: proto.String("client_version"), Val: proto.String(fmt.Sprintf("%v", Version))},
				{Key: proto.String("plist_version"), Val: proto.String(fmt.Sprintf("%v", Version))},
				{Key: proto.String("language"), Val: proto.String("zh")},
				{Key: proto.String("locale_country"), Val: proto.String("CN")},
				{Key: proto.String("screen_width"), Val: proto.String("834")},
				{Key: proto.String("screen_height"), Val: proto.String("1112")},
				{Key: proto.String("install_time"), Val: proto.String(installTime)},
				{Key: proto.String("kern_boottime"), Val: proto.String(kernBootTime)},
			},
		},
	}
	// pb包序列化为二进制流
	pb, _ := proto.Marshal(td)
	// zlib压缩
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(pb)
	w.Close()
	// zt异或加密
	zt := new(Algorithm.ZT)
	zt.Init()

	encData := zt.Encrypt(b.Bytes())

	randKey := make([]byte, 16)
	io.ReadFull(rand.Reader, randKey)

	fp := &mm.FPFresh{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    []byte{},
			Uin:           proto.Uint32(0),
			DeviceId:      append([]byte(DeviceId), 0),
			ClientVersion: proto.Int32(Version),
			DeviceType:    []byte(DeviceType),
			Scene:         proto.Uint32(0),
		},
		SessKey: randKey,
		Ztdata: &mm.ZTData{
			Version:   proto.String("00000003\x00"),
			Encrypted: proto.Uint32(1),
			Data:      encData,
			TimeStamp: proto.Uint32(uint32(time.Now().Unix())),
			Optype:    proto.Uint32(5),
			Uin:       proto.Uint32(0),
		},
	}

	reqdata, _ := proto.Marshal(fp)

	hec := &Algorithm.Client{}
	if Version == int32(Algorithm.AndroidPadVersion) {
		hec.Init("AndroidPad")
	}
	if Version == int32(Algorithm.AndroidPadVersionx) {
		hec.Init("AndroidPad")
	}
	if Version == int32(Algorithm.IPadVersion) {
		hec.Init("IOS")
	}
	if Version == int32(Algorithm.IPadVersionx) {
		hec.Init("IOS")
	}
	if Version == int32(Algorithm.WinVersion) {
		hec.Init("Windows")
	}
	if Version == int32(Algorithm.WinUwpVersion) {
		hec.Init("WindowsUwp")
	}
	if Version == int32(Algorithm.CarVersion) {
		hec.Init("Car")
	}
	if Version == int32(Algorithm.MacVersion) {
		hec.Init("MAC")
	}

	hecData := hec.HybridEcdhPackIosEn(3789, 0, nil, reqdata)
	recvData, err := httpclient.MMtlsPost(domain, "/cgi-bin/micromsg-bin/fpinitnl", hecData, proxy)
	if err != nil {
		return mm.TrustResponse{}, err
	}
	ph := hec.HybridEcdhPackIosUn(recvData)
	DTResp := &mm.TrustResponse{}
	_ = proto.Unmarshal(ph.Data, DTResp)
	return *DTResp, nil
}
