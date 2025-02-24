package OfficialAccounts

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"wechatdll/Algorithm"
	"wechatdll/Cilent/mm"
	"wechatdll/bts"
	"wechatdll/comm"
	"wechatdll/models"
	"wechatdll/models/Friend"
)

func Quit(Data DefaultParam) models.ResponseResult {
	D, err := comm.GetLoginata(Data.Wxid)
	if err != nil {
		return models.ResponseResult{
			Code:    -8,
			Success: false,
			Message: fmt.Sprintf("异常：%v", err.Error()),
			Data:    nil,
		}
	}

	//先获取用户的基本信息
	getContact := Friend.GetContractDetail(Friend.GetContractDetailparameter{
		Wxid:     Data.Wxid,
		Towxids:  Data.Appid,
		ChatRoom: "",
	})

	if getContact.Code != 0 {
		return getContact
	}

	Contact := bts.GetContactResponse(getContact.Data)

	if len(Contact.ContactList) > 0 {
		modContact := Contact.ContactList[0]
		ContactList := &mm.ModContact{
			UserName:        modContact.UserName,
			NickName:        modContact.NickName,
			PyInitial:       modContact.Pyinitial,
			QuanPin:         modContact.QuanPin,
			Sex:             modContact.Sex,
			ImgBuf:          modContact.ImgBuf,
			BitMask:         modContact.BitMask,
			BitVal:          proto.Uint32(2),
			ImgFlag:         modContact.ImgFlag,
			Remark:          modContact.Remark,
			RemarkPyinitial: modContact.RemarkPyinitial,
			RemarkQuanPin:   modContact.RemarkQuanPin,
			ContactType:     modContact.ContactType,
			ChatRoomNotify:  proto.Uint32(1),
			AddContactScene: modContact.AddContactScene,
			Extflag:         proto.Int32(int32(*modContact.ExtFlag)),
		}

		var cmdItems []*mm.CmdItem
		buffer, err := proto.Marshal(ContactList)
		if err != nil {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("系统异常：%v", err.Error()),
				Data:    nil,
			}
		}

		cmdItem := mm.CmdItem{
			CmdId: proto.Int32(2),
			CmdBuf: &mm.SKBuiltinBufferT{
				ILen:   proto.Uint32(uint32(len(buffer))),
				Buffer: buffer,
			},
		}
		cmdItems = append(cmdItems, &cmdItem)

		req := &mm.OpLogRequest{
			Cmd: &mm.CmdList{
				Count: proto.Uint32(uint32(len(cmdItems))),
				List:  cmdItems,
			},
		}

		reqdata, err := proto.Marshal(req)

		if err != nil {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("系统异常：%v", err.Error()),
				Data:    nil,
			}
		}

		//发包
		protobufdata, _, errtype, err := comm.SendRequest(comm.SendPostData{
			Ip:     D.Mmtlsip,
			Host:   D.MmtlsHost,
			Cgiurl: "/cgi-bin/micromsg-bin/oplog",
			Proxy:  D.Proxy,
			PackData: Algorithm.PackData{
				Reqdata:          reqdata,
				Cgi:              681,
				Uin:              D.Uin,
				Cookie:           D.Cooike,
				Sessionkey:       D.Sessionkey,
				EncryptType:      5,
				Loginecdhkey:     D.RsaPublicKey,
				Clientsessionkey: D.Clientsessionkey,
				UseCompress:      false,
			},
		}, D.MmtlsKey)

		if err != nil {
			return models.ResponseResult{
				Code:    errtype,
				Success: false,
				Message: err.Error(),
				Data:    nil,
			}
		}

		//解包
		Response := mm.OplogResponse{}
		err = proto.Unmarshal(protobufdata, &Response)

		if err != nil {
			return models.ResponseResult{
				Code:    -8,
				Success: false,
				Message: fmt.Sprintf("反序列化失败：%v", err.Error()),
				Data:    nil,
			}
		}

		return models.ResponseResult{
			Code:    0,
			Success: true,
			Message: "成功",
			Data:    Response,
		}

	}

	return models.ResponseResult{}
}
