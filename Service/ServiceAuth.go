package Service

import (
	"database/sql"
	"fmt"
)

type queryAuthSqlStruct struct {
	APPID string
	UID   string
}

// ApiAuth 对外提公开接口，自动从数据库中查询内容并返回,token为秘钥，authType为请求权限
func ApiAuth(db *sql.DB, token, authType string) (bool, int) {
	// 查询对应的Token和对应的权限是否存在
	querySqlStr := "SELECT APPID,UID FROM TOKEN_TABLE WHERE TOKEN=? AND TYPE=?"
	// 创建一个接收器
	returnData := new(queryAuthSqlStruct)
	queryTokenSqlErr := db.QueryRow(querySqlStr, token, authType).Scan(&returnData.APPID, &returnData.UID)
	if queryTokenSqlErr != nil {
		fmt.Println("查询该token不存在")
		fmt.Println(queryTokenSqlErr)
		return false, 0
	} else if returnData.APPID == "" {
		fmt.Println("查询该token不存在")
		return false, 0
	}
	// 查询剩余调用次数
	if authType == "random_pro" {
		queryTypeNumSqlStr := "SELECT RANDOM_PRO_NUM FROM USER_TABLE WHERE UID=? AND TYPE=1"
		apiResuidNum := new(int)
		queryTypeNumSqlErr := db.QueryRow(queryTypeNumSqlStr, returnData.UID).Scan(&apiResuidNum)
		if queryTypeNumSqlErr != nil {
			fmt.Println("查询该用户不存在")
			fmt.Println(queryTokenSqlErr)
			return false, 0
		} else if apiResuidNum == nil {
			fmt.Println("查询该用户不存在")
			fmt.Println(queryTokenSqlErr)
			return false, 0
		} else if *apiResuidNum <= 0 {
			fmt.Println("查询该用户不存在")
			fmt.Println(queryTokenSqlErr)
			return false, 0
		}
		expendRandomSqlStr := "update USER_TABLE SET RANDOM_PRO_NUM=? where UID=?"
		_, expendRandomSqlErr := db.Exec(expendRandomSqlStr, *apiResuidNum-1, returnData.UID)
		if expendRandomSqlErr != nil {
			fmt.Println("消耗次数失败")
			fmt.Println("后台错误")
			fmt.Println(expendRandomSqlErr)
			return false, 0
		}
		return true, *apiResuidNum
	}
	fmt.Println("没有匹配的API")
	return false, 0
	// 存在，返回处理的信息
}
