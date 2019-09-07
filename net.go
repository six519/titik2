package main

import (
	"net"
	"errors"
	"strconv"
	"io/ioutil"
	//"fmt"
)

func Netc_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[1].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		var err error
		(*globalSettings).netConnection[scopeName], err = net.Dial(arguments[1].StringValue, arguments[0].StringValue)

		if(err != nil) {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		}
	}

	return ret
}

func Netx_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if ((*globalSettings).netConnection[scopeName] == nil) {
		*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		err := (*globalSettings).netConnection[scopeName].Close()

		if(err != nil) {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		}
	}

	return ret
}

func Netw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if ((*globalSettings).netConnection[scopeName] == nil) {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			_, err := (*globalSettings).netConnection[scopeName].Write([]byte(escapeString(arguments[0].StringValue)))
	
			if(err != nil) {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}
		}
	}

	return ret
}

func Netr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if ((*globalSettings).netConnection[scopeName] == nil) {
		*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		msg, err := ioutil.ReadAll((*globalSettings).netConnection[scopeName])

		if(err != nil) {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			ret.StringValue = string(msg)
		}
	}

	return ret
}