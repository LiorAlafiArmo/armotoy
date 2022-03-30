package broadcasters

import (
	"fmt"
	"net/smtp"
	"regexp"

	"github.com/mitchellh/mapstructure"
)

const (
	EMAIL_REGEX = `^(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`
)

type SMTPInfo struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	IsTLS bool   `json:"TLS"`
}

type EmailAuth struct {
	Password string `json:"password"`
	User     string `json:"user"`
}
type EmailConfiguration struct {
	SMTP SMTPInfo  `json:"smtp"`
	Auth EmailAuth `json:"auth"`
}

type EmailRecipient struct {
	Email string
}
type EMailBroadcaster struct {
	Config     EmailConfiguration
	Recipients []EmailRecipient
	SMTPAuth   smtp.Auth
}

func (b *EMailBroadcaster) GetFullHost() string {
	return fmt.Sprintf("%s:%d", b.Config.SMTP.Host, b.Config.SMTP.Port)
}

func (b *EMailBroadcaster) GetTargets() []string {
	recipients := make([]string, 0, len(b.Recipients))
	for i := range b.Recipients {
		recipients = append(recipients, b.Recipients[i].Email)
	}

	return recipients
}

func (b *EMailBroadcaster) SendMessage(severity, title, message string) error {
	msg := fmt.Sprintf("[%s]%s:\n\n%s\n\n", severity, title, message)
	return smtp.SendMail(b.GetFullHost(),
		b.SMTPAuth,
		b.Config.Auth.User, b.GetTargets(), []byte(msg))
}

func (b *EMailBroadcaster) ExportConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (b *EMailBroadcaster) FindTarget(target string) int {
	for i := range b.Recipients {
		if b.Recipients[i].Email == target {
			return i
		}
	}

	return -1
}

func (b *EMailBroadcaster) AddTarget(target string) error {
	match, _ := regexp.MatchString(EMAIL_REGEX, target)

	pos := b.FindTarget(target)
	if match && pos == -1 {
		b.Recipients = append(b.Recipients, EmailRecipient{Email: target})
		return nil
	} else if pos > -1 {
		return fmt.Errorf("existing email address")
	}
	return fmt.Errorf("%s is invalid email address", target)
}

func (b *EMailBroadcaster) GetType() string {
	return "email"
}

func (b *EMailBroadcaster) RemoveTarget(target string) error {
	if pos := b.FindTarget(target); pos > -1 {
		b.Recipients = append(b.Recipients[:pos], b.Recipients[pos+1:]...)
		return nil
	}
	return fmt.Errorf("%s doesn't exist", target)
}

func EMailBroadcasterInit(config interface{}) (*EMailBroadcaster, error) {
	emailBroadcaster := &EMailBroadcaster{Recipients: make([]EmailRecipient, 0)}
	if err := mapstructure.Decode(config, &emailBroadcaster.Config); err != nil {
		return nil, err
	}

	emailBroadcaster.SMTPAuth = smtp.PlainAuth("", emailBroadcaster.Config.Auth.User, emailBroadcaster.Config.Auth.Password, emailBroadcaster.Config.SMTP.Host)

	return emailBroadcaster, nil
}
