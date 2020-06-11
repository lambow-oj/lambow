package http

import (
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
)

type MiddlewareHandler func(*APIContext) (interface{}, error)

type PostOpt struct {
	MiddlewareHandlers []MiddlewareHandler
}

//post handler
type postHandler struct {
	funcValue reflect.Value // handler函数
	reqType   reflect.Type  // handler请求参数的Elem的类型，为nil表示没有此参数
	rspType   reflect.Type  // handler回复参数的Elem的类型，为nil表示没有此参数
	opt       PostOpt
}

func (handler *postHandler) mkCinexHandler() func(*gin.Context) {
	f := func(ginCtx *gin.Context) {
		// 封装
		ctx := NewAPIContext(ginCtx)
		defer ctx.timeoutCtxCancel()
		defer func() {
			if err := recover(); err != nil {
				if gp, ok := err.(ginPanicType); ok {
					//是gin抛出的错误，扔给上层
					panic(gp.err)
				}
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
			}
		}()
	}

	return f
}
