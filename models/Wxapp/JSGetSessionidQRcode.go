package Wxapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"wechatdll/models"
)

type ResponseData struct {
	UUID string `json:"uuid"`
}

func JSGetSessionidQRcode(Data SessionidQRParam) models.ResponseResult {

	url := "https://open.weixin.qq.com/wxaruntime/getuuid?session_id=" + Data.Sessionid
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"appid":"%s","req_data":"{\"invokeData\":\"{\\\"miniprogramAppID\\\":\\\"%s\\\",\\\"args\\\":\\\"{\\\\\\\"adUxInfo\\\\\\\":\\\\\\\"\\\\\\\",\\\\\\\"grantMessageQuota\\\\\\\":true,\\\\\\\"signType\\\\\\\":\\\\\\\"MD5\\\\\\\",\\\\\\\"provider\\\\\\\":\\\\\\\"wxpay\\\\\\\",\\\\\\\"package\\\\\\\":\\\\\\\"%s\\\\\\\",\\\\\\\"timeStamp\\\\\\\":\\\\\\\"%s\\\\\\\",\\\\\\\"cookie\\\\\\\":\\\\\\\"busid=wxapp; appid=%s;; busid=wxapp; sessionid=;; busid=wxapp; scene=1001;; busid=wxapp; scene_note=\\\\\\\",\\\\\\\"nonceStr\\\\\\\":\\\\\\\"%s\\\\\\\",\\\\\\\"paySign\\\\\\\":\\\\\\\"%s\\\\\\\"}\\\",\\\"isBridgedJsApi\\\":true,\\\"name\\\":\\\"requestPayment\\\",\\\"jsapiType\\\":\\\"appservice\\\",\\\"transitiveData\\\":\\\"\\\"}\",\"pathType\":1,\"runtimeSessionId\":\"%s\",\"rumtimeAppid\":\"%s\",\"runtimeTicket\":\"Test\"}"}`, Data.Appid, Data.Appid, Data.Package, Data.timeStamp, Data.Appid, Data.nonceStr, Data.PaySign, Data.Sessionid, Data.Appid))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return models.ResponseResult{
			Code:    -1,
			Success: false,
			Message: fmt.Sprintf("创建请求失败：%v", err),
			Data:    nil,
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return models.ResponseResult{
			Code:    -2,
			Success: false,
			Message: fmt.Sprintf("执行请求失败：%v", err),
			Data:    nil,
		}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return models.ResponseResult{
			Code:    -3,
			Success: false,
			Message: fmt.Sprintf("读取响应失败：%v", err),
			Data:    nil,
		}
	}

	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return models.ResponseResult{
			Code:    -4,
			Success: false,
			Message: fmt.Sprintf("解析 JSON 失败：%v", err),
			Data:    nil,
		}
	}

	err = json.Unmarshal(body, &responseData)
	if responseData.UUID == "" {
		return models.ResponseResult{
			Code:    -4,
			Success: false,
			Message: fmt.Sprintf("获取UUID失败，请重试"),
			Data:    nil,
		}
	}

	newURL := "https://api.weixin.qq.com/wxaruntime/readqrcode?uuid=" + responseData.UUID
	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    newURL,
	}
}
