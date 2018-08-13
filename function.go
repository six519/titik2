package main

import (
	"fmt"
	"strconv"
	"os"
	"errors"
	"time"
	"bufio"
	"strings"
)

//function return type
const (
	RET_TYPE_NONE = iota
	RET_TYPE_STRING
	RET_TYPE_INTEGER
	RET_TYPE_FLOAT
	RET_TYPE_ARRAY
	RET_TYPE_BOOLEAN
)

//function argument types
const (
	ARG_TYPE_NONE = iota
	ARG_TYPE_STRING
	ARG_TYPE_INTEGER
	ARG_TYPE_FLOAT
	ARG_TYPE_ARRAY
	ARG_TYPE_BOOLEAN
)

type FunctionReturn struct {
	Type int
	StringValue string
	IntegerValue int
	FloatValue float64
	BooleanValue bool
}

type FunctionArgument struct {
	Type int
	StringValue string
	IntegerValue int
	FloatValue float64
	BooleanValue bool
}

type Execute func([]FunctionArgument, *error) FunctionReturn

type Function struct {
	Name string
	IsNative bool
	Tokens []Token
	Run Execute
	Arguments []Token
	ArgumentCount int
}

func DumpFunction(functions []Function) {
	fmt.Printf("====================================\n")

	for x := 0; x < len(functions); x++ {
		fmt.Printf("Function Name: %s\n", functions[x].Name)
		fmt.Printf("Argument Count: %d\n", functions[x].ArgumentCount)

		if(functions[x].IsNative) {
			fmt.Println("Is Native: Yes")
		} else {
			fmt.Println("Is Native: No")
		}

		fmt.Printf("====================================\n")
	}
}

func isFunctionExists(token Token, globalFunctionArray []Function) (bool, int) {

	for x := 0; x < len(globalFunctionArray); x++ {
		if(globalFunctionArray[x].Name == token.Value) {
			return true, x
		}
	}

	return false, 0
}

func isParamExists(token Token, functionParams []Token) bool {

	for x := 0; x < len(functionParams); x++ {
		if(functionParams[x].Value == token.Value) {
			return true
		}
	}

	return false
}

func defineFunction(globalFunctionArray *[]Function, funcName string, funcExec Execute, argumentCount int, isNative bool) {
	function := Function{Name: funcName, IsNative: isNative, Run: funcExec, ArgumentCount: argumentCount}
	//append to global functions
	*globalFunctionArray = append(*globalFunctionArray, function)
}

//native functions
func P_execute(arguments []FunctionArgument, errMessage *error) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(arguments[0].Type == ARG_TYPE_FLOAT) {
		fmt.Printf("%f\n", arguments[0].FloatValue)
		ret.StringValue = strconv.FormatFloat(arguments[0].FloatValue, 'f', -1, 64)
	} else if(arguments[0].Type == ARG_TYPE_STRING) {
		fmt.Printf("%s\n", arguments[0].StringValue)
		ret.StringValue = arguments[0].StringValue
	} else if(arguments[0].Type == ARG_TYPE_INTEGER) {
		//integer
		fmt.Printf("%d\n", arguments[0].IntegerValue)
		ret.StringValue = strconv.Itoa(arguments[0].IntegerValue)
	} else if(arguments[0].Type == ARG_TYPE_BOOLEAN) {
		//boolean
		fmt.Printf("%v\n", arguments[0].BooleanValue)
		if(arguments[0].BooleanValue) {
			ret.StringValue = "true"
		} else {
			ret.StringValue = "false"
		}	
	} else {
		//Nil
		fmt.Println("Nil")
	}

	return ret
}

func Ex_execute(arguments []FunctionArgument, errMessage *error) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be an integer type")
	} else {
		os.Exit(arguments[0].IntegerValue)
	}

	return ret
}

func Abt_execute(arguments []FunctionArgument, errMessage *error) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type")
	} else {
		*errMessage = errors.New(arguments[0].StringValue)
	}

	return ret
}

func ReverseBoolean_execute(arguments []FunctionArgument, errMessage *error) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: false}

	if(arguments[0].Type != ARG_TYPE_BOOLEAN) {
		*errMessage = errors.New("Error: Parameter must be a boolean type")
	} else {
		ret.BooleanValue = !arguments[0].BooleanValue
	}

	return ret
}

func Zzz_execute(arguments []FunctionArgument, errMessage *error) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be an integer type")
	} else {
		time.Sleep(time.Duration(arguments[0].IntegerValue) * time.Millisecond)
	}

	return ret
}

func R_execute(arguments []FunctionArgument, errMessage *error) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type")
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s", arguments[0].StringValue)
		text, _ := reader.ReadString('\n')
		ret.StringValue = strings.Trim(text, "\n")
	}

	return ret
}

func Toi_execute(arguments []FunctionArgument, errMessage *error) FunctionReturn {
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

func Tos_execute(arguments []FunctionArgument, errMessage *error) FunctionReturn {
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

func initNativeFunctions(globalFunctionArray *[]Function) {
	
	//p(<anyvar>)
	defineFunction(globalFunctionArray, "p", P_execute, 1, true)

	//ex(<integer>)
	defineFunction(globalFunctionArray, "ex", Ex_execute, 1, true)

	//abt(<string>)
	defineFunction(globalFunctionArray, "abt", Abt_execute, 1, true)

	//!(<bool>)
	defineFunction(globalFunctionArray, "!", ReverseBoolean_execute, 1, true)

	//zzz(<integer>)
	defineFunction(globalFunctionArray, "zzz", Zzz_execute, 1, true)

	//r(<string>)
	defineFunction(globalFunctionArray, "r", R_execute, 1, true)

	//toi(<anyvar>)
	defineFunction(globalFunctionArray, "toi", Toi_execute, 1, true)

	//tos(<anyvar>)
	defineFunction(globalFunctionArray, "tos", Tos_execute, 1, true)
}