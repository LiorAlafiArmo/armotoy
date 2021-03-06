package controller

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/armosec/armotoy/common"
	"github.com/armosec/armotoy/model"
	"github.com/armosec/armotoy/view"
	v1 "github.com/armosec/opa-utils/reporthandling/helpers/v1"
	"github.com/armosec/opa-utils/reporthandling/results/v1/reportsummary"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v2"
	"k8s.io/utils/strings/slices"
)

const (
	FRAMEWORK_TABLE_FOOTER_TEXT = "Ctrl+D - Dig into\tCtrl+B - Broadcast selection\tCtrl+E - Create an Exception from selection\tCtrl+S - Apply Filters\t]-switch focus to filters\tEsc-clear selections\tCtrl+A - show all(reset filters as well)"
	CONTROL_TABLE_FOOTER_TEXT   = "Ctrl+D - Dig into\tCtrl+K - Control Info.\tCtrl+B - Broadcast selection\tCtrl+E - Create an Exception from selection\tCtrl+S - Apply Filters\t]-switch focus to filters\tEsc-clear selections\tCtrl+A - show all(reset filters as well)"

	RESOURCE_TABLE_FOOTER_TEXT = "Ctrl+D - Dig into\tCtrl+B - Broadcast selection\tCtrl+E - Create an Exception from selection\tCtrl+S - Apply Filters\t]-switch focus to filters\tEsc-clear selections\tCtrl+A - show all(reset filters as well)"

	MAX_RULE_DISPLAY_SIZE = 550
)

func (controller *Controller) CreateFrameworkPage(data *model.PostureModel) (*common.State, *tview.Table, error) {
	state := &common.State{
		Name:  common.CONTEXT_FRAMEWORK,
		Index: 0,
	}

	frameworks, err := data.GetFrameworksTable()
	if err != nil {
		return nil, nil, err
	}

	selections := tview.NewTextView()
	filters := controller.StateFilters[common.CONTEXT_FRAMEWORK]
	view.SetFiltersText(selections, &filters)

	defaultFrameworkColumns := map[string]common.ColumnAttributes{"Name": *common.DefaultColumnAttributes(0), "Status": *view.StatusColumnAttributes(1), "Status info": *common.DefaultColumnAttributes(2), "Score": *common.DefaultColumnAttributes(3)}
	if len(filters.Equals) > 0 {
		frameworks = frameworks.FilterByColumns(filters.Equals)
	}
	controller.Data[common.CONTEXT_FRAMEWORK] = view.Data{ColumnAttributes: defaultFrameworkColumns, Model: frameworks}

	table := view.CreateTable(frameworks, defaultFrameworkColumns)

	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			view.ClearSelection(table)
			controller.ClearSelection(common.CONTEXT_FRAMEWORK)
		}
		if key == tcell.KeyEnter {

			table.SetSelectable(true, false)

		}

		// if key == tcell.KeyTab {

		// }
	}).SetSelectedFunc(func(row int, column int) {
		if row == 0 {
			return
		}
		frameworkName := table.GetCell(row, 0).Text

		if frameworkName == "" {
			return
		}

		controller.basicSelectionhandler(table, row, common.CONTEXT_FRAMEWORK, frameworkName, table.GetCell(row, 0).GetReference())
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyCtrlE:
			ex, err := controller.CreateExceptionPage()
			if err == nil {
				controller.AddState(ex)
				controller.SetState(ex.Name)
			}
			return nil
		case tcell.KeyCtrlB:
			controller.CreateBroadcastersPage()
			controller.SetState(common.CONTEXT_INTEGRATION_MESSAGE_SETTINGS)
			return nil
		case tcell.KeyCtrlD:
			row, col := table.GetSelection()
			if row == 0 {
				return nil
			}
			frameworkName := table.GetCell(row, col).Text

			newstate, _, err := controller.CreateControlPage([]string{frameworkName})
			if err == nil {
				controller.AddState(newstate)
				controller.SetState(newstate.Name)

			}

			return nil
		case tcell.KeyCtrlA:
			controller.ResetFilters(common.CONTEXT_FRAMEWORK)
			view.ClearSelection(table)
			controller.ClearSelection(common.CONTEXT_FRAMEWORK)

			astate, _, err := controller.CreateFrameworkPage(controller.model)
			if err == nil {
				controller.AddState(astate)
				controller.SetState(astate.Name)
			}

			selections.SetText("")
			return nil

		}

		return event
	})
	footer := tview.NewTextView().SetDynamicColors(true)
	footer.SetTextColor(tcell.ColorLightGoldenrodYellow).SetText(RESOURCE_TABLE_FOOTER_TEXT)

	filterForm := view.CreateTableFilterForm(defaultFrameworkColumns, &filters, selections)

	flex := view.CreateTableLayout(filterForm, selections, table, func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case ']':
			if table.HasFocus() {
				table.Blur()
				controller.app.SetFocus(filterForm)
				footer.SetText("Tab-jump between fields\t]-switch focus to table")
			} else {
				controller.app.SetFocus(table)
				footer.SetText(RESOURCE_TABLE_FOOTER_TEXT)
			}
			return nil
		}
		switch event.Key() {
		case tcell.KeyCtrlS:

			state, _, err := controller.CreateFrameworkPage(data)
			if err == nil {
				controller.AddState(state)
				controller.SetState(state.Name)
			}
		}
		return event
	})

	state.Content = flex

	state.Footer = footer
	return state, table, nil
}

