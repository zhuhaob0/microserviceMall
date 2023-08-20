package main

import (
	"final-design/pkg/bootstrap"
	conf "final-design/pkg/config"
	"final-design/pkg/mysql"
	"final-design/sk-app/setup"
)

func main() {
	mysql.InitMysql(conf.MysqlConfig.Host, conf.MysqlConfig.Port, conf.MysqlConfig.User,
		conf.MysqlConfig.Pwd, conf.MysqlConfig.Db) // conf.MysqlConfig.Db

	setup.InitZk()
	setup.InitRedis()

	setup.InitServer(bootstrap.HttpConfig.Host, bootstrap.HttpConfig.Port)
	// servicePort := flag.String("servicePort", "9031", "service port")
	// flag.Parse()
	// setup.InitServer(bootstrap.HttpConfig.Host, *servicePort)
}
