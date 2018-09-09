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

		(*globalSettings).stringSettings[scopeName] = make(map[string]string) //TODO: CLEAN UP (MAYBE AFTER FUNCTION CALL?)

		(*globalSettings).stringSettings[scopeName]["MYSQL_USER"] = arguments[3].StringValue
		(*globalSettings).stringSettings[scopeName]["MYSQL_PASSWORD"] = arguments[2].StringValue
		(*globalSettings).stringSettings[scopeName]["MYSQL_HOST"] = arguments[1].StringValue
		(*globalSettings).stringSettings[scopeName]["MYSQL_DATABASE"] = arguments[0].StringValue
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

		db, err := sql.Open("mysql", (*globalSettings).stringSettings[scopeName]["MYSQL_USER"] + ":" + (*globalSettings).stringSettings[scopeName]["MYSQL_PASSWORD"] + "@" + (*globalSettings).stringSettings[scopeName]["MYSQL_HOST"] + "/" + (*globalSettings).stringSettings[scopeName]["MYSQL_DATABASE"])
		defer db.Close()
	
		if(err != nil) {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			if(strSlices[0] == "select") {
				//select query
				rows, err := db.Query(arguments[0].StringValue)

				if(err != nil) {
					*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
				} else {
					columns, err2 := rows.Columns()
		
					if(err2 != nil) {
						*errMessage = errors.New("Error: " + err2.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
					} else {
						values := make([]sql.RawBytes, len(columns))
						scanArgs := make([]interface{}, len(values))
						for i := range values {
							scanArgs[i] = &values[i]
						}
		
						for rows.Next() {
							err3 := rows.Scan(scanArgs...)
							if err3 != nil {
								ret.BooleanValue = false
								*errMessage = errors.New("Error: " + err3.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
								break
							} else {

								var value string
								for i, col := range values {
		
									if col == nil {
										value = ""
									} else {
										value = string(col)
									}

									_, ok := (*globalSettings).mySQLResults[scopeName]

									if(!ok) {
										(*globalSettings).mySQLResults[scopeName] = make(map[string][]string) //TODO: THIS SHOULD BE CLEANUP (AFTER FUNCTION CALL?)
									}

									(*globalSettings).mySQLResults[scopeName][columns[i]] = append((*globalSettings).mySQLResults[scopeName][columns[i]], value)
									
								}

								ret.BooleanValue = true		
							}
						}
					}
				}

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

func Mysql_cr_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_NONE}

	delete((*globalSettings).stringSettings, scopeName)
	delete((*globalSettings).mySQLResults, scopeName)

	return ret
}