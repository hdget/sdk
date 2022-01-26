package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/hdget/hdsdk/lib/err"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type PageResponse struct {
	Response
	Total int64 `json:"total"`
}

const (
	_                     = err.ErrCodeModuleRoot + iota
	ErrCodeServerInternal // 内部错误
	ErrCodeUnauthorized   // 未授权
	ErrCodeInvalidRequest // 非法请求
	ErrCodeForbidden      // 拒绝访问

)

var (
	errInvalidRequest = err.New("invalid request", ErrCodeInvalidRequest)
	errForbidden      = err.New("forbidden", ErrCodeForbidden)
	errUnauthorized   = err.New("unauthorized", ErrCodeUnauthorized)
)

// Success respond with data
// empty args respond with 'ok' message
// args[0] is the response data
func Success(c *gin.Context, args ...interface{}) {
	var ret Response
	switch len(args) {
	case 0:
		ret.Data = "ok"
	case 1:
		ret.Data = args[0]
	}
	c.PureJSON(http.StatusOK, ret)
}

// SuccessPages respond with pagination data
func SuccessPages(c *gin.Context, total int64, pages interface{}) {
	c.PureJSON(http.StatusOK, PageResponse{
		Response: Response{
			Data: pages,
		},
		Total: total,
	})
}

func Failure(c *gin.Context, err error) {
	c.PureJSON(http.StatusInternalServerError, err2response(err))
}

func Redirect(c *gin.Context, location string) {
	c.Redirect(http.StatusFound, location)
}

func PermanentRedirect(c *gin.Context, location string) {
	c.Redirect(http.StatusMovedPermanently, location)
}

func InvalidRequest(c *gin.Context) {
	c.PureJSON(http.StatusBadRequest, err2response(errInvalidRequest))
}

func Forbidden(c *gin.Context) {
	c.PureJSON(http.StatusForbidden, err2response(errForbidden))
}

func Unauthorized(c *gin.Context) {
	c.PureJSON(http.StatusUnauthorized, err2response(errUnauthorized))
}

func err2response(e error) *Response {
	code := ErrCodeServerInternal
	codeErr, ok := e.(*err.CodeError)
	if ok {
		code = codeErr.Code()
	}
	return &Response{
		Msg:  e.Error(),
		Code: code,
	}
}