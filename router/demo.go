package router

import (
	v1 "easy-go-web/api/v1"
	"github.com/gin-gonic/gin"
)

// 消息机器路由
func InitDemoRouter(r *gin.RouterGroup) (R gin.IRoutes) {
	router1 := r.Group("/demo")
	{
		router1.GET("/list", v1.ListDemo)
		router1.GET("/single", v1.Single)
		router1.POST("/add", v1.Add)
	}
	return r
}
