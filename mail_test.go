package mailgun

import (
	"fmt"
	"testing"
)

const (
	_domain  = "***"
	_api     = "***"
	_to      = "***"
	_from    = "***"
	_subject = "Hey"
	_text    = "Body"
)

var (
	mailer Mail
)

func TestCreate(t *testing.T) {
	mailer = DebugMailClient(_domain, _api)
	mailer.Create(_to, _from, _subject, _text)
}

func TestSend(t *testing.T) {
	out, err := mailer.Send()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(out)
}
