package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func InternalServerError(writer http.ResponseWriter, msg string) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte(msg))
}

type WebObject struct {
	IsProcessing    bool
	URLs            map[string]string
	staticURLs      map[string]string
	globalSettings  *GlobalSettingsObject
	scopeName       string
	thisWriter      map[string]http.ResponseWriter
	thisRequest     map[string]*http.Request
	startedLine     int
	startedColumn   int
	startedFileName string
}

func (webObject *WebObject) Init(globalSettings *GlobalSettingsObject) {
	webObject.IsProcessing = false
	webObject.URLs = make(map[string]string)
	webObject.staticURLs = make(map[string]string)
	webObject.globalSettings = globalSettings
	webObject.thisWriter = make(map[string]http.ResponseWriter)
	webObject.thisRequest = make(map[string]*http.Request)
}

func (webObject *WebObject) AddURL(key string, value string) {
	webObject.URLs[key] = value
	//DumpVariable(*webObject.globalVariableArray)
}

func (webObject *WebObject) AddStaticURL(key string, value string) {
	webObject.staticURLs[key] = value
}

func (webObject *WebObject) handleHTTP(writer http.ResponseWriter, request *http.Request) {

	//default type to html
	writer.Header().Set("Content-Type", "text/html")
	//always process form POST
	request.ParseMultipartForm(32 << 20) //handle 32M upload (for now)

	thisPath := request.URL.Path[1:]

	if thisPath == "" {
		thisPath = "/"
	}
	_, ok := webObject.URLs[thisPath]

	if !ok {
		//no handler given
		http.NotFound(writer, request)
	} else {
		//try to load titik function
		t := Token{Value: webObject.URLs[thisPath]}
		isExists, funcIndex := isFunctionExists(t, *webObject.globalSettings.globalFunctionArray)

		if !isExists {
			InternalServerError(writer, "Error: Function handler doesn't exists on line number "+strconv.Itoa(webObject.startedLine)+" and column number "+strconv.Itoa(webObject.startedColumn)+", Filename: "+webObject.startedFileName)
		} else {

			array := *webObject.globalSettings.globalFunctionArray

			if array[funcIndex].ArgumentCount > 0 {
				InternalServerError(writer, "Error: Function argument is greater than zero on line number "+strconv.Itoa(webObject.startedLine)+" and column number "+strconv.Itoa(webObject.startedColumn)+", Filename: "+webObject.startedFileName)
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
				var getLastStackBool bool = false
				var lastStackBool bool = false

				//execute user defined function
				prsr := Parser{}
				parserErr := prsr.Parse(array[funcIndex].Tokens, webObject.globalSettings.globalVariableArray, webObject.globalSettings.globalFunctionArray, thisScopeName, webObject.globalSettings.globalNativeVarList, &thisGotReturn, &thisReturnToken, false, &thisNeedBreak, &thisStackReference, webObject.globalSettings, getLastStackBool, &lastStackBool)

				if parserErr != nil {
					InternalServerError(writer, parserErr.Error())
				}

				if thisGotReturn {
					//the function returns a value
					//newToken = thisReturnToken
					if thisReturnToken.Type != TOKEN_TYPE_STRING {
						InternalServerError(writer, "Error: Invalid return type on line number "+strconv.Itoa(webObject.startedLine)+" and column number "+strconv.Itoa(webObject.startedColumn)+", Filename: "+webObject.startedFileName)
					} else {
						//redirect page if there's a return (the return should be string)
						http.Redirect(writer, request, thisReturnToken.Value, http.StatusFound)
					}
				}

				delete(webObject.thisWriter, thisScopeName)  //cleanup map
				delete(webObject.thisRequest, thisScopeName) //cleanup map

			}
		}
	}

}

func Http_au_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		_, existing := (*globalSettings).webObject.staticURLs[arguments[1].StringValue]
		if existing {
			*errMessage = errors.New("Error: URL " + arguments[1].StringValue + " already exists as static URL on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			(*globalSettings).webObject.AddURL(arguments[1].StringValue, arguments[0].StringValue)
		}
	}

	return ret
}

