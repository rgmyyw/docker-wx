package bts

import (
	"encoding/json"
	"wechatdll/Cilent/mm"
)

func SearchContactResponse(Data interface{}) mm.SearchContactResponse {
	var Buff mm.SearchContactResponse
	result, err := json.Marshal(&Data)
	if err != nil {
		return mm.SearchContactResponse{}
	}
	_ = json.Unmarshal(result, &Buff)

	return Buff
}

func GetContactResponse(Data interface{}) mm.GetContactResponse {
	var Buff mm.GetContactResponse
	result, err := json.Marshal(&Data)
	if err != nil {
		return mm.GetContactResponse{}
	}
	_ = json.Unmarshal(result, &Buff)
	return Buff
}

func GetModContact(Data interface{}) mm.ModContact {
	var Buff mm.ModContact
	result, err := json.Marshal(&Data)
	if err != nil {
		return mm.ModContact{}
	}
	_ = json.Unmarshal(result, &Buff)
	return Buff
}

func GetA8KeyResponse(Data interface{}) mm.GetA8KeyResp {
	var Buff mm.GetA8KeyResp
	result, err := json.Marshal(&Data)
	if err != nil {
		return mm.GetA8KeyResp{}
	}
	_ = json.Unmarshal(result, &Buff)
	return Buff
}

func GetUserOpenIdResponse(Data interface{}) mm.GetUserOpenIdResp {
	var Buff mm.GetUserOpenIdResp
	result, err := json.Marshal(&Data)
	if err != nil {
		return mm.GetUserOpenIdResp{}
	}
	_ = json.Unmarshal(result, &Buff)
	return Buff
}



