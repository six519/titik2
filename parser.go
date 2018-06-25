package main

import (
	"errors"
	"strconv"
	//"fmt"
)

func expectedTokenTypes(token Token, tokenTypes ...int) error {
	isOK := false

	for x := 0; x < len(tokenTypes); x++ {
		if(token.Type == tokenTypes[x]) {
			isOK = true
			break
		}
	}

	if(!isOK) {
		return errors.New(SyntaxErrorMessage(token.Line, token.Column, "Invalid operand '" + token.Value + "'", token.FileName))
	}

	return nil
}

type Parser struct {
}

func (parser Parser) Parse(tokenArray []Token, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string) error {
	var tokensToEvaluate []Token
	operatorPrecedences := map[string] int{"function": 1, "=": 0, "+": 3, "-": 3, "/": 4, "*": 4, ",": 2} //operator order of precedences
	var operatorStack []Token
	var outputQueue []Token

	for x := 0; x < len(tokenArray); x++ {
		//ignore space, tab and comments
		if(tokenArray[x].Type != TOKEN_TYPE_SPACE && tokenArray[x].Type != TOKEN_TYPE_SINGLE_COMMENT && tokenArray[x].Type != TOKEN_TYPE_MULTI_COMMENT && tokenArray[x].Type != TOKEN_TYPE_CLOSE_MULTI_COMMENT && tokenArray[x].Type != TOKEN_TYPE_CLOSE_STRING && tokenArray[x].Type != TOKEN_TYPE_TAB) {
			if(tokenArray[x].Type == TOKEN_TYPE_NEWLINE) {
				//execute shunting yard
				if(len(tokensToEvaluate) > 0) {
					if(tokensToEvaluate[0].Type == TOKEN_TYPE_PLUS || tokensToEvaluate[0].Type == TOKEN_TYPE_MINUS || tokensToEvaluate[0].Type == TOKEN_TYPE_DIVIDE || tokensToEvaluate[0].Type == TOKEN_TYPE_MULTIPLY || tokensToEvaluate[0].Type == TOKEN_TYPE_EQUALS) {
						//syntax error if the first token is an operator
						return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '" + tokensToEvaluate[0].Value + "'", tokensToEvaluate[0].FileName))
					}
			
					if(tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_PLUS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_MINUS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_DIVIDE || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_MULTIPLY || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_EQUALS) {
						//syntax error if the last token is an operator
						return errors.New(SyntaxErrorMessage(tokensToEvaluate[len(tokensToEvaluate)-1].Line, tokensToEvaluate[len(tokensToEvaluate)-1].Column, "Unfinished operation", tokensToEvaluate[len(tokensToEvaluate)-1].FileName))
					}
			
					//shunting-yard
					for len(tokensToEvaluate) > 0 {
						currentToken := tokensToEvaluate[0]
						tokensToEvaluate = append(tokensToEvaluate[:0], tokensToEvaluate[1:]...) //pop the first element
						isValidToken := false
			
						if(currentToken.Type == TOKEN_TYPE_INTEGER || currentToken.Type == TOKEN_TYPE_FLOAT || currentToken.Type == TOKEN_TYPE_IDENTIFIER || currentToken.Type == TOKEN_TYPE_STRING) {
							//If it's a number or identifier, add it to queue, (ADD TOKEN_TYPE_KEYWORD AND string and other acceptable tokens LATER)
							outputQueue = append(outputQueue, currentToken)
							isValidToken = true
						}
			
						if(currentToken.Type == TOKEN_TYPE_PLUS || currentToken.Type == TOKEN_TYPE_MINUS || currentToken.Type == TOKEN_TYPE_DIVIDE || currentToken.Type == TOKEN_TYPE_MULTIPLY || currentToken.Type == TOKEN_TYPE_EQUALS || currentToken.Type == TOKEN_TYPE_FUNCTION || currentToken.Type == TOKEN_TYPE_COMMA) {
							//the token is operator
							for true {
								if(len(operatorStack) > 0) {

									if(currentToken.Type == TOKEN_TYPE_FUNCTION) {
										if(operatorPrecedences[operatorStack[len(operatorStack) - 1].Value] > operatorPrecedences["function"]) {
											outputQueue = append(outputQueue, operatorStack[len(operatorStack) - 1])
											operatorStack = operatorStack[:len(operatorStack)-1]
										} else {
											break
										}
									} else {
										if(operatorPrecedences[operatorStack[len(operatorStack) - 1].Value] > operatorPrecedences[currentToken.Value]) {
											outputQueue = append(outputQueue, operatorStack[len(operatorStack) - 1])
											operatorStack = operatorStack[:len(operatorStack)-1]
										} else {
											break
										}
									}
								} else {
									break
								}
							}
							operatorStack = append(operatorStack, currentToken)
							isValidToken = true
						}
			
						if(currentToken.Type == TOKEN_TYPE_OPEN_PARENTHESIS) {
							//if it's an open parenthesis '(' push it onto the stack
							operatorStack = append(operatorStack, currentToken)
							isValidToken = true
						}
			
						if(currentToken.Type == TOKEN_TYPE_CLOSE_PARENTHESIS) {
							isValidToken = true
							//close parenthesis
							if(len(operatorStack) > 0) {
								for true {
									if(operatorStack[len(operatorStack) - 1].Type != TOKEN_TYPE_OPEN_PARENTHESIS) {
										outputQueue = append(outputQueue, operatorStack[len(operatorStack) - 1])
										operatorStack = operatorStack[:len(operatorStack)-1]
									} else {
										operatorStack = operatorStack[:len(operatorStack)-1]
										break
									}
			
									if(len(operatorStack) == 0) {
										return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Operator expected", currentToken.FileName))		
									}
								}
							} else {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Operator expected", currentToken.FileName))
							}
						}
			
						if(!isValidToken) {
							return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Unexpected token '" + currentToken.Value + "'", currentToken.FileName))
						}
					}
			
					for len(operatorStack) > 0 {
						if(operatorStack[len(operatorStack) - 1].Type == TOKEN_TYPE_OPEN_PARENTHESIS) {
							return errors.New(SyntaxErrorMessage(operatorStack[len(operatorStack) - 1].Line, operatorStack[len(operatorStack) - 1].Column, "Unexpected token '" + operatorStack[len(operatorStack) - 1].Value + "'", operatorStack[len(operatorStack) - 1].FileName))
						}
						outputQueue = append(outputQueue, operatorStack[len(operatorStack) - 1])
						operatorStack = operatorStack[:len(operatorStack)-1]
					}
					//DumpToken(outputQueue)
					//the outputQueue contains the reverse polish notation
					if(len(outputQueue) > 0) {
						//read the reverse polish below
						var stack []Token
						var stack2 []Token
			
						for len(outputQueue) > 0 {
							currentToken := outputQueue[0]
							outputQueue = append(outputQueue[:0], outputQueue[1:]...) //pop the first element
			
							if(currentToken.Type == TOKEN_TYPE_PLUS || currentToken.Type == TOKEN_TYPE_MINUS || currentToken.Type == TOKEN_TYPE_DIVIDE || currentToken.Type == TOKEN_TYPE_MULTIPLY) {
								//arithmetic operation
								//NOTE: ASSUME THAT RIGHT OPERAND AND LEFT OPERAND ARE INTEGER AND FLOAT ONLY (NO IDENTIFIER, STRING ETC... (TEMPORARY ONLY)
								rightOperand := stack[len(stack)-1]
								stack = stack[:len(stack)-1]
								var tempRightInt int
								var tempRightFloat float64
								var tempRightString string
			
								leftOperand := stack[len(stack)-1]
								stack = stack[:len(stack)-1]
								var tempLeftInt int
								var tempLeftFloat float64
								var tempLeftString string

								var errConvert error
			
								result := leftOperand
								//convert the identifier to token below
								//left and right operand
								//will raise an error if not existing (of course)

								//TODO: CHECK FIRST IF LEFT/RIGHT OPERAND IS EXISTING AS FUNCTION, IF YES THEN RAISE AN ERROR
								//TODO: TRY TO FIX THE ASSIGNMENT OPERATION BEFORE THIS ONE

								if(leftOperand.Type == TOKEN_TYPE_IDENTIFIER) {
									leftOperand, errConvert = convertVariableToToken(leftOperand, *globalVariableArray, scopeName)
									if(errConvert != nil) {
										return errConvert
									}
								}
								if(rightOperand.Type == TOKEN_TYPE_IDENTIFIER) {
									rightOperand, errConvert = convertVariableToToken(rightOperand, *globalVariableArray, scopeName)
									if(errConvert != nil) {
										return errConvert
									}
								}
			
								//convert operands to its designated type
								if(leftOperand.Type == TOKEN_TYPE_INTEGER) {
									//convert to integer
									result.Type = TOKEN_TYPE_INTEGER
									tempLeftInt, _ = strconv.Atoi(leftOperand.Value)
									tempRightInt, _ = strconv.Atoi(rightOperand.Value)
								} else if(leftOperand.Type == TOKEN_TYPE_STRING) {
									//string
									result.Type = TOKEN_TYPE_STRING
									tempLeftString = leftOperand.Value
									tempRightString = rightOperand.Value
								} else {
									//let's assume that it should be converted to float (for now)
									result.Type = TOKEN_TYPE_FLOAT
									tempLeftFloat, _ = strconv.ParseFloat(leftOperand.Value, 32)
									tempRightFloat, _ = strconv.ParseFloat(rightOperand.Value, 32)
								}
			
								if(currentToken.Type == TOKEN_TYPE_PLUS) {
									//either addition or concatenation
			
									//validate left operand
									errLeft := expectedTokenTypes(leftOperand, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING)
									if (errLeft != nil) {
										return errLeft
									}
									//validate right operand
									errRight := expectedTokenTypes(rightOperand, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING)
									if (errRight != nil) {
										return errRight
									}
			
									if(leftOperand.Type == TOKEN_TYPE_INTEGER) {
										result.Value = strconv.Itoa(tempLeftInt + tempRightInt)
									} else if(leftOperand.Type == TOKEN_TYPE_STRING) {
										result.Value = tempLeftString + tempRightString //concatenate
									} else {
										//let's assume it's float
										result.Value = strconv.FormatFloat(tempLeftFloat + tempRightFloat, 'f', -1, 64)
									}
			
								} else {
									//substraction, division and multiplication
			
									//validate left operand (No String)
									errLeft := expectedTokenTypes(leftOperand, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT)
									if (errLeft != nil) {
										return errLeft
									}
									//validate right operand (No String)
									errRight := expectedTokenTypes(rightOperand, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT)
									if (errRight != nil) {
										return errRight
									}
			
									if(currentToken.Type == TOKEN_TYPE_MINUS) {
										//substraction
										if(leftOperand.Type == TOKEN_TYPE_INTEGER) {
											result.Value = strconv.Itoa(tempLeftInt - tempRightInt)
										} else {
											//let's assume it's float
											result.Value = strconv.FormatFloat(tempLeftFloat - tempRightFloat, 'f', -1, 64)
										}
									} else if(currentToken.Type == TOKEN_TYPE_MULTIPLY) {
										//multiplication
										if(leftOperand.Type == TOKEN_TYPE_INTEGER) {
											result.Value = strconv.Itoa(tempLeftInt * tempRightInt)
										} else {
											//let's assume it's float
											result.Value = strconv.FormatFloat(tempLeftFloat * tempRightFloat, 'f', -1, 64)
										}
									} else {
										//assume it's division
										if(leftOperand.Type == TOKEN_TYPE_INTEGER) {
											if(tempRightInt == 0) {
												return errors.New(SyntaxErrorMessage(rightOperand.Line, rightOperand.Column, "Division by zero", rightOperand.FileName))
											}
											result.Value = strconv.Itoa(tempLeftInt / tempRightInt)
										} else {
											//let's assume it's float
											if(tempRightInt == 0) {
												return errors.New(SyntaxErrorMessage(rightOperand.Line, rightOperand.Column, "Division by zero", rightOperand.FileName))
											}
											result.Value = strconv.FormatFloat(tempLeftFloat / tempRightFloat, 'f', -1, 64)
										}
									}
								}
			
								stack = append(stack, result)
			
							} else if(currentToken.Type == TOKEN_TYPE_EQUALS) {
								//assignment operation
								value := stack[len(stack)-1]
								stack = stack[:len(stack)-1]
								var errConvert error

								//if value is an identifier
								//the it's a variable
								if(value.Type == TOKEN_TYPE_IDENTIFIER) {
									value, errConvert = convertVariableToToken(value, *globalVariableArray, scopeName)
									if(errConvert != nil) {
										return errConvert
									}
								}
			
								variable := stack[len(stack)-1]
								stack = stack[:len(stack)-1]
			
								//validate value
								errVal := expectedTokenTypes(value, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING)
								if (errVal != nil) {
									return errVal
								}
								//validate variable
								errVar := expectedTokenTypes(variable, TOKEN_TYPE_IDENTIFIER)
								if (errVar != nil) {
									return errVar
								}

								//check if variable exists as a function
								//if yes then raise an error
								isExists, _ := isFunctionExists(variable, *globalFunctionArray)

								if(isExists) {
									return errors.New(SyntaxErrorMessage(variable.Line, variable.Column, "'" + variable.Value + "' exists as a function", variable.FileName))	
								}

								isExists, varIndex := isVariableExists(variable, *globalVariableArray, scopeName)
			
								if(!isExists) {
									//variable doesn't exists
									//create a new variable
									newVar := Variable{Name: variable.Value, ScopeName: scopeName}
									*globalVariableArray = append(*globalVariableArray, newVar)
									varIndex = len(*globalVariableArray) - 1 //last to execute
								}
			
								//modify the value/type of variable below
								if(value.Type == TOKEN_TYPE_INTEGER) {
									(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_INTEGER
									(*globalVariableArray)[varIndex].IntegerValue, _ = strconv.Atoi(value.Value)
								} else if(value.Type == TOKEN_TYPE_STRING) {
									(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_STRING
									(*globalVariableArray)[varIndex].StringValue = value.Value
								} else {
									//assume it's float for now (add types later on like string etc...)
									(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_FLOAT
									(*globalVariableArray)[varIndex].FloatValue, _ = strconv.ParseFloat(value.Value, 32)
								}

								stack = append(stack, value)
							} else if(currentToken.Type == TOKEN_TYPE_FUNCTION) {
								//function execution here
								var functionArguments []FunctionArgument

								//check if function is existing below
								isExists, funcIndex := isFunctionExists(currentToken, *globalFunctionArray)
								if(!isExists) {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Function '" + currentToken.Value + "' doesn't exists", currentToken.FileName))
								}

								//check if function got arguments
								if((*globalFunctionArray)[funcIndex].ArgumentCount > 0) {
									//function need parameters
									/* //TODO: DISABLED TEMPORARILY, NEED VALIDATOR FOR THE PARAMS
									if(len(stack) == 0 || len(stack) < (*globalFunctionArray)[funcIndex].ArgumentCount) {
										return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, currentToken.Value + " takes exactly " + strconv.Itoa((*globalFunctionArray)[funcIndex].ArgumentCount) + " argument", currentToken.FileName))
									}
									*/

									//add arguments from stack below
									processedArg := 0
									for true {
										var param Token
										//add to functionargument one by one
										if(len(stack2) > 0) {
											//if stack2 is available
											//get the parameters there
											param = stack2[0]
											stack2 = append(stack2[:0], stack2[1:]...)
										} else {
											param = stack[len(stack)-1]
											stack = stack[:len(stack)-1]
										}

										var errConvert error
										if(param.Type == TOKEN_TYPE_IDENTIFIER) {
											param, errConvert = convertVariableToToken(param, *globalVariableArray, scopeName)
											if(errConvert != nil) {
												return errConvert
											}
										}

										fa := FunctionArgument{}
										//convert token to param (TODO: create a function for this one?)
										if(param.Type == TOKEN_TYPE_INTEGER) {
											fa.Type = ARG_TYPE_INTEGER
											fa.IntegerValue, _ = strconv.Atoi(param.Value)
										} else if(param.Type == TOKEN_TYPE_STRING) {
											fa.Type = ARG_TYPE_STRING
											fa.StringValue = param.Value
										} else {
											//assume it's float for now (add types later on like string etc...)
											fa.Type = ARG_TYPE_FLOAT
											fa.FloatValue, _ = strconv.ParseFloat(param.Value, 32)
										}

										functionArguments = append(functionArguments, fa)

										processedArg += 1
										if (processedArg == (*globalFunctionArray)[funcIndex].ArgumentCount) {
											break
										}
									}
								}

								if((*globalFunctionArray)[funcIndex].IsNative) {
									//execute native function
									funcReturn := (*globalFunctionArray)[funcIndex].Run(functionArguments)
									//convert FunctionReturn to Token and append to stack (TODO: Create a function for conversion?)
									newToken := currentToken
									if(funcReturn.Type == RET_TYPE_INTEGER) {
										newToken.Type = TOKEN_TYPE_INTEGER
										newToken.Value = strconv.Itoa(funcReturn.IntegerValue)
									} else if(funcReturn.Type == RET_TYPE_STRING) {
										newToken.Type = TOKEN_TYPE_STRING
										newToken.Value = funcReturn.StringValue
									} else {
										//let's assume it's float
										newToken.Value = strconv.FormatFloat(funcReturn.FloatValue, 'f', -1, 64)
									}
									var newSlice []Token
									newSlice = append(newSlice, newToken)
									stack = append(newSlice, stack...)
								} else {
									//execute function from token
								}
							} else if(currentToken.Type == TOKEN_TYPE_COMMA) {
								//TODO: add last stack to special queue
								for true {
									if(len(stack) > 0) {
										//currentToken := stack[len(stack)-1]
										currentToken := stack[0]
										if(currentToken.Type == TOKEN_TYPE_FLOAT || currentToken.Type == TOKEN_TYPE_INTEGER || currentToken.Type == TOKEN_TYPE_IDENTIFIER || currentToken.Type == TOKEN_TYPE_STRING) {
											//stack = stack[:len(stack)-1]
											stack = append(stack[:0], stack[1:]...)
											stack2 = append(stack2, currentToken)
										} else {
											break
										}
									} else {
										break
									}
								}
							} else {
								stack = append(stack, currentToken)
							}
						}
			
						//fmt.Println(stack[0].Value)
					}
				}

			} else {
				//put the token to stack for shunting yard process later
				tokensToEvaluate = append(tokensToEvaluate, tokenArray[x])
			}
		}
	}

	return nil
}