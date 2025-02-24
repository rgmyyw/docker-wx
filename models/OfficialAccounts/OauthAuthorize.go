package OfficialAccounts

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/bts"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Tools"
)

func OauthAuthorize(Data GetkeyParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	A8key := Tools.GetA8Key(Tools.GetA8KeyParam{
		Wxid:        Data.Wxid,
		OpCode:      2,
		Scene:       4,
		CodeType:    19,
		CodeVersion: 10,
		ReqUrl:      Data.Url,
	})

	GetA8key := bts.GetA8KeyResponse(A8key.Data)

	//组url
	url := `https://open.weixin.qq.com/connect/oauth2/authorize?appid=`+Data.Appid+`&redirect_uri=`+GetA8key.GetFullURL()+`&response_type=code&scope=snsapi_base&state=123#wechat_redirect`

	req := &mm.OauthAuthorizeReq{
		BaseRequest:          &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		OauthUrl:             proto.String(url),
		BizUsername:          proto.String(""),
		Scene:                proto.Uint32(4),
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
		Cgiurl: "/cgi-bin/mmbiz-bin/oauth_authorize",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              1254,
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
	Response := mm.OauthAuthorizeResp{}
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
