package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/lib"
	"wechatdll/models"
	"wechatdll/models/Login"

	"github.com/bitly/go-simplejson"
)

// 登陆模块 支持二次 唤醒 62数据登陆(注意：代理必须使用SOCKS)
type LoginController struct {
	BaseController
}

// @Summary 获取二维码(iPad)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQR [post]
func (c *LoginController) LoginGetQR() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "iPad"
	}

	WXDATA := Login.GetQRCODE(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取二维码(iPad-绕过验证码)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRx [post]
func (c *LoginController) LoginGetQRx() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "iPad"
	}

	WXDATA := Login.GetQRCODEx(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取二维码(安卓Pad)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRPad [post]
func (c *LoginController) LoginGetQRPad() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "HUAWEI MatePad Pro"
	}

	WXDATA := Login.GetQRCODEPad(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取二维码(新版云函数)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRPad [post]
func (c *LoginController) LoginGetQRPad1() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "HUAWEI MatePad Pro"
	}

	WXDATA := Login.GetQRCODEPad1(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取二维码(安卓Pad-绕过验证码)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRPadx [post]
func (c *LoginController) LoginGetQRPadx() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "HUAWEI MatePad Pro"
	}

	WXDATA := Login.GetQRCODEPadx(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取二维码(安卓Pad-ppmt专用)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRPad [post]
func (c *LoginController) LoginGetQRPadppmt() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "Xiaomi-M2012K11AC"
	}

	WXDATA := Login.GetQRCODEPadppmt(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取二维码(Windows)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRWin [post]
func (c *LoginController) LoginGetQRWin() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "DESKTOP-P0QLAW8"
	}

	WXDATA := Login.GetQRCODEWin(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取二维码(WindowsUwp-绕过验证码)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRWinUwp [post]
func (c *LoginController) LoginGetQRWinUwp() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "DESKTOP-P0QLAW8"
	}

	WXDATA := Login.GetQRCODEWinUwp(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}


// @Summary 获取二维码(WinUnified-统一PC版)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRWinUnified [post]
func (c *LoginController) LoginGetQRWinUnified() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "DESKTOP-P0QLAW8"
	}

	WXDATA := Login.GetQRCODEWinUnified(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}


