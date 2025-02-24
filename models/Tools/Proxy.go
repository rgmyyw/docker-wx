package Tools

import (
	"fmt"
	"wechatdll/comm"
	"wechatdll/models"
)

type SetProxyParam struct {
	Wxid  string
	Proxy models.ProxyInfo
}

func SetProxy(Data SetProxyParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//初始化Mmtls
	_, MmtlsClient, err := comm.MmtlsInitialize(Data.Proxy, D.MmtlsHost)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("MMTLS初始化失败：%v", err.Error()),
			Data:    nil,
		}
	}

	D.Proxy = Data.Proxy
	D.MmtlsKey = MmtlsClient

	if Data.Proxy.ProxyIp == "" {
		D.Proxy = models.ProxyInfo{}
	}

	err = comm.CreateLoginData(*D, Data.Wxid, 0)

	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: "失败",
			Data:    err.Error(),
		}
	}

	return models.ResponseResult{
		Code:    1,
		Success: true,
		Message: "成功",
		Data:    nil,
	}

}
