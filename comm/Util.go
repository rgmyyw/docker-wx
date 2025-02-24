package comm

import (
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/models"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

// 设定代理中间件, 如proxyAddr为空则不适用代理
func Socks5Transport(proxyAddr string, proxyUser string, proxyPass string) (client *http.Transport, err error) {

	//设定代理
	var proxyAuth *proxy.Auth
	var transport *http.Transport
	if proxyAddr != "" && proxyAddr != "string" {
		//设定账号和用户名
		if proxyUser != "" && proxyUser != "string" && proxyPass != "" && proxyPass != "string" {
			proxyAuth = &proxy.Auth{
				User:     proxyUser,
				Password: proxyPass,
			}
		} else {
			proxyAuth = nil
		}
		dialer, err := proxy.SOCKS5("tcp", proxyAddr,
			proxyAuth,
			&net.Dialer{
				Timeout:  15 * time.Second,
				Deadline: time.Now().Add(time.Second * 15),
			},
		)
		if err != nil {
			return nil, err
		}
		transport = &http.Transport{
			Proxy:               nil,
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: 15 * time.Second,
			MaxIdleConnsPerHost: -1,   //连接池禁用缓存
			DisableKeepAlives:   true, //禁用客户端连接缓存到连接池
		}
	} else {
		transport = &http.Transport{
			Proxy:               nil,
			TLSHandshakeTimeout: 15 * time.Second,
			MaxIdleConnsPerHost: -1,   //连接池禁用缓存
			DisableKeepAlives:   true, //禁用客户端连接缓存到连接池
		}
	}

	return transport, nil
}

func GenDefaultIpadUA() string {
	code := Algorithm.IPadVersion
	major := 0x0f & (code >> 24)
	minor := 0xff & (code >> 16)
	patch := 0xff & (code >> 8)
	//build := 0xff & (code >> 0)
	wxVersion := strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
	iPadOsVersionS := strings.Replace(Algorithm.IPadOsVersion, ".", "_", -1)
	wechatUserAgent := fmt.Sprintf("Mozilla/5.0 (iPad; CPU iPad iPhone OS %s like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/%s(%s) NetType/WIFI Language/zh_CN", iPadOsVersionS, wxVersion, strconv.Itoa(code))
	return wechatUserAgent
}

func GenDefaultAndroidUA() string {
	code := Algorithm.AndroidVersion
	major := 0x0f & (code >> 24)
	minor := 0xff & (code >> 16)
	patch := 0xff & (code >> 8)
	build := 0xff & (code >> 0)
	wxVersion := strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
	wechatUserAgent := fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; %s %s Build/%s AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/78.0.3904.62 XWEB/2692 MMWEBSDK/201001 Mobile Safari/537.36 MMWEBID/1135 MicroMessenger/%s.%s(0x%s) Process/toolsmp WeChat/arm32 Weixin NetType/WIFI Language/zh_CN ABI/arm64", GetAndroidOsVersion(Algorithm.AndroidDeviceType), Algorithm.AndroidManufacture, Algorithm.AndroidModel, Algorithm.AndroidDeviceType, wxVersion, strconv.Itoa(build), strconv.FormatInt(int64(code), 16))
	return wechatUserAgent
}

func GenUserAgent(loginData *LoginData) string {
	var wechatUserAgent string
	code := loginData.ClientVersion
	major := 0x0f & (code >> 24)
	minor := 0xff & (code >> 16)
	patch := 0xff & (code >> 8)
	build := 0xff & (code >> 0)
	wxVersion := strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
	iPadOsVersionS := strings.Replace(Algorithm.IPadOsVersion, ".", "_", -1)
	if strings.Index(loginData.DeviceType, "android") >= 0 {
		wechatUserAgent = fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; %s %s Build/%s AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/78.0.3904.62 XWEB/2692 MMWEBSDK/201001 Mobile Safari/537.36 MMWEBID/1135 MicroMessenger/%s.%s(0x%s) Process/toolsmp WeChat/arm32 Weixin NetType/WIFI Language/zh_CN ABI/arm64", GetAndroidOsVersion(loginData.DeviceType), Algorithm.AndroidManufacture, Algorithm.AndroidModel, loginData.DeviceType, wxVersion, strconv.Itoa(build), strconv.FormatInt(int64(code), 16))
	} else {
		wechatUserAgent = fmt.Sprintf("Mozilla/5.0 (iPad; CPU iPad iPhone OS %s like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/%s(0x%s) NetType/WIFI Language/zh_CN", iPadOsVersionS, wxVersion, strconv.FormatInt(int64(code), 16))
	}
	return wechatUserAgent
}

func GetAndroidOsVersion(deviceType string) string {
	switch deviceType {
	case "android-18":
		return "4.3.1"
	case "android-19":
		return "4.4"
	case "android-20":
		return "4.4.4"
	case "android-21":
		return "5.0"
	case "android-22":
		return "5.1.1"
	case "android-23":
		return "6.0.1"
	case "android-24":
		return "7.0"
	case "android-25":
		return "7.1.2"
	case "android-26":
		return "8.0"
	case "android-27":
		return "8.1"
	case "android-28":
		return "9"
	case "android-29":
		return "10"
	case "android-30":
		return "11"
	case "android-31":
		return "12"
	default:
		return "8.0"
	}
}

// 模拟微信提交申请, 如果需要设置Cookie, 将其放入headers即可
func WxHttpRequest(Url string, action string, headers *map[string]string, body url.Values, ua string, proxyAddr string, proxyUser string, proxyPass string) (*http.Response, error) {
	var req *http.Request
	var err error
	if body != nil {
		reqBody := strings.NewReader(body.Encode())
		req, err = http.NewRequest(action, Url, reqBody)
	} else {
		req, err = http.NewRequest(action, Url, nil)
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if ua == "" {
		ua = GenDefaultIpadUA()
	}
	postUri, err := url.Parse(Url)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
	if action == "POST" {
		req.Header.Set("Host", postUri.Host)
		req.Header.Set("Origin", postUri.Scheme+"://"+postUri.Host)
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Requested-With", "com.tencent.mm")
	}
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}
	transport, err := Socks5Transport(proxyAddr, proxyUser, proxyPass)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	return resp, err
}

func WxHttpGetBody(resp *http.Response) (string, error) {
	// 解析body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		bodyReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
		defer bodyReader.Close()
		data, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			return "", err
		} else {
			res := string(data)
			return res, nil
		}
	} else {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		} else {
			res := string(data)
			return res, nil
		}
	}
}

