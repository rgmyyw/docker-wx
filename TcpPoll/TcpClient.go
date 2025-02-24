// +build linux

package TcpPoll

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/url"
	"strings"
	"time"
	"wechatdll/Algorithm"
	"wechatdll/Mmtls"
	"wechatdll/comm"
	"wechatdll/lib"
)

type TcpClient struct {
	establishCallback *func()              // 收到数据的回调事件
	model             *comm.LoginData      // wx缓存
	receivedBytes     []byte               // 接收数据缓存区
	pskData           [][]byte             // 握手队列
	handshaking       bool                 // 是否握手中
	mmtlsPrivateKey1  []byte               // 1次握手密钥
	mmtlsPublicKey1   []byte               // 1次握手公钥
	mmtlsPrivateKey2  []byte               // 1次握手密钥
	mmtlsPublicKey2   []byte               // 1次握手公钥
	handshakeHash     []byte               // 握手hash值
	longLinkEncodeIv  []byte               // 握手加密iv
	longLinkEncodeKey []byte               // 握手加密key
	longLinkDecodeIv  []byte               // 握手解密iv
	longLinkDecodeKey []byte               // 握手解密key
	lastHeartbeatTime int                  // 最后一次心跳时间
	shareKey          []byte               // 协商密钥
	serverSequence    int                  // 服务端发包序列
	clientSequence    int                  // 客户端发包序列
	mmtlsSequence     int                  // mmtls组包序列
	pskKey            []byte               // psk密钥
	queue             map[int]func([]byte) // 回调序列
	messageCache      []byte               // 消息缓存
	conn              net.Conn             // 长连接
}

var CORRECT_HEADER = []byte{0x00, 0xf1, 0x03}

func NewTcpClient(model *comm.LoginData) *TcpClient {
	instance := TcpClient{
		model: model,
	}
	instance.handshaking = true // 新建的连接先需要握手
	instance.mmtlsPrivateKey1, instance.mmtlsPublicKey1 = Algorithm.GetECDH415Key()
	instance.mmtlsPrivateKey2, instance.mmtlsPublicKey2 = Algorithm.GetECDH415Key()
	instance.serverSequence = 1
	instance.clientSequence = 1
	instance.mmtlsSequence = 1
	instance.queue = make(map[int]func([]byte))

	return &instance
}

// 创建长连接. remoteAddr string: 目标地址, proxyAddr string: 代理地址, proxyUser string: 代理用户, proxyPassword string: 代理密码
func CreateConnection(remoteAddr string, proxyAddr string, proxyUser string, proxyPassword string) (net.Conn, error) {
	// 判断proxy是否需要验证密码, 如果不需要则proxyAuth为空指针, 否则指向初始化的proxy.Auth结构
	var proxyAuth *proxy.Auth = nil
	if proxyUser != "" {
		proxyAuth = &proxy.Auth{
			User:     proxyUser,
			Password: proxyPassword,
		}
	}
	var conn net.Conn
	var connErr error
	if proxyAddr != "" {
		// 创建proxy连接到远程服务器
		dialer, dialErr := proxy.SOCKS5("tcp", proxyAddr, proxyAuth, proxy.Direct)
		if dialErr != nil {
			return nil, dialErr
		}
		conn, connErr = dialer.Dial("tcp", remoteAddr+":80")
		if connErr != nil {
			return nil, connErr
		}
	} else {
		// 直接连接远程服务器
		conn, connErr = net.Dial("tcp", remoteAddr)
		if connErr != nil {
			return nil, connErr
		}
	}
	return conn, nil
}

// 根据wx缓存创建并连接到长连接
func (client *TcpClient) Connect() error {
	var connErr error
	client.conn, connErr = CreateConnection(
		client.model.MarsHost,
		client.model.Proxy.ProxyIp,
		client.model.Proxy.ProxyUser,
		client.model.Proxy.ProxyPassword)
	if connErr != nil {
		if strings.Contains(connErr.Error(), "missing port in address") {
			client.conn, connErr = CreateConnection(
				client.model.MarsHost + ":80",
				client.model.Proxy.ProxyIp,
				client.model.Proxy.ProxyUser,
				client.model.Proxy.ProxyPassword)
			if connErr != nil {
				return connErr
			}
		} else {
			return connErr
		}
	}

	helloWrapper, helloBuff := client.BuildClientHello()
	client.handshakeHash = append(client.handshakeHash, helloBuff...)
	client.Send(helloWrapper, "客户端握手开始")
	return nil
}

