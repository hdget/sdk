package lib_ws

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hdget/utils"
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

func Error(c *gin.Context, code int, msg string) {
	ret := &Response{
		Code: code,
		Msg:  msg,
	}
	c.PureJSON(http.StatusOK, ret)
}

// SuccessRaw respond with raw data
func SuccessRaw(c *gin.Context, result interface{}) {
	var content string
	switch t := result.(type) {
	case string:
		content = t
	case []byte:
		content = utils.BytesToString(t)
	default:
		v, _ := json.Marshal(result)
		content = utils.BytesToString(v)
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Type", "application/json; charset=utf-8")
	_, _ = c.Writer.WriteString(content)
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

// Failure grpc http错误
func Failure(c *gin.Context, err error) {
	c.PureJSON(http.StatusInternalServerError, err)
}

func InvalidRequest(c *gin.Context, err error) {
	c.PureJSON(http.StatusBadRequest, err)
}

func Forbidden(c *gin.Context, err error) {
	c.PureJSON(http.StatusForbidden, err)
}

func Unauthorized(c *gin.Context, err error) {
	c.PureJSON(http.StatusUnauthorized, err)
}

func Redirect(c *gin.Context, location string) {
	c.Redirect(http.StatusFound, location)
}

func PermanentRedirect(c *gin.Context, location string) {
	c.Redirect(http.StatusMovedPermanently, location)
}
