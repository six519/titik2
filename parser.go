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

func (parser Parser) Parse(tokenArray []Token, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, gotReturn *bool, returnToken *Token, isLoop bool, needBreak *bool, stackReference *[]Token) error {
	var tokensToEvaluate []Token
	operatorPrecedences := map[string] int{"function_return": 0, "=": 1, "+": 2, "-": 2, "&": 2, "|": 2, "==": 2, "<>": 2, ">": 2, "<": 2, ">=": 2, "<=": 2, "/": 3, "*": 3} //operator order of precedences
	currentContext := "main_context"
	var operatorStack map[string][]Token
	operatorStack = make(map[string][]Token)
	var arrayArgCount map[string]int
	arrayArgCount = make(map[string]int)
	var functionStack []Token
	var outputQueue []Token
	var ignoreNewline bool = false
	var justAddTokens bool = false
	var isFunctionDefinition bool = false
	var isLoopStatement bool = false
	var isIfStatement bool = false
	var openLoopCount int = 0
	var openIfCount int = 0

	for x := 0; x < len(tokenArray); x++ {
		if(tokenArray[x].Type == TOKEN_TYPE_NEWLINE) {

			if(ignoreNewline) {
				//put the token to stack for shunting yard process later
				tokensToEvaluate = append(tokensToEvaluate, tokenArray[x])

				if(len(tokenArray) == (x+1)) {
					return errors.New(SyntaxErrorMessage(tokenArray[x].Line, tokenArray[x].Column, "Unfinished statement", tokenArray[x].FileName))
				}

				continue
			}

			//execute shunting yard
			if(len(tokensToEvaluate) > 0) {
				if(tokensToEvaluate[0].Type == TOKEN_TYPE_PLUS || tokensToEvaluate[0].Type == TOKEN_TYPE_MINUS || tokensToEvaluate[0].Type == TOKEN_TYPE_DIVIDE || tokensToEvaluate[0].Type == TOKEN_TYPE_MULTIPLY || tokensToEvaluate[0].Type == TOKEN_TYPE_EQUALS || tokensToEvaluate[0].Type == TOKEN_TYPE_AMPERSAND || tokensToEvaluate[0].Type == TOKEN_TYPE_OR || tokensToEvaluate[0].Type == TOKEN_TYPE_EQUALITY || tokensToEvaluate[0].Type == TOKEN_TYPE_INEQUALITY || tokensToEvaluate[0].Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS || tokensToEvaluate[0].Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS || tokensToEvaluate[0].Type == TOKEN_TYPE_GREATER_THAN || tokensToEvaluate[0].Type == TOKEN_TYPE_LESS_THAN) {
					//syntax error if the first token is an operator
					return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '" + tokensToEvaluate[0].Value + "'", tokensToEvaluate[0].FileName))
				}
		
				if(tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_PLUS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_MINUS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_DIVIDE || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_MULTIPLY || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_EQUALS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_AMPERSAND || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_OR || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_EQUALITY || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_INEQUALITY || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_GREATER_THAN || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_LESS_THAN) {
					//syntax error if the last token is an operator
					return errors.New(SyntaxErrorMessage(tokensToEvaluate[len(tokensToEvaluate)-1].Line, tokensToEvaluate[len(tokensToEvaluate)-1].Column, "Unfinished statement", tokensToEvaluate[len(tokensToEvaluate)-1].FileName))
				}
		
				justAddTokens = false
				isFunctionDefinition = false
				isLoopStatement = false
				isIfStatement = false
				openIfCount = 0
				openLoopCount = 0
				//shunting-yard
				for len(tokensToEvaluate) > 0 {
					currentToken := tokensToEvaluate[0]
					tokensToEvaluate = append(tokensToEvaluate[:0], tokensToEvaluate[1:]...) //pop the first element
					isValidToken := false
					currentContext = currentToken.Context

					if(justAddTokens) {
						//function body, just add to outputqueue
						outputQueue = append(outputQueue, currentToken)

						if(currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_END) {
							justAddTokens = false
							isFunctionDefinition = false
						}
						if(!isFunctionDefinition) {
							if(!isIfStatement) {
								if(currentToken.Type == TOKEN_TYPE_FOR_LOOP_START) {
									openLoopCount += 1
								}
								if(currentToken.Type == TOKEN_TYPE_FOR_LOOP_END) {
									if(openLoopCount == 0) {
										justAddTokens = false
									}
									if(openLoopCount > 0) {
										openLoopCount = openLoopCount - 1
									}
								}
							}
							if(!isLoopStatement) {
								if(currentToken.Type == TOKEN_TYPE_IF_START) {
									openIfCount += 1
								}
								if(currentToken.Type == TOKEN_TYPE_IF_END) {
									if(openIfCount == 0) {
										justAddTokens = false
									}
									if(openIfCount > 0) {
										openIfCount = openIfCount - 1
									}
								}
							}
						}
						continue
					}

					if(currentToken.Type == TOKEN_TYPE_NEWLINE) {
						//just ignore newline
						continue
					}
		
					if(currentToken.Type == TOKEN_TYPE_INTEGER || currentToken.Type == TOKEN_TYPE_FLOAT || currentToken.Type == TOKEN_TYPE_IDENTIFIER || currentToken.Type == TOKEN_TYPE_STRING || currentToken.Type == TOKEN_TYPE_LOOP_BREAK) {
						//If it's a number or identifier, add it to queue, (ADD TOKEN_TYPE_KEYWORD AND string and other acceptable tokens LATER)
						outputQueue = append(outputQueue, currentToken)
						isValidToken = true
					}

					//dontIgnorePopping := true
					if(currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION || currentToken.Type == TOKEN_TYPE_FUNCTION || currentToken.Type == TOKEN_TYPE_COMMA || currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START || currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END || currentToken.Type == TOKEN_TYPE_FOR_LOOP_START || currentToken.Type == TOKEN_TYPE_FOR_LOOP_PARAM_END || currentToken.Type == TOKEN_TYPE_IF_START || currentToken.Type == TOKEN_TYPE_IF_PARAM_END || currentToken.Type == TOKEN_TYPE_CLOSE_BRACES || currentToken.Type == TOKEN_TYPE_OPEN_BRACES) {
						isValidToken = true
						/*
						if(len(tokensToEvaluate) > 0) {
							if(currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION) {
								if(tokensToEvaluate[0].Type == TOKEN_TYPE_INVOKE_FUNCTION || tokensToEvaluate[0].Type == TOKEN_TYPE_COMMA || tokensToEvaluate[0].Type == TOKEN_TYPE_FOR_LOOP_PARAM_END) {
									dontIgnorePopping = false
								}
							}
						}
						*/

						if(currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION || currentToken.Type == TOKEN_TYPE_COMMA  || currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END || currentToken.Type == TOKEN_TYPE_FOR_LOOP_PARAM_END || currentToken.Type == TOKEN_TYPE_IF_PARAM_END || currentToken.Type == TOKEN_TYPE_CLOSE_BRACES) {
							//pop all operators from operator stack to output queue before the function
							//NOTE: don't include '=' (NOT SURE)
							for true {
								if(len(operatorStack[currentContext]) > 0) {
									if(operatorStack[currentContext][len(operatorStack[currentContext]) - 1].Type == TOKEN_TYPE_EQUALS) {
										break
									} else {
										outputQueue = append(outputQueue, operatorStack[currentContext][len(operatorStack[currentContext]) - 1])
										operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
									}
								} else {
									break
								}
							}
						}

						if(currentToken.Type == TOKEN_TYPE_FUNCTION || currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START || currentToken.Type == TOKEN_TYPE_FOR_LOOP_START || currentToken.Type == TOKEN_TYPE_IF_START || currentToken.Type == TOKEN_TYPE_OPEN_BRACES) {
							functionStack = append(functionStack, currentToken)
							if(currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START) {
								isFunctionDefinition = true
							}
							if(!isFunctionDefinition) {
								if(!isIfStatement) {
									if(currentToken.Type == TOKEN_TYPE_FOR_LOOP_START) {
										isLoopStatement = true
									}
								}
								if(!isLoopStatement) {
									if(currentToken.Type == TOKEN_TYPE_IF_START) {
										isIfStatement = true
									}
								}
							}
							if(currentToken.Type == TOKEN_TYPE_OPEN_BRACES) {
								if(len(tokensToEvaluate) > 0) {
									if(tokensToEvaluate[0].Type != TOKEN_TYPE_CLOSE_BRACES) {
										arrayArgCount[currentContext] = 1
										if(tokensToEvaluate[0].Type == TOKEN_TYPE_OPEN_BRACES) {
											return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '" + tokensToEvaluate[0].Value + "'", tokensToEvaluate[0].FileName))
										}
									}
								} else {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Unfinished statement", currentToken.FileName))
								}
							}
						} else if(currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION || currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END || currentToken.Type == TOKEN_TYPE_FOR_LOOP_PARAM_END || currentToken.Type == TOKEN_TYPE_IF_PARAM_END || currentToken.Type == TOKEN_TYPE_CLOSE_BRACES) {
							tokenToAppend := functionStack[len(functionStack) - 1]
							
							if(currentToken.Type == TOKEN_TYPE_CLOSE_BRACES) {
								tokenToAppend.OtherInt = arrayArgCount[currentContext]
							}

							outputQueue = append(outputQueue, tokenToAppend)
							functionStack = functionStack[:len(functionStack)-1]

							if(currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END || currentToken.Type == TOKEN_TYPE_FOR_LOOP_PARAM_END || currentToken.Type == TOKEN_TYPE_IF_PARAM_END) {
								//next is function body
								justAddTokens = true
							}
						} else {
							//comma
							//count parameter (currently for array only, not sure in the future (for function params?))
							arrayArgCount[currentContext] += 1
							//validate separator
							if(len(tokensToEvaluate) > 0) {
								if(tokensToEvaluate[0].Type != TOKEN_TYPE_FLOAT && tokensToEvaluate[0].Type != TOKEN_TYPE_INTEGER && tokensToEvaluate[0].Type != TOKEN_TYPE_STRING && tokensToEvaluate[0].Type != TOKEN_TYPE_IDENTIFIER && tokensToEvaluate[0].Type != TOKEN_TYPE_FUNCTION && tokensToEvaluate[0].Type != TOKEN_TYPE_OPEN_PARENTHESIS) {
									return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '" + tokensToEvaluate[0].Value + "'", tokensToEvaluate[0].FileName))
								}
							}
						}

						/*
						//dirty fix (not sure)
						if(!dontIgnorePopping) {
							//pop all operators from operator stack to output queue before the function
							//NOTE: don't include '=' (NOT SURE)
							for true {
								if(len(operatorStack[currentContext]) > 0) {
									if(operatorStack[currentContext][len(operatorStack[currentContext]) - 1].Type == TOKEN_TYPE_EQUALS) {
										break
									} else {
										outputQueue = append(outputQueue, operatorStack[currentContext][len(operatorStack[currentContext]) - 1])
										operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
									}
								} else {
									break
								}
							}
						}
						*/

					}

					if(currentToken.Type == TOKEN_TYPE_PLUS || currentToken.Type == TOKEN_TYPE_MINUS || currentToken.Type == TOKEN_TYPE_DIVIDE || currentToken.Type == TOKEN_TYPE_MULTIPLY || currentToken.Type == TOKEN_TYPE_EQUALS || currentToken.Type == TOKEN_TYPE_AMPERSAND || currentToken.Type == TOKEN_TYPE_OR || currentToken.Type == TOKEN_TYPE_FUNCTION_RETURN || currentToken.Type == TOKEN_TYPE_EQUALITY || currentToken.Type == TOKEN_TYPE_INEQUALITY || currentToken.Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS || currentToken.Type == TOKEN_TYPE_LESS_THAN || currentToken.Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS || currentToken.Type == TOKEN_TYPE_GREATER_THAN) {
						//the token is operator
						for true {
							if(len(operatorStack[currentContext]) > 0) {

								if(currentToken.Type == TOKEN_TYPE_FUNCTION_RETURN) {
									if(operatorPrecedences[operatorStack[currentContext][len(operatorStack[currentContext]) - 1].Value] >= operatorPrecedences["function_return"]) {
										outputQueue = append(outputQueue, operatorStack[currentContext][len(operatorStack[currentContext]) - 1])
										operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
									} else {
										break
									}
								} else {
									if(operatorPrecedences[operatorStack[currentContext][len(operatorStack[currentContext]) - 1].Value] >= operatorPrecedences[currentToken.Value]) {
										outputQueue = append(outputQueue, operatorStack[currentContext][len(operatorStack[currentContext]) - 1])
										operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
									} else {
										break
									}
								}

							} else {
								break
							}
						}
						operatorStack[currentContext] = append(operatorStack[currentContext], currentToken)
						isValidToken = true
					}
		
					if(currentToken.Type == TOKEN_TYPE_OPEN_PARENTHESIS) {
						//if it's an open parenthesis '(' push it onto the stack
						operatorStack[currentContext] = append(operatorStack[currentContext], currentToken)
						isValidToken = true
					}
		
					if(currentToken.Type == TOKEN_TYPE_CLOSE_PARENTHESIS) {
						isValidToken = true
						//close parenthesis
						if(len(operatorStack[currentContext]) > 0) {
							for true {
								if(operatorStack[currentContext][len(operatorStack[currentContext]) - 1].Type != TOKEN_TYPE_OPEN_PARENTHESIS) {
									outputQueue = append(outputQueue, operatorStack[currentContext][len(operatorStack[currentContext]) - 1])
									operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
								} else {
									operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
									break
								}
		
								if(len(operatorStack[currentContext]) == 0) {
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
		
				for len(operatorStack["main_context"]) > 0 {
					if(operatorStack["main_context"][len(operatorStack["main_context"]) - 1].Type == TOKEN_TYPE_OPEN_PARENTHESIS) {
						return errors.New(SyntaxErrorMessage(operatorStack["main_context"][len(operatorStack["main_context"]) - 1].Line, operatorStack["main_context"][len(operatorStack["main_context"]) - 1].Column, "Unexpected token '" + operatorStack["main_context"][len(operatorStack["main_context"]) - 1].Value + "'", operatorStack["main_context"][len(operatorStack["main_context"]) - 1].FileName))
					}
					outputQueue = append(outputQueue, operatorStack["main_context"][len(operatorStack["main_context"]) - 1])
					operatorStack["main_context"] = operatorStack["main_context"][:len(operatorStack["main_context"])-1]
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
		
						} else if(currentToken.Type == TOKEN_TYPE_AMPERSAND || currentToken.Type == TOKEN_TYPE_OR || currentToken.Type == TOKEN_TYPE_EQUALITY || currentToken.Type == TOKEN_TYPE_INEQUALITY || currentToken.Type == TOKEN_TYPE_GREATER_THAN || currentToken.Type == TOKEN_TYPE_LESS_THAN || currentToken.Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS || currentToken.Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS) {
							//logical or comparison operation
							rightOperand := stack[len(stack)-1]
							stack = stack[:len(stack)-1]
		
							leftOperand := stack[len(stack)-1]
							stack = stack[:len(stack)-1]

							var errConvert error
		
							result := leftOperand

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

							if(currentToken.Type == TOKEN_TYPE_AMPERSAND || currentToken.Type == TOKEN_TYPE_OR) {
								//validation for logical operations
								//validate left operand
								errLeft := expectedTokenTypes(leftOperand, TOKEN_TYPE_BOOLEAN)
								if (errLeft != nil) {
									return errLeft
								}
								//validate right operand
								errRight := expectedTokenTypes(rightOperand, TOKEN_TYPE_BOOLEAN)
								if (errRight != nil) {
									return errRight
								}
							} else {
								//validation for comparison operations
								//validate left operand
								errLeft := expectedTokenTypes(leftOperand, TOKEN_TYPE_BOOLEAN, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE)
								if (errLeft != nil) {
									return errLeft
								}
								//validate right operand
								errRight := expectedTokenTypes(rightOperand, TOKEN_TYPE_BOOLEAN, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE)
								if (errRight != nil) {
									return errRight
								}
							}

							result.Type = TOKEN_TYPE_BOOLEAN

							if(currentToken.Type == TOKEN_TYPE_AMPERSAND) {
								//LOGICAL AND operation
								leftBool := convertTokenToBool(leftOperand)
								rightBool := convertTokenToBool(rightOperand)
								if(leftBool && rightBool) {
									result.Value = "true"
								} else {
									result.Value = "false"
								}
							} else if(currentToken.Type == TOKEN_TYPE_OR) {
								//LOGICAL OR operation
								leftBool := convertTokenToBool(leftOperand)
								rightBool := convertTokenToBool(rightOperand)
								if(leftBool || rightBool) {
									result.Value = "true"
								} else {
									result.Value = "false"
								}
							} else if(currentToken.Type == TOKEN_TYPE_EQUALITY) {
								//COMPARISON EQUALITY operation
								switch leftOperand.Type {
									case TOKEN_TYPE_INTEGER:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
											var rightOperandInt int = 0

											if(rightOperand.Type == TOKEN_TYPE_INTEGER) {
												rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
											} else {
												rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
												rightOperandInt = int(rightOperandFloat)
											}
											
											if(leftOperandInt == rightOperandInt) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_FLOAT:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 32)
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
											
											if(leftOperandFloat == rightOperandFloat) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_STRING:
										if(rightOperand.Type == TOKEN_TYPE_STRING) {
											if(leftOperand.Value == rightOperand.Value) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_BOOLEAN:
										if(rightOperand.Type == TOKEN_TYPE_BOOLEAN) {
											leftBool := convertTokenToBool(leftOperand)
											rightBool := convertTokenToBool(rightOperand)

											if(leftBool == rightBool) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									default:
										//TOKEN_TYPE_NONE
										if(rightOperand.Type == TOKEN_TYPE_NONE) {
											result.Value = "true"
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											result.Value = "false"
										}
								}
							} else if(currentToken.Type == TOKEN_TYPE_INEQUALITY) {
								//COMPARISON INEQUALITY operation
								switch leftOperand.Type {
									case TOKEN_TYPE_INTEGER:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
											var rightOperandInt int = 0

											if(rightOperand.Type == TOKEN_TYPE_INTEGER) {
												rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
											} else {
												rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
												rightOperandInt = int(rightOperandFloat)
											}
											
											if(leftOperandInt != rightOperandInt) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "true"
										}
									case TOKEN_TYPE_FLOAT:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 32)
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
											
											if(leftOperandFloat != rightOperandFloat) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "true"
										}
									case TOKEN_TYPE_STRING:
										if(rightOperand.Type == TOKEN_TYPE_STRING) {
											if(leftOperand.Value != rightOperand.Value) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "true"
										}
									case TOKEN_TYPE_BOOLEAN:
										if(rightOperand.Type == TOKEN_TYPE_BOOLEAN) {
											leftBool := convertTokenToBool(leftOperand)
											rightBool := convertTokenToBool(rightOperand)

											if(leftBool != rightBool) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_NONE
											result.Value = "true"
										}
									default:
										//TOKEN_TYPE_NONE
										if(rightOperand.Type == TOKEN_TYPE_NONE) {
											result.Value = "false"
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											result.Value = "true"
										}
								}
							} else if(currentToken.Type == TOKEN_TYPE_GREATER_THAN) {
								//COMPARISON GREATER THAN operation
								switch leftOperand.Type {
									case TOKEN_TYPE_INTEGER:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
											var rightOperandInt int = 0

											if(rightOperand.Type == TOKEN_TYPE_INTEGER) {
												rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
											} else {
												rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
												rightOperandInt = int(rightOperandFloat)
											}
											
											if(leftOperandInt > rightOperandInt) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_FLOAT:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 32)
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
											
											if(leftOperandFloat > rightOperandFloat) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_STRING:
										result.Value = "false"
									case TOKEN_TYPE_BOOLEAN:
										result.Value = "false"
									default:
										//TOKEN_TYPE_NONE
										result.Value = "false"
								}
							} else if(currentToken.Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS) {
								//COMPARISON GREATER THAN OR EQUALS operation
								switch leftOperand.Type {
									case TOKEN_TYPE_INTEGER:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
											var rightOperandInt int = 0

											if(rightOperand.Type == TOKEN_TYPE_INTEGER) {
												rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
											} else {
												rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
												rightOperandInt = int(rightOperandFloat)
											}
											
											if(leftOperandInt >= rightOperandInt) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_FLOAT:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 32)
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
											
											if(leftOperandFloat >= rightOperandFloat) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_STRING:
										if(rightOperand.Type == TOKEN_TYPE_STRING) {
											if(leftOperand.Value == rightOperand.Value) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_BOOLEAN:
										if(rightOperand.Type == TOKEN_TYPE_BOOLEAN) {
											leftBool := convertTokenToBool(leftOperand)
											rightBool := convertTokenToBool(rightOperand)

											if(leftBool == rightBool) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									default:
										//TOKEN_TYPE_NONE
										if(rightOperand.Type == TOKEN_TYPE_NONE) {
											result.Value = "true"
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											result.Value = "false"
										}
								}
							} else if(currentToken.Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS) {
								//COMPARISON LESS THAN OR EQUALS operation
								switch leftOperand.Type {
									case TOKEN_TYPE_INTEGER:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
											var rightOperandInt int = 0

											if(rightOperand.Type == TOKEN_TYPE_INTEGER) {
												rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
											} else {
												rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
												rightOperandInt = int(rightOperandFloat)
											}
											
											if(leftOperandInt <= rightOperandInt) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_FLOAT:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 32)
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
											
											if(leftOperandFloat <= rightOperandFloat) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_STRING:
										if(rightOperand.Type == TOKEN_TYPE_STRING) {
											if(leftOperand.Value == rightOperand.Value) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_BOOLEAN:
										if(rightOperand.Type == TOKEN_TYPE_BOOLEAN) {
											leftBool := convertTokenToBool(leftOperand)
											rightBool := convertTokenToBool(rightOperand)

											if(leftBool == rightBool) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									default:
										//TOKEN_TYPE_NONE
										if(rightOperand.Type == TOKEN_TYPE_NONE) {
											result.Value = "true"
										} else {
											//TOKEN_TYPE_INTEGER
											//TOKEN_TYPE_FLOAT
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											result.Value = "false"
										}
								}
							} else {
								//COMPARISON LESS THAN operation
								switch leftOperand.Type {
									case TOKEN_TYPE_INTEGER:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
											var rightOperandInt int = 0

											if(rightOperand.Type == TOKEN_TYPE_INTEGER) {
												rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
											} else {
												rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
												rightOperandInt = int(rightOperandFloat)
											}
											
											if(leftOperandInt < rightOperandInt) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_FLOAT:
										if(rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT) {
											leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 32)
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 32)
											
											if(leftOperandFloat < rightOperandFloat) {
												result.Value = "true"
											} else {
												result.Value = "false"
											}
										} else {
											//TOKEN_TYPE_STRING
											//TOKEN_TYPE_BOOLEAN
											//TOKEN_TYPE_NONE
											result.Value = "false"
										}
									case TOKEN_TYPE_STRING:
										result.Value = "false"
									case TOKEN_TYPE_BOOLEAN:
										result.Value = "false"
									default:
										//TOKEN_TYPE_NONE
										result.Value = "false"
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
							errVal := expectedTokenTypes(value, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE, TOKEN_TYPE_BOOLEAN, TOKEN_TYPE_ARRAY)
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

							//if not main scope then check if a system constant
							//if yes then raise an error
							if(scopeName != "main") {
								if(isSystemVariable(variable.Value, *globalNativeVarList)) {
									return errors.New(SyntaxErrorMessage(variable.Line, variable.Column, "Cannot assign to " + variable.Value, variable.FileName))
								}
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
							} else if(value.Type == TOKEN_TYPE_BOOLEAN) {
								(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_BOOLEAN
								if(value.Value == "true") {
									(*globalVariableArray)[varIndex].BooleanValue = true
								} else {
									(*globalVariableArray)[varIndex].BooleanValue = false
								}
							} else if(value.Type == TOKEN_TYPE_ARRAY) {
								(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_ARRAY
								for arrayIndex := 0; arrayIndex < len(value.Array); arrayIndex++ {
									thisVar := Variable{}

									if(value.Array[arrayIndex].Type == TOKEN_TYPE_INTEGER) {
										thisVar.Type = VARIABLE_TYPE_INTEGER
										thisVar.IntegerValue, _ = strconv.Atoi(value.Array[arrayIndex].Value)
									} else if(value.Array[arrayIndex].Type == TOKEN_TYPE_STRING) {
										thisVar.Type = VARIABLE_TYPE_STRING
										thisVar.StringValue = value.Array[arrayIndex].Value
									} else if(value.Array[arrayIndex].Type == TOKEN_TYPE_FLOAT) {
										thisVar.Type = VARIABLE_TYPE_FLOAT
										thisVar.FloatValue, _ = strconv.ParseFloat(value.Array[arrayIndex].Value, 32)
									} else if(value.Array[arrayIndex].Type == TOKEN_TYPE_BOOLEAN) {
										thisVar.Type = VARIABLE_TYPE_BOOLEAN
										if(value.Array[arrayIndex].Value == "true") {
											thisVar.BooleanValue = true
										} else {
											thisVar.BooleanValue = false
										}
									} else {
										thisVar.Type = VARIABLE_TYPE_NONE
									}

									(*globalVariableArray)[varIndex].ArrayValue = append((*globalVariableArray)[varIndex].ArrayValue, thisVar)
								}
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
									} else if(param.Type == TOKEN_TYPE_BOOLEAN) {
										fa.Type = ARG_TYPE_BOOLEAN
										if(param.Value == "true") {
											fa.BooleanValue = true
										} else {
											fa.BooleanValue = false
										}
									} else if(param.Type == TOKEN_TYPE_ARRAY) {
										fa.Type = ARG_TYPE_ARRAY
										for thisArrayIndex := 0; thisArrayIndex < len(param.Array); thisArrayIndex++ {
											thisArgument := FunctionArgument{}
											if(param.Array[thisArrayIndex].Type == TOKEN_TYPE_INTEGER) {
												thisArgument.Type = ARG_TYPE_INTEGER
												thisArgument.IntegerValue, _ = strconv.Atoi(param.Array[thisArrayIndex].Value)
											} else if(param.Array[thisArrayIndex].Type == TOKEN_TYPE_STRING) {
												thisArgument.Type = ARG_TYPE_STRING
												thisArgument.StringValue = param.Array[thisArrayIndex].Value
											} else if(param.Array[thisArrayIndex].Type == TOKEN_TYPE_FLOAT) {
												thisArgument.Type = ARG_TYPE_FLOAT
												thisArgument.FloatValue, _ = strconv.ParseFloat(param.Array[thisArrayIndex].Value, 32)
											} else if(param.Array[thisArrayIndex].Type == TOKEN_TYPE_BOOLEAN) {
												thisArgument.Type = ARG_TYPE_BOOLEAN
												if(param.Array[thisArrayIndex].Value == "true") {
													thisArgument.BooleanValue = true
												} else {
													thisArgument.BooleanValue = false
												}
											} else {
												thisArgument.Type = ARG_TYPE_NONE
											}
											fa.ArrayValue = append(fa.ArrayValue, thisArgument)
										}
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
								var thisError error
								funcReturn := (*globalFunctionArray)[funcIndex].Run(functionArguments, &thisError)
								if(thisError != nil) {
									return thisError
								}
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
								} else if(funcReturn.Type == RET_TYPE_BOOLEAN) {
									newToken.Type = TOKEN_TYPE_BOOLEAN
									if(funcReturn.BooleanValue) {
										//true
										newToken.Value = "true"
									} else {
										//false
										newToken.Value = "false"
									}
								} else if(funcReturn.Type == RET_TYPE_ARRAY) {
									newToken.Type = TOKEN_TYPE_ARRAY
									newToken.OtherInt = len(funcReturn.ArrayValue)

									for thisArrayIndex := 0; thisArrayIndex < len(funcReturn.ArrayValue); thisArrayIndex++ {
										thisNewToken := Token{}

										if(funcReturn.ArrayValue[thisArrayIndex].Type == RET_TYPE_INTEGER) {
											thisNewToken.Type = TOKEN_TYPE_INTEGER
											thisNewToken.Value = strconv.Itoa(funcReturn.ArrayValue[thisArrayIndex].IntegerValue)
										} else if(funcReturn.ArrayValue[thisArrayIndex].Type == RET_TYPE_STRING) {
											thisNewToken.Type = TOKEN_TYPE_STRING
											thisNewToken.Value = funcReturn.ArrayValue[thisArrayIndex].StringValue
										} else if(funcReturn.ArrayValue[thisArrayIndex].Type == RET_TYPE_FLOAT) {
											thisNewToken.Type = TOKEN_TYPE_FLOAT
											thisNewToken.Value = strconv.FormatFloat(funcReturn.ArrayValue[thisArrayIndex].FloatValue, 'f', -1, 64)
										} else if(funcReturn.ArrayValue[thisArrayIndex].Type == RET_TYPE_BOOLEAN) {
											thisNewToken.Type = TOKEN_TYPE_BOOLEAN
											if(funcReturn.ArrayValue[thisArrayIndex].BooleanValue) {
												thisNewToken.Value = "true"
											} else {
												thisNewToken.Value = "false"
											}
										} else {
											thisNewToken.Type = TOKEN_TYPE_NONE
										}

										newToken.Array = append(newToken.Array, thisNewToken)
									}
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
									} else if(functionArguments[ind].Type == ARG_TYPE_BOOLEAN) {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_BOOLEAN
										(*globalVariableArray)[varIndex].BooleanValue = functionArguments[ind].BooleanValue
									} else if(functionArguments[ind].Type == ARG_TYPE_ARRAY) {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_ARRAY
										for thisArrayIndex := 0; thisArrayIndex < len(functionArguments[ind].ArrayValue); thisArrayIndex++ {
											thisVar := Variable{}

											if(functionArguments[ind].ArrayValue[thisArrayIndex].Type == ARG_TYPE_INTEGER) {
												thisVar.Type = VARIABLE_TYPE_INTEGER
												thisVar.IntegerValue = functionArguments[ind].ArrayValue[thisArrayIndex].IntegerValue
											} else if(functionArguments[ind].ArrayValue[thisArrayIndex].Type == ARG_TYPE_STRING) {
												thisVar.Type = VARIABLE_TYPE_STRING
												thisVar.StringValue = functionArguments[ind].ArrayValue[thisArrayIndex].StringValue
											} else if(functionArguments[ind].ArrayValue[thisArrayIndex].Type == ARG_TYPE_FLOAT) {
												thisVar.Type = VARIABLE_TYPE_FLOAT
												thisVar.FloatValue = functionArguments[ind].ArrayValue[thisArrayIndex].FloatValue
											} else if(functionArguments[ind].ArrayValue[thisArrayIndex].Type == ARG_TYPE_BOOLEAN) {
												thisVar.Type = VARIABLE_TYPE_BOOLEAN
												thisVar.BooleanValue = functionArguments[ind].ArrayValue[thisArrayIndex].BooleanValue
											} else {
												thisVar.Type = VARIABLE_TYPE_NONE
											}

											(*globalVariableArray)[varIndex].ArrayValue = append((*globalVariableArray)[varIndex].ArrayValue, thisVar)
										}
									} else {
										//Nil
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_NONE
									}
								}

								var thisGotReturn bool = false
								var thisReturnToken Token
								var thisNeedBreak bool = false
								var thisStackReference []Token
								
								//execute user defined function
								prsr := Parser{}
								parserErr := prsr.Parse((*globalFunctionArray)[funcIndex].Tokens, globalVariableArray, globalFunctionArray, thisScopeName, globalNativeVarList, &thisGotReturn, &thisReturnToken, false, &thisNeedBreak, &thisStackReference)
						
								if(parserErr != nil) {
									return parserErr
								}

								if(thisGotReturn) {
									//the function returns a value
									newToken = thisReturnToken
								}

								//TODO: NEED CLEANUP OF VARIABLES BELOW
								//DELETE GENERATED VARIABLES with thisScopeName
							}

							stack = append(stack, newToken)
						
						} else if(currentToken.Type == TOKEN_TYPE_OPEN_BRACES) {
							//array declaration
							processedArg := 0
							var tempArray []Token
							newToken := currentToken
							newToken.Type = TOKEN_TYPE_ARRAY
							if(currentToken.OtherInt > 0) {
								if(len(stack) == 0) {
									if(len(outputQueue) > 0) {
										return errors.New(SyntaxErrorMessage(outputQueue[0].Line, outputQueue[0].Column, "Unexpected token '" + outputQueue[0].Value + "'", outputQueue[0].FileName))
									} else {
										return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Unexpected token '" + currentToken.Value + "'", currentToken.FileName))
									}
								}

								for true {
									var param Token

									param = stack[len(stack)-1]
									stack = stack[:len(stack)-1]

									var errConvert error
									if(param.Type == TOKEN_TYPE_IDENTIFIER) {
										param, errConvert = convertVariableToToken(param, *globalVariableArray, scopeName)
										if(errConvert != nil) {
											return errConvert
										}
									}

									if(param.Type == TOKEN_TYPE_ARRAY) {
										return errors.New(SyntaxErrorMessage(param.Line, param.Column, "Unexpected token '" + param.Value + "'", param.FileName))
									}

									tempArray = append(tempArray, param)

									processedArg += 1
									if (processedArg == currentToken.OtherInt) {
										break
									}
								}

								if(len(tempArray) > 0) {
									//reverse the array (so it's in proper position)
									arrayLength := len(tempArray)
									for thisArrayIndex := 0; thisArrayIndex < arrayLength; thisArrayIndex++ {
										thisToken := tempArray[len(tempArray) - 1]
										tempArray = tempArray[:len(tempArray)-1]
										newToken.Array = append(newToken.Array, thisToken)
									}
								}
							}
							stack = append(stack, newToken)
							//DumpToken(stack)
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
						
						} else if(currentToken.Type == TOKEN_TYPE_FUNCTION_RETURN) {
							if(scopeName == "main") {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "'rtn' outside function", currentToken.FileName))
							}

							returnValue := stack[len(stack)-1]
							stack = stack[:len(stack)-1]
							var errConvert error

							if(returnValue.Type == TOKEN_TYPE_IDENTIFIER) {
								returnValue, errConvert = convertVariableToToken(returnValue, *globalVariableArray, scopeName)
								if(errConvert != nil) {
									return errConvert
								}
							}

							errValue := expectedTokenTypes(returnValue, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE, TOKEN_TYPE_BOOLEAN, TOKEN_TYPE_ARRAY)
							if (errValue != nil) {
								return errValue
							}

							*gotReturn = true
							*returnToken = returnValue
							return nil
						} else if(currentToken.Type == TOKEN_TYPE_FOR_LOOP_START) {
							//get looping parameters
							//param validation
							if(len(stack) == 0 || len(stack) < 2) {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, currentToken.Value + " takes exactly 2 argument", currentToken.FileName))
							}
							var errConvert error

							param2 := stack[len(stack)-1]
							stack = stack[:len(stack)-1]

							//if param2 is an identifier
							//then it's a variable
							if(param2.Type == TOKEN_TYPE_IDENTIFIER) {
								param2, errConvert = convertVariableToToken(param2, *globalVariableArray, scopeName)
								if(errConvert != nil) {
									return errConvert
								}
							}
		
							param1 := stack[len(stack)-1]
							stack = stack[:len(stack)-1]

							//if param1 is an identifier
							//then it's a variable
							if(param1.Type == TOKEN_TYPE_IDENTIFIER) {
								param1, errConvert = convertVariableToToken(param1, *globalVariableArray, scopeName)
								if(errConvert != nil) {
									return errConvert
								}
							}
		
							//validate param2
							errParam2 := expectedTokenTypes(param2, TOKEN_TYPE_INTEGER)
							if (errParam2 != nil) {
								return errParam2
							}
							//validate param1
							errParam1 := expectedTokenTypes(param1, TOKEN_TYPE_INTEGER)
							if (errParam1 != nil) {
								return errParam1
							}

							var tempTokens []Token
							openLoopCount = 0
							//append all tokens to temporary token
							for true {
								currentToken := outputQueue[0]
								outputQueue = append(outputQueue[:0], outputQueue[1:]...)

								tempTokens = append(tempTokens, currentToken)

								if(currentToken.Type == TOKEN_TYPE_FOR_LOOP_START) {
									openLoopCount += 1
								}

								if(currentToken.Type == TOKEN_TYPE_FOR_LOOP_END) {
									if(openLoopCount == 0) {
										break
									}
									openLoopCount = openLoopCount - 1
								}

								if(len(outputQueue) == 0) {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Invalid statement", currentToken.FileName))
								}
							}

							iParam1, _ := strconv.Atoi(param1.Value)
							iParam2, _ := strconv.Atoi(param2.Value)

							for loop_index := iParam1 ; loop_index <= iParam2; loop_index++ {

								var loopGotReturn bool = false
								var loopReturnToken Token
								var loopNeedBreak bool = false
								var loopStackReference []Token
								
								prsr := Parser{}
								parserErr := prsr.Parse(tempTokens, globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, &loopGotReturn, &loopReturnToken, true, &loopNeedBreak, &loopStackReference)
						
								if(parserErr != nil) {
									return parserErr
								}
								if(loopGotReturn) {
									*gotReturn = loopGotReturn
									*returnToken = loopReturnToken
									return nil
								}

								if(loopNeedBreak) {
									//break the loop
									break
								}
							}

							stack = append(stack, currentToken) //append TOKEN_TYPE_FOR_LOOP_END
						
						} else if(currentToken.Type == TOKEN_TYPE_LOOP_BREAK) {
							if(!isLoop) {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "'brk' outside loop", currentToken.FileName))
							}
							*needBreak = true
							return nil
						} else if(currentToken.Type == TOKEN_TYPE_IF_START) {
							var params []Token
							var gotElse bool = false

							//param := stack[len(stack)-1]

							params = append(params, stack[len(stack)-1])
							stack = stack[:len(stack)-1]

							var tempTokens []TokenArray
							openIfCount = 0
							//append all tokens to temporary token

							tempTokens = append(tempTokens, TokenArray{})
							for true {
								currentToken := outputQueue[0]
								outputQueue = append(outputQueue[:0], outputQueue[1:]...)

								tempTokens[len(tempTokens)-1].Tokens = append(tempTokens[len(tempTokens)-1].Tokens, currentToken)

								if(currentToken.Type == TOKEN_TYPE_IF_START) {
									openIfCount += 1
								}

								if(currentToken.Type == TOKEN_TYPE_IF_END || currentToken.Type == TOKEN_TYPE_ELSE || currentToken.Type == TOKEN_TYPE_ELIF_START) {
									if(openIfCount == 0) {
										if(currentToken.Type == TOKEN_TYPE_IF_END) {
											break
										} else if (currentToken.Type == TOKEN_TYPE_ELIF_START) {
											//else if
											//get parameters until end
											var paramTokens []Token
											var currentToken2 Token
											var will_change_to_main string = currentToken.Context
											for true {
												currentToken2 = outputQueue[0]
												outputQueue = append(outputQueue[:0], outputQueue[1:]...)

												if(currentToken2.Type == TOKEN_TYPE_ELIF_PARAM_END) {
													break
												} else {
													if(currentToken2.Context == will_change_to_main) {
														currentToken2.Context = "main_context"
													}
													paramTokens = append(paramTokens, currentToken2)
												}

												if(len(outputQueue) == 0) {
													return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Invalid statement", currentToken.FileName))
												}
											}
											//evaluate paramTokens
											var efGotReturn bool = false
											var efReturnToken Token
											var efNeedBreak bool = false
											var efStackReference []Token
											
											//add newline to paramTokens so parser can execute it
											paramTokens = append(paramTokens, Token{Value: "\n", FileName: currentToken2.FileName, Type: TOKEN_TYPE_NEWLINE, Line: currentToken2.Line, Column: currentToken2.Column, Context: "main_context" })

											prsr := Parser{}
											parserErr := prsr.Parse(paramTokens, globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, &efGotReturn, &efReturnToken, isLoop, &efNeedBreak, &efStackReference)
									
											if(parserErr != nil) {
												return parserErr
											}
											//end of evaluate paramTokens
											params = append(params, efStackReference[0])
											tempTokens = append(tempTokens, TokenArray{})
											continue											
										} else {
											//else
											gotElse = true
											tempTokens = append(tempTokens, TokenArray{})
											continue
										}
									}
									if(currentToken.Type == TOKEN_TYPE_IF_END) {
										openIfCount = openIfCount - 1
									}
								}

								if(len(outputQueue) == 0) {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Invalid statement", currentToken.FileName))
								}
							}


							var currentStatementIndex int = 0
							var executeIfStatement bool = false

							for paramIndex := 0; paramIndex < len(params); paramIndex++ {
								var thisParam Token = params[paramIndex]
								var errConvert error

								if(thisParam.Type == TOKEN_TYPE_IDENTIFIER) {
									thisParam, errConvert = convertVariableToToken(thisParam, *globalVariableArray, scopeName)
									if(errConvert != nil) {
										return errConvert
									}
								}

								//validate param
								errParam := expectedTokenTypes(thisParam, TOKEN_TYPE_BOOLEAN)
								if (errParam != nil) {
									return errParam
								}

								paramBool := convertTokenToBool(thisParam)

								if(paramBool) {
									executeIfStatement = true
									currentStatementIndex = paramIndex
									break
								}
							}

							if(!executeIfStatement) {
								if(gotElse && len(tempTokens) > 1) {
									executeIfStatement = true
									currentStatementIndex = len(tempTokens) - 1
								}
							}

							if(executeIfStatement) {
								var ifGotReturn bool = false
								var ifReturnToken Token
								var ifNeedBreak bool = false
								var ifStackReference []Token
								
								prsr := Parser{}
								parserErr := prsr.Parse(tempTokens[currentStatementIndex].Tokens, globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, &ifGotReturn, &ifReturnToken, isLoop, &ifNeedBreak, &ifStackReference)
						
								if(parserErr != nil) {
									return parserErr
								}
								if(ifGotReturn) {
									*gotReturn = ifGotReturn
									*returnToken = ifReturnToken
									return nil
								}

								if(isLoop && ifNeedBreak) {
									*needBreak = ifNeedBreak
									return nil
								}
							}

							stack = append(stack, currentToken) //TOKEN_TYPE_IF_END
						} else {
							stack = append(stack, currentToken)
						}
					}

					*stackReference = stack

					if(len(stack) > 1) {
						return errors.New(SyntaxErrorMessage(stack[0].Line, stack[0].Column, "Invalid statement", stack[0].FileName))
					} /* else {
						if(stack[0].Type == TOKEN_TYPE_IDENTIFIER) {
							return errors.New(SyntaxErrorMessage(stack[0].Line, stack[0].Column, "Unexpected token '" + stack[0].Value + "'", stack[0].FileName))
						}
					}
					*/
				}
				
			}

		} else if(tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_START || tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_END || tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_START || tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_END || tokenArray[x].Type == TOKEN_TYPE_IF_START || tokenArray[x].Type == TOKEN_TYPE_IF_END) {
			if(tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_START) {
				ignoreNewline = true
				isFunctionDefinition = true
			} else if(tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_END) {
				//TOKEN_TYPE_FUNCTION_DEF_END
				ignoreNewline = false
				isFunctionDefinition = false
				
				//NOTE: not sure if the code below is temporary
				//append newline (to make the one liner definition of function works)
				tokensToEvaluate = append(tokensToEvaluate, Token{Value: "\n", FileName: tokenArray[x].FileName, Type: TOKEN_TYPE_NEWLINE, Line: tokenArray[x].Line, Column: tokenArray[x].Column })
			}
			if(!isFunctionDefinition) {
				if(!isIfStatement) {
					if(tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_START) {
						openLoopCount += 1
						ignoreNewline = true
						isLoopStatement = true
					} else if(tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_END) {
						openLoopCount = openLoopCount - 1

						if(openLoopCount == 0) {
							ignoreNewline = false
							isLoopStatement = false
							tokensToEvaluate = append(tokensToEvaluate, Token{Value: "\n", FileName: tokenArray[x].FileName, Type: TOKEN_TYPE_NEWLINE, Line: tokenArray[x].Line, Column: tokenArray[x].Column })
						}
					}
				}
				if(!isLoopStatement) {
					if(tokenArray[x].Type == TOKEN_TYPE_IF_START) {
						openIfCount += 1
						ignoreNewline = true
						isIfStatement = true
					} else if(tokenArray[x].Type == TOKEN_TYPE_IF_END) {
						openIfCount = openIfCount - 1

						if(openIfCount == 0) {
							ignoreNewline = false
							isIfStatement = false
							tokensToEvaluate = append(tokensToEvaluate, Token{Value: "\n", FileName: tokenArray[x].FileName, Type: TOKEN_TYPE_NEWLINE, Line: tokenArray[x].Line, Column: tokenArray[x].Column })
						}
					}
				}
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