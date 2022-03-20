package controller

import (
	"github.com/armosec/armotoy/broadcasters"
	"github.com/armosec/armotoy/common"
	"github.com/armosec/armotoy/model"
	"github.com/rivo/tview"
)

type Integration struct {
	Broadcaster broadcasters.IBroadcaster
	Enabled     bool
}
type Controller struct {
	model        *model.PostureModel
	app          *tview.Application
	CurrentState string
	StateMap     map[string]common.State
	root         *tview.Grid

	Selections map[string][]Selection

	DiggingSelection []Selection
}

type Selection struct {
	Value         string
	relatedObject interface{}
}
