package main

import (
	"fmt"
	"bufio"
	"os"
)

func InteractiveShell(globalVariableArray *[]Variable, globalFunctionArray *[]Function, globalNativeVarList *[]string) {
	var indicator string
	var isContinue bool = false
	var stringContainer string
	var fdCount int = 0

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
		lxr := Lexer{FileName: "interactive_shell"}
		stringContainer += text
		lxr.ReadString(stringContainer)

		tokenArray, tokenErr := lxr.GenerateToken()
		fdCount = 0

		if (tokenErr != nil) {
			fmt.Println(tokenErr)
		} else {
			for x := 0; x < len(tokenArray); x++ {
				if(tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_START) {
					fdCount += 1
				}
				if(tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_END) {
					fdCount -= 1
				}
			}

			if(fdCount > 0) {
				isContinue = true
			} else {
				isContinue = false
			}

			if(!isContinue) {
				prsr := Parser{}
				parserErr := prsr.Parse(tokenArray, globalVariableArray, globalFunctionArray, "main", globalNativeVarList)
		
				if(parserErr != nil) {
					fmt.Println(parserErr)
				}

				stringContainer = ""
			}
		}
	}
}