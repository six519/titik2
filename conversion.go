package main

import (
	"strconv"
)

func Toi_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: 0}

	if(arguments[0].Type == ARG_TYPE_FLOAT) {
		ret.IntegerValue = int(arguments[0].FloatValue) + 0
	} else if(arguments[0].Type == ARG_TYPE_STRING) {
		ret.IntegerValue, _ = strconv.Atoi(arguments[0].StringValue)
	} else if(arguments[0].Type == ARG_TYPE_INTEGER) {
		ret.IntegerValue = arguments[0].IntegerValue
	}

	return ret
}

func Tos_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(arguments[0].Type == ARG_TYPE_FLOAT) {
		ret.StringValue = strconv.FormatFloat(arguments[0].FloatValue, 'f', -1, 64)
	} else if(arguments[0].Type == ARG_TYPE_STRING) {
		ret.StringValue = arguments[0].StringValue
	} else if(arguments[0].Type == ARG_TYPE_INTEGER) {
		ret.StringValue = strconv.Itoa(arguments[0].IntegerValue)
	} else if(arguments[0].Type == ARG_TYPE_BOOLEAN) {
		if(arguments[0].BooleanValue) {
			ret.StringValue = "true"
		} else {
			ret.StringValue = "false"
		}	
	} else {
		ret.StringValue = ""
	}

	return ret
}