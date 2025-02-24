package Friend

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strings"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

type GetContractDetailparameter struct {
	Wxid     string
	Towxids  string
	ChatRoom string
}

func GetContractDetail(Data GetContractDetailparameter) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	TowxdsSplit := strings.Split(Data.Towxids, ",")

	Towxds := make([]*mm.SKBuiltinStringT, len(TowxdsSplit))

	if len(TowxdsSplit) >= 1 {
		for i, v := range TowxdsSplit {
			Towxds[i] = &mm.SKBuiltinStringT{
				String_: proto.String(v),
			}
		}
	}

	ChatRoom := &mm.SKBuiltinStringT{}
	ChatRoomCount := uint32(1)

	if Data.ChatRoom != "" {
		ChatRoomCount = 1
		ChatRoom = &mm.SKBuiltinStringT{
			String_: proto.String(Data.ChatRoom),
		}
	} else {
		ChatRoom = nil
		ChatRoomCount = uint32(0)
	}

	req := &mm.GetContactRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		UserCount:         proto.Int32(int32(len(Towxds))),
		UserNameList:      Towxds,
		FromChatRoomCount: proto.Int32(int32(ChatRoomCount)),
		FromChatRoom:      ChatRoom,
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
		Cgiurl: "/cgi-bin/micromsg-bin/getcontact",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              182,
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
	Response := mm.GetContactResponse{}
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
