package info

import (
	"fmt"
)

func Help(exeName string) {
	fmt.Printf("Usage: %s [-options] filename.ttk\n", exeName)
	fmt.Printf("\nwhere options include:\n")
	fmt.Printf("\t-v\tget current version\n")
	fmt.Printf("\t-i\topen interactive shell\n")
	fmt.Printf("\t-h\tprint this usage info\n")
}