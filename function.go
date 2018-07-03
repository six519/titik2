package main

import (
	"fmt"
	"strconv"
)

//function return type
const RET_TYPE_NONE int = 0
const RET_TYPE_STRING int = 1
const RET_TYPE_INTEGER int = 2
const RET_TYPE_FLOAT int = 3
const RET_TYPE_ARRAY int = 4

//function argument types
const ARG_TYPE_NONE int = 0
const ARG_TYPE_STRING int = 1
const ARG_TYPE_INTEGER int = 2
const ARG_TYPE_FLOAT int = 3
const ARG_TYPE_ARRAY int = 4

type FunctionReturn struct {
	Type int
	StringValue string
	IntegerValue int
	FloatValue float64
}

type FunctionArgument struct {
	Type int
	StringValue string
	IntegerValue int
	FloatValue float64
}

type Execute func([]FunctionArgument) FunctionReturn

type Function struct {
	Name string
	IsNative bool
	Tokens []Token
	Run Execute
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

func defineFunction(globalFunctionArray *[]Function, funcName string, funcExec Execute, argumentCount int, isNative bool) {
	function := Function{Name: funcName, IsNative: isNative, Run: funcExec, ArgumentCount: argumentCount}
	//append to global functions
	*globalFunctionArray = append(*globalFunctionArray, function)
}

//native functions
func P_execute(arguments []FunctionArgument) FunctionReturn {
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
	} else {
		//Nil
		fmt.Println("Nil")
	}

	return ret
}

func initNativeFunctions(globalFunctionArray *[]Function) {
	
	//p(<anyvar>)
	defineFunction(globalFunctionArray, "p", P_execute, 1, true)
}