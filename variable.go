package main

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
)

// variable types (basic)
const (
	VARIABLE_TYPE_NONE = iota
	VARIABLE_TYPE_INTEGER
	VARIABLE_TYPE_STRING
	VARIABLE_TYPE_FLOAT
	VARIABLE_TYPE_ARRAY
	VARIABLE_TYPE_ASSOCIATIVE_ARRAY
	VARIABLE_TYPE_BOOLEAN
)

var VARIABLE_TYPES_STRING = []string{
	"VARIABLE_TYPE_NONE",
	"VARIABLE_TYPE_INTEGER",
	"VARIABLE_TYPE_STRING",
	"VARIABLE_TYPE_FLOAT",
	"VARIABLE_TYPE_ARRAY",
	"VARIABLE_TYPE_ASSOCIATIVE_ARRAY",
	"VARIABLE_TYPE_BOOLEAN",
}

type Variable struct {
	Name                  string
	ScopeName             string
	Type                  int
	StringValue           string
	IntegerValue          int
	FloatValue            float64
	BooleanValue          bool
	IsConstant            bool
	ArrayValue            []Variable
	AssociativeArrayValue map[string]*Variable
}

func cleanupVariables(variables *[]Variable, scopeName string) {
	varCount := 0
	//count the variables first
	for x := 0; x < len(*variables); x++ {
		if (*variables)[x].ScopeName == scopeName {
			varCount = varCount + 1
		}
	}
	if varCount > 0 {
		for true {

			for x := 0; x < len(*variables); x++ {
				if (*variables)[x].ScopeName == scopeName {
					copy((*variables)[x:], (*variables)[x+1:])
					//(*variables)[len((*variables))-1] = nil
					(*variables) = (*variables)[:len((*variables))-1]
					varCount = varCount - 1
					break
				}
			}

			if varCount == 0 {
				break
			}
		}
	}
}

func DumpVariable(variables []Variable) {
	fmt.Printf("====================================\n")

	for x := 0; x < len(variables); x++ {
		fmt.Printf("Variable Name: %s\n", variables[x].Name)

		//NOTE: ASSUME VARIABLE TYPE AS FLOAT AND INTEGER FOR NOW (TEMPORARY)
		if variables[x].Type == VARIABLE_TYPE_FLOAT {
			fmt.Printf("Variable Value: %f\n", variables[x].FloatValue)
		} else if variables[x].Type == VARIABLE_TYPE_STRING {
			fmt.Printf("Variable Value: %s\n", variables[x].StringValue)
		} else if variables[x].Type == VARIABLE_TYPE_INTEGER {
			fmt.Printf("Variable Value: %d\n", variables[x].IntegerValue)
		} else if variables[x].Type == VARIABLE_TYPE_BOOLEAN {
			fmt.Printf("Variable Value: %v\n", variables[x].BooleanValue)
		} else if variables[x].Type == VARIABLE_TYPE_ARRAY {
			strVal := ""

			for x2 := 0; x2 < len(variables[x].ArrayValue); x2++ {
				//strVal = strVal + variables[x].Array[x2].Value
				if variables[x].ArrayValue[x2].Type == VARIABLE_TYPE_FLOAT {
					strVal = strVal + strconv.FormatFloat(variables[x].ArrayValue[x2].FloatValue, 'f', -1, 64)
				} else if variables[x].ArrayValue[x2].Type == VARIABLE_TYPE_STRING {
					strVal = strVal + variables[x].ArrayValue[x2].StringValue
				} else if variables[x].ArrayValue[x2].Type == VARIABLE_TYPE_INTEGER {
					strVal = strVal + strconv.Itoa(variables[x].ArrayValue[x2].IntegerValue)
				} else if variables[x].ArrayValue[x2].Type == VARIABLE_TYPE_BOOLEAN {
					if variables[x].ArrayValue[x2].BooleanValue {
						strVal = strVal + "true"
					} else {
						strVal = strVal + "false"
					}
				} else {
					strVal = strVal + "Nil"
				}

				if (x2 + 1) != len(variables[x].ArrayValue) {
					strVal = strVal + " , "
				}
			}

			fmt.Printf("Variable Value: [ %s ]\n", strVal)
		} else {
			//Nil type
			fmt.Println("Variable Value: Nil")
		}

		fmt.Printf("Variable Scope: %s\n", variables[x].ScopeName)
		fmt.Printf("Variable Type: %s\n", VARIABLE_TYPES_STRING[variables[x].Type])

		if variables[x].IsConstant {
			fmt.Println("Variable Constant: Yes")
		} else {
			fmt.Println("Variable Constant: No")
		}

		fmt.Printf("====================================\n")
	}
}

func isVariableExists(token Token, globalVariableArray []Variable, scopeName string) (bool, int) {

	for x := 0; x < len(globalVariableArray); x++ {
		if globalVariableArray[x].Name == token.Value && globalVariableArray[x].ScopeName == scopeName {
			return true, x
		}
	}

	return false, 0
}

