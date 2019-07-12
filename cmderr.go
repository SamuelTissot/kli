package kli

import (
	"fmt"
	errw "github.com/pkg/errors"
)

// unix command error code as per
// https://www.tldp.org/LDP/abs/html/exitcodes.html
const (
	OK              = 0   //ok
	GeneralError    = 1   //Catchall for general errors	let "var1 = 1/0"	Miscellaneous errors, such as "divide by zero" and other impermissible operations
	MisuseError     = 2   //Misuse of shell builtins (according to Bash documentation)	empty_function() {}	Missing keyword or command, or permission problem (and diff return code on a failed binary file comparison).
	CannotExecute   = 126 //Command invoked cannot execute	/dev/null	Permission problem or command is not an executable
	NotFount        = 127 //"command not found"	illegal_command	Possible problem with $PATH or a typo
	InvalidArgument = 128 //Invalid argument to exit	exit 3.14159	exit takes only integer args in the range 0 - 255 (see first footnote)
	UserTermination = 130 //Script terminated by Control-C	Ctl-C	Control-C is fatal error signal 2, (130 = 128 + 2, see above)
	OutOfRange      = 255 //*	Exit status out of range	exit -1	exit takes only integer args in the range 0 - 255
)

type Error interface {
	error
	Code() int
}

type KError struct {
	e error
	c int
}

func ErrorWrap(err error, text string, code int) *KError {
	return &KError{
		e: errw.Wrap(err, text),
		c: code,
	}
}

func NewError(text string, code int) *KError {
	return NewErrorf(code, text)
}

func NewErrorf(code int, format string, a ...interface{}) *KError {
	return &KError{
		e: fmt.Errorf(format, a...),
		c: code,
	}
}

func (ke KError) Error() string {
	return ke.e.Error()
}

func (ke KError) Code() int {
	return ke.c
}
