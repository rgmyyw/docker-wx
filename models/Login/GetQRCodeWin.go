package Login

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/lib"

	"github.com/golang/protobuf/proto"

	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
)

func GetQRCODEWin(Data GetQRReq) models.ResponseResult2 {
	//初始化Mmtls
	httpclient, MmtlsClient, err := comm.MmtlsInitialize(Data.Proxy, Algorithm.MmtlsShortHost)
	if err != nil {
		return models.ResponseResult2{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("MMTLS初始化失败：%v", err.Error()),
			Data:    nil,
		}
	}
	Algorithm.IPadDeviceType = Algorithm.WinDeviceType
	Algorithm.IPadVersion = Algorithm.WinVersion
	aeskey := []byte(lib.RandSeq(16)) //获取随机密钥
	deviceid := Data.DeviceID
	devicelIdByte, _ := hex.DecodeString(deviceid)

	DeviceToken, err := IPadGetDeviceToken(deviceid, Algorithm.WinModel, Data.DeviceName, Algorithm.WinDeviceType, int32(Algorithm.WinVersion), *httpclient, Data.Proxy, Algorithm.MmtlsShortHost)
	if err != nil {
		DeviceToken = mm.TrustResponse{}
	}

	req := &mm.GetLoginQRCodeRequest{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    []byte{},
			Uin:           proto.Uint32(0),
			DeviceId:      devicelIdByte,
			ClientVersion: proto.Int32(int32(Algorithm.WinVersion)),
			DeviceType:    []byte(Algorithm.WinDeviceType),
			Scene:         proto.Uint32(0),
		},
		RandomEncryKey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(aeskey))),
			Buffer: aeskey,
		},
		Opcode:           proto.Uint32(0),
		MsgContextPubKey: nil,
	}

	reqdata, err := proto.Marshal(req)

	if err != nil {
		return models.ResponseResult2{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}

	hec := &Algorithm.Client{}
	hec.Init("Windows")
	hypack := hec.HybridEcdhPackIosEn(502, 0, nil, reqdata)
	recvData, err := httpclient.MMtlsPost(Algorithm.MmtlsShortHost, "/cgi-bin/micromsg-bin/getloginqrcode", hypack, Data.Proxy)
	if err != nil {
		return models.ResponseResult2{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}
	ph1 := hec.HybridEcdhPackIosUn(recvData)
	getloginQRRes := mm.GetLoginQRCodeResponse{}

	err = proto.Unmarshal(ph1.Data, &getloginQRRes)

	if err != nil {
		return models.ResponseResult2{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	if getloginQRRes.GetBaseResponse().GetRet() == 0 {
		if getloginQRRes.Uuid == nil || *getloginQRRes.Uuid == "" {
			return models.ResponseResult2{
				Code:    -9,
				Success: false,
				Message: "取码过于频繁",
				Data:    getloginQRRes.GetBaseResponse(),
			}
		}

		//保存redis
		err := comm.CreateLoginData(comm.LoginData{
			Uuid:          getloginQRRes.GetUuid(),
			Aeskey:        aeskey,
			NotifyKey:     getloginQRRes.GetNotifyKey().GetBuffer(),
			Deviceid_str:  deviceid,
			Deviceid_byte: devicelIdByte,
			DeviceName:    Data.DeviceName,
			ClientVersion: Algorithm.WinVersion,
			Cooike:        ph1.Cookies,
			Proxy:         Data.Proxy,
			MmtlsKey:      MmtlsClient,
			DeviceToken:   DeviceToken,
		}, "", 300)

		if err == nil {
			return models.ResponseResult2{
				Code:    1,
				Success: true,
				Message: "成功",
				Data: GetQRRes{
					baseResponse: GetQRResErr{
						Ret:   getloginQRRes.GetBaseResponse().GetRet(),
						Error: getloginQRRes.GetBaseResponse().GetErrMsg().GetString_(),
					},
					QrBase64:    fmt.Sprintf("data:image/jpg;base64,%v", base64.StdEncoding.EncodeToString(getloginQRRes.GetQrcode().GetBuffer())),
					Uuid:        getloginQRRes.GetUuid(),
					QrUrl:       "https://api.qrserver.com/v1/create-qr-code/?data=http://weixin.qq.com/x/" + getloginQRRes.GetUuid(),
					ExpiredTime: time.Unix(int64(getloginQRRes.GetExpiredTime()), 0).Format("2006-01-02 15:04:05"),
				},
				Data62:   lib.Get62Data(deviceid),
				DeviceId: deviceid,
			}
		}
	}

	return models.ResponseResult2{
		Code:    -0,
		Success: false,
		Message: "未知的错误",
		Data:    getloginQRRes,
	}
}
