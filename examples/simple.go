package main

import (
	"flag"
	"fmt"
	"github.com/SamuelTissot/kli"
)

func main() {
	// declare the root command
	root := kli.NewCommand("root", flag.ExitOnError)
	root.String("foo", "", "the echoed string")
	root.Execute(func(cmd *kli.Command, _ map[string]*kli.Arg) kli.CmdError {
		foo, _ := cmd.Flag("foo").String()
		fmt.Println(foo)
		return nil
	})

	sub := kli.NewCommand("sub", flag.ExitOnError)
	sub.String("str", "", "the echoed string value")
	sub.Execute(func(command *kli.Command, globals map[string]*kli.Arg) kli.CmdError {
		for n, arg := range globals {
			v, _ := arg.String()
			fmt.Printf("%s: %s\n", n, v)
		}

		foo, _ := command.Flag("str").String()
		fmt.Println(foo)
		return nil
	})

	err := root.SetChildren(sub)
	if err != nil {
		panic(err)
	}

	app := &kli.App{}
	app.SetRoot(root)

	app.Run(kli.NewContext().Default())
}