func isGlobalVariableExists(varToCheck string, globalVariableArray []Variable) (bool, int) {

	for x := 0; x < len(globalVariableArray); x++ {
		if globalVariableArray[x].Name == varToCheck && globalVariableArray[x].ScopeName == "main" {
			return true, x
		}
	}

	return false, 0
}

func isSystemVariable(name string, globalNativeVarList []string) bool {

	for x := 0; x < len(globalNativeVarList); x++ {
		if name == globalNativeVarList[x] {
			return true
		}
	}

	return false
}

func convertTokenToBool(token Token) bool {
	if token.Value == "true" {
		return true
	}
	return false
}

func convertVariableToToken(token Token, variables []Variable, scopeName string) (Token, error) {

	isExists, indx := isVariableExists(token, variables, scopeName)

	if scopeName == "main" {
		//if scope is main and not existing then raise an error right away
		if !isExists {
			return token, errors.New(SyntaxErrorMessage(token.Line, token.Column, "Variable doesn't exists '"+token.Value+"'", token.FileName))
		}
	} else {
		if !isExists {
			//if doesnt exists in function scope then check in main scope
			isExists, indx = isVariableExists(token, variables, "main")
			if !isExists {
				return token, errors.New(SyntaxErrorMessage(token.Line, token.Column, "Variable doesn't exists '"+token.Value+"'", token.FileName))
			}
		}
	}

	if variables[indx].Type == VARIABLE_TYPE_INTEGER {
		token.Type = TOKEN_TYPE_INTEGER
		token.Value = strconv.Itoa(variables[indx].IntegerValue)
	} else if variables[indx].Type == VARIABLE_TYPE_STRING {
		token.Type = TOKEN_TYPE_STRING
		token.Value = variables[indx].StringValue
	} else if variables[indx].Type == VARIABLE_TYPE_FLOAT {
		token.Type = TOKEN_TYPE_FLOAT
		token.Value = strconv.FormatFloat(variables[indx].FloatValue, 'f', -1, 64)
	} else if variables[indx].Type == VARIABLE_TYPE_BOOLEAN {
		token.Type = TOKEN_TYPE_BOOLEAN
		if variables[indx].BooleanValue {
			//true
			token.Value = "true"
		} else {
			//false
			token.Value = "false"
		}
	} else if variables[indx].Type == VARIABLE_TYPE_ASSOCIATIVE_ARRAY {
		token.Type = TOKEN_TYPE_ASSOCIATIVE_ARRAY
		token.OtherInt = len(variables[indx].AssociativeArrayValue)
		token.AssociativeArray = make(map[string]Token)

		for k, v := range variables[indx].AssociativeArrayValue {
			newToken := Token{}
			if v.Type == VARIABLE_TYPE_INTEGER {
				newToken.Type = TOKEN_TYPE_INTEGER
				newToken.Value = strconv.Itoa(v.IntegerValue)
			} else if v.Type == VARIABLE_TYPE_STRING {
				newToken.Type = TOKEN_TYPE_STRING
				newToken.Value = v.StringValue
			} else if v.Type == VARIABLE_TYPE_FLOAT {
				newToken.Type = TOKEN_TYPE_FLOAT
				newToken.Value = strconv.FormatFloat(v.FloatValue, 'f', -1, 64)
			} else if v.Type == VARIABLE_TYPE_BOOLEAN {
				newToken.Type = TOKEN_TYPE_BOOLEAN
				if v.BooleanValue {
					newToken.Value = "true"
				} else {
					newToken.Value = "false"
				}
			} else {
				newToken.Type = TOKEN_TYPE_NONE
			}

			token.AssociativeArray[k] = newToken
		}

	} else if variables[indx].Type == VARIABLE_TYPE_ARRAY {
		token.Type = TOKEN_TYPE_ARRAY
		token.OtherInt = len(variables[indx].ArrayValue)
		for x := 0; x < len(variables[indx].ArrayValue); x++ {
			newToken := Token{}
			if variables[indx].ArrayValue[x].Type == VARIABLE_TYPE_INTEGER {
				newToken.Type = TOKEN_TYPE_INTEGER
				newToken.Value = strconv.Itoa(variables[indx].ArrayValue[x].IntegerValue)
			} else if variables[indx].ArrayValue[x].Type == VARIABLE_TYPE_STRING {
				newToken.Type = TOKEN_TYPE_STRING
				newToken.Value = variables[indx].ArrayValue[x].StringValue
			} else if variables[indx].ArrayValue[x].Type == VARIABLE_TYPE_FLOAT {
				newToken.Type = TOKEN_TYPE_FLOAT
				newToken.Value = strconv.FormatFloat(variables[indx].ArrayValue[x].FloatValue, 'f', -1, 64)
			} else if variables[indx].ArrayValue[x].Type == VARIABLE_TYPE_BOOLEAN {
				newToken.Type = TOKEN_TYPE_BOOLEAN
				if variables[indx].ArrayValue[x].BooleanValue {
					newToken.Value = "true"
				} else {
					newToken.Value = "false"
				}
			} else {
				newToken.Type = TOKEN_TYPE_NONE
			}

			token.Array = append(token.Array, newToken)
		}
	} else {
		//Nil
		token.Type = TOKEN_TYPE_NONE
	}

	return token, nil
}

