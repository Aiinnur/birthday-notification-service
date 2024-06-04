package sendEmail

import (
	"birthday-notification-service/internal/config"
	"birthday-notification-service/internal/entity"
	"context"
	"fmt"
	"log"
	"net/smtp"
)

type BirthdayNotificationRepo interface {
	GetSubscribersForTodayBirthdays(ctx context.Context) ([]entity.SubscriberInfo, error)
}

func SendBirthdayNotifications(cfg *config.Config, repo BirthdayNotificationRepo) {
	ctx := context.Background()
	subscribers, err := repo.GetSubscribersForTodayBirthdays(ctx)
	if err != nil {
		log.Printf("Error fetching subscribers: %v", err)
		return
	}

	for _, info := range subscribers {
		message := fmt.Sprintf("Subject: Поздравь коллегу\n\nСегодня день рождения у коллеги %s:%s. Не забудь отправить поздравительное письмо!", info.BirthdayName, info.BirthdayEmail)
		sendEmail(cfg, info.SubscriberEmail, message)
	}
}

func sendEmail(cfg *config.Config, to, message string) {
	auth := smtp.PlainAuth("", cfg.SmtpServer.User, cfg.SmtpServer.Password, cfg.SmtpServer.Address)
	addr := fmt.Sprintf("%s:%s", cfg.SmtpServer.Address, cfg.SmtpServer.Port)

	if err := smtp.SendMail(addr, auth, cfg.SmtpServer.User, []string{to}, []byte(message)); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
	} else {
		log.Printf("Birthday email sent to %s", to)
	}
}
