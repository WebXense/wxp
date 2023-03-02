package wxp

type Error interface {
	Error() string
	UUID() string
}

type _error struct {
	uuid string
	msg  string
}

func NewError(uuid string) Error {
	return &_error{
		uuid: uuid,
		msg:  errMap[uuid],
	}
}

func NewCustom(uuid string, msg string) Error {
	return &_error{
		uuid: uuid,
		msg:  msg,
	}
}

func (e *_error) Error() string {
	return e.msg
}

func (e *_error) UUID() string {
	return e.uuid
}

var errMap = make(map[string]string)

func RegisterError(uuid string, msg string) any {
	errMap[uuid] = msg
	return nil
}
