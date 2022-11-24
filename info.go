package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

const TITIK_APP_NAME string = "Titik"
const TITIK_STRING_VERSION string = "3.5.0"
const TITIK_AUTHOR string = "Ferdinand E. Silva"

func Help(exeName string) {
	fmt.Printf("Usage: %s [-options] filename.ttk\n", exeName)
	fmt.Printf("\nwhere options include:\n")
	fmt.Printf("\t-v\tget current version\n")
	fmt.Printf("\t-i\topen interactive shell\n")
	fmt.Printf("\t-h\tprint this usage info\n")
}

func Version() {
	fmt.Printf("%s %s\n", TITIK_APP_NAME, TITIK_STRING_VERSION)
	fmt.Printf("By: %s\n", TITIK_AUTHOR)
	fmt.Printf("Operating System: %s\n", runtime.GOOS)

	if SDL_ENABLED {
		fmt.Println("SDL Support: Yes")
	} else {
		fmt.Println("SDL Support: No")
	}

	if RAYLIB_ENABLED {
		fmt.Println("Raylib Support: Yes")
	} else {
		fmt.Println("Raylib Support: No")
	}
}

func SyntaxErrorMessage(lineNumber int, columnNumber int, description string, fileName string) string {
	var strBuffer bytes.Buffer

	strBuffer.WriteString("Syntax error on line number ")
	strBuffer.WriteString(strconv.Itoa(lineNumber))
	strBuffer.WriteString(" and column number ")
	strBuffer.WriteString(strconv.Itoa(columnNumber))
	strBuffer.WriteString(", Error description: ")
	strBuffer.WriteString(description)
	strBuffer.WriteString(", Filename: ")
	strBuffer.WriteString(fileName)

	return strBuffer.String()
}
