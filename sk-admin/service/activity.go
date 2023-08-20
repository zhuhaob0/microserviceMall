package service

import (
	"context"
	"encoding/json"
	"final-design/sk-admin/model"
	"fmt"
	"log"
	"time"

	conf "final-design/pkg/config"

	"github.com/gohouse/gorose/v2"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/unknwon/com"
)

type ActivityService interface {
	GetActivityList() ([]gorose.Data, error)
	CreateActivity(activity *model.Activity) error
	UpdateActivity(activity *model.Activity) error
	DeleteActivity(activity *model.Activity) error
}

type ActivityServiceMiddleware func(ActivityService) ActivityService

type ActivityServiceImpl struct{}

// 从数据库中读取活动数据，将非正常的活动过滤掉
func (p ActivityServiceImpl) GetActivityList() ([]gorose.Data, error) {
	activityEntity := model.NewActivityModel()
	activityList, err := activityEntity.GetActivityList()
	if err != nil {
		log.Printf("ActivityEntity.GetActivityList, err: %v", err)
		return nil, err
	}

	for _, v := range activityList {
		startTime, _ := com.StrTo(fmt.Sprint(v["start_time"])).Int64()
		v["start_time_str"] = time.Unix(startTime, 0).Format("2006-01-02 15:04:05")

		endTime, _ := com.StrTo(fmt.Sprint(v["end_time"])).Int64()
		v["end_time_str"] = time.Unix(endTime, 0).Format("2006-01-02 15:04:05")

		nowTime := time.Now().Unix()
		if nowTime > endTime {
			v["status_str"] = "结束"
			continue
		}

		status, _ := com.StrTo(fmt.Sprint(v["status"])).Int()
		if status == model.ActivityStatusNormal {
			v["status_str"] = "正常"
		} else if status == model.ActivityStatusDisable {
			v["status_str"] = "禁用"
		}

	}
	log.Printf("get activity success, activity list is [%v]", activityList)
	return activityList, nil
}

func (p ActivityServiceImpl) UpdateActivity(activity *model.Activity) error {
	activityEntity := model.NewActivityModel()
	err := activityEntity.UpdateActivity(activity)
	if err != nil {
		log.Printf("ActivityModel.UpdateActivity, err: %v", err)
		return err
	}
	log.Println("updateSyncToZK......")
	err = p.updateSyncToZK(activity) // 更新ZK的数据
	if err != nil {
		log.Printf("updateSyncToZK is failed, err: %v", err)
		return err
	}
	return nil
}

func (p ActivityServiceImpl) DeleteActivity(activity *model.Activity) error {
	activityEntity := model.NewActivityModel()
	err := activityEntity.DeleteActivity(activity)
	if err != nil {
		log.Printf("ActivityModel.DeleteActivity, err: %v", err)
		return err
	}
	log.Println("deleteSyncToZK......")
	err = p.deleteSyncToZK(activity) // 删除ZK的数据
	if err != nil {
		log.Printf("deleteSyncToZK is failed, err: %v", err)
		return err
	}
	return nil
}

// 创建活动到数据库，将秒杀活动信息同步到zookeeper
func (p ActivityServiceImpl) CreateActivity(activity *model.Activity) error {
	activityEntity := model.NewActivityModel()
	err := activityEntity.CreateActivity(activity)
	if err != nil {
		log.Printf("ActivityModel.CreateActivity, err: %v", err)
		return err
	}

	log.Println("createSyncToZK......")
	err = p.createSyncToZK(activity) // 写入到Zk
	if err != nil {
		log.Printf("createSyncToZK is failed, err: %v", err)
		return err
	}
	return nil
}

