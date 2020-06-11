package http

import (
	"context"
	"git/lambow-oj/lambow/constant"

	"github.com/gin-gonic/gin"
)

// 新版本的 gin 有个坑爹的问题，ginCtx的Done和Err实现采用其Request.Context()，但是ginCtx是用pool维护的，并发的时候会因为复用而导致context包panic
// 这里将ginCtx包装成一个安全的形式，其Value方法虽然是从ginCtx获取，但是在ginCtx放入pool后不会再被调用，因此还是安全的
type safeGinCtx struct {
	context.Context
	ginCtx *gin.Context
}

// http 头部信息
type Header struct {
}

// http 头部信息类型
type HeaderType struct {
}

// APIContext gin.Context 封装
type APIContext struct {
	//匿名包含Context接口，即Ctx本身也是一个合法的Context
	context.Context
	timeoutCtxCancel context.CancelFunc // 超时取消函数，因为Ctx默认继承自WithTimeout的Context
	ginCtx           *gin.Context       // 原始的ginCtx，有些地方需要用到
	uid              int64              // 当前请求的uid，从req中解析得来
}

// NewAPIContext 初始化 APIContext
func NewAPIContext(ginCtx *gin.Context) *APIContext {
	ctx, cancel := context.WithTimeout(newSafeGinCtx(ginCtx), constant.PROC_TIMEOUT)
	apiCtx := &APIContext{
		Context:          ctx,
		timeoutCtxCancel: cancel,
		ginCtx:           ginCtx,
	}

	return apiCtx
}

func newSafeGinCtx(ginCtx *gin.Context) *safeGinCtx {
	return &safeGinCtx{
		Context: ginCtx.Request.Context(),
		ginCtx:  ginCtx,
	}
}
