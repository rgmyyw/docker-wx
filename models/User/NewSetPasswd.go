package User

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/lib"
	"wechatdll/models"
)

func NewSetPasswd(Data NewSetPasswdParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	req := &mm.SetPwdRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Password: proto.String(lib.MD5ToLower(Data.NewPassword)),
		Ticket:   proto.String(Data.Ticket),
		AutoAuthKey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(D.Autoauthkey))),
			Buffer: D.Autoauthkey,
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

	S2801, _ := hex.DecodeString("2801")

	reqdataA := new(bytes.Buffer)
	reqdataA.Write(reqdata)

	if len(D.Deviceid_byte) <= 16 {
		reqdataA.Write(S2801)
	}

	//发包
	protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
		Ip:     D.Mmtlsip,
		Host:   D.MmtlsHost,
		Cgiurl: "/cgi-bin/micromsg-bin/newsetpasswd",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdataA.Bytes(),
			Cgi:              383,
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
	Response := mm.SetPwdResponse{}
	err = proto.Unmarshal(protobufdata, &Response)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	//更新成功就保存autoAuthKey
	if Response.AutoAuthKey != nil && len(Response.AutoAuthKey.Buffer) > 10 {
		D.Autoauthkey = Response.AutoAuthKey.Buffer
		err = comm.CreateLoginData(*D, D.Wxid, 0)
		if err != nil {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("AutoAuthKey保存失败：%v", err.Error()),
				Data:    nil,
			}
		}
	}

	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    Response,
	}

}