func (p ActivityServiceImpl) createSyncToZK(activity *model.Activity) error {
	zkPath := conf.Zk.SecProductKey
	secProductInfoList, stat, err := p.LoadProductFromZk(zkPath)
	if err != nil {
		secProductInfoList = []*model.SecProductInfoConf{}
	}

	var secProductInfo = &model.SecProductInfoConf{
		ActivityName: activity.ActivityName,
		ProductId:    activity.ProductId,
		ProductName:  activity.ProductName,

		StartTime: activity.StartTime,
		EndTime:   activity.EndTime,
		Status:    activity.Status,
		Total:     activity.Total,
		LeftNum:   activity.LeftNum,

		MaxBuyPerPerson:  activity.MaxBuyPerPerson,
		MaxSoldPerSecond: activity.MaxSoldPerSecond,
		ActivityPrice:    activity.ActivityPrice,
		BuyRate:          1,
	}
	secProductInfoList = append(secProductInfoList, secProductInfo)

	return p.SyncToZK(secProductInfoList, zkPath, stat)
}

func (p ActivityServiceImpl) updateSyncToZK(activity *model.Activity) error {
	zkPath := conf.Zk.SecProductKey
	secProductInfoList, stat, _ := p.LoadProductFromZk(zkPath)
	existed := false
	for i, item := range secProductInfoList {
		if activity.ActivityName == item.ActivityName {
			existed = true
			secProductInfoList[i] = &model.SecProductInfoConf{
				ActivityName: secProductInfoList[i].ActivityName,
				ProductId:    secProductInfoList[i].ProductId,
				ProductName:  secProductInfoList[i].ProductName,

				StartTime: activity.StartTime,
				EndTime:   activity.EndTime,
				Total:     activity.Total,
				LeftNum:   activity.LeftNum,
				Status:    activity.Status,

				MaxBuyPerPerson:  activity.MaxBuyPerPerson,
				MaxSoldPerSecond: activity.MaxSoldPerSecond,
				ActivityPrice:    activity.ActivityPrice,
				BuyRate:          1,
			}
			break
		}
	}
	if existed {
		return p.SyncToZK(secProductInfoList, zkPath, stat)
	}
	return p.createSyncToZK(activity)
}

func (p ActivityServiceImpl) deleteSyncToZK(activity *model.Activity) error {
	zkPath := conf.Zk.SecProductKey
	secProductInfoList, stat, _ := p.LoadProductFromZk(zkPath)
	var total int = 0
	for _, item := range secProductInfoList {
		if item.ActivityName != activity.ActivityName {
			secProductInfoList[total] = item
			total++
		}
	}
	return p.SyncToZK(secProductInfoList[:total], zkPath, stat)
}

func (p ActivityServiceImpl) LoadProductFromZk(key string) ([]*model.SecProductInfoConf, *zk.Stat, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	v, s, err := conf.Zk.ZkConn.Get(key)
	if err != nil {
		log.Printf("get [%s] from zk failed, err: %v", key, err)
		return nil, nil, err
	}
	// log.Printf("get from zk success, resp: %+v\n", s)
	// log.Printf("value of path[%s]=[%s].\n", key, v)

	var secProductInfo []*model.SecProductInfoConf
	err = json.Unmarshal(v, &secProductInfo)
	if err != nil {
		log.Printf("Unmsharl secProductInfo failed, err: %v", err)
		return nil, nil, err
	}
	return secProductInfo, s, nil
}

func (p ActivityServiceImpl) SyncToZK(secProductInfoList []*model.SecProductInfoConf, zkPath string, stat *zk.Stat) error {
	data, err := json.Marshal(secProductInfoList)
	if err != nil {
		log.Printf("json marshal failed, err: %v", err)
		return err
	}

	conn := conf.Zk.ZkConn
	var byteData = []byte(string(data))

	// 原本是先判断zkPath是否存在，不存在进行create，再进行set。
	// 现将create部分在项目初始化部分完成
	_, err_set := conn.Set(zkPath, byteData, stat.Version)
	if err_set != nil {
		fmt.Println(err_set)
	}
	log.Println("\nput [secProductInfoList] to zk success")
	return nil
}
