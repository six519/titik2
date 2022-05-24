package main

import (
	"errors"
	"net"
	"os"
	"strconv"

	//"io/ioutil"
	"fmt"
)

func Netc_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		connection, err := net.Dial(arguments[1].StringValue, arguments[0].StringValue)

		if err != nil {
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

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		connection, err := net.Listen(arguments[1].StringValue, arguments[0].StringValue)

		if err != nil {
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

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).netConnectionListener[arguments[0].StringValue] == nil {
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
	if isExists {
		array := (*globalFunctionArray)

		if array[funcIndex].ArgumentCount == 1 {
			thisScopeName := array[funcIndex].Name + generateRandomNumbers()
			var thisGotReturn bool = false
			var thisReturnToken Token
			var thisNeedBreak bool = false
			var thisStackReference []Token
			var getLastStackBool bool = false
			var lastStackBool bool = false

			newVar := Variable{Name: array[funcIndex].Arguments[0].Value, ScopeName: thisScopeName, Type: VARIABLE_TYPE_STRING, StringValue: connection_reference}
			*globalSettings.globalVariableArray = append(*globalSettings.globalVariableArray, newVar)

			//execute user defined function
			prsr := Parser{}
			parserErr := prsr.Parse(array[funcIndex].Tokens, globalSettings.globalVariableArray, globalSettings.globalFunctionArray, thisScopeName, globalSettings.globalNativeVarList, &thisGotReturn, &thisReturnToken, false, &thisNeedBreak, &thisStackReference, globalSettings, getLastStackBool, &lastStackBool)

			if parserErr != nil {
				//error
				fmt.Println(parserErr)
				os.Exit(2)
			}

			if thisGotReturn {
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

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).netConnectionListener[arguments[1].StringValue] == nil {
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

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).netConnectionListener[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := (*globalSettings).netConnectionListener[arguments[0].StringValue].Close()

			if err != nil {
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

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).netConnection[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := (*globalSettings).netConnection[arguments[0].StringValue].Close()

			if err != nil {
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

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).netConnection[arguments[1].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			_, err := (*globalSettings).netConnection[arguments[1].StringValue].Write([]byte(escapeString(arguments[0].StringValue)))

			if err != nil {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}
		}
	}

	return ret
}

func Netr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		if (*globalSettings).netConnection[arguments[1].StringValue] == nil {
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

func Netul_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		connection, err := net.ListenUDP(arguments[2].StringValue, &net.UDPAddr{
			Port: arguments[0].IntegerValue,
			IP:   net.ParseIP(arguments[1].StringValue),
		})

		if err != nil {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			connection_reference := "conul_" + generateRandomNumbers()
			(*globalSettings).netUDPConnectionListener[connection_reference] = connection
			ret.StringValue = connection_reference
		}
	}

	return ret
}

func Netulf_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_INTEGER) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		connection, err := net.ListenUDP(arguments[3].StringValue, &net.UDPAddr{
			Port: arguments[1].IntegerValue,
			IP:   net.ParseIP(arguments[2].StringValue),
		})

		if err != nil {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			connection_reference := "conul_" + generateRandomNumbers()
			(*globalSettings).netUDPConnectionListener[connection_reference] = connection
			ret.StringValue = connection_reference
			go netHandleRequest(connection_reference, arguments[0].StringValue, globalVariableArray, globalFunctionArray, globalNativeVarList, globalSettings, line_number, column_number, file_name)
		}
	}

	return ret
}

func Netur_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		if (*globalSettings).netUDPConnectionListener[arguments[1].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			buf := make([]byte, arguments[0].IntegerValue)
			_, _, err := (*globalSettings).netUDPConnectionListener[arguments[1].StringValue].ReadFromUDP(buf[0:])
			if err == nil {
				ret.StringValue = string(buf)
			}
		}
	}

	return ret
}

func Netus_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_INTEGER) {
		if (*globalSettings).netUDPConnectionListener[arguments[3].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			_, err := (*globalSettings).netUDPConnectionListener[arguments[3].StringValue].WriteToUDP([]byte(escapeString(arguments[2].StringValue)), &net.UDPAddr{
				Port: arguments[0].IntegerValue,
				IP:   net.ParseIP(arguments[1].StringValue),
			})

			if err != nil {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}
		}
	}

	return ret
}
