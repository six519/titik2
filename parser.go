package main

import (
	"errors"
	"strconv"
	"unicode"
	"math/rand"
	"time"
	//"fmt"
)

func generateRandomNumbers() string {
	rand.Seed(time.Now().UnixNano())

	return strconv.Itoa(rand.Int())
}

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
	operatorPrecedences := map[string] int{"=": 0, "+": 1, "-": 1, "/": 2, "*": 2} //operator order of precedences
	var operatorStack []Token
	var functionStack []Token
	var outputQueue []Token
	var ignoreNewline bool = false
	var justAddTokens bool = false

	for x := 0; x < len(tokenArray); x++ {
		if(tokenArray[x].Type == TOKEN_TYPE_NEWLINE) {

			if(ignoreNewline) {
				//put the token to stack for shunting yard process later
				tokensToEvaluate = append(tokensToEvaluate, tokenArray[x])
				continue
			}

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
		
				justAddTokens = false
				//shunting-yard
				for len(tokensToEvaluate) > 0 {
					currentToken := tokensToEvaluate[0]
					tokensToEvaluate = append(tokensToEvaluate[:0], tokensToEvaluate[1:]...) //pop the first element
					isValidToken := false

					if(justAddTokens) {
						//function body, just add to outputqueue
						outputQueue = append(outputQueue, currentToken)

						if(currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_END) {
							justAddTokens = false
						}
						continue
					}

					if(currentToken.Type == TOKEN_TYPE_NEWLINE) {
						//just ignore newline
						continue
					}
		
					if(currentToken.Type == TOKEN_TYPE_INTEGER || currentToken.Type == TOKEN_TYPE_FLOAT || currentToken.Type == TOKEN_TYPE_IDENTIFIER || currentToken.Type == TOKEN_TYPE_STRING) {
						//If it's a number or identifier, add it to queue, (ADD TOKEN_TYPE_KEYWORD AND string and other acceptable tokens LATER)
						outputQueue = append(outputQueue, currentToken)
						isValidToken = true
					}

					dontIgnorePopping := true
					if(currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION || currentToken.Type == TOKEN_TYPE_FUNCTION || currentToken.Type == TOKEN_TYPE_COMMA || currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START || currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END) {
						isValidToken = true

						if(len(tokensToEvaluate) > 0) {
							if(currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION) {
								if(tokensToEvaluate[0].Type == TOKEN_TYPE_INVOKE_FUNCTION || tokensToEvaluate[0].Type == TOKEN_TYPE_COMMA) {
									dontIgnorePopping = false
								}
							}
						}

						if(currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION || currentToken.Type == TOKEN_TYPE_COMMA  || currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END) {
							//pop all operators from operator stack to output queue before the function
							//NOTE: don't include '=' (NOT SURE)
							if(dontIgnorePopping) {
								for true {
									if(len(operatorStack) > 0) {
										if(operatorStack[len(operatorStack) - 1].Type == TOKEN_TYPE_EQUALS) {
											break
										} else {
											outputQueue = append(outputQueue, operatorStack[len(operatorStack) - 1])
											operatorStack = operatorStack[:len(operatorStack)-1]
										}
									} else {
										break
									}
								}
							}
						}

						if(currentToken.Type == TOKEN_TYPE_FUNCTION || currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START) {
							functionStack = append(functionStack, currentToken)
						} else if(currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION || currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END) {
							outputQueue = append(outputQueue, functionStack[len(functionStack) - 1])
							functionStack = functionStack[:len(functionStack)-1]

							if(currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END) {
								//next is function body
								justAddTokens = true
							}
						} else {
							//comma
							//validate separator
							if(len(tokensToEvaluate) > 0) {
								if(tokensToEvaluate[0].Type != TOKEN_TYPE_FLOAT && tokensToEvaluate[0].Type != TOKEN_TYPE_INTEGER && tokensToEvaluate[0].Type != TOKEN_TYPE_STRING && tokensToEvaluate[0].Type != TOKEN_TYPE_IDENTIFIER && tokensToEvaluate[0].Type != TOKEN_TYPE_FUNCTION && tokensToEvaluate[0].Type != TOKEN_TYPE_OPEN_PARENTHESIS) {
									return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '" + tokensToEvaluate[0].Value + "'", tokensToEvaluate[0].FileName))
								}
							}
						}

						//dirty fix (not sure)
						if(!dontIgnorePopping) {
							//pop all operators from operator stack to output queue before the function
							//NOTE: don't include '=' (NOT SURE)
							for true {
								if(len(operatorStack) > 0) {
									if(operatorStack[len(operatorStack) - 1].Type == TOKEN_TYPE_EQUALS) {
										break
									} else {
										outputQueue = append(outputQueue, operatorStack[len(operatorStack) - 1])
										operatorStack = operatorStack[:len(operatorStack)-1]
									}
								} else {
									break
								}
							}
						}

					}

					if(currentToken.Type == TOKEN_TYPE_PLUS || currentToken.Type == TOKEN_TYPE_MINUS || currentToken.Type == TOKEN_TYPE_DIVIDE || currentToken.Type == TOKEN_TYPE_MULTIPLY || currentToken.Type == TOKEN_TYPE_EQUALS) {
						//the token is operator
						for true {
							if(len(operatorStack) > 0) {

								if(operatorPrecedences[operatorStack[len(operatorStack) - 1].Value] > operatorPrecedences[currentToken.Value]) {
									outputQueue = append(outputQueue, operatorStack[len(operatorStack) - 1])
									operatorStack = operatorStack[:len(operatorStack)-1]
								} else {
									break
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
				//end of shunting-yard

				//validate end of function
				if(len(functionStack) > 0) {
					return errors.New(SyntaxErrorMessage(functionStack[0].Line, functionStack[0].Column, "End of function call expected", functionStack[0].FileName))
				}

				//DumpToken(outputQueue)
				//the outputQueue contains the reverse polish notation
				
				if(len(outputQueue) > 0) {
					//read the reverse polish below
					var stack []Token
		
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
							errVal := expectedTokenTypes(value, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE)
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
								varIndex = len(*globalVariableArray) - 1 

								//check if the first letter of variable name is in uppercase
								//if yes then tag it as constant
								firstChar := string((*globalVariableArray)[varIndex].Name[0])
								if(unicode.IsUpper([]rune(firstChar)[0])) {
									(*globalVariableArray)[varIndex].IsConstant = true
								}
							} else {
								//if variable exists
								//check if constant, if yes then raise an error
								if((*globalVariableArray)[varIndex].IsConstant) {
									return errors.New(SyntaxErrorMessage(variable.Line, variable.Column, "Cannot override constant '" + variable.Value + "'", variable.FileName))
								}
							}
		
							//modify the value/type of variable below
							if(value.Type == TOKEN_TYPE_INTEGER) {
								(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_INTEGER
								(*globalVariableArray)[varIndex].IntegerValue, _ = strconv.Atoi(value.Value)
							} else if(value.Type == TOKEN_TYPE_STRING) {
								(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_STRING
								(*globalVariableArray)[varIndex].StringValue = value.Value
							} else if(value.Type == TOKEN_TYPE_FLOAT) {
								(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_FLOAT
								(*globalVariableArray)[varIndex].FloatValue, _ = strconv.ParseFloat(value.Value, 32)
							} else {
								//Nil
								(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_NONE
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
								//function parameter validation below
								if(len(stack) == 0 || len(stack) < (*globalFunctionArray)[funcIndex].ArgumentCount) {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, currentToken.Value + " takes exactly " + strconv.Itoa((*globalFunctionArray)[funcIndex].ArgumentCount) + " argument", currentToken.FileName))
								}

								//add arguments from stack below
								processedArg := 0
								for true {
									var param Token
									//add to functionargument one by one
									param = stack[len(stack)-1]
									stack = stack[:len(stack)-1]

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
									} else if(param.Type == TOKEN_TYPE_FLOAT) {
										fa.Type = ARG_TYPE_FLOAT
										fa.FloatValue, _ = strconv.ParseFloat(param.Value, 32)
									} else {
										//Nil
										fa.Type = ARG_TYPE_NONE
									}

									functionArguments = append(functionArguments, fa)

									processedArg += 1
									if (processedArg == (*globalFunctionArray)[funcIndex].ArgumentCount) {
										break
									}
								}
							}

							newToken := currentToken
							if((*globalFunctionArray)[funcIndex].IsNative) {
								//execute native function
								funcReturn := (*globalFunctionArray)[funcIndex].Run(functionArguments)
								//convert FunctionReturn to Token and append to stack (TODO: Create a function for conversion?)
								if(funcReturn.Type == RET_TYPE_INTEGER) {
									newToken.Type = TOKEN_TYPE_INTEGER
									newToken.Value = strconv.Itoa(funcReturn.IntegerValue)
								} else if(funcReturn.Type == RET_TYPE_STRING) {
									newToken.Type = TOKEN_TYPE_STRING
									newToken.Value = funcReturn.StringValue
								} else if(funcReturn.Type == RET_TYPE_FLOAT) {
									newToken.Type = TOKEN_TYPE_FLOAT
									newToken.Value = strconv.FormatFloat(funcReturn.FloatValue, 'f', -1, 64)
								} else {
									//Nil
									newToken.Type = TOKEN_TYPE_NONE
								}
								
							} else {
								//execute function from token
								//set the default return to Nil
								newToken.Type = TOKEN_TYPE_NONE
								thisScopeName := (*globalFunctionArray)[funcIndex].Name + generateRandomNumbers()

								//set the arguments below
								for ind := 0; ind < len(functionArguments); ind++ {

									newVar := Variable{Name: (*globalFunctionArray)[funcIndex].Arguments[ind].Value, ScopeName: thisScopeName}
									*globalVariableArray = append(*globalVariableArray, newVar)
									varIndex := len(*globalVariableArray) - 1

									if(functionArguments[ind].Type == ARG_TYPE_INTEGER) {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_INTEGER
										(*globalVariableArray)[varIndex].IntegerValue = functionArguments[ind].IntegerValue
									} else if(functionArguments[ind].Type == ARG_TYPE_STRING) {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_STRING
										(*globalVariableArray)[varIndex].StringValue = functionArguments[ind].StringValue
									} else if(functionArguments[ind].Type == ARG_TYPE_FLOAT) {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_FLOAT
										(*globalVariableArray)[varIndex].FloatValue = functionArguments[ind].FloatValue
									} else {
										//Nil
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_NONE
									}
								}
								
								//execute user defined function
								prsr := Parser{}
								parserErr := prsr.Parse((*globalFunctionArray)[funcIndex].Tokens, globalVariableArray, globalFunctionArray, thisScopeName)
						
								if(parserErr != nil) {
									return parserErr
								}

								//TODO: NEED CLEANUP OF VARIABLES BELOW
								//DELETE GENERATED VARIABLES with thisScopeName
							}

							stack = append(stack, newToken)
						
						} else if(currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START) {
							//check if function already exists
							//if yes then raise an error
							isExists, _ := isFunctionExists(currentToken, *globalFunctionArray)
							if(isExists) {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Function '" + currentToken.Value + "' already exists", currentToken.FileName))
							}

							//get function arguments
							var functionParams []Token
							for true {
								var param Token

								if(len(stack) > 0) {
									param = stack[len(stack)-1]
									stack = stack[:len(stack)-1]

									//validate param type
									errParam := expectedTokenTypes(param, TOKEN_TYPE_IDENTIFIER)
									if (errParam != nil) {
										return errParam
									}

									//check if param is a constant
									firstChar := string(param.Value[0])
									if(unicode.IsUpper([]rune(firstChar)[0])) {
										return errors.New(SyntaxErrorMessage(param.Line, param.Column, "Argument cannot be a constant", param.FileName))
									}

									//check if param already exists
									if(isParamExists(param, functionParams)) {
										return errors.New(SyntaxErrorMessage(param.Line, param.Column, "Duplicate argument '" + param.Value + "' in function definition", param.FileName))
									}

									functionParams = append(functionParams, param)
								} else {
									break
								}	

							}

							//define function below
							newFunction := Function{Name: currentToken.Value, IsNative: false, ArgumentCount: len(functionParams), Arguments: functionParams}
							
							//append all tokens to function (body of function)
							for true {
								currentToken := outputQueue[0]
								outputQueue = append(outputQueue[:0], outputQueue[1:]...)

								newFunction.Tokens = append(newFunction.Tokens, currentToken)

								if(currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_END) {
									break
								}

								if(len(outputQueue) == 0) {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Invalid statement", currentToken.FileName))
								}
							}
							
							//append to global functions
							*globalFunctionArray = append(*globalFunctionArray, newFunction)
							stack = append(stack, currentToken) //TODO: not sure if it should append the TOKEN_TYPE_FUNCTION_DEF_END
						} else {
							stack = append(stack, currentToken)
						}
					}

					if(len(stack) > 1) {
						return errors.New(SyntaxErrorMessage(stack[0].Line, stack[0].Column, "Invalid statement", stack[0].FileName))
					} else {
						if(stack[0].Type == TOKEN_TYPE_IDENTIFIER) {
							return errors.New(SyntaxErrorMessage(stack[0].Line, stack[0].Column, "Unexpected token '" + stack[0].Value + "'", stack[0].FileName))
						}
					}
				}
				
			}

		} else if(tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_START || tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_END) {
			if(tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_START) {
				ignoreNewline = true
			} else {
				//TOKEN_TYPE_FUNCTION_DEF_END
				ignoreNewline = false
			}
			//put the token to stack for shunting yard process later
			tokensToEvaluate = append(tokensToEvaluate, tokenArray[x])
		} else {
			//put the token to stack for shunting yard process later
			tokensToEvaluate = append(tokensToEvaluate, tokenArray[x])
		}
	}

	return nil
}