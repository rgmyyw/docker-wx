package Group

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Tools"
)

func SetChatRoomName(Data OperateChatRoomInfoParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	GetContact := Tools.GetContact(Tools.GetContactParam{
		Wxid:         Data.Wxid,
		UserNameList: Data.QID,
	})

	if GetContact.Data == nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", GetContact.Message),
			Data:    nil,
		}
	}

	ModContact := &mm.ModContact{
		UserName: &mm.SKBuiltinStringT{
			String_: proto.String(Data.QID),
		},
		NickName: &mm.SKBuiltinStringT{
			String_: proto.String(Data.Content),
		},
	}

	buffer, err := proto.Marshal(ModContact)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}

	cmdItem := mm.CmdItem{
		CmdId: proto.Int32(27),
		CmdBuf: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(buffer))),
			Buffer: buffer,
		},
	}

	var cmdItems []*mm.CmdItem
	cmdItems = append(cmdItems, &cmdItem)

	req := &mm.OpLogRequest{
		Cmd: &mm.CmdList{
			Count: proto.Uint32(uint32(len(cmdItems))),
			List:  cmdItems,
		},
	}

	reqdata, err := proto.Marshal(req)

	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//发包
	protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/micromsg-bin/oplog",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              681,
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
	GetContactResponse := mm.OplogResponse{}
	err = proto.Unmarshal(protobufdata, &GetContactResponse)

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
		Data:    GetContactResponse,
	}
}
