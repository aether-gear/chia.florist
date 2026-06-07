# Suggested Shared Mailer Structure

Since your mailer is a shared infrastructure capability, I’d place it under:

```text
internal/shared/mailer
```

Recommended structure:

```text
internal/shared/mailer
├── contract.go
├── sender.go
├── template.go        (optional later)
├── resend.go          (optional later)
└── smtp.go            (optional later)
```

---

# contract.go

```go
package mailer

// Sender is infrastructure capability for delivering emails.
//
// This interface should stay generic because multiple domains
// may use email delivery:
//
// - authentication
// - merchant onboarding
// - invoices
// - notifications
// - support
// - marketing
//
// Business-specific methods should live in usecase/service layer.
type Sender interface {
	Send(input SendInput) error
}

// SendInput represents generic email delivery payload.
type SendInput struct {
	To string

	Subject string

	// Plain text content.
	Text string

	// Optional HTML content.
	HTML *string
}
```

---

# sender.go

Example implementation using SMTP.

```go
package mailer

import (
	"fmt"
	"net/smtp"
)

// SMTPSender implements Sender using SMTP transport.
type SMTPSender struct {
	host string
	port string

	username string
	password string

	from string
}

func NewSMTPSender(
	host string,
	port string,
	username string,
	password string,
	from string,
) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) Send(input SendInput) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	auth := smtp.PlainAuth(
		"",
		s.username,
		s.password,
		s.host,
	)

	message := s.buildMessage(input)

	if err := smtp.SendMail(
		addr,
		auth,
		s.from,
		[]string{input.To},
		[]byte(message),
	); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *SMTPSender) buildMessage(input SendInput) string {
	message := ""

	message += fmt.Sprintf("From: %s\r\n", s.from)
	message += fmt.Sprintf("To: %s\r\n", input.To)
	message += fmt.Sprintf("Subject: %s\r\n", input.Subject)
	message += "MIME-Version: 1.0\r\n"

	if input.HTML != nil {
		message += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
		message += "\r\n"
		message += *input.HTML

		return message
	}

	message += "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
	message += "\r\n"
	message += input.Text

	return message
}
```

---

# Recommended Usage In Auth Layer

Auth should NOT own transport implementation.

Instead:

```go
package auth

import "chia.florist/service-core/internal/shared/mailer"

type RegisterUsecase struct {
	mailer mailer.Sender
}
```

Then inside usecase:

```go
err := u.mailer.Send(mailer.SendInput{
	To: params.Email,
	Subject: "Verify your account",
	Text: fmt.Sprintf("Your OTP is %s", otp),
})
```

---

# Why This Structure Is Better

This keeps:

```text
transport
≠
authentication behavior
```

separated.

Meaning:

* auth decides WHAT to send
* shared mailer decides HOW to send

Very scalable architecture.
