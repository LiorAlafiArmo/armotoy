package controller

import "strings"

func (c *Controller) IndexOf(key, value string) int {
	data, ok := c.Selections[key]
	if !ok {
		return -1
	}

	for i := range data {
		if strings.EqualFold(value, data[i].Value) {
			return i
		}
	}
	return -1
}

func (c *Controller) AddSelection(key, value string, robj interface{}) int {
	data, ok := c.Selections[key]
	if !ok {
		data = make([]Selection, 0)

	}
	tmp := Selection{Value: value}
	if robj != nil {
		tmp.relatedObject = robj
	}
	data = append(data, tmp)
	c.Selections[key] = data
	return len(data) - 1
}

func (c *Controller) RemoveSelection(key, value string) {
	pos := c.IndexOf(key, value)
	if pos == -1 {
		return
	}
	data := c.Selections[key]
	c.Selections[key] = append(data[:pos], data[pos+1:]...)
	if len(c.Selections[key]) == 0 {
		c.ClearSelection(key)
	}
}

func (c *Controller) ClearSelection(key string) {
	delete(c.Selections, key)
}

func (c *Controller) GetSelection(key string) []Selection {
	return c.Selections[key]
}

func (c *Controller) GetSelectionValues(key string) ([]string, []interface{}) {
	values := make([]string, 0)
	data := make([]interface{}, 0)
	selections := c.GetSelection(key)
	for i := range selections {
		values = append(values, selections[i].Value)
		data = append(data, selections[i].relatedObject)

	}

	return values, data
}
