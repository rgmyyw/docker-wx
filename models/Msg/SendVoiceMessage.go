package Msg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strings"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Report"
)

type SendVoiceMessageParam struct {
	Wxid      string
	ToWxid    string
	Base64    string
	VoiceTime int32
	Type      int32
}

func SendVoiceMessage(Data SendVoiceMessageParam) models.ResponseResult {
	var err error
	var protobufdata []byte
	var errtype int64

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

	VoiceData := strings.Split(Data.Base64, ",")

	var VoiceBase64 []byte

	if len(VoiceData) > 1 {
		VoiceBase64, _ = base64.StdEncoding.DecodeString(VoiceData[1])
	} else {
		VoiceBase64, _ = base64.StdEncoding.DecodeString(Data.Base64)
	}

	// 原始版本
	//VoiceStream := bytes.NewBuffer(VoiceBase64)
	//
	//Startpos := 0
	//datalen := 65000
	//datatotalength := VoiceStream.Len()
	//
	//ClientImgId := fmt.Sprintf("%v_%v", Data.Wxid, time.Now().Unix())
	//
	//I := 0
	//
	//for {
	//	Startpos = I * datalen
	//	count := 0
	//	if datatotalength-Startpos > datalen {
	//		count = datalen
	//	} else {
	//		count = datatotalength - Startpos
	//	}
	//	if count < 0 {
	//		break
	//	}
	//
	//	Databuff := make([]byte, count)
	//	_, _ = VoiceStream.Read(Databuff)
	//
	//	req := &mm.UploadVoiceRequest{
	//		FromUserName: proto.String(Data.Wxid),
	//		ToUserName:   proto.String(Data.ToWxid),
	//		Offset:       proto.Uint32(uint32(Startpos)),
	//		Length:       proto.Int32(int32(datatotalength)),
	//		ClientMsgId:  proto.String(ClientImgId),
	//		MsgId:        proto.Uint32(0),
	//		VoiceLength:  proto.Int32(Data.VoiceTime),
	//		Data: &mm.SKBuiltinBufferT{
	//			ILen:   proto.Uint32(uint32(len(Databuff))),
	//			Buffer: Databuff,
	//		},
	//		EndFlag: proto.Uint32(1),
	//		BaseRequest: &mm.BaseRequest{
	//			SessionKey:    D.Sessionkey,
	//			Uin:           proto.Uint32(D.Uin),
	//			DeviceId:      D.Deviceid_byte,
	//			ClientVersion: proto.Int32(int32(D.ClientVersion)),
	//			DeviceType:    []byte(D.DeviceType),
	//			Scene:         proto.Uint32(0),
	//		},
	//		CancelFlag:  proto.Uint32(0),
	//		Msgsource:   proto.String(""),
	//		VoiceFormat: proto.Int32(Data.Type),
	//		ForwardFlag: proto.Uint32(0),
	//		NewMsgId:    proto.Uint64(0),
	//		Offst:       proto.Uint32(0),
	//	}
	//
	//	//序列化
	//	reqdata, _ := proto.Marshal(req)
	//
	//	//发包
	//	protobufdata, _, errtype, err = comm.SendRequest(comm.SendPostData{
	//		Ip:     D.Mmtlsip,
	//		Cgiurl: "/cgi-bin/micromsg-bin/uploadvoice",
	//		Proxy:  D.Proxy,
	//		PackData: Algorithm.PackData{
	//			Reqdata:          reqdata,
	//			Cgi:              127,
	//			Uin:              D.Uin,
	//			Cookie:           D.Cooike,
	//			Sessionkey:       D.Sessionkey,
	//			EncryptType:      5,
	//			Loginecdhkey:     D.Loginecdhkey,
	//			Clientsessionkey: D.Clientsessionkey,
	//			UseCompress:      true,
	//		},
	//	}, D.MmtlsKey)
	//
	//	if err != nil {
	//		break
	//	}
	//
	//	I++
	//}


	// HT版本循环版
	VoiceStream := bytes.NewBuffer(VoiceBase64)

	Startpos := 0
	datalen := 65000
	datatotalength := VoiceStream.Len()

	I := 0
	for {
		Startpos = I * datalen
		count := 0
		if datatotalength-Startpos > datalen {
			count = datalen
		} else {
			count = datatotalength - Startpos
		}
		if count < 0 {
			break
		}

		Databuff := make([]byte, count)
		n, err := VoiceStream.Read(Databuff)
		fmt.Println("=======n======:",n)

		buffer := &mm.SKBuiltinBufferT{
			ILen:                 proto.Uint32(uint32(len(Databuff))),
			Buffer:               Databuff,
		}
		ClientImgId := fmt.Sprintf("%v", time.Now().Unix())
		req := &mm.UploadVoiceRequest{
			FromUserName:         proto.String(Data.Wxid),
			ToUserName:           proto.String(Data.ToWxid),
			Offset:               proto.Uint32(uint32(Startpos)),
			Length:               proto.Int32(int32(len(Databuff))),
			//Length:               proto.Int32(int32(datatotalength)),
			ClientMsgId:          proto.String(ClientImgId),
			MsgId:                proto.Uint32(0),
			VoiceLength:          proto.Int32(Data.VoiceTime),
			Data:                 buffer,
			EndFlag:              proto.Uint32(1),
			BaseRequest:          &mm.BaseRequest{
				SessionKey:    D.Sessionkey,
				Uin:           proto.Uint32(D.Uin),
				DeviceId:      D.Deviceid_byte,
				ClientVersion: proto.Int32(int32(D.ClientVersion)),
				DeviceType:    []byte(D.DeviceType),
				Scene:         proto.Uint32(0),
			},
			CancelFlag:  proto.Uint32(0),
			Msgsource:            proto.String(""),
			VoiceFormat:          proto.Int32(Data.Type),
			//UicreateTime:         nil,
			ForwardFlag: proto.Uint32(0),
			NewMsgId:    proto.Uint64(0),
			//ReqTime:              nil,
			//VoiceId:              nil,
			Offst:                proto.Uint32(0),
		}

		//序列化
		reqdata, _ := proto.Marshal(req)

		//发包
		protobufdata, _, errtype, err = comm.SendRequest(comm.SendPostData{
			Ip:     D.Mmtlsip,
			Host:   D.MmtlsHost,
			Cgiurl: "/cgi-bin/micromsg-bin/uploadvoice",
			Proxy:  D.Proxy,
			PackData: Algorithm.PackData{
				Reqdata:          reqdata,
				Cgi:              127,
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
			break
		}
		I++
	}




	// HT版本不循环版（不能超过65000）
	//VoiceStream := bytes.NewBuffer(VoiceBase64)
	//datatotalength := VoiceStream.Len()
	//Databuff := make([]byte, datatotalength)
	//n, err := VoiceStream.Read(Databuff)
	//fmt.Println("=======n======:",n)
	//
	//buffer := &mm.SKBuiltinBufferT{
	//	ILen:                 proto.Uint32(uint32(len(Databuff))),
	//	Buffer:               Databuff,
	//}
	//ClientImgId := fmt.Sprintf("%v", time.Now().Unix())
	//req := &mm.UploadVoiceRequest{
	//	FromUserName:         proto.String(Data.Wxid),
	//	ToUserName:           proto.String(Data.ToWxid),
	//	Offset:               proto.Uint32(0),
	//	Length:               proto.Int32(int32(len(Databuff))),
	//	ClientMsgId:          proto.String(ClientImgId),
	//	MsgId:                proto.Uint32(0),
	//	VoiceLength:          proto.Int32(Data.VoiceTime),
	//	Data:                 buffer,
	//	EndFlag:              proto.Uint32(1),
	//	BaseRequest:          &mm.BaseRequest{
	//		SessionKey:    D.Sessionkey,
	//		Uin:           proto.Uint32(D.Uin),
	//		DeviceId:      D.Deviceid_byte,
	//		ClientVersion: proto.Int32(int32(D.ClientVersion)),
	//		DeviceType:    []byte(D.DeviceType),
	//		Scene:         proto.Uint32(0),
	//	},
	//	CancelFlag:  proto.Uint32(0),
	//	Msgsource:            proto.String(""),
	//	VoiceFormat:          proto.Int32(Data.Type),
	//	//UicreateTime:         nil,
	//	ForwardFlag: proto.Uint32(0),
	//	NewMsgId:    proto.Uint64(0),
	//	//ReqTime:              nil,
	//	//VoiceId:              nil,
	//	Offst:                proto.Uint32(0),
	//}
	//
	////序列化
	//reqdata, _ := proto.Marshal(req)
	//
	////发包
	//protobufdata, _, errtype, err = comm.SendRequest(comm.SendPostData{
	//	Ip:     D.Mmtlsip,
	//	Cgiurl: "/cgi-bin/micromsg-bin/uploadvoice",
	//	Proxy:  D.Proxy,
	//	PackData: Algorithm.PackData{
	//		Reqdata:          reqdata,
	//		Cgi:              127,
	//		Uin:              D.Uin,
	//		Cookie:           D.Cooike,
	//		Sessionkey:       D.Sessionkey,
	//		EncryptType:      5,
	//		Loginecdhkey:     D.Loginecdhkey,
	//		Clientsessionkey: D.Clientsessionkey,
	//		UseCompress:      true,
	//	},
	//}, D.MmtlsKey)
	//
	//fmt.Println("111111",protobufdata)


	// 到这
	if err != nil {
		return models.ResponseResult{
			Code:    errtype,
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}
	}

	//解包
	Response := mm.UploadVoiceResponse{}
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
