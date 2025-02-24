package Wxapp

type DefaultParam struct {
	Wxid  string
	Appid string
}

type JSOperateWxParam struct {
	Appid string
	Data  string
	Opt   int
	Wxid  string
}

type CloudCallParam struct {
	Appid string
	Data  string
	Wxid  string
}

type GetWxAppRecordParam struct {
	Wxid string
}

type AddWxAppRecordParam struct {
	Wxid     string
	Username string
}

type SessionidQRParam struct {
	Wxid      string
	Appid     string
	Sessionid string
	timeStamp string
	nonceStr  string
	Package   string
	PaySign   string
}

type QrcodeAuthLoginParam struct {
	Wxid string
	Uuid  string
}

type CheckVerifyCodeData struct {
	Appid      string
	Mobile     string
	VerifyCode string
	Wxid       string
}

type CheckVerifyCodeNData struct {
	Appid      string
	Mobile     string
	VerifyCode string
	Wxid       string
	Opcode     int
}

type PostVerifyCodeParam struct {
	Appid  string
	Mobile string
	Opcode int
	Wxid   string
}

type DelMobileData struct {
	Appid  string
	Mobile string
	Opcode int
	Wxid   string
}

type AddAvatarParam struct {
	Wxid     string
	Appid    string
	NickName string
	AFilekey string
}

type AddAvatarImgParam struct {
	Wxid    string
	Appid   string
	JPGlink string
}
