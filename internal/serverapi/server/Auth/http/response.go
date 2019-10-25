package http

type Answer interface {
	SetStatus(string)
}

type RegisterStruct struct {
	status     string `json:"status"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName"`
	BirthDate  string `json:"birthDate"`
	Sex        string `json:"sex"`
}

func (r *RegisterStruct) SetStatus(status string) {
	r.status = status
}


type OkStruct struct {
	Status string `json:"status"`
}

func (o *OkStruct) SetStatus(status string) {
	o.Status = status
}

