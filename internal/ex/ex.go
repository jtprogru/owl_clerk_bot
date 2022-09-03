package ex

import "context"

type SM struct{}

type answer struct {
	msg     string
	buttons []string
}

func (a answer) GetMessage() string {
	return ""
}
func (a answer) GetKeyboard() []string {
	return nil
}

type Profile interface {
	GetUID() int64
	GetFirstName() string
	GetLastName() string
	GetUsername() string
}

type Message interface {
	GetMessage() string
}

func (s SM) SaveOrUpdateState(ctx context.Context, p Profile, m Message) (answer, error) {
	return answer{}, nil
}
