package kli

import (
	"fmt"
	"log"
	"os"
)

type App struct {
	root Command
	seen []Command
}

func (a *App) SetRoot(root Command) {
	a.root = root
}

// Run runs the app with the given argument list
// usually it's the os.Args[1:]
func (a *App) Run(ctx *Context) {

	if len(ctx.Args()) == 0 && !a.root.IsExecutable() {
		// only one argument, call the excute the root fnc right away
		fmt.Println("no arguments, printing default")
		a.root.PrintDefaults()
		os.Exit(0)
		return
	}
	// os.Arg[0] is the path
	// os.Arg[1] is the command (the root) -- we don't care for it's name
	// since we want to be able to rename the command without changing
	// the name of the root command
	// always parse to root element flag since they are the globals
	e := a.root.Parse(ctx.Args())
	if e != nil {
		log.Printf("could not parse arguments: %s", e.Error())
		os.Exit(1)
	}
	args := a.root.Args()
	a.seen = []Command{a.root}

	if len(args) >= 1 {
		a.compute(a.root.Children(), args)
	}
	//args of the first command are the global
	first := a.seen[0]
	// the last command is the one to execute
	last := a.seen[len(a.seen)-1]
	if !last.IsExecutable() {
		log.Printf("command %s does not have an executing method", last.Name())
		last.PrintDefaults()
		os.Exit(1)
	}

	err := last.Execute(last, first.GetKFlag())
	if err != nil {
		log.Println(err.Error())
		os.Exit(err.Code())
	}

	os.Exit(0)
}

func (a *App) compute(cmds []Command, args []string) {
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