// 接收数据, 回调至ReceiveMessage处理
func (client *TcpClient) Once() {
	buf := make([]byte, 0x10)
	len, err := client.conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			fmt.Println("TcpClient: Once() read error:", err)
		}
		return
	}
	// log.Infof("TcpClient: Once执行, 收到消息[%v]", hex.EncodeToString(buf[:len]))
	client.ReceiveMessage(buf[:len])
}

// 识别并预处理消息字节
func (client *TcpClient) ReceiveMessage(buf []byte) {
	client.receivedBytes = append(client.receivedBytes, buf...)
	for len(client.receivedBytes) > 5 {
		// 正常的消息报文头为5字节, 小于5字节的报文说明没有接收完成, 不做预处理, 继续接收, 大于5则解析消息头, 逐条取出消息并回调
		// 取出包头, 5字节长度
		headBytes := client.receivedBytes[:5]
		// 校验报文头的合法性
		if headBytes[1] != CORRECT_HEADER[1] || headBytes[2] != headBytes[2] {
			log.Errorf("数据包头格式错误: [%x]", headBytes)
			// 清除缓冲区并退出循环
			client.receivedBytes = []byte{}
			break
		}
		// 消息头3-5为完整消息长度
		messageLen := binary.BigEndian.Uint16(headBytes[3:5])
		if len(client.receivedBytes) < int(messageLen+5) {
			// 数据包长度小于报文长度, 说明还没有接收完成, 退出循环继续接收消息
			//log.Infof("TcpClient: ReceiveMessage执行, 消息[%v]完成[%v]", messageLen + 5, len(client.receivedBytes))
			break
		}
		// 取出该条消息
		messageBytes := client.receivedBytes[:messageLen+5]
		// 从缓存中移除该条消息
		client.receivedBytes = client.receivedBytes[messageLen+5:]
		// 回调消息处理
		client.ProcessMessage(messageBytes)
	}
}

// 处理接收到的消息, 包括mmtls解包和ecdh密钥交换
func (client *TcpClient) ProcessMessage(messageBytes []byte) {
	//log.Infof("TcpClient: ProcessMessage执行, 处理消息{%v}[%x]", len(messageBytes), messageBytes)
	if client.handshaking {
		// 握手阶段的消息处理
		if len(client.pskData) < 5 {
			// 握手交互小于5次, 将消息放入握手队列
			client.pskData = append(client.pskData, messageBytes)
		}
		if len(client.pskData) == 4 {
			// 前四次握手交互
			if err := client.ProcessServerHello(client.pskData[0]); err != nil {
				log.Errorf("第一次握手密钥协商错误")
				return
			}
			if err := client.GetCertificateVerify(client.pskData[1]); err != nil {
				log.Errorf("第二次握手密钥协商错误")
				return
			}
			if err := client.BuildServerTicket(client.pskData[2]); err != nil {
				log.Errorf("第三次握手密钥协商错误")
				return
			}
			if err := client.ServerFinish(client.pskData[3]); err != nil {
				log.Errorf("第四次握手密钥协商错误")
				return
			}
			clientFinishPackage := client.ClientFinish(client.pskData[3])
			//log.Infof("ClientFinish: %x", clientFinishPackage)
			if clientFinishPackage == nil {
				log.Errorf("ClientFinish失败")
				return
			}
			client.Send(clientFinishPackage, "客户端握手结束")
			client.handshaking = false
			go client.SendTcpHeartBeat()
		}
	} else {
		if messageBytes[0] == 0x17 {
			// mmtls解包并回调
			message := client.UnpackMmtlsLong(messageBytes)
			//log.Infof("TcpClient: ProcessMessage执行, 消息[%x]mmtls解密为[%x]", messageBytes, message)
			client.HandleMessage(message)
		} else if messageBytes[0] == 0x15 {
			log.Infof("TcpClient: ProcessMessage执行, 终止长连接[%x]", messageBytes)
			client.Terminate()
		}
	}
}

