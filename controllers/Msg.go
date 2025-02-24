package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"wechatdll/models"
	"wechatdll/models/Msg"
)

// 消息模块
type MsgController struct {
	BaseController
}

// @Summary 同步消息
// @Param	body		body 	Msg.SyncParam   true	"Scene填写0,Synckey留空"
// @Success 200
// @router /Sync [post]
func (c *MsgController) Sync() {
	var ParamData Msg.SyncParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.Sync(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

/*// @Summary 引用文本消息
// @Param	body		body 	Msg.Quote   true	"Fromusr == 被引用人Wxid, Displayname == 被引用人名称, NewMsgId == 引用人消息返回的NewMsgId, QuoteContent == 引用内容, MsgContent == 消息内容"
// @Success 200
// @router /Quote [post]
func (c *MsgController) Quote() {
	var ParamData Msg.Quote
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	AppXml := `<appmsg appid=""  sdkver="0"><title>`+ParamData.MsgContent+`</title><des></des><action></action><type>57</type><showtype>0</showtype><soundtype>0</soundtype><mediatagname></mediatagname><messageext></messageext><messageaction></messageaction><content></content><contentattr>0</contentattr><url></url><lowurl></lowurl><dataurl></dataurl><lowdataurl></lowdataurl><songalbumurl></songalbumurl><songlyric></songlyric><appattach><totallen>0</totallen><attachid></attachid><emoticonmd5></emoticonmd5><fileext></fileext><cdnthumbaeskey></cdnthumbaeskey><aeskey></aeskey></appattach><extinfo></extinfo><sourceusername></sourceusername><sourcedisplayname></sourcedisplayname><thumburl></thumburl><md5></md5><statextstr></statextstr><directshare>0</directshare><refermsg><type>1</type><svrid>`+ParamData.NewMsgId+`</svrid><fromusr>`+ParamData.Fromusr+`</fromusr><chatusr>`+ParamData.Wxid+`</chatusr><displayname>`+ParamData.Displayname+`</displayname><content>`+ParamData.QuoteContent+`</content><msgsource>&lt;msgsource&gt;&lt;sequence_id&gt;`+ParamData.MsgSeq+`&lt;/sequence_id&gt;&lt;/msgsource&gt;</msgsource></refermsg></appmsg><frsername></fromusername>`

	fmt.Println(AppXml)

	WXDATA := Msg.SendAppMsg(Msg.SendAppMsgParam{
		Wxid:   ParamData.Wxid,
		ToWxid: ParamData.ToWxid,
		Xml:    AppXml,
		Type:   57,
	})
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}*/

// @Summary 发送App消息
// @Param	body		body 	Msg.SendAppMsgParam   true	"Type请根据场景设置,xml请自行构造"
// @Success 200
// @router /SendApp [post]
func (c *MsgController) SendApp() {
	var ParamData Msg.SendAppMsgParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendAppMsg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送分享链接消息
// @Param	body		body 	Msg.SendAppMsgParam   true	"Type==类型 Desc==描述 Xml==发送xml内容 ToWxid==接受者"
// @Success 200
// @router /ShareLink [post]
func (c *MsgController) ShareLink() {
	var ParamData Msg.SendAppMsgParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendAppMsg(Msg.SendAppMsgParam{
		Wxid:   ParamData.Wxid,
		ToWxid: ParamData.ToWxid,
		Xml:    ParamData.Xml,
		Type:   ParamData.Type,
	})
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送文本消息
// @Param	body		body 	Msg.SendNewMsgParam   true	"Type请填写1 At == 群@,多个wxid请用,隔开"
// @Success 200
// @router /SendTxt [post]
func (c *MsgController) SendTxt() {
	var ParamData Msg.SendNewMsgParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendNewMsg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送图片
// @Param	body		body 	Msg.SendImageMsgParam   true	"请注意base64格式"
// @Success 200
// @router /UploadImg [post]
func (c *MsgController) UploadImg() {
	var ParamData Msg.SendImageMsgParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendImageMsg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送Emoji
// @Param	body		body 	Msg.SendEmojiParam   true	""
// @Success 200
// @router /SendEmoji [post]
func (c *MsgController) SendEmoji() {
	var ParamData Msg.SendEmojiParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendEmojiMsg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送Cdn视频(转发视频)
// @Param	body		body 	Msg.DefaultParam   true	"Content==消息xml"
// @Success 200
// @router /SendCDNVideo [post]
func (c *MsgController) SendCDNVideo() {
	var ParamData Msg.DefaultParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendCDNVideoMsg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送Cdn图片(转发图片)
// @Param	body		body 	Msg.DefaultParam   true	"Content==消息xml"
// @Success 200
// @router /SendCDNImg [post]
func (c *MsgController) SendCDNImg() {
	var ParamData Msg.DefaultParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendCDNImgMsg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送文件(转发,并非上传)
// @Param	body		body 	Msg.DefaultParam   true	"Content==收到文件消息xml"
// @Success 200
// @router /SendCDNFile [post]
func (c *MsgController) SendCDNFile() {
	var ParamData Msg.DefaultParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	XML := strings.Replace(ParamData.Content, "<msg>", "", -1)
	XML = strings.Replace(XML, "</msg>", "", -1)

	WXDATA := Msg.SendAppMsg(Msg.SendAppMsgParam{
		Wxid:   ParamData.Wxid,
		ToWxid: ParamData.ToWxid,
		Xml:    XML,
		Type:   6,
	})
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送语音
// @Param	body		body 	Msg.SendVoiceMessageParam   true	"Type： AMR = 0, MP3 = 2, SILK = 4, SPEEX = 1, WAVE = 3 VoiceTime ：音频长度 1000为一秒"
// @Success 200
// @router /SendVoice [post]
func (c *MsgController) SendVoice() {
	var ParamData Msg.SendVoiceMessageParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendVoiceMessage(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送视频
// @Param	body		body 	Msg.SendVideoMsgParam   true	""
// @Success 200
// @router /SendVideo [post]
func (c *MsgController) SendVideo() {
	var ParamData Msg.SendVideoMsgParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendVideoMsg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 分享名片
// @Param	body		body 	Msg.ShareCardParam   true	"ToWxid==接收的微信ID CardWxId==名片wxid CardNickName==名片昵称 CardAlias==名片别名 "
// @Success 200
// @router /ShareCard [post]
func (c *MsgController) ShareCard() {
	var ParamData Msg.ShareCardParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendNewMsg(Msg.SendNewMsgParam{
		Wxid:    ParamData.Wxid,
		ToWxid:  ParamData.ToWxid,
		Content: fmt.Sprintf("<msg username=\"%v\" nickname=\"%v\" fullpy=\"\" shortpy=\"\" alias=\"%v\" imagestatus=\"3\" scene=\"17\" province=\"\" city=\"\" sign=\"\" sex=\"1\" certflag=\"0\" certinfo=\"\" brandIconUrl=\"\" brandHomeUrl=\"\" brandSubscriptConfigUrl=\"\" brandFlags=\"0\" regionCode=\"CN\" ></msg>", ParamData.CardWxId, ParamData.CardNickName, ParamData.CardAlias),
		Type:    42,
	})
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 分享位置
// @Param	body		body 	Msg.ShareLocationParam   true	" "
// @Success 200
// @router /ShareLocation [post]
func (c *MsgController) ShareLocation() {
	var ParamData Msg.ShareLocationParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.SendNewMsg(Msg.SendNewMsgParam{
		Wxid:    ParamData.Wxid,
		ToWxid:  ParamData.ToWxid,
		Content: fmt.Sprintf("<msg><location x=\"%v\" y=\"%v\" scale=\"%v\" label=\"%v\" poiname=\"%v\" maptype=\"roadmap\" infourl=\"\" fromusername=\"\" poiid=\"City\" /></msg>", ParamData.X, ParamData.Y, ParamData.Scale, ParamData.Label, ParamData.Poiname),
		Type:    42,
	})
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 撤回消息
// @Param	body		body 	Msg.RevokeMsgParam   true	"请注意参数"
// @Success 200
// @router /Revoke [post]
func (c *MsgController) Revoke() {
	var ParamData Msg.RevokeMsgParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}

	WXDATA := Msg.RevokeMsg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送分享视频消息
// @Param	body		body 	Msg.ShareVideoMsgParam   true	"xml：微信返回的视频xml"
// @Success 200
// @router /ShareVideo [post]
func (c *MsgController) ShareVideo() {
	var ParamData Msg.ShareVideoMsgParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &ParamData)
	if err != nil {
		Result := models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("系统异常：%v", err.Error()),
			Data:    nil,
		}
		c.Data["json"] = &Result
		c.ServeJSON()
		return
	}
	WXDATA := Msg.ShareVideoMsg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}
