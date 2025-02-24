package Tools

type UploadParam struct {
	Wxid   string
	Base64 string
}

type DownloadVoiceParam struct {
	Wxid         string
	FromUserName string
	MsgId        uint32
	Bufid        string
	Length       int
}

type DownloadData struct {
	Base64 []byte
	Length uint32
}

type DownloadParam struct {
	Wxid         string
	ToWxid       string
	MsgId        uint32
	DataLen      int
	Section      Section
	CompressType int
}

type CdnDownloadImageParam struct {
	Wxid       string
	FileNo     string
	FileAesKey string
}

type CdnImageBase64 struct {
	Image string
}

type DownloadAppAttachParam struct {
	Wxid     string
	AppID    string
	AttachId string
	UserName string
	DataLen  int
	Section  Section
}

type Section struct {
	StartPos uint32
	DataLen  uint32
}

type GetCertParam struct {
	Wxid    string
	Version uint32
}

type ThirdAppGrantParam struct {
	Wxid  string
	Url   string
	Appid string
}
