package controllers

import (
	"encoding/json"
	"fmt"
	"wechatdll/models"
	"wechatdll/models/Friend"
)

// 朋友模块
type FriendController struct {
	BaseController
}

// @Summary 搜索联系人
// @Param	body			body	Friend.SearchParam	 true		"爆粉情况下特殊通道请自行填写,默认时FromScene=0,SearchScene=1"
// @Failure 200
// @router /Search [post]
func (c *FriendController) Search() {
	var Data Friend.SearchParam
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
	WXDATA := Friend.Search(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 添加联系人(发送好友请求)
// @Param	body			body	Friend.SendRequestParam	 true		"V1 V2是必填项"
// @Failure 200
// @router /SendRequest [post]
func (c *FriendController) SendRequest() {
	var Data Friend.SendRequestParam
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
	WXDATA := Friend.SendRequest(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 通过好友请求
// @Param	body			body	Friend.PassVerifyParam	 true		"Scene：代表来源,请在消息中的xml中获取"
// @Failure 200
// @router /PassVerify [post]
func (c *FriendController) PassVerify() {
	var Data Friend.PassVerifyParam
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
	WXDATA := Friend.PassVerify(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 上传通讯录
// @Param	body			body	Friend.UploadParam	 true		"PhoneNo多个手机号请用,隔开   CurrentPhoneNo自己的手机号  Opcode == 1上传 2删除"
// @Failure 200
// @router /Upload [post]
func (c *FriendController) Upload() {
	var Data Friend.UploadParam
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
	WXDATA := Friend.Upload(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取手机通讯录
// @Param	wxid		query 	string	true		"请输入登陆后的wxid"
// @Failure 200
// @router /GetMFriend [post]
func (c *FriendController) GetMFriend() {
	wxid := c.GetString("wxid")
	WXDATA := Friend.GetMFriend(wxid)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取通讯录好友
// @Param	body			body	Friend.GetContractListparameter	 true		"CurrentWxcontactSeq和CurrentChatRoomContactSeq没有的情况下请填写0"
// @Failure 200
// @router /GetContractList [post]
func (c *FriendController) GetContractList() {
	var Data Friend.GetContractListparameter
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
	WXDATA := Friend.GetContractList(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 获取通讯录好友详情
// @Param	body			body	Friend.GetContractDetailparameter	 true		"多个微信请用,隔开(最多20个),ChatRoom请留空"
// @Failure 200
// @router /GetContractDetail [post]
func (c *FriendController) GetContractDetail() {
	var Data Friend.GetContractDetailparameter
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
	WXDATA := Friend.GetContractDetail(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 设置好友备注
// @Param	body			body	Friend.SetRemarksParam	 true		""
// @Failure 200
// @router /SetRemarks [post]
func (c *FriendController) SetRemarks() {
	var Data Friend.SetRemarksParam
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
	WXDATA := Friend.SetRemarks(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 删除好友
// @Param	body			body	Friend.DefaultParam	 true		""
// @Failure 200
// @router /Delete [post]
func (c *FriendController) Delete() {
	var Data Friend.DefaultParam
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
	WXDATA := Friend.Delete(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 添加/移除黑名单
// @Param	body			body	Friend.BlacklistParam	 true		"Val == 15添加  7移除"
// @Failure 200
// @router /Blacklist [post]
func (c *FriendController) Blacklist() {
	var Data Friend.BlacklistParam
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
	WXDATA := Friend.Blacklist(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 附近人
// @Param	body			body	Friend.LbsFindParam	 true		"OpCode == 1"
// @Failure 200
// @router /LbsFind [post]
func (c *FriendController) LbsFind() {
	var Data Friend.LbsFindParam
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
	WXDATA := Friend.LbsFind(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

// @Summary 查询好友状态
// @Param	body			body	Friend.GetFriendstate1Param	 true		"OpCode == 1"
// @Failure 200
// @router /GetFriendstate [post]
func (c *FriendController) GetFriendstate() {
	var Data Friend.FriendRelationParam
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
	WXDATA := Friend.FriendRelation(Data)
	c.Data["json"] = &WXDATA
	c.ServeJSON()
}

