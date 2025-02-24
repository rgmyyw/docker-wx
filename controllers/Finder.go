package controllers

import (
	"wechatdll/models/Finder"
)

// 视频号模块
type FinderController struct {
	BaseController
}

// @Summary 用户中心
// @Param	wxid		query 	string	true		"请输登陆后的wxid"
// @Success 200
// @router /UserPrepare [post]
func (c *FinderController) UserPrepare() {
	wxid := c.GetString("wxid")
	WXDATA := Finder.UserPrepare(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}
