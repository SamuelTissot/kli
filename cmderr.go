package kli

import "fmt"

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
