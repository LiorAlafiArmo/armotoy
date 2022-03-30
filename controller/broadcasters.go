package controller

import (
	"fmt"
	"sort"

	"github.com/armosec/armotoy/broadcasters"
	"github.com/armosec/armotoy/common"
	"github.com/armosec/armotoy/view"

	// postureapis "github.com/armosec/opa-utils/reporthandling/apis"
	// "github.com/armosec/opa-utils/reporthandling/results/v1/reportsummary"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (c *Controller) AddBroadcasters(broadcasterconfigs map[string]interface{}) {
	for broadcastertype := range broadcasterconfigs {
		if broadcaster, err := broadcasters.Factory(broadcastertype, broadcasterconfigs[broadcastertype]); err == nil {
			c.Broadcasters = append(c.Broadcasters, common.Integration{Broadcaster: broadcaster, IsActive: true})
		}
	}
}

func (c *Controller) HandleIntegrationSettings() {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	newState := common.State{Name: common.CONTEXT_INTEGRATION_SETTINGS, Menu: view.MakeMenu(common.CONTEXT_INTEGRATION_SETTINGS), Content: flex}
	if len(c.IntegrationVMngr.Integrations) == 0 {
		flex.AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[red]Error: Integration section was empty/invalid in config.json, add integration configurations to enable/disable them and restart ARMOTOY[white]"), 0, 5, false)
		newState.Footer = tview.NewTextView().SetText("")

	} else {
		newState.Footer = c.IntegrationVMngr.IntegrationFooter
		c.SetCurrentIntegrationSettingsView(flex, &newState)
		flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			return c.IntegrationSettings(event, flex, &newState)
		})

	}
	c.AddState(&newState)
}

func (c *Controller) IntegrationSettings(event *tcell.EventKey, flex *tview.Flex, newState *common.State) *tcell.EventKey {
	switch event.Rune() {
	case ']':

		numOfIntegrations := len(c.IntegrationVMngr.Integrations)
		c.IntegrationVMngr.CurrentIntegrationPos++
		c.IntegrationVMngr.CurrentIntegrationPos %= numOfIntegrations
		c.SetCurrentIntegrationSettingsView(flex, newState)
		c.SetState(newState.Name)

		return nil

	case '[':
		numOfIntegrations := len(c.IntegrationVMngr.Integrations)
		if numOfIntegrations >= 1 {

			if c.IntegrationVMngr.CurrentIntegrationPos > 0 {
				c.IntegrationVMngr.CurrentIntegrationPos--
				numOfIntegrations %= numOfIntegrations
			} else if c.IntegrationVMngr.CurrentIntegrationPos == -1 {
				c.IntegrationVMngr.CurrentIntegrationPos = 0
			} else {
				c.IntegrationVMngr.CurrentIntegrationPos = numOfIntegrations - 1
			}
			c.SetCurrentIntegrationSettingsView(flex, newState)
			c.SetState(newState.Name)

		}
		return nil
	}

	switch event.Key() {
	case tcell.KeyCtrlB:
		c.SetState(common.CONTEXT_INTEGRATION_MESSAGE_SETTINGS)
		return nil
	}
	return event
}

func (c *Controller) ClearIntegrationFlex(flex *tview.Flex) {
	flex.RemoveItem(c.IntegrationVMngr.Integrations[c.IntegrationVMngr.CurrentIntegrationPos].Form)
	flex.RemoveItem(c.IntegrationVMngr.Integrations[c.IntegrationVMngr.CurrentIntegrationPos].Title)
	flex.RemoveItem(c.IntegrationVMngr.Integrations[c.IntegrationVMngr.CurrentIntegrationPos].Error)
	flex.RemoveItem(c.IntegrationVMngr.Integrations[c.IntegrationVMngr.CurrentIntegrationPos].Info)
}
func (c *Controller) SetCurrentIntegrationSettingsView(flex *tview.Flex, state *common.State) {
	flex2 := tview.NewFlex().SetDirection(tview.FlexRow)

	flex2.AddItem(c.IntegrationVMngr.Integrations[c.IntegrationVMngr.CurrentIntegrationPos].Title, 0, 5, false)
	flex2.AddItem(c.IntegrationVMngr.Integrations[c.IntegrationVMngr.CurrentIntegrationPos].Form, 0, 20, true)
	flex2.AddItem(c.IntegrationVMngr.Integrations[c.IntegrationVMngr.CurrentIntegrationPos].Info, 0, 10, false)
	flex2.AddItem(c.IntegrationVMngr.Integrations[c.IntegrationVMngr.CurrentIntegrationPos].Error, 0, 10, false)

	flex2.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return c.IntegrationSettings(event, flex2, state)
	})
	state.Content = flex2
	c.AddState(state)
}

