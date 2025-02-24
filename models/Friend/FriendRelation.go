package Friend

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strings"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/Cilent/wechat"
	"wechatdll/comm"
	"wechatdll/models"
)

type FriendRelationParam struct {
	Wxid     string
	UserName string
}

func GetBaseRequest(D *comm.LoginData) *mm.BaseRequest {
	ret := &mm.BaseRequest{}
	ret.SessionKey = []byte(D.Sessionkey)
	ret.Uin = &D.Uin
	if !strings.HasPrefix(D.Deviceid_str, "A") {
		ret.DeviceId = D.Deviceid_byte
		ret.ClientVersion = proto.Int32(int32(D.ClientVersion))
		ret.Scene = proto.Uint32(0)
		//log.Info("ios is base request")
	} else {
		ret.ClientVersion = proto.Int32(int32(D.ClientVersion))
		ret.DeviceId = D.Deviceid_byte
		ret.Scene = proto.Uint32(1)
	}
	return ret
}

// 好友关系检测
func FriendRelation(Data FriendRelationParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//构建请求  MMBizJsApiGetUserOpenIdRequest  SearchContactRequest
	req := &wechat.MMBizJsApiGetUserOpenIdRequest{
		BaseRequest: &wechat.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(1),
		},
		AppId:    proto.String("wx7c8d593b2c3a7703"),
		UserName: proto.String(Data.UserName),
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
		Cgiurl: "/cgi-bin/mmbiz-bin/usrmsg/mmbizjsapi_getuseropenid",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              1177,
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
	Response := wechat.MMBizJsApiGetUserOpenIdResponse{}
	err = proto.Unmarshal(protobufdata, &Response)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	//判断好友关系，1//删除 4/自己拉黑 5/被拉黑 0/正常
	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    Response,
	}
}
