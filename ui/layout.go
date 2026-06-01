package ui

import (
	"image"
	"image/color"
	"strconv"

	"lottery.app/logic"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type LotteryUI struct {
	theme     *material.Theme
	randomBtn widget.Clickable
	dantuoBtn widget.Clickable

	// 胆码拖码输入框
	frontDare1 widget.Editor
	frontDare2 widget.Editor
	frontDare3 widget.Editor
	frontDare4 widget.Editor

	frontDrag1 widget.Editor
	frontDrag2 widget.Editor
	frontDrag3 widget.Editor
	frontDrag4 widget.Editor
	frontDrag5 widget.Editor
	frontDrag6 widget.Editor
	frontDrag7 widget.Editor
	frontDrag8 widget.Editor

	backDare1 widget.Editor
	backDrag1 widget.Editor
	backDrag2 widget.Editor
	backDrag3 widget.Editor

	currentResult string
	showDanTuoUI  bool
	errorMsg      string
	generateBtn   widget.Clickable

	list layout.List // 添加 List 作为字段
}

func NewLotteryUI() *LotteryUI {
	ui := &LotteryUI{
		showDanTuoUI: false,
		list: layout.List{
			Axis: layout.Vertical,
		},
	}
	ui.initEditors()
	return ui
}

func (ui *LotteryUI) initEditors() {
	editors := []*widget.Editor{
		&ui.frontDare1, &ui.frontDare2, &ui.frontDare3, &ui.frontDare4,
		&ui.frontDrag1, &ui.frontDrag2, &ui.frontDrag3, &ui.frontDrag4,
		&ui.frontDrag5, &ui.frontDrag6, &ui.frontDrag7, &ui.frontDrag8,
		&ui.backDare1, &ui.backDrag1, &ui.backDrag2, &ui.backDrag3,
	}
	for _, ed := range editors {
		ed.SingleLine = true
	}
}

func (ui *LotteryUI) Run(window *app.Window) error {
	ui.theme = material.NewTheme()
	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ui.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (ui *LotteryUI) layout(gtx layout.Context) layout.Dimensions {
	// 使用字段中的 List
	return ui.list.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(ui.drawHeader),
			layout.Rigid(ui.drawButtons),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if ui.showDanTuoUI {
					return ui.drawDanTuoInput(gtx)
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(ui.drawResult),
		)
	})
}

func (ui *LotteryUI) drawHeader(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.H1(ui.theme, "大乐透随机号码生成器")
		label.Color = color.NRGBA{R: 255, G: 100, B: 100, A: 255}
		return label.Layout(gtx)
	})
}

func (ui *LotteryUI) drawButtons(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(ui.theme, &ui.randomBtn, "随机一注")
			if ui.randomBtn.Clicked(gtx) {
				ui.generateRandom()
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(ui.theme, &ui.dantuoBtn, "胆拖选号")
			if ui.dantuoBtn.Clicked(gtx) {
				ui.showDanTuoUI = !ui.showDanTuoUI
				ui.errorMsg = ""
			}
			return btn.Layout(gtx)
		}),
	)
}

func (ui *LotteryUI) drawDanTuoInput(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(ui.drawDanTuoTitle),
			layout.Rigid(ui.drawFrontArea),
			layout.Rigid(ui.drawBackArea),
			layout.Rigid(ui.drawGenerateBtn),
			layout.Rigid(ui.drawErrorMsg),
		)
	})
}

func (ui *LotteryUI) drawDanTuoTitle(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.H5(ui.theme, "胆拖选号（前区胆码1-4个，后区胆码1个）")
		return label.Layout(gtx)
	})
}

func (ui *LotteryUI) drawFrontArea(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(ui.theme, "前区胆码（1-35，每个框一个数字）")
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(ui.editorWrapper(&ui.frontDare1, "胆码1")),
				layout.Rigid(ui.editorWrapper(&ui.frontDare2, "胆码2")),
				layout.Rigid(ui.editorWrapper(&ui.frontDare3, "胆码3")),
				layout.Rigid(ui.editorWrapper(&ui.frontDare4, "胆码4")),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(ui.theme, "前区拖码（1-35，每个框一个数字）")
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(ui.editorWrapper(&ui.frontDrag1, "拖码1")),
				layout.Rigid(ui.editorWrapper(&ui.frontDrag2, "拖码2")),
				layout.Rigid(ui.editorWrapper(&ui.frontDrag3, "拖码3")),
				layout.Rigid(ui.editorWrapper(&ui.frontDrag4, "拖码4")),
				layout.Rigid(ui.editorWrapper(&ui.frontDrag5, "拖码5")),
				layout.Rigid(ui.editorWrapper(&ui.frontDrag6, "拖码6")),
				layout.Rigid(ui.editorWrapper(&ui.frontDrag7, "拖码7")),
				layout.Rigid(ui.editorWrapper(&ui.frontDrag8, "拖码8")),
			)
		}),
	)
}

func (ui *LotteryUI) drawBackArea(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(ui.theme, "后区胆码（1-12，1个）")
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(ui.editorWrapper(&ui.backDare1, "胆码")),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(ui.theme, "后区拖码（1-12，每个框一个数字）")
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(ui.editorWrapper(&ui.backDrag1, "拖码1")),
				layout.Rigid(ui.editorWrapper(&ui.backDrag2, "拖码2")),
				layout.Rigid(ui.editorWrapper(&ui.backDrag3, "拖码3")),
			)
		}),
	)
}

func (ui *LotteryUI) editorWrapper(editor *widget.Editor, hint string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// 设置最小高度便于触摸
		gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(48))

		border := widget.Border{
			Color:        color.NRGBA{A: 255},
			CornerRadius: unit.Dp(4),
			Width:        unit.Dp(1),
		}
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				ed := material.Editor(ui.theme, editor, hint)
				ed.TextSize = unit.Sp(14)
				return ed.Layout(gtx)
			})
		})
	}
}

