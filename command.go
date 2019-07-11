package kli

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

type Command struct {
	*flag.FlagSet
	f        map[string]*Arg
	desc     string
	usage    io.Reader
	parent   *Command
	children []*Command
	fn       func(cmd *Command, globals map[string]*Arg) CmdError
}

// Description sets the command's description
func (c *Command) Description(desc string) {
	c.desc = desc
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

func (c *Command) PrintDefaults() {
	b := bytes.Buffer{}
	w := bufio.NewWriter(&b)
	c.FlagSet.SetOutput(w)
	padding := " "
	divider := strings.Repeat("-", 78)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "|%s\n", divider)
	fmt.Fprintf(w, "| %s - %s\n", strings.ToUpper(c.Name()), c.desc)
	fmt.Fprintf(w, "|%s\n", divider)
	fmt.Fprintf(w, "\n%-1sARGS:\n%-1[1]s⎺⎺⎺\n", padding)
	c.PrintDefaults()
	if c.usage != nil {
		fmt.Fprintf(w, "\n%-1sUSAGE:\n%-1[1]s⎺⎺⎺\n", padding)
		scanner := bufio.NewScanner(c.usage)
		for scanner.Scan() {
			fmt.Fprintf(w, "%-2s%s\n", padding, scanner.Text())
		}
	}
	w.Flush()
	fmt.Println(b.String())

	//print the child default
	for _, child := range c.Children() {
		child.PrintDefaults()
	}
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
