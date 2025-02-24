package Tools

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

func DownloadVoice(Data DownloadVoiceParam) models.ResponseResult {

	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//NewMsgId, _ := strconv.ParseInt(Data.NewMsgId, 10, 64)
	MasterBufId, _ := strconv.ParseInt(Data.Bufid, 10, 64)

	req := &mm.DownloadVoiceRequest{
		MsgId:  proto.Uint32(Data.MsgId),
		Offset: proto.Uint32(0),
		Length: proto.Uint32(uint32(Data.Length)),
		BaseRequest: &mm.BaseRequest{
			SessionKey:    []byte{},
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		NewMsgId:     proto.Int64(0),
		ChatRoomName: proto.String(Data.FromUserName),
		MasterBufId:  proto.Int64(MasterBufId),
	}

	//序列化
	reqdata, _ := proto.Marshal(req)

	//发包
	protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/micromsg-bin/downloadvoice",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              128,
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
	Response := mm.DownloadVoiceResponse{}
	err = proto.Unmarshal(protobufdata, &Response)
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
		Data:    Response,
	}
}
