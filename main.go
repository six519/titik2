package main

import (
	"fmt"
	"os"
)

func main() {
	var globalVariableArray []Variable
	var globalFunctionArray []Function
	var globalNativeVarList []string
	var globalSettings GlobalSettingsObject

	//init global settings
	globalSettings = GlobalSettingsObject{}
	globalSettings.Init(&globalVariableArray, &globalFunctionArray, &globalNativeVarList)

	//initialize native functions
	initNativeFunctions(&globalFunctionArray)
	//initialize built-in variables
	initBuiltInVariables(&globalVariableArray, &globalNativeVarList)

	if len(os.Args) < 2 {
		Help(os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "-v" {
		Version()
	} else if os.Args[1] == "-h" {
		Help(os.Args[0])
	} else if os.Args[1] == "-i" {
		InteractiveShell(&globalVariableArray, &globalFunctionArray, &globalNativeVarList, &globalSettings)
	} else {
		var gotReturn bool = false
		var returnToken Token
		var needBreak bool = false
		var stackReference []Token
		var getLastStackBool bool = false
		var lastStackBool bool = false
		//open titik file
		lxr := Lexer{FileName: os.Args[1]}
		fileErr := lxr.ReadSourceFile()

		if fileErr != nil {
			fmt.Println(fileErr)
			os.Exit(1)
		}

		//generate token below
		tokenArray, tokenErr := lxr.GenerateToken()
		//DumpToken(tokenArray)
		if tokenErr != nil {
			fmt.Println(tokenErr)
			os.Exit(2)
		}
		//parser object
		prsr := Parser{}
		parserErr := prsr.Parse(tokenArray, &globalVariableArray, &globalFunctionArray, "main", &globalNativeVarList, &gotReturn, &returnToken, false, &needBreak, &stackReference, &globalSettings, getLastStackBool, &lastStackBool)

		if parserErr != nil {
			fmt.Println(parserErr)
			os.Exit(2)
		}

		//DumpToken(tokenArray) //:TEMPORARY
		//DumpVariable(globalVariableArray)
		//DumpFunction(globalFunctionArray)
	}
}
