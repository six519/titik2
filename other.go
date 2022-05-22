package main

import (
	"errors"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func ReverseBoolean_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: false}

	if arguments[0].Type != ARG_TYPE_BOOLEAN {
		*errMessage = errors.New("Error: Parameter must be a boolean type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		ret.BooleanValue = !arguments[0].BooleanValue
	}

	return ret
}

func Len_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: 0}

	if arguments[0].Type != ARG_TYPE_ARRAY && arguments[0].Type != ARG_TYPE_ASSOCIATIVE_ARRAY && arguments[0].Type != ARG_TYPE_STRING {
		*errMessage = errors.New("Error: Parameter must be a lineup or glossary or string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if arguments[0].Type == ARG_TYPE_ARRAY {
			ret.IntegerValue = len(arguments[0].ArrayValue)
		} else if arguments[0].Type == ARG_TYPE_STRING {
			ret.IntegerValue = len(arguments[0].StringValue)
		} else {
			ret.IntegerValue = len(arguments[0].AssociativeArrayValue)
		}
	}

	return ret
}

func I_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if arguments[0].Type != ARG_TYPE_STRING {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if scopeName != "main" {
			*errMessage = errors.New("Error: You cannot include file inside a function on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			suffix := ""
			fileToLoad := arguments[0].StringValue

			if runtime.GOOS == "windows" {
				suffix = "\\"
				fileToLoad = strings.Replace(fileToLoad, "/", "\\", -1)
			} else {
				suffix = "/"
			}

			dir, _ := filepath.Abs(filepath.Dir(file_name))

			//open titik file to include
			lxr := Lexer{FileName: dir + suffix + fileToLoad + ".ttk"}
			fileErr := lxr.ReadSourceFile()
			if fileErr != nil {
				*errMessage = errors.New(fileErr.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {
				//generate token below
				tokenArray, tokenErr := lxr.GenerateToken()
				if tokenErr != nil {
					*errMessage = tokenErr
				} else {
					var gotReturn bool = false
					var returnToken Token
					var needBreak bool = false
					var stackReference []Token
					var getLastStackBool bool = false
					var lastStackBool bool = false

					//parser object
					prsr := Parser{}
					parserErr := prsr.Parse(tokenArray, globalVariableArray, globalFunctionArray, "main", globalNativeVarList, &gotReturn, &returnToken, false, &needBreak, &stackReference, globalSettings, getLastStackBool, &lastStackBool)
					if parserErr != nil {
						*errMessage = parserErr
					}
				}
			}
		}
	}

	return ret
}

func In_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: false}

	if arguments[0].Type == ARG_TYPE_NONE {
		ret.BooleanValue = true
	}

	return ret
}
