package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// BookCard 代表书籍卡片的数据与交互状态
type BookCard struct {
	Title    string  // 书名
	Format   string  // 格式，如 EPUB, PDF
	Progress float32 // 进度 (0.0 ~ 1.0)

	click      widget.Clickable
	hoverAnim  float32   // 动画插值进度: 0.0 ~ 1.0
	lastUpdate time.Time // 用于计算两帧间的 Delta Time
}

// updateHoverAnim 处理平滑过渡动画
func (b *BookCard) updateHoverAnim(gtx layout.Context) {
	now := gtx.Now
	if b.lastUpdate.IsZero() {
		b.lastUpdate = now
	}
	dt := float32(now.Sub(b.lastUpdate).Seconds())
	b.lastUpdate = now

	// 目标值：悬浮时为 1.0，离开时为 0.0
	target := float32(0.0)
	if b.click.Hovered() {
		target = 1.0
	}

	// 动画过渡速度 (数值越大过渡越快)
	const speed = 10.0
	if b.hoverAnim != target {
		if b.hoverAnim < target {
			b.hoverAnim += speed * dt
			if b.hoverAnim > target {
				b.hoverAnim = target
			}
		} else {
			b.hoverAnim -= speed * dt
			if b.hoverAnim < target {
				b.hoverAnim = target
			}
		}
		// 动画未完成时，通知 Gio 请求重绘下一帧
		gtx.Execute(op.InvalidateCmd{})
	}
}

// Layout 渲染书籍卡片
func (b *BookCard) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	b.updateHoverAnim(gtx)

	// 卡片尺寸限制
	cardWidth := gtx.Dp(unit.Dp(230))
	cardHeight := gtx.Dp(unit.Dp(140))
	gtx.Constraints.Min = image.Pt(cardWidth, cardHeight)
	gtx.Constraints.Max = image.Pt(cardWidth, cardHeight)

	// 1. 悬浮微动画：向上平移 (Y 轴位移 -5dp)
	offsetY := int(float32(gtx.Dp(unit.Dp(-5))) * b.hoverAnim)
	defer op.Offset(image.Pt(0, offsetY)).Push(gtx.Ops).Pop()

	return b.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// 基础色彩（现代暗调简约风）
		baseBg := color.NRGBA{R: 0x1E, G: 0x22, B: 0x2B, A: 0xFF}
		hoverBg := color.NRGBA{R: 0x26, G: 0x2B, B: 0x36, A: 0xFF}
		cardBg := lerpColor(baseBg, hoverBg, b.hoverAnim)

		baseBorder := color.NRGBA{R: 0x2C, G: 0x32, B: 0x3F, A: 0xFF}
		hoverBorder := color.NRGBA{R: 0x4F, G: 0x80, B: 0xE1, A: 0xFF} // 悬浮时微蓝高亮
		borderColor := lerpColor(baseBorder, hoverBorder, b.hoverAnim)

		radius := gtx.Dp(unit.Dp(12))

		// 2. 绘制卡片背景与圆角剪裁
		bgClip := clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, radius).Push(gtx.Ops)
		paint.Fill(gtx.Ops, cardBg)

		// 3. 绘制 1dp 微弱边框
		strokeWidth := float32(gtx.Dp(unit.Dp(1)))
		paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
			Path:  clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, radius).Path(gtx.Ops),
			Width: strokeWidth,
		}.Op())

		bgClip.Pop()

		// 4. 内部内容排版 (带内边距)
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				// 顶部：格式 Badge
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return b.layoutBadge(gtx, th)
				}),

				// 中间：书本名称
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					titleLabel := material.Body1(th, b.Title)
					titleLabel.Color = color.NRGBA{R: 0xF0, G: 0xF3, B: 0xF6, A: 0xFF}
					titleLabel.TextSize = unit.Sp(15)
					titleLabel.MaxLines = 2
					titleLabel.Truncator = "..."
					return titleLabel.Layout(gtx)
				}),

				// 底部：阅读进度条与百分比
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return b.layoutProgress(gtx, th)
				}),
			)
		})
	})
}

