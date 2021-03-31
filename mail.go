package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"time"
	"strconv"

	"github.com/emersion/go-smtp"
	"gitlab.com/avarf/getenvs"
)

// The Backend implements SMTP server methods.
type Backend struct{}

type user struct {
    username string
    password string
}

// Login handles a login command with username and password.
func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {

	env_username := getenvs.GetEnvString("LOGIN_USERNAME", "john.doe@example.tld")
	env_password := getenvs.GetEnvString("LOGIN_PASSWORD", "s3cr3t")

	if username != env_username || password != env_password {
		return nil, errors.New("Invalid username or password")
	}
	return &Session{}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

// A Session is returned after successful login.
type Session struct{}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		log.Println("Data:", string(b))
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}


func newUser(name, password string) *user {
    u := user{username: name, password: password}
    return &u
}


func main() {
	be := &Backend{}

	s := smtp.NewServer(be)

	ReadTimeout, err := time.ParseDuration(getenvs.GetEnvString("SMTP_READTIMEOUT", "1m30s"))
	if err != nil {
		log.Fatal(err)
	}

	WriteTimeout, err := time.ParseDuration(getenvs.GetEnvString("WRITE_TIMEOUT", "1m30s"))
	if err != nil {
		log.Fatal(err)
	}

	MaxMessageKBytesInt, err := strconv.Atoi(getenvs.GetEnvString("MAX_MESSAGE_KBYTES", "1024"))
	if err != nil {
		log.Fatal(err)
	}
	MaxRecipientsInt, err := strconv.Atoi(getenvs.GetEnvString("MAX_RECIPIENTS", "50"))
	if err != nil {
		log.Fatal(err)
	}

	s.Addr = getenvs.GetEnvString("STMP_PORT", ":1025")
	s.Domain = getenvs.GetEnvString("SMTP_DOMAIN", "localhost")
	s.ReadTimeout = ReadTimeout
	s.WriteTimeout = WriteTimeout
	s.MaxMessageBytes = MaxMessageKBytesInt * 1024
	s.MaxRecipients = MaxRecipientsInt
	s.AllowInsecureAuth = true



	log.Println("Starting server at")
	log.Println("")
	log.Println("Address: ", s.Addr)
	log.Println("ReadTimeout: ", s.ReadTimeout)
	log.Println("WriteTimeout: ", s.WriteTimeout)
	log.Println("MaxMessageBytes: ", s.MaxMessageBytes)
	log.Println("MaxRecipients: ", s.MaxRecipients)
	log.Println("AllowInsecureAuth: ", s.AllowInsecureAuth)

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
