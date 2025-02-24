package controllers

import (
	"encoding/json"
	"fmt"
	"wechatdll/models"
	"wechatdll/models/Wxapp"
)

// 微信小程序模块
type WxappController struct {
	BaseController
}

// @Tags Wxapp
// @Summary 授权小程序(返回授权后的code)
// @Param	body			body	Wxapp.DefaultParam	 true		""
// @Failure 200
// @router /JSLogin [post]
func (c *WxappController) JSLogin() {
	var Data Wxapp.DefaultParam
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
	WXDATA := Wxapp.JSLogin(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Tags Wxapp
// @Summary 小程序操作
// @Param	body			body	Wxapp.JSOperateWxParam	 true		""
// @Failure 200
// @router /JSOperateWxData [post]
func (c *WxappController) JSOperateWxData() {
	var Data Wxapp.JSOperateWxParam
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
	WXDATA := Wxapp.JSOperateWx(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Tags Wxapp
// @Summary 获取小程序支付sessionid
// @Param	body			body	Wxapp.DefaultParam	 true		""
// @Failure 200
// @router /JSGetSessionid [post]
func (c *WxappController) JSGetSessionid() {
	var Data Wxapp.DefaultParam
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
	WXDATA := Wxapp.JSGetSessionid(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Tags Wxapp
// @Summary 小程序云函数
// @Param	body			body	Wxapp.CloudCallParam	 true		""
// @Failure 200
// @router /CloudCallFunction [post]
func (c *WxappController) CloudCallFunction() {
	var Data Wxapp.CloudCallParam
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
	WXDATA := Wxapp.CloudCallFunction(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Tags Wxapp
// @Summary 新增小程序记录
// @Param	body			body	Wxapp.AddWxAppRecordParam	 true		" "
// @Failure 200
// @router /Wxapp/AddWxAppRecord [post]
func (c *WxappController) AddWxAppRecord() {
	var Data Wxapp.AddWxAppRecordParam
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
	WXDATA := Wxapp.AddWxAppRecord(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Tags Wxapp
// @Summary 获取付小程序款二维码
// @Param	body			body	Wxapp.JSGetSessionidQRcode	 true		" "
// @Failure 200
// @router /Wxapp/JSGetSessionidQRcode [post]
func (c *WxappController) JSGetSessionidQRcode() {
	var Data Wxapp.SessionidQRParam
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
	WXDATA := Wxapp.JSGetSessionidQRcode(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Tags Wxapp
// @Summary 扫码授权登录app或网页
// @Param	body			body	Wxapp.QrcodeAuthLogin	 true		" "
// @Failure 200
// @router /Wxapp/QrcodeAuthLogin [post]
func (c *WxappController) QrcodeAuthLogin() {
	var Data Wxapp.QrcodeAuthLoginParam
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
	WXDATA := Wxapp.QrcodeAuthLogin(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 小程序绑定增加手机号
// @Param	body			body	Wxapp.CheckVerifyCodeData	 true		""
// @Failure 200
// @router /AddMobile [post]
func (c *WxappController) AddMobile() {
	var Data Wxapp.CheckVerifyCodeData
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
	if Data.VerifyCode == "" {
		datas := Wxapp.PostVerifyCodeParam{
			Appid:  Data.Appid,
			Mobile: Data.Mobile,
			Wxid:   Data.Wxid,
			Opcode: 0,
		}
		WXDATA := Wxapp.PostVerifyCode(datas)
		c.Data["json"] = &WXDATA
		c.ServeJSON()
	} else {
		WXDATA := Wxapp.AddMobile(Data)
		c.Data["json"] = &WXDATA
		c.ServeJSON()
	}

}

// @Summary 小程序删除手机号
// @Param	body			body	Wxapp.DelMobileData	 true		""
// @Failure 200
// @router /DelMobile [post]
func (c *WxappController) DelMobile() {
	var Data Wxapp.DelMobileData
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
	WXDATA := Wxapp.DelMobile(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()

}

// @Summary GetRandomAvatar
// @Param	body			body	Wxapp.DefaultParam	 true		""
// @Failure 200
// @router /GetRandomAvatar [post]
func (c *WxappController) GetRandomAvatar() {
	var Data Wxapp.DefaultParam
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
	WXDATA := Wxapp.GetRandomAvatar(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary UploadAvatarImg
// @Param	body			body	Wxapp.AddAvatarImgParam	 true		""
// @Failure 200
// @router /UploadAvatarImg [post]
func (c *WxappController) UploadAvatarImg() {
	var Data Wxapp.AddAvatarImgParam
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
	WXDATA := Wxapp.UploadAvatarImg(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary AddAvatar
// @Param	body			body	Wxapp.AddAvatarParam	 true		""
// @Failure 200
// @router /AddAvatar [post]
func (c *WxappController) AddAvatar() {
	var Data Wxapp.AddAvatarParam
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
	WXDATA := Wxapp.AddAvatar(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary GetUserOpenId
// @Param	body			body	Wxapp.GetUserOpenId	 true		""
// @Failure 200
// @router /GetUserOpenId [post]
func (c *WxappController) GetUserOpenId() {
	var Data Wxapp.GetUserOpenIdParam
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
	WXDATA := Wxapp.GetUserOpenId(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}



// @Summary GetAllMobile
// @Param	body			body	Wxapp.GetAllMobile	 true		""
// @Failure 200
// @router /GetAllMobile [post]
func (c *WxappController) GetAllMobile() {
	var Data Wxapp.JSOperateWxParam
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
	WXDATA := Wxapp.GetAllMobile(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}