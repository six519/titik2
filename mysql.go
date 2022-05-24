package main

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	//"fmt"
)

func Mysql_set_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_STRING, StringValue: ""}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 3, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 2, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		db, err := sql.Open("mysql", arguments[3].StringValue+":"+arguments[2].StringValue+"@tcp("+arguments[1].StringValue+")/"+arguments[0].StringValue)

		if err != nil {
			*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			connection_reference := "con_" + generateRandomNumbers()
			(*globalSettings).mySQLConnection[connection_reference] = db
			ret.StringValue = connection_reference
		}

	}

	return ret
}

func Mysql_q_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_BOOLEAN, BooleanValue: false}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		strQuery := strings.ToLower(arguments[0].StringValue)
		strQuery = strings.Trim(strQuery, " ")
		strQuery = strings.Trim(strQuery, "\n")

		strSlices := strings.Split(strQuery, " ")

		if (*globalSettings).mySQLConnection[arguments[1].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			if strSlices[0] == "select" {
				//select query
				rows, err := (*globalSettings).mySQLConnection[arguments[1].StringValue].Query(arguments[0].StringValue)

				if err != nil {
					*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
				} else {
					columns, err2 := rows.Columns()

					if err2 != nil {
						*errMessage = errors.New("Error: " + err2.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
					} else {
						values := make([]sql.RawBytes, len(columns))
						scanArgs := make([]interface{}, len(values))
						for i := range values {
							scanArgs[i] = &values[i]
						}

						_, ok := (*globalSettings).mySQLResults[arguments[1].StringValue]

						if ok {
							delete((*globalSettings).mySQLResults, arguments[1].StringValue)
						}

						(*globalSettings).mySQLResults[arguments[1].StringValue] = make(map[string][]string) //TODO: THIS SHOULD BE CLEANUP (AFTER FUNCTION CALL?)
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

									(*globalSettings).mySQLResults[arguments[1].StringValue][columns[i]] = append((*globalSettings).mySQLResults[arguments[1].StringValue][columns[i]], value)

								}

								ret.BooleanValue = true
							}
						}
					}
				}

			} else {
				//other query (insert/delete)
				stmt, err2 := (*globalSettings).mySQLConnection[arguments[1].StringValue].Prepare(arguments[0].StringValue)

				if err2 != nil {
					*errMessage = errors.New("Error: " + err2.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
				} else {
					_, err3 := stmt.Exec()

					if err3 != nil {
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

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		if (*globalSettings).mySQLConnection[arguments[0].StringValue] == nil {
			*errMessage = errors.New("Error: Uninitialized connection on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
		} else {
			err := (*globalSettings).mySQLConnection[arguments[0].StringValue].Close()

			if err != nil {
				*errMessage = errors.New("Error: " + err.Error() + " on line number " + strconv.Itoa(line_number) + " and column number " + strconv.Itoa(column_number) + ", Filename: " + file_name)
			} else {
				delete((*globalSettings).mySQLConnection, arguments[0].StringValue)

				_, ok := (*globalSettings).mySQLResults[arguments[0].StringValue]
				if ok {
					delete((*globalSettings).mySQLResults, arguments[0].StringValue)
				}

			}
		}
	}

	return ret
}

func Mysql_fa_execute(arguments []FunctionArgument, errMessage *error, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, globalSettings *GlobalSettingsObject, line_number int, column_number int, file_name string) FunctionReturn {
	ret := FunctionReturn{Type: RET_TYPE_ARRAY}

	if validateParameters(arguments, errMessage, line_number, column_number, file_name, 1, ARG_TYPE_STRING) &&
		validateParameters(arguments, errMessage, line_number, column_number, file_name, 0, ARG_TYPE_STRING) {
		_, ok := (*globalSettings).mySQLResults[arguments[1].StringValue]

		if ok {
			val, ok2 := (*globalSettings).mySQLResults[arguments[1].StringValue][arguments[0].StringValue]

			if ok2 {
				for x := 0; x < len(val); x++ {
					funcReturn := FunctionReturn{Type: RET_TYPE_STRING}
					funcReturn.StringValue = val[x]
					ret.ArrayValue = append(ret.ArrayValue, funcReturn)
				}
			}
		}
	}

	return ret
}
