package comm

import (
	"encoding/binary"
	"errors"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/Mmtls"
	"wechatdll/models"
)

type SendPostData struct {
	Ip       string
	Host     string
	Cgiurl   string
	Proxy    models.ProxyInfo
	PackData Algorithm.PackData
}

func MmtlsInitialize(Proxy models.ProxyInfo, domain string) (*Mmtls.HttpClientModel, *Mmtls.MmtlsClient, error) {
	//生成mmtls公私钥
	V1PrivKey, V1pubKey := Algorithm.GetECDH415Key()
	V2PrivKey, V2pubKey := Algorithm.GetECDH415Key()

	Shakehandpubkey := Mmtls.Shakehandpubkey{
		V1PrivKey: V1PrivKey,
		V1pubKey:  V1pubKey,
		V2PrivKey: V2PrivKey,
		V2pubKey:  V2pubKey,
	}

	//初始化Mmtls
	httpclient := Mmtls.GenNewHttpClient(nil, domain)
	MmtlsClient, err := httpclient.InitMmtlsShake(domain, Proxy, Shakehandpubkey)

	if err != nil {
		return nil, &Mmtls.MmtlsClient{}, err
	}

	return httpclient, MmtlsClient, nil
}

func SendRequest(SENP SendPostData, MmtlsClient *Mmtls.MmtlsClient) (protobufdata, Cookie []byte, errtype int64, err error) {
	//logInfoS, _ := json.Marshal(SENP)
	//log.Infof("http发送前pb组包(%d): [%x]", len(SENP.PackData.Reqdata), SENP.PackData.Reqdata)
	senddata := Algorithm.Pack(SENP.PackData.Reqdata, SENP.PackData.Cgi, SENP.PackData.Uin, SENP.PackData.Sessionkey, SENP.PackData.Cookie, SENP.PackData.Clientsessionkey, SENP.PackData.Loginecdhkey, SENP.PackData.EncryptType, SENP.PackData.UseCompress)
	//logInfoR, _ := json.Marshal(senddata)
	//log.Infof("http发送时组包(%d): [%s]", len(senddata), hex.EncodeToString(senddata))
	if SENP.Host == "" {
		SENP.Host = Algorithm.MmtlsShortHost
	}

	//初始化Mmtls
	var httpclient *Mmtls.HttpClientModel
	if MmtlsClient == nil {
		httpclient, MmtlsClient, err = MmtlsInitialize(SENP.Proxy, SENP.Host)
		if err != nil {
			return nil, nil, 0, err
		}
	} else {
		httpclient = Mmtls.GenNewHttpClient(MmtlsClient, SENP.Host)
	}

	var response []byte
	if SENP.PackData.MMtlsClose == true {
		response, err = httpclient.POST(SENP.Cgiurl, senddata, SENP.Host, SENP.Proxy)
	} else {
		response, err = httpclient.MMtlsPost(SENP.Host, SENP.Cgiurl, senddata, SENP.Proxy)
	}

	if err != nil {
		return nil, nil, -1, err
	}

	if len(response) <= 31 {
		Ret, err := RetConst(response)
		if Ret == -13 {
			return nil, nil, Ret, errors.New("用户可能退出")
		}
		return nil, nil, Ret, errors.New("微信服务返回信息：" + err.Error())
	}

	if SENP.Cgiurl == "/cgi-bin/micromsg-bin/newsync" {
		protobufdata = Algorithm.UnpackBusinessPacketWithAesGcm(response, SENP.PackData.Uin, &Cookie, SENP.PackData.Serversessionkey)
	} else {
		protobufdata = Algorithm.UnpackBusinessPacket(response, SENP.PackData.Sessionkey, SENP.PackData.Uin, &Cookie)
	}

	if protobufdata != nil {
		return
	} else {
		return nil, nil, -8, errors.New("数据解密失败")
	}
}

func RetConst(data []byte) (int64, error) {
	var Ret int32
	Ret = BytesToInt32(data[2:10])
	return int64(Ret), errors.New(mm.RetConst_name[BytesToInt32(data[2:10])])
}

func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

func Int32ToBytes(i int32) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}
