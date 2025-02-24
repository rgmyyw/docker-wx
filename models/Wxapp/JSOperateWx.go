package Wxapp

import (
	"encoding/base64"
	"fmt"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"

	"github.com/golang/protobuf/proto"
)

func JSOperateWx(Data JSOperateWxParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	//var AndroidDeviceType = "android-34"
	//var AndroidManufacture = "HUAWEI Mate 60"
	//var AndroidModel = "BRA-AL00"
	//var AndroidRelease = "8"
	//var AndroidIncremental = "1"
	//DataBase64 := base64.StdEncoding.EncodeToString([]byte(Data.Data))
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	//sts, _ := base64.StdEncoding.DecodeString(DataBase64)
	if Data.Data == "" {
		Data.Data = "eyJ3aXRoX2NyZWRlbnRpYWxzIjp0cnVlLCJkYXRhIjp7ImxhbmciOiJlbiJ9LCJhcGlfbmFtZSI6IndlYmFwaV9nZXR1c2VyaW5mbyIsImZyb21fY29tcG9uZW50Ijp0cnVlfQ=="
		tmp, _ := base64.StdEncoding.DecodeString(Data.Data)
		Data.Data = string(tmp)
	}

	req := &mm.JSOperateWxDataRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Appid:       proto.String(Data.Appid),
		Data:        []byte(Data.Data),
		GrantScope:  proto.String(""),
		Opt:         proto.Int(Data.Opt),
		VersionType: proto.Int32(0),

		ExtInfo: &mm.WxaExternalInfo{
			HostAppid: proto.String(""),
			//Scene:     proto.Int32(1089),
			Scene: proto.Int32(0),
			//SourceEnv: proto.Int32(1),
			SourceEnv: proto.Int32(0),
		},

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
		Cgiurl: "/cgi-bin/mmbiz-bin/js-operatewxdata",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              1133,
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
	Response := mm.JSOperateWxDataResponse{}
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
