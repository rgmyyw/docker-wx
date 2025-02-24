package baseinfo

import (
	"wechatdll/Cilent/wechat"
	"wechatdll/models/baseutils"
)

// SDKVersion 本协议SDK版本号
var SDKVersion = string("1.0.0")

// ClientVersion 微信版本号
var ClientVersion = uint32(0x1800312A) //0x1800312A   0x1800312A 0x1800312A
var PlistVersion = uint32(0x1800312A)  //plist-version
var Md5OfMachOHeader = string("d05a80a94b6c2e3c31424403437b6e18")

type SKBuiltinString struct {
	Str                  *string  `protobuf:"bytes,1,opt,name=str" json:"str,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type BaseResponse struct {
	Ret                  *int32           `protobuf:"varint,1,opt,name=ret" json:"ret,omitempty"`
	ErrMsg               *SKBuiltinString `protobuf:"bytes,2,opt,name=errMsg" json:"errMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

// DeviceInfo 62设备信息
type DeviceInfo struct {
	UUIDOne            string `json:"uuidone"`
	UUIDTwo            string `json:"uuidtwo"`
	Imei               string `json:"imei"`
	DeviceID           []byte `json:"deviceid"`
	DeviceName         string `json:"devicename"`
	DeviceMac          string `json:"Devicemac"`
	TimeZone           string `json:"timezone"`
	Language           string `json:"language"`
	DeviceBrand        string `json:"devicebrand"`
	RealCountry        string `json:"realcountry"`
	IphoneVer          string `json:"iphonever"`
	BundleID           string `json:"boudleid"`
	OsType             string `json:"ostype"`
	AdSource           string `json:"adsource"`
	OsTypeNumber       string `json:"ostypenumber"`
	CoreCount          uint32 `json:"corecount"`
	CarrierName        string `json:"carriername"`
	SoftTypeXML        string `json:"softtypexml"`
	ClientCheckDataXML string `json:"clientcheckdataxml"`
	// extInfo
	GUID2       string `json:"GUID2"`
	DeviceToken *wechat.TrustResp
}

func (d *DeviceInfo) SetDeviceId(deviceId string) {
	d.Imei = deviceId
	d.DeviceID = baseutils.HexStringToBytes(deviceId)
	d.DeviceID[0] = 0x49
}

// LoginDataInfo 62/16 数据登陆
type LoginDataInfo struct {
	Type     byte
	UserName string
	PassWord string
	//伪密码
	NewPassWord string
	//登录数据 62/A16
	LoginData string
	Ticket    string
	NewType   int
	Language  string
}

// WifiInfo WifiInfo
type WifiInfo struct {
	Name      string
	WifiBssID string
}

// ModifyItem 修改用户信息项
type ModifyItem struct {
	CmdID uint32
	Len   uint32
	Data  []byte
}

// HeadImgItem 头像数据项
type HeadImgItem struct {
	ImgPieceData []byte
	TotalLen     uint32
	StartPos     uint32
	ImgHash      string
}

// RevokeMsgItem 撤回消息项
type RevokeMsgItem struct {
	FromUserName   string
	ToUserName     string
	NewClientMsgID uint32
	CreateTime     uint32
	SvrNewMsgID    uint64
	IndexOfRequest uint32
}

// DownMediaItem 下载图片/视频/文件项
type DownMediaItem struct {
	AesKey   string
	FileURL  string
	FileType uint32
}

// DownVoiceItem 下载音频信息项
type DownVoiceItem struct {
	TotalLength  uint32
	NewMsgID     uint64
	ChatRoomName string
	MasterBufID  uint64
}

// VerifyUserItem 添加好友/验证好友/打招呼 项
type VerifyUserItem struct {
	OpType           uint32 // 1免验证发送请求, 2发送验证申请, 3通过好友验证
	FromType         byte   // 1来源QQ，2来源邮箱，3来源微信号，14群聊，15手机号，18附近的人，25漂流瓶，29摇一摇，30二维码，13来源通讯录
	VerifyContent    string // 验证信息
	VerifyUserTicket string // 通过验证UserTicket(同步到的)
	AntispamTicket   string // searchcontact请求返回
	UserValue        string // searchcontact请求返回
	ChatRoomUserName string // 通过群来添加好友 需要设置此值为群id
	NeedConfirm      uint32 // 是否确认
}

// StatusNotifyItem 状态通知项
type StatusNotifyItem struct {
	Code         uint32
	ToUserName   string
	ClientMsgID  string
	FunctionName string
	FunctionArg  string
}

type CheckFavCdnItem struct {
	DataId         string
	DataSourceId   string
	DataSourceType uint32
	FullMd5        string
	FullSize       uint32
	Head256Md5     string
	IsThumb        uint32
}

// SnsLocationInfo 朋友圈地址项
type SnsLocationInfo struct {
	City               string
	Longitude          string
	Latitude           string
	PoiName            string
	PoiAddress         string
	PoiScale           int32
	PoiInfoURL         string
	PoiClassifyID      string
	PoiClassifyType    uint32
	PoiClickableStatus uint32
}

// SnsMediaItem 朋友圈媒体项
type SnsMediaItem struct {
	EncKey        string
	EncValue      uint32
	ID            uint32
	Type          uint32
	Title         string
	Description   string
	Private       uint32
	UserData      string
	SubType       uint32
	URL           string
	URLType       string
	Thumb         string
	ThumType      string
	SizeWidth     string
	SizeHeight    string
	TotalSize     string
	VideoWidth    string
	VideoHeight   string
	MD5           string
	VideoMD5      string
	VideoDuration float64
}

// SnsPostItem 发送朋友圈需要的信息
type SnsPostItem struct {
	Xml           bool   //Content 是否纯xml
	ContentStyle  uint32 // 纯文字/图文/引用/视频
	Description   string
	ContentUrl    string
	Privacy       uint32           // 是否仅自己可见
	Content       string           // 文本内容
	MediaList     []*SnsMediaItem  // 图片/视频列表
	WithUserList  []string         // 提醒好友看列表
	GroupUserList []string         // 可见好友列表
	BlackList     []string         // 不可见好友列表
	LocationInfo  *SnsLocationInfo // 发送朋友圈的位置信息
}

// SnsObjectOpItem SnsObjectOpItem
type SnsObjectOpItem struct {
	SnsObjID string // 朋友圈ID
	OpType   uint32 // 操作码
	DataLen  uint32 // 其它数据长度
	Data     []byte // 其它数据
	Ext      uint32
}

// ReplyCommentItem 回覆的评论项
type ReplyCommentItem struct {
	UserName string // 评论的微信ID
	NickName string // 发表评论的昵称
	OpType   uint32 // 操作类型：评论/点赞
	Source   uint32 // source
}

// SnsCommentItem 朋友圈项：发表评论/点赞
type SnsCommentItem struct {
	OpType         uint32            // 操作类型：评论/点赞
	ItemID         uint64            // 朋友圈项ID
	ToUserName     string            // 好友微信ID
	Content        string            // 评论内容
	CreateTime     uint32            // 创建时间
	ReplyCommentID uint32            // 回复的评论ID
	ReplyItem      *ReplyCommentItem // 回覆的评论项
}

// GetLbsLifeListItem 获取地址列表项
type GetLbsLifeListItem struct {
	Opcode    uint32
	Buffer    []byte
	Longitude float32
	Latitude  float32
	KeyWord   string
}

// UploadVoiceItem 上传语音项
type UploadVoiceItem struct {
	ToUser      string
	Data        []byte
	VoiceLength uint32
	ClientMsgID string
	EndFlag     uint32
}

// LabelItem 标签项
type LabelItem struct {
	Name string
	ID   uint32
}

// UserLabelInfoItem 好友标签信息
type UserLabelInfoItem struct {
	UserName    string
	LabelIDList string
}

// ThumbItem 缩略图数据
type ThumbItem struct {
	Data   []byte
	Width  int32
	Height int32
}

// PackHeader 请求数据包头
type PackHeader struct {
	ReqData        []byte
	RetCode        int32
	Signature      byte
	HeadLength     byte
	CompressType   byte
	EncodeType     byte
	ServerVersion  uint32
	Uin            uint32
	Session        []byte
	SeqId          uint32
	URLID          uint32
	SrcLen         uint32
	ZipLen         uint32
	EncodeVersion  uint32
	HeadDeviceType byte
	CheckSum       uint32
	RunState       byte
	RqtCode        uint32
	EndFlag        byte
	Data           []byte
	HybridKeyVer   byte
}

func (p PackHeader) GetRetCode() int32 {
	return p.RetCode
}

func (p PackHeader) CheckSessionOut() bool {
	return p.RetCode == MMErrSessionTimeOut || p.RetCode == MMRequestRetSessionTimeOut
}
