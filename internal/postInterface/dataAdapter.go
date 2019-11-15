package postinterface

import (
	"2019_2_Next_Level/internal/post"
	pb "2019_2_Next_Level/internal/post/MessageQueue/service"
)

type ParcelAdapter struct {
}

func (a *ParcelAdapter) ToEmail(from *pb.Email) post.Email {
	return post.Email{From:from.From, To:from.To, Subject:from.Subject, Body:from.Body}
}

func (a *ParcelAdapter) FromEmail(from *post.Email) pb.Email {
	return pb.Email{
		From: from.From,
		To:   from.To,
		Body: from.Body,
		Subject:from.Subject,
	}
}
