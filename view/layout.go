package view

import (
	"github.com/armosec/armoapi-go/armotypes"
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

func CreateExceptionLayout(selections string, except *armotypes.PostureExceptionPolicy) (*tview.Flex, *tview.InputField, *tview.TextView) {
	selected := tview.NewTextView().SetText(selections)
	fileinput := tview.NewInputField().SetLabel("save exception json as (specify path)").SetFieldWidth(40)
	review := tview.NewTextView()

	flex := tview.NewFlex().
		AddItem(fileinput, 10, 3, true).
		AddItem(selected, 0, 12, false).
		AddItem(review, 0, 11, false)
	flex.SetDirection(tview.FlexRow)
	return flex, fileinput, review
}
