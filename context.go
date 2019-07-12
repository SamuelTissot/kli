package kli

import "os"

// todo create a proper context with timeout
// that has similar functionality as the http.context
// else it's a bit confusing
type Context struct {
	args []string
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) SetArgs(args []string) *Context {
	c.args = args
	return c
}

func (c *Context) Default() *Context {
	//only take the arguments we don't care about the name of the command
	if len(os.Args) > 1 {
		c.args = os.Args[1:]
	}

	return c
}

func (c *Context) Args() []string {
	return c.args
}
