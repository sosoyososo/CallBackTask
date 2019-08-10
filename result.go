package main

type Result struct {
	Succeed bool
	Code    ErrCode
	Msg     string
	Data    interface{}
}

func SucceedResult(data interface{}) *Result {
	return &Result{Succeed: true, Code: ErrCodeNone, Data: data, Msg: "成功"}
}

func ErrResult(code ErrCode, msg string) *Result {
	r := Result{Succeed: false, Code: code, Data: nil, Msg: "请求失败"}
	if len(msg) > 0 {
		r.Msg = msg
	}
	return &r
}

func RedisErrResult() *Result {
	return ErrResult(ErrCodeRedisErr, "Redis操作错误")
}

func DBErrResult() *Result {
	return ErrResult(ErrCodeMysqlErr, "Mysql操作错误")
}

func ParameterErrResult() *Result {
	return ErrResult(ErrCodeWrongParameter, "参数错误")
}
