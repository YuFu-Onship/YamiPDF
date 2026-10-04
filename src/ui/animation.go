package ui

import (
	"math"
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
}

func NewAniSystem() *AniSystem {
	self := AniSystem{
		EnableAni: true,
		AniList:   map[any]AniDriver{},
	}
	return &self
}

func (self *AniSystem) Run(gtx C) {
	finish := true
	// log.Println("当前的动画数量:", len(self.AniList))
	for tag, item := range self.AniList {
		result := item.Update()
		finish = result && finish
		if result {
			delete(self.AniList, tag)
		}
	}

	// log.Println(len(self.AniList))

	if !finish {
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
	}
}

// 动画曲线 -----------------------------------------------------
type AniCurve struct {
	values []float64
	limit  int
	track  *float64
}

// 绘制动画曲线在窗口上
//   - track: 需要跟踪的目标值
//   - limit: 记录的数据量
//   - interval: 时间间隔(帧)
func NewAniCurve(track *float64, limit int, interval int) *AniCurve {
	self := AniCurve{
		track: track,
		limit: limit,
	}
	return &self
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
}

func (self *EaseGroupDriver2) Update() bool {
	result := true
	for i := range self.cur {
		self.cur[i] += (self.tar[i] - self.cur[i]) * self.coe
		if math.Abs(self.cur[i]-self.tar[i]) < self.precision {
			self.cur[i] = self.tar[i]
		} else {
			result = false
		}
	}

	return result
}

// 组缓动,不使用指针
func (self *AniSystem) EaseGroup(tag any, cur []float64, tar []float64, precision float64, coe float64) []float64 {
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
		}

		return cur
	}
}

// 缓出64 计算增量 -----------------------------------------------
type Ease64v2Driver struct {
	cur       *float64
	tar       float64
	coe       float64
	precision float64
}

func (self *Ease64v2Driver) Update() bool {
	if math.Abs(*self.cur-self.tar) <= self.precision {
		*self.cur = self.tar
		return true
	}
	*self.cur += (self.tar - *self.cur) * self.coe
	return false
}

// 通过增量进行处理
//   - precision 精度
//   - cur 当前值
//   - delta 增量
//   - coe 系数
func (self *AniSystem) Ease64v2(tag any, precision float64, cur *float64, delta float64, coe float64) {
	if delta == 0.0 {
		return
	}

	if driver, exist := self.AniList[tag].(*Ease64v2Driver); exist {
		// if *driver.cur == driver.tar {
		// 	return
		// }
		driver.cur = cur
		driver.coe = coe
		driver.precision = precision
		driver.tar = *cur + delta
	} else {
		self.AniList[tag] = &Ease64v2Driver{
			cur:       cur,
			tar:       *cur + delta,
			precision: precision,
			coe:       coe,
		}
	}

}

// 缓出32 ---------------------------------------------------------
type Ease32Driver struct {
	Cur       *float32
	Tar       float32
	Coe       float32
	Precision float64
}

func (self *Ease32Driver) Update() bool {
	if math.Abs(float64(*self.Cur-self.Tar)) <= self.Precision {
		*self.Cur = self.Tar
		return true
	} else {
		*self.Cur += (self.Tar - *self.Cur) * self.Coe
		return false
	}
}

// 缓出曲线
//   - precision 精度
func (self *AniSystem) Ease32(tag any, precision float64, cur *float32, tar float32, coe float32) {
	if *cur == tar {
		return
	}

	if driver, exists := self.AniList[tag].(*Ease32Driver); exists {
		driver.Cur = cur
		driver.Tar = tar
		driver.Coe = coe
		driver.Precision = precision
	} else {
		self.AniList[tag] = &Ease32Driver{
			Cur:       cur,
			Tar:       tar,
			Coe:       coe,
			Precision: precision,
		}
	}
}

// 贝塞尔曲线 ------------------------------------------------------
type CubicBezierDriver struct {
	Cur *float64
	Tar float64

	timer     float64
	limit     float64
	params    [4]float64 // [x1,y1,x2,y2]
	Precision float64
}

func (self *CubicBezierDriver) Update() bool {
	if math.Abs(self.Tar-*self.Cur) < self.Precision {
		*self.Cur = self.Tar
		return true
	}

	self.timer += 1.0
	t := math.Min(1.0, math.Max(0.0, self.timer/self.limit))
	value := 3*(1-t)*(1-t)*t*self.params[1] + 3*(1-t)*t*t*self.params[3] + t*t*t
	*self.Cur += (self.Tar - *self.Cur) * value
	return false
}

var (
	Bezier_ease_in_out_1   = [4]float64{0.3, 0.0, 0.0, 1.0}
	Bezier_ease_in_out     = [4]float64{0.65, 0.0, 0.35, 0.85}
	Bezier_back_out_light  = [4]float64{0.34, 1.3, 0.64, 1.0}
	Bezier_back_out        = [4]float64{0.175, 0.885, 0.32, 1.275}
	Bezier_back_out_strong = [4]float64{0.68, -0.55, 0.265, 1.55}
)

// 三次贝塞尔曲线
//   - limit 动画持续时间(单位:帧)
//   - params [x1,y1,x2,y2](0~1)
//   - params : Bezier_ease_in_out ...
func (self *AniSystem) AniCubicBezier(tag any, precision float64, cur *float64, tar float64, limit float64, params [4]float64) {
	if tar == *cur {
		return
	}
	driver, exist := self.AniList[tag].(*CubicBezierDriver)
	if !exist {
		self.AniList[tag] = &CubicBezierDriver{
			Cur:       cur,
			Tar:       tar,
			limit:     limit,
			timer:     0,
			params:    params,
			Precision: precision,
		}
	} else {
		driver.Cur = cur
		driver.Tar = tar
		driver.limit = limit
		driver.params = params
		driver.Precision = precision
	}
}
