package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/armotoy/common"
	"github.com/armosec/armotoy/model"
	"github.com/armosec/armotoy/view"
	"github.com/armosec/opa-utils/reporthandling/results/v1/reportsummary"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"k8s.io/utils/strings/slices"
)

func (c *Controller) CreateAnException(selectedFrameworks []Selection, selectedControls []Selection, selectedResources []Selection) *armotypes.PostureExceptionPolicy {
	now := time.Now().String()
	name := fmt.Sprintf("time_%s_customer_%s_cluster_%s_%s", now, c.model.GetCustomerGUID(), c.model.GetCluster(), c.model.GetReportGUID())
	exception := &armotypes.PostureExceptionPolicy{
		PortalBase:      armotypes.PortalBase{Name: name, Attributes: map[string]interface{}{"generatedby": "armotoy", "cluster": c.model.GetCluster(), "customerGUID": c.model.GetCustomerGUID()}},
		CreationTime:    now,
		PolicyType:      "postureExceptionPolicy",
		Actions:         []armotypes.PostureExceptionPolicyActions{armotypes.AlertOnly, armotypes.Disable},
		Resources:       make([]armotypes.PortalDesignator, 0),
		PosturePolicies: make([]armotypes.PosturePolicy, 0),
	}

	for _, fw := range selectedFrameworks {
		for _, ctrl := range selectedControls {
			framework, fwok := fw.relatedObject.(*reportsummary.FrameworkSummary)
			if fwok {
				ctrlids := framework.ListControls().ListControlsIDs().All()
				if slices.Contains(ctrlids, ctrl.Value) {
					policy := armotypes.PosturePolicy{FrameworkName: fw.Value, ControlID: ctrl.Value}
					exception.PosturePolicies = append(exception.PosturePolicies, policy)
				}

			}

		}
	}

	if len(selectedFrameworks) == 0 {
		for _, ctrl := range selectedControls {
			policy := armotypes.PosturePolicy{ControlID: ctrl.Value}
			exception.PosturePolicies = append(exception.PosturePolicies, policy)

		}
	}
	if len(selectedResources) > 0 {
		for _, resource := range selectedResources {
			rsrc, ok := resource.relatedObject.(*model.ResourceReference)
			if ok {
				r := armotypes.PortalDesignator{DesignatorType: "attribute", Attributes: map[string]string{"cluster": c.model.GetCluster(), "customerGUID": c.model.GetCustomerGUID(), "kind": rsrc.Raw.GetKind(), "namespace": rsrc.Raw.GetNamespace(), "name": rsrc.Raw.GetName(), "apiVersion": rsrc.Raw.GetApiVersion(), "resourceID": rsrc.Raw.ResourceID}}
				exception.Resources = append(exception.Resources, r)
			}
		}
	} else {
		r := armotypes.PortalDesignator{DesignatorType: "attribute", Attributes: map[string]string{"cluster": c.model.GetCluster(), "customerGUID": c.model.GetCustomerGUID()}}

		exception.Resources = append(exception.Resources, r)
	}

	return exception
}

func (c *Controller) CreateExceptionPage() (*common.State, error) {
	selectedResources := c.GetSelection(common.CONTEXT_RESOURCE)
	selectedFrameworks := c.GetSelection(common.CONTEXT_FRAMEWORK)
	selectedControls := c.GetSelection(common.CONTEXT_CONTROLS)

	s := ""

	s = view.SelectionText(GetValuesFromSelection(selectedFrameworks), "frameworks", s)
	s = view.SelectionText(GetValuesFromSelection(selectedControls), "controls", s)
	if len(selectedResources) > 0 {
		s = fmt.Sprintf("%sSelected resources:\n", s)
		for i := range selectedResources {
			r, ok := selectedResources[i].relatedObject.(*model.ResourceReference)
			if ok {
				s = fmt.Sprintf("%sAPI ver.: %s\t-\tNamespace: %s\t-\tKind: %s\t-\tName: %s\n", s, r.Raw.GetApiVersion(), r.Raw.GetNamespace(), r.Raw.GetKind(), r.Raw.GetName())

			}

		}
	}

	flex, input, review := view.CreateExceptionLayout(s, nil)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			file := input.GetText()
			ex := c.CreateAnException(selectedFrameworks, selectedControls, selectedResources)
			if ex != nil && strings.HasSuffix(file, ".json") {
				if b, err := json.Marshal(*ex); err == nil {
					os.WriteFile(file, b, 0777)
					review.SetText(string(b))
				}

			}
			return nil
		}
		return event
	})
	state := &common.State{Name: "exception", Index: 5, Menu: view.MakeMenu("exception")}
	state.Content = flex
	state.Footer = tview.NewTextView().SetText("ctrl+S - save exception")
	return state, nil

}
