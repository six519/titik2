package main

import (
	"os"
	"net"
	"errors"
	"strconv"
	//"io/ioutil"
	"fmt"
)

func Netc_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[1].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		connection, err := net.Dial(arguments[1].StringValue, arguments[0].StringValue)

		if(err != nil) {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			connection_reference := "con_" + generateRandomNumbers()
			(*globalSettings).netConnection[connection_reference] = connection
			ret.StringValue = connection_reference
		}
	}

	return ret
}

func Netl_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[1].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		connection, err := net.Listen(arguments[1].StringValue, arguments[0].StringValue)

		if(err != nil) {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			connection_reference := "conl_" + generateRandomNumbers()
			(*globalSettings).netConnectionListener[connection_reference] = connection
			ret.StringValue = connection_reference
		}
	}

	return ret
}

func Netla_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if ((*globalSettings).netConnectionListener[arguments[0].StringValue] == nil) {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			connection, err := (*globalSettings).netConnectionListener[arguments[0].StringValue].Accept()
			
			if err == nil {
				connection_reference := "con_" + generateRandomNumbers()
				(*globalSettings).netConnection[connection_reference] = connection
				ret.StringValue = connection_reference	
			}
		}
	}

	return ret
}

//helper for Netlaf_execute
func netHandleRequest(connection_reference string, function_name string, globalVariableArray *[]Variable, globalFunctionArray *[]Function, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) {
	t := Token{Value: function_name}
	isExists, funcIndex := isFunctionExists(t, (*globalFunctionArray))
	if(isExists) {
		array := (*globalFunctionArray)

		if(array[funcIndex].ArgumentCount == 1) {
			thisScopeName := array[funcIndex].Name + generateRandomNumbers()
			var thisGotReturn bool = false
			var thisReturnToken Token
			var thisNeedBreak bool = false
			var thisStackReference []Token
			
			newVar := Variable{Name: array[funcIndex].Arguments[0].Value, ScopeName: thisScopeName, Type: VARIABLE_TYPE_STRING, StringValue: connection_reference}
			*globalSettings.globalVariableArray = append(*globalSettings.globalVariableArray, newVar)

			//execute user defined function
			prsr := Parser{}
			parserErr := prsr.Parse(array[funcIndex].Tokens, globalSettings.globalVariableArray, globalSettings.globalFunctionArray, thisScopeName, globalSettings.globalNativeVarList, &thisGotReturn, &thisReturnToken, false, &thisNeedBreak, &thisStackReference, globalSettings)
	
			if(parserErr != nil) {
				//error
				fmt.Println(parserErr)
				os.Exit(2)
			}

			if(thisGotReturn) {
				//ignore return
			}
		} else {
			fmt.Println("Error: Function argument takes exactly one on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			os.Exit(2)
		}
	} else {
		fmt.Println("Error: Function handler doesn't exists on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		os.Exit(2)
	}
}
//end helper for Netlaf_execute

func Netlaf_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[1].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if ((*globalSettings).netConnectionListener[arguments[1].StringValue] == nil) {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			connection, err := (*globalSettings).netConnectionListener[arguments[1].StringValue].Accept()
			
			if err == nil {
				connection_reference := "con_" + generateRandomNumbers()
				(*globalSettings).netConnection[connection_reference] = connection
				ret.StringValue = connection_reference
				go netHandleRequest(connection_reference, arguments[0].StringValue, globalVariableArray, globalFunctionArray, globalNativeVarList, globalSettings, line_number, column_number, file_name)
			}
		}
	}

	return ret
}

func Netlx_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if ((*globalSettings).netConnectionListener[arguments[0].StringValue] == nil) {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := (*globalSettings).netConnectionListener[arguments[0].StringValue].Close()
	
			if(err != nil) {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {
				delete((*globalSettings).netConnectionListener, arguments[0].StringValue)
			}
		}
	}

	return ret
}

func Netx_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if ((*globalSettings).netConnection[arguments[0].StringValue] == nil) {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := (*globalSettings).netConnection[arguments[0].StringValue].Close()
	
			if(err != nil) {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {
				delete((*globalSettings).netConnection, arguments[0].StringValue)
			}
		}
	}

	return ret
}

func Netw_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[1].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if ((*globalSettings).netConnection[arguments[1].StringValue] == nil) {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			_, err := (*globalSettings).netConnection[arguments[1].StringValue].Write([]byte(escapeString(arguments[0].StringValue)))
	
			if(err != nil) {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}
		}
	}

	return ret
}

func Netr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter 2 must be an integer type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else 	if(arguments[1].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if ((*globalSettings).netConnection[arguments[1].StringValue] == nil) {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			buf := make([]byte, arguments[0].IntegerValue)
			_, err := (*globalSettings).netConnection[arguments[1].StringValue].Read(buf[0:])
			if err == nil {
				ret.StringValue = string(buf)
			}
		}
	}

	return ret
}