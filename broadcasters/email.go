package broadcasters

import (
	"fmt"
	"net/smtp"

	"github.com/mitchellh/mapstructure"
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
	Config               EmailConfiguration
	RecipientsBySeverity map[string][]EmailRecipient
	SMTPAuth             smtp.Auth
}

func (b *EMailBroadcaster) GetFullHost() string {
	return fmt.Sprintf("%s:%d", b.Config.SMTP.Host, b.Config.SMTP.Port)
}

func (b *EMailBroadcaster) GetRecipientsEmail(severity string) []string {
	recipients := make([]string, 0, len(b.RecipientsBySeverity[severity]))
	for i := range b.RecipientsBySeverity[severity] {
		recipients = append(recipients, b.RecipientsBySeverity[severity][i].Email)
	}

	return recipients
}

func (b *EMailBroadcaster) SendMessage(severity, title, message string) error {
	msg := fmt.Sprintf("[%s]%s:\n\n%s\n\n", severity, title, message)
	return smtp.SendMail(b.GetFullHost(),
		b.SMTPAuth,
		b.Config.Auth.User, b.GetRecipientsEmail(severity), []byte(msg))
}

func (b *EMailBroadcaster) ExportConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func EMailBroadcasterInit(config interface{}) (*EMailBroadcaster, error) {
	emailBroadcaster := &EMailBroadcaster{RecipientsBySeverity: make(map[string][]EmailRecipient)}
	if err := mapstructure.Decode(config, &emailBroadcaster.Config); err != nil {
		return nil, err
	}

	emailBroadcaster.SMTPAuth = smtp.PlainAuth("", emailBroadcaster.Config.Auth.User, emailBroadcaster.Config.Auth.Password, emailBroadcaster.Config.SMTP.Host)

	return emailBroadcaster, nil
}
