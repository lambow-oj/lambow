package cinex

import (
	"fmt"
	"time"
)

const (
	COMMON_ERROR      = -100
	PANIC_ERROR       = -101
	BAD_REQUEST_ERROR = -102
	ECODE_UNSET_ERROR = -103
)

var errMsgMap = map[int]string{
	COMMON_ERROR:      "Common error",
	PANIC_ERROR:       "Code panic error",
	BAD_REQUEST_ERROR: "Bad Request error",
	ECODE_UNSET_ERROR: "Error code unset",
}

type RspBase struct {
	Code    string `json:"ret"`
	Msg     string `json:"errmsg"`
	SrvTime string `json:"systime"`
}

type RspBaseBased interface {
	SetCode(code int)
	GetRspBase() *RspBase
}

func NewRspBase(code int) *RspBase {
	rsp := &RspBase{}
	rsp.SetCode(code)
	return rsp
}

func (rsp *RspBase) SetCode(code int) {
	msg := errMsgMap[code]
	if msg == "" {
		msg = "Unknown Error"
	}
	rsp.Code = fmt.Sprintf("%d", code)
	rsp.Msg = msg
	rsp.SrvTime = fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond/time.Nanosecond))
}
