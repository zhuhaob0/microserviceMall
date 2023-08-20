package service

import (
	conf "final-design/pkg/config"
	"final-design/sk-app/config"
	"final-design/sk-app/model"
	"final-design/sk-app/service/srv_err"
	"final-design/sk-app/service/srv_limit"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Service define a service interface
type Service interface {
	// HealthCheck check service health status
	HealthCheck() bool
	SecInfo(productId int) map[string]interface{}
	SecKill(req *model.SecRequest) (map[string]interface{}, int, error)
	SecInfoList() ([]map[string]interface{}, int, error)
}

type ServiceMiddleware func(Service) Service

// UserService implement Service interface
type SkAppService struct {
}

// HealthCheck implement Service interface
// 用于检查服务的健康状态，这里仅仅返回true
func (s SkAppService) HealthCheck() bool {
	return true
}

// 根据productId在conf.SecKill.SecProductInfoMap中获取商品数据
func (s SkAppService) SecInfo(productId int) map[string]interface{} {
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()
	v, ok := conf.SecKill.SecProductInfoMap[productId]
	if !ok {
		return nil
	}
	data := make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTime
	data["status"] = v.Status

	return data
}

func (s SkAppService) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	//对Map加锁处理
	//config.SkAppContext.RWSecProductLock.RLock()
	//defer config.SkAppContext.RWSecProductLock.RUnlock()

	var code int
	err := srv_limit.AntiSpam(req) // 防作弊，id、ip黑名单校验，秒级限制、分级限制
	if err != nil {
		code = srv_err.ErrUserServiceBusy
		log.Printf("userId [%d] antiSpam is failed, err: [%v]", req.UserId, err)
		return nil, code, err
	}

	data, code, err := SecInfoById(req.ProductId) // 判断商品是否因各种原因不再销售
	if err != nil {
		log.Printf("userId [%d] secInfoById is failed, err: [%v]", req.UserId, err)
		return nil, code, err
	}

	userKey := fmt.Sprintf("%d_%d", req.UserId, req.ProductId)
	ResultChan := make(chan *model.SecResult, 1)
	config.SkAppContext.UserConnMapLock.Lock()
	// 通过userKey映射一个channel，结果返回时也要通过这个userKey找到这个channel，
	// 然后传过来，源码在redis_proc.go的ReadHandle函数
	config.SkAppContext.UserConnMap[userKey] = ResultChan
	config.SkAppContext.UserConnMapLock.Unlock()

	// 将请求送入通道并推入到redis队列当中
	config.SkAppContext.SecReqChan <- req
	// 启动定时器等待结果
	ticker := time.NewTicker(time.Millisecond * time.Duration(conf.SecKill.AppWaitResultTimeout))

	defer func() {
		ticker.Stop()
		config.SkAppContext.UserConnMapLock.Lock()
		delete(config.SkAppContext.UserConnMap, userKey) // 抢购结束了，无论成功与否，需将userKey删除
		config.SkAppContext.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = srv_err.ErrProcessTimeout
		err = fmt.Errorf("SecKill request timeout")
		return nil, code, err

	case <-req.CloseNotify:
		code = srv_err.ErrClientClosed
		err = fmt.Errorf("client already closed")
		return nil, code, err

	case result := <-ResultChan:
		code = result.Code
		if code != srv_err.ErrSecKillSucc { // 因其他原因抢购失败
			return data, code, srv_err.GetErrMsg(code)
		}
		log.Printf("secKill success\n")
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId
		return data, code, nil
	}
}

// 将conf.SecKill.SecProductInfoMap中的所有数据，通过SecInfoById(int)函数筛选一遍，
// 将秒杀结束、停售、售罄以及超过购买频率限制的商品过滤掉，将剩余商品返回
func (s SkAppService) SecInfoList() ([]map[string]interface{}, int, error) {
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	var data []map[string]interface{}
	for _, v := range conf.SecKill.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			log.Printf("get sec info, err: %v", err)
			continue
		}
		data = append(data, item)
	}
	return data, len(data), nil // 活动数据，数目，nil
}

func SecInfoById(productId int) (map[string]interface{}, int, error) {
	// 对Map加锁处理
	// config.SkAppContext.RWSecProductLock.RLock()
	// defer config.SkAppContext.RWSecProductLock.RUnlock()

	var code int
	v, ok := conf.SecKill.SecProductInfoMap[productId]
	if !ok {
		return nil, srv_err.ErrNotFoundProductId, fmt.Errorf("not found product_id: %d", productId)
	}
	fmt.Printf("v= %+v\n", v)
	start, end, status := false, false, "success" // 活动是否开始、结束；状态
	var err error
	nowTime := time.Now().Unix()

	// 秒杀活动没有开始
	if nowTime < v.StartTime {
		start = false
		end = false
		status = "second kill not start"
		code = srv_err.ErrActiveNotStart
		err = fmt.Errorf(status)
	}
	// 秒杀活动已经开始
	if nowTime > v.StartTime {
		start = true
	}
	// 秒杀活动已经结束
	if nowTime > v.EndTime {
		start = false
		end = true
		status = "second kill is already end"
		code = srv_err.ErrActiveAlreadyEnd
		err = fmt.Errorf(status)
		// fmt.Println("结束: err=", err)
		return nil, code, err
	}
	// 活动商品售罄
	if v.Status == config.ProductStatusSoldOut {
		start = false
		end = false
		status = "product is sold out"
		code = srv_err.ErrActiveSaleOut
		err = fmt.Errorf(status)
		return nil, code, err
	}

	curRate := rand.Float64() //为啥用随机？
	// fmt.Println("curRate=", curRate)
	// fmt.Println("v.BuyRate=", v.BuyRate)
	// 放小于等于购买比率的1.5倍的请求进入core层
	v.BuyRate = 1 // 临时添加
	if curRate > v.BuyRate*1.5 {
		start = false
		end = false
		status = "retry"
		code = srv_err.ErrRetry
		err = fmt.Errorf(status)
		return nil, code, err
	}

	// 组装数据
	data := map[string]interface{}{
		"product_id": productId,
		"start":      start,
		"end":        end,
		"status":     status,
	}
	return data, code, err
}

func NewSecRequest() *model.SecRequest {
	secRequest := &model.SecRequest{
		ResultChan: make(chan *model.SecResult, 1),
	}
	return secRequest
}
