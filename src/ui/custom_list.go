// SPDX-License-Identifier: Unlicense OR MIT

package ui

import (
	"image"
	"math"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type iterationDir uint8

const (
	iterateNone iterationDir = iota
	iterateForward
	iterateBackward
)

type scrollChild struct {
	size image.Point
	call op.CallOp
}

const inf = 1e6

// Position is a List scroll offset represented as an offset from the top edge
// of a child element.
type Position struct {
	// BeforeEnd tracks whether the List position is before the very end. We
	// use "before end" instead of "at end" so that the zero value of a
	// Position struct is useful.
	//
	// When laying out a list, if ScrollToEnd is true and BeforeEnd is false,
	// then First and Offset are ignored, and the list is drawn with the last
	// item at the bottom. If ScrollToEnd is false then BeforeEnd is ignored.
	BeforeEnd bool
	// First is the index of the first visible child.
	First int
	// Offset is the distance in pixels from the leading edge to the child at index
	// First.
	Offset int
	// OffsetLast is the signed distance in pixels from the trailing edge to the
	// bottom edge of the child at index First+Count.
	OffsetLast int
	// Count is the number of visible children.
	Count int
	// Length is the estimated total size of all children, measured in pixels.
	Length int
}

// type scrollChild struct {
// 	size image.Point
// 	call op.CallOp
// }

// List displays a subsection of a potentially infinitely
// large underlying list. List accepts user input to scroll
// the subsection.
type CustomListWidget struct {
	Axis layout.Axis
	// ScrollToEnd instructs the list to stay scrolled to the far end position
	// once reached. A List with ScrollToEnd == true and Position.BeforeEnd ==
	// false draws its content with the last item at the bottom of the list
	// area.
	ScrollToEnd bool
	// Alignment is the cross axis alignment of list elements.
	Alignment layout.Alignment
	// ScrollAnyAxis allows any scroll axis to scroll the list, not just the main axis.
	ScrollAnyAxis bool
	// Gap is the space in pixels between children.
	Gap int

	cs          layout.Constraints
	scroll      gesture.Scroll
	scrollDelta int

	// Position is updated during Layout. To save the list scroll position,
	// just save Position after Layout finishes. To scroll the list
	// programmatically, update Position (e.g. restore it from a saved value)
	// before calling Layout.
	Position Position

	len int

	// maxSize is the total size of visible children.
	maxSize  int
	children []scrollChild
	dir      iterationDir

	pending float64
}

// ListElement is a function that computes the D of
// a list element.
type ListElement func(gtx C, index int) D

// type iterationDir uint8

// init prepares the list for iterating through its children with next.
func (l *CustomListWidget) init(gtx C, len int) {
	if l.more() {
		panic("unfinished child")
	}
	l.cs = gtx.Constraints
	l.maxSize = 0
	l.children = l.children[:0]
	l.len = len
	l.update(gtx)
	if l.Position.First < 0 {
		l.Position.Offset = 0
		l.Position.First = 0
	}
	if l.scrollToEnd() || l.Position.First > len {
		l.Position.Offset = 0
		l.Position.First = len
	}
}

// Layout a List of len items, where each item is implicitly defined
// by the callback w. Layout can handle very large lists because it only calls
// w to fill its viewport and the distance scrolled, if any.
func (l *CustomListWidget) Layout(gtx C, len int, w ListElement) D {
	l.init(gtx, len)
	crossMin, crossMax := crossConstraint(l.Axis, gtx.Constraints)
	gtx.Constraints = axisConstraints(l.Axis, 0, inf, crossMin, crossMax)
	macro := op.Record(gtx.Ops)
	laidOutTotalLength := 0
	numLaidOut := 0

	for l.next(); l.more(); l.next() {
		child := op.Record(gtx.Ops)
		dims := w(gtx, l.index())
		call := child.Stop()
		l.end(dims, call)
		laidOutTotalLength += l.Axis.Convert(dims.Size).X
		numLaidOut++
	}

	if numLaidOut > 0 {
		l.Position.Length = laidOutTotalLength*len/numLaidOut + l.Gap*(len-1)
	} else {
		l.Position.Length = 0
	}
	return l.layout(gtx.Ops, macro)
}

func (l *CustomListWidget) scrollToEnd() bool {
	return l.ScrollToEnd && !l.Position.BeforeEnd
}

// Dragging reports whether the List is being dragged.
func (l *CustomListWidget) Dragging() bool {
	return l.scroll.State() == gesture.StateDragging
}

func (l *CustomListWidget) update(gtx C) {
	min, max := int(-inf), int(inf)
	if l.Position.First == 0 {
		// Use the size of the invisible part as scroll boundary.
		min = -l.Position.Offset
		if min > 0 {
			min = 0
		}
	}
	if l.Position.First+l.Position.Count == l.len {
		max = -l.Position.OffsetLast
		if max < 0 {
			max = 0
		}
	}

	xrange := pointer.ScrollRange{Min: min, Max: max}
	yrange := pointer.ScrollRange{}

	axis := gesture.Axis(l.Axis)
	if l.ScrollAnyAxis {
		axis = gesture.Both
		yrange = xrange
	} else if l.Axis == layout.Vertical {
		xrange, yrange = yrange, xrange
	}
	d := l.scroll.Update(gtx.Metric, gtx.Source, gtx.Now, axis, xrange, yrange)

	l.scrollDelta = d

	l.pending += float64(d)
	l.pending = math.Max(float64(min), math.Min(float64(max), l.pending))

	df := l.pending * 0.15
	l.pending -= df

	l.Position.Offset += int(math.Round(df))

	if math.Abs(df) >= 0.1 {
		gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(time.Second / 60)})
	} else {
		l.pending = 0
	}
}

