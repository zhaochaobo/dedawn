package mask

import (
	"dedawn/card"
	"dedawn/check"
	"dedawn/control"
	"dedawn/countdown"
	"dedawn/tips"
	"errors"
	"fmt"
	"gioui.org/app"
	"gioui.org/io/key"
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

type Mask struct {
	Server   string
	card     card.Card
	DeductAt time.Time
}

func (m *Mask) Run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops
	var closeButton widget.Clickable
	var cardNoInput widget.Editor
	var passwordInput widget.Editor
	messageLabel := material.Label(theme, 20, "")
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			fmt.Printf("destroy")
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			if closeButton.Clicked(gtx) {
				messageLabel.Text = ""
				cn := cardNoInput.Text()
				pw := passwordInput.Text()

				fmt.Printf("the card no. is, %s, password is %s\n", cn, pw)
				c, err := check.Check("http://localhost:8081", cn, pw)
				if err != nil {
					fmt.Printf("check card failed, %v", err)
					messageLabel.Text = err.Error()
				} else {
					// 记录点卡时长，并关闭窗口
					m.card = c
					fmt.Printf("che card info %v", c)
					window.Perform(system.ActionClose)
				}

			}

			// Define an large label with an appropriate text:
			title := material.H1(theme, "Avoid Dawn Surfing System")

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon

			// Change the position of the label.
			title.Alignment = text.Middle

			// Draw the label to the graphics context.
			title.Layout(gtx)

			th := material.NewTheme()

			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				// We insert two rigid elements:
				// First one to hold a button ...
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						messageLabel.Alignment = text.Middle
						messageLabel.Color = color.NRGBA{R: 204, G: 10, B: 10, A: 255}
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
								return border.Layout(gtx, messageLabel.Layout)
							},
						)

					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						cn := material.Editor(th, &cardNoInput, "Card NO.")
						cardNoInput.LineHeight = 32
						cardNoInput.InputHint = key.HintNumeric
						cn.TextSize = 32
						cn.LineHeight = 32

						margins := layout.Inset{
							Top:    unit.Dp(0),
							Right:  unit.Dp(170),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(170),
						}
						// ... and borders ...
						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(1),
						}
						// ... before laying it out, one inside the other
						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return border.Layout(gtx, cn.Layout)
							},
						)

					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {

						pw := material.Editor(th, &passwordInput, "Secret")
						passwordInput.Mask = '*'
						passwordInput.InputHint = key.HintPassword
						passwordInput.LineHeight = 32
						passwordInput.SingleLine = true
						pw.TextSize = 32
						pw.LineHeight = 32

						margins := layout.Inset{
							Top:    unit.Dp(0),
							Right:  unit.Dp(170),
							Bottom: unit.Dp(10),
							Left:   unit.Dp(170),
						}
						// ... and borders ...
						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(1),
						}
						// ... before laying it out, one inside the other
						return margins.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return border.Layout(gtx, pw.Layout)
							},
						)

					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {

						margins := layout.Inset{
							Top:    unit.Dp(25),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(170),
							Left:   unit.Dp(170),
						}
						// TWO: ... then we lay out those margins ...
						return margins.Layout(gtx,
							// THREE: ... and finally within the margins, we ddefine and lay out the button
							func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &closeButton, "start")
								return btn.Layout(gtx)
							},
						)

					},
				),
				// ... then one to hold an empty spacer
				layout.Rigid(
					// The height of the spacer is 25 Device independent pixels
					layout.Spacer{Height: unit.Dp(25)}.Layout,
				),
			)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}

func (m *Mask) Wait() {
	defer func() {
		m.card = card.Card{} // clean card
	}()
	if len(m.card.No) == 0 {
		return
	}
	defer func() {
		countdown.Tips("%s 后将锁屏，尽快保存好个人数据并下机")
	}()
	ticker := time.NewTicker(5 * time.Second)
	hasTipped := false
	for {
		select {
		case <-ticker.C:
			if _, ok := card.Admin(m.card.No, m.card.Secret); ok {
				continue
			}
			deductAt := time.Now()
			amount := 5 * time.Second
			if !m.DeductAt.IsZero() {
				amount = deductAt.Sub(m.DeductAt)
			}
			c, err := check.Deduct(m.Server, m.card.No, amount)
			if err != nil {
				if errors.Is(err, check.ErrCardDepleted) {
					m.card = card.Card{}
					log.Println("card depleted")
					return
				}
				log.Printf("deduct card amount failed, %v", err)
				continue
			}
			m.card = c
			m.DeductAt = deductAt
			if c.Amount < 5*time.Minute && !hasTipped {
				hasTipped = true
				tips.Tips(fmt.Sprintf("玩耍时间还剩%s，即将锁屏，请保存好个人数据", c.Amount))
			}
			log.Printf("status is %v", c)
		}
	}
}

func (m *Mask) Pause() {
	m.card = card.Card{}
}

func Run(server string) {
	go func() {
		m := Mask{Server: server}
		for {
			window := new(app.Window)
			//window.Option(app.Fullscreen.Option())
			err := m.Run(window)
			if err != nil {
				log.Fatal(err)
			}
			control.Start(&m.card)
			m.Wait()
		}

		//os.Exit(0)
	}()
}
