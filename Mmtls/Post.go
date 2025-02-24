package Mmtls

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"wechatdll/models"

	"golang.org/x/net/proxy"
)

func (httpclient *HttpClientModel) POST(cgiurl string, data []byte, host string, P models.ProxyInfo) ([]byte, error) {
	var iphost string
	var err error
	iphost = "http://"
	iphost += host
	iphost += cgiurl
	body := bytes.NewReader(data)

	var Client *http.Client
	//设定代理
	if P.ProxyIp != "" && P.ProxyIp != "string" {
		var ProxyUser *proxy.Auth
		//设定账号和用户名
		if P.ProxyUser != "" && P.ProxyUser != "string" && P.ProxyPassword != "" && P.ProxyPassword != "string" {
			ProxyUser = &proxy.Auth{
				User:     P.ProxyUser,
				Password: P.ProxyPassword,
			}
		} else {
			ProxyUser = nil
		}
		Client, err = Socks5Client(P.ProxyIp, ProxyUser)
		if err != nil {
			return []byte{}, err
		}
	} else {
		Client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*15) //设置建立连接超时
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * 15)) //设置发送接受数据超时
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * 15,
				MaxIdleConnsPerHost:   -1,   //禁用连接池缓存
				DisableKeepAlives:     true, //禁用客户端连接缓存到连接池
			},
		}
	}

	request, err := http.NewRequest("POST", iphost, body)
	if err != nil {
		return []byte(""), err
	}
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Connection", "close")
	request.Header.Set("Content-type", "application/octet-stream")
	request.Header.Set("User-Agent", "MicroMessenger Client")
	request.Close = true
	var resp *http.Response
	resp, err = Client.Do(request)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return b, nil
}

func Socks5Client(addr string, auth *proxy.Auth) (client *http.Client, err error) {
	dialer, err := proxy.SOCKS5("tcp", addr,
		auth,
		&net.Dialer{
			Timeout:  15 * time.Second,
			Deadline: time.Now().Add(time.Second * 15),
		},
	)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy:               nil,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 15 * time.Second,
		MaxIdleConnsPerHost: -1,   //连接池禁用缓存
		DisableKeepAlives:   true, //禁用客户端连接缓存到连接池
	}

	client = &http.Client{Transport: transport}
	return
}
