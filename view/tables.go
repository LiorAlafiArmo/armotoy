package view

import (
	"github.com/armosec/armotoy/common"
	"github.com/armosec/armotoy/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	ESeverityCritical = "Critical"
	ESeverityHigh     = "High"
	ESeverityMedium   = "Medium"
	ESeverityLow      = "Low"
)

func ClearSelection(table *tview.Table) {
	for i := 1; i < table.GetRowCount(); i++ {
		table.GetCell(i, 0).SetTextColor(tcell.ColorWhite)
	}
}
func CreateTable(dm *model.DataModel, columnsAttributes map[string]common.ColumnAttributes) *tview.Table {
	table := tview.NewTable().SetBorders(true)
	CreateHeader(dm, columnsAttributes, table)
	for row := range dm.Data {
		for colName := range columnsAttributes {
			if columnsAttributes[colName].Hidden {
				continue
			}
			if _, idx := dm.GetColumn(colName); idx > -1 {
				color := tcell.ColorWhite
				if columnsAttributes[colName].Color != nil {
					color = columnsAttributes[colName].Color(dm.Data[row].Fields[idx])
				}
				table.SetCell(row+1, columnsAttributes[colName].Index,
					tview.NewTableCell(dm.Data[row].Fields[idx]).
						SetTextColor(color).
						SetReference(dm.Data[row].Data).
						SetAlign(tview.AlignCenter)).SetBorder(true)

			}
		}
	}
	return table
}

func CreateHeader(dm *model.DataModel, columnsAttributes map[string]common.ColumnAttributes, table *tview.Table) {
	for colName := range columnsAttributes {
		if columnsAttributes[colName].Hidden {
			continue
		}

		if col, i := dm.GetColumn(colName); i > -1 {
			label := columnsAttributes[colName].Label
			if label == "" {
				label = col.Label
			}
			color := tcell.ColorYellow
			if columnsAttributes[colName].ColumnColor != 0 {
				color = columnsAttributes[colName].ColumnColor
			}
			table.SetCell(0, columnsAttributes[colName].Index,
				tview.NewTableCell(label).
					SetTextColor(color).
					SetAlign(tview.AlignCenter)).SetBorder(true)
		}
	}
}

func SeverityColumnAttributes(col int) *common.ColumnAttributes {
	colAttr := common.DefaultColumnAttributes(col)
	colAttr.Color = func(value string) tcell.Color {
		return Severity2Color(value)
	}
	return colAttr
}

func StatusColumnAttributes(col int) *common.ColumnAttributes {
	colAttr := common.DefaultColumnAttributes(col)
	colAttr.Color = func(status string) tcell.Color {
		switch status {
		case "failed":
			return tcell.ColorRed
		case "passed":
			return tcell.ColorGreen
		case "excluded":
			return tcell.ColorYellow

		case "skipped":
			return tcell.ColorGray

		}
		return tcell.ColorWhite
	}
	return colAttr
}

func Severity2Color(severity string) tcell.Color {

	switch severity {
	case ESeverityCritical:
		return tcell.ColorDarkRed
	case ESeverityHigh:
		return tcell.ColorRed
	case ESeverityMedium:
		return tcell.ColorOrange
	case ESeverityLow:
		return tcell.ColorYellow

	}
	return tcell.ColorWhite
}
