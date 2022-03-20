package model

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/armosec/armotoy/common"
	"github.com/armosec/opa-utils/gitregostore"
	postureapis "github.com/armosec/opa-utils/reporthandling/apis"
	v1 "github.com/armosec/opa-utils/reporthandling/helpers/v1"
	"github.com/armosec/opa-utils/reporthandling/results/v1/reportsummary"
	"github.com/armosec/opa-utils/reporthandling/results/v1/resourcesresults"
	v2 "github.com/armosec/opa-utils/reporthandling/v2"
	"k8s.io/utils/strings/slices"
)

func UpdateLeaf(parent interface{}, o interface{}, keys []string, update string, leafValue string, reset string) (reflect.Kind, interface{}, interface{}) {
	switch v := reflect.ValueOf(o); v.Kind() {
	case reflect.Slice:
		tmp, _ := o.([]interface{})
		if len(keys) > 0 {
			k := keys[0][1 : len(keys[0])-1]
			i, _ := strconv.Atoi(k)

			_, _, leaf := UpdateLeaf(tmp, tmp[i], keys[1:], update, leafValue, reset)

			tmp[i] = leaf
			return reflect.Slice, parent, tmp
		}
		return reflect.Slice, parent, tmp

	case reflect.Map:
		tmp, _ := o.(map[string]interface{})
		if len(keys) > 0 {
			kind, _, leaf := UpdateLeaf(tmp, tmp[keys[0]], keys[1:], update, leafValue, reset)
			switch kind {
			case reflect.Bool, reflect.String, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
				tmp[update+keys[0]] = leaf
				delete(tmp, keys[0])
			}
			return reflect.Map, parent, tmp

		}
		return reflect.Map, parent, tmp

	default:

		data := fmt.Sprintf("%s%v%s", update, o, reset)
		return v.Kind(), parent, data
	}
}

func GetLeaf(parent interface{}, o interface{}, keys []string) (reflect.Kind, interface{}, interface{}) {
	switch v := reflect.ValueOf(o); v.Kind() {
	case reflect.Slice:
		tmp, _ := o.([]interface{})
		if len(keys) > 0 {
			k := keys[0][1 : len(keys[0])-1]
			i, _ := strconv.Atoi(k)

			return GetLeaf(tmp, tmp[i], keys[1:])
		}
		return reflect.Slice, parent, tmp

	case reflect.Map:
		tmp, _ := o.(map[string]interface{})
		if len(keys) > 0 {

			return GetLeaf(tmp, tmp[keys[0]], keys[1:])
		}
		return reflect.Map, parent, tmp

	default:
		return v.Kind(), parent, o
	}
}

func PostureModelInit(source, path, version string) (*PostureModel, error) {
	var pm *PostureModel = nil
	var err error
	switch source {
	case "JSON":
		pm, err = jsonPostureModelMaker(path, version)
	default:
		err = fmt.Errorf("not implemented")
	}
	if err != nil {
		return nil, err
	}

	pm.gitrego.SetRegoObjects()
	return pm, err

}

