package kli

import (
	"log"
	"os"
)

type App struct {
	root *Command
	seen []*Command
}

func (a *App) SetRoot(root *Command) {
	a.root = root
}

// Run runs the app with the given argument list
// usually it's the os.Args[1:]
func (a *App) Run(ctx *Context) {
	if len(ctx.Args) == 0 {
		log.Println("no argument list")
		os.Exit(1)
	}

	if len(ctx.Args) == 1 {
		// only one argument, call the excute the root fnc right away
		if err := a.root.fn(a.root, map[string]*Arg{}); err != nil {
			log.Println(err.Error())
			os.Exit(err.Code())
		}
		return
	}
	// os.Arg[0] is the path
	// os.Arg[1] is the command (the root) -- we don't care for it's name
	// since we want to be able to rename the command without changing
	// the name of the root command
	// always parse to root element flag since they are the globals
	e := a.root.Parse(ctx.Args[1:])
	if e != nil {
		log.Printf("could not parse arguments: %s", e.Error())
		os.Exit(1)
	}
	args := a.root.Args()
	a.seen = []*Command{a.root}

	if len(args) > 2 {
		a.compute(a.root.Children(), args)
	}
	//args of the first command are the global
	first := a.seen[0]
	// the last command is the one to execute
	last := a.seen[len(a.seen)-1]
	if last.fn == nil {
		log.Printf("command %s does not have an executing method", last.Name())
		os.Exit(1)
	}

	err := last.fn(last, first.f)
	if err != nil {
		log.Println(err.Error())
		os.Exit(err.Code())
	}

	os.Exit(0)
}

func (a *App) compute(cmds []*Command, args []string) {
	if len(args) < 1 {
		return
	}

	// pop-front
	// todo maybe it would be more performance to pop
	arg, args := args[0], args[1:]

	for _, c := range cmds {
		if arg == c.Name() {
			//parse arguments
			err := c.Parse(args)
			if err != nil {
				log.Printf("could not parse argument: %s", err.Error())
			}
			args = c.Args()
			//add to seen
			a.seen = append(a.seen, c)
			a.compute(c.Children(), args)
			return
		}
	}

	a.compute(cmds, args)
}