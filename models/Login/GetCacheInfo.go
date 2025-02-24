package Login

import (
	"fmt"
	"wechatdll/comm"
	"wechatdll/models"
)

func CacheInfo(Wxid string) models.ResponseResult {
	D, err := comm.GetLoginata(Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	return models.ResponseResult{
		Code:    1,
		Success: true,
		Message: "成功",
		Data:    D,
	}
}
