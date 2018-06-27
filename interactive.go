package main

import (
	"fmt"
	"bufio"
	"os"
)

func InteractiveShell(globalVariableArray *[]Variable, globalFunctionArray *[]Function) {
	var indicator string
	var isContinue bool = false

	fmt.Printf("%s %s\n", TITIK_APP_NAME, TITIK_STRING_VERSION)
	fmt.Println("To exit, press ^C")

	for true {
		if(!isContinue) {
			indicator = ">>>"
		} else {
			indicator = "..."
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s ", indicator)
		text, _ := reader.ReadString('\n')
		lxr := Lexer{}
		lxr.ReadString(text)
	}
}