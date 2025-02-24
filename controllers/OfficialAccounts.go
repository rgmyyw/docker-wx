package controllers

import (
	"encoding/json"
	"fmt"
	"wechatdll/models"
	"wechatdll/models/OfficialAccounts"
)

// 公众号模块
type OfficialAccountsController struct {
	BaseController
}

// @Summary 关注
// @Param	body		body 	OfficialAccounts.DefaultParam   true	""
// @Success 200
// @router /Follow [post]
func (c *OfficialAccountsController) Follow() {
	var ParamData OfficialAccounts.DefaultParam
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

	WXDATA := OfficialAccounts.Follow(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 取消关注
// @Param	body		body 	OfficialAccounts.DefaultParam   true	""
// @Success 200
// @router /Quit [post]
func (c *OfficialAccountsController) Quit() {
	var ParamData OfficialAccounts.DefaultParam
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

	WXDATA := OfficialAccounts.Quit(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary MpGetA8Key(获取文章key和uin)
// @Param	body		body 	OfficialAccounts.ReadParam   true	""
// @Success 200
// @router /MpGetA8Key [post]
func (c *OfficialAccountsController) MpGetA8Key() {
	var ParamData OfficialAccounts.ReadParam
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

	WXDATA := OfficialAccounts.MpGetA8Key(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 阅读文章,返回 分享、看一看、阅读数据
// @Param	body		body 	OfficialAccounts.ReadParam   true	""
// @Success 200
// @router /GetAppMsgExt [post]
func (c *OfficialAccountsController) GetAppMsgExt() {
	var ParamData OfficialAccounts.ReadParam
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

	WXDATA := OfficialAccounts.GetAppMsgExt(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 点赞文章,返回 分享、看一看、阅读数据
// @Param	body		body 	OfficialAccounts.ReadParam   true	""
// @Success 200
// @router /GetAppMsgExtLike [post]
func (c *OfficialAccountsController) GetAppMsgExtLike() {
	var ParamData OfficialAccounts.ReadParam
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

	WXDATA := OfficialAccounts.GetAppMsgExtLike(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary OauthAuthorize
// @Param	body		body 	OfficialAccounts.GetkeyParam   true	""
// @Success 200
// @router /OauthAuthorize [post]
func (c *OfficialAccountsController) OauthAuthorize() {
	var ParamData OfficialAccounts.GetkeyParam
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

	WXDATA := OfficialAccounts.OauthAuthorize(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary JSAPIPreVerify
// @Param	body		body 	OfficialAccounts.GetkeyParam   true	""
// @Success 200
// @router /JSAPIPreVerify [post]
func (c *OfficialAccountsController) JSAPIPreVerify() {
	var ParamData OfficialAccounts.GetkeyParam
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

	WXDATA := OfficialAccounts.JSAPIPreVerify(ParamData)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}
