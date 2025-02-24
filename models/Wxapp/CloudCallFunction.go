package Wxapp

import (
	"fmt"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"

	"github.com/golang/protobuf/proto"
)

func CloudCallFunction(Data CloudCallParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	D.DeviceType = "Android-pad-34"
	//sts,_ := base64.StdEncoding.DecodeString("eyJhcGlfbmFtZSI6InFiYXNlX2NvbW1hcGkiLCJkYXRhIjp7InFiYXNlX2FwaV9uYW1lIjoidGNiYXBpX3Nsb3djYWxsZnVuY3Rpb25fdjIiLCJxYmFzZV9yZXEiOiJ7XFwiZnVuY3Rpb25fbmFtZVxcIjpcXCJjcnlwdG9cXCIsXFwiZGF0YVxcIjpcXCJ7XFxcXFxcImlcXFxcXFwiOlxcXFxcXCI3NzQwNlxcXFxcXCIsXFxcXFxcImNcXFxcXFwiOlxcXFxcXCI1MzAxMTFcXFxcXFwiLFxcXFxcXCJ1XFxcXFxcIjpcXFxcXFwidjE1ODk2MzMwNDQzNDcxXFxcXFxcIixcXFxcXFwib1xcXFxcXCI6XFxcXFxcIjQ5NDgyOFxcXFxcXCIsXFxcXFxcInRcXFxcXFwiOjE2MjQwNDI5NjJ9XFwiLFxcImFjdGlvblxcIjoxLFxcInNjZW5lXFwiOjEsXFwiY2FsbF9pZFxcIjpcXCIxNjI0MDM1OTIyMzk2LTAuMzc3MTIyODk4NjExMTEwOFxcIixcXCJjbG91ZGlkX2xpc3RcXCI6W119IiwicWJhc2Vfb3B0aW9ucyI6e30sInFiYXNlX21ldGEiOnsic2Vzc2lvbl9pZCI6IjE2MjQwMzU4MjkyMTgiLCJmaWx0ZXJfdXNlcl9pbmZvIjpmYWxzZSwic2RrX3ZlcnNpb24iOiJ3eC1taW5pcHJvZ3JhbS1zZGtcXC8yLjE3LjMgKDE2MjMzNzg0NzIwMDApIn0sImNsaV9yZXFfaWQiOiIxNjI0MDM1OTIyNDA4XzAuNzQzNjUyMzI4NDg2NzEyMyJ9LCJvcGVyYXRlX2RpcmVjdGx5IjpmYWxzZX0=")

	sts := Data.Data
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
		Data:        []byte(sts),
		GrantScope:  proto.String("scope.userInfo"),
		Opt:         proto.Int(1),
		VersionType: proto.Int32(0),
		ExtInfo: &mm.WxaExternalInfo{
			HostAppid: proto.String(""),
			//Scene:     proto.Int32(1089),
			Scene:     proto.Int32(1047),
			SourceEnv: proto.Int32(2),
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
