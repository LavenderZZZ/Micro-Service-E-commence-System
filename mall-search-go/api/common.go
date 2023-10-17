package api

type CommonResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResultCode struct {
	SUCCESS         int
	FAILED          int
	VALIDATE_FAILED int
	UNAUTHORIZED    int
	FORBIDDEN       int
	// Add other result codes if needed...
}

var resultCode = ResultCode{
	SUCCESS:         200,
	FAILED:          500,
	VALIDATE_FAILED: 400,
	UNAUTHORIZED:    401,
	FORBIDDEN:       403,
}

func Success(data interface{}) *CommonResult {
	return &CommonResult{
		Code:    resultCode.SUCCESS,
		Message: "Success",
		Data:    data,
	}
}

func SuccessWithMessage(data interface{}, message string) *CommonResult {
	return &CommonResult{
		Code:    resultCode.SUCCESS,
		Message: message,
		Data:    data,
	}
}

func Failed(message string) *CommonResult {
	return &CommonResult{
		Code:    resultCode.FAILED,
		Message: message,
		Data:    nil,
	}
}

func ValidateFailed(message string) *CommonResult {
	return &CommonResult{
		Code:    resultCode.VALIDATE_FAILED,
		Message: message,
		Data:    nil,
	}
}

func Unauthorized(data interface{}) *CommonResult {
	return &CommonResult{
		Code:    resultCode.UNAUTHORIZED,
		Message: "Unauthorized",
		Data:    data,
	}
}

func Forbidden(data interface{}) *CommonResult {
	return &CommonResult{
		Code:    resultCode.FORBIDDEN,
		Message: "Forbidden",
		Data:    data,
	}
}
