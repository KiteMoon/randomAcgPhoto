package Service

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fastrand"
)

//
type randomImgReturnSql struct {
	Title  string `json:"title" db:"TITLE"`
	Path   string `json:"path" db:"PATH"`
	CdnUrl string `json:"CdnUrl" db:"CDNURL"`
}

// RandomImg 随机返回图片
func RandomImg(DB *sql.DB, tableLen int) (string, error) {
	//生成随机数
	// 查询行数
	//查询数据库中对应ID的信息
	sqlReutrnData := new(randomImgReturnSql)
	sqlQueryErr := DB.QueryRow("SELECT TITLE,URL,CDN_URL FROM RANDOM_PHOTO_TABLE WHERE PID=?", fastrand.Uint32n(uint32(tableLen))).Scan(&sqlReutrnData.Title, &sqlReutrnData.Path, &sqlReutrnData.CdnUrl)
	if sqlQueryErr != nil {
		fmt.Println("查询失败")
		fmt.Println(sqlQueryErr)
		return "查询失败\n错误码:RandomImg_001", sqlQueryErr
	}
	return sqlReutrnData.CdnUrl + sqlReutrnData.Path, nil
	//返回
}

// QueryRandomImg 精确筛选
func QueryRandomImg(context *gin.Context, DB *sql.DB, tableLen int) (string, error) {
	workName := context.Query("imgBClass")    //ACG类下游戏或者作品
	officalType := context.Query("imgCClass") // 官方/三方
	// 如果用户没有匹配任何作品
	if workName == "" {
		return RandomImg(DB, tableLen)
	}
	// 查询为官方作品
	if officalType == "offical" {
		sqlReutrnData := new(randomImgReturnSql)
		sqlQueryErr := DB.QueryRow("SELECT TITLE,URL,CDN_URL FROM RANDOM_PHOTO_TABLE WHERE PID=? AND OFFICAL=1 AND WORK=?", fastrand.Uint32n(uint32(tableLen)), workName).Scan(&sqlReutrnData.Title, &sqlReutrnData.Path, &sqlReutrnData.CdnUrl)
		if sqlQueryErr != nil {
			fmt.Println("查询失败")
			fmt.Println(sqlQueryErr)
			return "查询失败\n错误码:RandomImg_001", sqlQueryErr
		}
		return sqlReutrnData.CdnUrl + sqlReutrnData.Path, nil
	} else if officalType == "parties" { //识别为三方
		sqlReutrnData := new(randomImgReturnSql)
		sqlQueryErr := DB.QueryRow("SELECT TITLE,URL,CDN_URL FROM RANDOM_PHOTO_TABLE WHERE PID=? AND OFFICAL=0 AND WORK=?", fastrand.Uint32n(uint32(tableLen)), workName).Scan(&sqlReutrnData.Title, &sqlReutrnData.Path, &sqlReutrnData.CdnUrl)
		if sqlQueryErr != nil {
			fmt.Println("查询失败")
			fmt.Println(sqlQueryErr)
			return "查询失败\n错误码:RandomImg_001", sqlQueryErr
		}
		return sqlReutrnData.CdnUrl + sqlReutrnData.Path, nil
	} else {
		sqlReutrnData := new(randomImgReturnSql)
		sqlQueryErr := DB.QueryRow("SELECT TITLE,URL,CDN_URL FROM RANDOM_PHOTO_TABLE WHERE PID=? AND WORK=?", fastrand.Uint32n(uint32(tableLen)), workName).Scan(&sqlReutrnData.Title, &sqlReutrnData.Path, &sqlReutrnData.CdnUrl)
		if sqlQueryErr != nil {
			fmt.Println("查询失败")
			fmt.Println(sqlQueryErr)
			return "查询失败\n错误码:RandomImg_001", sqlQueryErr
		}
		return sqlReutrnData.CdnUrl + sqlReutrnData.Path, nil
	}
}
