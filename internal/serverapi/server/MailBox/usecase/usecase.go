package usecase

import "fmt"

type MailBoxUsecase struct {
}

func (u *MailBoxUsecase) SendMail(from, to, body string) error {
	fmt.Printf("Send mail:\n From: %s\n To: %s\nBody: %s\n", from, to, body)
	return nil
}
