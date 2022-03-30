package view

import (
	"fmt"
	"strings"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armotoy/common"
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

func CreatBroadcastingLayout(selections string, broadcastOptions *common.BroadcastOptions) (*tview.Flex, *tview.Form) {
	selected := tview.NewTextView().SetText(selections)
	selected.SetDynamicColors(true)
	sendform := broadcastOptionForm(broadcastOptions)

	flex := tview.NewFlex().
		AddItem(selected, 0, 12, false).
		AddItem(sendform, 0, 7, true)
	flex.SetDirection(tview.FlexRow)

	return flex, sendform
}

func CreateIntegrations(b []common.Integration) []IntegrationView {
	integrationforms := make([]IntegrationView, 0)

	for i := range b {

		integrationView := IntegrationView{Title: tview.NewTextView(), Error: tview.NewTextView().SetTextColor(tcell.ColorRed), Info: tview.NewTextView(), Form: tview.NewForm(), Broadcaster: b[i].Broadcaster}
		integrationView.Title.SetDynamicColors(true).SetText(fmt.Sprintf("[green][[white]\t  %s \t[green]][white]", b[i].Broadcaster.GetType()))
		integrationView.Form.AddCheckbox("Active", b[i].IsActive, func(checked bool) {
			b[i].IsActive = checked
		})
		integrationView.Form.AddInputField("Target", "", 40, nil, nil)
		integrationView.Form.AddButton("Add/Remove target", func() {
			inp := integrationView.Form.GetFormItemByLabel("Target")
			targetinp, ok := inp.(*tview.InputField)
			if ok {
				pos := integrationView.Broadcaster.FindTarget(targetinp.GetText())
				var err error
				if pos == -1 {
					err = integrationView.Broadcaster.AddTarget(targetinp.GetText())
				} else {

					err = integrationView.Broadcaster.RemoveTarget(targetinp.GetText())
				}
				if err != nil {
					integrationView.Error.SetTextColor(tcell.ColorRed).SetText(err.Error())
				} else {
					integrationView.Info.SetText(strings.Join(integrationView.Broadcaster.GetTargets(), ","))
					targetinp.SetText("")
				}

			} else {
				integrationView.Error.SetText("not ok")
			}

		})

		integrationforms = append(integrationforms, integrationView)

	}

	return integrationforms
}

func broadcastOptionForm(broadcastOptions *common.BroadcastOptions) *tview.Form {
	sendform := tview.NewForm().AddDropDown("Severity", []string{"INFO", "WARNING", "HIGH", "CRITICAL"}, 0, func(option string, optionIndex int) {
		broadcastOptions.Severity = option
	}).
		AddInputField("Title", "", 40, nil, func(text string) {
			broadcastOptions.Title = text
		}).
		AddCheckbox(common.CONTEXT_FRAMEWORK, true, func(checked bool) {
			broadcastOptions.FrameworksEnabled = checked
		}).
		AddCheckbox(common.CONTEXT_CONTROLS, true, func(checked bool) {
			broadcastOptions.ControlsEnabled = checked
		}).
		AddCheckbox(common.CONTEXT_RESOURCE, true, func(checked bool) {
			broadcastOptions.ResourcesEnabled = checked
		})
	return sendform
}
