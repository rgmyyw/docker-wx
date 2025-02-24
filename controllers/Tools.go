package controllers

import (
	"encoding/json"
	"fmt"
	"wechatdll/Cilent/mm"
	"wechatdll/models"
	"wechatdll/models/Tools"
)

// 工具箱模块
type ToolsController struct {
	BaseController
}

// @Summary 设置/删除代理IP
// @Param	body		body 	Tools.SetProxyParam   true	"删除代理ip时直接留空即可"
// @Success 200
// @router /setproxy [post]
func (c *ToolsController) SetProxy() {
	var ParamData Tools.SetProxyParam
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

	WXDATA := Tools.SetProxy(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary GetA8Key
// @Param	body		body 	Tools.GetA8KeyParam   true	"OpCode == 2 Scene == 4 CodeType == 19 CodeVersion == 5 以上是默认参数,如有需求自行修改"
// @Success 200
// @router /GetA8Key [post]
func (c *ToolsController) GetA8Key() {
	var ParamData Tools.GetA8KeyParam
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

	WXDATA := Tools.GetA8Key(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 高清图片下载
// @Param	body		body 	Tools.DownloadParam    true	"DataLen == 图片大小, xml中获取"
// @Success 200
// @router /DownloadImg [post]
func (c *ToolsController) DownloadImg() {
	var ParamData Tools.DownloadParam
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

	WXDATA := Tools.DownloadImg(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取CDN服务器dns信息
// @Param wxid		query 	string	true		"请输入登录后的wxid"
// @Success 200
// @router /GetCdnDns [post]
func (c *ToolsController) GetCdnDns() {
	wxid := c.GetString("wxid")
	WXDATA := Tools.GetCdnDns(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary CDN下载高清图片
// @Param	body		body 	Tools.CdnDownloadImageParam
// @Success 200
// @router /CdnDownloadImage [post]
func (c *ToolsController) CdnDownloadImage() {
	var ParamData Tools.CdnDownloadImageParam
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

	WXDATA := Tools.GetCdnDns(ParamData.Wxid)
	if !WXDATA.Success {
		c.Data["json"] = &WXDATA
		c.ServeJSON()
	}

	// 连接tcp
	dnsResponse := WXDATA.Data.(mm.GetCDNDnsResponse)

	ResData := Tools.CdnDownloadImg(ParamData, dnsResponse)

	c.Data["json"] = &ResData
	c.ServeJSON()
}

// @Summary 视频下载
// @Param	body		body 	Tools.DownloadParam    true	"DataLen == 视频大小, xml中获取"
// @Success 200
// @router /DownloadVideo [post]
func (c *ToolsController) DownloadVideo() {
	var ParamData Tools.DownloadParam
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

	WXDATA := Tools.DownloadVideo(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 文件下载
// @Param	body		body 	Tools.DownloadAppAttachParam    true	"DataLen == 文件大小, xml中获取"
// @Success 200
// @router /DownloadFile [post]
func (c *ToolsController) DownloadFile() {
	var ParamData Tools.DownloadAppAttachParam
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

	WXDATA := Tools.DownloadAppAttach(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 语音下载
// @Param	body		body 	Tools.DownloadVoiceParam    true	"注意参数"
// @Success 200
// @router /DownloadVoice [post]
func (c *ToolsController) DownloadVoice() {
	var ParamData Tools.DownloadVoiceParam
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

	WXDATA := Tools.DownloadVoice(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary GetBoundHardDevices
// @Param	wxid		query 	string	true		"请输入登录后的wxid"
// @Success 200
// @router /GetBoundHardDevices [post]
func (c *ToolsController) GetBoundHardDevices() {
	wxid := c.GetString("wxid")
	WXDATA := Tools.GetBoundHardDevices(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

/*// @Summary GetCert
// @Param	body		body 	Tools.GetCertParam    true	"Version == Rsa版本"
// @Success 200
// @router /GetCert [post]
func (c *ToolsController) GetCert() {
	var ParamData Tools.GetCertParam
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

	WXDATA := Tools.GetCert(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}*/


// @Summary 生成支付二维码
// @Param	wxid		query 	string	true		"请输入登录后的wxid"
// @Success 200
// @router /GeneratePayQCode [get]
func (c *ToolsController) GeneratePayQCode() {
	wxid := c.GetString("wxid")
	WXDATA := Tools.GeneratePayQCode(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 第三方APP授权
// @Param	body		body 	Tools.ThirdAppGrantParam    true	"注意参数"
// @Success 200
// @router /ThirdAppGrant [post]
func (c *ToolsController) ThirdAppGrant() {
	var ParamData Tools.ThirdAppGrantParam
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
	WXDATA := Tools.ThirdAppGrant(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}


// @Summary 获取余额以及银行卡信息
// @Param	wxid		query 	string	true		"请输入登录后的wxid"
// @Success 200
// @router /GetBandCardList [POST]
func (c *ToolsController) GetBandCardList() {
	wxid := c.GetString("wxid")
	WXDATA := Tools.GetBandCardList(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 修改微信步数
// @Param	body		body 	Tools.SetStepParam   true	"步数，最高支持98000"
// @Success 200
// @router /setproxy [post]
func (c *ToolsController) UpdateStepNumberApi() {
	var ParamData Tools.SetStepParam
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

	WXDATA := Tools.UpdateStepNumberApi(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary OauthSdkApp
// @Param	body			body	Tools.OauthSdkApp	 true		""
// @Failure 200
// @router /OauthSdkApp [post]
func (c *ToolsController) OauthSdkApp() {
	var Data Tools.OauthSdkAppParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &Data)
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
	WXDATA := Tools.OauthSdkApp(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}