func defineConstantString(variableName string, variableValue string, globalVariableArray *[]Variable, globalNativeVarList *[]string) {
	strVar := Variable{Name: variableName, ScopeName: "main", Type: VARIABLE_TYPE_STRING, IsConstant: true, StringValue: variableValue}
	*globalVariableArray = append(*globalVariableArray, strVar)
	*globalNativeVarList = append(*globalNativeVarList, variableName)
}

func defineConstantBoolean(variableName string, variableValue bool, globalVariableArray *[]Variable, globalNativeVarList *[]string) {
	boolVar := Variable{Name: variableName, ScopeName: "main", Type: VARIABLE_TYPE_BOOLEAN, IsConstant: true, BooleanValue: variableValue}
	*globalVariableArray = append(*globalVariableArray, boolVar)
	*globalNativeVarList = append(*globalNativeVarList, variableName)
}

func defineConstantInteger(variableName string, variableValue int, globalVariableArray *[]Variable, globalNativeVarList *[]string) {
	intVar := Variable{Name: variableName, ScopeName: "main", Type: VARIABLE_TYPE_INTEGER, IsConstant: true, IntegerValue: variableValue}
	*globalVariableArray = append(*globalVariableArray, intVar)
	*globalNativeVarList = append(*globalNativeVarList, variableName)
}

func initBuiltInVariables(globalVariableArray *[]Variable, globalNativeVarList *[]string) {
	//add Nil Variable
	nilVar := Variable{Name: "Nil", ScopeName: "main", Type: VARIABLE_TYPE_NONE, IsConstant: true}
	*globalVariableArray = append(*globalVariableArray, nilVar)
	*globalNativeVarList = append(*globalNativeVarList, "Nil")

	//define string constants
	defineConstantString("__AUTHOR__", TITIK_AUTHOR, globalVariableArray, globalNativeVarList)
	defineConstantString("__VERSION_STRING__", TITIK_STRING_VERSION, globalVariableArray, globalNativeVarList)
	defineConstantString("__OS__", runtime.GOOS, globalVariableArray, globalNativeVarList)

	//define boolean constants
	defineConstantBoolean("T", true, globalVariableArray, globalNativeVarList)
	defineConstantBoolean("F", false, globalVariableArray, globalNativeVarList)

	if SDL_ENABLED {
		//sdl init
		defineConstantInteger("S_E", 0, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_T", 1, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_A", 2, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_V", 3, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_J", 4, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_H", 5, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_G", 6, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E", 7, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_N", 8, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_S", 9, globalVariableArray, globalNativeVarList)

		//sdl window positions
		defineConstantInteger("S_WP_U", 0, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_WP_C", 1, globalVariableArray, globalNativeVarList)

		//sdl window flags
		defineConstantInteger("S_W_F", 0, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_S", 1, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_H", 2, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_B", 3, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_R", 4, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_MI", 5, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_MA", 6, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_FD", 7, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_A", 8, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_T", 9, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_W_P", 10, globalVariableArray, globalNativeVarList)

		//sdl events
		defineConstantInteger("S_E_F", 0, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_Q", 1, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_D", 2, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_W", 3, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_S", 4, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_KD", 5, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_KU", 6, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_MM", 7, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_MBD", 8, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_MBU", 9, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_E_MW", 10, globalVariableArray, globalNativeVarList)

		//mix init
		defineConstantInteger("SM_F", 0, globalVariableArray, globalNativeVarList)
		defineConstantInteger("SM_MO", 1, globalVariableArray, globalNativeVarList)
		defineConstantInteger("SM_MP", 2, globalVariableArray, globalNativeVarList)
		defineConstantInteger("SM_O", 3, globalVariableArray, globalNativeVarList)

		//mix defaults
		defineConstantInteger("SM_D_FR", 0, globalVariableArray, globalNativeVarList)
		defineConstantInteger("SM_D_FO", 1, globalVariableArray, globalNativeVarList)
		defineConstantInteger("SM_D_CHA", 2, globalVariableArray, globalNativeVarList)
		defineConstantInteger("SM_D_CHU", 3, globalVariableArray, globalNativeVarList)

		//sdl renderer flags
		defineConstantInteger("S_R_S", 0, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_R_A", 1, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_R_P", 2, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_R_T", 3, globalVariableArray, globalNativeVarList)

		//sdl scan codes
		defineConstantInteger("S_SC_RE", 0, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_SC_E", 1, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_SC_S", 2, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_SC_RI", 3, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_SC_L", 4, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_SC_D", 5, globalVariableArray, globalNativeVarList)
		defineConstantInteger("S_SC_U", 6, globalVariableArray, globalNativeVarList)
	}

}
