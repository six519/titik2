package main

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

type Execute func([]FunctionArgument, FunctionReturn) FunctionReturn

type Function struct {
	Name string
	IsSystem bool
	Tokens []Token
	Run Execute
}