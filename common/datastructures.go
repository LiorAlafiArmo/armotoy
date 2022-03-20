package common

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ColumnAttributes struct {
	Hidden      bool
	Color       func(value string) tcell.Color
	ColumnColor tcell.Color
	Label       string
	Index       int
}
type State struct {
	Name    string
	Index   int
	Menu    tview.Primitive
	Content tview.Primitive
	Footer  tview.Primitive
}