func Http_su_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		_, ok := (*globalSettings).webObject.URLs[arguments[1].StringValue]
		if ok {
			*errMessage = errors.New("Error: URL " + arguments[1].StringValue + " already exists on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			_, ok = (*globalSettings).webObject.staticURLs[arguments[1].StringValue]
			if ok {
				*errMessage = errors.New("Error: URL " + arguments[1].StringValue + " already exists on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {
				(*globalSettings).webObject.AddStaticURL(arguments[1].StringValue, arguments[0].StringValue)
				http.Handle(arguments[1].StringValue, http.StripPrefix(arguments[1].StringValue, http.FileServer(http.Dir(arguments[0].StringValue))))
			}
		}
	}

	return ret
}

func Http_run_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		_, ok := (*globalSettings).webObject.URLs["/"]

		if !ok {
			*errMessage = errors.New("Error: Please handle web root's HTTP requests on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			return ret
		}

		(*globalSettings).webObject.IsProcessing = true
		(*globalSettings).webObject.scopeName = scopeName
		(*globalSettings).webObject.startedLine = line_number
		(*globalSettings).webObject.startedColumn = column_number
		(*globalSettings).webObject.startedFileName = file_name

		http.HandleFunc("/", (*globalSettings).webObject.handleHTTP)
		err := http.ListenAndServe(arguments[0].StringValue, nil)

		if err != nil {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		}
	}

	return ret
}

func Http_p_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if !(*globalSettings).webObject.IsProcessing {
			*errMessage = errors.New("Error: Web server should be running on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			return ret
		}

		fmt.Fprintln((*globalSettings).webObject.thisWriter[scopeName], arguments[0].StringValue)
	}

	return ret
}

func Http_h_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if !(*globalSettings).webObject.IsProcessing {
			*errMessage = errors.New("Error: Web server should be running on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			return ret
		}

		(*globalSettings).webObject.thisWriter[scopeName].Header().Set(arguments[1].StringValue, arguments[0].StringValue)
	}

	return ret
}

func Http_gm_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if !(*globalSettings).webObject.IsProcessing {
		*errMessage = errors.New("Error: Web server should be running on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		return ret
	}

	ret.StringValue = (*globalSettings).webObject.thisRequest[scopeName].Method

	return ret
}

func Http_gq_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ARRAY}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if !(*globalSettings).webObject.IsProcessing {
			*errMessage = errors.New("Error: Web server should be running on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			return ret
		}

		thisQuery := (*globalSettings).webObject.thisRequest[scopeName].URL.Query()

		val, ok := thisQuery[arguments[0].StringValue]

		if ok {
			for x := 0; x < len(val); x++ {
				funcReturn := FunctionReturn{Type: RET_TYPE_STRING}
				funcReturn.StringValue = val[x]
				ret.ArrayValue = append(ret.ArrayValue, funcReturn)
			}
		}
	}

	return ret
}

func Http_gfp_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ARRAY}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if !(*globalSettings).webObject.IsProcessing {
			*errMessage = errors.New("Error: Web server should be running on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			return ret
		}

		thisPostForm := (*globalSettings).webObject.thisRequest[scopeName].Form[arguments[0].StringValue]

		for x := 0; x < len(thisPostForm); x++ {
			funcReturn := FunctionReturn{Type: RET_TYPE_STRING}
			funcReturn.StringValue = thisPostForm[x]
			ret.ArrayValue = append(ret.ArrayValue, funcReturn)
		}
	}

	return ret
}

func Http_lt_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_ASSOCIATIVE_ARRAY) {
		//get all associative arrays
		//and convert all to string
		//and add to map
		var stringMap map[string]string
		stringMap = make(map[string]string)

		for k, v := range arguments[0].AssociativeArrayValue {
			if v.Type == ARG_TYPE_FLOAT {
				stringMap[k] = strconv.FormatFloat(v.FloatValue, 'f', -1, 64)
			} else if v.Type == ARG_TYPE_STRING {
				stringMap[k] = v.StringValue
			} else if v.Type == ARG_TYPE_INTEGER {
				stringMap[k] = strconv.Itoa(v.IntegerValue)
			} else if v.Type == ARG_TYPE_BOOLEAN {
				if v.BooleanValue {
					stringMap[k] = "true"
				} else {
					stringMap[k] = "false"
				}
			} else {
				stringMap[k] = ""
			}
		}

		t, err := template.ParseFiles(arguments[1].StringValue)

		if err != nil {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			t.Execute((*globalSettings).webObject.thisWriter[scopeName], stringMap)
		}

	}

	return ret
}

