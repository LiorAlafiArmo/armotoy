package broadcasters

import (
	"fmt"

	"github.com/armosec/slacker/slacker"
	"github.com/mitchellh/mapstructure"
)

type SlackerConfiguration struct {
	Token string `json:"token"`
}

type SlackChannel struct {
	Id   string
	Name string
}
type SlackBroadcaster struct {
	Channels []SlackChannel
	Config   SlackerConfiguration
	slack    *slacker.SlackBot
}

func (b *SlackBroadcaster) SendMessage(severity, title, message string) error {
	for c := range b.Channels {
		err := b.slack.SendMessage(severity, b.Channels[c].Id, title+"->\n"+message)
		if err != nil {
			continue
		}

	}
	return nil
}
func (b *SlackBroadcaster) ExportConfig() map[string]interface{} {
	return nil
}

func (b *SlackBroadcaster) FindTarget(target string) int {
	for i := range b.Channels {
		if b.Channels[i].Name == target {
			return i
		}
	}

	return -1
}

func (b *SlackBroadcaster) AddTarget(target string) error {
	if len(target) == 0 {
		return fmt.Errorf("must specify a channel")
	}
	if target[0] != '#' {
		return fmt.Errorf("all slack channels begin with #")
	}
	pos := b.FindTarget(target)
	if pos != -1 {
		return nil
	}
	id, er := b.slack.FindChannel(target)
	if er != nil {
		return er
	}
	er = b.slack.JoinChannel(id)
	if er != nil {
		return er
	}
	b.Channels = append(b.Channels, SlackChannel{Id: id, Name: target})

	return nil
}

func (b *SlackBroadcaster) RemoveTarget(target string) error {
	pos := b.FindTarget(target)
	if pos > -1 {
		b.Channels = append(b.Channels[:pos], b.Channels[pos+1:]...)
	}

	return nil
}

func (b *SlackBroadcaster) GetTargets() []string {
	channels := make([]string, 0)
	for i := range b.Channels {
		channels = append(channels, b.Channels[i].Name)
	}

	return channels
}
func (b *SlackBroadcaster) GetType() string {
	return "slack"
}

func SlackBroadcasterInit(config interface{}) (*SlackBroadcaster, error) {

	slackbroadcer := &SlackBroadcaster{}
	if err := mapstructure.Decode(config, &slackbroadcer.Config); err != nil {
		return nil, err
	}
	bot, err := slacker.SlackBotInit("", slackbroadcer.Config.Token, false)
	if err == nil {
		slackbroadcer.slack = bot
		slackbroadcer.Channels = make([]SlackChannel, 0)
	}
	return slackbroadcer, err

}
