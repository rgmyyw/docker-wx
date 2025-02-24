package models



type ResponseResult struct {
	Code    int64
	Success bool
	Message string
	Data    interface{}
	Data62  string
	Debug   string
}

type ResponseResult2 struct {
	Code     int64
	Success  bool
	Message  string
	Data     interface{}
	Data62   string
	DeviceId string
}

type SessionidQRParam struct {
	Code     int64
	Success  bool
	Message  string
	Data     interface{}
	Data62   string
	DeviceId string
}