// @Summary 获取二维码(Car)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRWatch [post]
func (c *LoginController) LoginGetQRCar() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "Xiaomi-M2012K11AC"
	}

	WXDATA := Login.GetQRCODECar(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取二维码(Mac)
// @Param	body		body 	Login.GetQRReq	true		"不使用代理请留空"
// @Success 200
// @router /GetQRMac [post]
func (c *LoginController) LoginGetQRMac() {
	var GetQR Login.GetQRReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &GetQR)
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
	//如果没有指定deviceId,生成设备ID
	if GetQR.DeviceID == "" || GetQR.DeviceID == "string" {
		GetQR.DeviceID = lib.CreateDeviceId(GetQR.DeviceID)
	}

	if GetQR.DeviceName == "" || GetQR.DeviceName == "string" {
		GetQR.DeviceName = "MacBook Pro"
	}

	WXDATA := Login.GetQRCODEMac(GetQR)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 检测二维码
// @Param	uuid		query 	string	true		"请输入取码时返回的UUID"
// @Success 200
// @router /CheckQR [post]
func (c *LoginController) LoginCheckQR() {
	uuid := c.GetString("uuid")
	WXDATA := Login.CheckUuid(uuid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}



// @Summary 二次登陆
// @Param	wxid			query 	string	true		"请输入登陆成功的wxid"
// @Failure 200
// @router /TwiceAutoAuth [post]
func (c *LoginController) LoginTwiceAutoAuth() {
	wxid := c.GetString("wxid")
	WXDATA := Login.Secautoauth(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 62登陆(账号或密码)
// @Param	body			body 	Login.Data62LoginReq	true		"不使用代理请留空"
// @Failure 200
// @router /62data [post]
func (c *LoginController) Data62Login() {
	var reqdata Login.Data62LoginReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)

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

	WXDATA := Login.Data62(reqdata, Algorithm.MmtlsShortHost)

	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 62登陆(账号或密码), 并申请使用SMS验证
// @Param	body			body 	Login.Data62LoginReq	true		"不使用代理请留空"
// @Failure 200
// @router /62dataSMSApply [post]
func (c *LoginController) Data62SMSApply() {
	var reqdata Login.Data62LoginReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)

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

	// 生成62随机数据
	if reqdata.Data62 == "" || reqdata.Data62 == "string" {
		deviceId := lib.CreateDeviceId(reqdata.Data62)
		reqdata.Data62 = lib.Get62Data(deviceId)
	}
	if reqdata.DeviceName == "" || reqdata.DeviceName == "string" {
		reqdata.DeviceName = "iPad"
	}
	// 使用62数据登录并自动滑块
	WXDATA := Login.Data62(reqdata, Algorithm.MmtlsShortHost)

	// 记录62
	WXDATA.Data62 = reqdata.Data62

	// 二次验证使用短信验证
	message, transed := WXDATA.Data.(mm.UnifyAuthResponse)
	if transed && strings.Index(message.GetBaseResponse().GetErrMsg().GetString_(), "&ticket=") >= 0 {
		checkUrl, againUrl, setCookie := Login.WechatSMS1(message.GetBaseResponse().GetErrMsg().GetString_(), comm.GenRealWeixinUA(), reqdata.Proxy)
		WXDATA = models.ResponseResult{
			Code:    0,
			Success: true,
			Message: "已申请短信验证",
			Data: &map[string]string{
				"CheckUrl": checkUrl,
				"AgainUrl": againUrl,
				"Cookie":   setCookie,
			},
			Data62: reqdata.Data62,
		}
	}

	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 62登陆(账号或密码), 重发验证码
// @Param	body			body 	Login.Data62SMSAgainReq	true		"不使用代理请留空"
// @Failure 200
// @router /62dataSMSAgain [post]
func (c *LoginController) Data62SMSAgain() {
	var reqdata Login.Data62SMSAgainReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)
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

	// 重发短信
	headers := &map[string]string{
		"Cookie": reqdata.Cookie,
	}
	res := comm.HttpGet1(reqdata.Url, headers, comm.GenRealWeixinUA(), reqdata.Proxy)
	resJson, err := simplejson.NewJson([]byte(res))
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
	title, err := resJson.Get("resultData").Get("title").String()
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
	WXDATA := models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "",
		Data:    title,
		Data62:  "",
	}
	c.Data["json"] = &WXDATA
	c.ServeJSON()
	return
}

// @Summary 62登陆(账号或密码), 短信验证
// @Param	body			body 	Login.Data62SMSVerifyReq	true		"不使用代理请留空"
// @Failure 200
// @router /62dataSMSVerify [post]
func (c *LoginController) Data62SMSVerify() {
	var reqdata Login.Data62SMSVerifyReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)
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

	// 验证短信
	verifyUrl := strings.Replace(reqdata.Url, "[[[verifycode]]]", reqdata.Sms, -1)
	verifyUrl = strings.Replace(verifyUrl, "[[[currentMilliseStamp]]]", string(time.Now().Unix()), -1)
	headers := &map[string]string{
		"Cookie": reqdata.Cookie,
	}
	res := comm.HttpGet1(verifyUrl, headers, comm.GenRealWeixinUA(), reqdata.Proxy)
	resJson, err := simplejson.NewJson([]byte(res))
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
	title, err := resJson.Get("resultData").Get("title").String()
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
	WXDATA := models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "",
		Data:    title,
		Data62:  "",
	}
	c.Data["json"] = &WXDATA
	c.ServeJSON()
	return
}

// @Summary 62登陆(账号或密码), 并申请使用二维码验证
// @Param	body			body 	Login.Data62LoginReq	true		"不使用代理请留空"
// @Failure 200
// @router /62dataQRCodeApply [post]
func (c *LoginController) Data62QRCodeApply() {
	var reqdata Login.Data62LoginReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)

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

	// 生成62随机数据
	if reqdata.Data62 == "" || reqdata.Data62 == "string" {
		deviceId := lib.CreateDeviceId(reqdata.Data62)
		reqdata.Data62 = lib.Get62Data(deviceId)
	}
	if reqdata.DeviceName == "" || reqdata.DeviceName == "string" {
		reqdata.DeviceName = "iPad"
	}
	// 使用62数据登录并自动滑块
	WXDATA := Login.Data62(reqdata, Algorithm.MmtlsShortHost)

	// 记录62
	WXDATA.Data62 = reqdata.Data62

	// 二次验证使用短信验证
	message, transed := WXDATA.Data.(mm.UnifyAuthResponse)
	if transed && strings.Index(message.GetBaseResponse().GetErrMsg().GetString_(), "&ticket=") >= 0 {
		qrUrl, checkUrl := Login.WeChatQrCode1(message.GetBaseResponse().GetErrMsg().GetString_(), comm.GenRealWeixinUA(), reqdata.Proxy)
		WXDATA = models.ResponseResult{
			Code:    0,
			Success: true,
			Message: "已申请短信验证",
			Data: &map[string]string{
				"QrUrl":    qrUrl,
				"CheckUrl": checkUrl,
			},
			Data62: reqdata.Data62,
		}
	}

	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 62登陆(账号或密码), 二维码验证校验
// @Param	body			body 	Login.Data62SMSVerifyReq	true		"不使用代理请留空"
// @Failure 200
// @router /62dataSMSVerify [post]
func (c *LoginController) Data62QRCodeVerify() {
	var reqdata Login.Data62QRCodeVerifyReq
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)
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

	// 验证短信
	verifyUrl := reqdata.Url
	verifyUrl = strings.Replace(verifyUrl, "[[[currentMilliseStamp]]]", string(time.Now().Unix()), -1)
	res := comm.HttpGet1(verifyUrl, nil, comm.GenRealWeixinUA(), reqdata.Proxy)
	WXDATA := models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "",
		Data:    res,
		Data62:  "",
	}
	c.Data["json"] = &WXDATA
	c.ServeJSON()
	return
}

