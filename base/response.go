package base

import (
	"github.com/gin-gonic/gin"
	"goweb/entity"
	"net/http"
)

func Response(context *gin.Context, code DataCode, data interface{}) {
	context.JSON(http.StatusOK, gin.H{
		"code":        int(code),
		"message":     code.String(),
		"data":        data,
	})
}


func ResponseWithPage(context *gin.Context, code DataCode, data interface{}, page entity.Page) {
	context.JSON(http.StatusOK, gin.H{
		"code":        int(code),
		"message":     code.String(),
		"data":        data,
		"currentPage": page.CurrentPage,
		"pageSize":    page.PageSize,
	})
}

func ResponseWithMessage(context *gin.Context, code DataCode, message string, data interface{}) {
	context.JSON(http.StatusOK, gin.H{
		"code":    int(code),
		"message": message,
		"data":    data,
	})
}
