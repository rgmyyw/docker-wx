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

func SendVideoMsg(Data SendVideoMsgParam) models.ResponseResult {
	var err error
	var protobufdata []byte
	var errtype int64
	var imgbase64 []byte
	var playbase64 []byte

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

	ImgData := strings.Split(Data.ImageBase64, ",")
	PlayData := strings.Split(Data.Base64, ",")

	if len(ImgData) > 1 {
		imgbase64, _ = base64.StdEncoding.DecodeString(ImgData[1])
	} else {
		imgbase64, _ = base64.StdEncoding.DecodeString(Data.Base64)
	}

	if len(PlayData) > 1 {
		playbase64, _ = base64.StdEncoding.DecodeString(PlayData[1])
	} else {
		playbase64, _ = base64.StdEncoding.DecodeString(Data.Base64)
	}

	imgStream := bytes.NewBuffer(imgbase64)
	PlayStream := bytes.NewBuffer(playbase64)

	imgStartpos := 0
	imgdatalen := 50000
	imgdatatotalength := imgStream.Len()

	PlayStartpos := 0
	Playdatalen := 50000
	Playdatatotalength := PlayStream.Len()

	ClientImgId := fmt.Sprintf("%v_%v", Data.Wxid, time.Now().Unix())
	Time := time.Now().Unix()

	//上传缩略图
	I := 0
	for {
		imgStartpos = I * imgdatalen
		count := 0
		if imgdatatotalength-imgStartpos > imgdatalen {
			count = imgdatalen
		} else {
			count = imgdatatotalength - imgStartpos
		}

		if count < 0 {
			break
		}

		Databuff := make([]byte, count)
		_, _ = imgStream.Read(Databuff)
		req := &mm.UploadVideoReques{
			BaseRequest: &mm.BaseRequest{
				SessionKey:    D.Sessionkey,
				Uin:           proto.Uint32(D.Uin),
				DeviceId:      D.Deviceid_byte,
				ClientVersion: proto.Int32(int32(D.ClientVersion)),
				DeviceType:    []byte(D.DeviceType),
				Scene:         proto.Uint32(0),
			},
			ClientMsgId: proto.String(ClientImgId),
			ThumbData: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(Databuff))),
				Buffer: Databuff,
			},
			ThumbStartPos: proto.Uint32(uint32(imgStartpos)),
			ThumbTotalLen: proto.Uint32(uint32(imgdatatotalength)),
			FromUserName:  proto.String(Data.Wxid),
			ToUserName:    proto.String(Data.ToWxid),
			VideoTotalLen: proto.Uint32(uint32(Playdatatotalength)),
			VideoStartPos: proto.Uint32(0),
			VideoData: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(0),
				Buffer: []byte{},
			},
			PlayLength: proto.Uint32(Data.PlayLength),
			FuncFlag:   proto.Uint32(2),
			NetworkEnv: proto.Uint32(1),
			CameraType: proto.Uint32(2),
			EncryVer:   proto.Uint32(0),
			VideoFrom:  proto.Int(0),
			ReqTime:    proto.Int64(Time),
		}

		//序列化
		reqdata, _ := proto.Marshal(req)

		//发包
		protobufdata, _, errtype, err = comm.SendRequest(comm.SendPostData{
			Ip:     D.Mmtlsip,
			Host:   D.MmtlsHost,
			Cgiurl: "/cgi-bin/micromsg-bin/uploadvideo",///cgi-bin/micromsg-bin/uploadvideo
			Proxy:  D.Proxy,
			PackData: Algorithm.PackData{
				Reqdata:          reqdata,
				Cgi:              149,
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
	Response := mm.UploadVideoResponse{}
	err = proto.Unmarshal(protobufdata, &Response)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	if Response.GetBaseResponse().GetRet() != 0 {
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "缩略图上传失败",
			Data:    Response,
		}
	}

	//上传视频
	B := 0
	for {
		PlayStartpos = B * Playdatalen
		counts := 0
		if Playdatatotalength-PlayStartpos > Playdatalen {
			counts = Playdatalen
		} else {
			counts = Playdatatotalength - PlayStartpos
		}
		if counts < 0 {
			break
		}

		Databuffs := make([]byte, counts)
		_, _ = PlayStream.Read(Databuffs)
		reqA := &mm.UploadVideoReques{
			BaseRequest: &mm.BaseRequest{
				SessionKey:    D.Sessionkey,
				Uin:           proto.Uint32(D.Uin),
				DeviceId:      D.Deviceid_byte,
				ClientVersion: proto.Int32(int32(D.ClientVersion)),
				DeviceType:    []byte(D.DeviceType),
				Scene:         proto.Uint32(0),
			},
			ClientMsgId: proto.String(ClientImgId),
			ThumbData: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(0),
				Buffer: []byte{},
			},
			ThumbStartPos: proto.Uint32(uint32(imgdatatotalength)),
			ThumbTotalLen: proto.Uint32(uint32(imgdatatotalength)),
			FromUserName:  proto.String(Data.Wxid),
			ToUserName:    proto.String(Data.ToWxid),
			VideoTotalLen: proto.Uint32(uint32(Playdatatotalength)),
			VideoStartPos: proto.Uint32(uint32(PlayStartpos)),
			VideoData: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(Databuffs))),
				Buffer: Databuffs,
			},
			PlayLength: proto.Uint32(Data.PlayLength),
			FuncFlag:   proto.Uint32(2),
			NetworkEnv: proto.Uint32(1),
			CameraType: proto.Uint32(2),
			EncryVer:   proto.Uint32(0),
			VideoFrom:  proto.Int(0),
			ReqTime:    proto.Int64(Time),
		}

		//序列化
		reqdataA, _ := proto.Marshal(reqA)
		//fmt.Println(hex.EncodeToString(reqdataA))

		//发包
		protobufdata, _, errtype, err = comm.SendRequest(comm.SendPostData{
			Ip:     D.Mmtlsip,
			Host:   D.MmtlsHost,
			Cgiurl: "/cgi-bin/micromsg-bin/uploadvideo",
			Proxy:  D.Proxy,
			PackData: Algorithm.PackData{
				Reqdata:          reqdataA,
				Cgi:              149,
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

		B++
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
	Responses := mm.UploadVideoResponse{}
	err = proto.Unmarshal(protobufdata, &Responses)
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
