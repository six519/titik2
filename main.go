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

	if(os.Args[1] == "-v") {
		info.Version()
	} else if (os.Args[1] == "-h") {
		info.Help(os.Args[0])
	} else if (os.Args[1] == "-i") {
	} else {
		//open titik file
	}
}