// ----------------------------- mmtls握手内容 --------------------------------------
// 发起握手组包
func (client *TcpClient) BuildClientHello() ([]byte, []byte) {
	var helloBuff = new(bytes.Buffer)
	// part_1: 固定头, 10 - 0 = 10位
	helloBuff.Write([]byte{0x0, 0x0, 0x0, 0xd0, 0x1, 0x3, 0xf1, 0x1, 0xc0, 0x2b})
	// part_2: 随机32位bytes, 42 - 10 = 32位
	helloBuff.Write([]byte(lib.RandSeq(32)))
	// part_3: 时间戳, 46 - 42 = 4位
	time := time.Now().Unix()
	binary.Write(helloBuff, binary.BigEndian, (int32)(time))
	// part_4: 未知数据, 这里写死, 58 - 46 = 12位, TODO: 搞明白这里怎么模拟
	helloBuff.Write([]byte{0x00, 0x00, 0x00, 0xA2, 0x01, 0x00, 0x00, 0x00, 0x9D, 0x00, 0x10, 0x02})
	// part_5: 第一个公钥, 序列和长度一共6字节,所以这里长度是公钥长度加6, 62 - 58 = 4位
	binary.Write(helloBuff, binary.BigEndian, (int32)(len(client.mmtlsPublicKey1)+6))
	// part_6: 公钥序号, 66 - 62 = 4位
	binary.Write(helloBuff, binary.BigEndian, (int32)(1))
	// part_7: 公钥长度, 68 - 66 = 2位
	binary.Write(helloBuff, binary.BigEndian, (int16)(len(client.mmtlsPublicKey1)))
	// part_8: 第一个公钥值, 133 - 68 = 65位
	helloBuff.Write(client.mmtlsPublicKey1)
	// part_9: 第二个公钥, 137 - 133 = 4位
	binary.Write(helloBuff, binary.BigEndian, (int32)(len(client.mmtlsPublicKey2)+6))
	// part_10: 公钥序号, 141 - 137 = 4位
	binary.Write(helloBuff, binary.BigEndian, (int32)(2))
	// part_11: 公钥长度, 143 - 141 = 2位
	binary.Write(helloBuff, binary.BigEndian, (int16)(len(client.mmtlsPublicKey2)))
	// part_12: 第二个公钥值, 208 - 143 = 65位
	helloBuff.Write(client.mmtlsPublicKey2)
	// part_13: 尾部1, 212 - 208 = 4位
	binary.Write(helloBuff, binary.BigEndian, (int32)(1))
	// 外层包装, 215 - 212 = 3位
	var helloWrapper = new(bytes.Buffer)
	helloWrapper.Write([]byte{0x16, 0xf1, 0x03})
	binary.Write(helloWrapper, binary.BigEndian, (int16)(len(helloBuff.Bytes())))
	helloWrapper.Write(helloBuff.Bytes())
	return helloWrapper.Bytes(), helloBuff.Bytes()
}

// mmtls握手第一个收到的消息是ServerHello
func (client *TcpClient) ProcessServerHello(message []byte) error {
	keyPairSequence := binary.BigEndian.Uint32(message[57:61]) // 确定是使用密钥对1还是密钥对2
	serverPublicKey := message[63:128]                         // 从63位置取出65位服务端公钥
	privateKey := client.mmtlsPrivateKey1
	if keyPairSequence == 2 {
		privateKey = client.mmtlsPrivateKey2
	}
	shareKey := Algorithm.DoECDH415Key(privateKey, serverPublicKey)
	if shareKey == nil {
		return errors.New("Mmtls: 秘钥交互失败")
	} else if len(shareKey) != 32 {
		return errors.New("Mmtls: 交互秘钥长度存在异常")
	}
	client.shareKey = Mmtls.Getsha256(shareKey) // 协商再hash得到共享密钥
	client.handshakeHash = append(client.handshakeHash, message[5:]...)

	var infoBytes = new(bytes.Buffer)
	infoBytes.Write(Mmtls.Utf8ToBytes("handshake key expansion"))
	dataHash := Mmtls.Getsha256(client.handshakeHash)
	infoBytes.Write(dataHash)
	expandKey := Algorithm.Hkdf_Expand(sha256.New, client.shareKey, infoBytes.Bytes(), 56)
	//fmt.Println(len(hkdfexpandkey))
	client.longLinkEncodeKey = expandKey[:16]
	client.longLinkDecodeKey = expandKey[16:32]
	client.longLinkEncodeIv = expandKey[32:44]
	client.longLinkDecodeIv = expandKey[44:]
	return nil
}

