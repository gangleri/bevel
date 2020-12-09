package main

import (
	"fmt"
	"github.com/gangleri/bevel"
)

func handleHelloCmd(args []string, ff *bevel.Flags) error {
	DoSomething(ff.String("first"), ff.String("last"))
	return nil
}

func DoSomething( first, last string) {
	fmt.Print("Hello " + first + " " + last)
}

func handleInfoCmd(_ []string, _*bevel.Flags) error {
	fmt.Println("info")
	return nil
}

func main() {
	bevel.CommandFunc("hello", "The hello command says hello",
		handleHelloCmd,
		bevel.StringFlag("first", "", "First name. (Required)"),
		bevel.StringFlag("last", "", "Last name. (Required)"),
		bevel.BoolFlag("upper", false, "Output in uppercase."))

	bevel.CommandFunc("info", "Display some info", handleInfoCmd)

	bevel.Execute()
}

// example #1
// example hello -name Alan -upper
//output: HELLO ALAN

// example #2
// example hello -name Alan
//output: Hello Alan

// example #3
// example info
// output: info
