package error

type I18Error struct {
	error   error
	message string
}

func (i *I18Error) Error() string {
	return i.message

}
func (i *I18Error) Unwrap() error {
	return i.error
}

func NewI18nError(err error, message string) error {
	return &I18Error{
		error:   err,
		message: message,
	}

}
