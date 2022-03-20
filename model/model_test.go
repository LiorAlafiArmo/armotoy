package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLeaf(t *testing.T) {
	r, e := PostureModelInit("JSON", "../myres.json", "v2")
	assert := assert.New(t)
	assert.Nil(e, "unable to load data")

	o := r.pr.Resources[0].GetObject()
	kind, parent, leaf := UpdateLeaf(nil, o, []string{"rules", "[0]", "verbs", "[0]"}, "[red]", "", "[white]")

	t.Errorf("%v %v %v", kind, parent, leaf)
}

func TestJSONLoading(t *testing.T) {
	r, e := PostureModelInit("JSON", "./clusterscan.json", "v2")
	if e != nil || r == nil {
		t.Errorf("could not load posture properly")
	}

	rs, _ := r.GetResourcesTable([]string{}, []string{}, map[string][]string{})
	col, i := rs.GetColumn("namespace")
	if i < 0 {
		t.Errorf("missing column")
	}

	d := rs.FilterByColumns(map[string][]string{col.Label: {"kube-system"}})
	if e != nil || r == nil {
		t.Errorf("could not load posture properly")
	}

	for _, v := range d.Data {
		fmt.Printf("%v\n", v.Fields)
	}
	t.Errorf("1")
}

func TestModelSingleFilter(t *testing.T) {
	r, e := PostureModelInit("JSON", "./clusterscan.json", "v2")
	if e != nil || r == nil {
		t.Errorf("could not load posture properly")
	}

	rs, _ := r.GetResourcesTable([]string{}, []string{}, map[string][]string{})
	col, i := rs.GetColumn("namespace")
	if i < 0 {
		t.Errorf("missing column")
	}

	d := rs.FilterByColumns(map[string][]string{col.Label: {"kube-system"}})
	if e != nil || r == nil {
		t.Errorf("could not load posture properly")
	}

	for _, v := range d.Data {
		fmt.Printf("%v\n", v.Fields)
	}
	t.Errorf("1")
}

func TestModelMultiFilter(t *testing.T) {
	r, e := PostureModelInit("JSON", "./clusterscan.json", "v2")
	if e != nil || r == nil {
		t.Errorf("could not load posture properly")
	}

	rs, _ := r.GetResourcesTable([]string{}, []string{}, map[string][]string{})
	col, i := rs.GetColumn("namespace")
	if i < 0 {
		t.Errorf("missing column")
	}

	d := rs.FilterByColumns(map[string][]string{col.Label: {"kube-system"}, "Kind": {"ConfigMap"}})
	if e != nil || r == nil {
		t.Errorf("could not load posture properly")
	}

	for _, v := range d.Data {
		fmt.Printf("%v\n", v.Fields)
	}
	t.Errorf("1")
}

func TestModelMultiSort(t *testing.T) {
	r, e := PostureModelInit("JSON", "./clusterscan.json", "v2")
	if e != nil || r == nil {
		t.Errorf("could not load posture properly")
	}

	rs, _ := r.GetResourcesTable([]string{}, []string{}, map[string][]string{})

	rs.SortByColumns([]string{"Namespace", "Kind", "Name", "Status"})
	if e != nil || r == nil {
		t.Errorf("could not load posture properly")
	}

	for _, v := range rs.Data {
		fmt.Printf("%v\n", v.Fields)
	}
	t.Errorf("1")
}
