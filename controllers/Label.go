package controllers

import (
	"encoding/json"
	"fmt"
	"wechatdll/models"
	"wechatdll/models/Label"
)

// 标签模块
type LabelController struct {
	BaseController
}

// @Summary 获取标签列表
// @Param	wxid			query 	string	true		"请输入登陆成功的wxid"
// @Failure 200
// @router /GetList [post]
func (c *LabelController) GetList() {
	wxid := c.GetString("wxid")
	WXDATA := Label.GetList(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 添加标签
// @Param	body			body	Label.AddParam	true		""
// @Failure 200
// @router /Add [post]
func (c *LabelController) Add() {
	var Data Label.AddParam
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
	WXDATA := Label.Add(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 更新标签列表
// @Param	body			body	Label.UpdateListParam	true		"ToWxid:多个请用,隔开"
// @Failure 200
// @router /UpdateList [post]
func (c *LabelController) UpdateList() {
	var Data Label.UpdateListParam
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
	WXDATA := Label.UpdateList(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 修改标签
// @Param	body			body	Label.UpdateNameParam	true		""
// @Failure 200
// @router /UpdateName [post]
func (c *LabelController) UpdateName() {
	var Data Label.UpdateNameParam
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
	WXDATA := Label.UpdateName(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 删除标签
// @Param	body			body	Label.DeleteParam	true		""
// @Failure 200
// @router /Delete [post]
func (c *LabelController) Delete() {
	var Data Label.DeleteParam
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
	WXDATA := Label.Delete(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}
