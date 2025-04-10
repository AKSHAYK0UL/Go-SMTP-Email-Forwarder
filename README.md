
# Go SMTP Email Forwarder

A simple Go-based SMTP email forwarding service that maps custom email aliases to Gmail addresses. This service receives emails on a custom domain and forwards them to a designated Gmail account using Gmail’s SMTP server.

## Overview

This project implements an SMTP server using [go-smtp](https://github.com/emersion/go-smtp) for handling incoming emails and [godotenv](https://github.com/joho/godotenv) for managing environment variables. The core functionality reads the incoming email, verifies that the recipient matches the configured alias, and then forwards the email to a Gmail account via the Gmail SMTP server.

The incoming email address is built using your custom alias and domain, for example:  
`<CUSTOM_EMAIL>@<CUSTOM_DOMAIN>`

Emails are forwarded to:  
`<FORWARD_EMAIL>@gmail.com`

## Features

- **Custom Alias Mapping:** Only emails addressed to your configured alias are processed.
- **Gmail SMTP Forwarding:** Uses Gmail’s secure SMTP server to send the forwarded email.
- **Environment-Based Configuration:** Easily manage credentials and settings using a `.env` file.
- **Basic Test Suite:** Includes a test to verify the email forwarding logic on a local SMTP server.

## Requirements

- [Go](https://golang.org/dl/) 1.13 or higher
- A Gmail account with [App Password](https://support.google.com/accounts/answer/185833) configured
- A valid custom domain to receive emails (e.g. `example.com`)
- Basic SMTP knowledge and access

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/AKSHAYK0UL/Go-SMTP-Email-Forwarder.git
   cd Go-SMTP-Email-Forwarder
   ```

2. **Install dependencies:**

   Ensure your Go modules are enabled. Run:

   ```bash
   go mod tidy
   ```

3. **Prepare your environment:**

   Create a `.env` file in the project root with the following content:

   ```ini
   # Gmail credentials for forwarding
   GMAIL_USER="your_gmail_username"
   GMAIL_APP_PASSWORD="your_gmail_app_password"

   # SMTP server settings
   SMTP_PORT=":2525"
   CUSTOM_DOMAIN="yourcustomdomain.com"   # Can use any custom domain (e.g., xyz.com, abc.com)
   CUSTOM_EMAIL="your_alias"                # The alias you wish to accept mails for
   FORWARD_EMAIL="destination_email_prefix" # The Gmail prefix to forward the mail (without @gmail.com)
   FORWARD_SERVER="smtp.gmail.com:587"
   ```

   Replace the placeholder values with your actual credentials and settings.

## Configuration

The project utilizes environment variables to control the following configurations:

- **GMAIL_USER** and **GMAIL_APP_PASSWORD:** Used for authenticating with Gmail’s SMTP server.
- **SMTP_PORT:** The port on which the custom SMTP server will listen (e.g., `:2525`).
- **CUSTOM_DOMAIN:** Your custom domain used for receiving emails.
- **CUSTOM_EMAIL:** The local part (alias) of the email address.
- **FORWARD_EMAIL:** The Gmail username (prefix) that will receive forwarded emails. The final email becomes `<FORWARD_EMAIL>@gmail.com`.
- **FORWARD_SERVER:** The SMTP server used for sending forwarded emails (`smtp.gmail.com:587` for Gmail).

## Usage

### Running the Forwarder

To run the SMTP forwarder, execute:

```bash
go build .
```
```bash
go run main.go
```

This will start an SMTP server on the port defined by `SMTP_PORT` (e.g., `:2525`). The server accepts incoming emails addressed to `<CUSTOM_EMAIL>@<CUSTOM_DOMAIN>` and forwards them to `<FORWARD_EMAIL>@gmail.com`.

### How It Works

1. **Incoming Email Handling:**  
   The server uses an alias mapping to check if the incoming email’s recipient matches `<CUSTOM_EMAIL>@<CUSTOM_DOMAIN>`. If it does, the email data is accepted.

2. **Email Forwarding:**  
   The accepted email is forwarded using the SMTP credentials provided for Gmail. The forwarding process utilizes `smtp.SendMail` with proper authentication.

## Testing

A basic test case is included to verify email forwarding functionality. It sends an email to the configured test SMTP server (using localhost and your defined port) and checks for errors. To run the tests:

```bash
go test
```

Make sure the necessary environment variables (`CUSTOM_EMAIL` and `CUSTOM_DOMAIN`) are set in your `.env` file before running the test.

## Troubleshooting

- **Missing Environment Variables:**  
  If required environment variables are missing, the application will exit with a corresponding error. Verify your `.env` file exists and is correctly configured.

- **Gmail SMTP Errors:**  
  Make sure that your Gmail account has [App Passwords](https://support.google.com/accounts/answer/185833) enabled and that you use the correct credentials.

- **Port Conflicts:**  
  Ensure that the `SMTP_PORT` specified is not in use by another service on your machine.

## Project Status

**Note:** This project is still under development.
