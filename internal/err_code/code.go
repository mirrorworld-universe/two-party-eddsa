package error_code

const (
	Success       = 200
	InternalError = 500
	LogicalError  = 2 // 逻辑错误
	ParamsError   = 3 // 参数错误
	ThirdError    = 4 // 第三方调用错误
	OtherError    = 5 // 其他错误
	SignError     = 6 // 签名错误
)

const (
	SuccessCode = 0
	SuccessMsg  = "success"
)

type BaseResp struct {
	Code int
	MSg  string
}

func NewBaseResp() *BaseResp {
	return &BaseResp{
		Code: SuccessCode,
		MSg:  SuccessMsg,
	}
}

func (e *BaseResp) SetMsg(code int, msg string) {
	e.Code = code
	e.MSg = msg
}
