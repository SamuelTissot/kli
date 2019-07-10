package kli

import "os"

type Context struct {
	args []string
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) SetArgs(args []string) *Context {
	args = append([]string{"CDM_NAME"}, args...)
	c.args = args
	return c
}

func (c *Context) Default() *Context {
	c.args = os.Args[2:]
	return c
}

func (c *Context) Args() []string {
	return c.args
}
