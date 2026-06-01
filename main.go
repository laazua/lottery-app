package main

import (
	"log"
	"os"

	"lottery.app/ui"

	"gioui.org/app"
	"gioui.org/unit"
)

func main() {
	window := new(app.Window)
	window.Option(
		app.Title("大乐透随机号码生成器"),
		app.Size(unit.Dp(400), unit.Dp(600)),
	)

	lotteryUI := ui.NewLotteryUI()
	if err := lotteryUI.Run(window); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
