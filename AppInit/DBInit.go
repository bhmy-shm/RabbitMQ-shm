package AppInit

import (
	"fmt"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func DBInit(database string) error {
	var err error
	//数据库连接信息
	dsn := fmt.Sprintf("root:123.com@tcp(127.0.0.1:3306)/%s?" +
		"charset=utf8mb4&parseTime=True&loc=Local",database)

	//连接数据库
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("dsn:%s failed err :%v\n", dsn, err)
		return err
	}

	//设置数据库的最大连接数,设置为10
	db.SetMaxOpenConns(10)
	//最大空闲连接数
	db.SetMaxIdleConns(5)

	return nil
}

func GetDB() *sqlx.DB {
	return db
}
