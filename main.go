package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
)

// 定义一个全局的db
var (
	DB *gorm.DB
)

// paper module
type paper struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

// 初始化mysql连接
func initMysql() (err error) {
	dsn := "root:abc123@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	// =是赋值全局的DB,:=是声明的新的变量
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})
	if err != nil {
		return err
	}
	return
}

func main() {

	// 连接数据库
	err := initMysql()
	if err != nil {
		panic(err)
	}

	// 模型绑定
	_ = DB.AutoMigrate(&paper{})

	// gin默认路由
	engine := gin.Default()

	// v1路由组
	v1Group := engine.Group("/v1")

	// 具体路由的代码块
	{
		// 添加
		v1Group.POST("/add", func(context *gin.Context) {
			var p paper
			// 把传过来的参数，映射到paper的模型上
			context.BindJSON(&p)

			if err = DB.Debug().Create(&p).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				context.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "success",
					"data": p,
				})
			}
		})

		// 查看
		v1Group.GET("/get", func(context *gin.Context) {

		})

		// 修改
		v1Group.POST("/update", func(context *gin.Context) {

		})

		// 删除
		v1Group.GET("/delete", func(context *gin.Context) {

		})
	}

	// 暴露端口
	_ = engine.Run(":8000")
}
