package main

import (
	"flag"
	"fmt"
	"github.com/SamuelTissot/kli"
	"strings"
)

func main() {
	// declare the root command
	root := kli.NewCommand("cow", flag.ExitOnError)
	root.Bool("eat", false, "informs the cow to eat")
	root.Do(func(cmd kli.Command, _ kli.KFlag) kli.Error {
		isEating, _ := cmd.BoolFlag("eat")
		if isEating {
			fmt.Println(strings.Repeat("munch ", 3))
		} else {
			fmt.Println("the cow stands there looking smug")
		}
		return nil
	})

	//declare the sub command
	sub := kli.NewCommand("say", flag.ExitOnError)
	sub.String("what", "mooooo", "what the cow will say")
	sub.Int("repeat", 1, "how many time it repeats the word")
	sub.Do(func(cmd kli.Command, globals kli.KFlag) kli.Error {
		if isEating, ok := globals.BoolFlag("eat"); ok {
			if isEating {
				fmt.Println("munch... can't say anything, I'm eating")
				return nil
			}
		}

		what, _ := cmd.StringFlag("what")
		repeat, _ := cmd.IntFlag("repeat")
		for i := 1; i <= repeat; i++ {
			fmt.Println(what)
		}
		return nil
	})

	//add the subcommand to the root command
	err := root.SetChildren(sub)
	if err != nil {
		panic(err)
	}

	//create the app with the root command
	app := &kli.App{}
	app.SetRoot(root)

	//run the app with the context default
	// the context default are the os.Args
	app.Run(kli.NewContext().Default())
}
