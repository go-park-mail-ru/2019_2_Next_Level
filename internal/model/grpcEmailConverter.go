package model

import (
	"2019_2_Next_Level/internal/post"
	pb "2019_2_Next_Level/internal/post/messagequeue/service"
)

type ParcelAdapter struct {
}

func (a *ParcelAdapter) ToEmail(from *pb.Email) post.Email {
	return post.Email{from.From, from.To, from.Body}
}

func (a *ParcelAdapter) FromEmail(from *post.Email) pb.Email {
	return pb.Email{
		From: from.From,
		To:   from.To,
		Body: from.Body,
	}
}
