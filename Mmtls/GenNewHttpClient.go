package Mmtls

import (
	"net/http"
)

func GenNewHttpClient(Data *MmtlsClient, domain string) (httpclient *HttpClientModel) {

	mmtlsClient := &MmtlsClient{
		//不需要发送队列。
		ServerSeq: 1,
		ClientSeq: 1,
	}

	if Data != nil {
		mmtlsClient = Data
	}

	httpclientModel := &HttpClientModel{
		mmtlsClient: mmtlsClient,
		httpClient:  &http.Client{},
		curShortip:  domain,
	}

	return httpclientModel
}

/*func GenNewTcpClient(Data *MmtlsClient) (tcpClient *TcpClientModel) {

}*/
