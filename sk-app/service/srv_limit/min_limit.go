package srv_limit

type TimeLimit interface {
	Count(nowTime int64) (curCount int)
	Check(nowTime int64) int
}

// 分钟限制
type MinLimit struct {
	count   int
	curTime int64
}

// 在1分钟内访问次数
func (p *MinLimit) Count(nowTime int64) (curCount int) {
	if nowTime-p.curTime > 60 {
		p.count = 1
		p.curTime = nowTime
		curCount = p.count
		return
	}
	p.count++
	curCount = p.count
	return
}

// 检查用户访问次数
func (p *MinLimit) Check(nowTime int64) int {
	if nowTime-p.curTime > 60 {
		return 0
	}
	return p.count
}
