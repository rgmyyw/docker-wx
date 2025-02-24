package Report

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
)

func Statusnotify(Wxid, ToWxid string) {
	D, err := comm.GetLoginata(Wxid)
	if err != nil {
		return
	}

	req := &mm.StatusNotifyRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Code:         proto.Uint32(2),
		FromUserName: proto.String(Wxid),
		ToUserName:   proto.String(ToWxid),
		ClientMsgId:  proto.String(fmt.Sprintf("%v_%v", ToWxid, time.Now().Unix())),
	}

	reqdata, err := proto.Marshal(req)

	if err != nil {
		return
	}

	//发包
	_, _, _, err = comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/micromsg-bin/statusnotify",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              251,
			Uin:              D.Uin,
			Cookie:           D.Cooike,
			Sessionkey:       D.Sessionkey,
			EncryptType:      5,
			Loginecdhkey:     D.Loginecdhkey,
			Clientsessionkey: D.Clientsessionkey,
			UseCompress:      false,
		},
	}, D.MmtlsKey)

	if err != nil {
		return
	}

	return
}
