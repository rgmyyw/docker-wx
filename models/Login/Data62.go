package Login

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/Mmtls"
	"wechatdll/comm"
	"wechatdll/lib"
	"wechatdll/models"

	"github.com/golang/protobuf/proto"
)

func Data62(Data Data62LoginReq, domain string) models.ResponseResult {
	// 获取username为key的缓存
	D, _ := comm.GetLoginata(Data.UserName)
	if D == nil || D.Wxid == "" {
		// 没有缓存, 初始化新的账号环境
		D = GeniPhoneLoginData(Data)
	} else {
		D = UpdateiPhoneLoginData(D, Data)
	}
	// 更新MmtlsHost
	if D.MmtlsHost == "" {
		D.MmtlsHost = domain
	} else {
		domain = D.MmtlsHost
	}

	//初始化Mmtls
	var httpclient *Mmtls.HttpClientModel
	var err error
	if D.MmtlsKey == nil {
		httpclient, D.MmtlsKey, err = comm.MmtlsInitialize(Data.Proxy, domain)
		if err != nil {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("MMTLS初始化失败：%v", err.Error()),
				Data:    nil,
			}
		}
	} else {
		httpclient = Mmtls.GenNewHttpClient(D.MmtlsKey, domain)
	}

	// 请求设备device_token
	if D.DeviceToken.TrustResponseData == nil || D.DeviceToken.TrustResponseData.DeviceToken == nil || *D.DeviceToken.TrustResponseData.DeviceToken == "" {
		D.DeviceToken, err = IPadGetDeviceToken(D.Deviceid_str, D.RomModel, D.DeviceName, D.DeviceType, int32(D.ClientVersion), *httpclient, D.Proxy, D.MmtlsHost)
		if err != nil {
			// 请求失败则放空结构
			D.DeviceToken = mm.TrustResponse{}
		}
	}
	prikey, pubkey := Algorithm.GetEcdh713Key()

	Wcstf := Algorithm.IphoneWcstf(D.Wxid)
	Wcste := Algorithm.IphoneWcste(0, 0)

	accountRequest := &mm.ManualAuthRsaReqData{
		RandomEncryKey: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(D.Aeskey))),
			Buffer: D.Aeskey,
		},
		CliPubEcdhkey: &mm.ECDHKey{
			Nid: proto.Int32(713),
			Key: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(pubkey))),
				Buffer: pubkey,
			},
		},
		UserName: proto.String(Data.UserName),
		Pwd:      proto.String(lib.MD5ToLower(Data.Password)),
	}

	ccData := &mm.CryptoData{
		Version:     []byte("00000006"),
		Type:        proto.Uint32(1),
		EncryptData: Algorithm.GetiPhoneNewSpamData(D.Deviceid_str, Data.DeviceName, D.DeviceToken),
		Timestamp:   proto.Uint32(uint32(time.Now().Unix())),
		Unknown5:    proto.Uint32(5),
		Unknown6:    proto.Uint32(0),
	}

	ccDataseq, _ := proto.Marshal(ccData)

	DeviceTokenCCD := &mm.DeviceToken{
		Version:   proto.String(""),
		Encrypted: proto.Uint32(1),
		Data: &mm.SKBuiltinStringT{
			String_: proto.String(D.DeviceToken.GetTrustResponseData().GetDeviceToken()),
		},
		TimeStamp: proto.Uint32(uint32(time.Now().Unix())),
		Optype:    proto.Uint32(2),
		Uin:       proto.Uint32(0),
	}
	DeviceTokenCCDPB, _ := proto.Marshal(DeviceTokenCCD)

	WCExtInfo := &mm.WCExtInfo{
		Wcstf: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(Wcstf))),
			Buffer: Wcstf,
		},
		Wcste: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(Wcste))),
			Buffer: Wcste,
		},
		CcData: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(ccDataseq))),
			Buffer: ccDataseq,
		},
		DeviceToken: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(DeviceTokenCCDPB))),
			Buffer: DeviceTokenCCDPB,
		},
	}

	WCExtInfoseq, _ := proto.Marshal(WCExtInfo)

	ClientSeqId := fmt.Sprintf("%v_%v", D.Deviceid_str, time.Now().Unix())
	uuid1, _ := Algorithm.IOSUuid(D.Deviceid_str)

	deviceRequest := &mm.ManualAuthAesReqData{
		BaseRequest: &mm.BaseRequest{
			SessionKey:    []byte{},
			Uin:           proto.Uint32(0),
			DeviceId:      D.Deviceid_byte,
			ClientVersion: proto.Int32(int32(D.ClientVersion)),
			DeviceType:    []byte(D.DeviceType),
			Scene:         proto.Uint32(1),
		},
		BaseReqInfo:  &mm.BaseAuthReqInfo{},
		Imei:         &D.Imei,
		SoftType:     &D.SoftType,
		BuiltinIpseq: proto.Uint32(0),
		ClientSeqId:  &ClientSeqId,
		DeviceName:   proto.String(Data.DeviceName),
		DeviceType:   proto.String(D.DeviceType),
		Language:     proto.String("zh_CN"),
		TimeZone:     proto.String("8.0"),
		Channel:      proto.Int(0),
		TimeStamp:    proto.Uint32(uint32(time.Now().Unix())),
		DeviceBrand:  proto.String("Apple"),
		RealCountry:  proto.String("CN"),
		BundleId:     proto.String("com.tencent.xin"),
		AdSource:     &uuid1,
		IphoneVer:    &D.RomModel,
		InputType:    proto.Uint32(2),
		ExtSpamInfo: &mm.SKBuiltinBufferT{
			ILen:   proto.Uint32(uint32(len(WCExtInfoseq))),
			Buffer: WCExtInfoseq,
		},
	}

	requset := &mm.SecManualLoginRequest{
		RsaReqData: accountRequest,
		AesReqData: deviceRequest,
	}

	reqdata, _ := proto.Marshal(requset)

	// 更新代理IP
	if Data.Proxy.ProxyIp != "" && Data.Proxy.ProxyIp != "String" {
		D.Proxy = Data.Proxy
	}
	// 更新代理IP
	if Data.Proxy.ProxyIp != "" && Data.Proxy.ProxyIp != "String" {
		D.Proxy = Data.Proxy
	}
	// 更新代理IP
	if Data.Proxy.ProxyIp != "" && Data.Proxy.ProxyIp != "String" {
		D.Proxy = Data.Proxy
	}
	// 临时缓存上下文, 以便被异常中断后再次保持上下文重试
	err = comm.CreateLoginData(*D, D.Wxid, 7200)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}

	hec := &Algorithm.Client{}
	hec.Init("IOS")
	hecData := hec.HybridEcdhPackIosEn(252, 0, nil, reqdata)
	// 遇到mmtls失败, 则重新握手
	retrys := 2
	doRetry := true
	var recvData []byte
	for retrys > 0 && doRetry == true {
		doRetry = false
		retrys--
		recvData, err = httpclient.MMtlsPost(domain, "/cgi-bin/micromsg-bin/secmanualauth", hecData, Data.Proxy)
		if err != nil && strings.Contains(err.Error(), "MMTLS") {
			// mmtls异常, 重新握手
			httpclient, D.MmtlsKey, err = comm.MmtlsInitialize(Data.Proxy, domain)
			if err != nil {
				return models.ResponseResult{
					Code:    -8,
					Success: false,
					Message: fmt.Sprintf("MMTLS初始化失败：%v", err.Error()),
					Data:    nil,
				}
			}
			// 重新提交
			doRetry = true
		}
	}
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}
	if len(recvData) <= 31 {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("组包异常, 返回31字节"),
			Data:    nil,
		}
	}

	ph1 := hec.HybridEcdhPackIosUn(recvData)
	loginRes := mm.UnifyAuthResponse{}
	err = proto.Unmarshal(ph1.Data, &loginRes)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}

	if loginRes.GetBaseResponse().GetRet() == 0 && loginRes.GetUnifyAuthSectFlag() > 0 {
		D.Cooike = ph1.Cookies
		D.MmtlsHost = domain

		Wx_loginecdhkey := Algorithm.DoECDH713Key(prikey, loginRes.GetAuthSectResp().GetSvrPubEcdhkey().GetKey().GetBuffer())
		Wx_loginecdhkeylen := int32(len(Wx_loginecdhkey))
		m := md5.New()
		m.Write(Wx_loginecdhkey[:Wx_loginecdhkeylen])
		ecdhdecrptkey := m.Sum(nil)
		D.Loginecdhkey = Wx_loginecdhkey //TODO: 用于计算组包的checkSum, 使用什么做盐?
		D.Uin = loginRes.GetAuthSectResp().GetUin()
		D.Wxid = loginRes.GetAcctSectResp().GetUserName()
		D.Alais = loginRes.GetAcctSectResp().GetAlias()
		D.Mobile = loginRes.GetAcctSectResp().GetBindMobile()
		D.NickName = loginRes.GetAcctSectResp().GetNickName()
		D.Email = loginRes.GetAcctSectResp().GetBindEmail()
		D.Sessionkey = Algorithm.AesDecrypt(loginRes.GetAuthSectResp().GetSessionKey().GetBuffer(), ecdhdecrptkey)
		D.Sessionkey_2 = loginRes.GetAuthSectResp().GetSessionKey().GetBuffer()
		D.Autoauthkey = loginRes.GetAuthSectResp().GetAutoAuthKey().GetBuffer()
		D.Autoauthkeylen = int32(loginRes.GetAuthSectResp().GetAutoAuthKey().GetILen())
		D.Serversessionkey = loginRes.GetAuthSectResp().GetServerSessionKey().GetBuffer()
		D.Clientsessionkey = loginRes.GetAuthSectResp().GetClientSessionKey().GetBuffer()
		D.MmtlsHost = comm.Rmu0000(*loginRes.NetworkSectResp.BuiltinIplist.ShortConnectIplist[0].Host)
		D.MarsHost = comm.Rmu0000(*loginRes.NetworkSectResp.BuiltinIplist.LongConnectIplist[0].Host)
		D.RsaPublicKey = pubkey
		D.RsaPrivateKey = prikey
		// 更新代理IP
		if Data.Proxy.ProxyIp != "" && Data.Proxy.ProxyIp != "String" {
			D.Proxy = Data.Proxy
		}
		err := comm.CreateLoginData(*D, D.Wxid, 0)

		if err != nil {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("系统异常：%v", err.Error()),
				Data:    nil,
			}
		}

		return models.ResponseResult{
			Code:    0,
			Success: true,
			Message: "成功",
			Data:    loginRes,
			Data62:  lib.Get62Data(D.Deviceid_str),
		}
	}

	//30系列转向
	if loginRes.GetBaseResponse().GetRet() == -301 {
		D.MmtlsHost = comm.Rmu0000(*loginRes.NetworkSectResp.BuiltinIplist.ShortConnectIplist[0].Host)
		D.MarsHost = comm.Rmu0000(*loginRes.NetworkSectResp.BuiltinIplist.LongConnectIplist[0].Host)
		D.MmtlsKey = nil
		// 更新代理IP
		if Data.Proxy.ProxyIp != "" && Data.Proxy.ProxyIp != "String" {
			D.Proxy = Data.Proxy
		}
		err := comm.CreateLoginData(*D, D.Wxid, 7200)
		if err != nil {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("系统异常：%v", err.Error()),
				Data:    nil,
			}
		}

		return Data62(Data, D.MmtlsHost)
	}
	/*

		// 自动过滑块
		if strings.Index(loginRes.GetBaseResponse().GetErrMsg().GetString_(), "环境存在异常") >= 0 || strings.Contains(loginRes.GetBaseResponse().GetErrMsg().GetString_(), "The system has detected an abnormal login environment.") {
			// 过滑块并根据结果判断
			if err := LoginOCR(loginRes.GetBaseResponse().GetErrMsg().GetString_()); err == nil {
				// 滑块成功, 再次登录
				return Data62(Data, domain)
			} else {
				// 返回异常
				return models.ResponseResult{
					Code:    -8,
					Success: false,
					Message: fmt.Sprintf("系统异常：%v", err.Error()),
					Data:    nil,
				}
			}
		}
	*/

	return models.ResponseResult{
		Code:    -8,
		Success: false,
		Message: "失败",
		Data:    loginRes,
	}
}

