package tool

// 单线程互斥锁
type ChanMutex struct {
	ch chan struct{}
}

// 单线程互斥锁
func NewChanMutex() *ChanMutex {
	self := &ChanMutex{
		ch: make(chan struct{}, 1),
	}
	self.ch <- struct{}{} // 加入令牌,解锁
	return self
}

// 上锁
func (self *ChanMutex) Lock() {
	<-self.ch // 拿走令牌,上锁
}

// 解锁
func (self *ChanMutex) Unlock() {
	self.ch <- struct{}{}
}

// 尝试上锁
//   - true 上锁
func (self *ChanMutex) TryLock() bool {
	select {
	case <-self.ch: // 尝试拿走
		return true
	default:
		return false
	}
}
