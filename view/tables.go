package view

import (
	"fmt"
	"strings"

	"github.com/armosec/armotoy/common"
	"github.com/armosec/armotoy/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"k8s.io/utils/strings/slices"
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

func CreateTableFilterForm(columnsAttributes map[string]common.ColumnAttributes, filters *common.Filters, txt *tview.TextView) *tview.Form {
	keys := make([]string, 0, len(columnsAttributes))
	for k := range columnsAttributes {
		keys = append(keys, k)
	}
	if filters.Equals == nil {
		filters.Equals = make(map[string][]string)
	}
	form := tview.NewForm().
		AddDropDown("Field", keys, 0, nil).
		AddInputField("Value", "", 30, nil, nil)

	// form.SetFocus()
	form.AddButton("Add", func() {
		ddown := form.GetFormItemByLabel("Field")
		if ddown == nil {
			return
		}

		dropdown, ok := ddown.(*tview.DropDown)
		if ok {
			_, column := dropdown.GetCurrentOption()
			inp := form.GetFormItemByLabel("Value")
			input, ok := inp.(*tview.InputField)
			if !ok {
				return
			}
			data, ok := filters.Equals[column]
			if !ok {
				data = make([]string, 0)
			}
			if len(input.GetText()) > 0 {
				pos := slices.Index(data, input.GetText())
				if pos > -1 {
					data = append(data[:pos], data[pos+1:]...)

				} else {
					data = append(data, input.GetText())
				}
			}
			if len(data) > 0 {
				filters.Equals[column] = data
			} else {
				delete(filters.Equals, column)
			}
			input.SetText("")
			SetFiltersText(txt, filters)

		}

	})
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	// form.AddButton("Filter", func() {

	// })

	return form
}

func SetFiltersText(txt *tview.TextView, filters *common.Filters) {
	if txt != nil {
		s := "Filters:\n"
		for column := range filters.Equals {
			s = fmt.Sprintf("%s%s: %s\n", s, column, strings.Join(filters.Equals[column], ", "))
		}
		txt.SetText(s)
	}
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
