package main

import (
	"final-design/pkg/bootstrap"
	conf "final-design/pkg/config"
	"final-design/pkg/mysql"
	"final-design/sk-admin/setup"
)

func main() {
	mysql.InitMysql(conf.MysqlConfig.Host, conf.MysqlConfig.Port, conf.MysqlConfig.User,
		conf.MysqlConfig.Pwd, conf.MysqlConfig.Db)
	setup.InitZk()
	setup.InitSever(bootstrap.HttpConfig.Host, bootstrap.HttpConfig.Port, bootstrap.DiscoverConfig.Host, bootstrap.DiscoverConfig.Port)
}