func GenIpadLoginData(request Data62LoginReq) *comm.LoginData {
	deviceId := lib.Get62Key(request.Data62)
	if deviceId[:2] != "49" {
		deviceId = "49" + deviceId[2:]
	}
	deviceIdByte, _ := hex.DecodeString(deviceId)
	if request.DeviceName == "" || request.DeviceName == "string" {
		request.DeviceName = "iPad Pro 13(M4)"
	}
	D := &comm.LoginData{
		Wxid:          request.UserName,
		Pwd:           request.Password,
		Aeskey:        []byte(lib.RandSeq(16)), //随机密钥
		Deviceid_str:  deviceId,
		Deviceid_byte: deviceIdByte,
		DeviceType:    Algorithm.IPadDeviceType,
		ClientVersion: Algorithm.IPadVersion,
		DeviceName:    request.DeviceName,
		MmtlsHost:     Algorithm.MmtlsShortHost,
		MarsHost:      Algorithm.MmtlsLongHost,
		Proxy:         request.Proxy,
		DeviceToken:   mm.TrustResponse{},
		RomModel:      Algorithm.IPadModel,
		Imei:          Algorithm.IOSImei(deviceId),
		SoftType:      Algorithm.SoftType_iPad(deviceId, Algorithm.IPadOsVersion, Algorithm.IPadModel),
		OsVersion:     Algorithm.IPadOsVersion,
	}
	return D
}