func Http_gp_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if !(*globalSettings).webObject.IsProcessing {
		*errMessage = errors.New("Error: Web server should be running on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		return ret
	}

	ret.StringValue = (*globalSettings).webObject.thisRequest[scopeName].URL.Path

	return ret
}

func Http_cr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ASSOCIATIVE_ARRAY}
	ret.AssociativeArrayValue = make(map[string]FunctionReturn)
	ret.AssociativeArrayValue["status"] = FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: 0}
	ret.AssociativeArrayValue["body"] = FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	//http_cr("POST", "http://test.com", {'Content-Type': 'application/json'}, '{"name":"ferdinand"}')
	//http_cr("POST", "http://test.com", {'Content-Type': 'application/x-www-form-urlencoded'}, {"age": 23})
	//http_cr("GET", "http://test.com", Nil, Nil)

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) {
		if arguments[1].Type != ARG_TYPE_ASSOCIATIVE_ARRAY && arguments[1].Type != ARG_TYPE_NONE {
			*errMessage = errors.New("Error: Parameter 3 must be a glossary/Nil type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else if arguments[0].Type != ARG_TYPE_ASSOCIATIVE_ARRAY && arguments[0].Type != ARG_TYPE_NONE && arguments[0].Type != ARG_TYPE_STRING {
			*errMessage = errors.New("Error: Parameter 4 must be a glossary/Nil/string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			var b io.Reader
			if arguments[0].Type == ARG_TYPE_ASSOCIATIVE_ARRAY {
				data := url.Values{}
				for k, v := range arguments[0].AssociativeArrayValue {
					if v.Type == ARG_TYPE_FLOAT {
						data.Set(k, strconv.FormatFloat(v.FloatValue, 'f', -1, 64))
					} else if v.Type == ARG_TYPE_STRING {
						data.Set(k, v.StringValue)
					} else if v.Type == ARG_TYPE_INTEGER {
						data.Set(k, strconv.Itoa(v.IntegerValue))
					} else if v.Type == ARG_TYPE_BOOLEAN {
						if v.BooleanValue {
							data.Set(k, "true")
						} else {
							data.Set(k, "false")
						}
					} else {
						data.Set(k, "Nil")
					}
				}

				b = bytes.NewBufferString(data.Encode())

			} else if arguments[0].Type == ARG_TYPE_STRING {
				var jsonReq = []byte(arguments[0].StringValue)
				b = bytes.NewBuffer(jsonReq)
			}

			req, err2 := http.NewRequest(arguments[3].StringValue, arguments[2].StringValue, b)

			//set http header
			if arguments[1].Type == ARG_TYPE_ASSOCIATIVE_ARRAY {
				for k, v := range arguments[0].AssociativeArrayValue {
					if v.Type == ARG_TYPE_FLOAT {
						req.Header.Set(k, strconv.FormatFloat(v.FloatValue, 'f', -1, 64))
					} else if v.Type == ARG_TYPE_STRING {
						req.Header.Set(k, v.StringValue)
					} else if v.Type == ARG_TYPE_INTEGER {
						req.Header.Set(k, strconv.Itoa(v.IntegerValue))
					} else if v.Type == ARG_TYPE_BOOLEAN {
						if v.BooleanValue {
							req.Header.Set(k, "true")
						} else {
							req.Header.Set(k, "false")
						}
					} else {
						req.Header.Set(k, "Nil")
					}
				}
			}

			if err2 != nil {
				*errMessage = errors.New("Error: " + err2.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}

			c := &http.Client{}
			resp, err3 := c.Do(req)

			if err3 != nil {
				*errMessage = errors.New("Error: " + err3.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}

			defer resp.Body.Close()
			r, err4 := ioutil.ReadAll(resp.Body)
			if err4 != nil {
				*errMessage = errors.New("Error: " + err4.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			}

			ret.AssociativeArrayValue["status"] = FunctionReturn{Type: RET_TYPE_INTEGER, IntegerValue: resp.StatusCode}
			ret.AssociativeArrayValue["body"] = FunctionReturn{Type: RET_TYPE_STRING, StringValue: string(r)}

		}
	}

	return ret
}
