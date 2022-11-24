//go:build win
// +build win

package main

import (
	"errors"
	"strconv"
	"syscall"
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
