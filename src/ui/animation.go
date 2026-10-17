package ui

import (
	"math"
	"slices"
	"time"

	"gioui.org/op"
)

// 动画系统接口 --------------------------------------------------
type AniDriver interface {
	Update() bool
}

// 动画系统 ------------------------------------------------------
type AniSystem struct {
	EnableAni bool
	AniList   map[any]AniDriver
	nextTick  time.Time

	last_time time.Time
	timer     time.Duration
}

func NewAniSystem() *AniSystem {
	self := AniSystem{
		EnableAni: true,
		AniList:   map[any]AniDriver{},
	}
	return &self
}

func (self *AniSystem) Run(gtx C) {
	// log.Println(len(self.AniList))

	if self.AniList == nil || len(self.AniList) == 0 {
		return
	}

	now := time.Now()
	self.last_time = Ternary(self.last_time.IsZero(), now, self.last_time) // 防止每次启动动画时, 第一帧动画出现异常增量

	dt := min(now.Sub(self.last_time), time.Millisecond*200)

	self.last_time = now
	self.timer += dt

	const step = time.Second / 60
	if self.timer >= step {
		self.timer -= step
		for tag, item := range self.AniList {
			if item.Update() {
				delete(self.AniList, tag)
			}
		}
	}

	if len(self.AniList) > 0 {
		remaining := max(step-self.timer, 0)
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(remaining)})
	} else {
		self.last_time = time.Time{} // 设置为绝对时间点
		self.timer = 0
	}
}

// 缓出 ----------------------------------------------------------
type Ease64Driver struct {
	Cur       *float64
	Tar       float64
	Coe       float64
	Precision float64
}

func (self *Ease64Driver) Update() bool {
	if math.Abs(*self.Cur-self.Tar) <= self.Precision {
		*self.Cur = self.Tar
		return true
	} else {
		*self.Cur += (self.Tar - *self.Cur) * self.Coe
		return false
	}
}

// 缓出曲线
//   - precision 精度
//   - 当前值
//   - 目标值
//   - 缓动系数
func (self *AniSystem) Ease64(tag any, precision float64, cur *float64, tar float64, coe float64) {
	if *cur == tar {
		if driver, exist := self.AniList[tag].(*Ease64Driver); exist {
			*driver.Cur = driver.Tar
		}
		return
	}

	if driver, exists := self.AniList[tag].(*Ease64Driver); exists {
		driver.Cur = cur
		driver.Tar = tar
		driver.Coe = coe
		driver.Precision = precision
	} else {
		self.AniList[tag] = &Ease64Driver{
			Cur:       cur,
			Tar:       tar,
			Coe:       coe,
			Precision: precision,
		}
	}
}

// 缓出组 ----------------------------------------------------------
//
//	多项值共享同一驱动同步缓动: 每帧所有项按相同系数向各自目标收敛,
//	全部达到各自精度后才在同一帧结束, 避免某项提前 snap 导致数值间
//	失步 (如缩放时 scale 已停而 offset 仍在追, 造成锚点漂移晃动)。
type EaseItem struct {
	Cur       *float64
	Tar       float64
	Precision float64
}

type EaseGroupDriver struct {
	Items []EaseItem
	Coe   float64
}

func (self *EaseGroupDriver) Update() bool {
	for i := range self.Items {
		it := &self.Items[i]
		*it.Cur += (it.Tar - *it.Cur) * self.Coe
	}

	for i := range self.Items {
		it := &self.Items[i]
		if math.Abs(*it.Cur-it.Tar) > it.Precision {
			return false
		}
	}
	for i := range self.Items {
		it := &self.Items[i]
		*it.Cur = it.Tar
	}

	return true
}

// 组缓动
//   - items 需要同步缓动的值集合(每项携带各自精度)
//   - coe 缓动系数, 所有项共用
func (self *AniSystem) Ease64Group(tag any, items []EaseItem, coe float64) {
	if driver, exists := self.AniList[tag].(*EaseGroupDriver); exists {
		driver.Items = items
		driver.Coe = coe
	} else {
		self.AniList[tag] = &EaseGroupDriver{
			Items: items,
			Coe:   coe,
		}
	}
}

type EaseGroupDriver2 struct {
	cur       []float64
	tar       []float64
	precision float64
	coe       float64
	value     int
}

func (self *EaseGroupDriver2) Update() bool {
	result := true

	for i := range self.cur {
		self.cur[i] += (self.tar[i] - self.cur[i]) * self.coe

		if math.Abs(self.cur[i]-self.tar[i]) < self.precision {
			self.cur[i] = self.tar[i]
			self.value += 1
		} else {
			result = false
		}
	}

	if self.value <= len(self.cur) {
		return false
	}

	return result
}

// 组缓动,不使用指针
func (self *AniSystem) EaseGroup(tag any, cur []float64, tar []float64, precision float64, coe float64) []float64 {
	if slices.Equal(cur, tar) {
		if _, exist := self.AniList[tag].(*EaseGroupDriver2); exist {
			delete(self.AniList, tag)
		}
		return tar
	}

	if driver, exist := self.AniList[tag].(*EaseGroupDriver2); exist {
		driver.tar = tar
		driver.coe = coe
		driver.precision = precision
		return driver.cur
	} else {
		self.AniList[tag] = &EaseGroupDriver2{
			cur:       cur,
			tar:       tar,
			coe:       coe,
			precision: precision,
			value:     0,
		}
		return cur
	}
}

// 缓出64 计算增量 -----------------------------------------------
// type Ease64v2Driver struct {
// 	cur       *float64
// 	tar       float64
// 	coe       float64
// 	precision float64
// }
//
// func (self *Ease64v2Driver) Update() bool {
// 	if math.Abs(*self.cur-self.tar) <= self.precision {
// 		*self.cur = self.tar
// 		return true
// 	}
// 	*self.cur += (self.tar - *self.cur) * self.coe
// 	return false
// }
//
// // 通过增量进行处理
// //   - precision 精度
// //   - cur 当前值
// //   - delta 增量
// //   - coe 系数
// func (self *AniSystem) Ease64v2(tag any, precision float64, cur *float64, delta float64, coe float64) {
// 	if delta == 0.0 {
// 		return
// 	}
//
// 	if driver, exist := self.AniList[tag].(*Ease64v2Driver); exist {
// 		// if *driver.cur == driver.tar {
// 		// 	return
// 		// }
// 		driver.cur = cur
// 		driver.coe = coe
// 		driver.precision = precision
// 		driver.tar = *cur + delta
// 	} else {
// 		self.AniList[tag] = &Ease64v2Driver{
// 			cur:       cur,
// 			tar:       *cur + delta,
// 			precision: precision,
// 			coe:       coe,
// 		}
// 	}
//
// }
