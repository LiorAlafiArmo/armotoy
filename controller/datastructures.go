package controller

import (
	"github.com/armosec/armotoy/common"
	"github.com/armosec/armotoy/model"
	"github.com/armosec/armotoy/view"
	"github.com/rivo/tview"
)

type IntegrationViewManager struct {
	CurrentIntegrationPos int
	Integrations          []view.IntegrationView
	IntegrationFooter     tview.Primitive

	MessageOptionsView   tview.Primitive
	MessageOptionsFooter tview.Primitive
}

type Controller struct {
	model        *model.PostureModel
	app          *tview.Application
	CurrentState string
	StateMap     map[string]common.State
	root         *tview.Grid
	StateFilters map[string]common.Filters
	Selections   map[string][]Selection
	Data         map[string]view.Data

	DiggingSelection []Selection

	Broadcasters                    []common.Integration
	IntegrationMessageOptions       common.BroadcastOptions
	IntegrationVMngr                IntegrationViewManager
	broadcasterColumnSelectionIndex int
	columnSelectionsForms           []*tview.Form
}

type TUIConfigurations struct {
	Integrations map[string]interface{} `json:"integrations"`
}

type Selection struct {
	Value         string
	relatedObject interface{}
}
