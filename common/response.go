package common

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

const (
	SuccessCode int = 0
	ErrorCode   int = -1
)

// RespResult 响应通用结构
type RespResult struct {
	Code    int         `json:"code"` //错误码(0:正常,非零错误)
	Message string      `json:"msg"`  //提示信息
	Data    interface{} `json:"data"` //返回数据
}

//DataResult 数据结构
type DataResult struct {
	Total int         `json:"total"` //数据量
	Rows  interface{} `json:"rows"`  //数据集合
}

func buildResp(r *ghttp.Request, code int, message string, data interface{}) {
	respData := interface{}(nil)
	if data != nil {
		respData = data
	}
	resp := RespResult{
		Code:    code,
		Message: message,
		Data:    respData,
	}
	if err := r.Response.WriteJsonExit(resp); err != nil {
		g.Log().Error(err)
	}
}

func BuildRespSuccess(r *ghttp.Request, data interface{}) {
	buildResp(r, SuccessCode, "成功", data)
}

func BuildRespWithError(r *ghttp.Request, data interface{}, err error) {
	if err != nil {
		BuildRespFailWithError(r, err)
		return
	}
	BuildRespSuccess(r, data)
}

func BuildPagingResult(r *ghttp.Request, total int, rows interface{}) {
	BuildRespSuccess(r, &DataResult{
		Total: total,
		Rows:  rows,
	})
}

func BuildPagingResultWithError(r *ghttp.Request, total int, rows interface{}, err error) {
	if err != nil {
		BuildRespFailWithError(r, err)
		return
	}
	BuildPagingResult(r, total, rows)
}

func BuildRespFail(r *ghttp.Request, err string) {
	BuildRespFailWithCode(r, ErrorCode, err)
}

func BuildRespFailWithCode(r *ghttp.Request, code int, err string) {
	if code == SuccessCode {
		BuildRespSuccess(r, nil)
		return
	}
	buildResp(r, code, err, nil)
}

func BuildRespFailWithError(r *ghttp.Request, err error) {
	BuildRespFailWithCode(r, ErrorCode, err.Error())
}
