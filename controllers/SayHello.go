package controllers

import (
	"encoding/json"
	"fmt"
	"wechatdll/models"
	"wechatdll/models/SayHello"
)

// 打招呼模块
type SayHelloController struct {
	BaseController
}

// @Summary 模式1-扫码
// @Param	body			body	SayHello.Model1Param	 true		"注意,请先执行1再执行2"
// @Failure 200
// @router /Modelv1 [post]
func (c *SayHelloController) ModelV1() {
	var Data SayHello.Model1Param
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

	fmt.Println("============",Data)
	WXDATA := SayHello.Model1(Data)

	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 模式2-一键打招呼
// @Param	body			body	SayHello.Model2Param	 true		"Scene 招呼通道 FromScene SearchScene 搜索联系人场景"
// @Failure 200
// @router /Modelv2 [post]
func (c *SayHelloController) Modelv2() {
	var Data SayHello.Model2Param
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
	WXDATA := SayHello.Model2(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 模式3-v3\v4打招呼
// @Param	body			body	SayHello.SendRequestParam1	 true		"Scene 招呼通道 v3v4通道，v4可空"
// @Failure 200
// @router /Modelv2 [post]
func (c *SayHelloController) Modelv3() {
	var Data SayHello.SendRequestParam1
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
	if (Data.V4 == "" || Data.V4 == "string") {
		Data.V4 = ""
	}
	WXDATA := SayHello.Model3(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}
