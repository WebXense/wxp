package errs

func Register(uuid string, msg string) {
	errMap[uuid] = msg
}

func New(uuid string) Error {
	return &_error{
		uuid: uuid,
		msg:  errMap[uuid],
	}
}

type Error interface {
	Error() string
	UUID() string
}

type _error struct {
	uuid string
	msg  string
}

func (e *_error) Error() string {
	return e.msg
}

func (e *_error) UUID() string {
	return e.uuid
}

var errMap = make(map[string]string)
