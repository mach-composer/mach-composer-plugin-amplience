package internal

import "fmt"

type NoSiteConfigError struct {
	msg string
}

func NewNoSiteConfigError(format string, a ...any) *NoSiteConfigError {
	return &NoSiteConfigError{msg: fmt.Sprintf(format, a...)}
}

func (n *NoSiteConfigError) Error() string {
	return n.msg
}

type InvalidSiteConfigError struct {
	msg string
}

func NewInvalidSiteConfigError(format string, a ...any) *InvalidSiteConfigError {
	return &InvalidSiteConfigError{msg: fmt.Sprintf(format, a...)}
}

func (n *InvalidSiteConfigError) Error() string {
	return n.msg
}
