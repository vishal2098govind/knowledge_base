#go #smtp

Example of simple SMTP Email
```
MIME-Version: 1.0 
Date: Sun, 22 Jan 2023 11:54:48 -0500 
To: jon@calhoun.io 
From: test@lenslocked.com 
Subject: This is a test email 
Content-Transfer-Encoding: quoted-printable 
Content-Type: text/plain; charset=UTF-8 

This is the body of the email
```

Example of complex SMTP Email
This includes both the plain text version of the body and HTML version of the body, separated by the boundary
Here, the content-type is `multipart/alternative` which indicates that there multiple parts of the body 
```txt
MIME-Version: 1.0 
Date: Sun, 22 Jan 2023 11:55:59 -0500 
Subject: This is a test email 
To: jon@calhoun.io 
From: test@lenslocked.com 
Content-Type: multipart/alternative; 
boundary=2df315b0bd754b2cea495f617b327626853ede3bcbf4608725384a95937f

--2df315b0bd754b2cea495f617b327626853ede3bcbf4608725384a95937f 
Content-Transfer-Encoding: quoted-printable 
Content-Type: text/plain; charset=UTF-8 

This is the body of the email 
--2df315b0bd754b2cea495f617b327626853ede3bcbf4608725384a95937f 
Content-Transfer-Encoding: quoted-printable 
Content-Type: text/html; charset=UTF-8

<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p> 
--2df315b0bd754b2cea495f617b327626853ede3bcbf4608725384a95937f--
```
Over time, we can end up having more and more part and this could get a little bit more complicated for instance, if we wanted to add attachments to emails. Thus, we use a library that is well tested and that works so that we can construct our emails and get back to focusing on our application.

Go has an SMTP package in the standard library, but it's pretty limited in what it does. Thus, we can use `github.com/go-mail/mail`.

### Building Mail (SMTP Request Body)

```go
package main

import (
	"os"
	
	"github.com/go-mail/mail/v2"
)

func main() {
	from := "test@lenslocked.com"
	to := "jon@calhoun.io"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`

	msg := mail.NewMessage()

	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)
}

// OUTPUT:
// MIME-Version: 1.0
// Date: Wed, 23 Jul 2025 20:04:49 +0530
// To: jon@calhoun.io
// From: test@lenslocked.com
// Subject: This is a test email
// Content-Type: multipart/alternative;
// boundary=15cab7f8525e116da166ae6ed368e4b346a647edca41ea67a85d9430ea05

// --15cab7f8525e116da166ae6ed368e4b346a647edca41ea67a85d9430ea05
// Content-Transfer-Encoding: quoted-printable
// Content-Type: text/plain; charset=UTF-8

// This is the body of the email
// --15cab7f8525e116da166ae6ed368e4b346a647edca41ea67a85d9430ea05
// Content-Transfer-Encoding: quoted-printable
// Content-Type: text/html; charset=UTF-8

// <h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>
// --15cab7f8525e116da166ae6ed368e4b346a647edca41ea67a85d9430ea05--
```

### Connecting to SMTP Server and Sending an Email
- We need the host, port, username and password of the SMTP Server we use (like the `Mailtrap`)
- To connect to a SMTP server, we use `mail.Dialer`
- `mail.Dialer.Dial()` authenticates an SMTP server and gives a way of sending or closing emails
```go
package main

import (
	"os"

	"github.com/go-mail/mail"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "dc214a83099356"
	password = "56deedf1bb4769"
)

func main() {
	from := "test@lenslocked.com"
	to := "jon@calhoun.io"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`
	
	msg := mail.NewMessage()
	
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)
	
	msg.WriteTo(os.Stdout)
	
	dialer := mail.NewDialer(host, port, username, password)
	senderCloser, err := dialer.Dial()
	if err != nil {
		panic(err)
	}
	defer senderCloser.Close()
	
	// to send multiple emails without having to re-dial
	senderCloser.Send(from, []string{to}, msg)
	senderCloser.Send(from, []string{to}, msg)
	senderCloser.Send(from, []string{to}, msg)
	
	// to send one mail
	err := dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message sent")
}
```