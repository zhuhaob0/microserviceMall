package srv_limit

import (
	conf "final-design/pkg/config"
	"final-design/sk-app/model"
	"fmt"
	"log"
	"sync"
)

// 限制管理
type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	lock         sync.Mutex
}

var SecLimitMgrVars = &SecLimitMgr{
	UserLimitMap: make(map[int]*Limit),
	IpLimitMap:   make(map[string]*Limit),
}

// 使用ip、id黑名单，以及对ip、id访问频率控制，防止作弊
func AntiSpam(req *model.SecRequest) (err error) {
	// 判断用户ID是否在黑名单
	_, ok := conf.SecKill.IDBlackMap[req.UserId]
	if ok {
		err = fmt.Errorf("invalid request")
		log.Printf("user[%v] is blocked by id black", req.UserId)
		return
	}

	// 判断用户IP是否在黑名单
	_, ok = conf.SecKill.IPBlackMap[req.ClientAddr]
	if ok {
		err = fmt.Errorf("invalid request")
		log.Printf("userId[%v] ip[%v] is blocked by ip black", req.UserId, req.ClientAddr)
		return
	}
	var secIdCount, minIdCount, secIpCount, minIpCount int
	// 加锁
	SecLimitMgrVars.lock.Lock()
	{
		// 用户ID频率控制
		limit, ok := SecLimitMgrVars.UserLimitMap[req.UserId]
		if !ok {
			limit = &Limit{
				secLimit: &SecLimit{},
				minLimit: &MinLimit{},
			}
			SecLimitMgrVars.UserLimitMap[req.UserId] = limit
		}
		secIdCount = limit.secLimit.Count(req.AccessTime) // 获取该秒内该用户访问次数
		minIdCount = limit.minLimit.Count(req.AccessTime) // 获取该分钟内该用户访问次数

		// 客户端IP频率控制
		limit, ok = SecLimitMgrVars.IpLimitMap[req.ClientAddr]
		if !ok {
			limit = &Limit{
				secLimit: &SecLimit{},
				minLimit: &MinLimit{},
			}
			SecLimitMgrVars.IpLimitMap[req.ClientAddr] = limit
		}
		secIpCount = limit.secLimit.Count(req.AccessTime) // 获取该秒内该IP访问次数
		minIpCount = limit.minLimit.Count(req.AccessTime) // 获取该秒内该IP访问次数
	}
	SecLimitMgrVars.lock.Unlock() // 释放锁

	// fmt.Println("userSecAccessLimit=", conf.SecKill.AccessLimitConf.UserSecAccessLimit)
	// fmt.Println("userMinAccessLimit=", conf.SecKill.AccessLimitConf.UserMinAccessLimit)
	// fmt.Println("ipSecAccessLimit=", conf.SecKill.AccessLimitConf.IPSecAccessLimit)
	// fmt.Println("ipMinAccessLimit=", conf.SecKill.AccessLimitConf.IPMinAccessLimit)
	// fmt.Println(conf.SecKill.AppWriteToHandleGoroutineNum)

	// 判断该用户一秒内访问次数是否大于配置的最大访问次数
	if secIdCount > conf.SecKill.AccessLimitConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request")
		conf.SecKill.IDBlackMap[req.UserId] = true
		return
	}

	// 判断该用户一分钟内访问次数是否大于配置的最大访问次数
	if minIdCount > conf.SecKill.AccessLimitConf.UserMinAccessLimit {
		err = fmt.Errorf("invalid request")
		conf.SecKill.IDBlackMap[req.UserId] = true
		return
	}

	// 判断该IP一秒内访问次数是否大于配置的最大访问次数
	if secIpCount > conf.SecKill.AccessLimitConf.IPSecAccessLimit {
		err = fmt.Errorf("invalid request")
		conf.SecKill.IPBlackMap[req.ClientAddr] = true
		return
	}

	// 判断该IP一分钟内访问次数是否大于配置的最大访问次数
	if minIpCount > conf.SecKill.AccessLimitConf.IPMinAccessLimit {
		err = fmt.Errorf("invalid request")
		conf.SecKill.IPBlackMap[req.ClientAddr] = true
		return
	}

	return
}
