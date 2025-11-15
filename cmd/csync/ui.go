package main

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func Success(content string) {
	pterm.Println()
	contentStyle := pterm.NewStyle(pterm.Bold)
	iconStyle := pterm.NewStyle(pterm.FgLightGreen, pterm.Bold)
	contentStyle.Println(iconStyle.Sprint("> ") + content + "  ")
	pterm.Println()
}

func Info(content string) {
	pterm.Println()
	contentStyle := pterm.NewStyle(pterm.Bold)
	iconStyle := pterm.NewStyle(pterm.FgLightBlue, pterm.Bold)
	contentStyle.Println(iconStyle.Sprint("> ") + content + "  ")
	pterm.Println()
}

func Warning(content string) {
	pterm.Println()
	contentStyle := pterm.NewStyle(pterm.Bold)
	iconStyle := pterm.NewStyle(pterm.FgLightYellow, pterm.Bold)
	contentStyle.Println(iconStyle.Sprint("> ") + content + "  ")
	pterm.Println()
}

func Fail(content string) {
	pterm.Println()
	contentStyle := pterm.NewStyle(pterm.Bold)
	iconStyle := pterm.NewStyle(pterm.FgLightRed, pterm.Bold)
	contentStyle.Println(iconStyle.Sprint("> ") + content + "  ")
	pterm.Println()
}

func Spinner(labels []string, showTimer bool) func() {
	multi, _ := pterm.DefaultMultiPrinter.Start()

	successPrinter := pterm.PrefixPrinter{
		Prefix: pterm.Prefix{
			Text:  "✓",
			Style: pterm.NewStyle(pterm.FgGreen),
		},
	}

	spinners := make([]*pterm.SpinnerPrinter, 0, len(labels))

	for _, label := range labels {
		spinner, _ := pterm.DefaultSpinner.
			WithSequence(" ⣾ ", " ⣽ ", " ⣻ ", " ⢿ ", " ⡿ ", " ⣟ ", " ⣯ ", " ⣷ ").
			WithStyle(pterm.NewStyle(pterm.FgCyan)).
			WithShowTimer(showTimer).
			WithWriter(multi.NewWriter()).
			Start(label)

		spinner.SuccessPrinter = &successPrinter

		spinners = append(spinners, spinner)
	}

	return func() {
		for i, s := range spinners {
			s.Success(labels[i])
		}
		multi.Stop()
	}
}

func Text(content string, icon string) {
	if icon == "" {
		pterm.Println(content)
		return
	}
	iconStyle := pterm.NewStyle(pterm.FgLightBlue)
	pterm.Println(iconStyle.Sprint(icon) + "  " + content)
}

func BreakLine() {
	pterm.Println()
}

func Tree(rootNode string, list pterm.LeveledList) {
	root := putils.TreeFromLeveledList(list)
	root.Text = rootNode

	pterm.DefaultTree.WithRoot(root).Render()
}

func List(rootNode string, list []string) {
	style := pterm.NewStyle(pterm.Bold)
	style.Println(rootNode)
	for i, item := range list {
		fmt.Println("  " + fmt.Sprintf("%d. ", i+1) + item)
	}
}