func GeniPhoneLoginData(request Data62LoginReq) *comm.LoginData {
	deviceId := lib.Get62Key(request.Data62)
	if deviceId[:2] != "49" {
		deviceId = "49" + deviceId[2:]
	}
	deviceIdByte, _ := hex.DecodeString(deviceId)
	if request.DeviceName == "" || request.DeviceName == "string" {
		request.DeviceName = "iPhone 16 Pro Max"
	}
	D := &comm.LoginData{
		Wxid:          request.UserName,
		Pwd:           request.Password,
		Aeskey:        []byte(lib.RandSeq(16)), //随机密钥
		Deviceid_str:  deviceId,
		Deviceid_byte: deviceIdByte,
		DeviceType:    Algorithm.IPhoneDeviceType,
		ClientVersion: Algorithm.IPhoneVersion,
		DeviceName:    request.DeviceName,
		MmtlsHost:     Algorithm.MmtlsShortHost,
		MarsHost:      Algorithm.MmtlsLongHost,
		Proxy:         request.Proxy,
		DeviceToken:   mm.TrustResponse{},
		RomModel:      Algorithm.IPhoneModel,
		Imei:          Algorithm.IOSImei(deviceId),
		SoftType:      Algorithm.SoftType_iPad(deviceId, Algorithm.IPhoneOsVersion, Algorithm.IPhoneModel),
		OsVersion:     Algorithm.IPhoneOsVersion,
	}
	return D
}

