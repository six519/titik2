//go:build !sdl
// +build !sdl

package main

const SDL_ENABLED bool = false

var SDL_INIT_TYPES = []uint32{}
var SDL_WPOSITIONS = []int32{}
var SDL_WFLAGS = []uint32{}
var SDL_EVENT_TYPES = map[uint32]int{}

func S_i_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_q_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_cw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_dw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_usw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_gsw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_cr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_frsw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_pe_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_ce_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_gte_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_d_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_it_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func S_qt_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_NONE}
}
