package notification

import (
	"fmt"
	"github.com/gregdel/pushover"
	log "github.com/sirupsen/logrus"
)

type pushoverHandler struct {
	pushover  *pushover.Pushover
	recipient *pushover.Recipient
}

func (p *pushoverHandler) Handle(pageInfo PageInfo) {
	m := pushover.Message{
		Message:  fmt.Sprintf("PS5 Available at %s", pageInfo.Name()),
		Title:    "PS5",
		Priority: 1,
		URL:      pageInfo.URL(),
		Sound:    "bugle",
	}

	_, err := p.pushover.SendMessage(&m, p.recipient)
	if err != nil {
		log.Error("Failed to send message")
		return
	}
}

func NewPushoverHandler(token string, user string) Handler {
	p := pushover.New(token)
	u := pushover.NewRecipient(user)
	return &pushoverHandler{
		pushover:  p,
		recipient: u,
	}
}
