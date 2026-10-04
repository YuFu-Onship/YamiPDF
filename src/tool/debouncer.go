package tool

import (
	"sync"
	"time"
)

// 防抖,在操作结束后,指定时间内执行回调
type Debouncer struct {
	mu    sync.Mutex
	delay time.Duration
	timer *time.Timer
}

func NewDebouncer(delay time.Duration) *Debouncer {
	return &Debouncer{delay: delay}
}

// 触发操作时调用
func (self *Debouncer) Do(call_fn func()) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.timer != nil {
		self.timer.Stop()
	}
	self.timer = time.AfterFunc(self.delay, call_fn)
}

// 取消任务
func (self *Debouncer) Stop() {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.timer != nil {
		self.timer.Stop()
	}
}
