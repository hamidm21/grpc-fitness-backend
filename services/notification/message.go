package notification

import (
	"github.com/kataras/iris/core/errors"
	"gitlab.com/mefit/mefit-server/utils/assert"
)

type notifType string

const (
	MailType notifType = "mail"
	SMSType  notifType = "sms"
)

type Notification struct {
	From    string
	Subject string
	To      []string
	Body    string
	tYpe    notifType
}

func NewNotification(tYpe notifType) *Notification {
	return &Notification{tYpe: tYpe}
}

func (n *Notification) Validate() bool {
	assert.True(n.tYpe != "", errors.New("notification type needs to be specified"))

	switch n.tYpe {
	case MailType:
		if n.From == "" {
			n.From = "youtab@u-v.ir"
		}

		if len(n.To) == 0 {
			return false
		}
	case SMSType:
		return true
	}

	return true
}
