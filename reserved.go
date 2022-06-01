package main

var RESERVED_WORDS = []string{
	"fl",  //for loop
	"lf",  //end of for loop
	"wl",  //while loop
	"lw",  //end of while loop
	"if",  //if
	"fi",  //end of if
	"ef",  //else if
	"el",  //else
	"to",  //use in looping statement
	"brk", //break loop
	"fd",  //function
	"df",  //end of function declaration
	"rtn", //return
	"fea", //for each
	"aef", //end of for each
}

func IsReservedWord(wrd string) bool {

	for x := 0; x < len(RESERVED_WORDS); x++ {
		if wrd == RESERVED_WORDS[x] {
			return true
		}
	}

	return false
}
