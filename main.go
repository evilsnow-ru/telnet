package main

import (
	"fmt"
	"os"
	"telnet/cmd"
	"telnet/system"
)

func main() {
	//Try to register CTRL+C handler. If it already registered - exit with code 0
	if !system.RegisterSignalHandler() {
		fmt.Println("Can't register signal handler")
		os.Exit(0)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
