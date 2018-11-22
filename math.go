package main
import (
	"math"
	"strconv"
	"errors"
)

func Abs_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_FLOAT, FloatValue: 0}

	if(arguments[0].Type != ARG_TYPE_FLOAT && arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be a float or integer type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		var output float64

		if(arguments[0].Type == ARG_TYPE_FLOAT) {
			output = math.Abs(arguments[0].FloatValue)
		} else {
			output = math.Abs(float64(arguments[0].IntegerValue))
		}

		/*
		if(math.IsNaN(output)) {
			ret.Type = RET_TYPE_NONE
		} else {
			ret.FloatValue = output
		}
		*/

		ret.FloatValue = output

	}

	return ret
}

func Acs_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_FLOAT, FloatValue: 0}

	if(arguments[0].Type != ARG_TYPE_FLOAT && arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be a float or integer type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		var output float64

		if(arguments[0].Type == ARG_TYPE_FLOAT) {
			output = math.Acos(arguments[0].FloatValue)
		} else {
			output = math.Acos(float64(arguments[0].IntegerValue))
		}

		if(math.IsNaN(output)) {
			ret.Type = RET_TYPE_NONE
		} else {
			ret.FloatValue = output
		}

		ret.FloatValue = output

	}

	return ret
}

func Acsh_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_FLOAT, FloatValue: 0}

	if(arguments[0].Type != ARG_TYPE_FLOAT && arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be a float or integer type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		var output float64

		if(arguments[0].Type == ARG_TYPE_FLOAT) {
			output = math.Acosh(arguments[0].FloatValue)
		} else {
			output = math.Acosh(float64(arguments[0].IntegerValue))
		}

		if(math.IsNaN(output)) {
			ret.Type = RET_TYPE_NONE
		} else {
			ret.FloatValue = output
		}

		ret.FloatValue = output

	}

	return ret
}

func Asn_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_FLOAT, FloatValue: 0}

	if(arguments[0].Type != ARG_TYPE_FLOAT && arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be a float or integer type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		var output float64

		if(arguments[0].Type == ARG_TYPE_FLOAT) {
			output = math.Asin(arguments[0].FloatValue)
		} else {
			output = math.Asin(float64(arguments[0].IntegerValue))
		}

		if(math.IsNaN(output)) {
			ret.Type = RET_TYPE_NONE
		} else {
			ret.FloatValue = output
		}

		ret.FloatValue = output

	}

	return ret
}

func Asnh_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_FLOAT, FloatValue: 0}

	if(arguments[0].Type != ARG_TYPE_FLOAT && arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be a float or integer type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		var output float64

		if(arguments[0].Type == ARG_TYPE_FLOAT) {
			output = math.Asinh(arguments[0].FloatValue)
		} else {
			output = math.Asinh(float64(arguments[0].IntegerValue))
		}

		if(math.IsNaN(output)) {
			ret.Type = RET_TYPE_NONE
		} else {
			ret.FloatValue = output
		}

		ret.FloatValue = output

	}

	return ret
}

func Atn_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_FLOAT, FloatValue: 0}

	if(arguments[0].Type != ARG_TYPE_FLOAT && arguments[0].Type != ARG_TYPE_INTEGER) {
		*errMessage = errors.New("Error: Parameter must be a float or integer type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		var output float64

		if(arguments[0].Type == ARG_TYPE_FLOAT) {
			output = math.Atan(arguments[0].FloatValue)
		} else {
			output = math.Atan(float64(arguments[0].IntegerValue))
		}

		if(math.IsNaN(output)) {
			ret.Type = RET_TYPE_NONE
		} else {
			ret.FloatValue = output
		}

		ret.FloatValue = output

	}

	return ret
}