package Friend

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strings"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

type UploadParam struct {
	Wxid           string
	PhoneNo        string
	CurrentPhoneNo string
	Opcode         int32
}

func Upload(Data UploadParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	if Data.PhoneNo == "" {
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "PhoneNo 手机号是必须的",
			Data:    nil,
		}
	}

	PhoneNoSplit := strings.Split(Data.PhoneNo, ",")

	var PhoneNoList []*mm.SKBuiltinStringT

	for _, v := range PhoneNoSplit {
		PhoneNoList = append(PhoneNoList, &mm.SKBuiltinStringT{
			String_: proto.String(v),
		})
	}

	req := &mm.UploadMContactRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		UserName:       proto.String(Data.Wxid),
		Opcode:         proto.Int32(Data.Opcode),
		Mobile:         proto.String(Data.CurrentPhoneNo),
		MobileListSize: proto.Int32(int32(len(PhoneNoList))),
		MobileList:     PhoneNoList,
		EmailListSize:  proto.Int32(0),
		EmailList:      nil,
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
		Cgiurl: "/cgi-bin/micromsg-bin/uploadmcontact",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              133,
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
	Response := mm.UploadMContactResponse{}
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
