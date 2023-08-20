package srv_user

import "sync"

// 用户购买记录
type UserBuyHistory struct {
	History map[int]int
	Lock    sync.RWMutex
}

// 用户购买某个产品的数量
func (u *UserBuyHistory) GetProductBuyCount(productId int) int {
	u.Lock.Lock()
	defer u.Lock.Unlock()

	count := u.History[productId]
	return count
}

// 买产品
func (u *UserBuyHistory) Add(productId, count int) {
	u.Lock.Lock()
	defer u.Lock.Unlock()

	cur, ok := u.History[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	u.History[productId] = cur
}
