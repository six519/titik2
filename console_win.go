//go:build win
// +build win

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"syscall"

	"golang.org/x/crypto/ssh/terminal"
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

const (
	FOREGROUND_BLACK  uint16 = 0x0000
	FOREGROUND_BLUE   uint16 = 0x0001
	FOREGROUND_GREEN  uint16 = 0x0002
	FOREGROUND_CYAN   uint16 = 0x0003
	FOREGROUND_RED    uint16 = 0x0004
	FOREGROUND_PURPLE uint16 = 0x0005
	FOREGROUND_YELLOW uint16 = 0x0006
	FOREGROUND_WHITE  uint16 = 0x0007
)

type (
	SHORT int16
	WORD  uint16

	SMALL_RECT struct {
		Left   SHORT
		Top    SHORT
		Right  SHORT
		Bottom SHORT
	}

	COORD struct {
		X SHORT
		Y SHORT
	}

	CONSOLE_SCREEN_BUFFER_INFO struct {
		Size              COORD
		CursorPosition    COORD
		Attributes        WORD
		Window            SMALL_RECT
		MaximumWindowSize COORD
	}
)

func P_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if arguments[0].Type == ARG_TYPE_FLOAT {
		fmt.Printf("%f", arguments[0].FloatValue)
		ret.StringValue = strconv.FormatFloat(arguments[0].FloatValue, 'f', -1, 64)
	} else if arguments[0].Type == ARG_TYPE_STRING {
		fmt.Printf("%s", escapeString(arguments[0].StringValue))
		ret.StringValue = arguments[0].StringValue
	} else if arguments[0].Type == ARG_TYPE_INTEGER {
		//integer
		fmt.Printf("%d", arguments[0].IntegerValue)
		ret.StringValue = strconv.Itoa(arguments[0].IntegerValue)
	} else if arguments[0].Type == ARG_TYPE_BOOLEAN {
		//boolean
		fmt.Printf("%v", arguments[0].BooleanValue)
		if arguments[0].BooleanValue {
			ret.StringValue = "true"
		} else {
			ret.StringValue = "false"
		}
	} else if arguments[0].Type == ARG_TYPE_ASSOCIATIVE_ARRAY {
		strVal := ""
		x := 0

		for k, v := range arguments[0].AssociativeArrayValue {

			if v.Type == ARG_TYPE_FLOAT {
				strVal = strVal + k + ":" + strconv.FormatFloat(v.FloatValue, 'f', -1, 64)
			} else if v.Type == ARG_TYPE_STRING {
				strVal = strVal + k + ":" + v.StringValue
			} else if v.Type == ARG_TYPE_INTEGER {
				strVal = strVal + k + ":" + strconv.Itoa(v.IntegerValue)
			} else if v.Type == ARG_TYPE_BOOLEAN {
				if v.BooleanValue {
					strVal = strVal + k + ":" + "true"
				} else {
					strVal = strVal + k + ":" + "false"
				}
			} else {
				strVal = strVal + k + ":" + "Nil"
			}

			if (x + 1) != len(arguments[0].AssociativeArrayValue) {
				strVal = strVal + ", "
			}

			x += 1
		}

		fmt.Printf("{%s}", strVal)
	} else if arguments[0].Type == ARG_TYPE_ARRAY {
		strVal := ""

		for x := 0; x < len(arguments[0].ArrayValue); x++ {
			if arguments[0].ArrayValue[x].Type == ARG_TYPE_FLOAT {
				strVal = strVal + strconv.FormatFloat(arguments[0].ArrayValue[x].FloatValue, 'f', -1, 64)
			} else if arguments[0].ArrayValue[x].Type == ARG_TYPE_STRING {
				strVal = strVal + arguments[0].ArrayValue[x].StringValue
			} else if arguments[0].ArrayValue[x].Type == ARG_TYPE_INTEGER {
				strVal = strVal + strconv.Itoa(arguments[0].ArrayValue[x].IntegerValue)
			} else if arguments[0].ArrayValue[x].Type == ARG_TYPE_BOOLEAN {
				if arguments[0].ArrayValue[x].BooleanValue {
					strVal = strVal + "true"
				} else {
					strVal = strVal + "false"
				}
			} else {
				strVal = strVal + "Nil"
			}

			if (x + 1) != len(arguments[0].ArrayValue) {
				strVal = strVal + ", "
			}
		}

		fmt.Printf("[%s]", strVal)
	} else {
		//Nil
		fmt.Println("Nil")
	}

	return ret
}

func R_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s", escapeString(arguments[0].StringValue))
		text, _ := reader.ReadString('\n')
		ret.StringValue = strings.Trim(text, "\n")
	}

	return ret
}

func Rp_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		fmt.Printf("%s", escapeString(arguments[0].StringValue))
		text, err := terminal.ReadPassword(0)
		if err == nil {
			ret.StringValue = strings.Trim(string(text), "\n")
		}
	}

	return ret
}

func Sc_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		if arguments[0].IntegerValue > (len(ANSI_COLORS)-1) || arguments[0].IntegerValue < 0 {
			*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {

			kernel32 := syscall.NewLazyDLL("kernel32.dll")
			setConsoleTextAttributeProc := kernel32.NewProc("SetConsoleTextAttribute")
			handle, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)

			switch arguments[0].IntegerValue {
			case 1:
				_, _, _ = setConsoleTextAttributeProc.Call(uintptr(handle), uintptr(FOREGROUND_BLACK), 0)
			case 2:
				_, _, _ = setConsoleTextAttributeProc.Call(uintptr(handle), uintptr(FOREGROUND_RED), 0)
			case 3:
				_, _, _ = setConsoleTextAttributeProc.Call(uintptr(handle), uintptr(FOREGROUND_GREEN), 0)
			case 4:
				_, _, _ = setConsoleTextAttributeProc.Call(uintptr(handle), uintptr(FOREGROUND_YELLOW), 0)
			case 5:
				_, _, _ = setConsoleTextAttributeProc.Call(uintptr(handle), uintptr(FOREGROUND_BLUE), 0)
			case 6:
				_, _, _ = setConsoleTextAttributeProc.Call(uintptr(handle), uintptr(FOREGROUND_PURPLE), 0)
			case 7:
				_, _, _ = setConsoleTextAttributeProc.Call(uintptr(handle), uintptr(FOREGROUND_CYAN), 0)
			case 8:
				_, _, _ = setConsoleTextAttributeProc.Call(uintptr(handle), uintptr(FOREGROUND_WHITE), 0)
			default:
				_, _, _ = setConsoleTextAttributeProc.Call(uintptr(handle), uintptr((*globalSettings).consoleInfo.Attributes), 0)
			}
		}
	}

	return ret
}
