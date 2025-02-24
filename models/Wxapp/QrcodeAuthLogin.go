package Wxapp

import (
	"fmt"
	"net/http"
	"strings"
	"wechatdll/bts"
	"io"
	"encoding/json"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Tools"
)

func QrcodeAuthLogin(Data QrcodeAuthLoginParam) models.ResponseResult {
	_, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	url := "https://open.weixin.qq.com/connect/confirm?uuid=" + Data.Uuid


	a8key := Tools.GetA8Key(Tools.GetA8KeyParam{
		Wxid:        Data.Wxid,
		OpCode:      2,
		Scene:       4,
		CodeType:    19,
		CodeVersion: 5,
		ReqUrl:      url,
	})

	getA8key := bts.GetA8KeyResponse(a8key.Data)

	if getA8key.FullURL == nil {
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "请求失败",
			Data:    nil,
		}
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("POST", *getA8key.FullURL, strings.NewReader("s=1"))
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	


	req.Header.Set("Origin", "https://open.weixin.qq.com")
	resp, err := client.Do(req)

	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}
	defer resp.Body.Close()

	str := *getA8key.FullURL
	delimiter := "https://open.weixin.qq.com/connect/confirm?"
	rightText := GetRightText(str, delimiter)
	url = "https://open.weixin.qq.com/connect/confirm_reply?" + rightText + "&snsapi_login=on&allow=allow"
	//fmt.Println(url)
	// 判断第一个请求成功后执行第二个请求
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("创建第二个 GET 请求失败:", err)
			return models.ResponseResult{
				Code:    -9,
				Success: false,
				Message: "第二个请求创建失败",
				Data:    nil,
			}
		}

		// 发送第二个请求
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("发送第二个 GET 请求失败:", err)
			return models.ResponseResult{
				Code:    -9,
				Success: false,
				Message: "第二个请求发送失败",
				Data:    nil,
			}
		}
		defer resp.Body.Close()
	} else {
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "第一个请求未成功，不执行第二个请求",
			Data:    nil,
		}
	}



	url = "https://long.open.weixin.qq.com/connect/l/qrconnect?uuid=" + Data.Uuid + "&last=0&f=json"
	//fmt.Println(url)
	// 判断第二个请求成功后执行第三个请求

		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("创建第三个 GET 请求失败:", err)
			return models.ResponseResult{
				Code:    -9,
				Success: false,
				Message: "第三个请求创建失败",
				Data:    nil,
			}
		}

		req.Header.Set("User-Agent", "User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.54(0x18003639) NetType/4G Language/zh_CN")
		req.Header.Set("Origin", "")

		// 发送第二个请求
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("发送第三个 GET 请求失败:", err)
			return models.ResponseResult{
				Code:    -9,
				Success: false,
				Message: "第三个请求发送失败",
				Data:    nil,
			}
		}
		defer resp.Body.Close()


	


 
    body, err := io.ReadAll(resp.Body)
	
	type responseDataParam struct {
		Wx_errcode   int `json:"wx_errcode"`
		Wx_code string `json:"wx_code"`
	}


	var responseData responseDataParam

    if err := json.Unmarshal(body, &responseData); err != nil {

    }
	Wx_code := responseData.Wx_code

if  Wx_code == "" {
	Wx_errcode := responseData.Wx_errcode
	url = "https://long.open.weixin.qq.com/connect/l/qrconnect?uuid=" + Data.Uuid + "&last=" + string(Wx_errcode) + "&f=json"
	//fmt.Println(url)
	// 判断第二个请求成功后执行第三个请求
	
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("创建第三个 GET 请求失败:", err)
			return models.ResponseResult{
				Code:    -9,
				Success: false,
				Message: "第三个请求创建失败",
				Data:    nil,
			}
		}
	
		req.Header.Set("User-Agent", "User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 16_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.54(0x18003639) NetType/4G Language/zh_CN")
		req.Header.Set("Origin", "")
	
		// 发送第二个请求
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("发送第三个 GET 请求失败:", err)
			return models.ResponseResult{
				Code:    -9,
				Success: false,
				Message: "第三个请求发送失败",
				Data:    nil,
			}
		}
		defer resp.Body.Close()
	

	body, err = io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &responseData); err != nil {
	
	}

}
	
	

	


	//str = string(body)
	



	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    responseData,
	}

}
// GetRightText 从给定的字符串中获取右边的文本
func GetRightText(s string, delimiter string) string {
	index := strings.LastIndex(s, delimiter)
	if index == -1 {
		return s
	}
	return s[index+len(delimiter):]
}

