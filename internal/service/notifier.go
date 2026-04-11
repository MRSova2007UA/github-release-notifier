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

// NewNotifier створює новий екземпляр Notifier
func NewNotifier(host, port, username, password string) *Notifier {
	return &Notifier{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

// SendReleaseEmail відправляє лист списку підписників (кожному окремо)
func (n *Notifier) SendReleaseEmail(emails []string, repoName, newTag string) error {
	// Якщо немає кому відправляти - просто виходимо
	if len(emails) == 0 {
		return nil
	}

	// Налаштовуємо авторизацію
	auth := smtp.PlainAuth("", n.username, n.password, n.host)
	addr := n.host + ":" + n.port

	// Відправляємо лист кожному підписнику ОКРЕМО
	for _, email := range emails {
		// 1. Формуємо правильні заголовки (Headers)
		headerFrom := fmt.Sprintf("From: %s\r\n", n.username)
		headerTo := fmt.Sprintf("To: %s\r\n", email)
		headerSubject := fmt.Sprintf("Subject: Новий реліз у репозиторії %s!\r\n", repoName)

		// Вказуємо кодування UTF-8, щоб українські літери відображалися правильно
		headerMIME := "MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"

		// 2. Формуємо тіло листа
		body := fmt.Sprintf("Привіт!\r\n\r\nУ репозиторії %s щойно вийшов новий реліз: %s.\r\n\r\nПеревір GitHub для деталей!", repoName, newTag)

		// 3. Збираємо весь лист до купи
		msg := []byte(headerFrom + headerTo + headerSubject + headerMIME + body)

		// 4. Відправляємо конкретному отримувачу
		err := smtp.SendMail(addr, auth, n.username, []string{email}, msg)
		if err != nil {
			// Якщо помилка для одного email, ми її логуємо, але не зупиняємо цикл
			fmt.Printf("Помилка відправки для %s: %v\n", email, err)
		} else {
			fmt.Printf("Успішно відправлено лист на %s\n", email)
		}
	}

	return nil
}
