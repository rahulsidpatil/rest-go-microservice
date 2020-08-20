package dal

// Interface ... is an interface to data access layer methods
type Interface interface {
	AddMessage(msg *Message) error
	GetMessage(msg *Message) error
	UpdateMessage(msg *Message) error
	DeleteMessage(msg *Message) error
	GetAll() ([]Message, error)
}
