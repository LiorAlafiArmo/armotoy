package controller

func GetValuesFromSelection(selection []Selection) []string {
	values := make([]string, 0, len(selection))

	for i := range selection {
		values = append(values, selection[i].Value)
	}
	return values
}
