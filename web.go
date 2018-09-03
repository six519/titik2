package main

import (
	"fmt"
	"errors"
	"net/http"
	"strconv"
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
	thisRequest map[string]*http.Request
	staticURL string
	startedLine int
	startedColumn int
	startedFileName string
}

func (webObject *WebObject) Init(globalVariableArray *[]Variable, globalFunctionArray *[]Function, globalNativeVarList *[]string) {
	webObject.IsProcessing = false
	webObject.URLs = make(map[string]string)
	webObject.globalVariableArray = globalVariableArray
	webObject.globalFunctionArray = globalFunctionArray
	webObject.globalNativeVarList = globalNativeVarList
	webObject.thisWriter = make(map[string]http.ResponseWriter)
	webObject.thisRequest = make(map[string]*http.Request)
	webObject.staticURL = ""
}

func (webObject *WebObject) AddURL(key string, value string) {
	webObject.URLs[key] = value
	//DumpVariable(*webObject.globalVariableArray)
}

func (webObject *WebObject) handleHTTP(writer http.ResponseWriter, request *http.Request) {
	
	//default type to html
	writer.Header().Set("Content-Type", "text/html")

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
			InternalServerError(writer, "Error: Function handler doesn't exists on line number " + strconv.Itoa(webObject.startedLine) + " and column number " + strconv.Itoa(webObject.startedColumn) + ", Filename: " + webObject.startedFileName)
		} else {

			array := *webObject.globalFunctionArray

			if(array[funcIndex].ArgumentCount > 0) {
				InternalServerError(writer, "Error: Function argument is greater than zero on line number " + strconv.Itoa(webObject.startedLine) + " and column number " + strconv.Itoa(webObject.startedColumn) + ", Filename: " + webObject.startedFileName)
			} else {
				//execute titik function
				//newToken := Token{}
				thisScopeName := array[funcIndex].Name + generateRandomNumbers()
				webObject.thisWriter[thisScopeName] = writer
				webObject.thisRequest[thisScopeName] = request

				var thisGotReturn bool = false
				var thisReturnToken Token
				var thisNeedBreak bool = false
				var thisStackReference []Token
				
				//execute user defined function
				prsr := Parser{}
				parserErr := prsr.Parse(array[funcIndex].Tokens, webObject.globalVariableArray, webObject.globalFunctionArray, thisScopeName, webObject.globalNativeVarList, &thisGotReturn, &thisReturnToken, false, &thisNeedBreak, &thisStackReference, webObject)
		
				if(parserErr != nil) {
					InternalServerError(writer, parserErr.Error())
				}

				if(thisGotReturn) {
					//the function returns a value
					//newToken = thisReturnToken
					if(thisReturnToken.Type != TOKEN_TYPE_STRING) {
						InternalServerError(writer, "Error: Invalid return type on line number " + strconv.Itoa(webObject.startedLine) + " and column number " + strconv.Itoa(webObject.startedColumn) + ", Filename: " + webObject.startedFileName)
					} else {
						//redirect page if there's a return (the return should be string)
						http.Redirect(writer, request, thisReturnToken.Value, http.StatusFound)
					}
				}

				delete(webObject.thisWriter, thisScopeName) //cleanup map
				delete(webObject.thisRequest, thisScopeName) //cleanup map

			}
		}
	}

}

func Http_au_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[1].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {

		if(arguments[1].StringValue == (*webObject).staticURL) {
			*errMessage = errors.New("Error: URL " + arguments[1].StringValue + " already exists as static URL on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*webObject).AddURL(arguments[1].StringValue, arguments[0].StringValue)
		}
	}

	return ret
}

func Http_su_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[1].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		//(*webObject).AddURL(arguments[1].StringValue, arguments[0].StringValue)
		if(len((*webObject).staticURL) == 0) {
			_, ok := (*webObject).URLs[arguments[1].StringValue]
			if(ok) {
				*errMessage = errors.New("Error: URL " + arguments[1].StringValue + " already exists on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {
				(*webObject).staticURL = arguments[1].StringValue
				http.Handle(arguments[1].StringValue, http.StripPrefix(arguments[1].StringValue, http.FileServer(http.Dir(arguments[0].StringValue))))
			}
		} else {
			*errMessage = errors.New("Error: Static URL already exists on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		}
	}

	return ret
}

func Http_run_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {

		_, ok := (*webObject).URLs["/"]

		if(!ok) {
			*errMessage = errors.New("Error: Please handle web root's HTTP requests on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			return ret
		}

		(*webObject).IsProcessing = true
		(*webObject).scopeName = scopeName
		(*webObject).startedLine = line_number
		(*webObject).startedColumn = column_number
		(*webObject).startedFileName = file_name

		http.HandleFunc("/", (*webObject).handleHTTP)
		err := http.ListenAndServe(arguments[0].StringValue, nil)

		if(err != nil) {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		}
	}

	return ret
}

func Http_p_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {

		if(!(*webObject).IsProcessing) {
			*errMessage = errors.New("Error: Web server should be running on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			return ret
		}

		fmt.Fprintln((*webObject).thisWriter[scopeName], arguments[0].StringValue)
	}

	return ret
}

func Http_gm_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if(!(*webObject).IsProcessing) {
		*errMessage = errors.New("Error: Web server should be running on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		return ret
	}

	ret.StringValue = (*webObject).thisRequest[scopeName].Method

	return ret
}

func Http_gq_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, webObject *WebObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ARRAY}


	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		if(!(*webObject).IsProcessing) {
			*errMessage = errors.New("Error: Web server should be running on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			return ret
		}
	
		thisQuery := (*webObject).thisRequest[scopeName].URL.Query()
	
		val, ok := thisQuery[arguments[0].StringValue]

		if(ok) {
			for x := 0;x < len(val); x++ {
				funcReturn := FunctionReturn{Type: RET_TYPE_STRING}
				funcReturn.StringValue = val[x]
				ret.ArrayValue = append(ret.ArrayValue, funcReturn)
			}
		}
	}

	return ret
}