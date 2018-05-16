package main

import (
	"os"
	"fmt"
	"github.com/six519/titik2/info"
	"github.com/six519/titik2/lexer"
	"github.com/six519/titik2/parser"
	"github.com/six519/titik2/variable"
)

func main() {

	if(len(os.Args) < 2) {
		info.Help(os.Args[0])
		os.Exit(1)
	}

	if(os.Args[1] == "-v") {
		info.Version()
	} else if (os.Args[1] == "-h") {
		info.Help(os.Args[0])
	} else if (os.Args[1] == "-i") {
	} else {
		//open titik file
		lxr := lexer.Lexer{FileName: os.Args[1]}
		fileErr := lxr.ReadSourceFile()
		var globalVariableArray []variable.Variable

		if (fileErr != nil) {
			fmt.Println(fileErr)
			os.Exit(1)
		}

		//generate token below
		tokenArray, tokenErr := lxr.GenerateToken()

		if (tokenErr != nil) {
			fmt.Println(tokenErr)
			os.Exit(info.TOKEN_ERROR)
		}
		//parser object
		prsr := parser.Parser{}
		parserErr := prsr.Parse(tokenArray, &globalVariableArray)

		if(parserErr != nil) {
			fmt.Println(parserErr)
			os.Exit(info.SYNTAX_ERROR)
		}

		//lexer.DumpToken(tokenArray) //:TEMPORARY
	}
}