package OfficialAccounts

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"time"
	"wechatdll/bts"
	"wechatdll/comm"
	"wechatdll/models"
)

func GetAppMsgExt(Data ReadParam) models.ResponseResult {
	//获取mp-geta8key
	geta8key := MpGetA8Key(Data)
	if geta8key.Success == false {
		return geta8key
	}

	GetA8KEY := bts.GetA8KeyResponse(geta8key.Data)
	if GetA8KEY.FullURL == nil || *GetA8KEY.HttpHeaderCount < 2 {
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "异常,文章解析失败。",
			Data:    GetA8KEY,
		}
	}
	// 模拟打开页面
	D, err := comm.GetLoginata(Data.Wxid)
	ua := comm.GenUserAgent(D)
	header := make(map[string]string)
	for _, v := range GetA8KEY.HttpHeader {
		header[*v.Key] = *v.Value
	}
	getBody, cookie, err := comm.HttpGetBodyAndCookie1(*GetA8KEY.FullURL, &header, ua, D.Proxy)
	log.Infof("阅读: %s", getBody[:64])
	header["Cookie"] = cookie
	// 解析
	Url2, err := url.Parse(*GetA8KEY.FullURL)
	if err != nil {
		return models.ResponseResult {
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    GetA8KEY,
		}
	}
	JX2 := Url2.Query()
	var PostUrl string
	var PostData string
	PostUrl = "https://mp.weixin.qq.com/mp/getappmsgext?f=json&mock=&fasttmplajax=1&f=json&uin=&key=&pass_ticket=" + JX2["pass_ticket"][0] + "&wxtoken=&devicetype=iOS13.3.1&clientversion=17000c2d&__biz=" + url.QueryEscape(JX2["__biz"][0]) + "&appmsg_token=&x5=0&f=json&wx_header=1&pass_ticket=" + JX2["pass_ticket"][0]
	PostData = "r=&__biz=" + url.QueryEscape(JX2["__biz"][0]) + "&appmsg_type=9&mid=" + JX2["mid"][0] + "&sn=" + JX2["sn"][0] + "&idx=" + JX2["idx"][0] + "&scene=126&title=&ct=" + time.Now().String() + "&abtest_cookie=&devicetype=iOS13.3.1&version=17000c2d&is_need_ticket=0&is_need_ad=0&comment_id=&is_need_reward=0&both_ad=0&reward_uin_count=0&send_time=&msg_daily_idx=1&is_original=0&is_only_read=1&req_id=&pass_ticket=" + JX2["pass_ticket"][0] + "&is_temp_url=0&item_show_type=0&tmp_version=1&more_read_type=0&appmsg_like_type=2&related_video_sn=&vid=&is_pay_subscribe=0&pay_subscribe_uin_count=0&has_red_packet_cover=0&album_id="
	postBody, err := url.ParseQuery(PostData)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    GetA8KEY,
		}
	}
	res := comm.HttpPost1(PostUrl, postBody, &header, ua, D.Proxy)
	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    res,
	}
}

func GetAppMsgExtLike(Data ReadParam) models.ResponseResult {
	//获取mp-geta8key
	geta8key := MpGetA8Key(Data)
	if geta8key.Success == false {
		return geta8key
	}

	GetA8KEY := bts.GetA8KeyResponse(geta8key.Data)
	if GetA8KEY.FullURL == nil || *GetA8KEY.HttpHeaderCount < 2 {
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "异常,文章解析失败。",
			Data:    GetA8KEY,
		}
	}
	// 模拟打开页面
	D, err := comm.GetLoginata(Data.Wxid)
	ua := comm.GenUserAgent(D)
	header := make(map[string]string)
	for _, v := range GetA8KEY.HttpHeader {
		header[*v.Key] = *v.Value
	}
	getBody, cookie, err := comm.HttpGetBodyAndCookie1(*GetA8KEY.FullURL, &header, ua, D.Proxy)
	// log.Infof("阅读: %s", getBody[:64])
	header["Cookie"] = cookie
	// 解析
	Url2, err := url.Parse(*GetA8KEY.FullURL)
	if err != nil {
		return models.ResponseResult {
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    GetA8KEY,
		}
	}
	JX2 := Url2.Query()
	var PostUrl string
	var PostData string
	PostUrl = "https://mp.weixin.qq.com/mp/getappmsgext?f=json&mock=&fasttmplajax=1&f=json&uin=&key=&pass_ticket=" + JX2["pass_ticket"][0] + "&wxtoken=&devicetype=iOS13.3.1&clientversion=17000c2d&__biz=" + url.QueryEscape(JX2["__biz"][0]) + "&appmsg_token=&x5=0&f=json&wx_header=1&pass_ticket=" + JX2["pass_ticket"][0]
	PostData = "r=&__biz=" + url.QueryEscape(JX2["__biz"][0]) + "&appmsg_type=9&mid=" + JX2["mid"][0] + "&sn=" + JX2["sn"][0] + "&idx=" + JX2["idx"][0] + "&scene=126&title=&ct=" + time.Now().String() + "&abtest_cookie=&devicetype=iOS13.3.1&version=" + strconv.FormatInt(int64(D.ClientVersion), 16) + "&is_need_ticket=0&is_need_ad=0&comment_id=&is_need_reward=0&both_ad=0&reward_uin_count=0&send_time=&msg_daily_idx=1&is_original=0&is_only_read=1&req_id=&pass_ticket=" + JX2["pass_ticket"][0] + "&is_temp_url=0&item_show_type=0&tmp_version=1&more_read_type=0&appmsg_like_type=2&related_video_sn=&vid=&is_pay_subscribe=0&pay_subscribe_uin_count=0&has_red_packet_cover=0&album_id="
	postBody, err := url.ParseQuery(PostData)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    GetA8KEY,
		}
	}
	res := comm.HttpPost1(PostUrl, postBody, &header, ua, D.Proxy)

	//点赞
	appmsgid := comm.RegExpGet(&getBody, `var appmsgid = ['"+][0-9]*['"+] \|\| ['"+][0-9]*['"+] \|\| ['"+][0-9]*['"+]`, `\D`)
	likeUrl := "https://mp.weixin.qq.com/mp/appmsg_like?__biz=" + url.QueryEscape(JX2["__biz"][0]) + "&mid=" + JX2["mid"][0] + "&idx=" + JX2["idx"][0] + "&like=1&f=json&appmsgid=" + appmsgid + "&itemidx=&fasttmplajax=1&f=json&uin=&key=&pass_ticket=" + JX2["pass_ticket"][0] + "&wxtoken=777&devicetype=iOS13.3&clientversion=" + strconv.FormatInt(int64(D.ClientVersion), 16) + "&appmsg_token=&x5=0&f=json&wx_header=1&pass_ticket=" + JX2["pass_ticket"][0]
	likeData := "scene=90&appmsg_like_type=1&item_show_type=0&client_version=" + strconv.FormatInt(int64(D.ClientVersion), 16) + "&is_temp_url=0&style=0&exptype=&expsessionid="
	likeBody, err := url.ParseQuery(likeData)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    GetA8KEY,
		}
	}
	res2 := comm.HttpPost1(likeUrl, likeBody, &header, ua, D.Proxy)
	log.Infof("点赞结果: %s", res2)

	return models.ResponseResult{
		Code:    0,
		Success: true,
		Message: "成功",
		Data:    res,
	}
}
