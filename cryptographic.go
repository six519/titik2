package main

import (
	"errors"
	"strconv"
	"crypto/md5"
	"io"
	"encoding/hex"
)

func M5_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		md5_hash := md5.New()
		io.WriteString(md5_hash, arguments[0].StringValue)
		ret.StringValue = hex.EncodeToString(md5_hash.Sum(nil))
	}

	return ret
}