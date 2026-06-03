package ui

import (
	"embed"
	"image"
	"image/color"
	"log"
	"strconv"

	"lottery.app/logic"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

//go:embed assets/*.otf
var fontsFS embed.FS

type LotteryUI struct {
	theme     *material.Theme
	randomBtn widget.Clickable
	dantuoBtn widget.Clickable

	// 胆拖输入控件
	frontDareCount    widget.Editor // 前区胆码个数
	frontDragCount    widget.Editor // 前区拖码个数
	backDareCount     widget.Editor // 后区胆码个数
	backDragCount     widget.Editor // 后区拖码个数
	generateDanTuoBtn widget.Clickable

	currentResult string
	showDanTuoUI  bool
	errorMsg      string
	isDanTuoMode  bool // 标记当前显示的是否为胆拖模式

	list layout.List
}

func NewLotteryUI() *LotteryUI {
	ui := &LotteryUI{
		showDanTuoUI: false,
		isDanTuoMode: false,
		list: layout.List{
			Axis: layout.Vertical,
		},
	}
	ui.initEditors()
	return ui
}

func (ui *LotteryUI) initEditors() {
	// 初始化编辑器
	ui.frontDareCount.SingleLine = true
	ui.frontDragCount.SingleLine = true
	ui.backDareCount.SingleLine = true
	ui.backDragCount.SingleLine = true
}

func (ui *LotteryUI) Run(window *app.Window) error {
	// 方法：从嵌入的文件系统中读取字体
	// 先尝试读取可能的字体文件
	var fontData []byte
	var err error

	// 尝试常见的字体文件名（根据你实际放置的文件修改）
	candidates := []string{
		"assets/SourceHanSans-Regular.otf",
		"assets/SourceHanSans-Bold.otf",
		"assets/SourceHanSans-Light.otf",
		"assets/SourceHanSans-Medium.otf",
		"assets/SourceHanSans-Heavy.otf",
		"assets/SourceHanSans-ExtraLight.otf",
		"assets/SourceHanSans-Normal.otf",
	}

	for _, candidate := range candidates {
		fontData, err = fontsFS.ReadFile(candidate)
		if err == nil {
			log.Printf("成功加载字体: %s", candidate)
			break
		}
	}
	if err != nil {
		log.Printf("未找到字体文件，使用默认字体: %v", err)
		fontData = nil // 让 text.NewCollection 使用默认字体
		return err
	}
	// 解析字体
	face, err := opentype.Parse(fontData)
	if err != nil {
		log.Printf("解析字体失败，使用默认字体: %v", err)
		fontData = nil // 让 text.NewCollection 使用默认字体
		return err
	}
	// 创建字体集合并配置主题
	fontCollection := []text.FontFace{{Face: face}}
	ui.theme = material.NewTheme()
	// 为这个主题显式地配置 Go 字体
	// text.WithCollection 用于指定字体来源，gofont.Collection() 提供了内嵌的 Go 字体
	ui.theme.Shaper = text.NewShaper(text.WithCollection(append(gofont.Collection(), fontCollection...)))
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
		label := material.H6(ui.theme, "lottery.app")
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
			layout.Rigid(ui.drawRuleHint),
		)
	})
}

func (ui *LotteryUI) drawDanTuoTitle(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.H5(ui.theme, "胆拖随机选号")
		return label.Layout(gtx)
	})
}

func (ui *LotteryUI) drawFrontArea(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(ui.theme, "前区（1-35）")
			return label.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return ui.editorWrapper(&ui.frontDareCount, "胆码个数(1-4)")(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Caption(ui.theme, "胆码个数")
							return label.Layout(gtx)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return ui.editorWrapper(&ui.frontDragCount, "拖码个数")(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Caption(ui.theme, "拖码个数")
							return label.Layout(gtx)
						}),
					)
				}),
			)
		}),
	)
}

func (ui *LotteryUI) drawBackArea(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(ui.theme, "后区（1-12）")
			return label.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return ui.editorWrapper(&ui.backDareCount, "胆码个数(1)")(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Caption(ui.theme, "胆码个数")
							return label.Layout(gtx)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return ui.editorWrapper(&ui.backDragCount, "拖码个数")(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Caption(ui.theme, "拖码个数")
							return label.Layout(gtx)
						}),
					)
				}),
			)
		}),
	)
}