// next advances to the next child.
func (l *CustomListWidget) next() {
	l.dir = l.nextDir()
	// The user scroll offset is applied after scrolling to
	// list end.
	if l.scrollToEnd() && !l.more() && l.scrollDelta < 0 {
		l.Position.BeforeEnd = true
		l.Position.Offset += l.scrollDelta
		l.dir = l.nextDir()
	}
}

// index is current child's position in the underlying list.
func (l *CustomListWidget) index() int {
	switch l.dir {
	case iterateBackward:
		return l.Position.First - 1
	case iterateForward:
		return l.Position.First + len(l.children)
	default:
		panic("Index called before Next")
	}
}

// more reports whether more children are needed.
func (l *CustomListWidget) more() bool {
	return l.dir != iterateNone
}

func (l *CustomListWidget) nextDir() iterationDir {
	_, vsize := mainConstraint(l.Axis, l.cs)
	last := l.Position.First + len(l.children)
	// Clamp offset.
	if l.maxSize-l.Position.Offset < vsize && last == l.len {
		l.Position.Offset = l.maxSize - vsize
	}
	if l.Position.Offset < 0 && l.Position.First == 0 {
		l.Position.Offset = 0
	}
	// Lay out an extra (invisible) child at each end to enable focus to
	// move to them, triggering automatic scroll.
	firstSize, lastSize := 0, 0
	if len(l.children) > 0 {
		if l.Position.First > 0 {
			firstChild := l.children[0]
			firstSize = l.Axis.Convert(firstChild.size).X + l.Gap
		}
		if last < l.len {
			lastChild := l.children[len(l.children)-1]
			lastSize = l.Axis.Convert(lastChild.size).X + l.Gap
		}
	}
	switch {
	case len(l.children) == l.len:
		return iterateNone
	case l.maxSize-l.Position.Offset-lastSize < vsize:
		return iterateForward
	case l.Position.Offset-firstSize < 0:
		return iterateBackward
	}
	return iterateNone
}

