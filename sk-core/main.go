package main

import (
	conf "final-design/pkg/config"
	"final-design/pkg/mysql"
	"final-design/sk-core/setup"
)

func main() {
	mysql.InitMysql(conf.MysqlConfig.Host, conf.MysqlConfig.Port, conf.MysqlConfig.User,
		conf.MysqlConfig.Pwd, conf.MysqlConfig.Db)
	setup.InitZk()
	setup.InitRedis()
	setup.RunService()
}
