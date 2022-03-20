package view

import (
	"github.com/armosec/armotoy/common"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func MakeMenu(selected string) tview.Primitive {
	menu := tview.NewTextView().SetDynamicColors(true)
	switch selected {
	case common.CONTEXT_FRAMEWORK:
		menu.SetTextColor(tcell.ColorLightGoldenrodYellow).SetText("([green]F[white])rameworks\t(C)ontrols\t(R)esources\tE(X)it")
	case common.CONTEXT_CONTROLS:
		menu.SetTextColor(tcell.ColorLightGoldenrodYellow).SetText("(F)rameworks\t([green]C[white])ontrols\t(R)esources\tE(X)it")
	case common.CONTEXT_RESOURCE:
		menu.SetTextColor(tcell.ColorLightGoldenrodYellow).SetText("(F)rameworks\t(C)ontrols\t([green]R[white])esources\tE(X)it")
	default:
		menu.SetTextColor(tcell.ColorLightGoldenrodYellow).SetText("(F)rameworks\t(C)ontrols\t(R)esources\tE(X)it")

	}

	return menu
}