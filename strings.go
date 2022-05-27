package main

import (
	"errors"
	"strconv"
	"strings"
)

func Str_rpl_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		ret.StringValue = strings.Replace(arguments[2].StringValue, arguments[1].StringValue, arguments[0].StringValue, -1)
	}

	return ret
}

func Str_spl_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ARRAY}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		splitted_str := strings.Split(arguments[1].StringValue, arguments[0].StringValue)

		for x := 0; x < len(splitted_str); x++ {
			funcReturn := FunctionReturn{Type: RET_TYPE_STRING}
			funcReturn.StringValue = splitted_str[x]
			ret.ArrayValue = append(ret.ArrayValue, funcReturn)
		}
	}

	return ret
}

func Str_l_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		ret.StringValue = strings.ToLower(arguments[0].StringValue)
	}

	return ret
}

func Str_u_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		ret.StringValue = strings.ToUpper(arguments[0].StringValue)
	}

	return ret
}

func Str_t_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		ret.StringValue = strings.Trim(arguments[0].StringValue, " ")
	}

	return ret
}

func Str_chr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		if arguments[0].IntegerValue >= 0 && arguments[0].IntegerValue <= 1114111 {
			ch := rune(arguments[0].IntegerValue)
			ret.StringValue = string(ch)
		} else {
			*errMessage = errors.New("Error: Not in range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		}
	}

	return ret
}

func Str_ord_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: -1}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if len(arguments[0].StringValue) == 1 {
			ret.IntegerValue = int(arguments[0].StringValue[0])
		} else {
			*errMessage = errors.New("Error: Expected a character on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		}
	}

	return ret
}

func Str_sub_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		if arguments[1].IntegerValue < 0 {
			*errMessage = errors.New("Error: Parameter 2 must be greater than 0 on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			ctr := 0
			for x := arguments[1].IntegerValue; x < len(arguments[2].StringValue); x++ {
				ret.StringValue += string(arguments[2].StringValue[x])
				ctr += 1
				if ctr == arguments[0].IntegerValue {
					break
				}
			}
		}
	}

	return ret
}

func Str_ind_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: -1}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		ret.IntegerValue = strings.Index(arguments[1].StringValue, arguments[0].StringValue)
	}

	return ret
}
