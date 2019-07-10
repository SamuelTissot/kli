package kli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"
)

type Command struct {
	*flag.FlagSet
	f        map[string]*Arg
	parent   *Command
	children []*Command
	fn       func(cmd *Command, globals map[string]*Arg) CmdError
}

func NewCommand(name string, handling flag.ErrorHandling) *Command {
	return &Command{
		FlagSet: flag.NewFlagSet(name, handling),
		f:       map[string]*Arg{},
	}
}

func (c *Command) Execute(fn func(*Command, map[string]*Arg) CmdError) {
	c.fn = fn
}

// SetChildren sets the children (sub-command)
// the method also sets the parent of the children command
// as the current command
func (c *Command) SetChildren(children ...*Command) error {
	for _, c := range children {
		err := c.setParent(c)
		if err != nil {
			return fmt.Errorf("attempting to reset the parent of a child command. %s", err.Error())
		}
	}
	c.children = append(c.children, children...)
	return nil
}

// Children returns the command's children command (it's sub-command)
func (c *Command) Children() []*Command {
	return c.children
}

// setParent
func (c *Command) setParent(parent *Command) error {
	if c.parent != nil {
		return fmt.Errorf("command %s already has the parent : %s", c.Name(), c.parent.Name())
	}

	c.parent = parent
	return nil
}

// Parent return the command's parent command
func (c *Command) Parent() *Command {
	return c.parent
}

// Flag returns the Arg for the given name
func (c *Command) Flag(name string) *Arg {
	if f, ok := c.f[name]; ok {
		return f
	}

	return nil
}

// Bool sets a flag of type Bool
func (c *Command) Bool(name string, value bool, usage string) {
	c.f[name] = &Arg{
		Kind: reflect.Bool,
		v:    c.FlagSet.Bool(name, value, usage),
	}
}

// Duration sets a flag of type time.Duration (int64)
func (c *Command) Duration(name string, value time.Duration, usage string) {
	c.Int64(name, int64(value), usage)
}

// Float64 sets a flag of type float64
func (c *Command) Float64(name string, value float64, usage string) {
	c.f[name] = &Arg{
		Kind: reflect.Float64,
		v:    c.FlagSet.Float64(name, value, usage),
	}
}

// Int sets a flag of type Int
func (c *Command) Int(name string, value int, usage string) {
	c.f[name] = &Arg{
		Kind: reflect.Int,
		v:    c.FlagSet.Int(name, value, usage),
	}
}

// Int64 sets a flag of type Int64
func (c *Command) Int64(name string, value int64, usage string) {
	c.f[name] = &Arg{
		Kind: reflect.Int64,
		v:    c.FlagSet.Int64(name, value, usage),
	}
}

// String sets a flag of type string
func (c *Command) String(name string, value string, usage string) {
	c.f[name] = &Arg{
		Kind: reflect.String,
		v:    c.FlagSet.String(name, value, usage),
	}
}

// Uint sets a flag of type Uint
func (c *Command) Uint(name string, value uint, usage string) {
	c.f[name] = &Arg{
		Kind: reflect.Uint,
		v:    c.FlagSet.Uint(name, value, usage),
	}
}

// Uint64 sets a flag of type Uint64
func (c *Command) Uint64(name string, value uint64, usage string) {
	c.f[name] = &Arg{
		Kind: reflect.Uint64,
		v:    c.FlagSet.Uint64(name, value, usage),
	}
}

type Context struct {
	Args []string
}

func (c *Context) Default() *Context {
	return &Context{
		Args: os.Args[1:],
	}
}
