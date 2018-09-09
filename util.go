package main

import (
	"strings"
)

type GlobalSettingsObject struct {
	webObject WebObject
	globalVariableArray *[]Variable
	globalFunctionArray *[]Function
	globalNativeVarList *[]string
	stringSettings map[string]map[string]string
	mySQLResults map[string]map[string][]string //NOTE: TEMPORARY ONLY
}

func (globalSettings *GlobalSettingsObject) Init(globalVariableArray *[]Variable, globalFunctionArray *[]Function, globalNativeVarList *[]string) {
	globalSettings.globalVariableArray = globalVariableArray
	globalSettings.globalFunctionArray = globalFunctionArray
	globalSettings.globalNativeVarList = globalNativeVarList
	globalSettings.webObject = WebObject{}
	globalSettings.webObject.Init(globalSettings)

	globalSettings.stringSettings = make(map[string]map[string]string) //TODO: NEED WAY TO CLEAN THIS UP //MAYBE END OF FUNCTION CALLS?
	globalSettings.mySQLResults = make(map[string]map[string][]string) //TODO: NEED WAY TO CLEAN THIS UP //MAYBE END OF FUNCTION CALLS?
}

func escapeString(str string) string {
	var retStr string

	//escape newline
	retStr = strings.Replace(str, "\\n", "\n", -1)

	//escape carriage return
	retStr = strings.Replace(retStr, "\\r", "\r", -1)

	//escape tab
	retStr = strings.Replace(retStr, "\\t", "\t", -1)

	//escape form feed
	retStr = strings.Replace(retStr, "\\f", "\f", -1)

	//escape bell
	retStr = strings.Replace(retStr, "\\a", "\a", -1)

	//escape backspace
	retStr = strings.Replace(retStr, "\\b", "\b", -1)

	return retStr
}

func unescapeString(str string) string {
	var retStr string

	//escape newline
	retStr = strings.Replace(str, "\n", "\\n", -1)

	//escape carriage return
	retStr = strings.Replace(retStr, "\r", "\\r", -1)

	//escape tab
	retStr = strings.Replace(retStr, "\t", "\\t", -1)

	//escape form feed
	retStr = strings.Replace(retStr, "\f", "\\f", -1)

	//escape bell
	retStr = strings.Replace(retStr, "\a", "\\a", -1)

	//escape backspace
	retStr = strings.Replace(retStr, "\b", "\\b", -1)

	return retStr
}