// +build !linux

package TcpPoll

import (
	"errors"
	"wechatdll/comm"
)

type TcpManager struct {
	running 		bool					// 控制是否消息loop
}

type TcpClient struct {
	running 		bool					// 控制是否消息loop
}

func GetTcpManager() (*TcpManager, error) {
	return nil, errors.New("windows不支持")
}

func (manager *TcpManager) Add(key string, client *TcpClient) error {
	return errors.New("windows不支持")
}

func (manager *TcpManager) Remove(client *TcpClient) {
	return
}

func (manager *TcpManager) GetClient(loginData *comm.LoginData) (*TcpClient, error) {
	return nil, errors.New("windows不支持")
}

func (manager *TcpManager) RunEventLoop() {
	return
}

func (client *TcpClient) MmtlsSend(data []byte, cmdId int, tag string) (*[]byte, error) {
	return nil, errors.New("windows不支持")
}