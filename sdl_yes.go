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

func S_cw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 5, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 4, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if arguments[4].IntegerValue < 0 || (arguments[4].IntegerValue-1) > len(SDL_WPOSITIONS) ||
			arguments[3].IntegerValue < 0 || (arguments[3].IntegerValue-1) > len(SDL_WPOSITIONS) ||
			arguments[0].IntegerValue < 0 || (arguments[0].IntegerValue-1) > len(SDL_WFLAGS) {
			*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {

			window, err := sdl.CreateWindow(arguments[5].StringValue, SDL_WPOSITIONS[arguments[4].IntegerValue], SDL_WPOSITIONS[arguments[3].IntegerValue], int32(arguments[2].IntegerValue), int32(arguments[1].IntegerValue), SDL_WFLAGS[arguments[0].IntegerValue])

			if err == nil {
				window_reference := "win_" + generateRandomNumbers()
				(*globalSettings).sdlWindow[window_reference] = window
				ret.Type = RET_TYPE_STRING
				ret.StringValue = window_reference
			}
		}

	}

	return ret
}

func S_dw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).sdlWindow[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized window on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).sdlWindow[arguments[0].StringValue].Destroy()
			delete((*globalSettings).sdlWindow, arguments[0].StringValue)
		}
	}

	return ret
}

func S_usw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).sdlWindow[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized window on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).sdlWindow[arguments[0].StringValue].UpdateSurface()
		}
	}

	return ret
}

func S_gsw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).sdlWindow[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized window on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			surface, err := (*globalSettings).sdlWindow[arguments[0].StringValue].GetSurface()

			if err == nil {
				surface_reference := "surf_" + generateRandomNumbers()
				(*globalSettings).sdlSurface[surface_reference] = surface
				ret.Type = RET_TYPE_STRING
				ret.StringValue = surface_reference
			}
		}
	}

	return ret
}

func S_cr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		rect := sdl.Rect{int32(arguments[3].IntegerValue), int32(arguments[2].IntegerValue), int32(arguments[1].IntegerValue), int32(arguments[0].IntegerValue)}

		rect_reference := "rect_" + generateRandomNumbers()
		(*globalSettings).sdlRect[rect_reference] = rect
		ret.Type = RET_TYPE_STRING
		ret.StringValue = rect_reference
	}

	return ret
}

func S_frsw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if (*globalSettings).sdlSurface[arguments[2].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized surface on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			if arguments[1].StringValue == "" {
				(*globalSettings).sdlSurface[arguments[2].StringValue].FillRect(nil, uint32(arguments[0].IntegerValue))
			} else {
				if val, ok := (*globalSettings).sdlRect[arguments[1].StringValue]; ok {
					(*globalSettings).sdlSurface[arguments[2].StringValue].FillRect(&val, uint32(arguments[0].IntegerValue))
				} else {
					*errMessage = errors.New("Error: Uninitialized rectangle on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
				}
			}
		}

	}

	return ret
}

func S_pe_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE} //TEMPORARY RETURN: NEED TO HANDLE EVENTS
	sdl.PollEvent()
	return ret
}
