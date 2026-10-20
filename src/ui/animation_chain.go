package ui

// 一个列表的缩进动画
type AniEaseList struct {
	Style *Style

	Num            int     // 动画总数
	Coe            float64 // 缩近系数
	Space          float64 // 初始间距
	Tigger_percent float64 // [0~1],当当前小于该百分比时,触发下一个

	list_cur []float64
}

func NewAniEaseList(style *Style, params AniEaseList) *AniEaseList {
	self := AniEaseList{
		Num:            params.Num,
		Coe:            params.Coe,
		Space:          params.Space,
		Tigger_percent: params.Tigger_percent,
		list_cur:       []float64{},
	}

	self.Style = style
	self.Coe = Ternary(self.Coe == 0, 0.4, self.Coe)
	self.Space = Ternary(self.Space == 0, 100, self.Space)
	self.Num = Ternary(self.Num == 0, 0, self.Num)
	self.Tigger_percent = Ternary(self.Tigger_percent == 0, 0.3, self.Tigger_percent)

	for range self.Num {
		self.list_cur = append(self.list_cur, self.Space)
	}

	return &self
}

// 更新参数, 不返回任何值, 习惯性使用Layout来更新, 而非Update
func (self *AniEaseList) Layout(gtx C) D {
	if len(self.list_cur) == 0 {
		return D{}
	}

	limit := self.Space * self.Tigger_percent

	for i := range len(self.list_cur) {
		if self.list_cur[i] == 0 {
			continue
		}

		if i == 0 {
			self.Style.AniSys.Ease64(&self.list_cur[i], 1, &self.list_cur[i], 0, self.Coe)
		} else {
			if self.list_cur[i-1] <= limit {
				self.Style.AniSys.Ease64(&self.list_cur[i], 1, &self.list_cur[i], 0, self.Coe)
			}
		}
	}

	return D{}
}

// 触发
func (self *AniEaseList) Reset() {
	if self.Num == 0 {
		return
	}

	for i := range len(self.list_cur) {
		self.list_cur[i] = self.Space
	}
}

// 返回特定值,未作分辨率处理
func (self *AniEaseList) API_get_value(index int) float64 {
	if index >= len(self.list_cur) {
		return 0
	}

	return self.list_cur[index]
}
