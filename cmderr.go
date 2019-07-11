package kli

import "fmt"

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

type CmdError interface {
	error
	fmt.Stringer
	Code() int
}

type CommandErr struct {
	s string
	c int
}

func NewCmdError(message string, code int) *CommandErr {
	return &CommandErr{
		s: message,
		c: code,
	}
}

func (e CommandErr) Error() string {
	return e.s
}

func (e CommandErr) Code() int {
	return e.c
}

func (e CommandErr) String() string {
	return fmt.Sprintf("error, code: %d, message: %s", e.c, e.s)
}
