package Friend

import (
	"fmt"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"

	"github.com/golang/protobuf/proto"
)

type SendRequestParam struct {
	Wxid          string
	V1            string
	V2            string
	Opcode        int32
	Scene         int
	VerifyContent string
}

func SendRequest(Data SendRequestParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	VerifyUserList := make([]*mm.VerifyUser, 0)

	VerifyUserList = append(VerifyUserList, &mm.VerifyUser{
		Value:               proto.String(Data.V1),
		VerifyUserTicket:    proto.String(""),
		AntispamTicket:      proto.String(Data.V2),
		FriendFlag:          proto.Uint32(0),
		ChatRoomUserName:    proto.String(""),
		SourceUserName:      proto.String(""),
		SourceNickName:      proto.String(""),
		ScanQrcodeFromScene: proto.Uint32(0),
		ReportInfo:          proto.String(""),
		OuterUrl:            proto.String(""),
		SubScene:            proto.Int32(0),
	})

	ccData := &mm.CryptoData{
		Version:     []byte("00000006"),
		Type:        proto.Uint32(1),
		EncryptData: Algorithm.GetiPadNewSpamData(D.Deviceid_str, D.DeviceName, D.DeviceToken),
		Timestamp:   proto.Uint32(uint32(time.Now().Unix())),
		Unknown5:    proto.Uint32(5),
		Unknown6:    proto.Uint32(0),
	}
	ccDataseq, _ := proto.Marshal(ccData)

	// TODO: 邀请好友使用01,03还是06加密?
	Wcstf := Algorithm.IpadWcstf(Data.Wxid)
	Wcste := Algorithm.IpadWcste(0, 0)

	DeviceTokenCCD := &mm.DeviceToken{
		Version:   proto.String(""),
		Encrypted: proto.Uint32(1),
		Data: &mm.SKBuiltinStringT{
			String_: proto.String(D.DeviceToken.GetTrustResponseData().GetDeviceToken()),
		},
		TimeStamp: proto.Uint32(uint32(time.Now().Unix())),
		Optype:    proto.Uint32(2),
		Uin:       proto.Uint32(0),
	}
	DeviceTokenCCDPB, _ := proto.Marshal(DeviceTokenCCD)

	WCExtInfo := &mm.WCExtInfo{
		Wcstf: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(Wcstf))),
			Buffer: Wcstf,
		},
		Wcste: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(Wcste))),
			Buffer: Wcste,
		},
		CcData: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(ccDataseq))),
			Buffer: ccDataseq,
		},
		DeviceToken: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(DeviceTokenCCDPB))),
			Buffer: DeviceTokenCCDPB,
		},
	}

	WCExtInfoseq, _ := proto.Marshal(WCExtInfo)
	/*

	if Data.Opcode != 1 || Data.Opcode != 2 {
		Data.Opcode = 2
	}
		*/

	req := &mm.VerifyUserRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		SceneListCount:     proto.Uint32(1),
		Opcode:             proto.Int32(Data.Opcode),//SceneListCount  
		SceneList:          []byte{byte(Data.Scene)},    
		VerifyContent:      proto.String(Data.VerifyContent),
		VerifyUserListSize: proto.Uint32(1),
		VerifyUserList:     VerifyUserList,
		ExtSpamInfo: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(WCExtInfoseq))),
			Buffer: WCExtInfoseq,
		},
		NeedConfirm: proto.Uint32(1),
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
		Cgiurl: "/cgi-bin/micromsg-bin/verifyuser",// /cgi-bin/micromsg-bin/verifyuser
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              137,
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
	Response := mm.VerifyUserResponse{}
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
