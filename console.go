//go:build !win
// +build !win

package main

import (
	"errors"
	"fmt"
	"strconv"
)

var ANSI_COLORS = []string{
	"\x1b[0m",    //Normal
	"\x1b[30;1m", //Black
	"\x1b[31;1m", //Red
	"\x1b[32;1m", //Green
	"\x1b[33;1m", //Yellow
	"\x1b[34;1m", //Blue
	"\x1b[35;1m", //Purple
	"\x1b[36;1m", //Cyan
	"\x1b[37;1m", //White
}

func Sc_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		if arguments[0].IntegerValue > (len(ANSI_COLORS)-1) || arguments[0].IntegerValue < 0 {
			*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			fmt.Printf("%s", ANSI_COLORS[arguments[0].IntegerValue])
		}
	}

	return ret
}
