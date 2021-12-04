package goVirtualHost

type errorWrapper struct {
	parent  error
	message string
}

func (wrapper errorWrapper) Error() string {
	return wrapper.message
}

func (wrapper errorWrapper) Unwrap() error {
	return wrapper.parent
}

func wrapError(parent error, message string) errorWrapper {
	return errorWrapper{parent, message}
}
