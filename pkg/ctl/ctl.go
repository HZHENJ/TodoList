package ctl

import (
	"net/http"
	"to-do-list/pkg/e"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Error string      `json:"error"`
}

type DataList struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type Wrapper struct {
	C *gin.Context
}

func NewWrapper(c *gin.Context) *Wrapper {
	return &Wrapper{C: c}
}

// Response 基础响应方法
func (w *Wrapper) Response(httpCode, errorCode int, data interface{}) {
	w.C.JSON(httpCode, Response{
		Code: errorCode,
		Msg:  e.GetMsg(errorCode),
		Data: data,
	})
}

// Success 成功
func (w *Wrapper) Success(data interface{}) {
	w.Response(http.StatusOK, e.SUCCESS, data)
}

// ResponseList 分页列表响应
// list: 数据切片, total: 总数
func (w *Wrapper) ResponseList(list interface{}, total int64) {
	w.Response(http.StatusOK, e.SUCCESS, DataList{
		List:  list,
		Total: total,
	})
}

// Error 错误响应
func (w *Wrapper) Error(errCode int, err error) {
	// 只有在某种Debug开关开启时，才把 err.Error() 放入响应
	// 平时只返回 errCode 对应的 Msg
	msg := e.GetMsg(errCode)

	// 如果想要 track_id 或者日志打印，在这里进行 log.Println(err)

	w.C.JSON(http.StatusOK, Response{
		Code: errCode,
		Msg:  msg,
		Data: nil,
		// Error: err.Error(), // 生产环境建议注释掉这行，或者加开关判断
	})
}
