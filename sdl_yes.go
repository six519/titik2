//go:build sdl
// +build sdl

package main

import (
	"github.com/veandco/go-sdl2/sdl"
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

func S_i_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: true}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		err := sdl.Init(SDL_INIT_TYPES[arguments[0].IntegerValue])

		if err != nil {
			ret.BooleanValue = false
		}
	}

	return ret
}

func S_q_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	sdl.Quit()
	return FunctionReturn{Type: RET_TYPE_NONE}
}
