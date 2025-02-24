package Login

import (
	"container/list"
	"crypto/md5"
	"fmt"
	"strings"
	"wechatdll/Algorithm"
	"wechatdll/comm"
	"wechatdll/models"
)

func CheckSecManualAuth(Data comm.LoginData, mmtlshost string) models.ResponseResult {
	//开始登陆
	loginRes, prikey, pubkey, Cookie, DeviceToken, err := SecManualAuth(Data, mmtlshost)

	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: "登陆异常",
			Data:    err.Error(),
		}
	}

	//登陆成功
	if loginRes.GetBaseResponse().GetRet() == 0 && loginRes.GetUnifyAuthSectFlag() > 0 {
		Wx_loginecdhkey := Algorithm.DoECDH713Key(prikey, loginRes.GetAuthSectResp().GetSvrPubEcdhkey().GetKey().GetBuffer())
		m := md5.New()
		m.Write(Wx_loginecdhkey)
		Data.Loginecdhkey = Wx_loginecdhkey
		ecdhdecrptkey := m.Sum(nil)
		Data.Uin = loginRes.GetAuthSectResp().GetUin()
		Data.Wxid = loginRes.GetAcctSectResp().GetUserName()
		Data.Alais = loginRes.GetAcctSectResp().GetAlias()
		Data.Mobile = loginRes.GetAcctSectResp().GetBindMobile()
		Data.Email = loginRes.GetAcctSectResp().GetBindEmail()
		Data.NickName = loginRes.GetAcctSectResp().GetNickName()
		Data.Cooike = Cookie
		Data.Sessionkey = Algorithm.AesDecrypt(loginRes.GetAuthSectResp().GetSessionKey().GetBuffer(), ecdhdecrptkey)
		Data.Sessionkey_2 = loginRes.GetAuthSectResp().GetSessionKey().GetBuffer()
		Data.Autoauthkey = loginRes.GetAuthSectResp().GetAutoAuthKey().GetBuffer()
		Data.Autoauthkeylen = int32(loginRes.GetAuthSectResp().GetAutoAuthKey().GetILen())
		Data.Serversessionkey = loginRes.GetAuthSectResp().GetServerSessionKey().GetBuffer()
		Data.Clientsessionkey = loginRes.GetAuthSectResp().GetClientSessionKey().GetBuffer()
		Data.MmtlsHost = mmtlshost
		Data.ClientVersion = Algorithm.IPadVersion
		Data.DeviceType = Algorithm.IPadDeviceType
		Data.DeviceToken = DeviceToken
		Data.MmtlsHost = comm.Rmu0000(*loginRes.NetworkSectResp.BuiltinIplist.ShortConnectIplist[0].Host)
		Data.MarsHost = comm.Rmu0000(*loginRes.NetworkSectResp.BuiltinIplist.LongConnectIplist[0].Host)
		Data.RsaPublicKey = pubkey
		Data.RsaPrivateKey = prikey

		err := comm.CreateLoginData(Data, Data.Wxid, 0)

		if err != nil {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("系统异常：%v", err.Error()),
				Data:    nil,
			}
		}

		type Suc struct {
			Wxid     string
			Alais    string
			NickName string
			Mobile   string
		}
		//return models.ResponseResult{
		//	Code:    0,
		//	Success: true,
		//	Message: "登陆成功",
		//	Data:    CheckLoginRes{
		//		Uuid:     Data.Uuid,
		//		WxId:     Data.Wxid,
		//		NickName: Data.NickName,
		//		Status:    2,
		//		Device:   Data.DeviceName,
		//		HeadUrl:  Data.HeadUrl,
		//		Mobile:   Data.Mobile,
		//		Email:    Data.Email,
		//		Alias:    Data.Alais,
		//		Data62:   "",
		//	},
		//}
		return models.ResponseResult{
			Code:    0,
			Success: true,
			Message: "登陆成功",
			Data:    loginRes,
		}
	}

	//30系列转向
	if loginRes.GetBaseResponse().GetRet() == -301 {
		var Wx_newLong_Host, Wx_newshort_Host, Wx_newshortext_Host list.List

		dns_info := loginRes.GetNetworkSectResp().GetNewHostList().GetList()
		for _, v := range dns_info {
			if v.GetHost() == "long.weixin.qq.com" {
				ip_info := loginRes.GetNetworkSectResp().GetBuiltinIplist().GetLongConnectIplist()
				for _, ip := range ip_info {
					host := ip.GetHost()
					host = strings.Replace(host, string(byte(0x00)), "", -1)
					if host == v.GetRedirect() {
						Wx_newLong_Host.PushBack(host)
					}
				}
			} else if v.GetHost() == "short.weixin.qq.com" {
				ip_info := loginRes.GetNetworkSectResp().GetBuiltinIplist().GetShortConnectIplist()
				for _, ip := range ip_info {
					host := ip.GetHost()
					host = strings.Replace(host, string(byte(0x00)), "", -1)
					if host == v.GetRedirect() {
						Wx_newshort_Host.PushBack(host)
					}
				}
			} else if v.GetHost() == "extshort.weixin.qq.com" {
				ip_info := loginRes.GetNetworkSectResp().GetBuiltinIplist().GetShortConnectIplist()
				for _, ip := range ip_info {
					host := ip.GetHost()
					host = strings.Replace(host, string(byte(0x00)), "", -1)
					if host == v.GetRedirect() {
						Wx_newshortext_Host.PushBack(host)
					}
				}
			}
		}
		return CheckSecManualAuth(Data, Wx_newshort_Host.Front().Value.(string))
	}

	//否则就是包有问题
	return models.ResponseResult{
		Code:    -8,
		Success: false,
		Message: "登陆异常",
		Data:    loginRes,
	}
}
