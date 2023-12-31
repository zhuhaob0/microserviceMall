package mysql

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
)

var engin *gorose.Engin
var err error

func InitMysql(hostMysql, portMysql, userMysql, pwdMysql, dbMysql string) {
	fmt.Println("InitMysql() => userMysql=", userMysql)
	fmt.Println("InitMysql() => dbMysql=", dbMysql)

	DbConfig := gorose.Config{
		// Default database configuration
		Driver: "mysql", // Database driver(mysql, sqlite, postgres, oracle, mssql)
		Dsn: userMysql + ":" + pwdMysql + "@tcp(" + hostMysql + ":" +
			portMysql + ")/" + dbMysql + "?charset=utf8&parseTime=true", // 数据库链接
		Prefix: "", // Table prefix
		// (Connection pool) Max open connections, default value 0 means unlimit.
		SetMaxOpenConns: 300,
		// (Connection pool) Max idle connections, default value is 1.
		SetMaxIdleConns: 10,
	}

	log.Println("DbConfig:", DbConfig)
	engin, err = gorose.Open(&DbConfig)
	if err != nil {
		fmt.Println("gorose.Open() happen an error:", err)
		return
	}
	fmt.Println("MySQL连接成功")
}

func DB() gorose.IOrm {
	return engin.NewOrm()
}
