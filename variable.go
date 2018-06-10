package main

//variable types (basic)
const VARIABLE_TYPE_NONE int = 0
const VARIABLE_TYPE_INTEGER int = 1
const VARIABLE_TYPE_STRING int = 3
const VARIABLE_TYPE_FLOAT int = 4
const VARIABLE_TYPE_ARRAY int = 0

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

func isVariableExists(token Token, globalVariableArray []Variable, scopeName string) (bool, int) {

	for x := 0; x < len(globalVariableArray); x++ {
		if(globalVariableArray[x].Name == token.Value && globalVariableArray[x].ScopeName == scopeName) {
			return true, x
		}
	}

	return false, 0
}