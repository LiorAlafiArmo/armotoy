package broadcasters

type IBroadcaster interface {
	SendMessage(severity, title, message string) error
	ExportConfig() map[string]interface{}
	AddTarget(target string) error
	RemoveTarget(target string) error
	FindTarget(target string) int
	GetTargets() []string
	GetType() string

	// AddTargets(target ...string)
}
