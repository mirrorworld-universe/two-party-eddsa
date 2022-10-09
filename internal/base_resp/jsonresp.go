package base_resp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	error_code "main/internal/err_code"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceId interface{} `json:"trace_id"`
}

func (r Response) String() string {
	return fmt.Sprintf("code: %v, msg: %v, data: %v", r.Code, r.Msg, r.Data)
}

func JsonResponse(c *gin.Context, bsp *error_code.BaseResp, data interface{}) {
	traceId := c.Value("trace_id")
	resp := Response{
		Code:    bsp.Code,
		Msg:     bsp.MSg,
		Data:    data,
		TraceId: traceId,
	}
	c.JSON(http.StatusOK, resp)
}
