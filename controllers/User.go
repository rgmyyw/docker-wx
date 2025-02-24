package controllers

import (
	"encoding/json"
	"fmt"
	"wechatdll/models"
	"wechatdll/models/User"
)

// 微信号管理模块
type UserController struct {
	BaseController
}

// @Summary 取个人信息
// @Param	wxid		query 	string	true		"请输入登陆后的wxid"
// @Success 200
// @router /GetContractProfile [post]
func (c *UserController) GetContractProfile() {
	wxid := c.GetString("wxid")
	WXDATA := User.GetContractProfile(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 取个人二维码
// @Param	body	body	User.GetQRCodeParam true		"Style == 二维码样式(请自行探索) 8默认"
// @Success 200
// @router /GetQRCode [post]
func (c *UserController) GetQRCode() {
	var Data User.GetQRCodeParam
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
	WXDATA := User.GetQRCode(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 绑定QQ
// @Param	body: {account: QQ账号,password: QQ密码}
// @Success 200
// @router /BindQQ [post]
func (c *UserController) BindQQ() {
	var Data User.BindQQParam
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
	WXDATA := User.BindQQ(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 隐私设置
// @Param	body	body	User.PrivacySettingsParam	 true		"核心参数请联系客服获取代码列表"
// @Success 200
// @router /PrivacySettings [post]
func (c *UserController) PrivacySettings() {
	var Data User.PrivacySettingsParam
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
	WXDATA := User.PrivacySettings(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 修改个人信息
// @Param	body	body	User.UpdateProfileParam	 true		"NickName ==名称  Sex == 性别（1:男 2：女） Country == 国家,例如：CH Province == 省份 例如:WuHan Signature == 个性签名"
// @Success 200
// @router /UpdateProfile [post]
func (c *UserController) UpdateProfile() {
	var Data User.UpdateProfileParam
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
	WXDATA := User.UpdateProfile(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 修改头像
// @Param	body	body	User.UploadHeadImageParam	 true		""
// @Success 200
// @router /UploadHeadImage [post]
func (c *UserController) UploadHeadImage() {
	var Data User.UploadHeadImageParam
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
	WXDATA := User.UploadHeadImage(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 验证密码
// @Param	body	body	User.NewVerifyPasswdParam	 true		""
// @Success 200
// @router /VerifyPasswd [post]
func (c *UserController) VerifyPasswd() {
	var Data User.NewVerifyPasswdParam
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
	WXDATA := User.NewVerifyPasswd(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 修改密码
// @Param	body	body	User.NewSetPasswdParam	 true		""
// @Success 200
// @router /SetPasswd [post]
func (c *UserController) SetPasswd() {
	var Data User.NewSetPasswdParam
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
	WXDATA := User.NewSetPasswd(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 设置微信号
// @Param	body	body	User.SetAlisaParam	 true		""
// @Success 200
// @router /SetAlisa [post]
func (c *UserController) SetAlisa() {
	var Data User.SetAlisaParam
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
	WXDATA := User.SetAlisa(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 绑定邮箱
// @Param	body	body	User.EmailParam	 true		""
// @Success 200
// @router /BindingEmail [post]
func (c *UserController) BindingEmail() {
	var Data User.EmailParam
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
	WXDATA := User.BindEmail(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 发送手机验证码
// @Param	body	body	User.SendVerifyMobileParam	 true		"Opcode == 场景(18代表绑手机号) Mobile == 格式：+8617399999999"
// @Success 200
// @router /SendVerifyMobile [post]
func (c *UserController) SendVerifyMobile() {
	var Data User.SendVerifyMobileParam
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
	WXDATA := User.SendVerifyMobile(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 换绑手机号
// @Param	body	body	User.BindMobileParam	 true		"Mobile == 格式：+8617399999999 Verifycode == 验证码请先通过(发送手机验证码)获取"
// @Success 200
// @router /BindingMobile [post]
func (c *UserController) BindingMobile() {
	var Data User.BindMobileParam
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
	WXDATA := User.BindMobile(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 登录设备管理
// @Param	wxid		query 	string	true		"请输入登陆后的wxid"
// @Success 200
// @router /GetSafetyInfo [post]
func (c *UserController) GetSafetyInfo() {
	wxid := c.GetString("wxid")
	WXDATA := User.GetSafetyInfo(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 删除登录设备
// @Param	body	body	User.DelSafetyInfoParam	 true		"UUID请在登录设备管理中获取"
// @Success 200
// @router /DelSafetyInfo [post]
func (c *UserController) DelSafetyInfo() {
	var Data User.DelSafetyInfoParam
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
	WXDATA := User.DelSafetyInfo(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary ReportMotion
// @Param	body	body	User.ReportMotionParam	 true		"具体用法请联系客服"
// @Success 200
// @router /ReportMotion [post]
func (c *UserController) ReviseMotion() {
	var Data User.ReportMotionParam
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
	WXDATA := User.ReportMotion(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}