// End the current child by specifying its D.
func (l *CustomListWidget) end(dims D, call op.CallOp) {
	child := scrollChild{dims.Size, call}
	mainSize := l.Axis.Convert(child.size).X
	if len(l.children) > 0 {
		l.maxSize += l.Gap
	}
	l.maxSize += mainSize
	switch l.dir {
	case iterateForward:
		l.children = append(l.children, child)
	case iterateBackward:
		l.children = append(l.children, scrollChild{})
		copy(l.children[1:], l.children)
		l.children[0] = child
		l.Position.First--
		l.Position.Offset += mainSize + l.Gap
	default:
		panic("call Next before End")
	}
	l.dir = iterateNone
}

// Layout the List and return its D.
func (l *CustomListWidget) layout(ops *op.Ops, macro op.MacroOp) D {
	if l.more() {
		panic("unfinished child")
	}
	mainMin, mainMax := mainConstraint(l.Axis, l.cs)
	children := l.children
	var first scrollChild
	// Skip invisible children.
	for len(children) > 0 {
		child := children[0]
		sz := child.size
		mainSize := l.Axis.Convert(sz).X
		if l.Position.Offset < mainSize {
			// First child is partially visible.
			break
		}
		l.Position.First++
		l.Position.Offset -= mainSize + l.Gap
		first = child
		children = children[1:]
	}
	size := -l.Position.Offset
	var maxCross int
	var last scrollChild
	for i, child := range children {
		sz := l.Axis.Convert(child.size)
		if c := sz.Y; c > maxCross {
			maxCross = c
		}
		if i > 0 {
			size += l.Gap
		}
		size += sz.X
		if size >= mainMax {
			if i < len(children)-1 {
				last = children[i+1]
			}
			children = children[:i+1]
			break
		}
	}
	l.Position.Count = len(children)
	l.Position.OffsetLast = mainMax - size
	// ScrollToEnd lists are end aligned.
	if space := l.Position.OffsetLast; l.ScrollToEnd && space > 0 {
		l.Position.Offset -= space
	}
	pos := -l.Position.Offset
	layout := func(child scrollChild) {
		sz := l.Axis.Convert(child.size)
		var cross int
		switch l.Alignment {
		case layout.End:
			cross = maxCross - sz.Y
		case layout.Middle:
			cross = (maxCross - sz.Y) / 2
		}
		childSize := sz.X
		pt := l.Axis.Convert(image.Pt(pos, cross))
		trans := op.Offset(pt).Push(ops)
		child.call.Add(ops)
		trans.Pop()
		pos += childSize
	}
	// Lay out leading invisible child.
	if first != (scrollChild{}) {
		sz := l.Axis.Convert(first.size)
		pos -= sz.X + l.Gap
		layout(first)
		pos += l.Gap
	}
	for i, child := range children {
		if i > 0 {
			pos += l.Gap
		}
		layout(child)
	}
	// Lay out trailing invisible child.
	if last != (scrollChild{}) {
		pos += l.Gap
		layout(last)
	}
	atStart := l.Position.First == 0 && l.Position.Offset <= 0
	atEnd := l.Position.First+len(children) == l.len && mainMax >= pos
	if atStart && l.scrollDelta < 0 || atEnd && l.scrollDelta > 0 {
		l.scroll.Stop()
	}
	l.Position.BeforeEnd = !atEnd
	if pos < mainMin {
		pos = mainMin
	}
	if pos > mainMax {
		pos = mainMax
	}
	if crossMin, crossMax := crossConstraint(l.Axis, l.cs); maxCross < crossMin {
		maxCross = crossMin
	} else if maxCross > crossMax {
		maxCross = crossMax
	}
	dims := l.Axis.Convert(image.Pt(pos, maxCross))
	call := macro.Stop()
	defer clip.Rect(image.Rectangle{Max: dims}).Push(ops).Pop()

	l.scroll.Add(ops)

	call.Add(ops)
	return D{Size: dims}
}

