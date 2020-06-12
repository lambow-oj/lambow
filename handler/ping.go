package handler

import (
	"fmt"
	"git/lambow-oj/lambow/cinex"
	"git/lambow-oj/lambow/logic"
	"git/lambow-oj/lambow/proto"
)

func pingHandler(ctx *cinex.APIContext, req *proto.PingReq) (rsp *proto.PingRsp) {
	defer func() {

	}()
	rsp = new(proto.PingRsp)
	if err := pingParamsCheck(ctx, req); err != nil {
		rsp.SetErrRsp(ctx, proto.ERR_PARAM)
		return
	}
	ans, err := logic.Add(ctx, req.A, req.B)
	if err != nil {
		rsp.SetErrRsp(ctx, proto.ERR_PING)
	}
	rsp.Data.Pong = ans
	rsp.SetErrRsp(ctx, proto.ERR_SUCCESS)
	return
}

func pingParamsCheck(ctx *cinex.APIContext, req *proto.PingReq) error {
	if req.A < -10 || req.A > 10 {
		return fmt.Errorf("invalid a: %v", req.A)
	}
	if req.B < -10 || req.B > 10 {
		return fmt.Errorf("invalid b: %v", req.B)
	}

	return nil
}
