package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateTableLayout(form tview.Primitive, selection tview.Primitive, table tview.Primitive, capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Flex {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(form, 0, 3, false)
	flex.AddItem(selection, 0, 2, false)
	flex.AddItem(table, 0, 11, true)
	if capture != nil {
		flex.SetInputCapture(capture)
	}
	return flex
}
