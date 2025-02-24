package Friend

type DefaultParam struct {
	Wxid   string
	ToWxid string
}

type BlacklistParam struct {
	Wxid   string
	ToWxid string
	Val    uint32
}

type LbsFindParam struct {
	Wxid   string
	Latitude float32
	Longitude float32
	OpCode uint32
}

