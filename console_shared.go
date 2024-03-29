package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func P_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if arguments[0].Type == ARG_TYPE_FLOAT {
		fmt.Printf("%f", arguments[0].FloatValue)
		ret.StringValue = strconv.FormatFloat(arguments[0].FloatValue, 'f', -1, 64)
	} else if arguments[0].Type == ARG_TYPE_STRING {
		fmt.Printf("%s", escapeString(arguments[0].StringValue))
		ret.StringValue = arguments[0].StringValue
	} else if arguments[0].Type == ARG_TYPE_INTEGER {
		//integer
		fmt.Printf("%d", arguments[0].IntegerValue)
		ret.StringValue = strconv.Itoa(arguments[0].IntegerValue)
	} else if arguments[0].Type == ARG_TYPE_BOOLEAN {
		//boolean
		fmt.Printf("%v", arguments[0].BooleanValue)
		if arguments[0].BooleanValue {
			ret.StringValue = "true"
		} else {
			ret.StringValue = "false"
		}
	} else if arguments[0].Type == ARG_TYPE_ASSOCIATIVE_ARRAY {
		strVal := ""
		x := 0

		for k, v := range arguments[0].AssociativeArrayValue {

			if v.Type == ARG_TYPE_FLOAT {
				strVal = strVal + k + ":" + strconv.FormatFloat(v.FloatValue, 'f', -1, 64)
			} else if v.Type == ARG_TYPE_STRING {
				strVal = strVal + k + ":" + v.StringValue
			} else if v.Type == ARG_TYPE_INTEGER {
				strVal = strVal + k + ":" + strconv.Itoa(v.IntegerValue)
			} else if v.Type == ARG_TYPE_BOOLEAN {
				if v.BooleanValue {
					strVal = strVal + k + ":" + "true"
				} else {
					strVal = strVal + k + ":" + "false"
				}
			} else {
				strVal = strVal + k + ":" + "Nil"
			}

			if (x + 1) != len(arguments[0].AssociativeArrayValue) {
				strVal = strVal + ", "
			}

			x += 1
		}

		fmt.Printf("{%s}", strVal)
	} else if arguments[0].Type == ARG_TYPE_ARRAY {
		strVal := ""

		for x := 0; x < len(arguments[0].ArrayValue); x++ {
			if arguments[0].ArrayValue[x].Type == ARG_TYPE_FLOAT {
				strVal = strVal + strconv.FormatFloat(arguments[0].ArrayValue[x].FloatValue, 'f', -1, 64)
			} else if arguments[0].ArrayValue[x].Type == ARG_TYPE_STRING {
				strVal = strVal + arguments[0].ArrayValue[x].StringValue
			} else if arguments[0].ArrayValue[x].Type == ARG_TYPE_INTEGER {
				strVal = strVal + strconv.Itoa(arguments[0].ArrayValue[x].IntegerValue)
			} else if arguments[0].ArrayValue[x].Type == ARG_TYPE_BOOLEAN {
				if arguments[0].ArrayValue[x].BooleanValue {
					strVal = strVal + "true"
				} else {
					strVal = strVal + "false"
				}
			} else {
				strVal = strVal + "Nil"
			}

			if (x + 1) != len(arguments[0].ArrayValue) {
				strVal = strVal + ", "
			}
		}

		fmt.Printf("[%s]", strVal)
	} else {
		//Nil
		fmt.Println("Nil")
	}

	return ret
}

func R_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s", escapeString(arguments[0].StringValue))
		text, _ := reader.ReadString('\n')
		ret.StringValue = strings.Trim(text, "\n")
	}

	return ret
}

func Rp_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		fmt.Printf("%s", escapeString(arguments[0].StringValue))
		text, err := terminal.ReadPassword(0)
		if err == nil {
			ret.StringValue = strings.Trim(string(text), "\n")
		}
	}

	return ret
}