// 使用IV和key解码第二次握手内容, 把结果放入hand_shake_hash
func (client *TcpClient) GetCertificateVerify(message []byte) error {
	pskLen := len(message)
	// 去除开头的5字节,包括后面16字节的tag
	dataSection := message[5:pskLen]
	var aad = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(aad, uint64(client.serverSequence))
	aad = append(aad, message[0:5]...)
	var iv []byte
	iv, client.serverSequence = Mmtls.GetDecryptIv(client.longLinkDecodeIv, client.serverSequence)
	serverSecret := Algorithm.NewAES_GCMDecrypter(client.longLinkDecodeKey, dataSection, iv, aad)
	if serverSecret == nil {
		err := errors.New("第二次握手密钥协商错误")
		return err
	}
	client.handshakeHash = append(client.handshakeHash, serverSecret...)
	return nil
}

// 解密第三次握手内容
func (client *TcpClient) BuildServerTicket(message []byte) error {
	pskLen := len(message)
	// 去除开头的5字节和后面16字节的tag
	dataSection := message[5:pskLen]
	var aad = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(aad, uint64(client.serverSequence))
	aad = append(aad, message[0:5]...)
	var iv []byte
	iv, client.serverSequence = Mmtls.GetDecryptIv(client.longLinkDecodeIv, client.serverSequence)
	client.pskKey = Algorithm.NewAES_GCMDecrypter(client.longLinkDecodeKey, dataSection, iv, aad)
	if client.pskKey == nil {
		err := errors.New("第三次握手密钥协商错误")
		return err
	}
	// client.build_short_link_ticket() 这里不需要短连接ticket了,短连接可自己握手
	client.handshakeHash = append(client.handshakeHash, client.pskKey...)
	return nil
}

// 解密第四次握手内容
func (client *TcpClient) ServerFinish(message []byte) error {
	pskLen := len(message)
	// 去除开头的5字节和后面16字节的tag
	dataSection := message[5:pskLen]
	var aad = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(aad, uint64(client.serverSequence))
	aad = append(aad, message[0:5]...)
	var iv []byte
	iv, client.serverSequence = Mmtls.GetDecryptIv(client.longLinkDecodeIv, client.serverSequence)
	serverFinishData := Algorithm.NewAES_GCMDecrypter(client.longLinkDecodeKey, dataSection, iv, aad)
	if serverFinishData == nil {
		err := errors.New("第四次握手密钥协商错误")
		return err
	}
	return nil
}

// mmtls握手完成组包
func (client *TcpClient) ClientFinish(message []byte) []byte {
	var infoBytes = new(bytes.Buffer)
	infoBytes.Write(Mmtls.Utf8ToBytes("client finished"))
	dataHash := Mmtls.Getsha256(client.handshakeHash)
	// infoBytes.Write(dataHash)
	clientFinishData := Algorithm.Hkdf_Expand(sha256.New, client.shareKey, infoBytes.Bytes(), 32)
	hmacHash := hmac.New(sha256.New, clientFinishData)
	hmacHash.Write(dataHash)
	hmacRet := hmacHash.Sum(nil)
	var aad = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(aad, uint64(client.clientSequence))
	aad = append(aad, message[0:5]...)
	sendData := []byte{0x00, 0x00, 0x00, 0x23, 0x14, 0x00, 0x20}
	sendData = append(sendData, hmacRet...)
	var iv []byte
	iv, client.clientSequence = Mmtls.GetEncryptIv(client.longLinkEncodeIv, client.clientSequence)
	clientFinishByte := Algorithm.NewAES_GCMEncrypter(client.longLinkEncodeKey, sendData, iv, aad)
	var clientFinishWrapperBuffer = new(bytes.Buffer)
	clientFinishWrapperBuffer.Write([]byte{0x16, 0xf1, 0x03})
	var finishLengthByte = []byte{0, 0}
	binary.BigEndian.PutUint16(finishLengthByte, uint16(len(clientFinishByte)))
	clientFinishWrapperBuffer.Write(finishLengthByte)
	clientFinishWrapperBuffer.Write(clientFinishByte)
	var secretBuffer = new(bytes.Buffer)
	secretBuffer.Write(Mmtls.Utf8ToBytes("expanded secret"))
	secretBuffer.Write(dataHash)
	secretData := Algorithm.Hkdf_Expand(sha256.New, client.shareKey, secretBuffer.Bytes(), 32)
	var expandBuffer = new(bytes.Buffer)
	expandBuffer.Write(Mmtls.Utf8ToBytes("application data key expansion"))
	expandBuffer.Write(dataHash)
	expandData := Algorithm.Hkdf_Expand(sha256.New, secretData, expandBuffer.Bytes(), 56)
	client.longLinkEncodeKey = expandData[0:16]
	client.longLinkDecodeKey = expandData[16:32]
	client.longLinkEncodeIv = expandData[32:44]
	client.longLinkDecodeIv = expandData[44:]
	return clientFinishWrapperBuffer.Bytes()
}

