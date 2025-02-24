package User

type NewSetPasswdParam struct {
	Wxid        string
	NewPassword string
	Ticket      string
}

type NewVerifyPasswdParam struct {
	Wxid     string
	Password string
}

type PrivacySettingsParam struct {
	Wxid     string
	Function int32
	Value    int32
}

type UpdateProfileParam struct {
	Wxid      string
	NickName  string
	Sex       int32
	Country   string
	Province  string
	City      string
	Signature string
}

type UploadHeadImageParam struct {
	Wxid   string
	Base64 string
}

type EmailParam struct {
	Wxid  string
	Email string
}

type SetAlisaParam struct {
	Wxid  string
	Alisa string
}

type SendVerifyMobileParam struct {
	Wxid   string
	Mobile string
	Opcode uint32
}

type BindMobileParam struct {
	Wxid       string
	Mobile     string
	Verifycode string
}

type GetQRCodeParam struct {
	Wxid  		string
	Style 		int32
}

type BindQQParam struct {
	Wxid  		string
	Account    	uint32
	Password   	string
}

type ReportMotionParam struct {
	Wxid       string
	DeviceId   string
	DeviceType string
	StepCount  int64
}

type DelSafetyInfoParam struct {
	Wxid string
	Uuid string
}
