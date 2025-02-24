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

type SendImageMsgParam struct {
	Wxid   string
	ToWxid string
	Base64 string
}

func SendImageMsg(Data SendImageMsgParam) models.ResponseResult {
	var err error
	var protobufdata []byte
	var errtype int64
	var imgbase64 []byte

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

	ImgData := strings.Split(Data.Base64, ",")

	if len(ImgData) > 1 {
		imgbase64, _ = base64.StdEncoding.DecodeString(ImgData[1])
	} else {
		imgbase64, _ = base64.StdEncoding.DecodeString(Data.Base64)
	}

	imgStream := bytes.NewBuffer(imgbase64)

	Startpos := 0
	datalen := 50000
	datatotalength := imgStream.Len()

	ClientImgId := fmt.Sprintf("%v_%v", Data.Wxid, time.Now().Unix())

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
		_, _ = imgStream.Read(Databuff)

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
				String_: proto.String(Data.Wxid),
			},
			ToUserNam: &mm.SKBuiltinStringT{
				String_: proto.String(Data.ToWxid),
			},
			TotalLen: proto.Uint32(uint32(datatotalength)),
			StartPos: proto.Uint32(uint32(Startpos)),
			DataLen:  proto.Uint32(uint32(len(Databuff))),
			Data: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(Databuff))),
				Buffer: Databuff,
			},
			MsgType:    proto.Uint32(3),
			EncryVer:   proto.Int32(0),
			ReqTime:    proto.Uint32(uint32(time.Now().Unix())),
			MessageExt: proto.String("png"),
		}

		//序列化
		reqdata, _ := proto.Marshal(req)

		//发包
		protobufdata, _, errtype, err = comm.SendRequest(comm.SendPostData{
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
				UseCompress:      true,
			},
		}, D.MmtlsKey)

		if err != nil {
			break
		}

		I++
	}

	if err != nil {
		return models.ResponseResult{
			Code:    errtype,
			Success: false,
			Message: err.Error(),
			Data:    nil,
		}
	}

	//解包
	Response := mm.UploadMsgImgResponse{}
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