func jsonPostureModelMaker(path, version string) (*PostureModel, error) {
	if version == "v2" {
		pr := &v2.PostureReport{}
		if path == "" {
			return nil, fmt.Errorf("empty path, must specify a valid path")
		}
		b, err := os.ReadFile(path)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(b, pr); err != nil {
			return nil, err
		}
		return &PostureModel{pr: pr, gitrego: gitregostore.NewDefaultGitRegoStore(0)}, nil
	}

	return nil, fmt.Errorf("not implemented")

}
func (pm *PostureModel) GetResourcesTable(frameworks []string, controls []string, filters map[string][]string) (*DataModel, error) {
	if pm.pr == nil {
		return nil, fmt.Errorf("missing posture status")
	}
	datamodel := &DataModel{Columns: []ColumnModel{{Label: "Namespace"},
		{Label: "Kind"},
		{Label: "Name"},
		{Label: "Status"},
		{Label: "Info"},
		{Label: "k8s Object"},
		{Label: "Failed Controls"},
		{Label: "ResourceID", Hidden: true},
	}, Data: make([]ElementModel, 0), Type: "Resources"}

	resources := make([]ResourceReference, 0)

	for i := range pm.pr.Results {
		for j := range pm.pr.Resources {
			if pm.pr.Results[i].ResourceID == pm.pr.Resources[j].ResourceID {
				if len(controls) == 0 || IsResultContainsControlID(&pm.pr.Results[i], frameworks, controls) {
					resources = append(resources, ResourceReference{Result: &pm.pr.Results[i], Raw: &pm.pr.Resources[j]})
					break
				}

			}
		}

	}

	for i := range resources {
		isK8sObj := ""
		if resources[i].Raw.GetObjectType() == "workload" {
			isK8sObj = "true"
		}

		element := ElementModel{Fields: []string{
			resources[i].Raw.GetNamespace(),
			resources[i].Raw.GetKind(),
			resources[i].Raw.GetName(),
			string(resources[i].Result.GetStatus(&v1.Filters{FrameworkNames: frameworks}).Status()),
			resources[i].Result.GetStatus(&v1.Filters{FrameworkNames: frameworks}).Info(),
			isK8sObj,
			strings.Join(resources[i].Result.ListControlsIDs(&v1.Filters{FrameworkNames: frameworks}).Failed(), ", "),
			resources[i].Raw.ResourceID,
		}, Data: &resources[i]}
		datamodel.Data = append(datamodel.Data, element)
	}
	datamodel.FilterByColumns(filters)
	return datamodel, nil
}

func IsResultContainsControlID(result *resourcesresults.Result, frameworks []string, controls []string) bool {
	failedNexcluded := result.ListControlsIDs(&v1.Filters{FrameworkNames: frameworks}).Failed()
	failedNexcluded = append(failedNexcluded, result.ListControlsIDs(&v1.Filters{FrameworkNames: frameworks}).Excluded()...)
	for _, ctrlid := range failedNexcluded {
		if slices.Contains(controls, ctrlid) {
			return true
		}
	}
	return false
}

func (pm *PostureModel) GetFrameworksTable() (*DataModel, error) {
	if pm.pr == nil {
		return nil, fmt.Errorf("missing posture status")
	}
	datamodel := &DataModel{Columns: []ColumnModel{{Label: "Name"},
		{Label: "Status"},
		{Label: "Status Info"},
		{Label: "Score", Type: DATA_TYPE_NUMBER},
	}, Data: make([]ElementModel, 0, len(pm.pr.SummaryDetails.Frameworks)), Type: "Frameworks"}
	for i := range pm.pr.SummaryDetails.Frameworks {
		element := ElementModel{Fields: []string{
			pm.pr.SummaryDetails.Frameworks[i].Name,
			string(pm.pr.SummaryDetails.Frameworks[i].GetStatus().Status()),
			pm.pr.SummaryDetails.Frameworks[i].GetStatus().Info(),
			fmt.Sprintf("%f", pm.pr.SummaryDetails.Frameworks[i].GetScore()),
		}, Data: &pm.pr.SummaryDetails.Frameworks[i]}
		datamodel.Data = append(datamodel.Data, element)
	}
	return datamodel, nil
}

