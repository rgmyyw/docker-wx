package Friend

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"

	"wechatdll/comm"
	"wechatdll/models"
)

func LbsFind(Data LbsFindParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	T := time.Now().Unix()
	WCExtInfoseq := Algorithm.GetIOSExtSpamInfoAndDeviceToken(D.Wxid, D.Deviceid_str, D.DeviceName, D.DeviceToken, T)

	req := &mm.LbsRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		OpCode:    proto.Uint32(Data.OpCode),
		Longitude: proto.Float32(Data.Longitude),
		Latitude:  proto.Float32(Data.Latitude),
		Precision: proto.Uint32(50),
		ExtSpamInfo: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(WCExtInfoseq))),
			Buffer: WCExtInfoseq,
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

	Host := D.MmtlsHost

	//发包
	protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
		Host:   Host,
		Cgiurl: "/cgi-bin/micromsg-bin/lbsfind",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              148,
			Uin:              D.Uin,
			Cookie:           D.Cooike,
			Sessionkey:       D.Sessionkey,
			EncryptType:      5,
			Loginecdhkey:     D.Loginecdhkey,
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
	Response := mm.LbsResponse{}
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
