package binding

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"main/internal/base_resp"
	error_code "main/internal/err_code"
)

func BindJson(c *gin.Context, params interface{}) error {
	bsp := error_code.NewBaseResp()

	err := c.ShouldBindJSON(&params)
	if err != nil {
		// 参数校验失败之后，处理校验错误
		transError, ok := err.(validator.ValidationErrors)
		if !ok {
			bsp.SetMsg(error_code.ParamsError, err.Error())
		} else {
			errMsg := transError.Error()
			bsp.SetMsg(error_code.ParamsError, fmt.Sprintf("%v", errMsg))
		}
		base_resp.JsonResponse(c, bsp, nil)
		return err
	}
	return nil
}
