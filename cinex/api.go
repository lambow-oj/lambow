package cinex

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

// 存放post handler的表，key是api path，value是handler信息
var postHandlerMap = map[string]*postHandler{}

func init() {
	gin.EnableJsonDecoderUseNumber() // 启用 json 解析的 use-number 选项
}

func (handler *postHandler) mkCinexHandler() func(*gin.Context) {
	f := func(ginCtx *gin.Context) {
		var errCode int
		var rspValue reflect.Value

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
				errCode = PANIC_ERROR
				rsp := NewRspBase(errCode)
				ginCtx.JSON(http.StatusInternalServerError, rsp)
			}
		}()
		// 参数列表
		args := make([]reflect.Value, 0, 2)
		args = append(args, reflect.ValueOf(ctx))

		// 执行通用 handler
		for _, middleware := range handler.opt.MiddlewareHandlers {
			rsp, err := middleware(ctx)
			if err != nil {
				ginPanicCall(func() {
					ginCtx.JSON(http.StatusOK, rsp)
				})
				errCode = COMMON_ERROR
				return
			}
		}

		// 解析请求参数
		if handler.reqType != nil {
			//生成请求
			reqValue := reflect.New(handler.reqType)
			req := reqValue.Interface()
			err := ctx.ginCtx.ShouldBindBodyWith(req, binding.JSON)
			if err != nil {
				ginPanicCall(func() {
					ginCtx.String(http.StatusBadRequest, fmt.Sprintf("Bad request [%s]", err.Error()))
				})
				errCode = BAD_REQUEST_ERROR
				return
			}
			args = append(args, reqValue)
		}

		// 解析响应参数
		if handler.rspType != nil {
			// 生成回复并将其转为 RspBaseBased，设置初始错误码
			rspValue = reflect.New(handler.rspType)
			rspValue.Interface().(RspBaseBased).SetCode(ECODE_UNSET_ERROR)
			args = append(args, rspValue)
		}

		// 调用handler
		if handler.rspType == nil {
			rspValue = handler.funcValue.Call(args)[0]
		} else {
			handler.funcValue.Call(args)
		}
		rsp := rspValue.Interface()

		rspBaseBased, ok := rsp.(RspBaseBased)
		if ok {
			ec, err := strconv.ParseInt(rspBaseBased.GetRspBase().Code, 10, 64)
			if err != nil {
				errCode = 0
			} else {
				errCode = int(ec)
			}
		} else {
			// 成功
			errCode = 0
		}

		// 回包
		ginPanicCall(func() {
			ginCtx.JSON(http.StatusOK, rsp)
		})
	}

	return f
}

// RegPostHandlerRegPostHandler 注册post handler，只应该在 init 的时候传入，如果输入非法，则直接 panic
func RegPostHandler(api string, handlerFunc interface{}, opt PostOpt) {
	// 禁止重复注册
	if _, ok := postHandlerMap[api]; ok {
		panic(fmt.Sprintf("duplicate api [%s]", api))
	}

	fv := reflect.ValueOf(handlerFunc)
	// handler不能是nil
	if !fv.IsValid() || fv.IsNil() {
		panic("handler is nil")
	}
	ft := fv.Type()
	argCount := ft.NumIn()
	// 只能有一到三个参数且参数表定长
	if argCount < 1 || argCount > 3 || ft.IsVariadic() {
		panic("arg count error")
	}
	// 第一参数必须是*Ctx
	if ft.In(0) != reflect.TypeOf((*APIContext)(nil)) {
		panic("arg #1 is not *APIContext")
	}
	reqType := reflect.Type(nil)
	rspType := reflect.Type(nil)
	argCheckFunc := func(idx int, isRsp bool) (struType reflect.Type) {
		argType := ft.In(idx)
		if argType.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("arg #%d is not a pointer", idx+1))
		}
		struType = argType.Elem()
		if struType.Kind() != reflect.Struct {
			panic(fmt.Sprintf("arg #%d is not a pointer to struct", idx+1))
		}
		if isRsp {
			// 如果是 RSP_STRU，还要检查下是否实现了 proto.RspBaseBased 接口
			_, ok := reflect.New(struType).Interface().(RspBaseBased)
			if !ok {
				panic(fmt.Sprintf("arg #%d is not a implementation of proto.RspBaseBased", idx+1))
			}
		}

		return
	}

	switch ft.NumOut() {
	case 1:
		switch argCount {
		case 1:
			// 前面检查过了
		case 2:
			// 第二参数 *REQ_STRU
			reqType = argCheckFunc(1, false)
		default:
			panic("arg count error")
		}
	case 0:
		switch argCount {
		case 2:
			// 第二参数 *RSP_STRU
			rspType = argCheckFunc(1, true)
		case 3:
			// 第二参数 *REQ_STRU
			reqType = argCheckFunc(1, false)
			// 第三参数 *RSP_STRU
			rspType = argCheckFunc(2, true)
		default:
			panic("arg count error")
		}
	default:
		panic("return value count error")
	}

	//注册
	postHandlerMap[api] = &postHandler{
		funcValue: fv,
		reqType:   reqType,
		rspType:   rspType,
		opt:       opt,
	}
}

func InitHandler(r *gin.Engine) {
	for api, handler := range postHandlerMap {
		r.POST(api, handler.mkCinexHandler())
	}
}
