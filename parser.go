package main

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
	//fmt
)

func generateRandomNumbers() string {
	rand.Seed(time.Now().UnixNano())

	return strconv.Itoa(rand.Int())
}

func expectedTokenTypes(token Token, tokenTypes ...int) error {
	isOK := false

	for x := 0; x < len(tokenTypes); x++ {
		if token.Type == tokenTypes[x] {
			isOK = true
			break
		}
	}

	if !isOK {
		return errors.New(SyntaxErrorMessage(token.Line, token.Column, "Invalid operand '"+token.Value+"'", token.FileName))
	}

	return nil
}

func PopStack(stack *[]Token) Token {
	ret := (*stack)[len((*stack))-1]
	*stack = (*stack)[:len((*stack))-1]
	return ret
}

type Parser struct {
}

func (parser Parser) Parse(tokenArray []Token, globalVariableArray *[]Variable, globalFunctionArray *[]Function, scopeName string, globalNativeVarList *[]string, gotReturn *bool, returnToken *Token, isLoop bool, needBreak *bool, stackReference *[]Token, globalSettings *GlobalSettingsObject, getLastStackBool bool, lastStackBool *bool) error {
	var tokensToEvaluate []Token
	operatorPrecedences := map[string]int{"function_return": 0, ":": 1, "=": 1, "+": 2, "-": 2, "&": 2, "|": 2, "==": 2, "<>": 2, ">": 2, "<": 2, ">=": 2, "<=": 2, "/": 3, "*": 3} //operator order of precedences
	currentContext := CONTEXT_NAME_MAIN
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
	var isWhileLoopStatement bool = false
	var isIfStatement bool = false
	var openLoopCount int = 0
	var openWhileLoopCount int = 0
	var openIfCount int = 0
	var whileLoopTokens map[string][]Token
	var whileLoopIsDone map[string]bool
	whileLoopTokens = make(map[string][]Token)
	whileLoopIsDone = make(map[string]bool)
	var lastToken Token

	for x := 0; x < len(tokenArray); x++ {
		if tokenArray[x].Type == TOKEN_TYPE_NEWLINE {

			if ignoreNewline {
				//put the token to stack for shunting yard process later
				tokensToEvaluate = append(tokensToEvaluate, tokenArray[x])

				if len(tokenArray) == (x + 1) {
					return errors.New(SyntaxErrorMessage(tokenArray[x].Line, tokenArray[x].Column, "Unfinished statement", tokenArray[x].FileName))
				}

				continue
			}

			//execute shunting yard
			if len(tokensToEvaluate) > 0 {
				if tokensToEvaluate[0].Type == TOKEN_TYPE_PLUS || tokensToEvaluate[0].Type == TOKEN_TYPE_MINUS || tokensToEvaluate[0].Type == TOKEN_TYPE_DIVIDE || tokensToEvaluate[0].Type == TOKEN_TYPE_MULTIPLY || tokensToEvaluate[0].Type == TOKEN_TYPE_EQUALS || tokensToEvaluate[0].Type == TOKEN_TYPE_AMPERSAND || tokensToEvaluate[0].Type == TOKEN_TYPE_OR || tokensToEvaluate[0].Type == TOKEN_TYPE_EQUALITY || tokensToEvaluate[0].Type == TOKEN_TYPE_INEQUALITY || tokensToEvaluate[0].Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS || tokensToEvaluate[0].Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS || tokensToEvaluate[0].Type == TOKEN_TYPE_GREATER_THAN || tokensToEvaluate[0].Type == TOKEN_TYPE_LESS_THAN || tokensToEvaluate[0].Type == TOKEN_TYPE_COLON {
					//syntax error if the first token is an operator
					return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '"+tokensToEvaluate[0].Value+"'", tokensToEvaluate[0].FileName))
				}

				if tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_PLUS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_MINUS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_DIVIDE || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_MULTIPLY || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_EQUALS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_AMPERSAND || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_OR || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_EQUALITY || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_INEQUALITY || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_GREATER_THAN || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_LESS_THAN || tokensToEvaluate[len(tokensToEvaluate)-1].Type == TOKEN_TYPE_COLON {
					//syntax error if the last token is an operator
					return errors.New(SyntaxErrorMessage(tokensToEvaluate[len(tokensToEvaluate)-1].Line, tokensToEvaluate[len(tokensToEvaluate)-1].Column, "Unfinished statement", tokensToEvaluate[len(tokensToEvaluate)-1].FileName))
				}

				justAddTokens = false
				isFunctionDefinition = false
				isLoopStatement = false
				isWhileLoopStatement = false
				isIfStatement = false
				openIfCount = 0
				openLoopCount = 0
				openWhileLoopCount = 0
				//shunting-yard
				for len(tokensToEvaluate) > 0 {
					currentToken := tokensToEvaluate[0]
					tokensToEvaluate = append(tokensToEvaluate[:0], tokensToEvaluate[1:]...) //pop the first element
					isValidToken := false
					currentContext = currentToken.Context

					if justAddTokens {
						//function body, just add to outputqueue
						outputQueue = append(outputQueue, currentToken)

						if currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_END {
							justAddTokens = false
							isFunctionDefinition = false
						}
						if !isFunctionDefinition {
							if !isIfStatement {
								if currentToken.Type == TOKEN_TYPE_FOR_LOOP_START {
									openLoopCount += 1
								}
								if currentToken.Type == TOKEN_TYPE_FOR_LOOP_END {
									if openLoopCount == 0 {
										justAddTokens = false
									}
									if openLoopCount > 0 {
										openLoopCount = openLoopCount - 1
									}
								}
								if currentToken.Type == TOKEN_TYPE_WHILE_LOOP_START {
									openWhileLoopCount += 1
								}
								if currentToken.Type == TOKEN_TYPE_WHILE_LOOP_END {
									if openWhileLoopCount == 0 {
										justAddTokens = false
									}
									if openWhileLoopCount > 0 {
										openWhileLoopCount = openWhileLoopCount - 1
									}
								}
							}
							if !isLoopStatement && !isWhileLoopStatement {
								if currentToken.Type == TOKEN_TYPE_IF_START {
									openIfCount += 1
								}
								if currentToken.Type == TOKEN_TYPE_IF_END {
									if openIfCount == 0 {
										justAddTokens = false
									}
									if openIfCount > 0 {
										openIfCount = openIfCount - 1
									}
								}
							}
						}
						continue
					}

					if currentToken.Type == TOKEN_TYPE_NEWLINE {
						//just ignore newline
						continue
					}

					if currentToken.Type == TOKEN_TYPE_INTEGER || currentToken.Type == TOKEN_TYPE_FLOAT || currentToken.Type == TOKEN_TYPE_IDENTIFIER || currentToken.Type == TOKEN_TYPE_STRING || currentToken.Type == TOKEN_TYPE_LOOP_BREAK {
						//If it's a number or identifier, add it to queue, (ADD TOKEN_TYPE_KEYWORD AND string and other acceptable tokens LATER)
						outputQueue = append(outputQueue, currentToken)
						isValidToken = true
					}

					//dontIgnorePopping := true
					if currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION || currentToken.Type == TOKEN_TYPE_FUNCTION || currentToken.Type == TOKEN_TYPE_COMMA || currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START || currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END || currentToken.Type == TOKEN_TYPE_FOR_LOOP_START || currentToken.Type == TOKEN_TYPE_FOR_LOOP_PARAM_END || currentToken.Type == TOKEN_TYPE_IF_START || currentToken.Type == TOKEN_TYPE_IF_PARAM_END || currentToken.Type == TOKEN_TYPE_CLOSE_BRACES || currentToken.Type == TOKEN_TYPE_OPEN_BRACES || currentToken.Type == TOKEN_TYPE_GET_ARRAY_START || currentToken.Type == TOKEN_TYPE_GET_ARRAY_END || currentToken.Type == TOKEN_TYPE_CLOSE_BRACKET || currentToken.Type == TOKEN_TYPE_OPEN_BRACKET || currentToken.Type == TOKEN_TYPE_WHILE_LOOP_START || currentToken.Type == TOKEN_TYPE_WHILE_LOOP_PARAM_END {
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

						if currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION || currentToken.Type == TOKEN_TYPE_COMMA || currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END || currentToken.Type == TOKEN_TYPE_FOR_LOOP_PARAM_END || currentToken.Type == TOKEN_TYPE_IF_PARAM_END || currentToken.Type == TOKEN_TYPE_CLOSE_BRACES || currentToken.Type == TOKEN_TYPE_GET_ARRAY_END || currentToken.Type == TOKEN_TYPE_CLOSE_BRACKET || currentToken.Type == TOKEN_TYPE_WHILE_LOOP_PARAM_END {
							//pop all operators from operator stack to output queue before the function
							//NOTE: don't include '=' (NOT SURE)
							for true {
								if len(operatorStack[currentContext]) > 0 {
									if operatorStack[currentContext][len(operatorStack[currentContext])-1].Type == TOKEN_TYPE_EQUALS {
										break
									} else {
										outputQueue = append(outputQueue, operatorStack[currentContext][len(operatorStack[currentContext])-1])
										operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
									}
								} else {
									break
								}
							}
						}

						if currentToken.Type == TOKEN_TYPE_FUNCTION || currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START || currentToken.Type == TOKEN_TYPE_FOR_LOOP_START || currentToken.Type == TOKEN_TYPE_IF_START || currentToken.Type == TOKEN_TYPE_OPEN_BRACES || currentToken.Type == TOKEN_TYPE_GET_ARRAY_START || currentToken.Type == TOKEN_TYPE_OPEN_BRACKET || currentToken.Type == TOKEN_TYPE_WHILE_LOOP_START {
							functionStack = append(functionStack, currentToken)
							if currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START {
								isFunctionDefinition = true
							}
							if !isFunctionDefinition {
								if !isIfStatement {
									if currentToken.Type == TOKEN_TYPE_FOR_LOOP_START {
										isLoopStatement = true
									}
									if currentToken.Type == TOKEN_TYPE_WHILE_LOOP_START {
										isWhileLoopStatement = true
									}
								}
								if !isLoopStatement && !isWhileLoopStatement {
									if currentToken.Type == TOKEN_TYPE_IF_START {
										isIfStatement = true
									}
								}
							}
							if currentToken.Type == TOKEN_TYPE_OPEN_BRACES || currentToken.Type == TOKEN_TYPE_OPEN_BRACKET {
								if len(tokensToEvaluate) > 0 {
									if currentToken.Type == TOKEN_TYPE_OPEN_BRACES {
										if tokensToEvaluate[0].Type != TOKEN_TYPE_CLOSE_BRACES {
											arrayArgCount[currentContext] = 1
											if tokensToEvaluate[0].Type == TOKEN_TYPE_OPEN_BRACES {
												return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '"+tokensToEvaluate[0].Value+"'", tokensToEvaluate[0].FileName))
											}
										} else {
											arrayArgCount[currentContext] = 0
										}
									}
									if currentToken.Type == TOKEN_TYPE_OPEN_BRACKET {
										if tokensToEvaluate[0].Type != TOKEN_TYPE_CLOSE_BRACKET {
											arrayArgCount[currentContext] = 1
											if tokensToEvaluate[0].Type == TOKEN_TYPE_OPEN_BRACKET {
												return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '"+tokensToEvaluate[0].Value+"'", tokensToEvaluate[0].FileName))
											}
										} else {
											arrayArgCount[currentContext] = 0
										}
									}
								} else {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Unfinished statement", currentToken.FileName))
								}
							}
							if currentToken.Type == TOKEN_TYPE_GET_ARRAY_START {
								if len(tokensToEvaluate) > 0 {
									if tokensToEvaluate[0].Type == TOKEN_TYPE_GET_ARRAY_END {
										return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '"+tokensToEvaluate[0].Value+"'", tokensToEvaluate[0].FileName))
									}
								} else {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Unfinished statement", currentToken.FileName))
								}
							}
						} else if currentToken.Type == TOKEN_TYPE_INVOKE_FUNCTION || currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END || currentToken.Type == TOKEN_TYPE_FOR_LOOP_PARAM_END || currentToken.Type == TOKEN_TYPE_IF_PARAM_END || currentToken.Type == TOKEN_TYPE_CLOSE_BRACES || currentToken.Type == TOKEN_TYPE_GET_ARRAY_END || currentToken.Type == TOKEN_TYPE_CLOSE_BRACKET || currentToken.Type == TOKEN_TYPE_WHILE_LOOP_PARAM_END {
							tokenToAppend := functionStack[len(functionStack)-1]

							if currentToken.Type == TOKEN_TYPE_CLOSE_BRACES || currentToken.Type == TOKEN_TYPE_CLOSE_BRACKET {
								tokenToAppend.OtherInt = arrayArgCount[currentContext]
							}

							outputQueue = append(outputQueue, tokenToAppend)
							functionStack = functionStack[:len(functionStack)-1]

							if currentToken.Type == TOKEN_TYPE_FUNCTION_PARAM_END || currentToken.Type == TOKEN_TYPE_FOR_LOOP_PARAM_END || currentToken.Type == TOKEN_TYPE_IF_PARAM_END || currentToken.Type == TOKEN_TYPE_WHILE_LOOP_PARAM_END {
								//next is function body
								justAddTokens = true
							}
						} else {
							//comma
							//count parameter (currently for array only, not sure in the future (for function params?))
							arrayArgCount[currentContext] += 1
							//validate separator
							if len(tokensToEvaluate) > 0 {
								if tokensToEvaluate[0].Type != TOKEN_TYPE_FLOAT && tokensToEvaluate[0].Type != TOKEN_TYPE_INTEGER && tokensToEvaluate[0].Type != TOKEN_TYPE_STRING && tokensToEvaluate[0].Type != TOKEN_TYPE_IDENTIFIER && tokensToEvaluate[0].Type != TOKEN_TYPE_FUNCTION && tokensToEvaluate[0].Type != TOKEN_TYPE_OPEN_PARENTHESIS && tokensToEvaluate[0].Type != TOKEN_TYPE_OPEN_BRACES && tokensToEvaluate[0].Type != TOKEN_TYPE_GET_ARRAY_START && tokensToEvaluate[0].Type != TOKEN_TYPE_OPEN_BRACKET {
									return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '"+tokensToEvaluate[0].Value+"'", tokensToEvaluate[0].FileName))
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

					if currentToken.Type == TOKEN_TYPE_PLUS || currentToken.Type == TOKEN_TYPE_MINUS || currentToken.Type == TOKEN_TYPE_DIVIDE || currentToken.Type == TOKEN_TYPE_MULTIPLY || currentToken.Type == TOKEN_TYPE_EQUALS || currentToken.Type == TOKEN_TYPE_AMPERSAND || currentToken.Type == TOKEN_TYPE_OR || currentToken.Type == TOKEN_TYPE_FUNCTION_RETURN || currentToken.Type == TOKEN_TYPE_EQUALITY || currentToken.Type == TOKEN_TYPE_INEQUALITY || currentToken.Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS || currentToken.Type == TOKEN_TYPE_LESS_THAN || currentToken.Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS || currentToken.Type == TOKEN_TYPE_GREATER_THAN || currentToken.Type == TOKEN_TYPE_COLON {
						//the token is operator

						//add validation
						if len(tokensToEvaluate) > 0 {
							//if next token is an operator then raise an error
							if tokensToEvaluate[0].Type == TOKEN_TYPE_PLUS || tokensToEvaluate[0].Type == TOKEN_TYPE_MINUS || tokensToEvaluate[0].Type == TOKEN_TYPE_DIVIDE || tokensToEvaluate[0].Type == TOKEN_TYPE_MULTIPLY || tokensToEvaluate[0].Type == TOKEN_TYPE_EQUALS || tokensToEvaluate[0].Type == TOKEN_TYPE_AMPERSAND || tokensToEvaluate[0].Type == TOKEN_TYPE_OR || tokensToEvaluate[0].Type == TOKEN_TYPE_FUNCTION_RETURN || tokensToEvaluate[0].Type == TOKEN_TYPE_EQUALITY || tokensToEvaluate[0].Type == TOKEN_TYPE_INEQUALITY || tokensToEvaluate[0].Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS || tokensToEvaluate[0].Type == TOKEN_TYPE_LESS_THAN || tokensToEvaluate[0].Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS || tokensToEvaluate[0].Type == TOKEN_TYPE_GREATER_THAN || tokensToEvaluate[0].Type == TOKEN_TYPE_COLON {
								return errors.New(SyntaxErrorMessage(tokensToEvaluate[0].Line, tokensToEvaluate[0].Column, "Unexpected token '"+tokensToEvaluate[0].Value+"'", tokensToEvaluate[0].FileName))
							}
						}

						for true {
							if len(operatorStack[currentContext]) > 0 {

								if currentToken.Type == TOKEN_TYPE_FUNCTION_RETURN {
									if operatorPrecedences[operatorStack[currentContext][len(operatorStack[currentContext])-1].Value] >= operatorPrecedences["function_return"] {
										outputQueue = append(outputQueue, operatorStack[currentContext][len(operatorStack[currentContext])-1])
										operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
									} else {
										break
									}
								} else {
									if operatorPrecedences[operatorStack[currentContext][len(operatorStack[currentContext])-1].Value] >= operatorPrecedences[currentToken.Value] {
										outputQueue = append(outputQueue, operatorStack[currentContext][len(operatorStack[currentContext])-1])
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

					if currentToken.Type == TOKEN_TYPE_OPEN_PARENTHESIS {
						//if it's an open parenthesis '(' push it onto the stack
						operatorStack[currentContext] = append(operatorStack[currentContext], currentToken)
						isValidToken = true
					}

					if currentToken.Type == TOKEN_TYPE_CLOSE_PARENTHESIS {
						isValidToken = true
						//close parenthesis
						if len(operatorStack[currentContext]) > 0 {
							for true {
								if operatorStack[currentContext][len(operatorStack[currentContext])-1].Type != TOKEN_TYPE_OPEN_PARENTHESIS {
									outputQueue = append(outputQueue, operatorStack[currentContext][len(operatorStack[currentContext])-1])
									operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
								} else {
									operatorStack[currentContext] = operatorStack[currentContext][:len(operatorStack[currentContext])-1]
									break
								}

								if len(operatorStack[currentContext]) == 0 {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Operator expected", currentToken.FileName))
								}
							}
						} else {
							return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Operator expected", currentToken.FileName))
						}
					}

					if !isValidToken {
						return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Unexpected token '"+currentToken.Value+"'", currentToken.FileName))
					}
				}

				for len(operatorStack[CONTEXT_NAME_MAIN]) > 0 {
					if operatorStack[CONTEXT_NAME_MAIN][len(operatorStack[CONTEXT_NAME_MAIN])-1].Type == TOKEN_TYPE_OPEN_PARENTHESIS {
						return errors.New(SyntaxErrorMessage(operatorStack[CONTEXT_NAME_MAIN][len(operatorStack[CONTEXT_NAME_MAIN])-1].Line, operatorStack[CONTEXT_NAME_MAIN][len(operatorStack[CONTEXT_NAME_MAIN])-1].Column, "Unexpected token '"+operatorStack[CONTEXT_NAME_MAIN][len(operatorStack[CONTEXT_NAME_MAIN])-1].Value+"'", operatorStack[CONTEXT_NAME_MAIN][len(operatorStack[CONTEXT_NAME_MAIN])-1].FileName))
					}
					outputQueue = append(outputQueue, operatorStack[CONTEXT_NAME_MAIN][len(operatorStack[CONTEXT_NAME_MAIN])-1])
					operatorStack[CONTEXT_NAME_MAIN] = operatorStack[CONTEXT_NAME_MAIN][:len(operatorStack[CONTEXT_NAME_MAIN])-1]
				}
				//end of shunting-yard

				//validate end of function
				if len(functionStack) > 0 {
					return errors.New(SyntaxErrorMessage(functionStack[0].Line, functionStack[0].Column, "End of function call expected", functionStack[0].FileName))
				}

				//DumpToken(outputQueue)
				//the outputQueue contains the reverse polish notation

				if len(outputQueue) > 0 {
					//read the reverse polish below
					var stack []Token

					for len(outputQueue) > 0 {
						currentToken := outputQueue[0]
						outputQueue = append(outputQueue[:0], outputQueue[1:]...) //pop the first element

						if currentToken.Type == TOKEN_TYPE_PLUS || currentToken.Type == TOKEN_TYPE_MINUS || currentToken.Type == TOKEN_TYPE_DIVIDE || currentToken.Type == TOKEN_TYPE_MULTIPLY {
							//arithmetic operation
							//NOTE: ASSUME THAT RIGHT OPERAND AND LEFT OPERAND ARE INTEGER AND FLOAT ONLY (NO IDENTIFIER, STRING ETC... (TEMPORARY ONLY)
							rightOperand := PopStack(&stack)
							var tempRightInt int
							var tempRightFloat float64
							var tempRightString string

							leftOperand := PopStack(&stack)
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

							if leftOperand.Type == TOKEN_TYPE_IDENTIFIER {
								leftOperand, errConvert = convertVariableToToken(leftOperand, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}
							if rightOperand.Type == TOKEN_TYPE_IDENTIFIER {
								rightOperand, errConvert = convertVariableToToken(rightOperand, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							//convert operands to its designated type
							if leftOperand.Type == TOKEN_TYPE_INTEGER {
								//convert to integer
								result.Type = TOKEN_TYPE_INTEGER
								tempLeftInt, _ = strconv.Atoi(leftOperand.Value)
								tempRightInt, _ = strconv.Atoi(rightOperand.Value)
							} else if leftOperand.Type == TOKEN_TYPE_STRING {
								//string
								result.Type = TOKEN_TYPE_STRING
								tempLeftString = leftOperand.Value
								tempRightString = rightOperand.Value
							} else {
								//let's assume that it should be converted to float (for now)
								result.Type = TOKEN_TYPE_FLOAT
								tempLeftFloat, _ = strconv.ParseFloat(leftOperand.Value, 64)
								tempRightFloat, _ = strconv.ParseFloat(rightOperand.Value, 64)
							}

							if currentToken.Type == TOKEN_TYPE_PLUS {
								//either addition or concatenation

								//validate left operand
								errLeft := expectedTokenTypes(leftOperand, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING)
								if errLeft != nil {
									return errLeft
								}
								//validate right operand
								errRight := expectedTokenTypes(rightOperand, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING)
								if errRight != nil {
									return errRight
								}

								if leftOperand.Type == TOKEN_TYPE_INTEGER {
									result.Value = strconv.Itoa(tempLeftInt + tempRightInt)
								} else if leftOperand.Type == TOKEN_TYPE_STRING {
									result.Value = tempLeftString + tempRightString //concatenate
								} else {
									//let's assume it's float
									result.Value = strconv.FormatFloat(tempLeftFloat+tempRightFloat, 'f', -1, 64)
								}

							} else {
								//substraction, division and multiplication

								//validate left operand (No String)
								errLeft := expectedTokenTypes(leftOperand, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT)
								if errLeft != nil {
									return errLeft
								}
								//validate right operand (No String)
								errRight := expectedTokenTypes(rightOperand, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT)
								if errRight != nil {
									return errRight
								}

								if currentToken.Type == TOKEN_TYPE_MINUS {
									//substraction
									if leftOperand.Type == TOKEN_TYPE_INTEGER {
										result.Value = strconv.Itoa(tempLeftInt - tempRightInt)
									} else {
										//let's assume it's float
										result.Value = strconv.FormatFloat(tempLeftFloat-tempRightFloat, 'f', -1, 64)
									}
								} else if currentToken.Type == TOKEN_TYPE_MULTIPLY {
									//multiplication
									if leftOperand.Type == TOKEN_TYPE_INTEGER {
										result.Value = strconv.Itoa(tempLeftInt * tempRightInt)
									} else {
										//let's assume it's float
										result.Value = strconv.FormatFloat(tempLeftFloat*tempRightFloat, 'f', -1, 64)
									}
								} else {
									//assume it's division
									if leftOperand.Type == TOKEN_TYPE_INTEGER {
										if tempRightInt == 0 {
											return errors.New(SyntaxErrorMessage(rightOperand.Line, rightOperand.Column, "Division by zero", rightOperand.FileName))
										}
										result.Value = strconv.Itoa(tempLeftInt / tempRightInt)
									} else {
										//let's assume it's float
										if tempRightInt == 0 {
											return errors.New(SyntaxErrorMessage(rightOperand.Line, rightOperand.Column, "Division by zero", rightOperand.FileName))
										}
										result.Value = strconv.FormatFloat(tempLeftFloat/tempRightFloat, 'f', -1, 64)
									}
								}
							}

							stack = append(stack, result)

						} else if currentToken.Type == TOKEN_TYPE_AMPERSAND || currentToken.Type == TOKEN_TYPE_OR || currentToken.Type == TOKEN_TYPE_EQUALITY || currentToken.Type == TOKEN_TYPE_INEQUALITY || currentToken.Type == TOKEN_TYPE_GREATER_THAN || currentToken.Type == TOKEN_TYPE_LESS_THAN || currentToken.Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS || currentToken.Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS {
							//logical or comparison operation
							rightOperand := PopStack(&stack)

							leftOperand := PopStack(&stack)

							var errConvert error

							result := leftOperand

							if leftOperand.Type == TOKEN_TYPE_IDENTIFIER {
								leftOperand, errConvert = convertVariableToToken(leftOperand, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}
							if rightOperand.Type == TOKEN_TYPE_IDENTIFIER {
								rightOperand, errConvert = convertVariableToToken(rightOperand, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							if currentToken.Type == TOKEN_TYPE_AMPERSAND || currentToken.Type == TOKEN_TYPE_OR {
								//validation for logical operations
								//validate left operand
								errLeft := expectedTokenTypes(leftOperand, TOKEN_TYPE_BOOLEAN)
								if errLeft != nil {
									return errLeft
								}
								//validate right operand
								errRight := expectedTokenTypes(rightOperand, TOKEN_TYPE_BOOLEAN)
								if errRight != nil {
									return errRight
								}
							} else {
								//validation for comparison operations
								//validate left operand
								errLeft := expectedTokenTypes(leftOperand, TOKEN_TYPE_BOOLEAN, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE)
								if errLeft != nil {
									return errLeft
								}
								//validate right operand
								errRight := expectedTokenTypes(rightOperand, TOKEN_TYPE_BOOLEAN, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE)
								if errRight != nil {
									return errRight
								}
							}

							result.Type = TOKEN_TYPE_BOOLEAN

							if currentToken.Type == TOKEN_TYPE_AMPERSAND {
								//LOGICAL AND operation
								leftBool := convertTokenToBool(leftOperand)
								rightBool := convertTokenToBool(rightOperand)
								if leftBool && rightBool {
									result.Value = "true"
								} else {
									result.Value = "false"
								}
							} else if currentToken.Type == TOKEN_TYPE_OR {
								//LOGICAL OR operation
								leftBool := convertTokenToBool(leftOperand)
								rightBool := convertTokenToBool(rightOperand)
								if leftBool || rightBool {
									result.Value = "true"
								} else {
									result.Value = "false"
								}
							} else if currentToken.Type == TOKEN_TYPE_EQUALITY {
								//COMPARISON EQUALITY operation
								switch leftOperand.Type {
								case TOKEN_TYPE_INTEGER:
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
										var rightOperandInt int = 0

										if rightOperand.Type == TOKEN_TYPE_INTEGER {
											rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
										} else {
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)
											rightOperandInt = int(rightOperandFloat)
										}

										if leftOperandInt == rightOperandInt {
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
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 64)
										rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)

										if leftOperandFloat == rightOperandFloat {
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
									if rightOperand.Type == TOKEN_TYPE_STRING {
										if unescapeString(leftOperand.Value) == unescapeString(rightOperand.Value) {
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
									if rightOperand.Type == TOKEN_TYPE_BOOLEAN {
										leftBool := convertTokenToBool(leftOperand)
										rightBool := convertTokenToBool(rightOperand)

										if leftBool == rightBool {
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
									if rightOperand.Type == TOKEN_TYPE_NONE {
										result.Value = "true"
									} else {
										//TOKEN_TYPE_INTEGER
										//TOKEN_TYPE_FLOAT
										//TOKEN_TYPE_STRING
										//TOKEN_TYPE_BOOLEAN
										result.Value = "false"
									}
								}
							} else if currentToken.Type == TOKEN_TYPE_INEQUALITY {
								//COMPARISON INEQUALITY operation
								switch leftOperand.Type {
								case TOKEN_TYPE_INTEGER:
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
										var rightOperandInt int = 0

										if rightOperand.Type == TOKEN_TYPE_INTEGER {
											rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
										} else {
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)
											rightOperandInt = int(rightOperandFloat)
										}

										if leftOperandInt != rightOperandInt {
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
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 64)
										rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)

										if leftOperandFloat != rightOperandFloat {
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
									if rightOperand.Type == TOKEN_TYPE_STRING {
										if leftOperand.Value != rightOperand.Value {
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
									if rightOperand.Type == TOKEN_TYPE_BOOLEAN {
										leftBool := convertTokenToBool(leftOperand)
										rightBool := convertTokenToBool(rightOperand)

										if leftBool != rightBool {
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
									if rightOperand.Type == TOKEN_TYPE_NONE {
										result.Value = "false"
									} else {
										//TOKEN_TYPE_INTEGER
										//TOKEN_TYPE_FLOAT
										//TOKEN_TYPE_STRING
										//TOKEN_TYPE_BOOLEAN
										result.Value = "true"
									}
								}
							} else if currentToken.Type == TOKEN_TYPE_GREATER_THAN {
								//COMPARISON GREATER THAN operation
								switch leftOperand.Type {
								case TOKEN_TYPE_INTEGER:
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
										var rightOperandInt int = 0

										if rightOperand.Type == TOKEN_TYPE_INTEGER {
											rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
										} else {
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)
											rightOperandInt = int(rightOperandFloat)
										}

										if leftOperandInt > rightOperandInt {
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
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 64)
										rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)

										if leftOperandFloat > rightOperandFloat {
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
							} else if currentToken.Type == TOKEN_TYPE_GREATER_THAN_OR_EQUALS {
								//COMPARISON GREATER THAN OR EQUALS operation
								switch leftOperand.Type {
								case TOKEN_TYPE_INTEGER:
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
										var rightOperandInt int = 0

										if rightOperand.Type == TOKEN_TYPE_INTEGER {
											rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
										} else {
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)
											rightOperandInt = int(rightOperandFloat)
										}

										if leftOperandInt >= rightOperandInt {
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
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 64)
										rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)

										if leftOperandFloat >= rightOperandFloat {
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
									if rightOperand.Type == TOKEN_TYPE_STRING {
										if leftOperand.Value == rightOperand.Value {
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
									if rightOperand.Type == TOKEN_TYPE_BOOLEAN {
										leftBool := convertTokenToBool(leftOperand)
										rightBool := convertTokenToBool(rightOperand)

										if leftBool == rightBool {
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
									if rightOperand.Type == TOKEN_TYPE_NONE {
										result.Value = "true"
									} else {
										//TOKEN_TYPE_INTEGER
										//TOKEN_TYPE_FLOAT
										//TOKEN_TYPE_STRING
										//TOKEN_TYPE_BOOLEAN
										result.Value = "false"
									}
								}
							} else if currentToken.Type == TOKEN_TYPE_LESS_THAN_OR_EQUALS {
								//COMPARISON LESS THAN OR EQUALS operation
								switch leftOperand.Type {
								case TOKEN_TYPE_INTEGER:
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
										var rightOperandInt int = 0

										if rightOperand.Type == TOKEN_TYPE_INTEGER {
											rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
										} else {
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)
											rightOperandInt = int(rightOperandFloat)
										}

										if leftOperandInt <= rightOperandInt {
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
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 64)
										rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)

										if leftOperandFloat <= rightOperandFloat {
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
									if rightOperand.Type == TOKEN_TYPE_STRING {
										if leftOperand.Value == rightOperand.Value {
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
									if rightOperand.Type == TOKEN_TYPE_BOOLEAN {
										leftBool := convertTokenToBool(leftOperand)
										rightBool := convertTokenToBool(rightOperand)

										if leftBool == rightBool {
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
									if rightOperand.Type == TOKEN_TYPE_NONE {
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
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandInt, _ := strconv.Atoi(leftOperand.Value)
										var rightOperandInt int = 0

										if rightOperand.Type == TOKEN_TYPE_INTEGER {
											rightOperandInt, _ = strconv.Atoi(rightOperand.Value)
										} else {
											rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)
											rightOperandInt = int(rightOperandFloat)
										}

										if leftOperandInt < rightOperandInt {
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
									if rightOperand.Type == TOKEN_TYPE_INTEGER || rightOperand.Type == TOKEN_TYPE_FLOAT {
										leftOperandFloat, _ := strconv.ParseFloat(leftOperand.Value, 64)
										rightOperandFloat, _ := strconv.ParseFloat(rightOperand.Value, 64)

										if leftOperandFloat < rightOperandFloat {
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
							lastToken = result //used in while loop to track the parameter
						} else if currentToken.Type == TOKEN_TYPE_EQUALS {
							//assignment operation
							value := PopStack(&stack)
							var errConvert error

							//if value is an identifier
							//the it's a variable
							if value.Type == TOKEN_TYPE_IDENTIFIER {
								value, errConvert = convertVariableToToken(value, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							variable := PopStack(&stack)

							//validate value
							errVal := expectedTokenTypes(value, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE, TOKEN_TYPE_BOOLEAN, TOKEN_TYPE_ARRAY, TOKEN_TYPE_ASSOCIATIVE_ARRAY)
							if errVal != nil {
								return errVal
							}

							if variable.Array_is_ref {
								//an array reference
								tmpIndex := variable.Array_ref_index
								tmpIndexStr := variable.Array_ref_index_str
								tmpIsAssoc := variable.Array_is_assoc
								tmpToken := Token{Value: variable.Array_ref_var_name, Line: variable.Line, Column: variable.Column, FileName: variable.FileName}
								variable, errConvert = convertVariableToToken(tmpToken, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
								variable.Array_is_assoc = tmpIsAssoc
								variable.Array_is_ref = true
								variable.Array_ref_index_str = tmpIndexStr
								variable.Array_ref_index = tmpIndex
							} else {
								//validate variable
								errVar := expectedTokenTypes(variable, TOKEN_TYPE_IDENTIFIER)
								if errVar != nil {
									return errVar
								}
							}

							//check if variable exists as a function
							//if yes then raise an error
							isExists, _ := isFunctionExists(variable, *globalFunctionArray)

							if isExists {
								return errors.New(SyntaxErrorMessage(variable.Line, variable.Column, "'"+variable.Value+"' exists as a function", variable.FileName))
							}

							//if not main scope then check if a system constant
							//if yes then raise an error
							if scopeName != "main" {
								if isSystemVariable(variable.Value, *globalNativeVarList) {
									return errors.New(SyntaxErrorMessage(variable.Line, variable.Column, "Cannot assign to "+variable.Value, variable.FileName))
								}
							}

							isExists, varIndex := isVariableExists(variable, *globalVariableArray, scopeName)

							if !isExists {
								//variable doesn't exists
								//create a new variable
								newVar := Variable{Name: variable.Value, ScopeName: scopeName}
								*globalVariableArray = append(*globalVariableArray, newVar)
								varIndex = len(*globalVariableArray) - 1

								//check if the first letter of variable name is in uppercase
								//if yes then tag it as constant
								firstChar := string((*globalVariableArray)[varIndex].Name[0])
								if unicode.IsUpper([]rune(firstChar)[0]) {
									(*globalVariableArray)[varIndex].IsConstant = true
								}
							} else {
								//if variable exists
								//check if constant, if yes then raise an error
								if (*globalVariableArray)[varIndex].IsConstant {
									return errors.New(SyntaxErrorMessage(variable.Line, variable.Column, "Cannot override constant '"+variable.Value+"'", variable.FileName))
								}
							}

							if variable.Array_is_ref {
								//modify the variable by its index
								if variable.Array_is_assoc {
									//associative array
									if value.Type == TOKEN_TYPE_INTEGER {
										(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].Type = VARIABLE_TYPE_INTEGER
										(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].IntegerValue, _ = strconv.Atoi(value.Value)
									} else if value.Type == TOKEN_TYPE_STRING {
										(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].Type = VARIABLE_TYPE_STRING
										(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].StringValue = value.Value
									} else if value.Type == TOKEN_TYPE_FLOAT {
										(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].Type = VARIABLE_TYPE_FLOAT
										(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].FloatValue, _ = strconv.ParseFloat(value.Value, 64)
									} else if value.Type == TOKEN_TYPE_BOOLEAN {
										(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].Type = VARIABLE_TYPE_BOOLEAN
										if value.Value == "true" {
											(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].BooleanValue = true
										} else {
											(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].BooleanValue = false
										}
									} else if value.Type == TOKEN_TYPE_ARRAY || value.Type == TOKEN_TYPE_ASSOCIATIVE_ARRAY {
										return errors.New(SyntaxErrorMessage(value.Line, value.Column, "Unexpected token '"+value.Value+"'", value.FileName))
									} else {
										//Nil
										(*globalVariableArray)[varIndex].AssociativeArrayValue[variable.Array_ref_index_str].Type = VARIABLE_TYPE_NONE
									}

								} else {
									//array
									if value.Type == TOKEN_TYPE_INTEGER {
										(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].Type = VARIABLE_TYPE_INTEGER
										(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].IntegerValue, _ = strconv.Atoi(value.Value)
									} else if value.Type == TOKEN_TYPE_STRING {
										(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].Type = VARIABLE_TYPE_STRING
										(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].StringValue = value.Value
									} else if value.Type == TOKEN_TYPE_FLOAT {
										(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].Type = VARIABLE_TYPE_FLOAT
										(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].FloatValue, _ = strconv.ParseFloat(value.Value, 64)
									} else if value.Type == TOKEN_TYPE_BOOLEAN {
										(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].Type = VARIABLE_TYPE_BOOLEAN
										if value.Value == "true" {
											(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].BooleanValue = true
										} else {
											(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].BooleanValue = false
										}
									} else if value.Type == TOKEN_TYPE_ARRAY || value.Type == TOKEN_TYPE_ASSOCIATIVE_ARRAY {
										return errors.New(SyntaxErrorMessage(value.Line, value.Column, "Unexpected token '"+value.Value+"'", value.FileName))
									} else {
										//Nil
										(*globalVariableArray)[varIndex].ArrayValue[variable.Array_ref_index].Type = VARIABLE_TYPE_NONE
									}
								}
							} else {
								//modify the value/type of variable below
								if value.Type == TOKEN_TYPE_INTEGER {
									(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_INTEGER
									(*globalVariableArray)[varIndex].IntegerValue, _ = strconv.Atoi(value.Value)
								} else if value.Type == TOKEN_TYPE_STRING {
									(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_STRING
									(*globalVariableArray)[varIndex].StringValue = value.Value
								} else if value.Type == TOKEN_TYPE_FLOAT {
									(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_FLOAT
									(*globalVariableArray)[varIndex].FloatValue, _ = strconv.ParseFloat(value.Value, 64)
								} else if value.Type == TOKEN_TYPE_BOOLEAN {
									(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_BOOLEAN
									if value.Value == "true" {
										(*globalVariableArray)[varIndex].BooleanValue = true
									} else {
										(*globalVariableArray)[varIndex].BooleanValue = false
									}
								} else if value.Type == TOKEN_TYPE_ASSOCIATIVE_ARRAY {
									(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_ASSOCIATIVE_ARRAY
									(*globalVariableArray)[varIndex].AssociativeArrayValue = make(map[string]*Variable)
									for k, v := range value.AssociativeArray {
										thisVar := new(Variable)

										if v.Type == TOKEN_TYPE_INTEGER {
											thisVar.Type = VARIABLE_TYPE_INTEGER
											thisVar.IntegerValue, _ = strconv.Atoi(v.Value)
										} else if v.Type == TOKEN_TYPE_STRING {
											thisVar.Type = VARIABLE_TYPE_STRING
											thisVar.StringValue = v.Value
										} else if v.Type == TOKEN_TYPE_FLOAT {
											thisVar.Type = VARIABLE_TYPE_FLOAT
											thisVar.FloatValue, _ = strconv.ParseFloat(v.Value, 64)
										} else if v.Type == TOKEN_TYPE_BOOLEAN {
											thisVar.Type = VARIABLE_TYPE_BOOLEAN
											if v.Value == "true" {
												thisVar.BooleanValue = true
											} else {
												thisVar.BooleanValue = false
											}
										} else {
											thisVar.Type = VARIABLE_TYPE_NONE
										}

										(*globalVariableArray)[varIndex].AssociativeArrayValue[k] = thisVar
									}
								} else if value.Type == TOKEN_TYPE_ARRAY {
									//reset slice
									(*globalVariableArray)[varIndex].ArrayValue = nil
									//set type
									(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_ARRAY
									for arrayIndex := 0; arrayIndex < len(value.Array); arrayIndex++ {
										thisVar := Variable{}

										if value.Array[arrayIndex].Type == TOKEN_TYPE_INTEGER {
											thisVar.Type = VARIABLE_TYPE_INTEGER
											thisVar.IntegerValue, _ = strconv.Atoi(value.Array[arrayIndex].Value)
										} else if value.Array[arrayIndex].Type == TOKEN_TYPE_STRING {
											thisVar.Type = VARIABLE_TYPE_STRING
											thisVar.StringValue = value.Array[arrayIndex].Value
										} else if value.Array[arrayIndex].Type == TOKEN_TYPE_FLOAT {
											thisVar.Type = VARIABLE_TYPE_FLOAT
											thisVar.FloatValue, _ = strconv.ParseFloat(value.Array[arrayIndex].Value, 64)
										} else if value.Array[arrayIndex].Type == TOKEN_TYPE_BOOLEAN {
											thisVar.Type = VARIABLE_TYPE_BOOLEAN
											if value.Array[arrayIndex].Value == "true" {
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
							}

							stack = append(stack, value)
						} else if currentToken.Type == TOKEN_TYPE_FUNCTION {
							//function execution here
							var functionArguments []FunctionArgument

							//check if function is existing below
							isExists, funcIndex := isFunctionExists(currentToken, *globalFunctionArray)
							if !isExists {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Function '"+currentToken.Value+"' doesn't exists", currentToken.FileName))
							}

							//check if function got arguments
							if (*globalFunctionArray)[funcIndex].ArgumentCount > 0 {
								//function parameter validation below

								arg_count := 0

								//add to counter if column is greater than current column
								for sx := 0; sx < len(stack); sx++ {
									if stack[sx].Column > currentToken.Column {
										arg_count += 1
									}
								}

								if (*globalFunctionArray)[funcIndex].ArgumentCount != arg_count {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, currentToken.Value+" takes exactly "+strconv.Itoa((*globalFunctionArray)[funcIndex].ArgumentCount)+" argument", currentToken.FileName))
								}

								//add arguments from stack below
								processedArg := 0
								for true {
									var param Token
									//add to functionargument one by one
									param = PopStack(&stack)

									var errConvert error
									if param.Type == TOKEN_TYPE_IDENTIFIER {
										param, errConvert = convertVariableToToken(param, *globalVariableArray, scopeName)
										if errConvert != nil {
											return errConvert
										}
									}

									fa := FunctionArgument{}
									//convert token to param (TODO: create a function for this one?)
									if param.Type == TOKEN_TYPE_INTEGER {
										fa.Type = ARG_TYPE_INTEGER
										fa.IntegerValue, _ = strconv.Atoi(param.Value)
									} else if param.Type == TOKEN_TYPE_STRING {
										fa.Type = ARG_TYPE_STRING
										fa.StringValue = param.Value
									} else if param.Type == TOKEN_TYPE_FLOAT {
										fa.Type = ARG_TYPE_FLOAT
										fa.FloatValue, _ = strconv.ParseFloat(param.Value, 64)
									} else if param.Type == TOKEN_TYPE_BOOLEAN {
										fa.Type = ARG_TYPE_BOOLEAN
										if param.Value == "true" {
											fa.BooleanValue = true
										} else {
											fa.BooleanValue = false
										}
									} else if param.Type == TOKEN_TYPE_ASSOCIATIVE_ARRAY {
										fa.Type = ARG_TYPE_ASSOCIATIVE_ARRAY
										fa.AssociativeArrayValue = make(map[string]FunctionArgument)
										for k, v := range param.AssociativeArray {
											thisArgument := FunctionArgument{}
											if v.Type == TOKEN_TYPE_INTEGER {
												thisArgument.Type = ARG_TYPE_INTEGER
												thisArgument.IntegerValue, _ = strconv.Atoi(v.Value)
											} else if v.Type == TOKEN_TYPE_STRING {
												thisArgument.Type = ARG_TYPE_STRING
												thisArgument.StringValue = v.Value
											} else if v.Type == TOKEN_TYPE_FLOAT {
												thisArgument.Type = ARG_TYPE_FLOAT
												thisArgument.FloatValue, _ = strconv.ParseFloat(v.Value, 64)
											} else if v.Type == TOKEN_TYPE_BOOLEAN {
												thisArgument.Type = ARG_TYPE_BOOLEAN
												if v.Value == "true" {
													thisArgument.BooleanValue = true
												} else {
													thisArgument.BooleanValue = false
												}
											} else {
												thisArgument.Type = ARG_TYPE_NONE
											}

											fa.AssociativeArrayValue[k] = thisArgument
										}
									} else if param.Type == TOKEN_TYPE_ARRAY {
										fa.Type = ARG_TYPE_ARRAY
										for thisArrayIndex := 0; thisArrayIndex < len(param.Array); thisArrayIndex++ {
											thisArgument := FunctionArgument{}
											if param.Array[thisArrayIndex].Type == TOKEN_TYPE_INTEGER {
												thisArgument.Type = ARG_TYPE_INTEGER
												thisArgument.IntegerValue, _ = strconv.Atoi(param.Array[thisArrayIndex].Value)
											} else if param.Array[thisArrayIndex].Type == TOKEN_TYPE_STRING {
												thisArgument.Type = ARG_TYPE_STRING
												thisArgument.StringValue = param.Array[thisArrayIndex].Value
											} else if param.Array[thisArrayIndex].Type == TOKEN_TYPE_FLOAT {
												thisArgument.Type = ARG_TYPE_FLOAT
												thisArgument.FloatValue, _ = strconv.ParseFloat(param.Array[thisArrayIndex].Value, 64)
											} else if param.Array[thisArrayIndex].Type == TOKEN_TYPE_BOOLEAN {
												thisArgument.Type = ARG_TYPE_BOOLEAN
												if param.Array[thisArrayIndex].Value == "true" {
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
									if processedArg == (*globalFunctionArray)[funcIndex].ArgumentCount {
										break
									}
								}
							} else {
								// no argument

								arg_count := 0

								//add to counter if column is greater than current column
								for sx := 0; sx < len(stack); sx++ {
									if stack[sx].Column > currentToken.Column {
										arg_count += 1
									}
								}

								if arg_count > 0 {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, currentToken.Value+" takes exactly 0 argument", currentToken.FileName))
								}
							}

							newToken := currentToken
							if (*globalFunctionArray)[funcIndex].IsNative {
								//execute native function
								var thisError error
								funcReturn := (*globalFunctionArray)[funcIndex].Run(functionArguments, &thisError, globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, globalSettings, currentToken.Line, currentToken.Column, currentToken.FileName)
								if thisError != nil {
									return thisError
								}
								//convert FunctionReturn to Token and append to stack (TODO: Create a function for conversion?)
								if funcReturn.Type == RET_TYPE_INTEGER {
									newToken.Type = TOKEN_TYPE_INTEGER
									newToken.Value = strconv.Itoa(funcReturn.IntegerValue)
								} else if funcReturn.Type == RET_TYPE_STRING {
									newToken.Type = TOKEN_TYPE_STRING
									newToken.Value = funcReturn.StringValue
								} else if funcReturn.Type == RET_TYPE_FLOAT {
									newToken.Type = TOKEN_TYPE_FLOAT
									newToken.Value = strconv.FormatFloat(funcReturn.FloatValue, 'f', -1, 64)
								} else if funcReturn.Type == RET_TYPE_BOOLEAN {
									newToken.Type = TOKEN_TYPE_BOOLEAN
									if funcReturn.BooleanValue {
										//true
										newToken.Value = "true"
									} else {
										//false
										newToken.Value = "false"
									}
								} else if funcReturn.Type == RET_TYPE_ASSOCIATIVE_ARRAY {
									newToken.Type = TOKEN_TYPE_ASSOCIATIVE_ARRAY
									newToken.AssociativeArray = make(map[string]Token)
									newToken.OtherInt = len(funcReturn.AssociativeArrayValue)

									for k, v := range funcReturn.AssociativeArrayValue {
										thisNewToken := Token{}

										if v.Type == RET_TYPE_INTEGER {
											thisNewToken.Type = TOKEN_TYPE_INTEGER
											thisNewToken.Value = strconv.Itoa(v.IntegerValue)
										} else if v.Type == RET_TYPE_STRING {
											thisNewToken.Type = TOKEN_TYPE_STRING
											thisNewToken.Value = v.StringValue
										} else if v.Type == RET_TYPE_FLOAT {
											thisNewToken.Type = TOKEN_TYPE_FLOAT
											thisNewToken.Value = strconv.FormatFloat(v.FloatValue, 'f', -1, 64)
										} else if v.Type == RET_TYPE_BOOLEAN {
											thisNewToken.Type = TOKEN_TYPE_BOOLEAN
											if v.BooleanValue {
												thisNewToken.Value = "true"
											} else {
												thisNewToken.Value = "false"
											}
										} else {
											thisNewToken.Type = TOKEN_TYPE_NONE
										}
										newToken.AssociativeArray[k] = thisNewToken
									}
								} else if funcReturn.Type == RET_TYPE_ARRAY {
									newToken.Type = TOKEN_TYPE_ARRAY
									newToken.OtherInt = len(funcReturn.ArrayValue)

									for thisArrayIndex := 0; thisArrayIndex < len(funcReturn.ArrayValue); thisArrayIndex++ {
										thisNewToken := Token{}

										if funcReturn.ArrayValue[thisArrayIndex].Type == RET_TYPE_INTEGER {
											thisNewToken.Type = TOKEN_TYPE_INTEGER
											thisNewToken.Value = strconv.Itoa(funcReturn.ArrayValue[thisArrayIndex].IntegerValue)
										} else if funcReturn.ArrayValue[thisArrayIndex].Type == RET_TYPE_STRING {
											thisNewToken.Type = TOKEN_TYPE_STRING
											thisNewToken.Value = funcReturn.ArrayValue[thisArrayIndex].StringValue
										} else if funcReturn.ArrayValue[thisArrayIndex].Type == RET_TYPE_FLOAT {
											thisNewToken.Type = TOKEN_TYPE_FLOAT
											thisNewToken.Value = strconv.FormatFloat(funcReturn.ArrayValue[thisArrayIndex].FloatValue, 'f', -1, 64)
										} else if funcReturn.ArrayValue[thisArrayIndex].Type == RET_TYPE_BOOLEAN {
											thisNewToken.Type = TOKEN_TYPE_BOOLEAN
											if funcReturn.ArrayValue[thisArrayIndex].BooleanValue {
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

									if functionArguments[ind].Type == ARG_TYPE_INTEGER {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_INTEGER
										(*globalVariableArray)[varIndex].IntegerValue = functionArguments[ind].IntegerValue
									} else if functionArguments[ind].Type == ARG_TYPE_STRING {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_STRING
										(*globalVariableArray)[varIndex].StringValue = functionArguments[ind].StringValue
									} else if functionArguments[ind].Type == ARG_TYPE_FLOAT {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_FLOAT
										(*globalVariableArray)[varIndex].FloatValue = functionArguments[ind].FloatValue
									} else if functionArguments[ind].Type == ARG_TYPE_BOOLEAN {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_BOOLEAN
										(*globalVariableArray)[varIndex].BooleanValue = functionArguments[ind].BooleanValue
									} else if functionArguments[ind].Type == ARG_TYPE_ASSOCIATIVE_ARRAY {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_ASSOCIATIVE_ARRAY
										(*globalVariableArray)[varIndex].AssociativeArrayValue = make(map[string]*Variable)

										for k, v := range functionArguments[ind].AssociativeArrayValue {
											thisVar := new(Variable)

											if v.Type == ARG_TYPE_INTEGER {
												thisVar.Type = VARIABLE_TYPE_INTEGER
												thisVar.IntegerValue = v.IntegerValue
											} else if v.Type == ARG_TYPE_STRING {
												thisVar.Type = VARIABLE_TYPE_STRING
												thisVar.StringValue = v.StringValue
											} else if v.Type == ARG_TYPE_FLOAT {
												thisVar.Type = VARIABLE_TYPE_FLOAT
												thisVar.FloatValue = v.FloatValue
											} else if v.Type == ARG_TYPE_BOOLEAN {
												thisVar.Type = VARIABLE_TYPE_BOOLEAN
												thisVar.BooleanValue = v.BooleanValue
											} else {
												thisVar.Type = VARIABLE_TYPE_NONE
											}
											(*globalVariableArray)[varIndex].AssociativeArrayValue[k] = thisVar
										}

									} else if functionArguments[ind].Type == ARG_TYPE_ARRAY {
										(*globalVariableArray)[varIndex].Type = VARIABLE_TYPE_ARRAY
										for thisArrayIndex := 0; thisArrayIndex < len(functionArguments[ind].ArrayValue); thisArrayIndex++ {
											thisVar := Variable{}

											if functionArguments[ind].ArrayValue[thisArrayIndex].Type == ARG_TYPE_INTEGER {
												thisVar.Type = VARIABLE_TYPE_INTEGER
												thisVar.IntegerValue = functionArguments[ind].ArrayValue[thisArrayIndex].IntegerValue
											} else if functionArguments[ind].ArrayValue[thisArrayIndex].Type == ARG_TYPE_STRING {
												thisVar.Type = VARIABLE_TYPE_STRING
												thisVar.StringValue = functionArguments[ind].ArrayValue[thisArrayIndex].StringValue
											} else if functionArguments[ind].ArrayValue[thisArrayIndex].Type == ARG_TYPE_FLOAT {
												thisVar.Type = VARIABLE_TYPE_FLOAT
												thisVar.FloatValue = functionArguments[ind].ArrayValue[thisArrayIndex].FloatValue
											} else if functionArguments[ind].ArrayValue[thisArrayIndex].Type == ARG_TYPE_BOOLEAN {
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
								var thisLastStackBool bool = false

								//execute user defined function
								prsr := Parser{}
								parserErr := prsr.Parse((*globalFunctionArray)[funcIndex].Tokens, globalVariableArray, globalFunctionArray, thisScopeName, globalNativeVarList, &thisGotReturn, &thisReturnToken, false, &thisNeedBreak, &thisStackReference, globalSettings, getLastStackBool, &thisLastStackBool)

								if parserErr != nil {
									return parserErr
								}

								if thisGotReturn {
									//the function returns a value
									newToken = thisReturnToken
								}

								//DELETE GENERATED VARIABLES with thisScopeName
								cleanupVariables(globalVariableArray, thisScopeName)
							}

							stack = append(stack, newToken)

						} else if currentToken.Type == TOKEN_TYPE_OPEN_BRACES || currentToken.Type == TOKEN_TYPE_OPEN_BRACKET {

							if currentToken.From_function_call {
								value := PopStack(&stack)

								if value.Type != TOKEN_TYPE_INTEGER {
									return errors.New(SyntaxErrorMessage(value.Line, value.Column, "Unexpected token '"+value.Value+"'", value.FileName))
								}

								// get the real value from array
								this_index, _ := strconv.Atoi(value.Value)
								old_value := value

								value = PopStack(&stack)

								if (this_index + 1) > len(value.Array) {
									return errors.New(SyntaxErrorMessage(old_value.Line, old_value.Column, "Index out of range", old_value.FileName))
								}
								stack = append(stack, value.Array[this_index])
							} else {
								//array declaration
								processedArg := 0
								newToken := currentToken
								var tempArray []Token //for TOKEN-TYPE_ARRAY

								if currentToken.Type == TOKEN_TYPE_OPEN_BRACES {
									newToken.Type = TOKEN_TYPE_ARRAY
								} else {
									//TOKEN_TYPE_OPEN_BRACKET , associative array
									newToken.Type = TOKEN_TYPE_ASSOCIATIVE_ARRAY
									newToken.AssociativeArray = make(map[string]Token)
								}

								if currentToken.OtherInt > 0 {
									if len(stack) == 0 {
										if len(outputQueue) > 0 {
											return errors.New(SyntaxErrorMessage(outputQueue[0].Line, outputQueue[0].Column, "Unexpected token '"+outputQueue[0].Value+"'", outputQueue[0].FileName))
										} else {
											return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Unexpected token '"+currentToken.Value+"'", currentToken.FileName))
										}
									}

									for true {
										var param Token

										if len(stack) == 0 {
											return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Unexpected token '"+currentToken.Value+"'", currentToken.FileName))
										}

										param = PopStack(&stack)

										if currentToken.Type == TOKEN_TYPE_OPEN_BRACES {
											var errConvert error
											if param.Type == TOKEN_TYPE_IDENTIFIER {
												param, errConvert = convertVariableToToken(param, *globalVariableArray, scopeName)
												if errConvert != nil {
													return errConvert
												}
											}

											if param.Type == TOKEN_TYPE_ARRAY {
												return errors.New(SyntaxErrorMessage(param.Line, param.Column, "Unexpected token '"+param.Value+"'", param.FileName))
											}

											tempArray = append(tempArray, param)
										} else {
											//associative array
											//param should be key-value pair
											errParam := expectedTokenTypes(param, TOKEN_TYPE_KEY_VALUE_PAIR)
											if errParam != nil {
												return errParam
											}

											for k, v := range param.AssociativeArray {
												newToken.AssociativeArray[k] = v
											}
										}

										processedArg += 1
										if processedArg == currentToken.OtherInt {
											break
										}
									}

									if currentToken.Type == TOKEN_TYPE_OPEN_BRACES {
										if len(tempArray) > 0 {
											//reverse the array (so it's in proper position)
											arrayLength := len(tempArray)
											for thisArrayIndex := 0; thisArrayIndex < arrayLength; thisArrayIndex++ {
												thisToken := tempArray[len(tempArray)-1]
												tempArray = tempArray[:len(tempArray)-1]
												newToken.Array = append(newToken.Array, thisToken)
											}
										}
									} else {
										newToken.OtherInt = len(newToken.AssociativeArray)
									}
								}
								stack = append(stack, newToken)
								//DumpToken(stack)
							}

						} else if currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_START {
							//check if function already exists
							//if yes then raise an error
							isExists, _ := isFunctionExists(currentToken, *globalFunctionArray)
							if isExists {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Function '"+currentToken.Value+"' already exists", currentToken.FileName))
							}

							//get function arguments
							var functionParams []Token
							for true {
								var param Token

								if len(stack) > 0 {
									param = PopStack(&stack)

									//validate param type
									errParam := expectedTokenTypes(param, TOKEN_TYPE_IDENTIFIER)
									if errParam != nil {
										return errParam
									}

									//check if param is a constant
									firstChar := string(param.Value[0])
									if unicode.IsUpper([]rune(firstChar)[0]) {
										return errors.New(SyntaxErrorMessage(param.Line, param.Column, "Argument cannot be a constant", param.FileName))
									}

									//check if param already exists
									if isParamExists(param, functionParams) {
										return errors.New(SyntaxErrorMessage(param.Line, param.Column, "Duplicate argument '"+param.Value+"' in function definition", param.FileName))
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

								if currentToken.Type == TOKEN_TYPE_FUNCTION_DEF_END {
									break
								}

								if len(outputQueue) == 0 {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Invalid statement", currentToken.FileName))
								}
							}

							//append to global functions
							*globalFunctionArray = append(*globalFunctionArray, newFunction)
							stack = append(stack, currentToken) //TODO: not sure if it should append the TOKEN_TYPE_FUNCTION_DEF_END

						} else if currentToken.Type == TOKEN_TYPE_FUNCTION_RETURN {
							if scopeName == "main" {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "'rtn' outside function", currentToken.FileName))
							}

							returnValue := PopStack(&stack)
							var errConvert error

							if returnValue.Type == TOKEN_TYPE_IDENTIFIER {
								returnValue, errConvert = convertVariableToToken(returnValue, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							errValue := expectedTokenTypes(returnValue, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE, TOKEN_TYPE_BOOLEAN, TOKEN_TYPE_ARRAY, TOKEN_TYPE_ASSOCIATIVE_ARRAY)
							if errValue != nil {
								return errValue
							}

							*gotReturn = true
							*returnToken = returnValue
							return nil
						} else if currentToken.Type == TOKEN_TYPE_FOR_LOOP_START {
							//get looping parameters
							//param validation
							if len(stack) == 0 || len(stack) < 2 {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, currentToken.Value+" takes exactly 2 argument", currentToken.FileName))
							}
							var errConvert error

							param2 := PopStack(&stack)

							//if param2 is an identifier
							//then it's a variable
							if param2.Type == TOKEN_TYPE_IDENTIFIER {
								param2, errConvert = convertVariableToToken(param2, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							param1 := PopStack(&stack)

							//if param1 is an identifier
							//then it's a variable
							if param1.Type == TOKEN_TYPE_IDENTIFIER {
								param1, errConvert = convertVariableToToken(param1, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							//validate param2
							errParam2 := expectedTokenTypes(param2, TOKEN_TYPE_INTEGER)
							if errParam2 != nil {
								return errParam2
							}
							//validate param1
							errParam1 := expectedTokenTypes(param1, TOKEN_TYPE_INTEGER)
							if errParam1 != nil {
								return errParam1
							}

							var tempTokens []Token
							openLoopCount = 0
							//append all tokens to temporary token
							for true {
								currentToken := outputQueue[0]
								outputQueue = append(outputQueue[:0], outputQueue[1:]...)

								tempTokens = append(tempTokens, currentToken)

								if currentToken.Type == TOKEN_TYPE_FOR_LOOP_START {
									openLoopCount += 1
								}

								if currentToken.Type == TOKEN_TYPE_FOR_LOOP_END {
									if openLoopCount == 0 {
										break
									}
									openLoopCount = openLoopCount - 1
								}
								if len(outputQueue) == 0 {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Invalid statement", currentToken.FileName))
								}
							}

							iParam1, _ := strconv.Atoi(param1.Value)
							iParam2, _ := strconv.Atoi(param2.Value)

							if iParam2 < iParam1 {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "First parameter should be lower than or equals to second parameter", currentToken.FileName))
							}

							var isInfiniteLoop = false
							if iParam2 == iParam1 {
								isInfiniteLoop = true
							}

							if isInfiniteLoop {
								for {

									var loopGotReturn bool = false
									var loopReturnToken Token
									var loopNeedBreak bool = false
									var loopStackReference []Token
									var thisLastStackBool bool = false

									prsr := Parser{}
									parserErr := prsr.Parse(tempTokens, globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, &loopGotReturn, &loopReturnToken, true, &loopNeedBreak, &loopStackReference, globalSettings, getLastStackBool, &thisLastStackBool)

									if parserErr != nil {
										return parserErr
									}
									if loopGotReturn {
										*gotReturn = loopGotReturn
										*returnToken = loopReturnToken
										return nil
									}

									if loopNeedBreak {
										//break the loop
										break
									}
								}
							} else {
								iParam2 -= 1
								for loop_index := iParam1; loop_index <= iParam2; loop_index++ {

									var loopGotReturn bool = false
									var loopReturnToken Token
									var loopNeedBreak bool = false
									var loopStackReference []Token
									var thisLastStackBool bool = false

									prsr := Parser{}
									parserErr := prsr.Parse(tempTokens, globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, &loopGotReturn, &loopReturnToken, true, &loopNeedBreak, &loopStackReference, globalSettings, getLastStackBool, &thisLastStackBool)

									if parserErr != nil {
										return parserErr
									}
									if loopGotReturn {
										*gotReturn = loopGotReturn
										*returnToken = loopReturnToken
										return nil
									}

									if loopNeedBreak {
										//break the loop
										break
									}
								}
							}

							stack = append(stack, currentToken) //append TOKEN_TYPE_FOR_LOOP_END

						} else if currentToken.Type == TOKEN_TYPE_WHILE_LOOP_START {

							if len(stack) == 0 || len(stack) > 1 {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, currentToken.Value+" takes exactly 1 argument", currentToken.FileName))
							}
							var errConvert error

							param1 := PopStack(&stack)

							//if param1 is an identifier
							//then it's a variable
							if param1.Type == TOKEN_TYPE_IDENTIFIER {
								param1, errConvert = convertVariableToToken(param1, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							//validate param1
							errParam1 := expectedTokenTypes(param1, TOKEN_TYPE_BOOLEAN)
							if errParam1 != nil {
								return errParam1
							}

							var tempTokens []Token
							openWhileLoopCount = 0
							//append all tokens to temporary token
							for true {
								currentToken := outputQueue[0]
								outputQueue = append(outputQueue[:0], outputQueue[1:]...)

								tempTokens = append(tempTokens, currentToken)

								if currentToken.Type == TOKEN_TYPE_WHILE_LOOP_START {
									openWhileLoopCount += 1
								}

								if currentToken.Type == TOKEN_TYPE_WHILE_LOOP_END {
									if openWhileLoopCount == 0 {
										break
									}
									openWhileLoopCount = openWhileLoopCount - 1
								}

								if len(outputQueue) == 0 {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Invalid statement", currentToken.FileName))
								}
							}

							whileCondition := convertTokenToBool(param1)

							for whileCondition {
								var loopGotReturn bool = false
								var loopReturnToken Token
								var loopNeedBreak bool = false
								var loopStackReference []Token
								var thisLastStackBool bool = false

								prsr := Parser{}
								// parse the body of the while loop
								parserErr := prsr.Parse(tempTokens, globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, &loopGotReturn, &loopReturnToken, true, &loopNeedBreak, &loopStackReference, globalSettings, getLastStackBool, &thisLastStackBool)

								if parserErr != nil {
									return parserErr
								}
								if loopGotReturn {
									*gotReturn = loopGotReturn
									*returnToken = loopReturnToken
									return nil
								}

								if loopNeedBreak {
									//break the loop
									break
								}

								//change the context of the token temporarily

								for w_i := 0; w_i < len(whileLoopTokens[currentToken.Context]); w_i++ {
									whileLoopTokens[currentToken.Context][w_i].Context = currentContext
								}

								// append newline so it can be parsed

								new_token := Token{
									Type:    TOKEN_TYPE_NEWLINE,
									Context: currentContext,
								}

								whileLoopTokens[currentToken.Context] = append(whileLoopTokens[currentToken.Context], new_token)

								//parse the while condition again
								prsr.Parse(whileLoopTokens[currentToken.Context], globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, &loopGotReturn, &loopReturnToken, true, &loopNeedBreak, &loopStackReference, globalSettings, true, &thisLastStackBool)

								whileCondition = thisLastStackBool
							}

							stack = append(stack, currentToken) //append TOKEN_TYPE_WHILE_LOOP_END
						} else if currentToken.Type == TOKEN_TYPE_LOOP_BREAK {
							if !isLoop {
								return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "'brk' outside loop", currentToken.FileName))
							}
							*needBreak = true
							return nil
						} else if currentToken.Type == TOKEN_TYPE_IF_START {
							var params []Token
							var gotElse bool = false

							//param := stack[len(stack)-1]

							params = append(params, PopStack(&stack))

							var tempTokens []TokenArray
							openIfCount = 0
							//append all tokens to temporary token

							tempTokens = append(tempTokens, TokenArray{})
							for true {
								currentToken := outputQueue[0]
								outputQueue = append(outputQueue[:0], outputQueue[1:]...)

								tempTokens[len(tempTokens)-1].Tokens = append(tempTokens[len(tempTokens)-1].Tokens, currentToken)

								if currentToken.Type == TOKEN_TYPE_IF_START {
									openIfCount += 1
								}

								if currentToken.Type == TOKEN_TYPE_IF_END || currentToken.Type == TOKEN_TYPE_ELSE || currentToken.Type == TOKEN_TYPE_ELIF_START {
									if openIfCount == 0 {
										if currentToken.Type == TOKEN_TYPE_IF_END {
											break
										} else if currentToken.Type == TOKEN_TYPE_ELIF_START {
											//else if
											//get parameters until end
											var paramTokens []Token
											var currentToken2 Token
											var will_change_to_main string = currentToken.Context
											for true {
												currentToken2 = outputQueue[0]
												outputQueue = append(outputQueue[:0], outputQueue[1:]...)

												if currentToken2.Type == TOKEN_TYPE_ELIF_PARAM_END {
													break
												} else {
													if currentToken2.Context == will_change_to_main {
														currentToken2.Context = CONTEXT_NAME_MAIN
													}
													paramTokens = append(paramTokens, currentToken2)
												}

												if len(outputQueue) == 0 {
													return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Invalid statement", currentToken.FileName))
												}
											}
											//evaluate paramTokens
											var efGotReturn bool = false
											var efReturnToken Token
											var efNeedBreak bool = false
											var efStackReference []Token
											var thisLastStackBool bool = false

											//add newline to paramTokens so parser can execute it
											paramTokens = append(paramTokens, Token{Value: "\n", FileName: currentToken2.FileName, Type: TOKEN_TYPE_NEWLINE, Line: currentToken2.Line, Column: currentToken2.Column, Context: CONTEXT_NAME_MAIN})

											prsr := Parser{}
											parserErr := prsr.Parse(paramTokens, globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, &efGotReturn, &efReturnToken, isLoop, &efNeedBreak, &efStackReference, globalSettings, getLastStackBool, &thisLastStackBool)

											if parserErr != nil {
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
									if currentToken.Type == TOKEN_TYPE_IF_END {
										openIfCount = openIfCount - 1
									}
								}

								if len(outputQueue) == 0 {
									return errors.New(SyntaxErrorMessage(currentToken.Line, currentToken.Column, "Invalid statement", currentToken.FileName))
								}
							}

							var currentStatementIndex int = 0
							var executeIfStatement bool = false

							for paramIndex := 0; paramIndex < len(params); paramIndex++ {
								var thisParam Token = params[paramIndex]
								var errConvert error

								if thisParam.Type == TOKEN_TYPE_IDENTIFIER {
									thisParam, errConvert = convertVariableToToken(thisParam, *globalVariableArray, scopeName)
									if errConvert != nil {
										return errConvert
									}
								}

								//validate param
								errParam := expectedTokenTypes(thisParam, TOKEN_TYPE_BOOLEAN)
								if errParam != nil {
									return errParam
								}

								paramBool := convertTokenToBool(thisParam)

								if paramBool {
									executeIfStatement = true
									currentStatementIndex = paramIndex
									break
								}
							}

							if !executeIfStatement {
								if gotElse && len(tempTokens) > 1 {
									executeIfStatement = true
									currentStatementIndex = len(tempTokens) - 1
								}
							}

							if executeIfStatement {
								var ifGotReturn bool = false
								var ifReturnToken Token
								var ifNeedBreak bool = false
								var ifStackReference []Token
								var thisLastStackBool bool = false

								prsr := Parser{}
								parserErr := prsr.Parse(tempTokens[currentStatementIndex].Tokens, globalVariableArray, globalFunctionArray, scopeName, globalNativeVarList, &ifGotReturn, &ifReturnToken, isLoop, &ifNeedBreak, &ifStackReference, globalSettings, getLastStackBool, &thisLastStackBool)

								if parserErr != nil {
									return parserErr
								}
								if ifGotReturn {
									*gotReturn = ifGotReturn
									*returnToken = ifReturnToken
									return nil
								}

								if isLoop && ifNeedBreak {
									*needBreak = ifNeedBreak
									return nil
								}
							}

							stack = append(stack, currentToken) //TOKEN_TYPE_IF_END
						} else if currentToken.Type == TOKEN_TYPE_GET_ARRAY_START {
							//array getter
							var errConvert error
							var thisArrayToken Token

							param := PopStack(&stack)

							if param.Type == TOKEN_TYPE_IDENTIFIER {
								param, errConvert = convertVariableToToken(param, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							//validate param
							errParam := expectedTokenTypes(param, TOKEN_TYPE_INTEGER, TOKEN_TYPE_STRING)
							if errParam != nil {
								return errParam
							}

							//validate array
							thisArrayToken, errConvert = convertVariableToToken(currentToken, *globalVariableArray, scopeName)
							if errConvert != nil {
								return errConvert
							}

							if param.Type == TOKEN_TYPE_INTEGER {
								//array
								if thisArrayToken.Type != TOKEN_TYPE_ARRAY {
									return errors.New(SyntaxErrorMessage(thisArrayToken.Line, thisArrayToken.Column, "Not a lineup type", thisArrayToken.FileName))
								}

								intIndex, _ := strconv.Atoi(param.Value)

								if intIndex < 0 || len(thisArrayToken.Array) == 0 || intIndex > (len(thisArrayToken.Array)-1) {
									return errors.New(SyntaxErrorMessage(param.Line, param.Column, "Index out of range", param.FileName))
								}

								//add reference for assigment operation
								thisArrayToken.Array[intIndex].Array_is_ref = true
								thisArrayToken.Array[intIndex].Array_is_assoc = false
								thisArrayToken.Array[intIndex].Array_ref_index = intIndex
								thisArrayToken.Array[intIndex].Array_ref_var_name = thisArrayToken.Value
								thisArrayToken.Array[intIndex].Line = currentToken.Line
								thisArrayToken.Array[intIndex].Column = currentToken.Column
								thisArrayToken.Array[intIndex].FileName = currentToken.FileName

								stack = append(stack, thisArrayToken.Array[intIndex])
							}

							if param.Type == TOKEN_TYPE_STRING {
								//associative array
								if thisArrayToken.Type != TOKEN_TYPE_ASSOCIATIVE_ARRAY {
									return errors.New(SyntaxErrorMessage(thisArrayToken.Line, thisArrayToken.Column, "Not a glossary type", thisArrayToken.FileName))
								}

								tokenItem, ok := thisArrayToken.AssociativeArray[param.Value]
								if !ok {
									return errors.New(SyntaxErrorMessage(param.Line, param.Column, "Index out of range", param.FileName))
								}

								//add reference for assigment operation
								tokenItem.Array_is_ref = true
								tokenItem.Array_is_assoc = true
								tokenItem.Array_ref_index_str = param.Value
								tokenItem.Array_ref_var_name = thisArrayToken.Value
								tokenItem.Line = currentToken.Line
								tokenItem.Column = currentToken.Column
								tokenItem.FileName = currentToken.FileName

								stack = append(stack, tokenItem)
							}

						} else if currentToken.Type == TOKEN_TYPE_COLON {
							//key-value pair
							value := PopStack(&stack)
							var errConvert error

							if value.Type == TOKEN_TYPE_IDENTIFIER {
								value, errConvert = convertVariableToToken(value, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							errVal := expectedTokenTypes(value, TOKEN_TYPE_INTEGER, TOKEN_TYPE_FLOAT, TOKEN_TYPE_STRING, TOKEN_TYPE_NONE, TOKEN_TYPE_BOOLEAN) //, TOKEN_TYPE_ARRAY (TODO: ADD THE ARRAY LATER)
							if errVal != nil {
								return errVal
							}

							key := PopStack(&stack)

							if key.Type == TOKEN_TYPE_IDENTIFIER {
								key, errConvert = convertVariableToToken(key, *globalVariableArray, scopeName)
								if errConvert != nil {
									return errConvert
								}
							}

							errVal = expectedTokenTypes(key, TOKEN_TYPE_STRING)
							if errVal != nil {
								return errVal
							}

							newToken := currentToken
							newToken.Type = TOKEN_TYPE_KEY_VALUE_PAIR
							newToken.AssociativeArray = make(map[string]Token)
							newToken.AssociativeArray[key.Value] = value

							stack = append(stack, newToken)

						} else {
							stack = append(stack, currentToken)
							lastToken = currentToken //used in while loop to track the parameter
						}

					}

					*stackReference = stack

					//REMOVE THE VALIDATION ERROR NOTE: not sure if this one should be removed!!!
					/*
						if(len(stack) > 1) {
							return errors.New(SyntaxErrorMessage(stack[0].Line, stack[0].Column, "Invalid statements", stack[0].FileName))
						}
					*/

					/* else {
						if(stack[0].Type == TOKEN_TYPE_IDENTIFIER) {
							return errors.New(SyntaxErrorMessage(stack[0].Line, stack[0].Column, "Unexpected token '" + stack[0].Value + "'", stack[0].FileName))
						}
					}
					*/
					if len(stack) > 0 {
						if stack[0].Type == TOKEN_TYPE_IDENTIFIER {
							//check if existing as function
							isExists, _ := isFunctionExists(stack[0], *globalFunctionArray)
							if isExists {
								return errors.New(SyntaxErrorMessage(stack[0].Line, stack[0].Column, "Unexpected token '"+stack[0].Value+"'", stack[0].FileName))
							} else {
								//variable
								isExists, _ := isVariableExists(stack[0], *globalVariableArray, scopeName)
								if !isExists {
									return errors.New(SyntaxErrorMessage(stack[0].Line, stack[0].Column, "Variable doesn't exists '"+stack[0].Value+"'", stack[0].FileName))
								}
							}
						}
					}
					//fmt.Println("test")
					//if getLastStackBool {
					//	*lastStackBool = convertTokenToBool(stack[0])
					//}
				}

			}

		} else if tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_START || tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_END || tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_START || tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_END || tokenArray[x].Type == TOKEN_TYPE_IF_START || tokenArray[x].Type == TOKEN_TYPE_IF_END || tokenArray[x].Type == TOKEN_TYPE_WHILE_LOOP_START || tokenArray[x].Type == TOKEN_TYPE_WHILE_LOOP_END {
			if tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_START {
				ignoreNewline = true
				isFunctionDefinition = true
			} else if tokenArray[x].Type == TOKEN_TYPE_FUNCTION_DEF_END {
				//TOKEN_TYPE_FUNCTION_DEF_END
				ignoreNewline = false
				isFunctionDefinition = false

				//NOTE: not sure if the code below is temporary
				//append newline (to make the one liner definition of function works)
				tokensToEvaluate = append(tokensToEvaluate, Token{Value: "\n", FileName: tokenArray[x].FileName, Type: TOKEN_TYPE_NEWLINE, Line: tokenArray[x].Line, Column: tokenArray[x].Column})
			}
			if !isFunctionDefinition {
				if !isIfStatement {
					if tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_START {
						openLoopCount += 1
						ignoreNewline = true
						isLoopStatement = true
					} else if tokenArray[x].Type == TOKEN_TYPE_FOR_LOOP_END {
						openLoopCount = openLoopCount - 1

						if openLoopCount == 0 {
							ignoreNewline = false
							isLoopStatement = false
							tokensToEvaluate = append(tokensToEvaluate, Token{Value: "\n", FileName: tokenArray[x].FileName, Type: TOKEN_TYPE_NEWLINE, Line: tokenArray[x].Line, Column: tokenArray[x].Column})
						}
					}
					if tokenArray[x].Type == TOKEN_TYPE_WHILE_LOOP_START {
						openWhileLoopCount += 1
						ignoreNewline = true
						isWhileLoopStatement = true
					} else if tokenArray[x].Type == TOKEN_TYPE_WHILE_LOOP_END {
						openWhileLoopCount = openWhileLoopCount - 1

						if openWhileLoopCount == 0 {
							ignoreNewline = false
							isWhileLoopStatement = false
							tokensToEvaluate = append(tokensToEvaluate, Token{Value: "\n", FileName: tokenArray[x].FileName, Type: TOKEN_TYPE_NEWLINE, Line: tokenArray[x].Line, Column: tokenArray[x].Column})
						}
					}
				}
				if !isLoopStatement && !isWhileLoopStatement {
					if tokenArray[x].Type == TOKEN_TYPE_IF_START {
						openIfCount += 1
						ignoreNewline = true
						isIfStatement = true
					} else if tokenArray[x].Type == TOKEN_TYPE_IF_END {
						openIfCount = openIfCount - 1

						if openIfCount == 0 {
							ignoreNewline = false
							isIfStatement = false
							tokensToEvaluate = append(tokensToEvaluate, Token{Value: "\n", FileName: tokenArray[x].FileName, Type: TOKEN_TYPE_NEWLINE, Line: tokenArray[x].Line, Column: tokenArray[x].Column})
						}
					}
				}
			}
			//put the token to stack for shunting yard process later
			tokensToEvaluate = append(tokensToEvaluate, tokenArray[x])
		} else {

			if isWhileLoopStatement && strings.Index(tokenArray[x].Context, CONTEXT_NAME_PREFIX_WHILE_LOOP) == 0 {
				//collect all the parameters from while loop
				//to be able to parse later on the body of while loop
				if _, ok := whileLoopIsDone[tokenArray[x].Context]; ok {
				} else {
					whileLoopIsDone[tokenArray[x].Context] = false
				}

				if tokenArray[x].Type == TOKEN_TYPE_WHILE_LOOP_PARAM_END {
					whileLoopIsDone[tokenArray[x].Context] = true
				}

				if !whileLoopIsDone[tokenArray[x].Context] {
					whileLoopTokens[tokenArray[x].Context] = append(whileLoopTokens[tokenArray[x].Context], tokenArray[x])
				}
			}

			//put the token to stack for shunting yard process later
			tokensToEvaluate = append(tokensToEvaluate, tokenArray[x])
		}
	}

	if getLastStackBool {

		if lastToken.Type == TOKEN_TYPE_IDENTIFIER {
			var errConvert error
			lastToken, errConvert = convertVariableToToken(lastToken, *globalVariableArray, scopeName)
			if errConvert != nil {
				return errConvert
			}
		}

		*lastStackBool = convertTokenToBool(lastToken)
	}

	return nil
}
