package main

import (
	"encoding/json"
	"final-design/sk-admin/model"
	"fmt"
	"log"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	
}

func InitZk() {
	var (
		hosts = []string{"127.0.0.1:2181"}
		path  = "/product"
	)
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return
	}

	v, _, err := conn.Get(path)
	if err != nil {
		log.Printf("get [%s] from zk failed, err: %v", path, err)
		return
	}

	var secProductInfo []*model.SecProductInfoConf
	err = json.Unmarshal(v, &secProductInfo)
	if err != nil {
		log.Printf("Unmsharl secProductInfo failed, err: %v", err)
		return
	}
	for _, v := range secProductInfo {
		fmt.Println(v)
	}

	fmt.Println("zookeeper 连接成功")
}
