package kli

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"
)

// todo review the interface ... it's quite big
type Command interface {
	KFlag

	// Description sets the shot description (except) of the Command
	Description(desc string)

	// Detail sets the Command details, it's the long description
	// like example
	Detail(detail io.Reader)

	// Do sets the function to be called on execution
	Do(fn func(Command, KFlag) CmdError)

	// Parse parses flag definitions from the argument
	// list, which should not include the command name.
	// Must be called after all flags in the
	// FlagSet are defined and before flags are accessed by the program.
	// The return value will be ErrHelp if -help or -h were set but not defined.
	Parse([]string) error

	// Args returns the non-flag arguments.
	Args() []string

	// Name returns the name of the command
	Name() string

	// Execute calls the function that what set by Command::Do
	// returns a CmdError function not found if the
	// execute function was not set
	Execute(Command, KFlag) CmdError

	// IsExecutable returns true if the command executing
	// function has been set
	IsExecutable() bool

	// GetKFlag returns the command Kflag
	GetKFlag() KFlag

	// SetChildren sets the Command Children
	SetChildren(children ...Command) error

	Children() []Command

	// returns the Command's parent Command
	Parent() Command

	// SetParent sets the Command's Parent
	SetParent(parent Command) error

	// PrintDefaults prints, to standard error unless configured otherwise,
	PrintDefaults()

	// Bool sets a flag of type Bool
	Bool(name string, value bool, usage string)

	// Duration sets a flag of type time.Duration (int64)
	Duration(name string, value time.Duration, usage string)

	// Float64 sets a flag of type float64
	Float64(name string, value float64, usage string)

	// Int sets a flag of type Int
	Int(name string, value int, usage string)

	// Int64 sets a flag of type Int64
	Int64(name string, value int64, usage string)

	// String sets a flag of type string
	String(name string, value string, usage string)

	// Uint sets a flag of type Uint
	Uint(name string, value uint, usage string)

	// Uint64 sets a flag of type Uint64
	Uint64(name string, value uint64, usage string)
}

type CMD struct {
	*flag.FlagSet
	KFlag
	desc     string
	detail   io.Reader
	parent   Command
	children []Command
	fn       func(cmd Command, globals KFlag) CmdError
}

// Description sets the command's description
func (c *CMD) Description(desc string) {
	c.desc = desc
}

func (c *CMD) Detail(detail io.Reader) {
	c.detail = detail
}

func NewCommand(name string, handling flag.ErrorHandling) *CMD {
	return &CMD{
		FlagSet: flag.NewFlagSet(name, handling),
		KFlag:   NewKflag(),
	}
}

func NewSubCommand(parent Command, name string, handling flag.ErrorHandling) *CMD {
	return &CMD{
		FlagSet: flag.NewFlagSet(name, handling),
		KFlag:   NewKflag(),
		parent:  parent,
	}
}

func (c *CMD) Do(fn func(Command, KFlag) CmdError) {
	c.fn = fn
}

func (c *CMD) Execute(cmd Command, f KFlag) CmdError {
	if !c.IsExecutable() {
		return NewCmdError("executable function not set", CannotExecute)
	}
	return c.fn(cmd, f)
}

func (c *CMD) IsExecutable() bool {
	return c.fn != nil
}

func (c *CMD) GetKFlag() KFlag {
	return c.KFlag
}

// SetChildren sets the children (sub-command)
// the method also sets the parent of the children command
// as the current command
func (c *CMD) SetChildren(children ...Command) error {
	for _, c := range children {
		err := c.SetParent(c)
		if err != nil {
			return fmt.Errorf("attempting to reset the parent of a child command. %s", err.Error())
		}
	}
	c.children = append(c.children, children...)
	return nil
}

// Children returns the command's children command (it's sub-command)
func (c *CMD) Children() []Command {
	return c.children
}

// setParent
func (c *CMD) SetParent(parent Command) error {
	if c.parent != nil {
		return fmt.Errorf("command %s already has the parent : %s", c.Name(), c.parent.Name())
	}

	c.parent = parent
	return nil
}

// Parent return the command's parent command
func (c *CMD) Parent() Command {
	return c.parent
}

func (c *CMD) PrintDefaults() {
	b := bytes.Buffer{}
	w := bufio.NewWriter(&b)
	c.FlagSet.SetOutput(w)
	padding := " "
	divider := strings.Repeat("-", 78)
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "|%s\n", divider)
	_, _ = fmt.Fprintf(w, "| %s - %s\n", strings.ToUpper(c.Name()), c.desc)
	_, _ = fmt.Fprintf(w, "|%s\n", divider)
	_, _ = fmt.Fprintf(w, "\n%-1sARGS:\n%-1[1]s⎺⎺⎺\n", padding)
	c.FlagSet.PrintDefaults()
	if c.detail != nil {
		_, _ = fmt.Fprintf(w, "\n%-1sUSAGE:\n%-1[1]s⎺⎺⎺\n", padding)
		scanner := bufio.NewScanner(c.detail)
		for scanner.Scan() {
			_, _ = fmt.Fprintf(w, "%-2s%s\n", padding, scanner.Text())
		}
	}
	_ = w.Flush()
	fmt.Println(b.String())

	//print the child default
	for _, child := range c.Children() {
		child.PrintDefaults()
	}
}

func (c *CMD) Bool(name string, value bool, usage string) {
	c.KFlag.SetFlag(name, c.FlagSet.Bool(name, value, usage))
}

func (c *CMD) Duration(name string, value time.Duration, usage string) {
	c.String(name, value.String(), usage)
}

func (c *CMD) Float64(name string, value float64, usage string) {
	c.KFlag.SetFlag(name, c.FlagSet.Float64(name, value, usage))
}

func (c *CMD) Int(name string, value int, usage string) {
	c.KFlag.SetFlag(name, c.FlagSet.Int(name, value, usage))
}

func (c *CMD) Int64(name string, value int64, usage string) {
	c.KFlag.SetFlag(name, c.FlagSet.Int64(name, value, usage))
}

func (c *CMD) String(name string, value string, usage string) {
	c.KFlag.SetFlag(name, c.FlagSet.String(name, value, usage))
}

func (c *CMD) Uint(name string, value uint, usage string) {
	c.KFlag.SetFlag(name, c.FlagSet.Uint(name, value, usage))
}

func (c *CMD) Uint64(name string, value uint64, usage string) {
	c.KFlag.SetFlag(name, c.FlagSet.Uint64(name, value, usage))
}