// @Summary A16登陆(账号或密码) - android == 8.0.50
// @Param	body			body 	Login.A16LoginParam	true		"不使用代理请留空"
// @Failure 200
// @router /A16Data [post]
func (c *LoginController) A16Data() {
	var reqdata Login.A16LoginParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)

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

	WXDATA := Login.AndroidA16Login(reqdata, Algorithm.MmtlsShortHost)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary A16登陆(账号或密码) - android == 新版云函数
// @Param	body			body 	Login.A16LoginParam	true		"不使用代理请留空"
// @Failure 200
// @router /A16Data848 [post]
func (c *LoginController) A16Data1() {
	var reqdata Login.A16LoginParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)

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

	WXDATA := Login.AndroidA16Login1(reqdata, Algorithm.MmtlsShortHost)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 心跳包
// @Param	wxid			query 	string	true		"请输入登陆成功的wxid"
// @Success 200
// @router /HeartBeat [post]
func (c *LoginController) HeartBeat() {
	wxid := c.GetString("wxid")
	WXDATA := Login.HeartBeat(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 心跳包
// @Param	wxid			query 	string	true		"请输入登陆成功的wxid"
// @Success 200
// @router /HeartBeatLong [post]
func (c *LoginController) HeartBeatLong() {
	wxid := c.GetString("wxid")
	WXDATA := Login.HeartBeatLong(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 初始化
// @Param	wxid			query 	string	true		"请输入登陆成功的wxid"
// @Param	MaxSynckey		query 	string	false		"二次同步需要带入"
// @Param	CurrentSynckey	query 	string	false		"二次同步需要带入"
// @Success 200
// @router /Newinit [post]
func (c *LoginController) Newinit() {
	wxid := c.GetString("wxid")
	MaxSynckey := c.GetString("MaxSynckey")
	CurrentSynckey := c.GetString("CurrentSynckey")
	WXDATA := Login.Newinit(wxid, MaxSynckey, CurrentSynckey)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 唤醒登陆(只限扫码登录)
// @Param	wxid		query 	string	true		"请输入登陆成功的wxid"
// @Success 200
// @router /Awaken [post]
func (c *LoginController) LoginAwaken() {
	wxid := c.GetString("wxid")
	WXDATA := Login.AwakenLogin(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取登陆缓存信息
// @Param	wxid		query 	string	true		"请输入登陆成功的wxid"
// @Success 200
// @router /GetCacheInfo [post]
func (c *LoginController) GetCacheInfo() {
	wxid := c.GetString("wxid")
	WXDATA := Login.CacheInfo(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取62数据
// @Param	wxid		query 	string	true		"请输入登陆成功的wxid"
// @Success 200
// @router /Get62Data [post]
func (c *LoginController) Get62Data() {
	wxid := c.GetString("wxid")
	Data62 := Login.Get62Data(wxid)
	Result := models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    Data62,
	}
	c.Data["json"] = &Result
	c.ServeJSON()
	return
}
// @Summary 获取A16数据
// @Param	wxid		query 	string	true		"请输入登陆成功的wxid"
// @Success 200
// @router /GetA16Data [post]
func (c *LoginController) GetA16Data() {
	wxid := c.GetString("wxid")
	Data62 := Login.GetA16Data(wxid)
	Result := models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    Data62,
	}
	c.Data["json"] = &Result
	c.ServeJSON()
	return
}


// @Summary 退出登录
// @Param	wxid			query 	string	true		"请输入登陆成功的wxid"
// @Success 200
// @router /LogOut [post]
func (c *LoginController) LogOut() {
	wxid := c.GetString("wxid")
	WXDATA := Login.LogOut(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 新设备扫码登录
// @Param	body			body 	Login.ExtDeviceLoginConfirmParam	true		"URL == MAC iPad Windows 的微信二维码解析出来的url"
// @Success 200
// @router /ExtDeviceLoginConfirmGet [post]
func (c *LoginController) ExtDeviceLoginConfirmGet() {
	var reqdata Login.ExtDeviceLoginConfirmParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)

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

	WXDATA := Login.ExtDeviceLoginConfirmGet(reqdata)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 新设备扫码确认登录
// @Param	body			body 	Login.ExtDeviceLoginConfirmParam	true		"URL == MAC iPad Windows 的微信二维码解析出来的url"
// @Success 200
// @router /ExtDeviceLoginConfirmOk [post]
func (c *LoginController) ExtDeviceLoginConfirmOk() {
	var reqdata Login.ExtDeviceLoginConfirmParam
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &reqdata)

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

	WXDATA := Login.ExtDeviceLoginConfirmOk(reqdata)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}
