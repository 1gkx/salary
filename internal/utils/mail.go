package utils

import(
	"fmt"
	"net/smtp"

	"github.com/1gkx/salary/internal/conf"
)

func Send(recipient string, password string) error{

	if !conf.Prod() {
		fmt.Printf("Сообщение отправлено на email-адрес %s\n", recipient)
		return nil
	}

	from := "admin@pskb.com" // TODO вынести в конфиг

	// TODO Переделать на шаблон
	msg := "From: " + from + "\n" +
		"Content-Type: text/html; charset=utf-8;\n" +
		"Subject: Регистрация в сервисе 'Зарплатный поект'\n\n" +
		"<html><body><h3>Добрый день!</h3>" +
		"<p>Вы зарегистрировались в сервисе 'Зарплатный проект'!</p>" +
		"Ваш пароль для входа: " + password +
		"</body></html>"

	if err := smtp.SendMail(
		fmt.Sprintf("%s:%s", conf.Cfg.Mail.Host, conf.Cfg.Mail.Port), // Адрес почтового серева
		nil, // Параметры авторизации // TODO вынести в конфиг
		from, // От кого
		[]string{recipient}, // Кому
		[]byte(msg)); // Тело письма
	err != nil {
		return err
	}
	fmt.Printf("Сообщение отправлено на email-адрес %s\n", recipient) // TODO логирование в файл

	return nil
}