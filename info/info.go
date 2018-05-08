package info

import (
	"fmt"
	"runtime"
)

const TITIK_APP_NAME string = "Titik"
const TITIK_STRING_VERSION string = "2.0.0"
const TITIK_AUTHOR string = "Ferdinand E. Silva"

func Help(exeName string) {
	fmt.Printf("Usage: %s [-options] filename.ttk\n", exeName)
	fmt.Printf("\nwhere options include:\n")
	fmt.Printf("\t-v\tget current version\n")
	fmt.Printf("\t-i\topen interactive shell\n")
	fmt.Printf("\t-h\tprint this usage info\n")
}

func Version() {
    fmt.Printf("%s %s\n", TITIK_APP_NAME, TITIK_STRING_VERSION);
	fmt.Printf("By: %s\n", TITIK_AUTHOR);
	fmt.Printf("Operating System: %s\n", runtime.GOOS);
}