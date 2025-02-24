package Group

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
	"wechatdll/Cilent/mm"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Tools"
)

type ScanIntoGroupParam struct {
	Wxid string
	Url  string
}

func ScanIntoGroup(Data ScanIntoGroupParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	firstA8keyRes := Tools.GetA8Key(Tools.GetA8KeyParam{
		Wxid:         Data.Wxid,
		OpCode:       2,
		Scene:        4,
		CodeType:     19,
		CodeVersion:  8,
		ReqUrl:       Data.Url,
		CookieBase64: "",
		NetType:      "",
	})
	var a8KeyRes mm.GetA8KeyResp
	if firstA8keyRes.Success && firstA8keyRes.Data != nil {
		a8KeyRes = firstA8keyRes.Data.(mm.GetA8KeyResp)
		if a8KeyRes.FullURL == nil && len(*(a8KeyRes.FullURL)) == 0 {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("异常：First扫码失败"),
				Data:    nil,
			}
		}
	} else {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：First扫码失败"),
			Data:    nil,
		}
	}
	firstCookie := base64.StdEncoding.EncodeToString(a8KeyRes.Cookie.GetBuffer())

	secondA8keyRes := Tools.GetA8Key(Tools.GetA8KeyParam{
		Wxid:         Data.Wxid,
		OpCode:       2,
		Scene:        4,
		CodeType:     0,
		CodeVersion:  0,
		ReqUrl:       Data.Url,
		CookieBase64: firstCookie,
		NetType:      "WIFI",
	})

	var protocolRes mm.GetA8KeyResp
	if secondA8keyRes.Success && secondA8keyRes.Data != nil {
		protocolRes = secondA8keyRes.Data.(mm.GetA8KeyResp)
		if protocolRes.FullURL == nil && len(*protocolRes.FullURL) == 0 {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("异常：Second扫码失败"),
				Data:    nil,
			}
		}
	} else {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：Second扫码失败"),
			Data:    nil,
		}
	}
	headers := make(map[string]string)
	for _, v := range protocolRes.HttpHeader {
		headers[*v.Key] = *v.Value
	}
	res, err := ScanIntoGrouppost(*protocolRes.FullURL, &headers, D)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("进群异常：%v", err.Error()),
			Data:    nil,
		}
	}
	if strings.Contains(res, "@chatroom") {
		return models.ResponseResult{
			Code:    0,
			Success: true,
			Message: "进群成功",
			Data:    res,
		}
	} else {
		if strings.Contains(res, "<h2 class=weui-msg__title>") {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
			if err != nil {
				return models.ResponseResult{
					Code:    -8,
					Success: false,
					Message: fmt.Sprintf("进群异常：%v", err.Error()),
					Data:    res,
				}
			}
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("进群异常：%v", doc.Find("h2.weui-msg__title").Text()),
				Data:    res,
			}
		} else if strings.Contains(res, "<div class=title>") {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
			if err != nil {
				return models.ResponseResult{
					Code:    -8,
					Success: false,
					Message: fmt.Sprintf("进群异常：%v", err.Error()),
					Data:    res,
				}
			}
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("进群异常：%v", doc.Find("div.title").Text()),
				Data:    res,
			}
		} else {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("进群异常"),
				Data:    res,
			}
		}
	}
}

func ScanIntoGroupEnterprise(Data ScanIntoGroupParam) models.ResponseResult {
	firstA8keyRes := Tools.GetA8Key3rd(Tools.GetA8KeyParam{
		Wxid:         Data.Wxid,
		OpCode:       2,
		Scene:        4,
		CodeType:     19,
		CodeVersion:  5,
		ReqUrl:       Data.Url,
		CookieBase64: "",
		NetType:      "",
	})
	var a8KeyRes mm.GetA8KeyResp
	if firstA8keyRes.Success && firstA8keyRes.Data != nil {
		a8KeyRes = firstA8keyRes.Data.(mm.GetA8KeyResp)
		if a8KeyRes.FullURL == nil && len(*(a8KeyRes.FullURL)) == 0 {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("异常：First扫码失败"),
				Data:    nil,
			}
		}
	} else {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：First扫码失败"),
			Data:    nil,
		}
	}

	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	headers := make(map[string]string)
	for _, v := range a8KeyRes.HttpHeader {
		headers[*v.Key] = *v.Value
	}
	res, err := ScanIntoGrouppost(*a8KeyRes.FullURL, &headers, D)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("进群异常：%v", err.Error()),
			Data:    nil,
		}
	}
	if strings.Contains(res, "@chatroom") {
		return models.ResponseResult{
			Code:    0,
			Success: true,
			Message: "进群成功",
			Data:    res,
		}
	} else {
		if strings.Contains(res, "<h2 class=weui-msg__title>") {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
			if err != nil {
				return models.ResponseResult{
					Code:    -8,
					Success: false,
					Message: fmt.Sprintf("进群异常：%v", err.Error()),
					Data:    res,
				}
			}
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("进群异常：%v", doc.Find("h2.weui-msg__title").Text()),
				Data:    res,
			}
		} else if strings.Contains(res, "<div class=title>") {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
			if err != nil {
				return models.ResponseResult{
					Code:    -8,
					Success: false,
					Message: fmt.Sprintf("进群异常：%v", err.Error()),
					Data:    res,
				}
			}
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("进群异常：%v", doc.Find("div.title").Text()),
				Data:    res,
			}
		} else {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("进群异常"),
				Data:    res,
			}
		}
	}
}

func ScanIntoGrouppost(URL string, headers *map[string]string, logindata *comm.LoginData) (string, error) {
	if strings.Index(URL, "weixin://jump/mainframe/") >= 0 {
		return strings.Replace(URL, "weixin://jump/mainframe/", "", -1), nil
	}

	userAgent := comm.GenUserAgent(logindata)
	body := comm.HttpPost1(URL, *new(url.Values), headers, userAgent, logindata.Proxy)

	if body != "" {
		return body, nil
	} else {
		return "", errors.New("POST请求失败")
	}
}
