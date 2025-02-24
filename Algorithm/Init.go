package Algorithm

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
)

func (h *Client) Init(Model string) {

	h.curve = elliptic.P256()
	h.clientHash = sha256.New()
	h.serverHash = sha256.New()

	if Model == "IOS" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = IPadVersion
		h.DeviceType = IPadDeviceType
		h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")
	}

	if Model == "iPhone" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = IPhoneVersion
		h.DeviceType = IPhoneDeviceType
		h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")
	}

	if Model == "MAC" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = MacVersion
		h.DeviceType = MacDeviceType
		h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")
	}

	if Model == "AndroidPad" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = AndroidPadVersion
		h.DeviceType = AndroidPadDeviceType
		h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")
	}

	if Model == "Windows" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = WinVersion
		h.DeviceType = WinDeviceType
		h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")
	}

	if Model == "WindowsUwp" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = WinUwpVersion
		h.DeviceType = WinDeviceType
		h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")
	}

	if Model == "WinUnified" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = WinUnifiedVersion
		h.DeviceType = WinunifiedDeviceType
		h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")
	}

	if Model == "Android" {
		h.Status = HYBRID_ENC
		h.Version = AndroidVersion
		h.DeviceType = AndroidDeviceType
		h.InitPubKey, _ = hex.DecodeString("0495BC6E5C1331AD172D0F35B1792C3CE63F91572ABD2DD6DF6DAC2D70195C3F6627CCA60307305D8495A8C38B4416C75021E823B6C97DFFE79C14CB7C3AF8A586")
	}

	if Model == "Car" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = CarVersion
		h.DeviceType = CarDeviceType
		h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")
	}

}
func (h *Client) Init2(Model string) {

	h.curve = elliptic.P256()
	h.clientHash = sha256.New()
	h.serverHash = sha256.New()

	//MicroMessenger/
	//7.0.12(0x17000c21)--wxid_4ucwxdg2896t222
	//MicroMessenger/7.0.14(0x17000e2e) --- 自由改版本
	//MicroMessenger/7.0.15(0x17000f26)
	//MicroMessenger/7.0.17(0x17001124)   iphone  0X1700112a
	//7.0.18  0x17001231
	if Model == "IOS" {
		h.Privkey, h.PubKey = GetECDH415Key()
		h.Version = IPadVersion
		h.DeviceType = IPadDeviceType                                                                                                                                            //"iPad iOS13.5"// "iPad iOS13.5"
		h.InitPubKey, _ = hex.DecodeString("044bb81879aff459ca8f1db3d38eea5d789afaed14765a859a6f70bf06b663f37c6bd9e05c9f5def4ab796ca2c45b9d9a0f553ac8be51c0f60e087faee24d14510") //10003
		//h.InitPubKey, _ = hex.DecodeString("047ebe7604acf072b0ab0177ea551a7b72588f9b5d3801dfd7bb1bca8e33d1c3b8fa6e4e4026eb38d5bb365088a3d3167c83bdd0bbb46255f88a16ede6f7ab43b5")//10000
		//h.InitPubKey,_=hex.DecodeString("0495bc6e5c1331ad172d0f35b1792c3ce63f91572abd2dd6df6dac2d70195c3f6627cca60307305d8495a8c38b4416c75021e823b6c97dffe79c14cb7c3af8a586")
	}

	////安卓版本号 -  7.019  654315572
	if Model == "Android" {
		h.Status = HYBRID_ENC
		h.Version = AndroidVersion //0x27001636// 654317110
		h.DeviceType = "android-28"
		h.InitPubKey, _ = hex.DecodeString("0495bc6e5c1331ad172d0f35b1792c3ce63f91572abd2dd6df6dac2d70195c3f6627cca60307305d8495a8c38b4416c75021e823b6c97dffe79c14cb7c3af8a586")
	}

}
