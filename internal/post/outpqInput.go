package post

type Sender interface {
	Put(Email) error
	Get() (Email, error)
	Init()
	Destroy()
}