func (controller *Controller) CreateControlPage(frameworks []string) (*common.State, *tview.Table, error) {
	state := &common.State{
		Name:  common.CONTEXT_CONTROLS,
		Index: 1,
	}

	controls, err := controller.model.GetControlsTable(frameworks)
	if err != nil {
		return nil, nil, err
	}

	selections := tview.NewTextView()
	filters := controller.StateFilters[common.CONTEXT_FRAMEWORK]
	view.SetFiltersText(selections, &filters)
	// {Label: "ID"},
	// 	{Label: "Name"},
	// 	{Label: "Status"},
	// 	{Label: "Status Info"},
	// 	{Label: "Score", Type: "number"},
	// 	{Label: "Severity"},
	// 	{Label: "Host Scan"},
	// 	{Label: "Customized Config"},
	// 	{Label: "Cloud Related"},

	defaultControlColumns := map[string]common.ColumnAttributes{"Name": *common.DefaultColumnAttributes(0), "Status": *view.StatusColumnAttributes(1), "Severity": *view.SeverityColumnAttributes(2), "Status info": *common.DefaultColumnAttributes(3), "Score": *common.DefaultColumnAttributes(4), "Host Scan": *common.DefaultColumnAttributes(5), "Customized Config": *common.DefaultColumnAttributes(6), "Cloud Related": *common.DefaultColumnAttributes(7), "ID": {Hidden: true, Color: nil}}

	if len(filters.Equals) > 0 {
		controls = controls.FilterByColumns(filters.Equals)
	}
	controls.SortByStatus("Status")
	controller.Data[common.CONTEXT_CONTROLS] = view.Data{ColumnAttributes: defaultControlColumns, Model: controls}

	table := view.CreateTable(controls, defaultControlColumns)

	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			view.ClearSelection(table)
			controller.ClearSelection(common.CONTEXT_CONTROLS)
		}
		if key == tcell.KeyEnter {

			table.SetSelectable(true, false)

		}

		// if key == tcell.KeyTab {

		// }
	}).SetSelectedFunc(func(row int, column int) {
		if row == 0 {
			return
		}
		control, ok := table.GetCell(row, 0).GetReference().(*reportsummary.ControlSummary)

		if !ok {
			return
		}
		controller.basicSelectionhandler(table, row, common.CONTEXT_CONTROLS, control.ControlID, control)
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyCtrlE:
			ex, err := controller.CreateExceptionPage()
			if err == nil {
				controller.AddState(ex)
				controller.SetState(ex.Name)
			}
			return nil
		case tcell.KeyCtrlB:
			controller.CreateBroadcastersPage()
			controller.SetState(common.CONTEXT_INTEGRATION_MESSAGE_SETTINGS)

			return nil
		case tcell.KeyCtrlD:

			row, _ := table.GetSelection()
			if row == 0 {
				return nil
			}
			control, ok := table.GetCell(row, 0).GetReference().(*reportsummary.ControlSummary)
			if !ok {
				return nil
			}
			controller.StateFilters[common.CONTEXT_RESOURCE].Equals["Status"] = []string{"failed", "excluded"}
			newState, _, err := controller.CreateResourcePage(frameworks, []string{control.ControlID})
			if err == nil {
				controller.AddState(newState)
				controller.SetState(newState.Name)

			}

			return nil
		case tcell.KeyCtrlK:
			{
				row, _ := table.GetSelection()
				if row == 0 {
					return nil
				}
				control, ok := table.GetCell(row, 0).GetReference().(*reportsummary.ControlSummary)
				if !ok {
					return nil
				}
				if newState, err := controller.CreateCtrlInfoPage(control.ControlID); err == nil {
					controller.AddState(newState)
					controller.SetState(newState.Name)
				}

				return nil

			}
		case tcell.KeyCtrlA:
			controller.ResetFilters(common.CONTEXT_CONTROLS)
			view.ClearSelection(table)
			controller.ClearSelection(common.CONTEXT_CONTROLS)

			astate, _, err := controller.CreateControlPage([]string{})
			if err == nil {
				controller.AddState(astate)
				controller.SetState(astate.Name)
			}

			return nil

		}

		return event

	})
	footer := tview.NewTextView().SetDynamicColors(true)
	footer.SetText(CONTROL_TABLE_FOOTER_TEXT)
	filterForm := view.CreateTableFilterForm(defaultControlColumns, &filters, selections)
	flex := view.CreateTableLayout(filterForm, selections, table, func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case ']':
			if table.HasFocus() {
				table.Blur()
				controller.app.SetFocus(filterForm)
				footer.SetText("Tab-jump between fields\t]-switch focus to table")
			} else {
				controller.app.SetFocus(table)
				footer.SetText(CONTROL_TABLE_FOOTER_TEXT)
			}
			return nil
		}
		switch event.Key() {
		case tcell.KeyCtrlS:

			state, _, err := controller.CreateControlPage(frameworks)
			if err == nil {
				controller.AddState(state)
				controller.SetState(state.Name)
			}
		}
		return event
	})

	// defaultColumnSettings := common.DefaultColumnAttributes()
	// statusCol := StatusColumnAttributes()
	state.Content = flex
	// footer.SetTextColor(tcell.ColorLightGoldenrodYellow).SetText("Tab- Dig into(must be selected)\tCtrl+B - Broadcast selection\tCtrl+E - Create an Exception")
	state.Footer = footer
	return state, table, nil
}

