# kli

An Expressive command line interface. 

#####features
- No creation of multiple flag pointers
- Easy child/parent relashionships (sub-commands)
- Support for global flags
- easy to test
- expressive command declaration 

### Usage

```go
root := kli.NewCommand("cow", flag.ExitOnError)
root.Bool("eat", false, "informs the cow to eat")
root.Execute(func(cmd *kli.Command, _ map[string]*kli.Arg) kli.CmdError {
    isEating, _ := cmd.Flag("eat").Bool()
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
sub.Execute(func(command *kli.Command, globals map[string]*kli.Arg) kli.CmdError {
    if eatFlg, ok := globals["eat"]; ok {
        isEating, _ := eatFlg.Bool()
        if isEating {
            fmt.Println("munch... can't say anything, I'm eating")
            return nil
        }
    }

    what, _ := command.Flag("what").String()
    repeat, _ := command.Flag("repeat").Int()
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

```