package main

import (
	"fmt"
	"errors"
	"strconv"
)

//variable types (basic)
const VARIABLE_TYPE_NONE int = 0
const VARIABLE_TYPE_INTEGER int = 1
const VARIABLE_TYPE_STRING int = 2
const VARIABLE_TYPE_FLOAT int = 3
const VARIABLE_TYPE_ARRAY int = 4

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

	if(!isExists) {
		return token, errors.New(SyntaxErrorMessage(token.Line, token.Column, "Variable doesn't exists '" + token.Value + "'", token.FileName))
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

func initBuiltInVariables(globalVariableArray *[]Variable) {
	//add Nil Variable
	nilVar := Variable{Name: "Nil", ScopeName: "main", Type: VARIABLE_TYPE_NONE, IsConstant: true}
	*globalVariableArray = append(*globalVariableArray, nilVar)
}