package setup

import (
	conf "final-design/pkg/config"
	"final-design/sk-app/service/srv_redis"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/unknwon/com"
)

// 初始化Redis
func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Connect redis failed. Error: %v", err)
	}
	log.Printf("init redis success")
	conf.Redis.RedisConn = client

	loadBlackList(client)
	initRedisProcess()
	// 每隔30s从zookeeper拉取数据更新conf.SecKill.SecProductInfoMap
	go UpdateSecProductInfoMap()
}

// 加载黑名单列表
func loadBlackList(conn *redis.Client) {
	conf.SecKill.IDBlackMap = make(map[int]bool, 10000)
	conf.SecKill.IPBlackMap = make(map[string]bool, 10000)

	// 用户ID
	idList, err := conn.HGetAll(conf.Redis.IdBlackListHash).Result()
	if err != nil {
		log.Printf("HGetAll failed. Error: %v", err)
		return
	}
	for _, v := range idList {
		id, err := com.StrTo(v).Int()
		if err != nil {
			log.Printf("invalid user id [%v]", id)
			continue
		}
		conf.SecKill.IDBlackMap[id] = true
	}

	// 用户IP
	ipList, err := conn.HGetAll(conf.Redis.IpBlackListHash).Result()
	if err != nil {
		log.Printf("HGetAll failed. Error: %v", err)
		return
	}
	for _, ip := range ipList {
		conf.SecKill.IPBlackMap[ip] = true
	}

	go syncIdBlackList(conn)
	go syncIpBlackList(conn)
}

// 将redis的ID黑名单队列的数据取出，同步用户ID黑名单
func syncIdBlackList(conn *redis.Client) {
	for {
		idArr, err := conn.BRPop(time.Minute, conf.Redis.IdBlackListQueue).Result()
		if err != nil {
			log.Printf("BRPop id failed, err: %v", err)
			continue
		}
		id, _ := com.StrTo(idArr[1]).Int()
		conf.SecKill.RWBlackLock.Lock()
		{
			conf.SecKill.IDBlackMap[id] = true
		}
		conf.SecKill.RWBlackLock.Unlock()
	}
}

// 将redis的IP黑名单队列的数据取出，同步用户IP黑名单
func syncIpBlackList(conn *redis.Client) {
	var ipList []string
	lastTime := time.Now().Unix()

	for {
		ipArr, err := conn.BRPop(time.Minute, conf.Redis.IpBlackListQueue).Result()
		if err != nil {
			log.Printf("BRPop ip failed, err: %v", err)
			continue
		}

		ip := ipArr[1]
		curTime := time.Now().Unix()
		ipList = append(ipList, ip)
		// ipList长度过长，或者长时间无法给ipList添加元素，那么就将ipList的数据更新到RWBlackList中
		if len(ipList) > 100 || curTime-lastTime > 5 {
			conf.SecKill.RWBlackLock.Lock()
			{
				for _, v := range ipList {
					conf.SecKill.IPBlackMap[v] = true
				}
			}
			conf.SecKill.RWBlackLock.Unlock()

			lastTime = curTime
			log.Printf("sync ip list from redis success, ip[%v]", ipList)
			ipList = ipList[:0] // 清空ipList
		}
	}
}

// 初始化redis进程
func initRedisProcess() {
	log.Printf("initRedisProcess %d %d", conf.SecKill.AppWriteToHandleGoroutineNum, conf.SecKill.AppReadFromHandleGoroutineNum)

	for i := 0; i < conf.SecKill.AppWriteToHandleGoroutineNum; i++ { // 默认开10个goroutine
		go srv_redis.WriteHandle()
	}

	for i := 0; i < conf.SecKill.AppReadFromHandleGoroutineNum; i++ { // 默认开10个goroutine
		go srv_redis.ReadHandle()
	}

	for i := 0; i < 5; i++ { // 默认开5个goroutine
		go srv_redis.WriteOrder2DB()
	}
}

func UpdateSecProductInfoMap() {
	t1 := time.NewTicker(30 * time.Second)
	for {
		<-t1.C
		LoadSecConf(conf.Zk.ZkConn)
	}
}
