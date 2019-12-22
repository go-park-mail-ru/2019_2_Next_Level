package model

import "2019_2_Next_Level/pkg/HttpTools"

type Folder struct {
	Name string `json:"name"`
	MessageCount int64 `json:"capacity"`
	IsSystem bool `json:"isSystem"`
}

func (folder *Folder) Sanitize() {
	HttpTools.Sanitizer([]*string{
		&folder.Name,
	})
}
