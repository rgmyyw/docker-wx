package controllers

import (
	"encoding/json"
	"fmt"
	"wechatdll/models"
	"wechatdll/models/QWContact"
)

// 企业联系人操作
type QWContactController struct {
	BaseController
}

// @Tags QWContact
// @Summary QWAddContact
// @Param	body			body	QWContact.QWAddContact	 true		" "
// @Failure 200
// @router /QWContact/QWAddContact [post]
func (c *QWContactController) QWAddContact() {
	var Data QWContact.QWAddContactParam
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
	WXDATA := QWContact.QWAddContact(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary QWApplyAddContact
// @Param	body			body	QWContact.QWApplyAddContactParam	 true		""
// @Failure 200
// @router /QWApplyAddContact [post]
func (c *QWContactController) QWApplyAddContact() {
	var Data QWContact.QWApplyAddContactParam
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
	WXDATA := QWContact.QWApplyAddContact(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary SearchQWContact
// @Param	body			body	QWContact.AddWxAppRecordParam	 true		""
// @Failure 200
// @router /SearchQWContact [post]
func (c *QWContactController) SearchQWContact() {
	var Data QWContact.AddWxAppRecordParam
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
	WXDATA := QWContact.SearchQWContact(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}
