package model

import "fmt"

type (
	BaseResponseInterface interface {
		GetCode() int
		GetMsg() string
		GetErr() error
		GetLogID() string
	}

	BaseResponse struct {
		Code  int    `format:"int" json:"code"`
		Msg   string `json:"msg"`
		LogID string `json:"log_id"`
		Err   error  `json:"err"`
	}

	BasePageReq struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}

	BaseHeaderReq struct {
		//
	}

	BaseListRsp struct {
		Total  int64 `json:"total"`
		IsNext bool  `json:"is_next"`
	}

	BasePageInfo struct {
		Page     int   `json:"page"`
		PageSize int   `json:"page_size"`
		Total    int64 `json:"total"`
		IsNext   bool  `json:"is_next"`
	}
)

func (r *BaseResponse) WithCode(code int) *BaseResponse {
	r.Code = code
	return r
}

func (r *BaseResponse) WithMsg(msg interface{}) *BaseResponse {
	r.Msg = fmt.Sprintf("%v", msg)
	return r
}

func (r *BaseResponse) WithErr(err error) *BaseResponse {
	r.Err = err
	return r
}

func (r *BaseResponse) GetCode() int {
	return r.Code
}

func (r *BaseResponse) GetMsg() string {
	return r.Msg
}

func (r *BaseResponse) GetLogID() string {
	return r.LogID
}

func (r *BaseResponse) GetErr() error {
	return r.Err
}
