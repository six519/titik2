package main

import (
	"errors"
	"strconv"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"fmt"
)

func Mysql_set_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		//database name
		*errMessage = errors.New("Error: Parameter 4 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[1].Type != ARG_TYPE_STRING) {
		//host
		*errMessage = errors.New("Error: Parameter 3 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[2].Type != ARG_TYPE_STRING) {
		//password
		*errMessage = errors.New("Error: Parameter 2 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else if(arguments[3].Type != ARG_TYPE_STRING) {
		//username
		*errMessage = errors.New("Error: Parameter 1 must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		(*globalSettings).stringSettings["MYSQL_USER"] = arguments[3].StringValue
		(*globalSettings).stringSettings["MYSQL_PASSWORD"] = arguments[2].StringValue
		(*globalSettings).stringSettings["MYSQL_HOST"] = arguments[1].StringValue
		(*globalSettings).stringSettings["MYSQL_DATABASE"] = arguments[0].StringValue
	}

	return ret
}

func Mysql_q_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: false}

	if(arguments[0].Type != ARG_TYPE_STRING) {
		*errMessage = errors.New("Error: Parameter must be a string type on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
	} else {
		strQuery := strings.ToLower(arguments[0].StringValue)
		strQuery = strings.Trim(strQuery, " ")
		strQuery = strings.Trim(strQuery, "\n")

		strSlices := strings.Split(strQuery, " ")

		db, err := sql.Open("mysql", (*globalSettings).stringSettings["MYSQL_USER"] + ":" + (*globalSettings).stringSettings["MYSQL_PASSWORD"] + "@" + (*globalSettings).stringSettings["MYSQL_HOST"] + "/" + (*globalSettings).stringSettings["MYSQL_DATABASE"])
		defer db.Close()
	
		if(err != nil) {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			if(strSlices[0] == "select") {
				//select query
			} else {
				//other query (insert/delete)
				stmt, err2 := db.Prepare(arguments[0].StringValue)

				if(err2 != nil) {
					*errMessage = errors.New("Error: " + err2.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
				} else {
					_, err3 := stmt.Exec()
		
					if(err3 != nil) {
						*errMessage = errors.New("Error: " + err3.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
					} else {
						//success execution
						ret.BooleanValue = true
					}
				}
			}
		}
	}

	return ret
}