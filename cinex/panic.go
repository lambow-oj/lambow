package cinex

type ginPanicType struct {
	err interface{}
}

// 调用传入的func，对于其panic的错误转成ginPanic类型重新抛出
func ginPanicCall(f func()) {
	defer func() {
		if err := recover(); err != nil {
			panic(ginPanicType{
				err: err,
			})
		}
	}()
	f()
}
