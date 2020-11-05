package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
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

	// 接口6 获取路径参数 127.0.0.1:8000/get-name/xiaodai/18
	engine.GET("/get-name/:username/:age", func(c *gin.Context) {

		// 获取路径参数
		username := c.Param("username")
		age := c.Param("age")

		data := gin.H{"name": username, "age": age}

		// 返回状态码
		c.JSON(http.StatusOK, data)

	})

	// 接口7 请求参数和自定义结构体进行绑定。POST 127.0.0.1:8000/login-binding {"username":"xiaodai", "password":"12345"}
	type UserInfo struct {
		Username string `json:"username"` // 首字母大写表示可以导出。tag表示导出后显示小写
		Password string
	}
	engine.POST("/login-binding", func(c *gin.Context) {
		var u UserInfo // 声明一个userInfo类型的结构体类型变量
		// 把请求参数绑定到我们的u类型上，注意，这里是为了修改u，所以不能传递值，需要传递引用。
		// 且变量首字母需要大写，小写表示不可导出。返回不到前端
		// c.ShouldBind(&u)不仅可以绑定post的body的参数，也可以绑定get请求的请求头参数:?username=a&password=12345。也可以绑定json的参数。按照字段名称进行绑定
		err := c.ShouldBind(&u)
		if err != nil { // 绑定失败则返回502错误信息
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err.Error(),
			})
		} else {
			// 返回状态码
			c.JSON(http.StatusOK, u)
		}
	})

	// 接口7 上传文件
	engine.LoadHTMLFiles("./upload.html")
	engine.GET("/file-page", func(c *gin.Context) {
		// 返回状态码
		c.HTML(http.StatusOK, "upload.html", nil)

	})

	engine.POST("/upload", func(c *gin.Context) {
		// 通过文件名从请求中接受文件，将接受到的文件保存在服务器本地
		file, err := c.FormFile("file")

		// 多文件上传。保存到files，循环files逐个保存
		//form, _ := c.MultipartForm()
		//files := form.File["file"]

		if err != nil { // 接受失败，返回错误
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err.Error(),
			})
		} else {
			// 将文件保存到本地
			p := path.Join("./%s", file.Filename)
			_ = c.SaveUploadedFile(file, p) // 可以通过设置上传文件缓冲区大小，来增加上传速度: engine.MaxMultipartMemory = 8 << 20
			c.JSON(http.StatusOK, gin.H{
				"ok": "OK",
			})
		}
	})

	engine.Run(":8000")

}
