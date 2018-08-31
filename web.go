package main

import (
	"fmt"
	"errors"
	"net/http"
)

func InternalServerError(writer http.ResponseWriter, msg string) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte(msg))
}

type WebObject struct {
	IsProcessing bool
	URLs map[string]string
	globalVariableArray *[]Variable
	globalFunctionArray *[]Function
	globalNativeVarList *[]string
	scopeName string
	thisWriter map[string]http.ResponseWriter
	thisRequest *http.Request
}

func (webObject *WebObject) Init(globalVariableArray *[]Variable, globalFunctionArray *[]Function, globalNativeVarList *[]string) {
	webObject.IsProcessing = false
	webObject.URLs = make(map[string]string)
	webObject.globalVariableArray = globalVariableArray
	webObject.globalFunctionArray = globalFunctionArray
	webObject.globalNativeVarList = globalNativeVarList
	webObject.thisWriter = make(map[string]http.ResponseWriter)
}

func (webObject *WebObject) AddURL(key string, value string) {
	webObject.URLs[key] = value
	//DumpVariable(*webObject.globalVariableArray)
}

func (webObject *WebObject) handleHTTP(writer http.ResponseWriter, request *http.Request) {
	
	//webObject.thisWriter = writer
	webObject.thisRequest = request

	thisPath := request.URL.Path[1:]

	if(thisPath == "") {
		thisPath = "/"
	}
	_, ok := webObject.URLs[thisPath]

	if(!ok) {
		//no handler given
		http.NotFound(writer, request)
	} else {
		//try to load titik function
		t := Token{Value: webObject.URLs[thisPath]}
		isExists, funcIndex := isFunctionExists(t, *webObject.globalFunctionArray)

		if(!isExists) {
			InternalServerError(writer, "Error: Function handler doesn't exists")
		} else {

			array := *webObject.globalFunctionArray

			if(array[funcIndex].ArgumentCount > 0) {
				InternalServerError(writer, "Error: Function argument is greater than zero")
			} else {
				//execute titik function
				//newToken := Token{}
				thisScopeName := array[funcIndex].Name + generateRandomNumbers()
				webObject.thisWriter[thisScopeName] = writer

				var thisGotReturn bool = false
				var thisReturnToken Token
				var thisNeedBreak bool = false
				var thisStackReference []Token
				
				//execute user defined function
				prsr := Parser{}
				parserErr := prsr.Parse(array[funcIndex].Tokens, webObject.globalVariableArray, webObject.globalFunctionArray, thisScopeName, webObject.globalNativeVarList, &thisGotReturn, &thisReturnToken, false, &thisNeedBreak, &thisStackReference, webObject)
		
				if(parserErr != nil) {
					InternalServerError(writer, "Error: " + parserErr.Error())
				}

				if(thisGotReturn) {
					//the function returns a value
					//newToken = thisReturnToken
					if(thisReturnToken.Type != TOKEN_TYPE_STRING) {
						InternalServerError(writer, "Error: Invalid return type")
					} else {
						//redirect page if there's a return (the return should be string)
						http.Redirect(writer, request, thisReturnToken.Value, http.StatusFound)
					}
				}

				delete(webObject.thisWriter, thisScopeName) //cleanup map

			}
		}
	}

}

func Http_au_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 2 must be a string type")
	} else if(arguments[1].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type")
	} else {
		(*webObject).AddURL(arguments[1].StringValue, arguments[0].StringValue)
	}

	return ret
}

func Http_run_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type")
	} else {

		_, ok := (*webObject).URLs["/"]

		if(!ok) {
			*errMessage = errors.New("Error: Please handle web root's HTTP requests")
			return ret
		}

		(*webObject).IsProcessing = true
		(*webObject).scopeName = scopeName

		http.HandleFunc("/", (*webObject).handleHTTP)
		*errMessage = http.ListenAndServe(arguments[0].StringValue, nil)
	}

	return ret
}

func Http_p_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type")
	} else {

		if(!(*webObject).IsProcessing) {
			*errMessage = errors.New("Error: Web server should be running")
			return ret
		}

		fmt.Fprintln((*webObject).thisWriter[scopeName], arguments[0].StringValue)
	}

	return ret
}