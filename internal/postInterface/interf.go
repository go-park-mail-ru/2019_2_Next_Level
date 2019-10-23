package postinterface

import (
	"2019_2_Next_Level/internal/post"
)

type IPostInterface interface {
	Init()
	Destroy()
	Put(post.Email) error
	Get() (post.Email, error)
}
