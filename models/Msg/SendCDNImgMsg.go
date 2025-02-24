package Msg

import (
	"encoding/xml"
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/Xml"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Report"
)

func SendCDNImgMsg(Data DefaultParam) models.ResponseResult {
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

	ClientImgId := fmt.Sprintf("%v", time.Now().Unix())

	//解析xml
	var Img Xml.ImgMsg
	xml.Unmarshal([]byte(Data.Content), &Img)

	TotalLen := Img.Img.Length

	req := &mm.UploadMsgImgRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		ClientImgId: &mm.SKBuiltinStringT{
			String_: proto.String(ClientImgId),
		},
		FromUserNam: &mm.SKBuiltinStringT{
			String_: proto.String(""),
		},
		ToUserNam: &mm.SKBuiltinStringT{
			String_: proto.String(Data.ToWxid),
		},
		TotalLen: proto.Uint32(uint32(TotalLen)),
		DataLen:  proto.Uint32(uint32(TotalLen)),
		StartPos: proto.Uint32(0),
		Data: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(0),
			Buffer: []byte{},
		},
		MsgType:           proto.Uint32(3),
		NetType:           proto.Uint32(0),
		CDNBigImgUrl:      proto.String(Img.Img.Cdnbigimgurl),
		CDNMidImgUrl:      proto.String(Img.Img.Cdnmidimgurl),
		AESKey:            proto.String(Img.Img.Aeskey),
		EncryVer:          proto.Int32(Img.Img.Encryver),
		CDNBigImgSize:     proto.Int32(Img.Img.Hdlength),
		CDNMidImgSize:     proto.Int32(TotalLen),
		CDNThumbImgUrl:    proto.String(Img.Img.Cdnthumburl),
		CDNThumbImgSize:   proto.Int32(Img.Img.Cdnthumblength),
		CDNThumbImgHeight: proto.Int32(Img.Img.Cdnthumbheight),
		CDNThumbImgWidth:  proto.Int32(Img.Img.Cdnthumbwidth),
		CDNThumbAESKey:    proto.String(Img.Img.Cdnthumbaeskey),
	}

	//序列化
	reqdata, _ := proto.Marshal(req)

	//发包
	protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/micromsg-bin/uploadmsgimg",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              110,
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
	NewSendMsgRespone := mm.UploadMsgImgResponse{}
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
