package main

import (
	"os"
	"fmt"
)

func main() {

	if(len(os.Args) < 2) {
		Help(os.Args[0])
		os.Exit(1)
	}

	if(os.Args[1] == "-v") {
		Version()
	} else if (os.Args[1] == "-h") {
		Help(os.Args[0])
	} else if (os.Args[1] == "-i") {
	} else {
		//open titik file
		lxr := Lexer{FileName: os.Args[1]}
		fileErr := lxr.ReadSourceFile()
		var globalVariableArray []Variable
		var globalFunctionArray []Function

		//initialize native functions
		initNativeFunctions()

		if (fileErr != nil) {
			fmt.Println(fileErr)
			os.Exit(1)
		}

		//generate token below
		tokenArray, tokenErr := lxr.GenerateToken()

		if (tokenErr != nil) {
			fmt.Println(tokenErr)
			os.Exit(2)
		}
		//parser object
		prsr := Parser{}
		parserErr := prsr.Parse(tokenArray, &globalVariableArray, &globalFunctionArray, "main")

		if(parserErr != nil) {
			fmt.Println(parserErr)
			os.Exit(2)
		}

		//DumpToken(tokenArray) //:TEMPORARY
		DumpVariable(globalVariableArray)
	}
}