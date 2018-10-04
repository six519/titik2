package main

import (
	"os"
	"errors"
	"time"
	"strconv"
	"os/exec"
	"strings"
	"fmt"
	"path/filepath"
)

func Ex_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be an integer type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		os.Exit(arguments[0].IntegerValue)
	}

	return ret
}

func Abt_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		*errMessage = errors.New(arguments[0].StringValue)
	}

	return ret
}

func Exe_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ASSOCIATIVE_ARRAY}

	ret.AssociativeArrayValue = make(map[string]FunctionReturn)
	isSuccess := false
	outString := ""

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		cmds := strings.Split(arguments[0].StringValue, " ")
	
		cmd := exec.Command(cmds[0], cmds[1:]...)
		out, err := cmd.Output()
	
		if(err != nil) {
			fmt.Println(err.Error())
			outString = err.Error()
		} else {
			isSuccess = true
			if(out != nil) {
				outString = string(out)
				fmt.Println(string(out))
			}
		}
	}

	ret.AssociativeArrayValue["success"] = FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: isSuccess}
	ret.AssociativeArrayValue["output"] = FunctionReturn{Type: RET_TYPE_STRING, StringValue: outString}

	return ret
}

func Zzz_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be an integer type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		time.Sleep(time.Duration(arguments[0].IntegerValue) * time.Millisecond)
	}

	return ret
}

func Sav_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ARRAY}

	for x := 0; x < len(os.Args); x++ {
		funcReturn := FunctionReturn{Type: RET_TYPE_STRING, StringValue: os.Args[x]}
		ret.ArrayValue = append(ret.ArrayValue, funcReturn)
	}

	return ret
}

func Gcp_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

    dir, err := filepath.Abs(filepath.Dir(file_name))
	
	if err == nil {
		ret.StringValue = dir
    }

	return ret
}