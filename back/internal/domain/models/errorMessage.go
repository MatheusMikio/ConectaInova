package models

type ErrorMessage struct {
	Property string
	Message  string
}

func NewErrorMessage(property, message string) *ErrorMessage {
	return &ErrorMessage{Property: property, Message: message}
}

func (e *ErrorMessage) Error() string {
	return e.Message
}
