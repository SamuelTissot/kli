package kli_test

import (
	"flag"
	"fmt"
	"github.com/SamuelTissot/kli"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestApp_Run(t *testing.T) {

	if os.Getenv("CRASH_PLAN") == "1" {
		// create an app
		root := kli.NewCommand("root", flag.ExitOnError)
		root.String("foo", "", "the echoed string")
		root.String("homer", "", "the echoed string")
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

		third := kli.NewCommand("third", flag.ExitOnError)
		third.String("cup", "", "the echoed string value")
		third.Execute(func(command *kli.Command, globals map[string]*kli.Arg) kli.CmdError {
			for n, arg := range globals {
				v, _ := arg.String()
				fmt.Printf("%s: %s\n", n, v)
			}
			foo, _ := command.Flag("cup").String()
			fmt.Print("the cup of " + foo)
			return nil
		})

		err = sub.SetChildren(third)
		if err != nil {
			panic(err)
		}

		app := &kli.App{}
		app.SetRoot(root)

		app.Run(&kli.Context{Args: []string{"root", "--foo", "bar", "-homer", "sub", "sub", "-str", "fuzzfizz", "third", "-cup", "tea"}})
	}

	tcmd := exec.Command(os.Args[0], "-test.run=TestApp_Run")
	tcmd.Env = append(os.Environ(), "CRASH_PLAN=1")
	out, err := tcmd.Output()
	if err != nil {
		t.Error(err)
	}

	mustFind := []string{
		"homer: sub",
		"foo: bar",
		"the cup of tea",
	}

	got := string(out)
	for _, str := range mustFind {
		if !strings.Contains(got, str) {
			t.Errorf("could not find \"%s\" in \"%s\"", str, got)
		}
	}
}
