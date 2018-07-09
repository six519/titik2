package main

import (
	"fmt"
	"errors"
	"strconv"
)

//variable types (basic)
const (
	VARIABLE_TYPE_NONE = iota
	VARIABLE_TYPE_INTEGER
	VARIABLE_TYPE_STRING
	VARIABLE_TYPE_FLOAT
	VARIABLE_TYPE_ARRAY
)

var VARIABLE_TYPES_STRING = []string {
	"VARIABLE_TYPE_NONE",
	"VARIABLE_TYPE_INTEGER",
	"VARIABLE_TYPE_STRING",
	"VARIABLE_TYPE_FLOAT",
	"VARIABLE_TYPE_ARRAY",
}

type Variable struct {
	Name string
	ScopeName string
	Type int
	StringValue string
	IntegerValue int
	FloatValue float64
	IsConstant bool
	ArrayValue []Variable
}

func DumpVariable(variables []Variable) {
	fmt.Printf("====================================\n")

	for x := 0; x < len(variables); x++ {
		fmt.Printf("Variable Name: %s\n", variables[x].Name)

		//NOTE: ASSUME VARIABLE TYPE AS FLOAT AND INTEGER FOR NOW (TEMPORARY)
		if(variables[x].Type == VARIABLE_TYPE_FLOAT) {
			fmt.Printf("Variable Value: %f\n", variables[x].FloatValue)
		} else if(variables[x].Type == VARIABLE_TYPE_STRING) {
			fmt.Printf("Variable Value: %s\n", variables[x].StringValue)
		} else if(variables[x].Type == VARIABLE_TYPE_INTEGER) {
			fmt.Printf("Variable Value: %d\n", variables[x].IntegerValue)
		} else {
			//Nil type
			fmt.Println("Variable Value: Nil")
		}

		fmt.Printf("Variable Scope: %s\n", variables[x].ScopeName)
		fmt.Printf("Variable Type: %s\n", VARIABLE_TYPES_STRING[variables[x].Type])
		
		if(variables[x].IsConstant) {
			fmt.Println("Variable Constant: Yes")
		} else {
			fmt.Println("Variable Constant: No")
		}

		fmt.Printf("====================================\n")
	}
}

func isVariableExists(token Token, globalVariableArray []Variable, scopeName string) (bool, int) {

	for x := 0; x < len(globalVariableArray); x++ {
		if(globalVariableArray[x].Name == token.Value && globalVariableArray[x].ScopeName == scopeName) {
			return true, x
		}
	}

	return false, 0
}

func convertVariableToToken(token Token, variables []Variable, scopeName string) (Token, error) {

	isExists, indx := isVariableExists(token, variables, scopeName)

	if(scopeName == "main"){
		//if scope is main and not existing then raise an error right away
		if(!isExists) {
			return token, errors.New(SyntaxErrorMessage(token.Line, token.Column, "Variable doesn't exists '" + token.Value + "'", token.FileName))
		}
	} else {
		if(!isExists) {
			//if doesnt exists in function scope then check in main scope
			isExists, indx = isVariableExists(token, variables, "main")
			if(!isExists) {
				return token, errors.New(SyntaxErrorMessage(token.Line, token.Column, "Variable doesn't exists '" + token.Value + "'", token.FileName))
			}
		}
	}

	if(variables[indx].Type == VARIABLE_TYPE_INTEGER) {
		token.Type = TOKEN_TYPE_INTEGER
		token.Value = strconv.Itoa(variables[indx].IntegerValue)
	} else if(variables[indx].Type == VARIABLE_TYPE_STRING) {
		token.Type = TOKEN_TYPE_STRING
		token.Value = variables[indx].StringValue
	} else if(variables[indx].Type == VARIABLE_TYPE_FLOAT) {
		token.Type = TOKEN_TYPE_FLOAT
		token.Value = strconv.FormatFloat(variables[indx].FloatValue, 'f', -1, 64)
	} else {
		//Nil
		token.Type = TOKEN_TYPE_NONE
	}

	return token, nil
}

func defineConstantString(variableName string, variableValue string, globalVariableArray *[]Variable) {
	strVar := Variable{Name: variableName, ScopeName: "main", Type: VARIABLE_TYPE_STRING, IsConstant: true, StringValue: variableValue}
	*globalVariableArray = append(*globalVariableArray, strVar)
}

func initBuiltInVariables(globalVariableArray *[]Variable) {
	//add Nil Variable
	nilVar := Variable{Name: "Nil", ScopeName: "main", Type: VARIABLE_TYPE_NONE, IsConstant: true}
	*globalVariableArray = append(*globalVariableArray, nilVar)

	//define string constants
	defineConstantString("__AUTHOR__", TITIK_AUTHOR, globalVariableArray)
	defineConstantString("__VERSION_STRING__", TITIK_STRING_VERSION, globalVariableArray)
}