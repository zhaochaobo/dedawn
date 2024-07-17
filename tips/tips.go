package tips

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	log "github.com/sirupsen/logrus"
	"image/color"
	"time"
)

type TipWindow struct {
	Message string
	w       *app.Window
}

func (t TipWindow) run(w *app.Window) error {
	theme := material.NewTheme()

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Define an large label with an appropriate text:
			title := material.H6(theme, "Tips")

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon

			// Change the position of the label.
			title.Alignment = text.Middle

			// Draw the label to the graphics context.
			title.Layout(gtx)

			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				// ... then one to hold an empty spacer
				layout.Rigid(
					// The height of the spacer is 25 Device independent pixels
					layout.Spacer{Height: unit.Dp(75)}.Layout,
				),
				// We insert two rigid elements:
				// First one to hold a button ...
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						msg := material.Label(theme, 15, t.Message)

						msg.Alignment = text.Middle
						msg.Color = color.NRGBA{R: 204, G: 10, B: 10, A: 255}
						margins := layout.Inset{
							Top:    unit.Dp(0),
							Right:  unit.Dp(10),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(10),
						}
						// ... and borders ...
						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(0),
						}
						// ... before laying it out, one inside the other
						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return border.Layout(gtx, msg.Layout)
							},
						)

					},
				),
			)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

func Tips(message string) {
	t := TipWindow{
		Message: message,
	}
	go func() {

		for {
			window := new(app.Window)
			window.Option(app.Size(300, 180))
			t.w = window
			err := t.run(window)
			if err != nil {
				log.Fatal(err)
			}
			return

		}

	}()

	<-time.After(5 * time.Second)

	t.w.Perform(system.ActionClose)

}
