package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/may20xx/booking/config"
	"github.com/may20xx/booking/pkg/log"
)

type Mail interface {
	SendMailConfirmAccount(to string, token string) error
}

type mail struct {
	Host       string
	Port       int
	User       string
	Pass       string
	ServerHost string
}

func NewMailService() *mail {
	setting := config.GetConfig()

	return &mail{
		Host:       setting.MailHost,
		Port:       setting.MailPort,
		User:       setting.MailUser,
		Pass:       setting.MailPass,
		ServerHost: "http://localhost:" + setting.Port,
	}
}

func (m *mail) sendMail(to string, subject string, body string) error {
	from := m.User
	password := m.Pass

	auth := smtp.PlainAuth("", from, password, m.Host)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject = "Subject: " + subject + "\n"
	msg := []byte(subject + mime + "\n" + body)

	smtpServer := fmt.Sprintf("%s:%d", m.Host, m.Port)

	err := smtp.SendMail(smtpServer, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (m *mail) SendMailConfirmAccount(to string, token string) error {
	_, currentFile, _, _ := runtime.Caller(0)

	rootDir := filepath.Join(filepath.Dir(currentFile), "..", "..")

	templatePath := filepath.Join(rootDir, "templates", "confirm_account.html")

	content, err := os.ReadFile(templatePath)
	if err != nil {
		log.Msg.Errorf("failed to read email template: %v", err)
		return fmt.Errorf("failed to read email template: %w", err)
	}

	htmlContent := string(content)

	confirmURL := m.ServerHost + "/api/v1/auth/confirm-account?token=" + token

	replacedContent := strings.Replace(htmlContent, "{{ ConfirmationURL }}", confirmURL, -1)

	subject := "Confirm Your Account"

	err = m.sendMail(to, subject, replacedContent)
	if err != nil {
		log.Msg.Errorf("failed to send confirmation email: %v", err)
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	log.Msg.Infof("Confirmation email sent successfully to %s", to)

	return nil
}
