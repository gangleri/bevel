package bevel

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type CmdHandler func([]string, *Flags) error

type SubCommand struct {
	Description string
	Handler     CmdHandler
	Flags       *Flags
}

type Cmd struct {
	SubCmdIndex int
	SubCmds     map[string]*SubCommand
	DisplaySubCommands func()
}

type Flags struct {
	flagSet   *flag.FlagSet
	intFlags  map[string]*int
	strFlags  map[string]*string
	boolFlags map[string]*bool
}

func (f *Flags) Int(name string) int {
	return *f.intFlags[name]
}

func (f *Flags) String(name string) string {
	return *f.strFlags[name]
}

func (f *Flags) Bool(name string) bool {
	return *f.boolFlags[name]
}

var DefaultCmd = &Cmd{
	SubCmdIndex: 1,
	SubCmds: map[string]*SubCommand{},
}

func CommandFunc(cmd, description string, handler CmdHandler, ff ...func(*Flags)) {
	DefaultCmd.CommandFunc(cmd, description, handler, ff...)
}

func PrintCommands() {
	fmt.Println("The commands are:")
	for key, cmd := range DefaultCmd.SubCmds {
		fmt.Println("\t\t" + key + "\t\t" + cmd.Description)
	}
	fmt.Println("Use \"" + filepath.Base(os.Args[0]) + " help <command>\" for more information about a command")
}

func init() {
	DefaultCmd.SubCmds["help"] = &SubCommand{
	Description: "Display help info for a sub command",
		Handler: func(args []string, ff *Flags) error { return HelpForCommand(DefaultCmd, args, ff) },
			Flags: &Flags{
			flagSet: flag.NewFlagSet("help", flag.ExitOnError),
		},
	}
}

func HelpForCommand(cmd *Cmd, args []string, _ *Flags) error {
	hc := "help"
	if len(args) > 0 {
		hc = args[0]
	}
	fmt.Println("Help for command: " + hc)
	fmt.Println(cmd.SubCmds[hc].Description)
	cmd.SubCmds[hc].Flags.flagSet.PrintDefaults()
	return nil
}

func Execute() {
	DefaultCmd.Execute()
}

func (c *Cmd) CommandFunc(cmd, description string, handler CmdHandler, ff ...func(*Flags)) {
	if handler == nil {
		panic("handler cannot be nil")
	}
	c.SubCmds[cmd] = &SubCommand{description, handler, &Flags{}}

	c.SubCmds[cmd].Flags = &Flags{
		flagSet:   flag.NewFlagSet(cmd, flag.ExitOnError),
		intFlags:  map[string]*int{},
		strFlags:  map[string]*string{},
		boolFlags: map[string]*bool{},
	}

	for _, f := range ff {
		f(c.SubCmds[cmd].Flags)
	}
}

func (c *Cmd) Execute() {
	if len(os.Args) < 2 {
		PrintCommands()
		return
	}
	cmd := os.Args[c.SubCmdIndex]
	if _, ok := c.SubCmds[cmd]; !ok {
		PrintCommands()
		return
	}
	c.SubCmds[cmd].Flags.flagSet.Parse(os.Args[2:])
	args := c.SubCmds[cmd].Flags.flagSet.Args()
	err := c.SubCmds[cmd].Handler(args, c.SubCmds[cmd].Flags)
	if err != nil {
		os.Exit(1)
	}
}

func IntFlag(name string, value int, usage string) func(*Flags) {
	return func(ff *Flags) {
		ff.intFlags[name] = ff.flagSet.Int(name, value, usage)
	}
}

func StringFlag(name string, value string, usage string) func(*Flags) {
	return func(ff *Flags) {
		ff.strFlags[name] = ff.flagSet.String(name, value, usage)
	}
}

func BoolFlag(name string, value bool, usage string) func(*Flags) {
	return func(ff *Flags) {
		ff.boolFlags[name] = ff.flagSet.Bool(name, value, usage)
	}
}
