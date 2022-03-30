package broadcasters

import "fmt"

func Factory(broadcaster string, config interface{}) (IBroadcaster, error) {
	switch broadcaster {
	case "email":
		return EMailBroadcasterInit(config)
	case "slack":
		return SlackBroadcasterInit(config)
	}

	return nil, fmt.Errorf("unknown broadcaster type")
}
