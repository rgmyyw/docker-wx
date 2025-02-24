package Group

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"wechatdll/bts"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Tools"
)

func ConsentToJoin(Data ConsentToJoinParam) models.ResponseResult {
	_, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	A8key := Tools.GetA8Key(Tools.GetA8KeyParam{
		Wxid:        Data.Wxid,
		OpCode:      2,
		Scene:       4,
		CodeType:    19,
		CodeVersion: 5,
		ReqUrl:      Data.Url,
	})

	GetA8key := bts.GetA8KeyResponse(A8key.Data)

	if GetA8key.FullURL == nil {
		return models.ResponseResult{
			Code:    -9,
			Success: false,
			Message: "GetA8key读取失败",
			Data:    nil,
		}
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("POST", *GetA8key.FullURL, strings.NewReader("s=1"))
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

	if resp.StatusCode == 302 {
		QID := strings.Replace(resp.Header.Get("Location"), "weixin://jump/mainframe/", "", -1)
		return models.ResponseResult{
			Code:    0,
			Success: true,
			Message: "进群成功",
			Data:    QID,
		}
	}

	str,_ := ioutil.ReadAll(resp.Body)

	return models.ResponseResult{
		Code:    -9,
		Success: false,
		Message: "失败",
		Data:   string(str),
	}

}