func (pm *PostureModel) GetControlsTable(frameworks []string) (*DataModel, error) {
	if pm.pr == nil {
		return nil, fmt.Errorf("missing posture status")
	}
	datamodel := &DataModel{Columns: []ColumnModel{{Label: "ID"},
		{Label: "Name"},
		{Label: "Status"},
		{Label: "Status Info"},
		{Label: "Score", Type: "number"},
		{Label: "Severity"},
		{Label: "Host Scan"},
		{Label: "Customized Config"},
		{Label: "Cloud Related"},
	}, Type: common.CONTEXT_CONTROLS}

	if len(frameworks) == 0 {
		pm.AddControls2DM(datamodel, "", pm.pr.SummaryDetails.Controls, nil)

		return datamodel, nil
	}
	controlExist := make(map[string]int)
	for i := range frameworks {
		for j := range pm.pr.SummaryDetails.Frameworks {
			if strings.EqualFold(pm.pr.SummaryDetails.Frameworks[j].Name, frameworks[i]) {
				pm.AddControls2DM(datamodel, pm.pr.SummaryDetails.Frameworks[j].GetName(), pm.pr.SummaryDetails.Frameworks[j].Controls, controlExist)
			}
		}
	}

	return datamodel, nil

}

func (pm *PostureModel) AddControls2DM(datamodel *DataModel, framework string, controls reportsummary.ControlSummaries, controlExist map[string]int) {
	for cidx := range controls {
		tmp := controls[cidx]
		pm.AddControlSummaryToModel(datamodel, &tmp, framework, controlExist)
	}
}

func (pm *PostureModel) GetControlPropertiesFromResults(id string) *ControlAttributes {
	attr := &ControlAttributes{CustomizedConfigurations: false, HostScan: true, RelevantCloudProviders: make([]string, 0)}
	for i := range pm.pr.Results {
		for j := range pm.pr.Results[i].AssociatedControls {
			if pm.pr.Results[i].AssociatedControls[j].ControlID == id {
				for _, rule := range pm.pr.Results[i].AssociatedControls[j].ResourceAssociatedRules {
					if len(rule.ControlConfigurations) > 0 {
						attr.CustomizedConfigurations = true
					}
				}
			}
		}
	}
	ctrl, err := pm.gitrego.GetOPAControlByID(id)

	if err == nil && ctrl != nil {
		for _, rule := range ctrl.Rules {
			if hostSensorRule, ok := rule.Attributes["hostSensorRule"]; !ok || hostSensorRule != "true" {
				attr.HostScan = false
			}

			for i := range rule.RelevantCloudProviders {
				if !slices.Contains(attr.RelevantCloudProviders, rule.RelevantCloudProviders[i]) {
					attr.RelevantCloudProviders = append(attr.RelevantCloudProviders, rule.RelevantCloudProviders[i])
				}
			}
		}
	} else {
		attr.HostScan = false
	}

	return attr
}

func (pm *PostureModel) AddControlSummaryToModel(datamodel *DataModel, control *reportsummary.ControlSummary, framework string, existance map[string]int) {
	if existance == nil {
		if pos, ok := existance[control.ControlID]; ok && framework != "" {
			for i := range datamodel.Data[pos].Contexts {
				if datamodel.Data[pos].Contexts[i].Type == common.CONTEXT_FRAMEWORK && !slices.Contains(datamodel.Data[pos].Contexts[i].Values, framework) {
					datamodel.Data[pos].Contexts[i].Values = append(datamodel.Data[pos].Contexts[i].Values, framework)
					return
				}
			}
		}
	}
	controlAttrib := pm.GetControlPropertiesFromResults(control.ControlID)
	ctrl := ElementModel{
		Fields: []string{control.ControlID,
			control.Name,
			string(control.GetStatus().Status()),
			control.GetStatus().Info(),
			fmt.Sprintf("%f", control.Score),
			postureapis.ControlSeverityToString(control.ScoreFactor),
			fmt.Sprintf("%v", controlAttrib.HostScan),
			fmt.Sprintf("%v", controlAttrib.CustomizedConfigurations),
			strings.Join(controlAttrib.RelevantCloudProviders, ","),
		},
		Data: control,
	}

	if framework != "" {
		ctrl.Contexts = append(ctrl.Contexts, Context{Type: common.CONTEXT_FRAMEWORK, Values: []string{framework}})
		if existance == nil {
			existance = map[string]int{}
		}
		existance[control.ControlID] = 1
	}

	datamodel.Data = append(datamodel.Data, ctrl)
}
