package config

import (
	"final-design/sk-core/service/srv_limit"
	"sync"

	"github.com/go-redis/redis"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	Redis       RedisConf
	SecKill     SecKillConf
	MysqlConfig MysqlConf
	TraceConfig TraceConf
	Zk          ZookeeperConf
)

type ZookeeperConf struct {
	ZkConn        *zk.Conn
	SecProductKey string // 商品键
}

type TraceConf struct {
	Host string
	Port string
	Url  string
}

type MysqlConf struct {
	Host string
	Port string
	User string
	Pwd  string
	Db   string
}

// redis配置
type RedisConf struct {
	RedisConn            *redis.Client // 链接
	Proxy2layerQueueName string        // 队列名称
	Layer2proxyQueueName string        // 队列名称
	Layer2DBQueueName    string        // 秒杀订单写入redis队列
	IdBlackListHash      string        // 用户黑名单hash表
	IpBlackListHash      string        // IP黑名单hash表
	IdBlackListQueue     string        // 用户黑名单队伍
	IpBlackListQueue     string        // IP黑名单队伍
	Host                 string
	Password             string
	Db                   int
}

type SecKillConf struct {
	RedisConf *RedisConf

	CookieSecretKey string
	ReferWhiteList  []string // 白名单
	AccessLimitConf AccessLimitConf

	RWBlackLock                  sync.RWMutex
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int

	IPBlackMap map[string]bool
	IDBlackMap map[int]bool

	SecProductInfoMap map[int]*SecProductInfoConf

	AppWriteToHandleGoroutineNum  int
	AppReadFromHandleGoroutineNum int
	CoreReadRedisGoroutineNum     int
	CoreWriteRedisGoroutineNum    int
	CoreHandleGoroutineNum        int

	AppWaitResultTimeout    int
	CoreWaitResultTimeout   int
	MaxRequestWaitTimeout   int
	SendToWriteChanTimeout  int
	SendToHandleChanTimeout int

	TokenPassWd string
}

// 商品信息配置
type SecProductInfoConf struct {
	ActivityName     string  `json:"activity_name"`       // 活动名
	ProductId        int     `json:"product_id"`          // 商品Id
	ProductName      string  `json:"product_name"`        // 商品名
	StartTime        int64   `json:"start_time"`          // 开始时间
	EndTime          int64   `json:"end_time"`            // 结束时间
	ActivityPrice    int     `json:"activity_price"`      // 活动价
	Status           int     `json:"status"`              // 状态
	Total            int     `json:"total"`               // 商品总数
	LeftNum          int     `json:"left_num"`            // 剩余商品数
	MaxBuyPerPerson  int     `json:"max_buy_per_person"`  // 单人购买限制
	MaxSoldPerSecond int     `json:"max_sold_per_second"` // 每秒最多能卖多少
	BuyRate          float64 `json:"buy_rate"`            // 买中几率
	// todo: error
	SecLimit *srv_limit.SecLimit `json:"sec_limit"` // 限速控制
}

// 访问限制
type AccessLimitConf struct {
	IPSecAccessLimit   int // IP每秒访问限制
	UserSecAccessLimit int // 用户每秒访问限制
	IPMinAccessLimit   int // IP每分钟访问限制
	UserMinAccessLimit int // 用户每分钟访问限制
}
