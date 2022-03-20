package model

import (
	"github.com/armosec/opa-utils/gitregostore"
	"github.com/armosec/opa-utils/reporthandling"
	"github.com/armosec/opa-utils/reporthandling/results/v1/resourcesresults"
	v2 "github.com/armosec/opa-utils/reporthandling/v2"
)

type ResourceReference struct {
	Result            *resourcesresults.Result
	Raw               *reporthandling.Resource
	ControlsRelevancy []string
}

type AcceptedStatuses struct {
	Passed   bool
	Failed   bool
	Excluded bool
	Skipped  bool
}

type PostureModel struct {
	pr        *v2.PostureReport
	gitrego   *gitregostore.GitRegoStore
	Resources []ResourceReference
}

type DataModel struct {
	Type    string
	Data    []ElementModel
	Columns []ColumnModel
}

type Context struct {
	Type   string
	Values []string
}
type ElementModel struct {
	Fields   []string
	Data     interface{}
	Contexts []Context
}

type ColumnModel struct {
	Label  string
	Type   string
	Hidden bool
}

type ControlAttributes struct {
	CustomizedConfigurations bool
	HostScan                 bool
	RelevantCloudProviders   []string
}
