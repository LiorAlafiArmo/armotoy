package view

import (
	"github.com/armosec/armotoy/broadcasters"
	"github.com/armosec/armotoy/common"
	"github.com/armosec/armotoy/model"
	"github.com/rivo/tview"
)

type Data struct {
	ColumnAttributes map[string]common.ColumnAttributes
	Model            *model.DataModel `json:",inline"`
}

type IntegrationView struct {
	Title       *tview.TextView
	Broadcaster broadcasters.IBroadcaster
	Error       *tview.TextView
	Info        *tview.TextView
	Form        *tview.Form
}
