package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"
	"randomimg/Service"
	"time"
)

var (
	randomLen int
	DB        *sql.DB
	dbLinkErr error
)

// 刷新函数，用于自动重载图片库
func refreshLen() {
	for true {
		queryRandomLenErr := DB.QueryRow("SELECT COUNT(PID) FROM random_photo_table").Scan(&randomLen)
		if queryRandomLenErr != nil {
			fmt.Println("获取随机图片库失败")
			fmt.Println(queryRandomLenErr)
			return
		}
		randomLen--
		fmt.Println("刷新图片库长度")
		time.Sleep(time.Second * 300)

	}

}
func init() {
	// 读取配置文件
	viper.SetConfigFile("config/config.yml")
	readConfigErr := viper.ReadInConfig()
	if readConfigErr != nil {
		fmt.Println("读取配置文件失效")
		fmt.Println(readConfigErr)
		panic("STOP")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", viper.GetString("config.database.username"), viper.GetString("config.database.password"), viper.GetString("config.database.host"), viper.GetString("config.database.dbname"))
	DB, dbLinkErr = sql.Open("mysql", dsn)
	if dbLinkErr != nil || DB.Ping() != nil {
		fmt.Println("数据库初始化失败，无法连接，终止系统")
		panic(dbLinkErr)
	}
	fmt.Println("初始化成功")
	fmt.Println(DB.Ping())

	queryRandomLenErr := DB.QueryRow("SELECT COUNT(PID) FROM random_photo_table").Scan(&randomLen)
	if queryRandomLenErr != nil {
		fmt.Println("获取随机图片库失败")
		panic(queryRandomLenErr)
	}
}

func main() {
	server := gin.Default()
	server.GET("/random/json", func(context *gin.Context) {
		token := context.Query("token")
		if token == "" {
			context.JSON(200, gin.H{
				"code": "402",
				"err":  "认证失败,如需使用，请前往https://shop.loli.fit购买服务",
			})
			return
		}
		apiAuthType, randomNum := Service.ApiAuth(DB, token, "random_pro")
		if apiAuthType != true || randomNum < 1 {
			context.JSON(200, gin.H{
				"code": "403",
				"err":  "认证失败,如需使用，请前往https://shop.loli.fit购买服务",
			})
			return
		}
		fmt.Println(randomLen)
		imgUrl, serviceRondomErr := Service.RandomImg(DB, randomLen)
		if serviceRondomErr != nil {
			fmt.Println("查询数据库失败，错误原因\n", serviceRondomErr)
			context.JSON(200, gin.H{
				"code": "500",
				"err":  serviceRondomErr.Error(),
			})
			return
		}
		context.JSON(200, gin.H{
			"code": 200,
			"url":  imgUrl,
		})
	})
	server.GET("/random", func(context *gin.Context) {
		token := context.Query("token")
		if token == "" {
			context.JSON(200, gin.H{
				"code": "402",
				"err":  "认证失败,如需使用，请前往https://shop.loli.fit购买服务",
			})
			return
		}
		apiAuthType, randomNum := Service.ApiAuth(DB, token, "random_pro")
		if apiAuthType != true || randomNum < 1 {
			context.JSON(200, gin.H{
				"code": "403",
				"err":  "认证失败,如需使用，请前往https://shop.loli.fit购买服务",
			})
			return
		}
		fmt.Println(randomLen)
		imgUrl, serviceRondomErr := Service.RandomImg(DB, randomLen)
		if serviceRondomErr != nil {
			fmt.Println("查询数据库失败，错误原因\n", serviceRondomErr)
			context.JSON(200, gin.H{
				"code": "500",
				"err":  serviceRondomErr.Error(),
			})
			return
		}
		context.Redirect(302, imgUrl+"/conversion.webp")
	})
	go refreshLen()
	server.Run(":18848")
}
