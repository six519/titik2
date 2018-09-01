package main

var RESERVED_WORDS = []string {
	"fl", //for loop
	"lf", //end of for loop
	"wl", //while loop
	"lw", //end of while loop
	"if", //if
	"fi", //end of if
	"ef", //else if
	"el", //else
	"to", //use in looping statement
	"brk", //break loop
	"fd", //function
	"df", //end of function declaration
	"rtn", //return
	/*
	"zzz", //sleep
	"p", //print
	"i", //include
	"tof", //convert to float
	"toi", //convert to integer
	"tos", //convert to string
	"ex", //exit
	"sc", //set color
	"flcp", //file copy
	"flmv", //file cut
	"flrm", //file remove
	"exe", //execute command
	"r", //readline
	"sac", //system arg count
	"sav", //get system arg values (must be an array)
	"savf", //get first system arg value
	"rndstr", //generate random string
	"rnd", //random
	"genv", //get environment variable
	"senv", //set environment variable
	"as", //get array size
	*/
}

func IsReservedWord(wrd string) bool {

	for x := 0; x < len(RESERVED_WORDS); x++ {
		if(wrd == RESERVED_WORDS[x]) {
			return true
		}
	}

	return false
}