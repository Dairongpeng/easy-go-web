package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	// 默认路由
	engine := gin.Default()

	// 接口1 json返回格式测试
	engine.GET("/json-test", func(c *gin.Context) {

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

	// 方法二：结构体。灵活使用tag,也就是字段后``标识，对字段做定制化操作
	type msg struct {
		Name string
		// 当用json进行解析该结构体时，变量转为小写导出给前端。
		Message string `json:"message"`
		Age     int
	}
	// 接口2 结构体数据返回json
	engine.GET("/struct-test", func(c *gin.Context) {
		data := msg{
			// 注意，当字段首字母小写时，该字段属于不可导出，json序列化使用的是反射，那么小写字母开头的变量无法反射获取。该字段返回不到前端。
			// 想要返回给前端字母小写的变量，可以参考上文，在结构体导出是，json解析字段为小写
			Name:    "xiaodai",
			Message: "hello go",
			Age:     18,
		}
		c.JSON(http.StatusOK, data) // json序列化结构体data
	})

	// 接口3 前端传?query=hello方式交互。接口为127.0.0.1:8000/query-test?query="hello"
	engine.GET("/query-test", func(c *gin.Context) {

		// 获取key为query的参数值
		// query := c.Query("query")

		// 获取key为query的参数值, 没找到则用默认值
		// c.DefaultQuery("query", "nothing")

		// 获取key为query的参数值, 取不到为false
		//query, ok := c.GetQuery("query")
		//
		//if !ok {
		//	query = "没取到值"
		//}

		// 多个query类型的参数，例如?query=abc&age=18 可以用多个key获取
		query := c.Query("query")
		age := c.Query("age")

		c.JSON(http.StatusOK, gin.H{
			"name": query,
			"age":  age,
		})

	})

	// 接口4 我们先用get请求返回给前端一个表单
	engine.LoadHTMLFiles("./login.html", "./index.html")
	engine.GET("/login", func(c *gin.Context) {

		c.HTML(http.StatusOK, "login.html", nil)

	})

	// 接口5 post请求获取请求体的数据, 不使用postman来模仿post请求，用接口4返回给前端的表单来触发post
	engine.POST("/login", func(c *gin.Context) {

		// 方式1获取form表单的body数据
		//username := c.PostForm("username")
		//password := c.PostForm("password")

		// 方式2获取form表单的body数据,没有则默认
		//username := c.DefaultPostForm("username", "tom")
		//password := c.DefaultPostForm("password", "abc123")

		// 方式3获取form表单的body数据,没有则默认
		username, ok := c.GetPostForm("username")
		if !ok {
			username = "tom"
		}
		password, ok := c.GetPostForm("password")
		if !ok {
			password = "***"
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Name":     username,
			"Password": password,
		})

	})

	engine.Run(":8000")

}
