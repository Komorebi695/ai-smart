package model

import (
	"fmt"
)

const (
	ErrOK          = 200
	ErrParam       = 201
	ErrOut         = 202
	ErrorWithToast = 201

	ErrNeedLogin = 3001

	ErrSystem = 5001
)

var (
	errMsg = map[int]string{
		ErrOK:    "success",
		ErrParam: "参数错误",
		ErrOut:   "Error out", // 重新登录处理错误
	}
)

var (
	ParamErrRsp = BaseResponse{
		Code: ErrorWithToast,
		Msg:  "参数错误",
	}

	NeedLoginRsp = BaseResponse{
		Code: ErrNeedLogin,
		Msg:  "请先登录",
	}

	ErrSystemRsp = BaseResponse{
		Code: ErrSystem,
		Msg:  "system busy",
	}
)

func (rsp *BaseResponse) WithStatus(code int) *BaseResponse {
	rsp.WithCode(code).WithMsg(errMsg[code])
	if rsp.Msg == "" {
		rsp.WithMsg("fail")
	}
	return rsp
}

func (rsp *BaseResponse) WithOK() *BaseResponse {
	return rsp.WithStatus(ErrOK)
}

func (rsp *BaseResponse) Error(msg interface{}) *BaseResponse {
	rsp.Msg = fmt.Sprintf("%v", msg)
	return rsp.WithCode(201)
}
