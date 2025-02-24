package Login

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"io"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/Mmtls"
	"wechatdll/models"

	"github.com/golang/protobuf/proto"
)

func AndroidGetDeviceToken(DeviceId string, info Algorithm.AndroidDeviceInfo, httpclient Mmtls.HttpClientModel, proxy models.ProxyInfo, domain string) (mm.TrustResponse, error) {
	td := &mm.TrustReq{
		Td: &mm.TrustData{
			Tdi: []*mm.TrustDeviceInfo{
				{Key: proto.String("IMEI"), Val: proto.String(info.AndriodImei(DeviceId))},
				{Key: proto.String("AndroidID"), Val: proto.String(info.AndriodID(DeviceId))},
				{Key: proto.String("PhoneSerial"), Val: proto.String(info.AndriodPhoneSerial(DeviceId))},
				{Key: proto.String("cid"), Val: proto.String("")},
				{Key: proto.String("WidevineDeviceID"), Val: proto.String(info.AndriodWidevineDeviceID(DeviceId))},
				{Key: proto.String("WidevineProvisionID"), Val: proto.String(info.AndriodWidevineProvisionID(DeviceId))},
				{Key: proto.String("GSFID"), Val: proto.String("")},
				{Key: proto.String("SoterID"), Val: proto.String("")},
				{Key: proto.String("SoterUid"), Val: proto.String("")},
				{Key: proto.String("FSID"), Val: proto.String(info.AndriodFSID(DeviceId))},
				{Key: proto.String("BootID"), Val: proto.String("")},
				{Key: proto.String("IMSI"), Val: proto.String("")},
				{Key: proto.String("PhoneNum"), Val: proto.String("")},
				{Key: proto.String("WeChatInstallTime"), Val: proto.String("1730105747")}, //1730105747
				{Key: proto.String("PhoneModel"), Val: proto.String(info.AndroidPhoneModel(DeviceId))},
				{Key: proto.String("BuildBoard"), Val: proto.String("bullhead")},
				{Key: proto.String("BuildBootloader"), Val: proto.String(info.AndroidBuildBoard(DeviceId))},
				{Key: proto.String("SystemBuildDate"), Val: proto.String("Fri Sep 28 23:37:27 UTC 2024")},
				{Key: proto.String("SystemBuildDateUTC"), Val: proto.String("1730103286")},
				{Key: proto.String("BuildFP"), Val: proto.String(info.AndroidBuildFP(DeviceId))},
				{Key: proto.String("BuildID"), Val: proto.String(info.AndroidBuildID(DeviceId))},
				{Key: proto.String("BuildBrand"), Val: proto.String("HUAWEI")},
				{Key: proto.String("BuildDevice"), Val: proto.String("bullhead")},
				{Key: proto.String("BuildProduct"), Val: proto.String("bullhead")},
				{Key: proto.String("Manufacturer"), Val: proto.String(info.AndroidManufacturer(DeviceId))},
				{Key: proto.String("RadioVersion"), Val: proto.String(info.AndroidRadioVersion(DeviceId))},
				{Key: proto.String("AndroidVersion"), Val: proto.String(info.AndroidVersion())},
				{Key: proto.String("SdkIntVersion"), Val: proto.String("34")},
				{Key: proto.String("ScreenWidth"), Val: proto.String("1080")},
				{Key: proto.String("ScreenHeight"), Val: proto.String("1794")},
				{Key: proto.String("SensorList"), Val: proto.String("BMI160 accelerometer#Bosch#0.004788#1,BMI160 gyroscope#Bosch#0.000533#1,BMM150 magnetometer#Bosch#0.000000#1,BMP280 pressure#Bosch#0.005000#1,BMP280 temperature#Bosch#0.010000#1,RPR0521 Proximity Sensor#Rohm#1.000000#1,RPR0521 Light Sensor#Rohm#10.000000#1,Orientation#Google#1.000000#1,BMI160 Step detector#Bosch#1.000000#1,Significant motion#Google#1.000000#1,Gravity#Google#1.000000#1,Linear Acceleration#Google#1.000000#1,Rotation Vector#Google#1.000000#1,Geomagnetic Rotation Vector#Google#1.000000#1,Game Rotation Vector#Google#1.000000#1,Pickup Gesture#Google#1.000000#1,Tilt Detector#Google#1.000000#1,BMI160 Step counter#Bosch#1.000000#1,BMM150 magnetometer (uncalibrated)#Bosch#0.000000#1,BMI160 gyroscope (uncalibrated)#Bosch#0.000533#1,Sensors Sync#Google#1.000000#1,Double Twist#Google#1.000000#1,Double Tap#Google#1.000000#1,Device Orientation#Google#1.000000#1,BMI160 accelerometer (uncalibrated)#Bosch#0.004788#1")},
				{Key: proto.String("DefaultInputMethod"), Val: proto.String("com.google.android.inputmethod.latin")},
				{Key: proto.String("InputMethodList"), Val: proto.String("Google \345\215\260\345\272\246\350\257\255\351\224\256\347\233\230#com.google.android.apps.inputmethod.hindi,Google \350\257\255\351\237\263\350\276\223\345\205\245#com.google.android.googlequicksearchbox,Google \346\227\245\350\257\255\350\276\223\345\205\245\346\263\225#com.google.android.inputmethod.japanese,Google \351\237\251\350\257\255\350\276\223\345\205\245\346\263\225#com.google.android.inputmethod.korean,Gboard#com.google.android.inputmethod.latin,\350\260\267\346\255\214\346\213\274\351\237\263\350\276\223\345\205\245\346\263\225#com.google.android.inputmethod.pinyin")},
				{Key: proto.String("DeviceID"), Val: proto.String(DeviceId)},
				{Key: proto.String("OAID"), Val: proto.String("")},
			},
		},
	}

	pb, _ := proto.Marshal(td)

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(pb)
	w.Close()

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
			ClientVersion: proto.Int32(int32(Algorithm.AndroidVersion)),
			DeviceType:    []byte(Algorithm.AndroidDeviceType),
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
	hec.Init("Android")
	hecData := hec.HybridEcdhPackAndroidEn(3789, 10002, 0, nil, reqdata)
	recvData, err := httpclient.MMtlsPost(domain, "/cgi-bin/micromsg-bin/fpinitnl", hecData, proxy)
	if err != nil {
		return mm.TrustResponse{}, err
	}
	ph := hec.HybridEcdhPackAndroidUn(recvData)
	DTResp := &mm.TrustResponse{}
	_ = proto.Unmarshal(ph.Data, DTResp)
	return *DTResp, nil
}
