package move

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ButtonAlign int

type ButtonIconPlacement int

const (
	ButtonAlignCenter ButtonAlign = iota
	ButtonAlignLeading
	ButtonAlignTrailing
)

const (
	ButtonIconLeadingText ButtonIconPlacement = iota
	ButtonIconTrailingText
)

type NeuButton struct {
	widget.DisableableWidget
	Text          string
	Icon          fyne.Resource
	Alignment     ButtonAlign
	IconPlacement ButtonIconPlacement

	OnTapped func() `json:"-"`

	hovered bool
	tapAnim *fyne.Animation
	tapBG   *canvas.Rectangle
}

func NewNeuButton(label string, tapped func()) *NeuButton {
	button := &NeuButton{
		Text:     label,
		OnTapped: tapped,
	}

	button.ExtendBaseWidget(button)
	return button
}

func NewNeuButtonWithIcon(label string, icon fyne.Resource, tapped func()) *NeuButton {
	button := &NeuButton{
		Text:     label,
		Icon:     icon,
		OnTapped: tapped,
	}

	button.ExtendBaseWidget(button)
	return button
}

func (b *NeuButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	text := canvas.NewText(b.Text, theme.ForegroundColor())
	text.TextStyle.Bold = true

	background := canvas.NewImageFromResource(ResourceButtonPng)
	b.tapBG = canvas.NewRectangle(color.Transparent)
	objects := []fyne.CanvasObject{
		background,
		b.tapBG,
		text,
	}
	r := &neuButtonRenderer{
		background: background,
		button:     b,
		label:      text,
		layout:     layout.NewHBoxLayout(),
		objects:    objects,
	}

	r.updateIconAndText()
	r.applyTheme()
	return r
}

func (b *NeuButton) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

func (b *NeuButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

func (b *NeuButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

func (b *NeuButton) MouseMoved(*desktop.MouseEvent) {
}

func (b *NeuButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

func (b *NeuButton) SetIcon(icon fyne.Resource) {
	b.Icon = icon

	b.Refresh()
}

func (b *NeuButton) SetText(text string) {
	b.Text = text

	b.Refresh()
}

func (b *NeuButton) Tapped(pe *fyne.PointEvent) {
	if b.Disabled() {
		return
	}

	b.tapAnimation()
	b.Refresh()

	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *NeuButton) tapAnimation() {
	if b.tapBG == nil {
		return
	}

	if b.tapAnim == nil {
		b.tapAnim = newButtonTapAnimation(b.tapBG, b)
		b.tapAnim.Curve = fyne.AnimationEaseOut
	} else {
		b.tapAnim.Stop()
	}

	b.tapAnim.Start()
}

type neuButtonRenderer struct {
	icon       *canvas.Image
	label      *canvas.Text
	background *canvas.Image
	button     *NeuButton
	layout     fyne.Layout
	objects    []fyne.CanvasObject
}

func (r *neuButtonRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	hasIcon := r.icon != nil
	hasLabel := r.label.Text != ""
	if !hasIcon && !hasLabel {
		// Nothing to layout
		return
	}
	iconSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	labelSize := r.label.MinSize()
	padding := r.padding()
	if hasLabel {
		if hasIcon {
			// Both
			var objects []fyne.CanvasObject
			if r.button.IconPlacement == ButtonIconLeadingText {
				objects = append(objects, r.icon, r.label)
			} else {
				objects = append(objects, r.label, r.icon)
			}
			r.icon.SetMinSize(iconSize)
			min := r.layout.MinSize(objects)
			r.layout.Layout(objects, min)
			pos := alignedPosition(r.button.Alignment, padding, min, size)
			r.label.Move(r.label.Position().Add(pos))
			r.icon.Move(r.icon.Position().Add(pos))
		} else {
			// Label Only
			r.label.Move(alignedPosition(r.button.Alignment, padding, labelSize, size))
			r.label.Resize(labelSize)
		}
	} else {
		// Icon Only
		r.icon.Move(alignedPosition(r.button.Alignment, padding, iconSize, size))
		r.icon.Resize(iconSize)
	}
}

func (r *neuButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(250, 250)
}

func (r *neuButtonRenderer) Refresh() {
	r.label.Text = r.button.Text
	r.updateIconAndText()
	r.applyTheme()
	r.background.Refresh()
	r.Layout(r.button.Size())
	canvas.Refresh(r.button)
}

func (r *neuButtonRenderer) Destroy() {
}

func (r *neuButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *neuButtonRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

func (r *neuButtonRenderer) applyTheme() {
	r.label.TextSize = theme.TextSize()
	r.label.Color = theme.ForegroundColor()
	switch {
	case r.button.Disabled():
		r.label.Color = theme.DisabledColor()
	}
}

func (r *neuButtonRenderer) buttonColor() color.Color {
	switch {
	case r.button.Disabled():
		return theme.DisabledButtonColor()
	case r.button.hovered:
		return theme.HoverColor()
	default:
		return theme.ButtonColor()
	}
}

func (r *neuButtonRenderer) padding() fyne.Size {
	if r.button.Text == "" {
		return fyne.NewSize(theme.Padding()*4, theme.Padding()*4)
	}
	return fyne.NewSize(theme.Padding()*6, theme.Padding()*4)
}

func (r *neuButtonRenderer) updateIconAndText() {
	if r.button.Icon != nil && r.button.Visible() {
		if r.icon == nil {
			r.icon = canvas.NewImageFromResource(r.button.Icon)
			r.icon.FillMode = canvas.ImageFillContain
			r.SetObjects([]fyne.CanvasObject{r.background, r.button.tapBG, r.label, r.icon})
		}
		if r.button.Disabled() {
			r.icon.Resource = theme.NewDisabledResource(r.button.Icon)
		} else {
			r.icon.Resource = r.button.Icon
		}
		r.icon.Refresh()
		r.icon.Show()
	} else if r.icon != nil {
		r.icon.Hide()
	}
	if r.button.Text == "" {
		r.label.Hide()
	} else {
		r.label.Show()
	}
}

func alignedPosition(align ButtonAlign, padding, objectSize, layoutSize fyne.Size) (pos fyne.Position) {
	pos.Y = (layoutSize.Height - objectSize.Height) / 2
	switch align {
	case ButtonAlignCenter:
		pos.X = (layoutSize.Width - objectSize.Width) / 2
	case ButtonAlignLeading:
		pos.X = padding.Width / 2
	case ButtonAlignTrailing:
		pos.X = layoutSize.Width - objectSize.Width - padding.Width/2
	}
	return
}

func newButtonTapAnimation(bg *canvas.Rectangle, w fyne.Widget) *fyne.Animation {
	return fyne.NewAnimation(canvas.DurationStandard, func(done float32) {
		mid := (w.Size().Width - theme.Padding()) / 2
		size := mid * done
		bg.Resize(fyne.NewSize(size*2, w.Size().Height-theme.Padding()))
		bg.Move(fyne.NewPos(mid-size, theme.Padding()/2))

		r, g, bb, a := theme.PressedColor().RGBA()
		aa := uint8(a)
		fade := aa - uint8(float32(aa)*done)
		bg.FillColor = &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(bb), A: fade}
		canvas.Refresh(bg)
	})
}
