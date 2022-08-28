package sm

// answer is simple answer for user with message and keyboard
type answer struct {
	msg string
	kb  []string
}

// GetMessage return message from answer
func (a answer) GetMessage() string {
	return a.msg
}

// GetKeyboard return message from answer
func (a answer) GetKeyboard() []string {
	return a.kb
}

// newAnswer create new answer
func newAnswer(msg string, kb []string) answer {
	return answer{
		msg: msg,
		kb:  kb,
	}
}
