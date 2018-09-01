package main

import (
	"os"
	"errors"
	"time"
)

func Ex_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be an integer type")
	} else {
		os.Exit(arguments[0].IntegerValue)
	}

	return ret
}

func Abt_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type")
	} else {
		*errMessage = errors.New(arguments[0].StringValue)
	}

	return ret
}

func Zzz_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be an integer type")
	} else {
		time.Sleep(time.Duration(arguments[0].IntegerValue) * time.Millisecond)
	}

	return ret
}

func Sav_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ARRAY}

	for x := 0; x < len(os.Args); x++ {
		funcReturn := FunctionReturn{Type: RET_TYPE_STRING, StringValue: os.Args[x]}
		ret.ArrayValue = append(ret.ArrayValue, funcReturn)
	}

	return ret
}