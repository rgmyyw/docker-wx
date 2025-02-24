package Login

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/Mmtls"
	"wechatdll/baseinfo"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Tools"
)

func AwakenLogin(Wxid string) models.ResponseResult {
	D, err := comm.GetLoginata(Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	hec := &Algorithm.Client{}
	hec.Init("IOS")

	if D.ClientVersion == Algorithm.AndroidPadVersion || D.ClientVersion == Algorithm.AndroidPadVersionx {
		hec.Init("AndroidPad")
	}
	if D.ClientVersion == Algorithm.IPadVersion || D.ClientVersion == Algorithm.IPadVersionx {
		hec.Init("IOS")
	}
	if D.ClientVersion == Algorithm.WinUwpVersion {
		hec.Init("WindowsUwp")
	}
	if D.ClientVersion == Algorithm.WinVersion {
		hec.Init("Windows")
	}
	if D.ClientVersion == Algorithm.MacVersion {
		hec.Init("MAC")
	}
	if D.ClientVersion == Algorithm.CarVersion {
		hec.Init("Car")
	}
	D.DeviceType = hec.DeviceType
	Autoauthkey := &mm.AutoAuthKey{}
	_ = proto.Unmarshal(D.Autoauthkey, Autoauthkey)

	req := &mm.PushLoginURLRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    D.Sessionkey,
			Uin:           proto.Uint32(D.Uin),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(0),
		},
		Autoauthticket: proto.String(""),
		Autoauthkey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(D.Autoauthkeylen)),
			Buffer: D.Autoauthkey,
		},
		ClientId:   proto.String("iPad-Push-" + strconv.Itoa(int(time.Now().Unix())) + ".110141"),
		Devicename: proto.String(D.DeviceName),
		Opcode:     proto.Int32(3),
		RandomEncryKey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(D.Sessionkey))),
			Buffer: D.Sessionkey,
		},
		Username: proto.String(D.Wxid),
	}

	reqdata, err := proto.Marshal(req)
	hecData := Tools.Pack(D, reqdata, 654, 1)

	//	hecData := hec.HybridEcdhPackIosEn(654, D.Uin, D.Cooike, reqdata)
	httpclient := Mmtls.GenNewHttpClient(D.MmtlsKey, D.MmtlsHost)
	recvData, err := httpclient.MMtlsPost(D.MmtlsHost, "/cgi-bin/micromsg-bin/pushloginurl", hecData, D.Proxy)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}
	if len(recvData) < 32 {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：协议返回少于32字节"),
			Data:    nil,
		}
	}

	packHeader, errRep := Tools.DecodePackHeader(recvData, nil)

	if errRep != nil {

		if packHeader != nil && packHeader.GetRetCode() == baseinfo.MMRequestRetSessionTimeOut {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("链接失效：%v", errRep.Error()),
				Data:    nil,
			}
		}

		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", errRep.Error()),
			Data:    nil,
		}
	}
	//解包
	qrCodeResponse := &mm.PushLoginURLResponse{}

	err = Tools.ParseResponseData(D, packHeader, qrCodeResponse)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}
	fmt.Println(qrCodeResponse.GetUuid())
	//保存redis
	err = comm.CreateLoginData(comm.LoginData{
		Uuid:                       qrCodeResponse.GetUuid(),
		Aeskey:                     D.Sessionkey,
		NotifyKey:                  qrCodeResponse.GetNotifyKey().GetBuffer(),
		Deviceid_str:               D.Deviceid_str,
		Deviceid_byte:              D.Deviceid_byte,
		DeviceName:                 D.DeviceName,
		HybridEcdhPrivkey:          D.HybridEcdhPrivkey,
		HybridEcdhPubkey:           D.HybridEcdhPubkey,
		HybridEcdhInitServerPubKey: D.HybridEcdhInitServerPubKey,
		Cooike:                     packHeader.Session,
		MmtlsKey:                   D.MmtlsKey,
	}, "", 300)

	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("Redis ERROR：%v", err.Error()),
			Data:    nil,
		}
	}

	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    qrCodeResponse,
	}
}
