package OfficialAccounts

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

func Follow(Data DefaultParam) models.ResponseResult {
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
		Value:               proto.String(Data.Appid),
		VerifyUserTicket:    proto.String(""),
		AntispamTicket:      proto.String(""),
		FriendFlag:          proto.Uint32(0),
		ChatRoomUserName:    proto.String(""),
		SourceUserName:      proto.String(""),
		SourceNickName:      proto.String(""),
		ScanQrcodeFromScene: proto.Uint32(0),
		ReportInfo:          proto.String(""),
		OuterUrl:            proto.String(""),
		SubScene:            proto.Int(0),
		BizReportInfo:       &mm.SKBuiltinBufferT{},
	})

	ccData := &mm.CryptoData{
		Version: []byte("00000003"),
		Type:    proto.Uint32(1),
		//EncryptData: wxCilent.GetiPadNewSpamData(D.Deviceid_str, D.DeviceName),
		Timestamp: proto.Uint32(uint32(time.Now().Unix())),
		Unknown5:  proto.Uint32(5),
		Unknown6:  proto.Uint32(0),
	}
	ccDataseq, _ := proto.Marshal(ccData)

	// 关注公众号使用01, 03, 还是06加密
	Wcstf := Algorithm.IpadWcstf(Data.Wxid)
	Wcste := Algorithm.IpadWcste(0, 0)

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
	}

	WCExtInfoseq, _ := proto.Marshal(WCExtInfo)

	req := &mm.VerifyUserRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Opcode:             proto.Int(1),
		VerifyUserListSize: proto.Uint32(1),
		VerifyUserList:     VerifyUserList,

		SceneListCount: proto.Uint32(1),
		SceneList:      []byte("W"),

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
		Cgiurl: "/cgi-bin/micromsg-bin/verifyuser",
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
