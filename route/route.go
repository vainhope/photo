package route

import (
	"github.com/gin-gonic/gin"
	"goweb/controller/cover"
	"goweb/controller/file"
	"goweb/controller/photo"
	"goweb/controller/user"
	"goweb/interceptor"
)

func Router(route *gin.Engine) {

	//前后拦截器
	route.Use(interceptor.Interceptor())

	group := route.Group("/vain/photo")
	group.POST("/user/login", user.Login)
	group.POST("/user/register", user.Register)

	group.POST("/upload", file.Upload)
	group.POST("/cover/add", cover.Add)
	group.POST("/cover/pageQuery", cover.PageQuery)
	group.GET("/coverCaptures", photo.CoverCaptures)
	group.POST("/addCaptures", photo.AddCaptures)

}