func (controller *Controller) CreateCtrlInfoPage(ctrlID string) (*common.State, error) {
	state := &common.State{
		Name:   common.CONTEXT_CONTROL_INFO,
		Index:  1,
		Menu:   view.MakeMenu(""),
		Footer: tview.NewBox(),
	}
	if git := controller.model.GetGitRegoStore(); git != nil {
		ctrl, err := git.GetOPAControlByID(ctrlID)
		if err != nil {
			return nil, err
		}
		txt := tview.NewTextView().SetWordWrap(true).SetWrap(true).SetDynamicColors(true)
		mytxt := fmt.Sprintf("[yellow]Control ID:[white] %s\n[yellow]Name: [white]%s\n[yellow]Description: [white]%s\n[yellow]Rules:[white]\n", ctrlID, ctrl.Name, ctrl.Description)
		for i := range ctrl.Rules {
			rule := ctrl.Rules[i].Rule
			if len(rule) > MAX_RULE_DISPLAY_SIZE {
				rule = rule[:MAX_RULE_DISPLAY_SIZE] + "..."
			}
			mytxt = fmt.Sprintf("%s[yellow]%s[white]\n---------\n%s\n", mytxt, ctrl.Rules[i].Name, rule)

		}
		txt.SetText(mytxt)
		state.Content = txt
		return state, nil
	}

	return nil, fmt.Errorf("missing control information")
}

func (controller *Controller) CreateResourcePage(frameworks []string, controls []string) (*common.State, *tview.Table, error) {
	state := &common.State{
		Name:  common.CONTEXT_RESOURCE,
		Index: 2,
	}

	resources, err := controller.model.GetResourcesTable(frameworks, controls)
	if err != nil {
		return nil, nil, err
	}

	resources.SortByColumns([]string{"Namespace", "Kind", "Name", "Status"})

	defaultResourceColumns := map[string]common.ColumnAttributes{"Namespace": *common.DefaultColumnAttributes(0), "Kind": *common.DefaultColumnAttributes(1), "Name": *common.DefaultColumnAttributes(2), "Status": *view.StatusColumnAttributes(3), "k8s Object": *common.DefaultColumnAttributes(4), "Failed Controls": *common.DefaultColumnAttributes(5)}

	selections := tview.NewTextView()
	filters := controller.StateFilters[common.CONTEXT_RESOURCE]
	view.SetFiltersText(selections, &filters)

	filterForm := view.CreateTableFilterForm(defaultResourceColumns, &filters, selections)

	if len(filters.Equals) > 0 {
		resources = resources.FilterByColumns(filters.Equals)
	}

	controller.Data[common.CONTEXT_RESOURCE] = view.Data{ColumnAttributes: defaultResourceColumns, Model: resources}

	table := view.CreateTable(resources, defaultResourceColumns)
	footer := tview.NewTextView().SetDynamicColors(true)
	footer.SetTextColor(tcell.ColorLightGoldenrodYellow).SetText("Ctrl+Y - inspect yaml\tCtrl+B - Broadcast selection\tCtrl+E - Create an Exception")

	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			view.ClearSelection(table)
			controller.ClearSelection(common.CONTEXT_RESOURCE)
			controller.ResetFilters(common.CONTEXT_RESOURCE)
		}
		if key == tcell.KeyEnter {

			table.SetSelectable(true, false)

		}
	}).SetSelectedFunc(func(row int, column int) {
		if row == 0 {
			return
		}
		resource, ok := table.GetCell(row, 0).GetReference().(*model.ResourceReference)

		if !ok {
			return
		}
		controller.basicSelectionhandler(table, row, common.CONTEXT_RESOURCE, resource.Raw.ResourceID, resource)
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyCtrlE:
			ex, err := controller.CreateExceptionPage()
			if err == nil {
				controller.AddState(ex)
				controller.SetState(ex.Name)
			}
			return nil
		case tcell.KeyCtrlB:
			controller.CreateBroadcastersPage()
			controller.SetState(common.CONTEXT_INTEGRATION_MESSAGE_SETTINGS)

			return nil

		case tcell.KeyCtrlY:
			row, _ := table.GetSelection()
			if row == 0 {
				return nil
			}
			resource, ok := table.GetCell(row, 0).GetReference().(*model.ResourceReference)

			if !ok {
				return nil
			}

			controller.YAMLInspect(resource, controls)

			return nil
		case tcell.KeyCtrlA:
			controller.ResetFilters(common.CONTEXT_RESOURCE)
			view.ClearSelection(table)
			controller.ClearSelection(common.CONTEXT_RESOURCE)

			if newstate, _, err := controller.CreateResourcePage([]string{}, []string{}); err == nil {
				controller.AddState(newstate)
				controller.SetState(newstate.Name)
			}

			selections.SetText("")
			return nil

		}

		return event

	})

	flex := view.CreateTableLayout(filterForm, selections, table, func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case ']':
			if table.HasFocus() {
				table.Blur()
				controller.app.SetFocus(filterForm)
				footer.SetText("Tab-jump between fields\t]-switch focus to table")
			} else {
				controller.app.SetFocus(table)
				footer.SetText(FRAMEWORK_TABLE_FOOTER_TEXT)
			}
			return nil
		}
		switch event.Key() {
		case tcell.KeyCtrlS:

			state, _, err := controller.CreateResourcePage(frameworks, controls)
			if err == nil {
				controller.AddState(state)
				controller.SetState(state.Name)
			}
		}
		return event
	})

	// defaultColumnSettings := common.DefaultColumnAttributes()
	// statusCol := StatusColumnAttributes()
	state.Content = flex

	state.Footer = footer
	return state, table, nil
}

