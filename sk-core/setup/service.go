package setup

import (
	"encoding/json"
	conf "final-design/pkg/config"
	"final-design/pkg/mysql"
	"final-design/sk-admin/service"
	"final-design/sk-core/service/srv_redis"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunService() {
	// 启动处理线程
	srv_redis.RunProcess()
	go store2Database()

	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	fmt.Println(error)
}

// 每隔30s同步内存Activity数据到 数据库 和 zookeeper
func store2Database() {
	conn := mysql.DB()
	t1 := time.NewTicker(30 * time.Second)
	for {
		<-t1.C
		conf.Logger.Log("Activity内存数据持久化到数据库中......")
		// 将conf.SecKill.SecProductInfoMap持久化到mysql
		leftNum := make(map[string]int, 20)
		for _, v := range conf.SecKill.SecProductInfoMap {
			fmt.Println("activity_name=", v.ActivityName, " left_num=", v.LeftNum)
			leftNum[v.ActivityName] = v.LeftNum
			_, err := conn.Execute("update activity set left_num=? where activity_name=?", v.LeftNum, v.ActivityName)
			// _, err := conn.Table("activity").Data(map[string]interface{}{
			// 	"left_num": v.LeftNum,
			// }).Where("activity_name", v.ActivityName).Update()

			if err != nil {
				log.Printf("UpdateActivity【%s】, Error: %v\n", v.ActivityName, err)
			}
		}
		conf.Logger.Log("Activity内存数据同步到Zookeeper中......")
		activityImpl := service.ActivityServiceImpl{}
		secProductInfoList, stat, _ := activityImpl.LoadProductFromZk(conf.Zk.SecProductKey)
		for i, item := range secProductInfoList {
			secProductInfoList[i].LeftNum = leftNum[item.ActivityName]
		}
		// 将conf.SecKill.SecProductInfoMap数据同步到zookeeper
		data, err := json.Marshal(secProductInfoList)
		if err != nil {
			log.Printf("json marshal failed, err: %v", err)
			return
		}

		conn := conf.Zk.ZkConn
		var byteData = []byte(string(data))

		_, err_set := conn.Set(conf.Zk.SecProductKey, byteData, stat.Version)
		if err_set != nil {
			fmt.Println("err_set=", err_set)
		}
	}
}
