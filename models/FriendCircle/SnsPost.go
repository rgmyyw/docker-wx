package FriendCircle

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strings"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

type Messagearameter struct {
	Wxid         string
	Content      string
	BlackList    string
	WithUserList string
}

func Messages(Data Messagearameter) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	B := strings.Split(Data.BlackList, ",")
	BS := make([]*mm.SKBuiltinStringT, len(B))

	if len(B) >= 1 {
		for i, v := range B {
			BS[i] = &mm.SKBuiltinStringT{
				String_: proto.String(v),
			}
		}
	}

	W := strings.Split(Data.WithUserList, ",")
	WS := make([]*mm.SKBuiltinStringT, len(W))

	if len(W) >= 1 {
		for i, v := range W {
			WS[i] = &mm.SKBuiltinStringT{
				String_: proto.String(v),
			}
		}
	}

	ccData := &mm.CryptoData{
		Version: []byte("00000003"),
		Type:    proto.Uint32(1),
		EncryptData: Algorithm.GetiPadNewSpamData(D.Deviceid_str, D.DeviceName, D.DeviceToken),
		Timestamp: proto.Uint32(uint32(time.Now().Unix())),
		Unknown5:  proto.Uint32(5),
		Unknown6:  proto.Uint32(0),
	}
	ccDataseq, _ := proto.Marshal(ccData)

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

	// TODO: 查看朋友圈是使用01, 03, 还是06加密
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
		DeviceToken: &mm.SKBuiltinBufferT{
			ILen:                 proto.Uint32(uint32(len(DeviceTokenCCDPB))),
			Buffer:               DeviceTokenCCDPB,
		},
	}
	_, _ = proto.Marshal(WCExtInfo)

	req := &mm.SnsPostRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		ObjectDesc: &mm.SKBuiltinString_S{
			ILen:   proto.Uint32(uint32(len(Data.Content))),
			Buffer: proto.String(Data.Content),
		},
		WithUserListNum: proto.Uint32(uint32(len(W))),
		WithUserList:    WS,
		ClientId:        proto.String(fmt.Sprintf("sns_post_%v_%v_0", D.Wxid, time.Now().Unix())),
		BlackListNum:    proto.Uint32(uint32(len(B))),
		BlackList:       BS,
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
		Cgiurl: "/cgi-bin/micromsg-bin/mmsnspost",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              209,
			Uin:              D.Uin,
			Cookie:           D.Cooike,
			Sessionkey:       D.Sessionkey,
			EncryptType:      5,
			Loginecdhkey:     D.RsaPublicKey,
			Clientsessionkey: D.Clientsessionkey,
			UseCompress:      true,
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
	SnsPostResponse := mm.SnsPostResponse{}
	err = proto.Unmarshal(protobufdata, &SnsPostResponse)

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
		Data:    SnsPostResponse,
	}

}
