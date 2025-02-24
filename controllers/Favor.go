package controllers

import (
	"encoding/json"
	"fmt"
	"wechatdll/models"
	"wechatdll/models/Favor"
)

// 收藏模块
type FavorController struct {
	BaseController
}

// @Summary 获取搜藏信息
// @Param	wxid		query 	string	true		"请输入登陆后的wxid"
// @Success 200
// @router /GetFavInfo [post]
func (c *FavorController) GetFavInfo() {
	wxid := c.GetString("wxid")
	WXDATA := Favor.GetFavInfo(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 同步收藏
// @Param	body			body 	Favor.SyncParam	true		"keybuf:第二次请求需要带上第一次返回的"
// @Success 200
// @router /Sync [post]
func (c *FavorController) Sync() {
	var reqdata Favor.SyncParam
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

	WXDATA := Favor.Sync(reqdata)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 读取收藏内容
// @Param	body			body 	Favor.GetFavItemParam	true		"FavId在同步收藏中获取"
// @Success 200
// @router /GetFavItem [post]
func (c *FavorController) GetFavItem() {
	var reqdata Favor.GetFavItemParam
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

	WXDATA := Favor.GetFavItem(reqdata)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 删除收藏
// @Param	body			body 	Favor.DelParam	true		"FavId在同步收藏中获取"
// @Success 200
// @router /Del [post]
func (c *FavorController) Del() {
	var reqdata Favor.DelParam
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

	WXDATA := Favor.Del(reqdata)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}
