package srv_redis

import (
	"crypto/md5"
	conf "final-design/pkg/config"
	"final-design/sk-core/config"
	"final-design/sk-core/service/srv_err"
	"final-design/sk-core/service/srv_user"
	"fmt"
	"log"
	"sync"
	"time"
)

func HandleUser() {
	log.Println("handle user running")

	for req := range config.SecLayerCtx.Read2HandleChan {
		log.Printf("begin process request: %v\n", req)
		res, err := HandleSeckill(req)

		if err != nil {
			log.Printf("process request %v failed, err: %v", req, err)
			res = &config.SecResult{
				Code: srv_err.ErrServiceBusy,
			}
		}

		fmt.Println("处理中~~", res)
		timer := time.NewTicker(time.Millisecond * time.Duration(conf.SecKill.SendToWriteChanTimeout))
		select {
		case config.SecLayerCtx.Handle2WriteChan <- res:
		case <-timer.C:
			log.Printf("send to response chan timeout, res: %v", res)
			// break
		}
	}
}

// 核心逻辑，对SecRequest处理，返回SecResult
func HandleSeckill(req *config.SecRequest) (res *config.SecResult, err error) {
	config.SecLayerCtx.RWSecProductLock.Lock()
	defer config.SecLayerCtx.RWSecProductLock.Unlock()

	res = &config.SecResult{}
	res.ProductId = req.ProductId
	res.UserId = req.UserId

	product, ok := conf.SecKill.SecProductInfoMap[req.ProductId] //找不到商品
	log.Println("product=", product)
	if !ok {
		log.Printf("not found product: %v\n", req.ProductId)
		res.Code = srv_err.ErrNotFoundProduct
		return
	}

	if product.Status == srv_err.ProductStatusSoldOut { //商品卖完了
		res.Code = srv_err.ErrSoldOut
		return
	}

	nowTime := time.Now().Unix()
	config.SecLayerCtx.HistoryMapLock.Lock()
	userHistory, ok := config.SecLayerCtx.HistoryMap[req.UserId] // 获取用户购买记录
	if !ok {
		userHistory = &srv_user.UserBuyHistory{ // 用户没有购买过，就新建一个HistoryMap
			History: make(map[int]int, 16),
		}
		config.SecLayerCtx.HistoryMap[req.UserId] = userHistory
	}
	historyCount := userHistory.GetProductBuyCount(req.ProductId) // 用户已经买了多少个
	config.SecLayerCtx.HistoryMapLock.Unlock()

	fmt.Println("max_buy_per_person=", product.MaxBuyPerPerson)
	if historyCount >= product.MaxBuyPerPerson { // 超过了一个人最大购买数量
		res.Code = srv_err.ErrAlreadyBuy
		return
	}

	curSoldCount := config.SecLayerCtx.ProductCountMgr.Count(req.ProductId) //商品已经卖出的数量
	if curSoldCount >= product.Total {                                      // 商品卖完了
		fmt.Println("curSoldCount=", curSoldCount, " total=", product.Total)
		res.Code = srv_err.ErrSoldOut
		product.Status = srv_err.ProductStatusSoldOut
		return
	}

	// curRate:=rand.Float64()
	// curRate := 0.1 // 为啥是0.1？
	// fmt.Println("curRate=", curRate, "product.BuyRate=", product.BuyRate)
	// if curRate > product.BuyRate { // 当前购买速率超过了产品购买频率速率
	// 	res.Code = srv_err.ErrRetry
	// 	return
	// }

	userHistory.Add(req.ProductId, 1)                        // 用户购买数量 + 1
	config.SecLayerCtx.ProductCountMgr.Add(req.ProductId, 1) // 产品售卖数量 + 1

	// product.LeftNum = product.LeftNum - 1 // 商品剩余数量 - 1

	var rwMutex sync.RWMutex
	rwMutex.Lock() // 防止超卖
	{
		if product.LeftNum > 0 {
			product.LeftNum = product.LeftNum - 1 // 商品剩余数量 - 1
		} else {
			rwMutex.Unlock()
			res.Code = srv_err.ErrSoldOut
			product.Status = srv_err.ProductStatusSoldOut
			return
		}
	}
	rwMutex.Unlock()

	// 组装order: req.ProductId req.ProductName req.SecTime req.Username req.ActivityPrice
	order := config.Order{
		ProductId:     req.ProductId,
		ProductName:   req.ProductName,
		OrderTime:     req.SecTime,
		Buyer:         req.Username,
		ActivityPrice: req.ActivityPrice,
	}
	config.SecLayerCtx.WriteOrder2RedisChan <- &order

	// 用户ID，商品ID，当前时间，密钥
	res.Code = srv_err.ErrSecKillSucc
	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s",
		req.UserId, req.ProductId, nowTime, conf.SecKill.TokenPassWd)
	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData))) // MD5加密
	res.TokenTime = nowTime

	return
}