// ------------------------------ 组包解包处理 -------------------------------------------
// mmtls加密
func (client *TcpClient) PackMmtlsLong(plainMessage []byte) ([]byte, error) {
	//log.Infof("TcpClient: PackMmtlsLong: message(%v)[%x]", len(plainMessage), plainMessage)
	var nonce []byte
	var aadBuffer = new(bytes.Buffer)
	binary.Write(aadBuffer, binary.BigEndian, uint64(client.clientSequence))
	nonce, client.clientSequence = Mmtls.GetNonce(client.longLinkEncodeIv, client.clientSequence)
	aadBuffer.Write([]byte{0x17, 0xf1, 0x03})
	binary.Write(aadBuffer, binary.BigEndian, uint16(len(plainMessage)+16))
	cipherMessage := Algorithm.NewAES_GCMEncrypter(client.longLinkEncodeKey, plainMessage, nonce, aadBuffer.Bytes())
	if cipherMessage == nil {
		return nil, errors.New("AESGCM加密消息失败")
	}
	var wrapBuffer = new(bytes.Buffer)
	wrapBuffer.Write([]byte{0x17, 0xf1, 0x03})
	binary.Write(wrapBuffer, binary.BigEndian, uint16(len(cipherMessage)))
	wrapBuffer.Write(cipherMessage)
	return wrapBuffer.Bytes(), nil
}

// mmtls解密
func (client *TcpClient) UnpackMmtlsLong(messageBytes []byte) []byte {
	packData := messageBytes[5:len(messageBytes)]
	var aad = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(aad, uint64(client.serverSequence))
	aad = append(aad, messageBytes[0:5]...)
	var iv []byte
	iv, client.serverSequence = Mmtls.GetDecryptIv(client.longLinkDecodeIv, client.serverSequence)
	ret := Algorithm.NewAES_GCMDecrypter(client.longLinkDecodeKey, packData, iv, aad)
	return ret
}

// 组包头
func BuildWrapper(message []byte, cmdId int, mmtlsSeq int) []byte {
	dataWrapper := new(bytes.Buffer)
	binary.Write(dataWrapper, binary.BigEndian, int32(len(message)+16)) // 包头组成1: 总长度
	binary.Write(dataWrapper, binary.BigEndian, int16(16))              // 包头组成2: 头部长度
	binary.Write(dataWrapper, binary.BigEndian, int16(1))               // 包头组成3: 组包版本号
	binary.Write(dataWrapper, binary.BigEndian, int32(cmdId))           // 包头组成4: 操作命令ID
	binary.Write(dataWrapper, binary.BigEndian, int32(mmtlsSeq))        // 包头组成5: 组包序列
	dataWrapper.Write(message)
	return dataWrapper.Bytes()
}

// 解包头
func StripWrapper(message []byte, length int) ([]byte, int, int, int, int, int) {
	totalLength := binary.BigEndian.Uint32(message[:4])
	headLength := binary.BigEndian.Uint16(message[4:6])
	packVersion := binary.BigEndian.Uint16(message[6:8])
	cmdId := binary.BigEndian.Uint32(message[8:12])
	packSequence := binary.BigEndian.Uint32(message[12:16])
	packData := message[16:]
	return packData, int(cmdId), int(packSequence), int(packVersion), int(headLength), int(totalLength)
}

