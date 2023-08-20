package setup

import (
	"encoding/json"
	"fmt"
	"time"

	conf "final-design/pkg/config"
	"final-design/sk-admin/model"
	"final-design/sk-admin/service"

	"github.com/gohouse/gorose/utils"
	"github.com/samuel/go-zookeeper/zk"
)

// 初始化Zk
func InitZk() {
	var hosts = []string{"127.0.0.1"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
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
		// 取出数据库中的Activity数据
		var activityimpl service.ActivityServiceImpl
		data, _ := activityimpl.GetActivityList()

		// 同步到zookeeper
		var activities []*model.SecProductInfoConf
		for _, v := range data {
			jsonStr, _ := utils.JsonEncode(v)
			var activity model.Activity
			json.Unmarshal([]byte(jsonStr), &activity)

			tmp := &model.SecProductInfoConf{
				ActivityName: activity.ActivityName,
				ProductId:    activity.ProductId,
				ProductName:  activity.ProductName,

				StartTime: activity.StartTime,
				EndTime:   activity.EndTime,
				Total:     activity.Total,
				LeftNum:   activity.LeftNum,
				Status:    activity.Status,

				MaxBuyPerPerson:  activity.MaxBuyPerPerson,
				MaxSoldPerSecond: activity.MaxSoldPerSecond,
				ActivityPrice:    activity.ActivityPrice,
				BuyRate:          activity.BuyRate,
			}
			activities = append(activities, tmp)
		}
		fmt.Println("a=", activities)
		// 同步到Zookeeper
		_, stat, _ := activityimpl.LoadProductFromZk(conf.Zk.SecProductKey)
		activityimpl.SyncToZK(activities, conf.Zk.SecProductKey, stat)
	}

	fmt.Println("zookeeper 连接成功")
}
