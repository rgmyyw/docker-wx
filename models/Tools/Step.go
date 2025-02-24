package Tools

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/Mmtls"
	"wechatdll/baseinfo"
	"wechatdll/comm"
	"wechatdll/models"
)



type SetStepParam struct {
	Wxid   string
	Number uint64
}

// 上传步数
func UploadStepSetRequestRequest(userInfo *comm.LoginData, deviceID string, deviceType string, number uint64) (*PackHeader, error) {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix()
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).Unix()
	var req = &mm.UploadDeviceStepRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    []byte{},
			Uin:           proto.Uint32(userInfo.Uin),
			DeviceId:      userInfo.Deviceid_byte,
			ClientVersion: proto.Int32(int32(userInfo.ClientVersion)),
			DeviceType:    []byte(userInfo.DeviceType),
			Scene:         proto.Uint32(0),
		},
		DeviceID:   proto.String(deviceID),
		DeviceType: proto.String(deviceType),
		FromTime:   proto.Uint32(uint32(startTime)),
		ToTime:     proto.Uint32(uint32(endTime)),
		StepCount:  proto.Uint32(uint32(number)),
	}
	// 打包发送数据
	srcData, _ := proto.Marshal(req)
	hecData := Pack(userInfo, srcData, 1261, 5)
	httpclient := Mmtls.GenNewHttpClient(userInfo.MmtlsKey, userInfo.MmtlsHost)
	resp, err := httpclient.MMtlsPost(userInfo.MmtlsHost, "/cgi-bin/mmoc-bin/hardware/uploaddevicestep", hecData, userInfo.Proxy)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(resp, nil)
}

// 获取设备
func GetBoundHardDeviceRequest(D *comm.LoginData) (*PackHeader, error) {

	// 获取设备信息
	req := &mm.GetBoundHardDevicesRequest835{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    []byte{},
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Version: proto.Uint32(0),
	}

	reqdata, err := proto.Marshal(req)

	if err != nil {
		return nil, err
	}

	hecData := Pack(D, reqdata, baseinfo.MMRequestTypeupdateStep, 5)
	httpclient := Mmtls.GenNewHttpClient(D.MmtlsKey, D.MmtlsHost)
	recvData, err := httpclient.MMtlsPost(D.MmtlsHost, "/cgi-bin/micromsg-bin/getboundharddevices", hecData, D.Proxy)
	if err != nil {
		return nil, err
	}
	return DecodePackHeader(recvData, nil)
}

// 修改微信步数
func UpdateStepNumberApi(Data SetStepParam) models.ResponseResult {
	// 得到登录人信息
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	packHeader, errRep := GetBoundHardDeviceRequest(D)

	if errRep != nil {

		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", errRep.Error()),
			Data:    nil,
		}
	}

	//解包
	response := &mm.GetBoundHardDevicesResponse{}

	err = ParseResponseData(D, packHeader, response)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}
	device := response.GetDeviceList()[0]

	//上传步数
	packHeaderUpd, errUpd := UploadStepSetRequestRequest(D, device.HardDevice.GetDeviceId(), device.HardDevice.GetDeviceType(), Data.Number)
	if errUpd != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}
	// TODO 在看看
	fmt.Println(packHeaderUpd)
	//responseUpdate := &mm.UploadDeviceStepResponse{}
	//err = ParseResponseData1(D, packHeaderUpd, responseUpdate)
	//if err != nil {
	//	return models.ResponseResult{
	//		Code:    -8,
	//		Success: false,
	//		Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
	//		Data:    nil,
	//	}
	//}

	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    "success",
	}
}

// 修改微信步数
func UpdateStepNumberApi1(Data SetStepParam) models.ResponseResult {
	// 得到登录人信息
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	packHeader, errRep := GetBoundHardDeviceRequest(D)

	if errRep != nil {

		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", errRep.Error()),
			Data:    nil,
		}
	}

	//解包
	response := &mm.GetBoundHardDevicesResponse{}

	err = ParseResponseData(D, packHeader, response)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}
	device := response.GetDeviceList()[0]

	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix()
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).Unix()
	var req = &mm.UploadDeviceStepRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    []byte{},
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		DeviceID:   proto.String(device.HardDevice.GetDeviceId()),
		DeviceType: proto.String(device.HardDevice.GetDeviceType()),
		FromTime:   proto.Uint32(uint32(startTime)),
		ToTime:     proto.Uint32(uint32(endTime)),
		StepCount:  proto.Uint32(uint32(Data.Number)),
	}
	// 打包发送数据
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
		Cgiurl: "/cgi-bin/mmoc-bin/hardware/uploaddevicestep",
		Proxy:  D.Proxy,
		PackData: Algorithm.PackData{
			Reqdata:          reqdata,
			Cgi:              1261,
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

	Response := &mm.UploadDeviceStepResponse{}
	err = proto.Unmarshal(protobufdata, Response)

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
