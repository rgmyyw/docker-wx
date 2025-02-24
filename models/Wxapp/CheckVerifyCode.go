package Wxapp

import (
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"

	"github.com/golang/protobuf/proto"
)

func CheckVerifyCode(Data CheckVerifyCodeData) *mm.CheckVerifyCodeResponse {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return nil
	}

	req := &mm.CheckVerifyCodeRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Appid:      proto.String(Data.Appid),
		Mobile:     proto.String(Data.Mobile),
		VerifyCode: proto.String(Data.VerifyCode),
	}

	reqdata, err := proto.Marshal(req)

	if err != nil {
		return nil
	}

	//发包
	protobufdata, _, _, err := comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/mmbiz-bin/wxaapp/customphone/checkverifycode",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              0xa7d,
			Uin:              D.Uin,
			Cookie:           D.Cooike,
			Sessionkey:       D.Sessionkey,
			EncryptType:      5,
			Loginecdhkey:     D.RsaPublicKey,
			Clientsessionkey: D.Clientsessionkey,
			UseCompress:      false,
		},
	}, D.MmtlsKey)

	if err != nil {
		return nil
	}

	//解包
	Response := mm.CheckVerifyCodeResponse{}
	err = proto.Unmarshal(protobufdata, &Response)

	if err != nil {
		return nil
	}

	return &Response

}
