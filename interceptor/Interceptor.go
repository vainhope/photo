package interceptor

import (
	"github.com/gin-gonic/gin"
	"goweb/base"
	"goweb/utils"
	"log"
	"time"
)

var notNeedLogin = map[string]int{
	"/vain/photo/user/login":      1,
	"/vain/photo/user/register":   1,
	"/vain/photo/upload":          0,
	"/vain/photo/cover/pageQuery": 1,
	"/vain/photo/coverCaptures":   1}

func Interceptor() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		header := context.GetHeader("Token")
		path := context.Request.URL.Path
		if notNeedLogin[path] != 1 {
			if "" == header {
				base.Response(context, base.NOT_LOGIN, nil)
				context.Abort()
				return
			}
			if claim, ok := utils.ParseToken(header); ok != nil {
				base.Response(context, base.INVALID_TOKEN, nil)
				context.Abort()
				return
			} else {
				context.Set("userId", claim.Id)
			}
		}
		log.Println("请求前 " + path)
		context.Next()
		latency := time.Since(t)
		log.Printf("请求后 %s %s ", path, latency)
	}
}
