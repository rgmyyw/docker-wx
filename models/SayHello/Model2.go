package SayHello

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"wechatdll/bts"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Friend"
)

type Model2Param struct {
	Wxid        string
	ToUserName  string
	Scene       int
	Content     string
	FromScene   uint32
	SearchScene uint32
}

//通道3 15 打招呼
func Model2(Data Model2Param) models.ResponseResult {
	_, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//不管如何将结果写入
	T := time.Now().Format("2006-01-02")

	logFile, err := os.OpenFile(fmt.Sprintf("log/SayHello_Model2_%v.log", T), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	loger := log.New(logFile, "", log.Ldate|log.Ltime)

	if Data.ToUserName == "" {
		loger.Printf("%v，%v，%v，%v，%v，%v，%v", Data.Wxid, Data.ToUserName, "null", "null", Data.Scene, Data.Content, "请输入微信号或者手机号")
		defer logFile.Close()
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "请输入微信号或者手机号",
			Data:    nil,
		}
	}

	//搜索V1V2
	S := Friend.Search(Friend.SearchParam{
		Wxid:        Data.Wxid,
		ToUserName:  Data.ToUserName,
		FromScene:   Data.FromScene,
		SearchScene: Data.SearchScene,
	})

	Search := bts.SearchContactResponse(S.Data)
	if Search.BaseResponse == nil {
		loger.Printf("%v，%v，%v，%v，%v，%v，%v", Data.Wxid, Data.ToUserName, "null", "null", Data.Scene, Data.Content, "离线")
		defer logFile.Close()
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "离线",
			Data:    nil,
		}
	}

	if *Search.BaseResponse.Ret != 0 {
		loger.Printf("%v，%v，%v，%v，%v，%v，%v", Data.Wxid, Data.ToUserName, "null", "null", Data.Scene, Data.Content, *Search.BaseResponse.ErrMsg.String_)
		defer logFile.Close()
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: *Search.BaseResponse.ErrMsg.String_,
			Data:    nil,
		}
	}

	if Search.AntispamTicket == nil {
		loger.Printf("%v，%v，%v，%v，%v，%v，%v", Data.Wxid, Data.ToUserName, *Search.UserName.String_, "null", Data.Scene, Data.Content, "未发现V4的存在")
		defer logFile.Close()
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "未发现V4的存在",
			Data:    nil,
		}
	}

	//开始提交验证
	Verify := Friend.SendRequest(Friend.SendRequestParam{
		Wxid:          Data.Wxid,
		V1:            *Search.UserName.String_,
		V2:            *Search.AntispamTicket,
		Opcode:        2,
		Scene:         Data.Scene,
		VerifyContent: Data.Content,
	})

	L, _ := json.Marshal(Verify.Data)
	loger.Printf("%v，%v，%v，%v，%v，%v，%v", Data.Wxid, Data.ToUserName, *Search.UserName.String_, *Search.AntispamTicket, Data.Scene, Data.Content, string(L))
	defer logFile.Close()

	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    Verify.Data,
	}
}