// ------------------------------ 消息执行 ----------------------------------------------
// 处理消息
func (client *TcpClient) HandleMessage(message []byte) {
	if len(message) > 0 {
		message = append(client.messageCache, message...)
	}
	messageBody, cmdId, packSequence, packVersion, headLength, totalLength := StripWrapper(message, len(message))
	// 超长的数据包会分包发送, 这里使用缓存将分包归总
	if totalLength > len(message) {
		// log.Infof("长包跟踪: 长度不对缓存不处理: 预计长度[%d], 实际长度[%d]", totalLength, len(message))
		client.messageCache = message
		return
	} else {
		client.messageCache = []byte{}
	}
	if cmdId != 1000000006 {
		log.Infof("TcpClient[%d]: handle_message收到消息{cmd_id: %v, pack_seq: %v, pack_version: %v, head_length: %v, total_length: %v}", socketFD(client.conn), cmdId, packSequence, packVersion, headLength, totalLength)
	}
	if cmdId == 24 {
		status := binary.BigEndian.Uint32(messageBody)
		// TODO: 回调SyncMessage
		log.Infof("收到24消息提醒, status[%d], 执行回调", status)
		syncUrl := strings.Replace(beego.AppConfig.String("syncmessagebusinessuri"), "{0}", client.model.Wxid, -1)
		go comm.HttpPost(syncUrl, *new(url.Values), nil, "", "", "", "")
		return
	}
	// 回调
	cb := client.queue[packSequence]
	if cb != nil {
		delete(client.queue, packSequence)
		cb(messageBody)
	}
}

// 发送数据
func (client *TcpClient) Send(data []byte, tag string) {
	if tag != "Tcp心跳" {
		log.Infof("TcpClient[%d]: Send执行[%s], 发送消息(%v)[%x...]", socketFD(client.conn), tag, len(data), data[:32])
	}
	client.conn.Write(data)
}

// mmtls发包
func (client *TcpClient) MmtlsSend(data []byte, cmdId int, tag string) (*[]byte, error) {
	// mmtls组包头
	dataWrapper := BuildWrapper(data, cmdId, client.mmtlsSequence) // mmtls加密
	sendData, err := client.PackMmtlsLong(dataWrapper)
	if err != nil {
		return nil, err
	}
	// 接收数据回调处理, 此时数据已mmtls解密
	var resp *[]byte
	client.queue[client.mmtlsSequence] = func(recv []byte) {
		resp = &recv // 闭包, 获取tcp传回的数据
	}
	// 组包序列自增
	client.mmtlsSequence++
	// 发包
	client.Send(sendData, tag)
	// 等待结果
	timeoutSpan, _ := time.ParseDuration(beego.AppConfig.String("longlinkconnecttimeout"))
	timeoutTime := time.Now().Add(timeoutSpan)
	// 进入循环等待, 完成握手或者超时都将退出循环
	for time.Now().Before(timeoutTime) {
		time.Sleep(100 * time.Millisecond)
		// 通过resp判断是否已经完成握手
		if resp != nil {
			break
		}
	}

	if resp == nil {
		// 超时没有完成, 报错
		return nil, errors.New("请求超时")
	} else if len(*resp) < 32 {
		// 长度小于32, 用户session已失效
		return nil, errors.New("用户可能退出")
	}
	unpackData := Algorithm.UnpackBusinessPacket(*resp, client.model.Sessionkey, client.model.Uin, &client.model.Cooike)
	// 将cookie更新保存到redis
	err = comm.CreateLoginData(*client.model, client.model.Wxid, 0)
	if err != nil {
		log.Errorf("TcpClient: MmtlsSend回调时更新redis失败[%v]", err.Error())
		return nil, err
	}
	return &unpackData, nil
}

func (client *TcpClient) Terminate() {
	client = nil
}

func (client *TcpClient) SendTcpHeartBeat() {
	for {
		sendData, _ := client.PackMmtlsLong(BuildWrapper([]byte{}, 6, -1))
		client.Send(sendData, "Tcp心跳")
		time.Sleep(20 * time.Second)
	}
}
