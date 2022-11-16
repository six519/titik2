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
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_BOOLEAN) {
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

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
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

func La_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ARRAY}
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_ARRAY) {
		if arguments[0].Type != ARG_TYPE_STRING && arguments[0].Type != ARG_TYPE_INTEGER && arguments[0].Type != ARG_TYPE_FLOAT && arguments[0].Type != ARG_TYPE_BOOLEAN && arguments[0].Type != ARG_TYPE_NONE {
			*errMessage = errors.New("Error: Parameter must be a string/integer/float/boolean/Nil type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			for x := 0; x < len(arguments[1].ArrayValue); x++ {
				ret.ArrayValue = append(ret.ArrayValue, FunctionReturn{
					Type:         arguments[1].ArrayValue[x].Type,
					StringValue:  arguments[1].ArrayValue[x].StringValue,
					IntegerValue: arguments[1].ArrayValue[x].IntegerValue,
					FloatValue:   arguments[1].ArrayValue[x].FloatValue,
					BooleanValue: arguments[1].ArrayValue[x].BooleanValue,
				})
			}
			ret.ArrayValue = append(ret.ArrayValue, FunctionReturn{
				Type:         arguments[0].Type,
				StringValue:  arguments[0].StringValue,
				IntegerValue: arguments[0].IntegerValue,
				FloatValue:   arguments[0].FloatValue,
				BooleanValue: arguments[0].BooleanValue,
			})
		}
	}
	return ret
}

func Lp_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ARRAY}
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_ARRAY) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		if (arguments[0].IntegerValue+1) > len(arguments[1].ArrayValue) || len(arguments[1].ArrayValue) == 0 || arguments[0].IntegerValue < 0 {
			*errMessage = errors.New("Error: Index out of range on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			for x := 0; x < len(arguments[1].ArrayValue); x++ {
				if x != arguments[0].IntegerValue {
					ret.ArrayValue = append(ret.ArrayValue, FunctionReturn{
						Type:         arguments[1].ArrayValue[x].Type,
						StringValue:  arguments[1].ArrayValue[x].StringValue,
						IntegerValue: arguments[1].ArrayValue[x].IntegerValue,
						FloatValue:   arguments[1].ArrayValue[x].FloatValue,
						BooleanValue: arguments[1].ArrayValue[x].BooleanValue,
					})
				}
			}
		}
	}
	return ret
}

func Sgv_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) {
		if isSystemVariable(arguments[1].StringValue, *globalNativeVarList) {
			*errMessage = errors.New("Error: Cannot assign to " + arguments[1].StringValue + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			isExists, varIndex := isGlobalVariableExists(arguments[1].StringValue, *globalVariableArray)

			if !isExists {
				*errMessage = errors.New("Error: Global variable '" + arguments[1].StringValue + "' doesn't exists on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {
				if (*globalVariableArray)[varIndex].IsConstant {
					*errMessage = errors.New("Error: Cannot override constant '" + arguments[1].StringValue + "' on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
				} else {
					// set global var
					if arguments[0].Type == ARG_TYPE_ARRAY {
						(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_ARRAY
						(*globalVariableArray)[varIndex].ArrayValue = nil
						for x := 0; x < len(arguments[0].ArrayValue); x++ {
							thisVar := Variable{}

							if arguments[0].ArrayValue[x].Type == ARG_TYPE_INTEGER {
								thisVar.Type = VARIABLE_TYPE_INTEGER
								thisVar.IntegerValue = arguments[0].ArrayValue[x].IntegerValue
							} else if arguments[0].ArrayValue[x].Type == ARG_TYPE_STRING {
								thisVar.Type = VARIABLE_TYPE_STRING
								thisVar.StringValue = arguments[0].ArrayValue[x].StringValue
							} else if arguments[0].ArrayValue[x].Type == ARG_TYPE_FLOAT {
								thisVar.Type = VARIABLE_TYPE_FLOAT
								thisVar.FloatValue = arguments[0].ArrayValue[x].FloatValue
							} else if arguments[0].ArrayValue[x].Type == ARG_TYPE_BOOLEAN {
								thisVar.Type = VARIABLE_TYPE_BOOLEAN
								thisVar.BooleanValue = arguments[0].ArrayValue[x].BooleanValue
							} else {
								thisVar.Type = VARIABLE_TYPE_NONE
							}

							(*globalVariableArray)[varIndex].ArrayValue = append((*globalVariableArray)[varIndex].ArrayValue, thisVar)
						}

					} else if arguments[0].Type == ARG_TYPE_ASSOCIATIVE_ARRAY {
						(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_ASSOCIATIVE_ARRAY
						(*globalVariableArray)[varIndex].AssociativeArrayValue = make(map[string]*Variable)

						for k, v := range arguments[0].AssociativeArrayValue {
							thisVar := new(Variable)

							if v.Type == ARG_TYPE_INTEGER {
								thisVar.Type = VARIABLE_TYPE_INTEGER
								thisVar.IntegerValue = v.IntegerValue
							} else if v.Type == ARG_TYPE_STRING {
								thisVar.Type = VARIABLE_TYPE_STRING
								thisVar.StringValue = v.StringValue
							} else if v.Type == ARG_TYPE_FLOAT {
								thisVar.Type = VARIABLE_TYPE_FLOAT
								thisVar.FloatValue = v.FloatValue
							} else if v.Type == ARG_TYPE_BOOLEAN {
								thisVar.Type = VARIABLE_TYPE_BOOLEAN
								thisVar.BooleanValue = v.BooleanValue
							} else {
								thisVar.Type = VARIABLE_TYPE_NONE
							}
							(*globalVariableArray)[varIndex].AssociativeArrayValue[k] = thisVar
						}

					} else if arguments[0].Type == ARG_TYPE_BOOLEAN {
						(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_BOOLEAN
						(*globalVariableArray)[varIndex].BooleanValue = arguments[0].BooleanValue
					} else if arguments[0].Type == ARG_TYPE_FLOAT {
						(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_FLOAT
						(*globalVariableArray)[varIndex].FloatValue = arguments[0].FloatValue
					} else if arguments[0].Type == ARG_TYPE_INTEGER {
						(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_INTEGER
						(*globalVariableArray)[varIndex].IntegerValue = arguments[0].IntegerValue
					} else if arguments[0].Type == ARG_TYPE_STRING {
						(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_STRING
						(*globalVariableArray)[varIndex].StringValue = arguments[0].StringValue
					} else {
						//none type
						(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_NONE
					}
				}
			}
		}
	}
	return FunctionReturn{Type: RET_TYPE_NONE}
}
