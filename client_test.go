package main

import (
	"net/smtp"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func TestEmailForwarding(t *testing.T) {
	testSMTPServer := "localhost:2525"

	customAlias := os.Getenv("CUSTOM_EMAIL")
	customDomain := os.Getenv("CUSTOM_DOMAIN")

	if customAlias == "" || customDomain == "" {
		t.Fatal("CUSTOM_ALIAS and CUSTOM_DOMAIN must be set in your .env file")
	}

	testEmail := customAlias + "@" + customDomain

	err := smtp.SendMail(
		testSMTPServer,
		nil,
		testEmail,
		[]string{testEmail},
		[]byte(strings.Join([]string{
			"To: " + testEmail,
			"Subject: Automated Test",
			"",
			"This is an automated test of the SMTP forwarding system.",
		}, "\r\n")),
	)

	if err != nil {
		t.Fatalf("Failed to send test email: %v", err)
	}
}
