package Msg

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Report"
)

func SendTransmitMsg(Data SendTransmitParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//提交个状态
	Report.Statusnotify(Data.Wxid, Data.ToWxid)

	T := time.Now().Unix()

	req := &mm.SendAppMsgRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Msg: &mm.AppMsg{
			FromUserName: proto.String(D.Wxid),
			AppId:        proto.String(""),
			SdkVersion:   proto.Int32(0),
			ToUserName:   proto.String(Data.ToWxid),
			Type:         proto.Int32(19),
			Content:      proto.String(Data.Xml),
			CreateTime:   proto.Int64(T),
			ClientMsgId:  proto.String(fmt.Sprintf("%v_%v", Data.ToWxid, T)),
			Source:       proto.Int32(3),
			MsgSource:    proto.String("<msgsource><bizflag>0</bizflag></msgsource>"),
		},
		MsgForwardType: proto.Int32(2),
		SendMsgTicket:  proto.String(""),
	}

	//序列化
	reqdata, _ := proto.Marshal(req)

	//发包
	protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/micromsg-bin/sendappmsg",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              222,
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
		return models.ResponseResult{
			Code:    errtype,
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}
	}

	//解包
	NewSendMsgRespone := mm.SendAppMsgResponse{}
	err = proto.Unmarshal(protobufdata, &NewSendMsgRespone)
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
		Data:    NewSendMsgRespone,
	}
}
