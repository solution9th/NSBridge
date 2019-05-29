package utils

import "runtime"

const (
	ErrSuccess int = 0
)

var (
	ErrCode = map[int]string{
		0: "success",
	}
)

// AddErrCodes 添加新的错误码
func AddErrCodes(codes map[int]string) {
	for k, v := range codes {
		ErrCode[k] = v
	}
}

type Response struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data"`
}

// ParseSuccessWithData return success whith data
func ParseSuccessWithData(data interface{}) Response {
	return ParseResult(ErrSuccess, "", data)
}

// ParseSuccess return success withon data
func ParseSuccess() Response {
	return ParseResult(ErrSuccess, "", "success")
}

// ParseResult parse new error
func ParseResult(errCode int, msgAdd string, data interface{}) Response {
	if errCode != 0 {
		pc, _, _, _ := runtime.Caller(1)
		Error(runtime.FuncForPC(pc).Name()+":", ErrCode[errCode]+msgAdd, data)
		if _, ok := ErrCode[errCode]; !ok {
			return Response{
				ErrCode: errCode,
				ErrMsg:  msgAdd,
				Data:    data,
			}
		}
	}

	return Response{
		ErrCode: errCode,
		ErrMsg:  ErrCode[errCode] + " " + msgAdd,
		Data:    data,
	}
}
