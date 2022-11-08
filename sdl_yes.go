//go:build sdl
// +build sdl

package main

import (
	"errors"
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
)

const SDL_ENABLED bool = true

var SDL_INIT_TYPES = []uint32{
	sdl.INIT_EVERYTHING,
	sdl.INIT_TIMER,
	sdl.INIT_AUDIO,
	sdl.INIT_VIDEO,
	sdl.INIT_JOYSTICK,
	sdl.INIT_HAPTIC,
	sdl.INIT_GAMECONTROLLER,
	sdl.INIT_EVENTS,
	sdl.INIT_NOPARACHUTE,
	sdl.INIT_SENSOR,
}

var SDL_WPOSITIONS = []int32{
	sdl.WINDOWPOS_UNDEFINED,
	sdl.WINDOWPOS_CENTERED,
}

var SDL_WFLAGS = []uint32{
	sdl.WINDOW_FULLSCREEN,
	sdl.WINDOW_SHOWN,
	sdl.WINDOW_HIDDEN,
	sdl.WINDOW_BORDERLESS,
	sdl.WINDOW_RESIZABLE,
	sdl.WINDOW_MINIMIZED,
	sdl.WINDOW_MAXIMIZED,
	sdl.WINDOW_FULLSCREEN_DESKTOP,
	sdl.WINDOW_ALWAYS_ON_TOP,
	sdl.WINDOW_TOOLTIP,
	sdl.WINDOW_POPUP_MENU,
}

func S_i_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: true}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if arguments[0].IntegerValue < 0 || (arguments[0].IntegerValue-1) > len(SDL_INIT_TYPES) {
			*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := sdl.Init(SDL_INIT_TYPES[arguments[0].IntegerValue])

			if err != nil {
				ret.BooleanValue = false
			}
		}
	}

	return ret
}

func S_q_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	sdl.Quit()
	return FunctionReturn{Type: RET_TYPE_NONE}
}
