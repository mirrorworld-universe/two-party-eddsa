package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"main/global"
	"main/internal/base_resp"
)

type newWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (nw newWriter) Write(b []byte) (int, error) {
	nw.bodyBuf.Write(b)
	return nw.ResponseWriter.Write(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		nw := newWriter{
			ResponseWriter: c.Writer,
			bodyBuf:        bytes.NewBufferString(""),
		}
		httpMethod := c.Request.Method
		httpPath := c.Request.URL.Path
		//httpHeader := c.Request.Header
		//SSO_Token := httpHeader["X-Sso-Token"]

		realIp := c.Request.Header.Get("X-Real-Ip")
		httpClientIP := c.Request.RemoteAddr

		reqBody := []byte{}
		if c.Request.Body != nil {
			reqBody, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		}

		// 替换默认的 writer
		c.Writer = nw
		c.Next()

		body := nw.bodyBuf.Bytes()
		resp := base_resp.Response{}
		json.Unmarshal(body, &resp)

		// add user address to log info
		//user, _ := c.Get(global.UserAddress)

		fields := logrus.Fields{
			"RequestPath": httpPath,
			"requestBody": reqBody,
			"httpMethod":  httpMethod,
			//"SSO-Token":       SSO_Token,
			//"user":            user,
			"RequestClientIP": httpClientIP,
			"X-Real-Ip":       realIp,
			"httpStatusCode":  c.Writer.Status(),
			"trace_id":        c.Value("trace_id"),
		}
		global.Logger.WithFields(fields).Info(fmt.Sprintf("%+v", resp))
	}
}
