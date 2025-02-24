package main

import (
	"fmt"
	"runtime"
	"time"
	"wechatdll/TcpPoll"
	"wechatdll/comm"
	_ "wechatdll/routers"

	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
)

func main() {
	longLinkEnabled, _ := beego.AppConfig.Bool("longlinkenabled")

	comm.RedisInitialize()
	_, err := comm.RedisClient.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("【Redis】连接失败，ERROR：%v", err.Error()))
	}

	sysType := runtime.GOOS

	if sysType == "linux" && longLinkEnabled {
		// LINUX系统
		tcpManager, err := TcpPoll.GetTcpManager()
		if err != nil {
			log.Errorf("TCP启动失败.")
		}
		go tcpManager.RunEventLoop()
	}

	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/"] = "swagger"

	beego.SetLogFuncCall(false)

	//// 添加日志拦截器
	//var FilterLog = func(ctx *context.Context) {
	//	url, _ := json.Marshal(ctx.Input.Data()["RouterPattern"])
	//	bodyParams := ctx.Input.RequestBody
	//	formParams, _ := json.Marshal(ctx.Request.Form)
	//	outputBytes, _ := json.Marshal(ctx.Input.Data()["json"])
	//	divider := " - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -"
	//	topDivider := "┌" + divider
	//	middleDivider := "├" + divider
	//	bottomDivider := "└" + divider
	//	outputStr := "\n" + topDivider + "\n│ 请求地址:" + string(url) + "\n" + middleDivider + "\n│ body参数: " + string(bodyParams) + "\n│ form参数: " + string(formParams) + "\n│ 返回数据:" + string(outputBytes[:64]) + "...\n" + bottomDivider
	//	log.Info(outputStr)
	//}
	//// 最后一个参数必须设置为false 不然无法打印数据
	//beego.InsertFilter("/*", beego.FinishRouter, FilterLog, false)

	// Zero: 测试BindQQ加密
	//qqCryptor := Algorithm.GetQQCryptor()
	//plainText, _ := hex.DecodeString("00090009001800160001000006001f1d5a7a280002373c5837dd000000000001001400014d6582203c5837dd6158ca5c000000000000010600681fe74febb4c4b17eb65e0dc30ce351b70822442483e0f47df41c32c55757df61945dc0e9a30ddb2e78a69e318b442431f935e8249f64bc51dc63d041f4fc82eb759ccb89b091ed0230660d1cc222b7ddb0a52436ea1c014d412c37b3aba128400d3b413ae740a3dc0116000a000000057c0001040000010000160001000000051f1d5a7a00000001280002370000204001070006000001900001014400989c085b69fec1e1af52f140be4badd6d27b43fcc17bf00af1b2cd166595ff63ac758975a90691c5646a7b19b255da8ab0cbf9091110f5173e5b24703d8dd1edc0e8c5e060c134b6e05cdb75b9784f6112fd192d7e50fb64f306f3d6808d87a39647c82ec0d21db899bf42be8ee6db98ed1366e5010b9ecf1cc8f73e5dc6338a882120024a9bac4f2007efcd346e2723a45bfe63c3e43e552a014200120000000e636f6d2e74656e63656e742e6d6d0145001041313535653936343730323564343130")
	//fmt.Printf("原值        : %x\n", plainText)
	//encryptKey, _ := hex.DecodeString("536A466263586F45466652735778504C")
	//encryptText := qqCryptor.Encrypt(plainText, encryptKey)
	//fmt.Printf("加密结果     : %x\n", encryptText)
	//goDecryptText := qqCryptor.Decrypt(encryptText, encryptKey)
	//fmt.Printf("go加密解密结果: %x\n", goDecryptText)

	beego.Run()
	return

	deadline := time.Date(2025, time.Month(1), 3, 0, 0, 0, 0, time.Local)
	if time.Now().Before(deadline) {
		//启动
		beego.Run()
	} else {
		fmt.Println("请联系相关客服")
	}

}