func WxHttpGetCookie(resp *http.Response) string {
	setCookie := ""
	for _, value := range resp.Cookies() {
		setCookie += strings.Split(value.String(), ";")[0] + ";"
	}
	return setCookie
}

func HttpGet(Url string, headers *map[string]string, ua string, proxyAddr string, proxyUser string, proxyPass string) string {
	resp, err := WxHttpRequest(Url, "GET", headers, nil, ua, proxyAddr, proxyUser, proxyPass)
	if err != nil {
		log.Errorf("HttpGet异常: %s", err.Error())
		return ""
	}
	body, err := WxHttpGetBody(resp)
	if err != nil {
		log.Errorf("Http异常: %s", err.Error())
		return ""
	}
	return body
}

func HttpGet1(Url string, headers *map[string]string, ua string, proxy models.ProxyInfo) string {
	return HttpGet(Url, headers, ua, proxy.ProxyIp, proxy.ProxyUser, proxy.ProxyPassword)
}

func HttpPost(Url string, data url.Values, headers *map[string]string, ua string, proxyAddr string, proxyUser string, proxyPass string) string {
	resp, err := WxHttpRequest(Url, "POST", headers, data, ua, proxyAddr, proxyUser, proxyPass)
	if err != nil {
		return err.Error()
	}
	if resp == nil {
		return "群页面无法打开"
	}
	body, err := WxHttpGetBody(resp)
	if err != nil {
		log.Errorf("Http异常: %s", err.Error())
		return ""
	}
	return body
}

func HttpPost1(Url string, data url.Values, headers *map[string]string, ua string, proxy models.ProxyInfo) string {
	return HttpPost(Url, data, headers, ua, proxy.ProxyIp, proxy.ProxyUser, proxy.ProxyPassword)
}

func HttpGetAndSetCookie(Url string, headers *map[string]string, ua string, proxyAddr string, proxyUser string, proxyPass string) string {
	resp, _ := WxHttpRequest(Url, "GET", headers, nil, ua, proxyAddr, proxyUser, proxyPass)
	setCookie := WxHttpGetCookie(resp)
	return setCookie
}

func HttpGetBodyAndCookie(url string, headers *map[string]string, ua string, proxyAddr string, proxyUser string, proxyPass string) (string, string, error) {
	resp, _ := WxHttpRequest(url, "GET", headers, nil, ua, proxyAddr, proxyUser, proxyPass)
	body, err := WxHttpGetBody(resp)
	if err != nil {
		log.Errorf("Http异常: %s", err.Error())
		return "", "", err
	}
	setCookie := WxHttpGetCookie(resp)
	return body, setCookie, nil
}

func HttpGetBodyAndCookie1(url string, headers *map[string]string, ua string, proxy models.ProxyInfo) (string, string, error) {
	return HttpGetBodyAndCookie(url, headers, ua, proxy.ProxyIp, proxy.ProxyUser, proxy.ProxyPassword)
}

// 使用chromedp获取网站上爬取的数据
func ChromeDpGetContent(url string, ua string, selector string, sel interface{}) (string, error) {
	if ua == "" {
		ua = GenDefaultIpadUA()
	}

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(ua),
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//创建一个上下文，超时时间为40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		log.Infof("Run err : %v\n", err)
		return "", err
	}
	return htmlContent, nil
}

func Rmu0000(s string) string {
	str := make([]rune, 0, len(s))
	for _, v := range []rune(s) {
		if v == 0 {
			continue
		}
		str = append(str, v)
	}
	return string(str)
}

// 从s中找到findReg表达式的字符串, 再从中替换掉replaceReg匹配的字符, 返回剩余的字符串
func RegExpGet(s *string, findReg string, replaceReg string) string {
	res := ""
	findRegExp := regexp.MustCompile(findReg)
	regResult := findRegExp.FindAllStringSubmatch(*s, -1)
	if regResult != nil {
		replaceRegExp := regexp.MustCompile(replaceReg)
		res = replaceRegExp.ReplaceAllString(regResult[0][0], "")
	}
	return res
}