// 格式 Badge 胶囊
func (b *BookCard) layoutBadge(gtx layout.Context, th *material.Theme) layout.Dimensions {
	badgeBg := color.NRGBA{R: 0x2A, G: 0x32, B: 0x41, A: 0xFF}
	badgeTextCol := color.NRGBA{R: 0x8C, G: 0xA0, B: 0xB8, A: 0xFF}

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			label := material.Caption(th, b.Format)
			label.Color = badgeTextCol
			label.TextSize = unit.Sp(10)
			return layout.Inset{
				Top: unit.Dp(2), Bottom: unit.Dp(2),
				Left: unit.Dp(6), Right: unit.Dp(6),
			}.Layout(gtx, label.Layout)
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rect := image.Rectangle{Max: gtx.Constraints.Min}
			rr := clip.UniformRRect(rect, gtx.Dp(unit.Dp(4)))
			defer rr.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, badgeBg)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
	)
}

// 进度条与百分比文字
func (b *BookCard) layoutProgress(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// 进度槽
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			height := gtx.Dp(unit.Dp(4))
			width := gtx.Constraints.Max.X
			fullRect := image.Rectangle{Max: image.Pt(width, height)}

			// 背景槽
			trackCol := color.NRGBA{R: 0x31, G: 0x37, B: 0x44, A: 0xFF}
			trackR := clip.UniformRRect(fullRect, height/2).Push(gtx.Ops)
			paint.Fill(gtx.Ops, trackCol)
			trackR.Pop()

			// 已读进度槽（主色调青蓝）
			fillWidth := int(float32(width) * b.Progress)
			if fillWidth > 0 {
				fillRect := image.Rectangle{Max: image.Pt(fillWidth, height)}
				fillR := clip.UniformRRect(fillRect, height/2).Push(gtx.Ops)
				paint.Fill(gtx.Ops, color.NRGBA{R: 0x4E, G: 0x88, B: 0xFF, A: 0xFF})
				fillR.Pop()
			}

			return layout.Dimensions{Size: image.Pt(width, height)}
		}),

		layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),

		// 百分比文字
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			percentStr := fmt.Sprintf("%d%%", int(b.Progress*100))
			pctLabel := material.Caption(th, percentStr)
			pctLabel.Color = color.NRGBA{R: 0x8C, G: 0x98, B: 0xA8, A: 0xFF}
			pctLabel.TextSize = unit.Sp(11)
			pctLabel.Alignment = text.End
			return pctLabel.Layout(gtx)
		}),
	)
}

// 颜色线性插值
func lerpColor(c1, c2 color.NRGBA, t float32) color.NRGBA {
	return color.NRGBA{
		R: uint8(float32(c1.R) + float32(int(c2.R)-int(c1.R))*t),
		G: uint8(float32(c1.G) + float32(int(c2.G)-int(c1.G))*t),
		B: uint8(float32(c1.B) + float32(int(c2.B)-int(c1.B))*t),
		A: uint8(float32(c1.A) + float32(int(c2.A)-int(c1.A))*t),
	}
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Book Card Demo"), app.Size(unit.Dp(600), unit.Dp(400)))

		th := material.NewTheme()
		th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

		// 准备两本书作为展示
		books := []*BookCard{
			{Title: "深入理解计算机系统 (CS:APP)", Format: "PDF", Progress: 0.68},
			{Title: "设计心理学 1：日常的设计", Format: "EPUB", Progress: 0.23},
		}

		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				// 整体背景色
				paint.Fill(gtx.Ops, color.NRGBA{R: 0x14, G: 0x16, B: 0x1B, A: 0xFF})

				// 居中展示卡片列表
				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:    layout.Horizontal,
						Spacing: layout.SpaceEvenly,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return books[0].Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return books[1].Layout(gtx, th)
						}),
					)
				})

				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
