package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"time"
)

// 案例11的中间件函数
func m1(c *gin.Context) {
	fmt.Println("中间件拦截请求=====")
	// 计时功能(统计请求函数的耗时)
	start := time.Now()
	c.Next() // 调用后续的处理函数，这里的后续处理函数是indexHandler
	//c.Abort() 阻止调用后续函数。常用来拦截阻止后续操作
	// return 加入return 停止该中间件，后续函数处理结束不会再回到该中间件函数
	cost := time.Since(start)
	fmt.Printf("cost:%v\n", cost)
}

// 案例11的全局中间件函数m2
func m2(c *gin.Context) {
	// 是否是登录用户
	// if 是登录用户
	// c.Next()
	// 不是登录用户
	// c.Abort()
}

func indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"test": "ok"})
}

func main() {

	// 默认路由
	engine := gin.Default()

	// 案例1 json返回格式测试
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

	// 案例2 结构体数据返回json
	// 方法二：结构体。灵活使用tag,也就是字段后``标识，对字段做定制化操作
	type msg struct {
		Name string
		// 当用json进行解析该结构体时，变量转为小写导出给前端。
		Message string `json:"message"`
		Age     int
	}
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

	// 案例3 前端传?query=hello方式交互。接口为127.0.0.1:8000/query-test?query="hello"
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

	// 案例4 我们先用get请求返回给前端一个表单
	engine.LoadHTMLFiles("./login.html", "./index.html")
	engine.GET("/login", func(c *gin.Context) {

		c.HTML(http.StatusOK, "login.html", nil)

	})

	// 案例5 post请求获取请求体的数据, 不使用postman来模仿post请求，用接口4返回给前端的表单来触发post
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

	// 案例6 获取路径参数 127.0.0.1:8000/get-name/xiaodai/18
	engine.GET("/get-name/:username/:age", func(c *gin.Context) {

		// 获取路径参数
		username := c.Param("username")
		age := c.Param("age")

		data := gin.H{"name": username, "age": age}

		// 返回状态码
		c.JSON(http.StatusOK, data)

	})

	// 案例7 请求参数和自定义结构体进行绑定。POST 127.0.0.1:8000/login-binding {"username":"xiaodai", "password":"12345"}
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

	// 案例8 上传文件
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

	// 案例9 接口重定向
	// 跳转到别的网站
	engine.GET("/redirect", func(c *gin.Context) {

		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")

	})
	// 站内跳转
	engine.GET("/redirect2", func(c *gin.Context) {

		c.Request.URL.Path = "/struct-test" // 把请求的url修改为"/struct-test"
		engine.HandleContext(c)             // 用重定向后的接口继续处理该次请求

	})

	// 案例10 gin框架中的路由和路由组的概念。
	// 普通的路由，比如一个请求路径我们对应一个处理函数。常用的GET、POST、PUT、DELETE等
	// 映射所有请求的路由，比如Any。点进入Any方法，可以看到它封装了各种路由请求，既能处理GET也能处理POST。
	// 比如下面"/router-test"方法，我们用get请求，和post请求都能识别，也能进入到相应的case分支中去
	// go语言自带break，每个分支中无需break
	engine.Any("/router-test", func(c *gin.Context) {

		switch c.Request.Method {

		case "GET":
			c.JSON(http.StatusOK, gin.H{"method": "GET"})
		case http.MethodPost:
			c.JSON(http.StatusOK, gin.H{"method": "POST"})

		}
	})
	// 默认路由，当请求服务端不存在的路由地址的时候，默认进入该路由
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "404"})
	})

	// 路由组的概念，处理路由的前缀相同的请求，进行归类。比如"/user/hello"和"/user/login"同属于"/user"前缀的路由组下
	// 把公用的前缀提取出来，创建公共的路由组
	shopGroup := engine.Group("/shop")

	{
		// 通过路由组GET,实质上url是带上路由组的路径：127.0.0.1:8000/shop/list 和 127.0.0.1:8000/shop/get
		shopGroup.GET("/list", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"apple": "10元", "brash": "12.66元"})
		})

		shopGroup.GET("/get", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"apple": "10元"})
		})

		// 路由组也支持嵌套。用shopGroup("/shop") 嵌套一个新的shopBusGroup("/bus")路由组
		shopBusGroup := shopGroup.Group("/bus")
		// "/shop/bus/view"
		shopBusGroup.GET("/view", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"view": "bingo!"})
		})
	}

	// 案例11 中间件，类似于钩子函数，类比于java的切面，拦截器等。golang的中间件函数适合处理一些多个接口公共的功能
	// 例如处理登录认证，权限检验，数据分页，记录日志，耗时统计等
	// 例如我们网站有"/index","/user","/shop"等路由组前缀
	// 1. 定义中间间函数，上文main函数外
	// 2. 定义处理函数，上文main函数外
	// 3.为接口添加中间件函数m1
	engine.GET("/index", m1, indexHandler)
	// 定义全局注册函数m2。此时不管什么接口，都会被该中间件函数拦截。部分有共性的接口使用一个中间件，可以运用路由组。
	// 可以engine.Use(m1,m2)注册多个全局中间件。多个全局中间件支持重入。例如m1和m2中间件，m1拦截->执行后续处理函数->m2->执行后续处理函数->index->m2结束->m1结束
	// engine.Use(m2)
	// c := gin.Default() 中默认包含了log和recovery两个中间件。如果想要一个不包含任何中间件的服务端，可以使用gin.New()
	// 其中log是日志中间件，服务启动或者报错会打印日志；recovery是错误处理，系统报错gin为cover住错误，报相应的状态码，例如500

	engine.Run(":8000")

}
