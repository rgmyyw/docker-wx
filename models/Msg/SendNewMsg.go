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

type SendNewMsgParam struct {
	Wxid    string
	ToWxid  string
	Content string
	Type    int64
	At		string
}

func SendNewMsg(Data SendNewMsgParam) models.ResponseResult {
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

	msgUtc := time.Now().Unix()
	msgId := msgUtc - 107961031

	//消息组包
	MsgRequest := &mm.NewSendMsgRequest{
		Cnt: proto.Int32(1),
		Info: &mm.ChatInfo{
			Toid: &mm.SKBuiltinStringT{
				String_: proto.String(Data.ToWxid),
			},
			Content:     proto.String(Data.Content),
			Type:        proto.Int64(Data.Type),
			Utc:         proto.Int64(msgUtc),
			ClientMsgId: proto.Uint64(uint64(msgId)),
			MsgSource:   nil,
		},
	}

	//群@
	if Data.At != "" {
		MsgRequest.Info.MsgSource = proto.String("<msgsource><atuserlist>"+Data.At+"</atuserlist><bizflag>0</bizflag></msgsource>")
	}

	//序列化
	reqdata, _ := proto.Marshal(MsgRequest)

	//fmt.Println(hex.EncodeToString(reqdata))

	//发包
	protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/micromsg-bin/newsendmsg",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              522,
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
	NewSendMsgRespone := mm.NewSendMsgRespone{}
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
