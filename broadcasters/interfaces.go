package broadcasters

type IBroadcaster interface {
	SendMessage(severity, title, message string)
}
