//go:build ray
// +build ray

package main

import (
	"errors"
	"github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"strconv"
)

const RAYLIB_ENABLED bool = true

var RAYLIB_COLORS = []color.RGBA{
	rl.RayWhite,
	rl.White,
	rl.Black,
	rl.Blank, //transparent
	rl.Magenta,
	rl.Blue,
	rl.Red,
	rl.Pink,
	rl.Orange,
	rl.Yellow,
	rl.Gray,
}

func Rl_iw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		rl.InitWindow(int32(arguments[2].IntegerValue), int32(arguments[1].IntegerValue), arguments[0].StringValue)
	}

	return FunctionReturn{Type: RET_TYPE_NONE}
}

func Rl_scw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	return FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: rl.WindowShouldClose()}
}

func Rl_cw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	rl.CloseWindow()
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func Rl_bd_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	rl.BeginDrawing()
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func Rl_ed_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	rl.EndDrawing()
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func Rl_cb_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {

		if arguments[0].IntegerValue < 0 || (arguments[0].IntegerValue-1) > len(RAYLIB_COLORS) {
			*errMessage = errors.New("Error: Parameter out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			rl.ClearBackground(RAYLIB_COLORS[arguments[0].IntegerValue])
		}
	}
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func Rl_stf_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		rl.SetTargetFPS(int32(arguments[0].IntegerValue))
	}
	return FunctionReturn{Type: RET_TYPE_NONE}
}

func Rl_li_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		image := rl.LoadImage(arguments[0].StringValue)
		image_reference := "rlimg_" + generateRandomNumbers()
		(*globalSettings).rayImage[image_reference] = image
		ret.StringValue = image_reference
	}

	return ret
}

func Rl_ui_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if _, ok := (*globalSettings).rayImage[arguments[0].StringValue]; ok {
			rl.UnloadImage((*globalSettings).rayImage[arguments[0].StringValue])
			delete((*globalSettings).rayImage, arguments[0].StringValue)
		} else {
			*errMessage = errors.New("Error: Uninitialized image on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		}
	}
	return FunctionReturn{Type: RET_TYPE_NONE}
}
