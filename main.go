package main

import (
	"./repl"
	"./runscript"
	"flag"
	"fmt"
	"os"
	"os/user"
)

var (
	scriptFileName = flag.String("f", "",
		"monkey script file name to run (run REPL instead if empty)")
)

func main() {
	flag.Parse()

	if *scriptFileName != "" {
		runScriptFile(*scriptFileName)
	} else {
		runRepl()
	}
}

func runRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

func runScriptFile(fileName string) {
	_, err := runscript.RunScript(fileName)
	if err != nil {
		fmt.Printf("%s", err)
	}
}