func specialSplit(str string) []string {
	resp := make([]string, 0)

	for _, s := range strings.Split(str, ".") {
		if !strings.Contains(s, "[") {
			resp = append(resp, s)
		} else {
			pos := strings.Index(s, "[")
			if pos > -1 {
				resp = append(resp, s[:pos])
				resp = append(resp, s[pos:])
			}
		}
	}
	return resp
}

func (controller *Controller) YAMLInspect(resource *model.ResourceReference, controls []string) error {
	newstate := &common.State{Name: common.CONTEXT_YAML_INSPECT, Index: 3}
	newstate.Footer = controller.StateMap[controller.CurrentState].Footer

	b, er := json.Marshal(resource.Raw.GetObject())
	if er != nil {
		return er
	}
	// original, err := yaml.Marshal(modified)
	var leaf interface{}
	er = json.Unmarshal(b, &leaf)
	if er != nil {
		return er
	}
	for i := range resource.Result.AssociatedControls {
		selectedFrameworks, _ := controller.GetSelectionValues(common.CONTEXT_FRAMEWORK)
		if len(controls) > 0 && !slices.Contains(controls, resource.Result.AssociatedControls[i].ControlID) {
			continue
		}
		if resource.Result.AssociatedControls[i].GetStatus(&v1.Filters{FrameworkNames: selectedFrameworks}).IsFailed() {
			for _, rule := range resource.Result.AssociatedControls[i].ListRules() {
				for _, path := range rule.Paths {
					splits := specialSplit(path.FailedPath)

					_, _, leaf = model.UpdateLeaf(nil, leaf, splits, "[red]", "", "[white]")

					splits = specialSplit(path.FixPath.Path)
					_, _, leaf = model.UpdateLeaf(nil, leaf, splits, "[green]", path.FixPath.Value, "[white]")
				}

			}
		}
	}

	yamlbytes, err := yaml.Marshal(leaf)
	if err != nil {
		return err
	}

	newstate.Content = tview.NewTextView().SetTextColor(tcell.ColorWhite).SetText(string(yamlbytes)).SetDynamicColors(true)
	if err == nil {
		controller.AddState(newstate)
		controller.SetState(newstate.Name)

	}
	return nil
}

func (controller *Controller) basicSelectionhandler(table *tview.Table, row int, selectionKey, value string, data interface{}) {
	if pos := controller.IndexOf(selectionKey, value); pos > -1 {
		controller.RemoveSelection(selectionKey, value)
		table.GetCell(row, 0).SetTextColor(tcell.ColorWhite)
	} else {
		table.GetCell(row, 0).SetTextColor(tcell.ColorRed)
		controller.AddSelection(selectionKey, value, data)
	}

	table.SetSelectable(false, false)

}
