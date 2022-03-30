package common

import "github.com/gdamore/tcell/v2"

func DefaultColumnAttributes(col int) *ColumnAttributes {
	return &ColumnAttributes{Color: nil, Hidden: false, ColumnColor: tcell.ColorYellow, Label: "", Index: col}
}

func (ca *ColumnAttributes) SetIndex(col int) *ColumnAttributes {
	ca.Index = col
	return ca
}

func (c *BroadcastOptions) IsValid() bool {
	return c.ControlsEnabled || c.FrameworksEnabled || c.ResourcesEnabled
}
