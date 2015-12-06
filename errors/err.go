package errors

type Error struct {
	Message string
	Origin  error
}

func (e *Error) Error() string {
	if e.Origin == nil {
		return e.Message
	}
	return e.Message + ": " + e.Origin.Error()
}

func New(msg string) error { return &Error{Message: msg} }
func Wrap(msg string, origin error) error {
	if origin == nil {
		return nil
	}
	return &Error{Message: msg, Origin: origin}
}
