package Group

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

type FacingCreateChatRoomParam struct {
	Wxid    	string
	Password    string
	OpCode      uint32
	Latitude 	float32
	Longitude	float32
}

func RandomLongitude() float32 {
	seed := rand.Intn(1999)
	longitude := 113.8 + float32(seed) / 10000
	return longitude
}

func RandomLatitude() float32 {
	seed := rand.Intn(1999)
	latitude := 27.5 + float32(seed) / 10000
	return latitude
}

func FacingCreateChatRoom(Data FacingCreateChatRoomParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	req := &mm.FacingCreateChatRoomRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		OpCode: proto.Uint32(Data.OpCode),
		ChatroomPassword: proto.String(Data.Password),
		Longitude: proto.Float32(Data.Longitude),
		Latitude: proto.Float32(Data.Latitude),
		Tag6: proto.Uint32(15),
		Tag9: proto.Uint32(1),
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
		Cgiurl: "/cgi-bin/micromsg-bin/mmfacingcreatechatroom",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              653,
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
	Response := mm.FacingCreateChatRoomResponse{}
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