// ScrollBy scrolls the list by a relative amount of items.
//
// Fractional scrolling may be inaccurate for items of differing
// D. This includes scrolling by integer amounts if the current
// l.Position.Offset is non-zero.
func (l *CustomListWidget) ScrollBy(num float32) {
	// Split number of items into integer and fractional parts
	i, f := math.Modf(float64(num))

	// Scroll by integer amount of items
	l.Position.First += int(i)

	// Adjust Offset to account for fractional items. If Offset gets so large that it amounts to an entire item, then
	// the layout code will handle that for us and adjust First and Offset accordingly.
	itemHeight := float64(l.Position.Length) / float64(l.len)
	l.Position.Offset += int(math.Round(itemHeight * f))

	// First and Offset can go out of bounds, but the layout code knows how to handle that.

	// Ensure that the list pays attention to the Offset field when the scrollbar drag
	// is started while the bar is at the end of the list. Without this, the scrollbar
	// cannot be dragged away from the end.
	l.Position.BeforeEnd = true
}

// ScrollTo scrolls to the specified item.
func (l *CustomListWidget) ScrollTo(n int) {
	l.Position.First = n
	l.Position.Offset = 0
	l.Position.BeforeEnd = true
}

// 填充某些函数调用

// mainConstraint returns the min and max main constraints for axis a.
func mainConstraint(a layout.Axis, cs layout.Constraints) (int, int) {
	if a == layout.Horizontal {
		return cs.Min.X, cs.Max.X
	}
	return cs.Min.Y, cs.Max.Y
}

// crossConstraint returns the min and max cross constraints for axis a.
func crossConstraint(a layout.Axis, cs layout.Constraints) (int, int) {
	if a == layout.Horizontal {
		return cs.Min.Y, cs.Max.Y
	}
	return cs.Min.X, cs.Max.X
}

// axisConstraints returns the constraints for axis a.
func axisConstraints(a layout.Axis, mainMin, mainMax, crossMin, crossMax int) layout.Constraints {
	if a == layout.Horizontal {
		return layout.Constraints{Min: image.Pt(mainMin, crossMin), Max: image.Pt(mainMax, crossMax)}
	}
	return layout.Constraints{Min: image.Pt(crossMin, mainMin), Max: image.Pt(crossMax, mainMax)}
}

// list 样式 -------------------------------------------------------------------------------------

// SPDX-License-Identifier: Unlicense OR MIT

// fromListPosition converts a layout.Position into two floats representing
// the location of the viewport on the underlying content. It needs to know
// the number of elements in the list and the major-axis size of the list
// in order to do this. The returned values will be in the range [0,1], and
// start will be less than or equal to end.
func fromListPosition(lp Position, elements int, majorAxisSize int) (start, end float32) {
	// Approximate the size of the scrollable content.
	lengthEstPx := float32(lp.Length)
	elementLenEstPx := lengthEstPx / float32(elements)

	// Determine how much of the content is visible.
	listOffsetF := float32(lp.Offset)
	listOffsetL := float32(lp.OffsetLast)

	// Compute the location of the beginning of the viewport using estimated element size and known
	// pixel offsets.
	viewportStart := clamp1((float32(lp.First)*elementLenEstPx + listOffsetF) / lengthEstPx)
	viewportEnd := clamp1((float32(lp.First+lp.Count)*elementLenEstPx + listOffsetL) / lengthEstPx)
	viewportFraction := viewportEnd - viewportStart

	// Compute the expected visible proportion of the list content based solely on the ratio
	// of the visible size and the estimated total size.
	visiblePx := float32(majorAxisSize)
	visibleFraction := visiblePx / lengthEstPx

	// Compute the error between the two methods of determining the viewport and diffuse the
	// error on either end of the viewport based on how close we are to each end.
	err := visibleFraction - viewportFraction
	adjStart := viewportStart
	adjEnd := viewportEnd
	if viewportFraction < 1 {
		startShare := viewportStart / (1 - viewportFraction)
		endShare := (1 - viewportEnd) / (1 - viewportFraction)
		startErr := startShare * err
		endErr := endShare * err

		adjStart -= startErr
		adjEnd += endErr
	}
	return adjStart, adjEnd
}

