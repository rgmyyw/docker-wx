package Tools

import (
	"fmt"
	"regexp"
	"strings"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"

	"github.com/golang/protobuf/proto"
)

func ThirdAppGrant(Data ThirdAppGrantParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	req := &mm.SdkOauthAuthorizeRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Appid:    proto.String(Data.Appid),
		Userinfo: proto.String("snsapi_message,snsapi_userinfo,snsapi_friend,snsapi_contact"),
		Tag4:     proto.String(""),
		Url:      proto.String(Data.Url),
		Tag8:     proto.String(""),
		Tag9:     proto.String(""),
		Tag10:    proto.String(""),
		Tag11:    proto.String(""),
		Tag12:    proto.Int32(0),
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
		Cgiurl: "/cgi-bin/mmbiz-bin/sdk_oauth_authorize",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              1388,
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

	code := string(protobufdata)
	str := GetInfoFromReg(code)
	if str != nil {
		return models.ResponseResult{
			Code:    0,
			Success: true,
			Message: "成功",
			Data:    str["code"],
		}
	}
	return models.ResponseResult{
		Code:    1,
		Success: false,
		Message: "失败",
		Data:    code,
	}
}
func GetMiddleString(SumString, LeftString, RightString string) string {
	if SumString == "" {
		return ""
	}
	LeftIndex := strings.Index(SumString, LeftString)
	if LeftIndex == -1 {
		return ""
	}
	//LeftIndex = LeftIndex + len(LeftString)
	RightIndex := strings.Index(SumString, RightString)
	if RightIndex == -1 {
		return ""
	}
	str := []rune(SumString)
	str1 := str[LeftIndex : LeftIndex+11]
	return string(str1)
}

func GetInfoFromReg(Xml string) map[string]string {
	re := regexp.MustCompile(`code=(?P<code>.*?)&state=`)
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
