package controller

import (
	"fmt"
	"os"

	"github.com/armosec/armotoy/common"
	"github.com/armosec/armotoy/model"
	"github.com/armosec/armotoy/view"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitController(path, source, version string) (*Controller, error) {
	controller := &Controller{StateMap: make(map[string]common.State),
		Selections: make(map[string][]Selection),
		StateFilters: map[string]common.Filters{
			common.CONTEXT_FRAMEWORK: {Equals: make(map[string][]string)},
			common.CONTEXT_CONTROLS:  {Equals: make(map[string][]string)},
			common.CONTEXT_RESOURCE:  {Equals: make(map[string][]string)}}}

	data, err := model.PostureModelInit(source, path, version)
	if err != nil {
		return nil, err
	}
	controller.model = data
	controller.app = tview.NewApplication()
	frameworkState, _, err := controller.CreateFrameworkPage(controller.model)
	if err != nil {
		return nil, err
	}

	controller.root = tview.NewGrid().
		SetRows(1, 0, 3).
		SetColumns(0).
		SetBorders(true)

	controller.AddState(frameworkState)
	controller.SetState(frameworkState.Name)

	controlsState, _, err := controller.CreateControlPage([]string{})
	if err != nil {
		return nil, err
	}

	controller.AddState(controlsState)
	resourceState, _, err := controller.CreateResourcePage([]string{}, []string{})
	if err != nil {
		return nil, err
	}

	controller.AddState(resourceState)
	controller.setupInputs()
	return controller, nil
}

func (controller *Controller) ResetFilters(section string) {
	filters := controller.StateFilters[common.CONTEXT_FRAMEWORK]

	for k := range filters.Equals {
		delete(filters.Equals, k)
	}

	controller.StateFilters[common.CONTEXT_FRAMEWORK] = filters

}
func (controller *Controller) setupInputs() {
	controller.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyCtrlF:

			controller.SetState(common.CONTEXT_FRAMEWORK)
			return nil

		case tcell.KeyCtrlC:
			controller.SetState(common.CONTEXT_CONTROLS)
			return nil

		case tcell.KeyCtrlR:

			controller.SetState(common.CONTEXT_RESOURCE)
			return nil

		case tcell.KeyCtrlX:
			controller.app.Stop()
			fmt.Printf("bye\n")
			os.Exit(0)

		case tcell.KeyF2:
			controller.app.SetFocus(controller.StateMap[controller.CurrentState].Content)
			return nil

		}

		return event
	})
}
func (controller *Controller) AddState(state *common.State) {
	state.Menu = view.MakeMenu(state.Name)
	controller.StateMap[state.Name] = *state
}
func (controller *Controller) SetState(stateName string) error {
	if state, ok := controller.StateMap[stateName]; ok {
		controller.CurrentState = stateName

		controller.root.Clear()
		controller.root.
			AddItem(state.Menu, 0, 0, 1, 3, 0, 0, false).
			AddItem(state.Footer, 2, 0, 1, 3, 0, 0, true).
			AddItem(state.Content, 1, 0, 1, 3, 0, 0, false)
		controller.app.SetFocus(state.Content)
		return nil
	}

	return fmt.Errorf("invalid state")
}

func (controller *Controller) Start() {
	controller.app.SetFocus(controller.StateMap[controller.CurrentState].Content)
	if err := controller.app.SetRoot(controller.root, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
