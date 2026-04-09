package service

import (
	"fmt"
	"net/smtp"
)

// Notifier відповідає за відправку email-повідомлень
type Notifier struct {
	host     string
	port     string
	username string
	password string
}

func NewNotifier(host, port, username, password string) *Notifier {
	return &Notifier{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

// SendReleaseEmail відправляє лист списку підписників
func (n *Notifier) SendReleaseEmail(emails []string, repoName, newTag string) error {
	if len(emails) == 0 {
		return nil
	}

	auth := smtp.PlainAuth("", n.username, n.password, n.host)

	subject := fmt.Sprintf("Subject: Новий реліз у репозиторії %s!\r\n", repoName)
	body := fmt.Sprintf("Привіт!\r\n\r\nУ репозиторії %s щойно вийшов новий реліз: %s.\r\n\r\nПеревір GitHub для деталей!", repoName, newTag)
	msg := []byte(subject + "\r\n" + body)

	addr := n.host + ":" + n.port

	// Відправляємо лист (в реальному продакшені краще відправляти кожному окремо або через bcc)
	err := smtp.SendMail(addr, auth, n.username, emails, msg)
	if err != nil {
		return fmt.Errorf("помилка відправки email: %v", err)
	}

	return nil
}
