package mailClient

import (
	"log"

	"github.com/wneessen/go-mail"
)

func MailClient(name, email, msg, htmlcode string) {
	m := mail.NewMsg()
	if err := m.From("rishisheshe@outlook.com"); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := m.To(email); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}
	m.Subject("This is from DLP ADMIN")
	m.SetBodyString(mail.TypeTextHTML, htmlcode)
	c, err := mail.NewClient("smtp-mail.outlook.com", mail.WithPort(587), mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithUsername("rishisheshe@outlook.com"), mail.WithPassword("rishi@sheshe"))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	if err := c.DialAndSend(m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
}
