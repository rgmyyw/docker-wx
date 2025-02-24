package User

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

func BindQQ(Data BindQQParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	qqReq := Algorithm.QQmobileRequest{
		ClientVersion:       uint32(D.ClientVersion),
		Account:             Data.Account,
		Password:            Data.Password,
		DeviceIdStr:         D.Deviceid_str,
		Ksid:                []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		TimeSpan:            uint32(time.Now().Unix()),
		MobileSystem:        "ios",
		MobileSystemVersion: "13.3",
		ISPName:             "中国移动",
		ISPType:             "wifi",
		MobileBrand:         "iPhone",
		WechatVersion:       "8.0.6",
	}

	qqLoginBuff := Algorithm.QQPackMessage(qqReq)

	req := &mm.BindQQRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Qq: proto.Uint32(Data.Account),
		Password: proto.String(""),
		Password2: proto.String(""),
		ImageSid: proto.String(""),
		ImageCode: proto.String(""),
		Opcode: proto.Uint32(1),
		ImageEncKey: &mm.SKBuiltinStringT{
			String_: proto.String(""),
		},
		Ksid: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(0),
			Buffer: []byte{},
		},
		SetAsMain: proto.Uint32(0),
		SafeDeviceName: proto.String("iPad"),
		SafeDeviceType: proto.String("iPad iOS13.3"),
		WtLoginReqBuff: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(qqLoginBuff))),
			Buffer: qqLoginBuff,
		},
	}

	reqdata, err := proto.Marshal(req)

	fmt.Printf("BindQQ请求: %x\n", reqdata)
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
		Cgiurl: "/cgi-bin/micromsg-bin/bindqq",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              144,
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
	Response := mm.BindQQResponse{}
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
