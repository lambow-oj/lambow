package proto

type (
	Base struct {
		Code    string `json:"ret"` // 返回码，0成功，其他失败
		Message string `json:"errmsg"`
		SvrTime int64  `json:"svr_time"` //服务器返回请求时间戳
	}

	PingRsp struct {
		Base
		Data struct {
			Pong string `json:"pong"`
		} `json:"data"`
	}
)
