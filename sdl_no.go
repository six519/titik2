//go:build !sdl
// +build !sdl

package main

const SDL_ENABLED bool = false

var SDL_INIT_TYPES = []uint32{}
var SDL_WPOSITIONS = []int32{}
var SDL_WFLAGS = []uint32{}

func S_i_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_q_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}