func UpdateIpadLoginData(D *comm.LoginData, Data Data62LoginReq) *comm.LoginData {
	D.Pwd = Data.Password
	D.Data62 = Data.Data62
	if D.DeviceType == "" {
		D.DeviceType = Algorithm.IPadDeviceType
	}
	if D.ClientVersion == 0 {
		D.ClientVersion = Algorithm.IPadVersion
	}
	if D.DeviceName == "" || D.DeviceName == "string" {
		if Data.DeviceName == "" || Data.DeviceName == "string" {
			D.DeviceName = "iPad Pro 13(M4)"
		} else {
			D.DeviceName = Data.DeviceName
		}
	}
	if D.MmtlsHost == "" {
		D.MmtlsHost = Algorithm.MmtlsShortHost
	}
	if D.MarsHost == "" {
		D.MarsHost = Algorithm.MmtlsLongHost
	}
	if D.RomModel == "" {
		D.RomModel = Algorithm.IPadModel
	}
	if D.Imei == "" {
		D.Imei = Algorithm.IOSImei(D.Deviceid_str)
	}
	if D.SoftType == "" {
		D.SoftType = Algorithm.SoftType_iPad(D.Deviceid_str, Algorithm.IPadOsVersion, Algorithm.IPadModel)
	}
	if D.OsVersion == "" {
		D.OsVersion = Algorithm.IPadOsVersion
	}
	return D
}

func UpdateiPhoneLoginData(D *comm.LoginData, Data Data62LoginReq) *comm.LoginData {
	D.Pwd = Data.Password
	D.Data62 = Data.Data62
	if D.DeviceType == "" {
		D.DeviceType = Algorithm.IPhoneDeviceType
	}
	if D.ClientVersion == 0 {
		D.ClientVersion = Algorithm.IPhoneVersion
	}
	if D.DeviceName == "" || D.DeviceName == "string" {
		if Data.DeviceName == "" || Data.DeviceName == "string" {
			D.DeviceName = "iPhone 16 Pro Max"
		} else {
			D.DeviceName = Data.DeviceName
		}
	}
	if D.MmtlsHost == "" {
		D.MmtlsHost = Algorithm.MmtlsShortHost
	}
	if D.MarsHost == "" {
		D.MarsHost = Algorithm.MmtlsLongHost
	}
	if D.RomModel == "" {
		D.RomModel = Algorithm.IPhoneModel
	}
	if D.Imei == "" {
		D.Imei = Algorithm.IOSImei(D.Deviceid_str)
	}
	if D.SoftType == "" {
		D.SoftType = Algorithm.SoftType_iPad(D.Deviceid_str, Algorithm.IPhoneOsVersion, Algorithm.IPhoneModel)
	}
	if D.OsVersion == "" {
		D.OsVersion = Algorithm.IPhoneOsVersion
	}
	return D
}
