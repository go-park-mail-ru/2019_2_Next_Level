package main

import (
	"2019_2_Next_Level/internal/serverapi/server/MailBox/models"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aonemd/margopher"
	"math"
	"os"
	"os/exec"
	"strings"
)

func main() {
	task := flag.String("task", "", "Choose what to generate: users, messages, folders");
	count := flag.Int("c", 0, "How many values to generate");
	user := flag.String("u", "", "Username");
	folder := flag.String("f", "inbox", "Choose folder");
	flag.Parse()
	switch *task {
	case "message":
		AddMessages(*count, *user, *folder, "67a99408-1cd3-11ea-abe3-98fa9b864510")
	}
}

func AddMessages (count int, user, folder string, session string){
	m := margopher.New()
	dir, _ := os.Getwd()
	filePath := dir+"/text.txt"
	for i:=0; i<count; i++ {
		subject := TrimString(m.ReadFile(filePath), 15)
		content := TrimString(m.ReadFile(filePath), 200)
		mail := models.MailToSend{To:[]string{user}, Subject:subject, Content:content}
		req := struct{
			Message models.MailToSend `json:"message"`
		}{mail}
		mailP := mail.ToMain()
		mailP.SetFrom(user)
		js, _ := json.Marshal(req)
		//-H "Content-Type: application/json" -H "X-Login: admin" -H "Origin: http://localhost:3000" --data @sendmail.json http://localhost:3001/api/messages/send
		args := []string{
			"-H", "\"Content-Type: application/json\"",
			"-H", "\"Origin: http://localhost:3000\"",
			"--cookie", "\"session-id="+session+"\"",
			"--data", "\""+strings.Replace(string(js), "\"", "\\\"", -1)+"\"",
			"http://localhost:3001/api/messages/put",
		}
		fmt.Println(strings.Join(args, " "))
		cmd := exec.Command("curl", args...);
		res := cmd.Run();
		fmt.Println(res)
	}
}

func TrimString(s string, length int) string {
	return s[:int(math.Min(float64(length), float64(len(s))))]
}