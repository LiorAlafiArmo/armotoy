package broadcasters

type IBroadcaster interface {
	SendMessage(severity, title, message string) error
	ExportConfig() map[string]interface{}
}
