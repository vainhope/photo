package config

import (
	"github.com/gin-gonic/gin"
	"goweb/base"
	"log"
	"net/http"
	"runtime"
)

func InitServer(r *gin.Engine) *gin.Engine {
	r.Use(LogSetUp(), ExceptionSetUp())
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "请求方法不存在",
			"data":    nil,
		})
		return
	})
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "请求方法不存在",
			"data":    nil,
		})
		return
	})
	return r
}

/**
全局异常处理
*/
func ExceptionSetUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var apiException *base.DataError
				h, ok := err.(*base.DataError)
				if ok {
					apiException = h
				} else {
					apiException = &base.DataError{500, "服务器错误"}
				}
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				logger := c.MustGet("log").(*log.Logger)
				logger.Printf("server error ------> %s", string(buf[:n]))
				c.JSON(http.StatusOK, gin.H{
					"code":    apiException.Code,
					"message": apiException.Message,
					"data":    nil,
				})
				if ok {
					return
				}
			}
		}()
		c.Next()
	}
}
