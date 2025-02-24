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

func SendCDNVideoMsg(Data DefaultParam) models.ResponseResult {
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

	ClientImgId := fmt.Sprintf("%v_%v", Data.Wxid, time.Now().Unix())

	//解析xml
	var Video Xml.VideoMsg
	xml.Unmarshal([]byte(Data.Content), &Video)

	req := &mm.UploadVideoReques{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		ClientMsgId:   proto.String(ClientImgId),
		FromUserName:  proto.String(Data.Wxid),
		ToUserName:    proto.String(Data.ToWxid),
		ThumbTotalLen: proto.Uint32(Video.Video.Cdnthumblength),
		ThumbStartPos: proto.Uint32(Video.Video.Cdnthumblength),
		ThumbData: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(0),
			Buffer: []byte{},
		},
		VideoTotalLen: proto.Uint32(Video.Video.Length),
		VideoStartPos: proto.Uint32(Video.Video.Length),
		VideoData: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(0),
			Buffer: []byte{},
		},
		PlayLength:        proto.Uint32(Video.Video.Playlength),
		NetworkEnv:        proto.Uint32(1),
		CameraType:        proto.Uint32(2),
		FuncFlag:          proto.Uint32(2),
		MsgSource:         proto.String(""),
		CdnvideoUrl:       proto.String(Video.Video.Cdnvideourl),
		Aeskey:            proto.String(Video.Video.Aeskey),
		EncryVer:          proto.Uint32(1),
		CdnthumbUrl:       proto.String(Video.Video.Cdnthumburl),
		CdnthumbImgSize:   proto.Uint32(Video.Video.Cdnthumblength),
		CdnthumbImgHeight: proto.Uint32(Video.Video.Cdnthumbheight),
		CdnthumbImgWidth:  proto.Uint32(Video.Video.Cdnthumbwidth),
		CdnthumbAeskey:    proto.String(Video.Video.Aeskey),
		VideoMd5:          proto.String(Video.Video.Md5),
		VideoNewMd5:       proto.String(Video.Video.Newmd5),
		MsgForwardType:    proto.Int64(2),
		Source:            proto.Int64(3),
	}

	//序列化
	reqdata, _ := proto.Marshal(req)

	//发包
	protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/micromsg-bin/uploadvideo",
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
	NewSendMsgRespone := mm.UploadVideoResponse{}
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
