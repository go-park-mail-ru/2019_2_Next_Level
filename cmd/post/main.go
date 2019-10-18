package main

import (
	"2019_2_Next_Level/cmd/post/outpq"
	"2019_2_Next_Level/cmd/post/smtpd"
)

func main() {
	outpq.Init()
	smtpd.Init()
}