// AnchorStrategy defines a means of attaching a scrollbar to content.
type AnchorStrategy uint8

const (
	// Occupy reserves space for the scrollbar, making the underlying
	// content region smaller on one axis.
	Occupy AnchorStrategy = iota
	// Overlay causes the scrollbar to float atop the content without
	// occupying any space. Content in the underlying area can be occluded
	// by the scrollbar.
	Overlay
)

// ListStyle configures the presentation of a layout.List with a scrollbar.
type CustomListStyle struct {
	state *CustomListStruct
	// CustomScrollStyle
	material.ScrollbarStyle
	AnchorStrategy
}

// List constructs a ListStyle using the provided theme and state.
func CustomList(th *material.Theme, state *CustomListStruct) CustomListStyle {
	return CustomListStyle{
		state:          state,
		ScrollbarStyle: material.Scrollbar(th, &state.Scrollbar),
	}
}

// Layout the list and its scrollbar.
func (l CustomListStyle) Layout(gtx layout.Context, length int, w ListElement) layout.Dimensions {
	originalConstraints := gtx.Constraints

	// Determine how much space the scrollbar occupies.
	barWidth := gtx.Dp(l.Width())

	if l.AnchorStrategy == Occupy {

		// Reserve space for the scrollbar using the gtx constraints.
		max := l.state.Axis.Convert(gtx.Constraints.Max)
		min := l.state.Axis.Convert(gtx.Constraints.Min)
		max.Y -= barWidth
		if max.Y < 0 {
			max.Y = 0
		}
		min.Y -= barWidth
		if min.Y < 0 {
			min.Y = 0
		}
		gtx.Constraints.Max = l.state.Axis.Convert(max)
		gtx.Constraints.Min = l.state.Axis.Convert(min)
	}

	listDims := l.state.Layout(gtx, length, w)
	gtx.Constraints = originalConstraints

	// Draw the scrollbar.
	anchoring := layout.E
	if l.state.Axis == layout.Horizontal {
		anchoring = layout.S
	}
	majorAxisSize := l.state.Axis.Convert(listDims.Size).X

	start, end := fromListPosition(l.state.Position, length, majorAxisSize)
	// layout.Direction respects the minimum, so ensure that the
	// scrollbar will be drawn on the correct edge even if the provided
	// layout.Context had a zero minimum constraint.
	gtx.Constraints.Min = listDims.Size
	if l.AnchorStrategy == Occupy {
		min := l.state.Axis.Convert(gtx.Constraints.Min)
		min.Y += barWidth
		gtx.Constraints.Min = l.state.Axis.Convert(min)
	}
	anchoring.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return l.ScrollbarStyle.Layout(gtx, l.state.Axis, start, end)
	})

	if delta := l.state.ScrollDistance(); delta != 0 {
		// Handle any changes to the list position as a result of user interaction
		// with the scrollbar.
		l.state.ScrollBy(delta * float32(length))
	}

	if l.AnchorStrategy == Occupy {
		// Increase the width to account for the space occupied by the scrollbar.
		cross := l.state.Axis.Convert(listDims.Size)
		cross.Y += barWidth
		listDims.Size = l.state.Axis.Convert(cross)
	}

	return listDims
}

// LayoutWidgets the widgets and its scrollbar.
func (l CustomListStyle) LayoutWidgets(gtx layout.Context, widgets ...layout.Widget) layout.Dimensions {
	return l.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return widgets[index](gtx)
	})
}

// 包含滚动条和列表自身
type CustomListStruct struct {
	widget.Scrollbar
	CustomListWidget
}
