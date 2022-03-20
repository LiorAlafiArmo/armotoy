package broadcasters

import (
	"fmt"

	"github.com/armosec/slacker/slacker"
)

func SlackSender(itemsType, alertLevel, channel, token, baseText, msg string) error {
	bot, _ := slacker.SlackBotInit("", token, false)

	id, err := bot.FindChannel(channel)
	if err != nil {
		return err
	}

	msg = fmt.Sprintf("%s->\n%s\n%s", itemsType, baseText, msg)
	return bot.SendMessage(alertLevel, id, msg)

}
