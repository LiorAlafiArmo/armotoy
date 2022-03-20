package model

type IPostureMode interface {
	GetFrameworksTable() (*DataModel, error)
	GetControlsTable(frameworks []string) (*DataModel, error)
	GetResourcesTable() (*DataModel, error)
}
