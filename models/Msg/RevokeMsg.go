package Msg

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

type RevokeMsgParam struct {
	Wxid        string
	NewMsgId    uint64
	ToUserName  string
	CreateTime  uint64
	ClientMsgId uint64
}

func RevokeMsg(Data RevokeMsgParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//组包
	req := &mm.RevokeMsgRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		ClientMsgId:    proto.String(""),
		NewClientMsgId: proto.Uint64(Data.ClientMsgId),
		CreateTime:     proto.Uint64(Data.CreateTime),
		IndexOfRequest: proto.Uint64(0),
		FromUserName:   proto.String(Data.Wxid),
		ToUserName:     proto.String(Data.ToUserName),
		MsgId:          proto.Uint64(0),
		NewMsgId:       proto.Uint64(Data.NewMsgId),
	}

	//序列化
	reqdata, _ := proto.Marshal(req)

	//发包
	protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/micromsg-bin/revokemsg",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              594,
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
		return models.ResponseResult{
			Code:    errtype,
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}
	}

	//解包
	RevokeMsgResponse := mm.RevokeMsgResponse{}
	err = proto.Unmarshal(protobufdata, &RevokeMsgResponse)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    RevokeMsgResponse,
	}

}
