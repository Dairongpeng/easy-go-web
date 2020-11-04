package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	// 默认路由
	engine := gin.Default()

	engine.GET("/json", func(c *gin.Context) {

		// 方法一：定义一个map，key为string，value为一个空接口，可以接受任意的值类型
		// 结尾花括号和最后一行不在同一行，需要在最后一行的代码后加上','
		//data := map[string]interface{} {
		//	"name" : "xiaodai",
		//	"age" : 18,
		//}

		// 方法二：借助gin封装好的值为string，value为空接口的结构
		data := gin.H{"name": "xiaodai", "age": 18}

		// 返回状态码
		c.JSON(http.StatusOK, data)

	})

	engine.Run(":8000")

}