func (ui *LotteryUI) drawGenerateBtn(gtx layout.Context) layout.Dimensions {
	btn := material.Button(ui.theme, &ui.generateBtn, "生成胆拖号码")
	if ui.generateBtn.Clicked(gtx) {
		ui.generateDanTuo()
	}
	return layout.Center.Layout(gtx, btn.Layout)
}

func (ui *LotteryUI) drawErrorMsg(gtx layout.Context) layout.Dimensions {
	if ui.errorMsg != "" {
		label := material.Caption(ui.theme, ui.errorMsg)
		label.Color = color.NRGBA{R: 255, A: 255}
		return layout.Center.Layout(gtx, label.Layout)
	}
	return layout.Dimensions{}
}

func (ui *LotteryUI) drawResult(gtx layout.Context) layout.Dimensions {
	if ui.currentResult == "" {
		return layout.Dimensions{}
	}

	// 使用带边框的矩形作为卡片效果
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// 绘制背景矩形
			paint.FillShape(gtx.Ops, color.NRGBA{R: 240, G: 240, B: 240, A: 255},
				clip.RRect{
					Rect: image.Rectangle{Min: gtx.Constraints.Min, Max: gtx.Constraints.Max},
					NW:   8, NE: 8, SW: 8, SE: 8,
				}.Op(gtx.Ops))

			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.H6(ui.theme, ui.currentResult)
				label.Color = color.NRGBA{R: 0, G: 150, B: 0, A: 255}
				return label.Layout(gtx)
			})
		})
	})
}

func (ui *LotteryUI) generateRandom() {
	result := logic.GenerateRandomTicket()
	ui.currentResult = logic.FormatNumbers(result.FrontNumbers, result.BackNumbers)
	ui.errorMsg = ""
}

func (ui *LotteryUI) generateDanTuo() {
	// 解析前区胆码
	frontDare := ui.parseNumbers([]*widget.Editor{
		&ui.frontDare1, &ui.frontDare2, &ui.frontDare3, &ui.frontDare4,
	})

	// 解析前区拖码
	frontDrag := ui.parseNumbers([]*widget.Editor{
		&ui.frontDrag1, &ui.frontDrag2, &ui.frontDrag3, &ui.frontDrag4,
		&ui.frontDrag5, &ui.frontDrag6, &ui.frontDrag7, &ui.frontDrag8,
	})

	// 解析后区胆码
	backDare := ui.parseNumbersBack([]*widget.Editor{&ui.backDare1})

	// 解析后区拖码
	backDrag := ui.parseNumbersBack([]*widget.Editor{
		&ui.backDrag1, &ui.backDrag2, &ui.backDrag3,
	})

	// 验证数据
	if len(frontDare) == 0 {
		ui.errorMsg = "请至少输入一个前区胆码"
		return
	}
	if len(frontDare) > 4 {
		ui.errorMsg = "前区胆码最多4个"
		return
	}
	if len(backDare) != 1 {
		ui.errorMsg = "后区胆码必须为1个"
		return
	}
	if len(frontDrag) < (5 - len(frontDare)) {
		ui.errorMsg = "前区拖码数量不足，至少需要" + strconv.Itoa(5-len(frontDare)) + "个"
		return
	}
	if len(backDrag) < (2 - len(backDare)) {
		ui.errorMsg = "后区拖码数量不足，至少需要" + strconv.Itoa(2-len(backDare)) + "个"
		return
	}

	// 检查号码是否重复
	if ui.hasDuplicate(frontDare) || ui.hasDuplicate(frontDrag) || ui.hasDuplicate(backDare) || ui.hasDuplicate(backDrag) {
		ui.errorMsg = "胆码或拖码中存在重复数字"
		return
	}

	// 检查胆码和拖码是否重复
	for _, dare := range frontDare {
		for _, drag := range frontDrag {
			if dare == drag {
				ui.errorMsg = "前区胆码和拖码不能重复"
				return
			}
		}
	}

	for _, dare := range backDare {
		for _, drag := range backDrag {
			if dare == drag {
				ui.errorMsg = "后区胆码和拖码不能重复"
				return
			}
		}
	}

	results := logic.GenerateDanTuoTicket(frontDare, frontDrag, backDare, backDrag)

	if len(results) == 0 {
		ui.errorMsg = "生成失败，请检查输入"
		return
	}

	// 显示第一注
	ui.currentResult = logic.FormatNumbers(results[0].FrontNumbers, results[0].BackNumbers)
	if len(results) > 1 {
		ui.currentResult += "\n共生成" + strconv.Itoa(len(results)) + "注"
	}
	ui.errorMsg = ""
}

func (ui *LotteryUI) parseNumbers(editors []*widget.Editor) []int {
	var numbers []int
	for _, editor := range editors {
		text := editor.Text()
		if text == "" {
			continue
		}
		if num, err := strconv.Atoi(text); err == nil {
			if num >= 1 && num <= 35 {
				numbers = append(numbers, num)
			}
		}
	}
	return numbers
}

func (ui *LotteryUI) parseNumbersBack(editors []*widget.Editor) []int {
	var numbers []int
	for _, editor := range editors {
		text := editor.Text()
		if text == "" {
			continue
		}
		if num, err := strconv.Atoi(text); err == nil {
			if num >= 1 && num <= 12 {
				numbers = append(numbers, num)
			}
		}
	}
	return numbers
}

func (ui *LotteryUI) hasDuplicate(numbers []int) bool {
	seen := make(map[int]bool)
	for _, num := range numbers {
		if seen[num] {
			return true
		}
		seen[num] = true
	}
	return false
}
