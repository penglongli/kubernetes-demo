package result

import "github.com/gin-gonic/gin"

type ErrorCode string

const (
	SUCCESS = "0"
	ERROR   = "1"
)

type Response struct {
	Code ErrorCode   `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Failed(ctx *gin.Context, code ErrorCode, errMsg string) {
	ctx.JSON(200, &Response{
		Code: code,
		Msg:  errMsg,
	})
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, &Response{
		Code: SUCCESS,
		Data: data,
	})
}
