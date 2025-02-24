package Login

import (
	"wechatdll/models"
)

type ExtDeviceLoginConfirmParam struct {
	Wxid string
	Url  string
}

type Data62LoginReq struct {
	UserName   string
	Password   string
	Data62     string
	DeviceName string
	Proxy      models.ProxyInfo
}

type Data62SMSAgainReq struct {
	Url			string
	Cookie		string
	Proxy      models.ProxyInfo
}

type Data62SMSVerifyReq struct {
	Url			string
	Cookie		string
	Sms			string
	Proxy      models.ProxyInfo
}

type Data62QRCodeVerifyReq struct {
	Url			string
	Proxy       models.ProxyInfo
}
