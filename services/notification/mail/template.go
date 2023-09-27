package mail

import (
	"bytes"
	"html/template"

	"gitlab.com/mefit/mefit-server/utils/assert"
)

const (
	WelcomeMail = iota
	InvoiceMail
)

var temps = map[int]string{
	WelcomeMail: "mail-welcome.html",
	InvoiceMail: "mail-invoice.html",
}

func ParseTemplate(temp int, in interface{}) string {
	t, err := template.New("mail").ParseFiles("templates/email/" + temps[temp])
	assert.Nil(err)

	var buffer bytes.Buffer
	err = t.ExecuteTemplate(&buffer, temps[temp], in)
	assert.Nil(err)

	return buffer.String()
}
