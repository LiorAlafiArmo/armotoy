package controller

// func addResourceYamlPage(table *tview.Table, controller *MainActionController) error {
// 	row, _ := table.GetSelection()

// 	resource := table.GetCell(row, 7)
// 	if resource.Text == "" {
// 		return fmt.Errorf("invalid resource selection")
// 	}
// 	//TODO use GetResourceContext
// 	p := controller.GetResourceContext(resource.Text)

// 	if p == nil {
// 		return fmt.Errorf("could not find relevant resource")
// 	}
// 	yamlbytes, err := yaml.Marshal(p.Raw.GetObject())
// 	if err != nil {
// 		return fmt.Errorf("could not extract yaml from resource")
// 	}

// 	list := tview.NewList()
// 	sc := 'a'
// 	list.AddItem("all", "apply all fixed & failed path controls", sc, nil)

// 	for _, c := range p.Result.AssociatedControls {
// 		if c.GetStatus(&v1.Filters{}).IsFailed() {

// 			list.AddItem(c.ControlID, c.Name, sc, nil)
// 		}
// 		sc++
// 	}

// 	yamlView := tview.NewTextView().
// 		SetDynamicColors(true).
// 		SetRegions(true).
// 		SetWordWrap(true).
// 		SetText(string(yamlbytes)).
// 		SetChangedFunc(func() {
// 			controller.app.Draw()
// 		})

// 	flex := tview.NewFlex()
// 	flex.AddItem(list, 0, 1, false)
// 	flex.AddItem(yamlView, 0, 2, false)
// 	flex.AddItem(yamlView, 0, 3, false) //smoke & mirrors

// 	button := tview.NewButton("Apply YAML changes").SetSelectedFunc(func() {
// 		delete(controller.contexts, "tempResource")
// 		AddControlsTable(controller, "Controls")
// 		controller.SetScene("Controls")
// 	})

// 	topelements := make([]tview.Primitive, 0)
// 	topelements = append(topelements, controller.contexts["Main Menu"].topElements...)
// 	topelements = append(topelements, button)
// 	controller.AddContext("tempResource", topelements, []tview.Primitive{flex})
// 	controller.SetScene("tempResource")
// 	return nil
// }
