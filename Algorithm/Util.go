package Algorithm

import (
	"crypto/elliptic"
	"hash"
)

// 0x17000841 IOS 708
// 0x17000C2B IOS 712

//浏览器版本
//[]byte("Windows-QQBrowser")

var MmtlsShortHost = "extshort.weixin.qq.com" // "extshort.weixin.qq.com"	// "szshort.weixin.qq.com"
var MmtlsLongHost = "long.weixin.qq.com"
var MmtlsLongPort = 80

// 设备类型
var IPadDeviceType = "iPad iOS18.0.1"
var IPadModel = "iPad16,6"
var IPadOsVersion = "18.0.1"

var IPhoneDeviceType = "iPhone iOS18.0.1"
var IPhoneModel = "iPhone17,2"
var IPhoneOsVersion = "18.0.1"

var AndroidDeviceType = "android-34"
var AndroidManufacture = "HUAWEI Mate XT"
var AndroidModel = "GRL-AL10"
var AndroidRelease = "12"
var AndroidIncremental = "1"

var AndroidPadDeviceType = "pad-android-34"
var AndroidPadModel = "HUAWEI MRO-W00" //HUAWEI MatePad Pro
var AndroidPadOsVersion = "10"

var WinunifiedDeviceType = "UnifiedPCWindows 11 x86_64"
var WinunifiedModel = "ASUS"
var WinunifiedOsVersion = "11"

var WinDeviceType = "Windows 11 x64"
var WinModel = "ASUS"
var WinOsVersion = "11"

var CarDeviceType = "car-31"
var CarModel = "Xiaomi-M2012K11AC"
var CarOsVersion = "10"

var MacDeviceType = "iMac MacBookPro16,1 OSX OSX11.5.2 build(20G95)"
var MacModel = "iMac MacBookPro16,1"
var MacOsVersion = "11.5.2"

// 版本号
var IPadVersion = 0x18003727  //ipad
var IPadVersionx = 0x17012A21 //ipad绕过验证码0x17000523

var IPhoneVersion = 0x18003727 //62IPhone

var AndroidVersion = 0x2800373B  //A16Android
var AndroidVersion1 = 0x28003035 //A16Android848

var AndroidPadVersion = 0x2800373B //安卓平板
var AndroidPadVersionx = 0x27001032 //安卓平板绕过验证码

var WinVersion = 0x63090C11        //win
var WinUwpVersion = 0x620603C8     //winuwp绕过验证码
var WinUnifiedVersion = 0x64000115 //WinUnified 4.0.1.21

var CarVersion = 0x2100091B //车载

var MacVersion = 0x1308080B //mac

var RSA182_N = "D153E8A2B314D2110250A0A550DDACDCD77F5801F3D1CC21CB1B477E4F2DE8697D40F10265D066BE8200876BB7135EDC74CDBC7C4428064E0CDCBE1B6B92D93CEAD69EC27126DEBDE564AAE1519ACA836AA70487346C85931273E3AA9D24A721D0B854A7FCB9DED49EE03A44C189124FBEB8B17BB1DBE47A534637777D33EEC88802CD56D0C7683A796027474FEBF237FA5BF85C044ADC63885A70388CD3696D1F2E466EB6666EC8EFE1F91BC9353F8F0EAC67CC7B3281F819A17501E15D03291A2A189F6A35592130DE2FE5ED8E3ED59F65C488391E2D9557748D4065D00CBEA74EB8CA19867C65B3E57237BAA8BF0C0F79EBFC72E78AC29621C8AD61A2B79B"
var RSA182_E = "010001"

type HYBRID_STATUS int32

const (
	HYBRID_ENC HYBRID_STATUS = 0
	HYBRID_DEC HYBRID_STATUS = 1
)

type Client struct {
	PubKey     []byte
	Privkey    []byte
	InitPubKey []byte
	Externkey  []byte

	Version    int
	DeviceType string

	clientHash hash.Hash
	serverHash hash.Hash

	curve elliptic.Curve

	Status HYBRID_STATUS
}

type PacketHeader struct {
	PacketCryptType byte
	Flag            uint16
	RetCode         uint32
	UICrypt         uint32
	Uin             uint32
	Cookies         []byte
	Data            []byte
}

type PackData struct {
	Reqdata          []byte
	Cgi              int
	Uin              uint32
	Cookie           []byte
	ClientVersion    int
	Sessionkey       []byte
	EncryptType      uint8
	Loginecdhkey     []byte
	Clientsessionkey []byte
	Serversessionkey []byte
	UseCompress      bool
	MMtlsClose       bool
}
