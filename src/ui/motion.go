package ui

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

type MotionActionKind int

const (
	ActionNone MotionActionKind = iota
	ActionScroll
	ActionZoom
	ActionMove
)

type MotionAction struct {
	Kind   MotionActionKind
	Scroll float64
	Zoom   float64
	Pos    f32.Point
	key    key.Name
}

type Motion struct {
	tag   bool
	focus bool
}

func NewMotion() *Motion {
	return &Motion{
		tag:   true,
		focus: false,
	}
}

func (self *Motion) Update(gtx layout.Context, area image.Point) MotionAction {
	defer clip.Rect{Max: area}.Push(gtx.Ops).Pop()

	// draw_rectangle(gtx, area, colorNRGBA(255, 0, 0, 64), 0)
	// draw_rectangle_line(gtx, area, COLOR_RED, gtx.Dp(8), 12)

	event.Op(gtx.Ops, &self.tag)

	var action MotionAction
	for {
		ev, ok := gtx.Event(
			pointer.Filter{
				Target:  &self.tag,
				Kinds:   pointer.Press | pointer.Release | pointer.Scroll,
				ScrollX: pointer.ScrollRange{Min: -1 << 30, Max: 1 << 30},
				ScrollY: pointer.ScrollRange{Min: -1 << 30, Max: 1 << 30},
			},
			key.Filter{
				Focus:    &self.tag,
				Required: key.ModCtrl,
			},
			key.FocusFilter{
				Target: &self.tag,
			},
		)
		if !ok {
			break
		}

		// mouse event
		if x, ok := ev.(pointer.Event); ok {
			action.Pos = x.Position
			switch x.Kind {
			case pointer.Press:
				gtx.Execute(key.FocusCmd{Tag: &self.tag})
				if x.Buttons.Contain(pointer.ButtonTertiary) {
					action.Kind = ActionMove
				}

			case pointer.Release:

			case pointer.Scroll:
				if x.Position.Y > 0 && x.Position.X > 0 {
					action.Scroll = float64(x.Scroll.Y)
					if x.Modifiers.Contain(key.ModCtrl) {
						action.Kind = ActionZoom
					} else {
						action.Kind = ActionScroll
					}
				}
			}
		}

		// key event
		if k, ok := ev.(key.Event); ok {
			if k.State == key.Press {
				switch k.Name {
				case key.NameReturn, key.NameEnter:
				case key.NameSpace:
				case key.NameEscape:
				}
			}
		}

		if k, ok := ev.(key.Event); ok && k.State == key.Press {
			if k.Modifiers.Contain(key.ModCtrl) {
				if k.Name == "S" {
				}
			}
		}

		if f, ok := ev.(key.FocusEvent); ok {
			self.focus = f.Focus
		}

	}

	return action
}
