package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"io"
)

func M5_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		md5_hash := md5.New()
		io.WriteString(md5_hash, arguments[0].StringValue)
		ret.StringValue = hex.EncodeToString(md5_hash.Sum(nil))
	}

	return ret
}

func S1_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		sha1_hash := sha1.New()
		io.WriteString(sha1_hash, arguments[0].StringValue)
		ret.StringValue = hex.EncodeToString(sha1_hash.Sum(nil))
	}

	return ret
}

func S256_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		sha256_hash := sha256.New()
		io.WriteString(sha256_hash, arguments[0].StringValue)
		ret.StringValue = hex.EncodeToString(sha256_hash.Sum(nil))
	}

	return ret
}

func S512_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		sha512_hash := sha512.New()
		io.WriteString(sha512_hash, arguments[0].StringValue)
		ret.StringValue = hex.EncodeToString(sha512_hash.Sum(nil))
	}

	return ret
}

func B64e_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		ret.StringValue = base64.StdEncoding.EncodeToString([]byte(escapeString(arguments[0].StringValue)))
	}

	return ret
}

func B64d_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		data, err := base64.StdEncoding.DecodeString(escapeString(arguments[0].StringValue))
		if err == nil {
			ret.StringValue = string(data)
		}
	}

	return ret
}

func Rsae_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		block, _ := pem.Decode([]byte(escapeString(arguments[1].StringValue)))
		if block != nil {
			pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err == nil {
				pub := pubInterface.(*rsa.PublicKey)
				encrypted_data, err2 := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(escapeString(arguments[0].StringValue)))
				if err2 == nil {
					ret.StringValue = string(encrypted_data)
				}
			}
		}
	}

	return ret
}
