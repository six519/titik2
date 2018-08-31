package main

import (
	"fmt"
	"errors"
)

//function return type
const (
	RET_TYPE_NONE = iota
	RET_TYPE_STRING
	RET_TYPE_INTEGER
	RET_TYPE_FLOAT
	RET_TYPE_ARRAY
	RET_TYPE_ASSOCIATIVE_ARRAY
	RET_TYPE_BOOLEAN
)

//function argument types
const (
	ARG_TYPE_NONE = iota
	ARG_TYPE_STRING
	ARG_TYPE_INTEGER
	ARG_TYPE_FLOAT
	ARG_TYPE_ARRAY
	ARG_TYPE_ASSOCIATIVE_ARRAY
	ARG_TYPE_BOOLEAN
)

type FunctionReturn struct {
	Type int
	StringValue string
	IntegerValue int
	FloatValue float64
	BooleanValue bool
	ArrayValue []FunctionReturn
	AssociativeArrayValue map[string]FunctionReturn
}

type FunctionArgument struct {
	Type int
	StringValue string
	IntegerValue int
	FloatValue float64
	BooleanValue bool
	ArrayValue []FunctionArgument
	AssociativeArrayValue map[string]FunctionArgument
}

type Execute func([]FunctionArgument, *error, *[]Variable, *[]Function, string, *[]string, *WebObject) FunctionReturn

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
func ReverseBoolean_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: false}

	if(arguments[0].Type != ARG_TYPE_BOOLEAN) {
		*errMessage = errors.New("Error: Parameter must be a boolean type")
	} else {
		ret.BooleanValue = !arguments[0].BooleanValue
	}

	return ret
}

func Len_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: 0}

	if(arguments[0].Type != ARG_TYPE_ARRAY && arguments[0].Type != ARG_TYPE_ASSOCIATIVE_ARRAY) {
		*errMessage = errors.New("Error: Parameter must be a lineup or glossary type")
	} else {
		if(arguments[0].Type == ARG_TYPE_ARRAY) {
			ret.IntegerValue = len(arguments[0].ArrayValue)
		} else {
			ret.IntegerValue = len(arguments[0].AssociativeArrayValue)
		}
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

	//len(<lineup>)
	defineFunction(globalFunctionArray, "len", Len_execute, 1, true)

	//sav()
	defineFunction(globalFunctionArray, "sav", Sav_execute, 0, true)

	//sc(<integer>)
	defineFunction(globalFunctionArray, "sc", Sc_execute, 1, true)


	//WEB FUNCTIONALITIES
	//http_au(<string>, <string>)
	defineFunction(globalFunctionArray, "http_au", Http_au_execute, 2, true)

	//http_run(<string>)
	defineFunction(globalFunctionArray, "http_run", Http_run_execute, 1, true)

	//http_p(<string>)
	defineFunction(globalFunctionArray, "http_p", Http_p_execute, 1, true)

	//http_gm()
	defineFunction(globalFunctionArray, "http_gm", Http_gm_execute, 0, true)

	//http_su(<string>, <string>)
	defineFunction(globalFunctionArray, "http_su", Http_su_execute, 2, true)
}