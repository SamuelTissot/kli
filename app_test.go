package kli_test

import (
	"flag"
	"fmt"
	"github.com/SamuelTissot/kli"
	"github.com/SamuelTissot/kli/ktest"
	"reflect"
	"strings"
	"testing"
)

func TestApp_ParseSubCommands_withGlobalFlags(t *testing.T) {

	kt := ktest.NewKT()
	kt.Exec(t, func(t *testing.T) {
		// create an app
		root := kli.NewCommand("root", flag.ExitOnError)
		root.String("foo", "", "the echoed string")
		root.String("homer", "", "the echoed string")
		root.Do(func(cmd kli.Command, _ kli.KFlag) kli.Error {
			foo, _ := cmd.StringFlag("foo")
			fmt.Println(foo)
			return nil
		})

		sub := kli.NewCommand("sub", flag.ExitOnError)
		sub.String("str", "", "the echoed string value")
		sub.Do(func(command kli.Command, globals kli.KFlag) kli.Error {
			for n, t := range globals.Store() {
				if t == reflect.String {
					v, _ := globals.StringFlag(n)
					fmt.Printf("%s: %s\n", n, v)
				}
			}

			foo, _ := command.StringFlag("str")
			fmt.Println(foo)
			return nil
		})

		err := root.SetChildren(sub)
		if err != nil {
			panic(err)
		}

		third := kli.NewCommand("third", flag.ExitOnError)
		third.String("cup", "", "the echoed string value")
		third.Do(func(command kli.Command, globals kli.KFlag) kli.Error {
			for n, t := range globals.Store() {
				if t == reflect.String {
					v, _ := globals.StringFlag(n)
					fmt.Printf("%s: %s\n", n, v)
				}
			}
			foo, _ := command.StringFlag("cup")
			fmt.Print("the cup of " + foo)
			return nil
		})

		err = sub.SetChildren(third)
		if err != nil {
			panic(err)
		}

		app := &kli.App{}
		app.SetRoot(root)

		app.Run(kli.NewContext().SetArgs([]string{"--foo", "bar", "-homer", "sub", "sub", "-str", "fuzzfizz", "third", "-cup", "tea"}))

	})

	if kt.Err != nil {
		t.Fatal(kt.Err)
	}

	mustFind := []string{
		"homer: sub",
		"foo: bar",
		"the cup of tea",
	}

	got := string(kt.Out)
	for _, str := range mustFind {
		if !strings.Contains(got, str) {
			t.Errorf("could not find \"%s\" in \n%s", str, got)
		}
	}
}

func TestApp_noArgs(t *testing.T) {

	kt := ktest.NewKT()
	kt.Exec(t, func(t *testing.T) {
		// create an app
		root := kli.NewCommand("root", flag.ExitOnError)
		root.Description("all of your root needs")
		root.String("foo", "", "the echoed string")
		root.String("homer", "", "the echoed string")

		sub := kli.NewCommand("sub", flag.ExitOnError)
		sub.Description("a sub command, yeah!")
		sub.String("str", "", "the echoed string value")
		sub.Do(func(command kli.Command, globals kli.KFlag) kli.Error {
			for n, t := range globals.Store() {
				if t == reflect.String {
					v, _ := globals.StringFlag(n)
					fmt.Printf("%s: %s\n", n, v)
				}
			}

			foo, _ := command.StringFlag("str")
			fmt.Println(foo)
			return nil
		})

		err := root.SetChildren(sub)
		if err != nil {
			panic(err)
		}

		app := &kli.App{}
		app.SetRoot(root)

		app.Run(kli.NewContext().SetArgs([]string{}))
	})

	if kt.Err != nil {
		t.Fatal(kt.Err)
	}

	mustFind := []string{
		"ROOT - all of your root needs",
		"SUB - a sub command, yeah!",
	}

	got := string(kt.Out)
	for _, str := range mustFind {
		if !strings.Contains(got, str) {
			t.Errorf("could not find \"%s\" in \"%s\"", str, got)
		}
	}
}

func TestApp_Help(t *testing.T) {

	kt := ktest.NewKT()
	kt.Exec(t, func(t *testing.T) {
		// create an app
		root := kli.NewCommand("root", flag.ExitOnError)
		root.Description("all of your root needs")
		root.String("foo", "", "the echoed string")
		root.String("homer", "", "the echoed string")
		root.Do(func(cmd kli.Command, _ kli.KFlag) kli.Error {
			foo, _ := cmd.StringFlag("foo")
			fmt.Println(foo)
			return nil
		})

		sub := kli.NewCommand("sub", flag.ExitOnError)
		sub.Description("a sub command, yeah!")
		sub.String("str", "", "the echoed string value")
		sub.Do(func(command kli.Command, globals kli.KFlag) kli.Error {
			for n, t := range globals.Store() {
				if t == reflect.String {
					v, _ := globals.StringFlag(n)
					fmt.Printf("%s: %s\n", n, v)
				}
			}

			foo, _ := command.StringFlag("str")
			fmt.Println(foo)
			return nil
		})

		err := root.SetChildren(sub)
		if err != nil {
			panic(err)
		}

		app := &kli.App{}
		app.SetRoot(root)

		app.Run(kli.NewContext().SetArgs([]string{"sub", "-h"}))
	})

	//todo test for exit code

	mustFind := []string{
		"Usage of sub:",
		"-str string",
	}

	got := kt.ErrOut.String()
	fmt.Println(got)

	for _, str := range mustFind {
		if !strings.Contains(got, str) {
			t.Errorf("could not find \"%s\" in \"%s\"", str, got)
		}
	}
}