func (c *Controller) CreateBroadcastersPage() {
	_, _, _, selected := c.GetSelectedItemsString()
	txt := tview.NewTextView()

	flex, msgform := view.CreatBroadcastingLayout(selected, &c.IntegrationMessageOptions)
	flex.AddItem(txt, 0, 5, false)
	c.columnSelectionsForms = []*tview.Form{msgform}
	c.columnSelectionsForms = append(c.columnSelectionsForms, msgform)

	for k, v := range c.Data {
		frm := view.ColumnActivationForm(k, &v)
		// frms = append(frms, frm)
		c.Data[k] = v
		c.columnSelectionsForms = append(c.columnSelectionsForms, frm)
		flex.AddItem(tview.NewTextView().SetText(k), 0, 1, false).AddItem(frm, 0, 2, false)
	}
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case ']':
			c.broadcasterColumnSelectionIndex = (c.broadcasterColumnSelectionIndex + 1) % len(c.columnSelectionsForms)
			c.app.SetFocus(c.columnSelectionsForms[c.broadcasterColumnSelectionIndex])
			return nil

		case '[':
			if c.broadcasterColumnSelectionIndex > 0 {
				c.broadcasterColumnSelectionIndex--
				c.app.SetFocus(c.columnSelectionsForms[c.broadcasterColumnSelectionIndex])
			} else {
				c.broadcasterColumnSelectionIndex = len(c.columnSelectionsForms) - 1
				c.app.SetFocus(c.columnSelectionsForms[c.broadcasterColumnSelectionIndex])
			}

			return nil
		}

		switch event.Key() {
		case tcell.KeyCtrlS:
			if c.IntegrationMessageOptions.IsValid() {
				c.CreateIntegrationMessage()
				c.SendMessage()
				c.IntegrationMessageOptions.Message = ""
				c.IntegrationMessageOptions.Title = ""
			}
			return nil

		case tcell.KeyCtrlK:
			c.SetState(common.CONTEXT_INTEGRATION_SETTINGS)
			return nil

		}
		return event
	})
	footer := tview.NewTextView()
	footer.SetText("Ctrl+S - broadcast to active integrations\tCtrl+K - integration settings")
	state := common.State{
		Name:    common.CONTEXT_INTEGRATION_MESSAGE_SETTINGS,
		Menu:    view.MakeMenu(common.CONTEXT_INTEGRATION_MESSAGE_SETTINGS),
		Content: flex,
		Footer:  footer,
	}
	c.AddState(&state)
}

func (c *Controller) SendMessage() []error {
	errors := make([]error, 0)
	for i := range c.Broadcasters {
		if c.Broadcasters[i].IsActive {
			if err := c.Broadcasters[i].Broadcaster.SendMessage(c.IntegrationMessageOptions.Severity, c.IntegrationMessageOptions.Title, c.IntegrationMessageOptions.Message); err != nil {
				errors = append(errors, err)
			}
		}
	}
	return errors
}

func (c *Controller) CreateIntegrationMessage() {
	fwmessage := ""
	ctrlsmessage := ""
	resourceMsg := ""

	//TODO: change to strategy pattern
	if c.IntegrationMessageOptions.FrameworksEnabled {
		fwmessage = c.defaultMessage(common.CONTEXT_FRAMEWORK, "Name", "Frameworks")
	}

	if c.IntegrationMessageOptions.FrameworksEnabled {
		ctrlsmessage = c.defaultMessage(common.CONTEXT_CONTROLS, "ID", "Controls")
	}

	if c.IntegrationMessageOptions.ResourcesEnabled {
		resourceMsg = c.defaultMessage(common.CONTEXT_RESOURCE, "ResourceID", "Resources")

	}
	// ctrlsmessage := c.defaultControlsMessage()

	c.IntegrationMessageOptions.Message = fwmessage + ctrlsmessage + resourceMsg

}

func (c *Controller) defaultMessage(key, filterKey, atype string) string {
	msg := atype + ":\n"

	// if c.IntegrationMessageOptions.FrameworksEnabled {
	frameworks, _ := c.GetSelectionValues(common.CONTEXT_FRAMEWORK)
	data := c.Data[common.CONTEXT_FRAMEWORK].Model.FilterByColumns(map[string][]string{filterKey: frameworks})
	cols := c.GetOrderedAttributes(common.CONTEXT_FRAMEWORK)

	for row := range data.Data {
		for j := range cols {
			if col, i := data.GetColumn(cols[j].Key); i > -1 && cols[j].Selected {
				if i < len(cols)-1 {
					msg = fmt.Sprintf("%s%s: %s\t-\t", msg, col.Label, data.Data[row].Fields[i])
				} else {
					msg = fmt.Sprintf("%s%s: %s\n", msg, col.Label, data.Data[row].Fields[i])
				}
			}
		}
	}

	// }
	return msg
}

func (c *Controller) GetOrderedAttributes(key string) []common.ColumnAttributes {
	attributesordered := make([]common.ColumnAttributes, 0, len(c.Data[key].ColumnAttributes))
	for k, col := range c.Data[key].ColumnAttributes {
		col.Key = k
		c.Data[key].ColumnAttributes[k] = col
		attributesordered = append(attributesordered, c.Data[key].ColumnAttributes[k])
	}
	sort.Slice(attributesordered, func(i, j int) bool {
		return attributesordered[i].Index < attributesordered[j].Index
	})

	return attributesordered

}