func (ui *LotteryUI) editorWrapper(editor *widget.Editor, hint string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// 设置最小宽度和高度
		gtx.Constraints.Min.X = gtx.Dp(unit.Dp(80))
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
	btn := material.Button(ui.theme, &ui.generateDanTuoBtn, "生成胆拖号码")
	if ui.generateDanTuoBtn.Clicked(gtx) {
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

func (ui *LotteryUI) drawRuleHint(gtx layout.Context) layout.Dimensions {
	// 显示规则提示
	if ui.showDanTuoUI {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			label := material.Caption(ui.theme, "提示：前区胆码1-4个，拖码个数需保证总数≥5；后区胆码必须1个，拖码个数需保证总数≥2")
			label.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
			return label.Layout(gtx)
		})
	}
	return layout.Dimensions{}
}

func (ui *LotteryUI) drawResult(gtx layout.Context) layout.Dimensions {
	if ui.currentResult == "" {
		return layout.Dimensions{}
	}

	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// 绘制背景矩形
			paint.FillShape(gtx.Ops, color.NRGBA{R: 240, G: 240, B: 240, A: 255},
				clip.RRect{
					Rect: image.Rectangle{Min: gtx.Constraints.Min, Max: gtx.Constraints.Max},
					NW:   8, NE: 8, SW: 8, SE: 8,
				}.Op(gtx.Ops))

			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// 使用带格式的文本显示，支持多行
				label := material.Label(ui.theme, unit.Sp(14), ui.currentResult)
				label.Color = color.NRGBA{R: 0, G: 150, B: 0, A: 255}
				return label.Layout(gtx)
			})
		})
	})
}

func (ui *LotteryUI) generateRandom() {
	result := logic.GenerateRandomTicket()
	ui.currentResult = logic.FormatNumbers(result.FrontNumbers, result.BackNumbers)
	ui.isDanTuoMode = false
	ui.errorMsg = ""
}

func (ui *LotteryUI) generateDanTuo() {
	// 解析输入
	frontDareCount, err1 := strconv.Atoi(ui.frontDareCount.Text())
	frontDragCount, err2 := strconv.Atoi(ui.frontDragCount.Text())
	backDareCount, err3 := strconv.Atoi(ui.backDareCount.Text())
	backDragCount, err4 := strconv.Atoi(ui.backDragCount.Text())

	// 验证输入
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		ui.errorMsg = "请输入有效的数字"
		return
	}

	// 验证范围
	if frontDareCount < 1 || frontDareCount > 4 {
		ui.errorMsg = "前区胆码个数必须在1-4之间"
		return
	}
	if backDareCount != 1 {
		ui.errorMsg = "后区胆码个数必须为1"
		return
	}
	if frontDareCount+frontDragCount < 5 {
		ui.errorMsg = "前区胆码+拖码总数不能少于5个"
		return
	}
	if frontDragCount < 0 {
		ui.errorMsg = "前区拖码个数不能为负数"
		return
	}
	if backDareCount+backDragCount < 2 {
		ui.errorMsg = "后区胆码+拖码总数不能少于2个"
		return
	}
	if backDragCount < 0 {
		ui.errorMsg = "后区拖码个数不能为负数"
		return
	}

	// 限制最大输入
	maxFrontDrag := 35 - frontDareCount
	if frontDragCount > maxFrontDrag {
		ui.errorMsg = "前区拖码个数不能超过" + strconv.Itoa(maxFrontDrag)
		return
	}

	maxBackDrag := 12 - backDareCount
	if backDragCount > maxBackDrag {
		ui.errorMsg = "后区拖码个数不能超过" + strconv.Itoa(maxBackDrag)
		return
	}

	// 生成胆拖号码
	result, err := logic.GenerateRandomDanTuo(frontDareCount, frontDragCount, backDareCount, backDragCount)
	if err != nil {
		ui.errorMsg = "生成失败，请检查输入是否符合规则"
		return
	}

	// 使用专门的胆拖格式化函数显示结果（带金额）
	ui.currentResult = logic.FormatDanTuoWithPrice(result)
	ui.isDanTuoMode = true
	ui.errorMsg = ""
}
