package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["wechatdll/controllers:FavorController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FavorController"],
		beego.ControllerComments{
			Method:           "Del",
			Router:           `/Del`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FavorController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FavorController"],
		beego.ControllerComments{
			Method:           "GetFavInfo",
			Router:           `/GetFavInfo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FavorController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FavorController"],
		beego.ControllerComments{
			Method:           "GetFavItem",
			Router:           `/GetFavItem`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FavorController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FavorController"],
		beego.ControllerComments{
			Method:           "Sync",
			Router:           `/Sync`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FinderController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FinderController"],
		beego.ControllerComments{
			Method:           "UserPrepare",
			Router:           `/UserPrepare`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"],
		beego.ControllerComments{
			Method:           "Comment",
			Router:           `/Comment`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"],
		beego.ControllerComments{
			Method:           "GetDetail",
			Router:           `/GetDetail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"],
		beego.ControllerComments{
			Method:           "GetIdDetail",
			Router:           `/GetIdDetail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"],
		beego.ControllerComments{
			Method:           "GetList",
			Router:           `/GetList`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"],
		beego.ControllerComments{
			Method:           "Messages",
			Router:           `/Messages`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"],
		beego.ControllerComments{
			Method:           "MmSnsSync",
			Router:           `/MmSnsSync`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"],
		beego.ControllerComments{
			Method:           "Operation",
			Router:           `/Operation`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"],
		beego.ControllerComments{
			Method:           "PrivacySettings",
			Router:           `/PrivacySettings`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendCircleController"],
		beego.ControllerComments{
			Method:           "Upload",
			Router:           `/Upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "Blacklist",
			Router:           `/Blacklist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/Delete`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "GetContractDetail",
			Router:           `/GetContractDetail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "GetContractList",
			Router:           `/GetContractList`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "GetMFriend",
			Router:           `/GetMFriend`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "PassVerify",
			Router:           `/PassVerify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "Search",
			Router:           `/Search`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "SendRequest",
			Router:           `/SendRequest`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "SetRemarks",
			Router:           `/SetRemarks`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "Upload",
			Router:           `/Upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
			

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "LbsFind",
			Router:           `/LbsFind`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
			

	beego.GlobalControllerRouter["wechatdll/controllers:FriendController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:FriendController"],
		beego.ControllerComments{
			Method:           "GetFriendstate",
			Router:           `/GetFriendRelation`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
			

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "AddChatRoomMember",
			Router:           `/AddChatRoomMember`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "ConsentToJoin",
			Router:           `/ConsentToJoin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "CreateChatRoom",
			Router:           `/CreateChatRoom`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "FacingCreateChatRoom",
			Router:           `/FacingCreateChatRoom`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "DelChatRoomMember",
			Router:           `/DelChatRoomMember`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "GetChatRoomInfo",
			Router:           `/GetChatRoomInfo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "GetChatRoomInfoDetail",
			Router:           `/GetChatRoomInfoDetail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "GetChatRoomMemberDetail",
			Router:           `/GetChatRoomMemberDetail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "GetQRCode",
			Router:           `/GetQRCode`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "InviteChatRoomMember",
			Router:           `/InviteChatRoomMember`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "MoveContractList",
			Router:           `/MoveContractList`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "OperateChatRoomAdmin",
			Router:           `/OperateChatRoomAdmin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "Quit",
			Router:           `/Quit`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "ScanIntoGroup",
			Router:           `/ScanIntoGroup`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "ScanIntoGroupEnterprise",
			Router:           `/ScanIntoGroupEnterprise`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "SetChatRoomAnnouncement",
			Router:           `/SetChatRoomAnnouncement`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "SetChatRoomName",
			Router:           `/SetChatRoomName`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:GroupController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:GroupController"],
		beego.ControllerComments{
			Method:           "SetChatRoomRemarks",
			Router:           `/SetChatRoomRemarks`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LabelController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LabelController"],
		beego.ControllerComments{
			Method:           "Add",
			Router:           `/Add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LabelController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LabelController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/Delete`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LabelController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LabelController"],
		beego.ControllerComments{
			Method:           "GetList",
			Router:           `/GetList`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LabelController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LabelController"],
		beego.ControllerComments{
			Method:           "UpdateList",
			Router:           `/UpdateList`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LabelController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LabelController"],
		beego.ControllerComments{
			Method:           "UpdateName",
			Router:           `/UpdateName`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "Data62Login",
			Router:           `/62data`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "Data62SMSApply",
			Router:           `/62dataSMSApply`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "Data62SMSAgain",
			Router:           `/62dataSMSAgain`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "Data62SMSVerify",
			Router:           `/62dataSMSVerify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "Data62QRCodeApply",
			Router:           `/62dataQrApply`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "Data62QRCodeVerify",
			Router:           `/62dataQrVerify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "A16Data",
			Router:           `/A16Data`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "A16Data1",
			Router:           `/A16Data1`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginAwaken",
			Router:           `/Awaken`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginCheckQR",
			Router:           `/CheckQR`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})


	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "ExtDeviceLoginConfirmGet",
			Router:           `/ExtDeviceLoginConfirmGet`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "ExtDeviceLoginConfirmOk",
			Router:           `/ExtDeviceLoginConfirmOk`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "Get62Data",
			Router:           `/Get62Data`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "GetA16Data",
			Router:           `/GetA16Data`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "GetCacheInfo",
			Router:           `/GetCacheInfo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQR",
			Router:           `/GetQR`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRx",
			Router:           `/GetQRx`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRPad",
			Router:           `/GetQRPad`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRPad1",
			Router:           `/GetQRPad1`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRPadx",
			Router:           `/GetQRPadx`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRPadppmt",
			Router:           `/GetQRPadppmt`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRWin",
			Router:           `/GetQRWin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRWinUwp",
			Router:           `/GetQRWinUwp`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
			

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRWinUnified",
			Router:           `/GetQRWinUnified`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
			

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRCar",
			Router:           `/GetQRCar`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginGetQRMac",
			Router:           `/GetQRMac`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "HeartBeat",
			Router:           `/HeartBeat`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "HeartBeatLong",
			Router:           `/HeartBeatLong`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LogOut",
			Router:           `/LogOut`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "Newinit",
			Router:           `/Newinit`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:LoginController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:LoginController"],
		beego.ControllerComments{
			Method:           "LoginTwiceAutoAuth",
			Router:           `/TwiceAutoAuth`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "Revoke",
			Router:           `/Revoke`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "SendApp",
			Router:           `/SendApp`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "SendCDNFile",
			Router:           `/SendCDNFile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "SendCDNImg",
			Router:           `/SendCDNImg`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "SendCDNVideo",
			Router:           `/SendCDNVideo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "SendEmoji",
			Router:           `/SendEmoji`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "SendTxt",
			Router:           `/SendTxt`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "SendVideo",
			Router:           `/SendVideo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "SendVoice",
			Router:           `/SendVoice`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "ShareCard",
			Router:           `/ShareCard`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "ShareLink",
			Router:           `/ShareLink`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "ShareLocation",
			Router:           `/ShareLocation`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "ShareVideo",
			Router:           `/ShareVideo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "Sync",
			Router:           `/Sync`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:MsgController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:MsgController"],
		beego.ControllerComments{
			Method:           "UploadImg",
			Router:           `/UploadImg`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"],
		beego.ControllerComments{
			Method:           "Follow",
			Router:           `/Follow`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"],
		beego.ControllerComments{
			Method:           "GetAppMsgExt",
			Router:           `/GetAppMsgExt`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"],
		beego.ControllerComments{
			Method:           "GetAppMsgExtLike",
			Router:           `/GetAppMsgExtLike`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"],
		beego.ControllerComments{
			Method:           "JSAPIPreVerify",
			Router:           `/JSAPIPreVerify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"],
		beego.ControllerComments{
			Method:           "MpGetA8Key",
			Router:           `/MpGetA8Key`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"],
		beego.ControllerComments{
			Method:           "OauthAuthorize",
			Router:           `/OauthAuthorize`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:OfficialAccountsController"],
		beego.ControllerComments{
			Method:           "Quit",
			Router:           `/Quit`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:SayHelloController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:SayHelloController"],
		beego.ControllerComments{
			Method:           "ModelV1",
			Router:           `/Modelv1`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:SayHelloController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:SayHelloController"],
		beego.ControllerComments{
			Method:           "Modelv2",
			Router:           `/Modelv2`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:SayHelloController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:SayHelloController"],
		beego.ControllerComments{
			Method:           "Modelv3",
			Router:           `/Modelv3`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "DownloadFile",
			Router:           `/DownloadFile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "DownloadImg",
			Router:           `/DownloadImg`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "CdnDownloadImage",
			Router:           `/CdnDownloadImage`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "DownloadVideo",
			Router:           `/DownloadVideo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "DownloadVoice",
			Router:           `/DownloadVoice`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "GeneratePayQCode",
			Router:           `/GeneratePayQCode`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "GetA8Key",
			Router:           `/GetA8Key`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "GetCdnDns",
			Router:           `/GetCdnDns`,
			AllowHTTPMethods: []string{"POST"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "GetBandCardList",
			Router:           `/GetBandCardList`,
			AllowHTTPMethods: []string{"POST"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "GetBoundHardDevices",
			Router:           `/GetBoundHardDevices`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "ThirdAppGrant",
			Router:           `/ThirdAppGrant`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "SetProxy",
			Router:           `/setproxy`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

			
	// 修改微信步数
	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "UpdateStepNumberApi",
			Router:           `/UpdateStep`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
			


	beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:ToolsController"],
		beego.ControllerComments{
			Method:           "OauthSdkApp",
			Router:           `/OauthSdkApp`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "BindingEmail",
			Router:           `/BindingEmail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "BindingMobile",
			Router:           `/BindingMobile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "BindQQ",
			Router:           `/BindQQ`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "DelSafetyInfo",
			Router:           `/DelSafetyInfo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetContractProfile",
			Router:           `/GetContractProfile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetQRCode",
			Router:           `/GetQRCode`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetSafetyInfo",
			Router:           `/GetSafetyInfo`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "PrivacySettings",
			Router:           `/PrivacySettings`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "ReviseMotion",
			Router:           `/ReportMotion`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "SendVerifyMobile",
			Router:           `/SendVerifyMobile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "SetAlisa",
			Router:           `/SetAlisa`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "SetPasswd",
			Router:           `/SetPasswd`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "UpdateProfile",
			Router:           `/UpdateProfile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "UploadHeadImage",
			Router:           `/UploadHeadImage`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:UserController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:UserController"],
		beego.ControllerComments{
			Method:           "VerifyPasswd",
			Router:           `/VerifyPasswd`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "JSLogin",
			Router:           `/JSLogin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "JSOperateWxData",
			Router:           `/JSOperateWxData`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "JSGetSessionid",
			Router:           `/JSGetSessionid`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "AddWxAppRecord",
			Router:           `/AddWxAppRecord`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "JSGetSessionidQRcode",
			Router:           `/JSGetSessionidQRcode`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "CloudCallFunction",
			Router:           `/CloudCallFunction`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "AddMobile",
			Router:           `/AddMobile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "DelMobile",
			Router:           `/DelMobile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "GetRandomAvatar",
			Router:           `/GetRandomAvatar`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "UploadAvatarImg",
			Router:           `/UploadAvatarImg`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "AddAvatar",
			Router:           `/AddAvatar`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "QrcodeAuthLogin",
			Router:           `/QrcodeAuthLogin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
			
	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "GetUserOpenId",
			Router:           `/GetUserOpenId`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})



	beego.GlobalControllerRouter["wechatdll/controllers:WxappController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:WxappController"],
		beego.ControllerComments{
			Method:           "GetAllMobile",
			Router:           `/GetAllMobile`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
			

	beego.GlobalControllerRouter["wechatdll/controllers:QWContactController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:QWContactController"],
		beego.ControllerComments{
			Method:           "QWAddContact",
			Router:           `/QWAddContact`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:QWContactController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:QWContactController"],
		beego.ControllerComments{
			Method:           "QWApplyAddContact",
			Router:           `/QWApplyAddContact`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["wechatdll/controllers:QWContactController"] = append(beego.GlobalControllerRouter["wechatdll/controllers:QWContactController"],
		beego.ControllerComments{
			Method:           "SearchQWContact",
			Router:           `/SearchQWContact`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
