package main

import (
	"bufio"
	"fmt"
	"os"
)

func InteractiveShell(globalVariableArray *[]Variable, globalFunctionArray *[]Function, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject) {
	var indicator string
	var isContinue bool = false
	var stringContainer string
	var fdCount int = 0
	var flCount int = 0
	var wlCount int = 0
	var ifCount int = 0

	fmt.Printf("%s %s\n", TITIK_APP_NAME, TITIK_STRING_VERSION)
	fmt.Println("To exit, press ^C")

	for true {
		if !isContinue {
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
		flCount = 0
		wlCount = 0
		ifCount = 0

		if tokenErr != nil {
			fmt.Println(tokenErr)
			stringContainer = ""
		} else {
			for x := 0; x < len(tokenArray); x++ {
				if tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_START {
					fdCount += 1
				}
				if tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_END {
					fdCount -= 1
				}
				if tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_START {
					flCount += 1
				}
				if tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_END {
					flCount -= 1
				}
				if tokenArray[x].Type == TOKEN_TYPE_WHILE_LOOP_START {
					wlCount += 1
				}
				if tokenArray[x].Type == TOKEN_TYPE_WHILE_LOOP_END {
					wlCount -= 1
				}
				if tokenArray[x].Type == TOKEN_TYPE_IF_START {
					ifCount += 1
				}
				if tokenArray[x].Type == TOKEN_TYPE_IF_END {
					ifCount -= 1
				}
			}

			if fdCount > 0 || flCount > 0 || ifCount > 0 || wlCount > 0 {
				isContinue = true
			} else {
				isContinue = false
			}

			if !isContinue {
				var gotReturn bool = false
				var returnToken Token
				var needBreak bool = false
				var stackReference []Token
				var getLastStackBool = false
				var lastStackBool = false

				prsr := Parser{}
				parserErr := prsr.Parse(tokenArray, globalVariableArray, globalFunctionArray, "main", globalNativeVarList, &gotReturn, &returnToken, false, &needBreak, &stackReference, globalSettings, getLastStackBool, &lastStackBool)

				if parserErr != nil {
					fmt.Println(parserErr)
				}

				stringContainer = ""
			}
		}
	}
}
