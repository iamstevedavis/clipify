package ui

import (
	"image"
	"image/color"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/iamstevedavis/clipify/internal/clipboard"
)

type HistoryItem struct {
	Content string
	Time    string
}

var (
	background     = color.NRGBA{R: 100, G: 100, B: 100, A: 255} // Gray background
	cardColor      = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White card
	accentColor    = color.NRGBA{R: 66, G: 133, B: 244, A: 255}  // Blue accent
	timeColor      = color.NRGBA{R: 130, G: 130, B: 130, A: 255} // Gray text for time
	hoverColor     = color.NRGBA{R: 240, G: 245, B: 255, A: 255} // Light blue hover
	borderColor    = color.NRGBA{R: 230, G: 230, B: 230, A: 255} // Light gray border
	selectedColor  = color.NRGBA{R: 232, G: 240, B: 254, A: 255} // Light blue selected
	notificationBg = color.NRGBA{R: 76, G: 175, B: 80, A: 230}   // Success notification
)

func NewWindow(history []HistoryItem) {
	// Create a new window
	w := new(app.Window)
	w.Option(app.Title("Clipboard History"))
	w.Option(app.Size(unit.Dp(400), unit.Dp(600)))
	w.Option(app.MinSize(unit.Dp(300), unit.Dp(400)))

	var ops op.Ops
	th := material.NewTheme()

	// Create a scrollable list
	list := &widget.List{
		List: layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
		},
	}

	// Create clickable widgets for each item
	clickables := make([]widget.Clickable, len(history))

	// Notification state
	showNotification := false
	notificationText := ""
	var notificationTimer time.Time

	// Listen for events in the window
	for {
		// Get the next event
		evt := w.Event()

		// Handle the event based on its type
		switch evt := evt.(type) {
		case app.DestroyEvent:
			// Exit the loop when the window is closed
			return
		case app.FrameEvent:
			gtx := app.NewContext(&ops, evt)

			// Check if any item was clicked and copy its content
			for i := range clickables {
				if clickables[i].Clicked(gtx) {
					err := clipboard.SetClipboardContent(history[i].Content)
					if err == nil {
						showNotification = true
						notificationText = "Copied to clipboard"
						notificationTimer = time.Now()
					}
				}
			}

			// Hide notification after 2 seconds
			if showNotification && time.Since(notificationTimer) > 2*time.Second {
				showNotification = false
			}

			// Fill the background
			paint.Fill(gtx.Ops, background)

			// Add some padding around the entire UI
			layout.Inset{
				Top:    unit.Dp(16),
				Bottom: unit.Dp(16),
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					// Header
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						title := material.H5(th, "Clipboard History")
						title.Color = accentColor
						title.Alignment = text.Middle
						return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, title.Layout)
					}),

					// List
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return list.Layout(gtx, len(history), func(gtx layout.Context, index int) layout.Dimensions {
							// Render each history item as a card
							item := history[index]

							// Use the clickable for the current item
							return layout.Inset{
								Bottom: unit.Dp(8),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								// Draw card with interactive elements
								return material.Clickable(gtx, &clickables[index], func(gtx layout.Context) layout.Dimensions {
									// Determine background color based on interaction state
									bgColor := cardColor
									if clickables[index].Pressed() {
										bgColor = selectedColor
									} else if clickables[index].Hovered() {
										bgColor = hoverColor
									}

									// Card background
									rect := clip.Rect{
										Min: image.Point{X: 0, Y: 0},
										Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y},
									}
									paint.FillShape(gtx.Ops, bgColor, rect.Op())

									// Draw border around the card
									borderRect := clip.Stroke{
										Path:  rect.Path(),
										Width: float32(unit.Dp(1)),
									}.Op()
									paint.FillShape(gtx.Ops, borderColor, borderRect)

									// Card content with proper text wrapping
									return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{
											Axis: layout.Vertical,
										}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												// Ensure the text can wrap and is not truncated
												gtx.Constraints.Min.X = 0
												lbl := material.Body1(th, item.Content)

												// Remove MaxLines to show full content
												lbl.MaxLines = 0

												return lbl.Layout(gtx)
											}),
											layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												timeLabel := material.Caption(th, item.Time)
												timeLabel.Color = timeColor
												return timeLabel.Layout(gtx)
											}),
										)
									})
								})
							})
						})
					}),

					// Notification overlay (when active)
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !showNotification {
							return layout.Dimensions{}
						}

						// Position notification at the bottom
						return layout.Stack{Alignment: layout.S}.Layout(gtx,
							layout.Stacked(func(gtx layout.Context) layout.Dimensions {
								// Notification background
								notificationRect := clip.Rect{
									Min: image.Point{X: 0, Y: 0},
									Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Dp(40)},
								}
								paint.FillShape(gtx.Ops, notificationBg, notificationRect.Op())

								// Notification text
								lbl := material.Body2(th, notificationText)
								lbl.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
								lbl.Alignment = text.Middle

								return layout.Center.Layout(gtx, lbl.Layout)
							}),
						)
					}),
				)
			})

			// Acknowledge the frame event
			evt.Frame(gtx.Ops)

			// Request a redraw if we have an active notification
			if showNotification {
				w.Invalidate()
			}
		}
	}
}
