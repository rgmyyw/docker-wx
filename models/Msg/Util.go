package Msg

type ShareCardParam struct {
	Wxid         string
	ToWxid       string
	CardWxId     string
	CardNickName string
	CardAlias    string
}

type ShareLocationParam struct {
	Wxid    string
	ToWxid  string
	X       float64
	Y       float64
	Scale   float64
	Label   string
	Poiname string
	Infourl string
}

type SendAppMsgParam struct {
	Wxid   string
	ToWxid string
	Xml    string
	Type   int32
}

type SendTransmitParam struct {
	Wxid   string
	ToWxid string
	Xml    string
}

type SendShareLinkMsgParam struct {
	Wxid     string
	ToWxid   string
	Title    string
	Desc     string
	Url      string
	ThumbUrl string
}

type SendVideoMsgParam struct {
	Wxid        string
	ToWxid      string
	PlayLength  uint32
	Base64      string
	ImageBase64 string
}

type SendEmojiParam struct {
	Wxid     string
	ToWxid   string
	TotalLen int32
	Md5      string
}

type DefaultParam struct {
	Wxid    string
	ToWxid  string
	Content string
}

type Quote struct {
	Wxid         string
	ToWxid       string
	Fromusr      string
	Displayname  string
	NewMsgId     string
	MsgContent   string
	QuoteContent string
	MsgSeq       string
}


type ShareVideoMsgParam struct {
	Wxid         string
	ToWxid       string
	Xml          string
}