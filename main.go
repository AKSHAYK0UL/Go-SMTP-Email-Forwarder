package main

import (
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"strings"

	smtpserver "github.com/emersion/go-smtp"
	"github.com/joho/godotenv"
)

var aliasMapping map[string]string

var (
	gmailUser     string
	gmailAppPass  string
	smtpPort      string
	customDomain  string // Use this as your custom domain for incoming emails.
	customEmail   string // The alias/local part for the custom email.
	forwardEmail  string //real address
	forwardServer string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using system environment variables")
	}

	gmailUser = os.Getenv("GMAIL_USER")
	gmailAppPass = os.Getenv("GMAIL_APP_PASSWORD")
	smtpPort = os.Getenv("SMTP_PORT")
	customDomain = os.Getenv("CUSTOM_DOMAIN")
	customEmail = os.Getenv("CUSTOM_EMAIL")
	forwardEmail = os.Getenv("FORWARD_EMAIL")
	forwardServer = os.Getenv("FORWARD_SERVER")

	if gmailUser == "" || gmailAppPass == "" || smtpPort == "" || customDomain == "" || customEmail == "" || forwardEmail == "" || forwardServer == "" {
		log.Fatalln("Missing one or more required environment variables. Please check your .env file.")
	}

	// The incoming email is constructed as: "<CUSTOM_ALIAS>@<CUSTOM_DOMAIN>"
	// And it will be forwarded to: "<CUSTOM_ALIAS>@gmail.com"
	incomingAddress := fmt.Sprintf("%s@%s", customEmail, customDomain)
	forwardAddress := fmt.Sprintf("%s@gmail.com", forwardEmail)

	aliasMapping = map[string]string{
		incomingAddress: forwardAddress,
	}

	log.Printf("Alias mapping loaded: %v", aliasMapping)
}

type Backend struct{}

func (be *Backend) NewSession(c *smtpserver.Conn) (smtpserver.Session, error) {
	return &Session{}, nil
}

type Session struct {
	from       string
	recipients []string
	data       []byte
}

func (s *Session) Mail(from string, opts *smtpserver.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtpserver.RcptOptions) error {
	toLower := strings.ToLower(to)
	if _, exists := aliasMapping[toLower]; !exists {
		return smtpserver.ErrAuthRequired
	}
	s.recipients = append(s.recipients, toLower)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	s.data = data

	for _, alias := range s.recipients {
		dest := aliasMapping[alias]
		err = smtp.SendMail(
			forwardServer,
			smtp.PlainAuth("", gmailUser, gmailAppPass, "smtp.gmail.com"),
			gmailUser,
			[]string{dest},
			s.data,
		)

		if err != nil {
			log.Printf("Failed to forward %s -> %s: %v", alias, dest, err)
		} else {
			log.Printf("Forwarded %s -> %s", alias, dest)
		}
	}
	return nil
}

func (s *Session) Reset()        {}
func (s *Session) Logout() error { return nil }

func main() {
	be := &Backend{}
	srv := smtpserver.NewServer(be)

	srv.Addr = smtpPort
	srv.Domain = customDomain
	srv.AllowInsecureAuth = true

	log.Printf("SMTP Forwarder running on %s (Domain: %s)", srv.Addr, srv.Domain)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
