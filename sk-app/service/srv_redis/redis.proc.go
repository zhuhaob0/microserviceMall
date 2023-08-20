package srv_redis

import (
	"encoding/json"
	conf "final-design/pkg/config"
	skadmin_model "final-design/sk-admin/model"
	"final-design/sk-app/config"
	"final-design/sk-app/model"
	"fmt"
	"log"
	"time"
)

// 写数据到redis
func WriteHandle() {
	for {
		fmt.Println("write data to redis.")
		req := <-config.SkAppContext.SecReqChan //SecKill(*model.SecRequest)放入的请求
		fmt.Println("accessTime =", req.AccessTime)
		conn := conf.Redis.RedisConn

		data, err := json.Marshal(req)
		if err != nil {
			log.Printf("json.Marshal req failed. Error: %v, req: %v", err, req)
			continue
		}
		fmt.Println("conf.Redis.Proxy2layerQueueName", conf.Redis.Proxy2layerQueueName)
		err = conn.LPush(conf.Redis.Proxy2layerQueueName, string(data)).Err() //放入redis队列中，让sk-core处理
		if err != nil {
			log.Printf("LPush req failed. Error: %v, req: %v", err, req)
			continue
		}
		log.Printf("Lpush req success. req: %v", string(data))
	}
}

// 从redis中读数据
func ReadHandle() {
	for {
		conn := conf.Redis.RedisConn
		// 阻塞弹出
		// fmt.Println("conf.Redis.Layer2proxyQueueName=", conf.Redis.Layer2proxyQueueName)
		data, err := conn.BRPop(time.Second, conf.Redis.Layer2proxyQueueName).Result() // 取出sk-core处理的结果
		if err != nil {
			// log.Printf("BRPop layer2proxy failed. Error: %v", err)
			continue
		}

		var result *model.SecResult
		err = json.Unmarshal([]byte(data[1]), &result) //将结果反序列化为model.SecResult
		if err != nil {
			log.Printf("json.Unmarshal failed. Error: %v\n", err)
			continue
		}

		userKey := fmt.Sprintf("%d_%d", result.UserId, result.ProductId)
		fmt.Println("userKey: ", userKey)

		config.SkAppContext.UserConnMapLock.Lock()
		resChan, ok := config.SkAppContext.UserConnMap[userKey] // 通过userKey找到要发送结果的channel
		config.SkAppContext.UserConnMapLock.Unlock()

		if !ok {
			log.Printf("user not found: %v", userKey)
			continue
		}
		log.Printf("request result send to chan")

		resChan <- result //结果发送回service.go的SecKill函数中
		log.Printf("request result send to chan success, userKey: %v", userKey)
	}
}

// 定时从redis中读取订单数据，写入到数据库中
func WriteOrder2DB() {
	t := time.NewTicker(time.Second * 30)
	conn := conf.Redis.RedisConn
	for {
		<-t.C
		for {
			data, err := conn.BRPop(time.Second, conf.Redis.Layer2DBQueueName).Result() // 取出sk-core返回的order
			if err != nil {                                                             // redis为空停止此次同步操作
				break
			}
			var order *skadmin_model.Order
			err = json.Unmarshal([]byte(data[1]), &order) //将结果反序列化为model.SecResult
			if err != nil {
				log.Printf("json.Unmarshal failed. Error: %v\n", err)
				continue
			}
			log.Println("order=", order)
			// 将order写入DB中
			orderModel := skadmin_model.NewOrderModel()
			err = orderModel.CreateOrder(order)
			if err != nil {
				log.Printf("order写入数据库时失败, Error = %v", err)
			}

		}
	}
}
