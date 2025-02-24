package Tools

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

type GetA8KeyParam struct {
	Wxid         string
	OpCode       uint32
	Scene        uint32
	CodeType     uint32
	CodeVersion  uint32
	ReqUrl       string
	CookieBase64 string
	NetType      string
	Flag		 uint32
}

func GetA8Key(Data GetA8KeyParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	cookieBytes, err := base64.StdEncoding.DecodeString(Data.CookieBase64)
	if err != nil {
		cookieBytes = []byte{}
	}

	req := &mm.GetA8KeyReq{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		CodeType:    proto.Uint32(Data.CodeType),
		CodeVersion: proto.Uint32(Data.CodeVersion),
		Flag:        proto.Uint32(0),
		FontScale:   proto.Uint32(100),
		NetType:     proto.String(Data.NetType),
		OpCode:      proto.Uint32(Data.OpCode),
		UserName:    proto.String(D.Wxid),
		ReqUrl: &mm.SKBuiltinStringT{
			String_: proto.String(Data.ReqUrl),
		},
		FriendQq: proto.Uint32(0),
		Scene:    proto.Uint32(Data.Scene),
		Cookie: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(cookieBytes))),
			Buffer: cookieBytes,
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
		Cgiurl: "/cgi-bin/micromsg-bin/geta8key",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              233,
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
	Response := mm.GetA8KeyResp{}
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

func GetA8Key3rd(Data GetA8KeyParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	cookieBytes, err := base64.StdEncoding.DecodeString(Data.CookieBase64)
	if err != nil {
		cookieBytes = []byte{}
	}

	req := &mm.GetA8KeyReq{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		CodeType:    proto.Uint32(Data.CodeType),
		CodeVersion: proto.Uint32(Data.CodeVersion),
		Flag:        proto.Uint32(Data.Flag),
		FontScale:   proto.Uint32(118),
		NetType:     proto.String(Data.NetType),
		OpCode:      proto.Uint32(Data.OpCode),
		UserName:    proto.String(D.Wxid),
		ReqUrl: &mm.SKBuiltinStringT{
			String_: proto.String(Data.ReqUrl),
		},
		FriendQq: 	 proto.Uint32(0),
		Scene:    	 proto.Uint32(Data.Scene),
		SubScene: 	 proto.Uint32(1),
		RequestId: 	 proto.Uint64(3304789257),
		Cookie: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(cookieBytes))),
			Buffer: cookieBytes,
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
		Cgiurl: "/cgi-bin/micromsg-bin/3rd-geta8key",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              226,
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
	Response := mm.GetA8KeyResp{}
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
