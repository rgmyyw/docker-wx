package Msg

import (
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/golang/protobuf/proto"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

type VideoMsg struct {
	VideoMsg xml.Name `xml:"videomsg"`
	Aeskey string    `xml:"aeskey,attr"`
	Cdnthumbaeskey string `xml:"cdnthumbaeskey,attr"`
	Cdnvideourl string `xml:"cdnvideourl,attr"`
	Cdnthumburl string `xml:"cdnthumburl,attr"`
	Length string `xml:"length,attr"`
	Playlength string `xml:"playlength,attr"`
	Cdnthumblength string `xml:"cdnthumblength,attr"`
	Cdnthumbwidth string `xml:"cdnthumbwidth,attr"`
	Cdnthumbheight string `xml:"cdnthumbheight,attr"`
	Fromusername string `xml:"fromusername,attr"`
	Md5 string `xml:"md5,attr"`
	Newmd5 string `xml:"newmd5,attr"`
}


func ShareVideoMsg(Data ShareVideoMsgParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	str := strings.Replace(Data.Xml,"<?xml version=\"1.0\"?>","",-1)

	str = strings.Replace(Data.Xml,"</msg>","",-1)
	str = strings.Replace(Data.Xml,"<msg>","",-1)
	//fmt.Println(str)
	ret,err := GetXml(str)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("Xml匹配失败"),
			Data:    nil,
		}
	}

	//Groups := GetInfoFromReg(Data.Xml)
	//if Groups == nil {
	//	return models.ResponseResult{
	//		Code:    -8,
	//		Success: false,
	//		Message: fmt.Sprintf("Xml匹配失败"),
	//		Data:    nil,
	//	}
	//}

	ToUserName := Data.ToWxid
	//ThumbTotalLen1,_ := strconv.Atoi(Groups["cdnthumblength"])
	//ThumbTotalLen := uint32(ThumbTotalLen1)
	//ThumbStartPos1,_ := strconv.Atoi(Groups["cdnthumblength"])
	//ThumbStartPos := uint32(ThumbStartPos1)
	//VideoTotalLen1,_ := strconv.Atoi(Groups["length"])
	//VideoTotalLen := uint32(VideoTotalLen1)
	//VideoStartPos2,_ := strconv.Atoi(Groups["length"])
	//VideoStartPos := uint32(VideoStartPos2)
	//PlayLength1,_ := strconv.Atoi(Groups["playlength"])
	//PlayLength := uint32(PlayLength1)
	//AESKey := Groups["aeskey"]
	//CDNVideoUrl := Groups["cdnvideourl"]

	ThumbTotalLen1,_ := strconv.Atoi(ret.Cdnthumblength)
	ThumbTotalLen := uint32(ThumbTotalLen1)
	ThumbStartPos1,_ := strconv.Atoi(ret.Cdnthumblength)
	ThumbStartPos := uint32(ThumbStartPos1)
	VideoTotalLen1,_ := strconv.Atoi(ret.Length)
	VideoTotalLen := uint32(VideoTotalLen1)
	VideoStartPos2,_ := strconv.Atoi(ret.Length)
	VideoStartPos := uint32(VideoStartPos2)
	PlayLength1,_ := strconv.Atoi(ret.Playlength)
	PlayLength := uint32(PlayLength1)
	AESKey := ret.Aeskey
	CDNVideoUrl := ret.Cdnvideourl

	ClientImgId := fmt.Sprintf("%v",time.Now().Unix())

	req := &mm.UploadVideoReques{
		BaseRequest:          &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		ClientMsgId:          proto.String(ClientImgId),
		FromUserName:         proto.String(""),
		ToUserName:           proto.String(ToUserName),
		ThumbTotalLen:        proto.Uint32(ThumbTotalLen),
		ThumbStartPos:        proto.Uint32(ThumbStartPos),
		ThumbData:            &mm.SKBuiltinBufferT{},
		VideoTotalLen:        proto.Uint32(VideoTotalLen),
		VideoStartPos:        proto.Uint32(VideoStartPos),
		VideoData:            &mm.SKBuiltinBufferT{},
		PlayLength:           proto.Uint32(PlayLength),
		NetworkEnv:           proto.Uint32(1),
		CameraType:           proto.Uint32(2),
		FuncFlag:             proto.Uint32(2),
		MsgSource:            nil,
		CdnvideoUrl:          proto.String(CDNVideoUrl),
		Aeskey:               proto.String(AESKey),
		EncryVer:             proto.Uint32(1),
		CdnthumbUrl:          proto.String(CDNVideoUrl),
		CdnthumbImgSize:      proto.Uint32(ThumbTotalLen),
		CdnthumbImgHeight:    proto.Uint32(960),
		CdnthumbImgWidth:     proto.Uint32(540),
		CdnthumbAeskey:       proto.String(AESKey),
		VideoFrom:            proto.Int(0),
		ReqTime:              proto.Int64(0),
		//VideoMd5:             nil,
		//StreamVideoUrl:       nil,
		//StreamVideoTotalTime: nil,
		//StreamVideoTitle:     nil,
		//StreamVideoWording:   nil,
		//StreamVideoWebUrl:    nil,
		//StreamVideoThumbUrl:  nil,
		//StreamVideoPublishId: nil,
		//StreamVideoAdUxInfo:  nil,
		//StatExtStr:           nil,
		//HitMd5:               nil,
		//VideoNewMd5:          nil,
		//Crc32:                nil,
		//MsgForwardType:       nil,
		//Source:               nil,
		//SendMsgTicket:        nil,
		//XXX_NoUnkeyedLiteral: struct{}{},
		//XXX_unrecognized:     nil,
		//XXX_sizecache:        0,
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

	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    Response,
	}
}


func GetInfoFromReg(Xml string) map[string]string {
	re := regexp.MustCompile(`<\?xml version="1.0"\?><msg><videomsg aeskey="(?P<aeskey>.*?)" cdnthumbaeskey="(?P<cdnthumbaeskey>.*?)" cdnvideourl="(?P<cdnvideourl>.*?)" cdnthumburl="(?P<cdnthumburl>.*?)" length="(?P<length>\d+)" playlength="(?P<playlength>\d+)" cdnthumblength="(?P<cdnthumblength>\d+)" cdnthumbwidth="(?P<cdnthumbwidth>\d+)" cdnthumbheight="(?P<cdnthumbheight>\d+)" fromusername="wxid_2555175551614" md5="c03f56866fd5698c54375a9874379c90" newmd5="fd4f62c4ca7b339451aeee45a6cbd3c2" isad="0" /></msg>`)
	res := re.MatchString(Xml)
	if res == false {
		return nil
	}
	infos := re.FindAllStringSubmatch(Xml, -1)
	results := re.SubexpNames()
	for _, info := range infos {
		m := make(map[string]string)
		for j, name := range results {
			if j != 0 && name != "" {
				m[name] = strings.TrimSpace(info[j])
			}
		}
		return m
	}
	return nil
}


func GetXml(Xml string) (VideoMsg, error) {
	v := VideoMsg{}
	err := xml.Unmarshal([]byte(Xml), &v)
	if err != nil {
		logs.Error("Unmarshal Failed: %v", err)
		return v,err
	}
	return v,nil
}