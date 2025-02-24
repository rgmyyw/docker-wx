package Login

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/lib"
	"wechatdll/models"
)

func RetConst(data []byte) (int64, string) {
	var Ret int32
	Ret = BytesToInt32(data[2:10])
	return int64(Ret), mm.RetConst_name[BytesToInt32(data[2:10])]
}
func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

// 最原始的二次登录版本
func Secautoauth(Wxid string) models.ResponseResult {
	D, err := comm.GetLoginata(Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//初始化Mmtls
	httpclient, MmtlsClient, err := comm.MmtlsInitialize(D.Proxy, "szshort.weixin.qq.com")
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("MMTLS初始化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	if len(D.Autoauthkey) <= 0 {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: "账号异常：Autoauthkey读取失败",
			Data:    nil,
		}
	}
	Autoauthkey := &mm.AutoAuthKey{}
	_ = proto.Unmarshal(D.Autoauthkey, Autoauthkey)

	prikey, pubkey := Algorithm.GetEcdh713Key()

	//基础设备信息
	Imei := Algorithm.IOSImei(D.Deviceid_str)
	SoftType := Algorithm.SoftType_iPad2(D.Deviceid_str)
	ClientSeqId := lib.GetClientSeqId(D.Deviceid_str)
	T := time.Now().Unix()
	WCExtInfoseq := Algorithm.GetIOSExtSpamInfoAndDeviceToken(D.Wxid, D.Deviceid_str, D.DeviceName, D.DeviceToken, T)

	req := &mm.AutoAuthRequest{
		RsaReqData: &mm.AutoAuthRsaReqData{
			AesEncryptKey: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(Autoauthkey.EncryptKey.Buffer))),
				Buffer: Autoauthkey.EncryptKey.Buffer,
			},
			CliPubEcdhkey: &mm.ECDHKey{
				Nid: proto.Int32(713),
				Key: &mm.SKBuiltinBufferT{
					ILen:   proto.Uint32(uint32(len(pubkey))),
					Buffer: pubkey,
				},
			},
		},
		AesReqData: &mm.AutoAuthAesReqData{
			BaseRequest: &mm.BaseRequest{
				SessionKey:    []byte{},
				Uin:           proto.Uint32(D.Uin),
				DeviceId:      D.Deviceid_byte,
				ClientVersion: proto.Int32(int32(D.ClientVersion)),
				DeviceType:    []byte(D.DeviceType),
				Scene:         proto.Uint32(0),
			},
			BaseReqInfo: &mm.BaseAuthReqInfo{},
			AutoAuthKey: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(D.Autoauthkey))),
				Buffer: D.Autoauthkey,
			},
			Channel:      proto.Int32(0),
			Imei:         &Imei,
			SoftType:     &SoftType,
			BuiltinIpseq: proto.Uint32(0),
			ClientSeqId:  &ClientSeqId, //
			//BundleId:     proto.String("com.tencent.xin"),
			DeviceName: proto.String(D.DeviceName), //9
			DeviceType: proto.String("iPad"),       //10
			Language:   proto.String("zh_CN"),      //11
			TimeZone:   proto.String("8.0"),        //12

			ExtSpamInfo: &mm.SKBuiltinBufferT{ //15
				ILen:   proto.Uint32(uint32(len(WCExtInfoseq))),
				Buffer: WCExtInfoseq,
			},
		},
	}

	reqdata, err := proto.Marshal(req)
	fmt.Println(hex.EncodeToString(reqdata))
	//Host := comm.GetIp(*D)

	hec := &Algorithm.Client{}
	hec.Init2("IOS")
	hecData := hec.HybridEcdhPackIosEn2(763, 0, D.Cooike, reqdata, D.Loginecdhkey)
	recvData, err := httpclient.MMtlsPost(D.MmtlsHost, "/cgi-bin/micromsg-bin/secautoauth", hecData, D.Proxy)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
	}

	if len(recvData) <= 31 {
		Ret, name := RetConst(recvData)
		error := "您已退出微信/session过期"
		if Ret != -13 {
			error = name
			if name == "" {
				error = "微信未知的错误信息"
			}
		}

		return models.ResponseResult{
			Code:    Ret,
			Success: false,
			Message: error,
			Data:    nil,
			Debug:   hex.EncodeToString(recvData),
		}
	}

	ph1 := hec.HybridEcdhPackIosUn(recvData)
	//解包
	loginRes := mm.UnifyAuthResponse{}
	err = proto.Unmarshal(ph1.Data, &loginRes)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	if loginRes.GetBaseResponse().GetRet() != 0 || loginRes.BaseResponse == nil || loginRes.AuthSectResp == nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: "登录失败：您可能已退出微信",
			Data:    loginRes,
			Debug:   hex.EncodeToString(recvData),
		}
	}

	if loginRes.GetAuthSectResp().GetSessionKey().GetBuffer() == nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: "登录失败：无法获取返回的Key",
			Data:    loginRes,
			Debug:   hex.EncodeToString(recvData),
		}
	}

	Wx_loginecdhkey := Algorithm.DoECDH713Key(prikey, loginRes.GetAuthSectResp().GetSvrPubEcdhkey().GetKey().GetBuffer())
	m := md5.New()
	m.Write(Wx_loginecdhkey)
	D.Loginecdhkey = Wx_loginecdhkey
	D.Uin = ph1.Uin
	ecdhdecrptkey := m.Sum(nil)
	D.Cooike = ph1.Cookies
	D.Sessionkey = Algorithm.AesDecrypt(loginRes.GetAuthSectResp().GetSessionKey().GetBuffer(), ecdhdecrptkey)
	D.Autoauthkey = loginRes.GetAuthSectResp().GetAutoAuthKey().GetBuffer()
	D.Autoauthkeylen = int32(loginRes.GetAuthSectResp().GetAutoAuthKey().GetILen())
	D.Serversessionkey = loginRes.GetAuthSectResp().GetServerSessionKey().GetBuffer()
	D.Clientsessionkey = loginRes.GetAuthSectResp().GetClientSessionKey().GetBuffer()
	D.MmtlsKey = MmtlsClient

	err = comm.CreateLoginData(*D, D.Wxid, 0)

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
		Message: "登陆成功",
		Data:    loginRes,
	}
}
