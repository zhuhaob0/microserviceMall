package setup

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	conf "final-design/pkg/config"

	"github.com/samuel/go-zookeeper/zk"
)

// 初始化Zk
func InitZk() {
	var hosts = []string{"127.0.0.1:2181"}
	option := zk.WithEventCallback(waitSecProductEvent)
	conn, _, err := zk.Connect(hosts, time.Second*5, option)
	if err != nil {
		fmt.Println(err)
		return
	}

	conf.Zk.ZkConn = conn
	conf.Zk.SecProductKey = "/product"
	{
		exists, _, _ := conn.Exists(conf.Zk.SecProductKey)
		if !exists {
			var byteData = []byte{}
			var flags int32 = 0
			// permission
			var acls = zk.WorldACL(zk.PermAll)

			_, err_create := conn.Create(conf.Zk.SecProductKey, byteData, flags, acls)
			if err_create != nil {
				fmt.Println(err_create)
			}
		}
		// 取出数据库中的活动数据

		// 同步到zookeeper

	}
	LoadSecConf(conn)
	// go watchZKEvent(conn, conf.Zk.SecProductKey) // 监听zookeeper变化，然后触发全局回调函数
}

// 加载秒杀商品信息
func LoadSecConf(conn *zk.Conn) {
	v, _, err := conn.Get(conf.Zk.SecProductKey)
	if err != nil {
		log.Printf("get product info failed, err: %v", err)
		return
	}
	// log.Printf("Connect zk success %s\n", conf.Zk.SecProductKey)
	log.Println("Get product info from zookeeper......")

	var secProductInfo []*conf.SecProductInfoConf
	err = json.Unmarshal(v, &secProductInfo)
	if err != nil {
		log.Printf("Unmarshal second product info failed, err: %v", err)
	}
	// log.Println("loadSecConf: secProductInfo=", secProductInfo)
	updateSecProductInfo(secProductInfo)
}

func waitSecProductEvent(event zk.Event) {
	log.Printf(">>>>>>>>>>>>>>>>>")
	log.Println("path:", event.Path)
	log.Println("type:", event.Type.String())
	log.Println("state:", event.State.String())
	log.Println("<<<<<<<<<<<<<<<<")
	if event.Path == conf.Zk.SecProductKey {
		log.Println("zookeeper中 [/product] 的数据发生更改")
		// 将zookeeper中的数据同步到conf.SecKill.SecProductInfoMap
		LoadSecConf(conf.Zk.ZkConn)
	}
}

// 更新秒杀商品信息
func updateSecProductInfo(secProductInfo []*conf.SecProductInfoConf) {
	tmp := make(map[int]*conf.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		log.Printf("updateSecProductInfo %v", v)
		tmp[v.ProductId] = v
	}
	conf.SecKill.RWBlackLock.Lock()
	conf.SecKill.SecProductInfoMap = tmp
	conf.SecKill.RWBlackLock.Unlock()
}

// 监听zookeeper的path
func watchZKEvent(conn *zk.Conn, path string) {
	for {
		_, _, event, err := conn.ExistsW(path)
		<-event // zookeeper没有事件时，阻塞在这里
		fmt.Println("watchAKEvent: path=", path)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
