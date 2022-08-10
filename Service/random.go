package Service

import (
	"database/sql"
	"fmt"
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
