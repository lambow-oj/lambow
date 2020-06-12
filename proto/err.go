package proto

import (
	"context"
	"time"
)

const (
	ERR_SUCCESS = "0"   // 成功
	ERR_COMMON  = "101" // 通用错误
	ERR_PARAM   = "102" // 参数校验失败
	ERR_PING    = "103" // ping 测试失败
)

type APIErrorInfo struct {
	ErrMsg string
}

func buildError(msg string) *APIErrorInfo {
	errInfo := new(APIErrorInfo)
	errInfo.ErrMsg = msg
	return errInfo
}

var errMsgMap = map[string]*APIErrorInfo{
	ERR_COMMON:  buildError("system busy"),
	ERR_SUCCESS: buildError("success"),
	ERR_PARAM:   buildError("system busy"),
	ERR_PING:    buildError("ping err"),
}

func (base *Base) SetErrRsp(ctx context.Context, errCode string) {
	base.Code = errCode
	base.Message = GetErrMsg(errCode)
	base.SvrTime = time.Now().Unix()
}

func GetErrMsg(errCode string) string {
	errMsg := "Unkown error"
	errInfo, ok := errMsgMap[errCode]
	if ok && errInfo != nil {
		errMsg = errInfo.ErrMsg
	}

	return errMsg
}

func SetErrRspBase(ctx context.Context, errCode string) *Base {
	rsp := new(Base)
	rsp.SetErrRsp(ctx, errCode)
	return rsp
}
