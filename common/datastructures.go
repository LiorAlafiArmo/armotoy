package common

import (
	"github.com/armosec/armotoy/broadcasters"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Integration struct {
	Broadcaster broadcasters.IBroadcaster
	IsActive    bool
}

type Filters struct {
	Equals map[string][]string
}
type ColumnAttributes struct {
	Hidden      bool
	Color       func(value string) tcell.Color
	ColumnColor tcell.Color
	Label       string
	Key         string
	Index       int
	Selected    bool
}
type State struct {
	Name    string
	Index   int
	Menu    tview.Primitive
	Content tview.Primitive
	Footer  tview.Primitive
}

type BroadcastOptions struct {
	FrameworksEnabled bool
	ControlsEnabled   bool
	ResourcesEnabled  bool
	Title             string
	Message           string
	Severity          string
}
