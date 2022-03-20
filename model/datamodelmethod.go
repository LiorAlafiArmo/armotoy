package model

import (
	"sort"
	"strconv"
	"strings"

	"k8s.io/utils/strings/slices"
)

func ShouldFilterByColumns(columns map[string][]string) bool {
	for _, values := range columns {
		if len(values) > 0 {
			return true
		}
	}
	return false
}

func (dm *DataModel) GetColumn(name string) (ColumnModel, int) {
	for i, col := range dm.Columns {
		if strings.EqualFold(name, col.Label) {
			return col, i
		}
	}

	return ColumnModel{}, -1
}

func (dm *DataModel) FilterByColumns(columns map[string][]string) *DataModel {
	if len(columns) == 0 {
		return dm
	}

	dmc := &DataModel{Type: dm.Type, Columns: dm.Columns, Data: make([]ElementModel, 0)}
	for i := len(dm.Data) - 1; i >= 0; i-- {
		if dm.isFilterable(columns, &dm.Data[i]) {
			dmc.Data = append(dmc.Data, dm.Data[i])
		}

	}

	return dmc
}

func (dm *DataModel) isFilterable(columns map[string][]string, element *ElementModel) bool {
	for columnName, values := range columns {
		if _, indx := dm.GetColumn(columnName); indx > len(element.Fields) {
			return false

			//not or to avoid element.Fields[indx] where indx oob
		} else if !slices.Contains(values, element.Fields[indx]) {
			return false
		}

	}
	return true
}

func (dm *DataModel) Sort(less func(i, j int) bool) {
	sort.Slice(dm.Data, less)
}

func (dm *DataModel) StableSort(less func(i, j int) bool) {
	sort.SliceStable(dm.Data, less)
}

func (dm *DataModel) SortByStatus(statusColumnName string) {
	sort.SliceStable(dm.Data, func(i, j int) bool {
		_, pos := dm.GetColumn(statusColumnName)
		if pos > -1 {
			return StatusPredicate(dm.Data[i].Fields[pos], dm.Data[j].Fields[pos])
		}

		return false
	})
}
func StatusPredicate(lhs, rhs string) bool {
	if lhs == "failed" {
		return true
	}

	if rhs == "failed" {
		return false
	}

	if lhs == "excluded" || lhs == "warning" {
		return true
	}

	if rhs == "excluded" || rhs == "warning" {
		return false
	}

	if lhs == "passed" {
		return true
	}

	if rhs == "passed" {
		return false
	}

	return false
}

// order matters
func (dm *DataModel) SortByColumns(columns []string) {

	sort.Slice(dm.Data, func(i, j int) bool {
		for k := range columns {
			if coldata, col := dm.GetColumn(columns[k]); col > -1 {
				if coldata.Type == DATA_TYPE_NUMBER {
					lhs, _ := strconv.ParseFloat(dm.Data[i].Fields[col], 32)
					rhs, _ := strconv.ParseFloat(dm.Data[j].Fields[col], 32)

					if lhs > rhs {
						return false
					}

					if lhs < rhs {
						return true
					}

				} else if coldata.Label == "Status" {
					return StatusPredicate(dm.Data[i].Fields[col], dm.Data[j].Fields[col])
				} else {
					if dm.Data[i].Fields[col] > dm.Data[j].Fields[col] {
						return false
					}

					if dm.Data[i].Fields[col] < dm.Data[j].Fields[col] {
						return true
					}
				}

			}
		}
		return true
	})

}
