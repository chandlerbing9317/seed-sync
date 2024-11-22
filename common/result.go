package common

// Result 结构体，使用泛型来支持不同类型的data
type Result[T any] struct {
	Success bool   `json:"success"` // 是否成功
	Msg     string `json:"msg"`     // 失败原因
	Data    T      `json:"data"`    // 返回数据
}

// Success 生成一个成功的结果
func SuccessResult[T any](data T) Result[T] {
	return Result[T]{
		Success: true,
		Data:    data,
	}
}

// Fail 生成一个失败的结果
func FailResult(msg string) Result[any] {
	return Result[any]{
		Success: false,
		Msg:     msg,
	}
}
