package main

import (
	"errors"
	"os"
	"strconv"
)

func Flrm_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: false}

	if arguments[0].Type != ARG_TYPE_STRING {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {

		err := os.RemoveAll(arguments[0].StringValue)

		if err == nil {
			ret.BooleanValue = true
		}
	}

	return ret
}

func Flmv_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: false}

	if arguments[0].Type != ARG_TYPE_STRING {
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if arguments[1].Type != ARG_TYPE_STRING {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {

		err := os.Rename(arguments[1].StringValue, arguments[0].StringValue)

		if err == nil {
			ret.BooleanValue = true
		}
	}

	return ret
}

func Flcp_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: false}

	if arguments[0].Type != ARG_TYPE_STRING {
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if arguments[1].Type != ARG_TYPE_STRING {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		err := FDCopy(arguments[1].StringValue, arguments[0].StringValue)

		if err == nil {
			ret.BooleanValue = true
		}
	}

	return ret
}

func Fo_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if arguments[0].Type != ARG_TYPE_STRING {
		//file mode
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if arguments[1].Type != ARG_TYPE_STRING {
		//file name
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {

		var file *os.File
		var err error
		proceed_to_open_file := true

		if arguments[0].StringValue == "r" {
			// read
			file, err = os.OpenFile(arguments[1].StringValue, os.O_RDONLY, 0755)
		} else if arguments[0].StringValue == "w" {
			// write
			file, err = os.OpenFile(arguments[1].StringValue, os.O_WRONLY|os.O_CREATE, 0755)
		} else if arguments[0].StringValue == "a" {
			// append
			file, err = os.OpenFile(arguments[1].StringValue, os.O_APPEND|os.O_WRONLY, 0755)
		} else {
			// invalid mode
			proceed_to_open_file = false
			*errMessage = errors.New("Error: Invalid file mode on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		}

		if proceed_to_open_file {
			if err != nil {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {
				file_handler := "file_" + generateRandomNumbers()
				(*globalSettings).fileHandler[file_handler] = file
				ret.StringValue = file_handler
			}
		}
	}

	return ret
}

func Fc_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if arguments[0].Type != ARG_TYPE_STRING {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if (*globalSettings).fileHandler[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Invalid file reference on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := (*globalSettings).fileHandler[arguments[0].StringValue].Close()

			if err != nil {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {
				delete((*globalSettings).fileHandler, arguments[0].StringValue)
			}
		}
	}

	return ret
}

func Fw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if arguments[0].Type != ARG_TYPE_STRING {
		//string to write
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if arguments[1].Type != ARG_TYPE_STRING {
		//file reference
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if (*globalSettings).fileHandler[arguments[1].StringValue] == nil {
			*errMessage = errors.New("Error: Invalid file reference on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			_, err := (*globalSettings).fileHandler[arguments[1].StringValue].WriteString(escapeString(arguments[0].StringValue))

			if err != nil {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}
		}
	}

	return ret
}

func Fr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if arguments[0].Type != ARG_TYPE_INTEGER {
		//bytes
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if arguments[1].Type != ARG_TYPE_STRING {
		//file reference
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if (*globalSettings).fileHandler[arguments[1].StringValue] == nil {
			*errMessage = errors.New("Error: Invalid file reference on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			data := make([]byte, arguments[0].IntegerValue)

			_, err := (*globalSettings).fileHandler[arguments[1].StringValue].Read(data)
			if err != nil {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}

			ret.StringValue = string(data)
		}
	}

	return ret
}
