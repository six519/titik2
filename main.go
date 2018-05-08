package main

import (
	"os"
	"github.com/six519/titik2/info"
)

func main() {
	if(len(os.Args) < 2) {
		info.Help(os.Args[0])
		os.Exit(1)
	}
}