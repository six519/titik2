package main

import (
	"os"
	"bufio"
	"errors"
	"unicode"
	"fmt"
)

//Tokenizer States
const TOKENIZER_STATE_GET_WORD int = 0
const TOKENIZER_STATE_GET_SINGLE_COMMENT int = 1
const TOKENIZER_STATE_GET_MULTI_COMMENT int = 2
const TOKENIZER_STATE_GET_STRING int = 3
const TOKENIZER_STATE_GET_FLOAT int = 4

//Token Types
const TOKEN_TYPE_IDENTIFIER int = 0
const TOKEN_TYPE_PERIOD int = 1
const TOKEN_TYPE_KEYWORD int = 2
const TOKEN_TYPE_SPACE int = 3
const TOKEN_TYPE_OPEN_PARENTHESIS int = 4
const TOKEN_TYPE_CLOSE_PARENTHESIS int = 5
const TOKEN_TYPE_SINGLE_COMMENT int = 6
const TOKEN_TYPE_MULTI_COMMENT int = 7
const TOKEN_TYPE_CLOSE_MULTI_COMMENT int = 8
const TOKEN_TYPE_NEWLINE int = 9
const TOKEN_TYPE_EQUALS int = 10
const TOKEN_TYPE_STRING int = 11
const TOKEN_TYPE_CLOSE_STRING int = 12
const TOKEN_TYPE_INTEGER int = 13
const TOKEN_TYPE_FLOAT int = 14
const TOKEN_TYPE_PLUS int = 15
const TOKEN_TYPE_MINUS int = 16
const TOKEN_TYPE_DIVIDE int = 17
const TOKEN_TYPE_MULTIPLY int = 18
const TOKEN_TYPE_COLON int = 19
const TOKEN_TYPE_SEMI_COLON int = 20
const TOKEN_TYPE_COMMA int = 21
const TOKEN_TYPE_OPEN_BRACKET int = 22
const TOKEN_TYPE_CLOSE_BRACKET int = 23
const TOKEN_TYPE_OPEN_BRACES int = 24
const TOKEN_TYPE_CLOSE_BRACES int = 25
const TOKEN_TYPE_AMPERSAND int = 26
const TOKEN_TYPE_GREATER_THAN int = 27
const TOKEN_TYPE_LESS_THAN int = 28
const TOKEN_TYPE_OR int = 29
const TOKEN_TYPE_EXCLAMATION int = 30
const TOKEN_TYPE_TAB int = 31
const TOKEN_TYPE_CARRIAGE_RETURN int = 32
const TOKEN_TYPE_NONE int = 33
const TOKEN_TYPE_FUNCTION int = 34
const TOKEN_TYPE_INVOKE_FUNCTION int = 35

//for debugging purpose only
var TOKEN_TYPES_STRING = []string {
	"TOKEN_TYPE_IDENTIFIER",
	"TOKEN_TYPE_PERIOD",
	"TOKEN_TYPE_KEYWORD",
	"TOKEN_TYPE_SPACE",
	"TOKEN_TYPE_OPEN_PARENTHESIS",
	"TOKEN_TYPE_CLOSE_PARENTHESIS",
	"TOKEN_TYPE_SINGLE_COMMENT",
	"TOKEN_TYPE_MULTI_COMMENT",
	"TOKEN_TYPE_CLOSE_MULTI_COMMENT",
	"TOKEN_TYPE_NEWLINE",
	"TOKEN_TYPE_EQUALS",
	"TOKEN_TYPE_STRING",
	"TOKEN_TYPE_CLOSE_STRING",
	"TOKEN_TYPE_INTEGER",
	"TOKEN_TYPE_FLOAT",
	"TOKEN_TYPE_PLUS",
	"TOKEN_TYPE_MINUS",
	"TOKEN_TYPE_DIVIDE",
	"TOKEN_TYPE_MULTIPLY",
	"TOKEN_TYPE_COLON",
	"TOKEN_TYPE_SEMI_COLON",
	"TOKEN_TYPE_COMMA",
	"TOKEN_TYPE_OPEN_BRACKET",
	"TOKEN_TYPE_CLOSE_BRACKET",
	"TOKEN_TYPE_OPEN_BRACES",
	"TOKEN_TYPE_CLOSE_BRACES",
	"TOKEN_TYPE_AMPERSAND",
	"TOKEN_TYPE_GREATER_THAN",
	"TOKEN_TYPE_LESS_THAN",
	"TOKEN_TYPE_OR",
	"TOKEN_TYPE_EXCLAMATION",
	"TOKEN_TYPE_TAB",
	"TOKEN_TYPE_CARRIAGE_RETURN",
	"TOKEN_TYPE_NONE",
	"TOKEN_TYPE_FUNCTION",
	"TOKEN_TYPE_INVOKE_FUNCTION",
}

//token object
type Token struct {
	Value string
	FileName string
	Type int
	Line int
	Column int
}

//function helpers
func setToken(initToken bool, tokenArray *[]Token, isTokenInit *bool, lineNumber int, colNumber int, tokenType int, tokenFileName string, tokenValue string) {
	token := Token{Value: tokenValue, FileName: tokenFileName, Type: tokenType, Line: lineNumber, Column: colNumber}
	*tokenArray = append(*tokenArray, token)
	*isTokenInit = initToken
}

func DumpToken(tokenArray []Token) {
	fmt.Printf("====================================\n")

	for x := 0; x < len(tokenArray); x++ {
        fmt.Printf("Token Type: %s\n", TOKEN_TYPES_STRING[tokenArray[x].Type])
        fmt.Printf("Line #: %d\n", tokenArray[x].Line)
        fmt.Printf("Column #: %d\n", tokenArray[x].Column)
        fmt.Printf("Value: %s\n", tokenArray[x].Value)
		fmt.Printf("====================================\n")
	}
}

//lexer object
type Lexer struct {
	FileName string
	fileContents []string
}

func (lexer *Lexer) ReadSourceFile() error {

	file, err := os.Open(lexer.FileName)
	defer file.Close()

	if (err == nil) {
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			lexer.fileContents = append(lexer.fileContents, scanner.Text() + "\n")
		}

	} else {
		return errors.New("Error: no such file: '" + lexer.FileName + "'")
	}

	return nil
}

func (lexer *Lexer) ReadString(inputString string) {
	lexer.fileContents = append(lexer.fileContents, inputString)
}

func (lexer Lexer) GenerateToken() ([]Token, error) {
	var tokenArray []Token
	var cleanTokenArray []Token
	var finalTokenArray []Token
	tokenizerState := TOKENIZER_STATE_GET_WORD
	isTokenInit := false
	usePeriod := true
	withPeriod := false
	stringOpener := "\""

	//read the file contents line by line
	for x := 0; x < len(lexer.fileContents); x++ {
		//read character by character
		for x2 := 0; x2 < len(lexer.fileContents[x]); x2++ {
			currentChar := string(lexer.fileContents[x][x2])

			switch_label:
			switch(tokenizerState) {
				case TOKENIZER_STATE_GET_WORD:
					//get word
					if(unicode.IsLetter([]rune(currentChar)[0]) || currentChar == "_") {
						//alphabetic or underscore
						if(len(tokenArray) > 0) {
							if(tokenArray[len(tokenArray) - 1].Type == TOKEN_TYPE_INTEGER) {
								isTokenInit = false
							}
						}
						if(!isTokenInit) {
							setToken(true, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_IDENTIFIER, lexer.FileName, "") //init token
						}
						tokenArray[len(tokenArray) - 1].Value += currentChar
					} else if(currentChar == "\n") {
						//new line
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_NEWLINE, lexer.FileName, currentChar) //set token
					} else if(currentChar == "\t") {
						//tab
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_TAB, lexer.FileName, currentChar) //set token
					} else if(currentChar == " ") {
						//space
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_SPACE, lexer.FileName, currentChar) //set token
					} else if(currentChar == "=") {
						//equals
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_EQUALS, lexer.FileName, currentChar) //set token
					} else if(currentChar == "+") {
						//plus
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_PLUS, lexer.FileName, currentChar) //set token
					} else if(currentChar == "-") {
						//minus
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_MINUS, lexer.FileName, currentChar) //set token
					} else if(currentChar == "/") {
						//divide
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_DIVIDE, lexer.FileName, currentChar) //set token
					} else if(currentChar == "*") {
						//multiply
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_MULTIPLY, lexer.FileName, currentChar) //set token
					} else if(currentChar == ".") {
						//period
						usePeriod = true
						withPeriod = false
						if(len(tokenArray) > 0) {
							if(tokenArray[len(tokenArray) - 1].Type == TOKEN_TYPE_INTEGER) {
								tokenizerState = TOKENIZER_STATE_GET_FLOAT
								usePeriod = false
							}
						}
						if(usePeriod) {
							setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_PERIOD, lexer.FileName, currentChar) //set token
						}
					} else if(currentChar == "(") {
						//open parenthesis
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_OPEN_PARENTHESIS, lexer.FileName, currentChar) //set token
					} else if(currentChar == ")") {
						//close parenthesis
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_CLOSE_PARENTHESIS, lexer.FileName, currentChar) //set token
					} else if(currentChar == "{") {
						//open bracket
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_OPEN_BRACKET, lexer.FileName, currentChar) //set token
					} else if(currentChar == "}") {
						//close brackeT
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_CLOSE_BRACKET, lexer.FileName, currentChar) //set token
					} else if(currentChar == "[") {
						//open braces
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_OPEN_BRACES, lexer.FileName, currentChar) //set token
					} else if(currentChar == "]") {
						//close braces
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_CLOSE_BRACES, lexer.FileName, currentChar) //set token
					} else if(currentChar == ":") {
						//colon
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_COLON, lexer.FileName, currentChar) //set token
					} else if(currentChar == ",") {
						//comma
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_COMMA, lexer.FileName, currentChar) //set token
					} else if(currentChar == ";") {
						//semi colon
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_SEMI_COLON, lexer.FileName, currentChar) //set token
					} else if(currentChar == "&") {
						//ampersand
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_AMPERSAND, lexer.FileName, currentChar) //set token
					} else if(currentChar == ">") {
						//greater than
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_GREATER_THAN, lexer.FileName, currentChar) //set token
					} else if(currentChar == "<") {
						//less than
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_LESS_THAN, lexer.FileName, currentChar) //set token
					} else if(currentChar == "|") {
						//or
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_OR, lexer.FileName, currentChar) //set token
					} else if(currentChar == "'" || currentChar == "\"") {
						//string
						stringOpener = currentChar
						setToken(true, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_STRING, lexer.FileName, "") //init token
						tokenizerState = TOKENIZER_STATE_GET_STRING
					} else if(currentChar == "#") {
						//start of single line comment
						setToken(true, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_SINGLE_COMMENT, lexer.FileName, "") //init token
						tokenizerState = TOKENIZER_STATE_GET_SINGLE_COMMENT
					} else if(currentChar == "\\") {
						//start of multiline comment
						setToken(true, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_MULTI_COMMENT, lexer.FileName, "") //init token
						tokenizerState = TOKENIZER_STATE_GET_MULTI_COMMENT
					} else if(unicode.IsDigit([]rune(currentChar)[0])) {
						//integer
						if(!isTokenInit) {
							setToken(true, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_INTEGER, lexer.FileName, "") //init token
						}
						tokenArray[len(tokenArray) - 1].Value += currentChar
					} else {
						return finalTokenArray, errors.New(SyntaxErrorMessage(x + 1, x2 + 1, "Invalid token", lexer.FileName))
					}
				case TOKENIZER_STATE_GET_SINGLE_COMMENT:
					//get single comment
					if(currentChar == "\n") {
						tokenizerState = TOKENIZER_STATE_GET_WORD
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_NEWLINE, lexer.FileName, currentChar) //set token
					} else {
						tokenArray[len(tokenArray) - 1].Value += currentChar
					}
				case TOKENIZER_STATE_GET_MULTI_COMMENT:
					//get multi comment
					if(currentChar == "\\") {
						tokenizerState = TOKENIZER_STATE_GET_WORD
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_CLOSE_MULTI_COMMENT, lexer.FileName, currentChar) //set token
					} else {
						tokenArray[len(tokenArray) - 1].Value += currentChar
					}
				case TOKENIZER_STATE_GET_STRING:
					//get string
					if(currentChar == stringOpener) {
						tokenizerState = TOKENIZER_STATE_GET_WORD
						setToken(false, &tokenArray, &isTokenInit, x + 1, x2 + 1, TOKEN_TYPE_CLOSE_STRING, lexer.FileName, currentChar) //set token
					} else {
						tokenArray[len(tokenArray) - 1].Value += currentChar
					}
				case TOKENIZER_STATE_GET_FLOAT:
					//get float
					if(unicode.IsDigit([]rune(currentChar)[0])) {
						if(!withPeriod) {
							withPeriod = true
							tokenArray[len(tokenArray) - 1].Value += "."
							tokenArray[len(tokenArray) - 1].Type = TOKEN_TYPE_FLOAT
						}

						tokenArray[len(tokenArray) - 1].Value += currentChar
					} else {
						if(!withPeriod) {
							setToken(false, &tokenArray, &isTokenInit, x + 1, x2, TOKEN_TYPE_PERIOD, lexer.FileName, ".") //set token
						}
						tokenizerState = TOKENIZER_STATE_GET_WORD
						isTokenInit = false
						goto switch_label
					}
				default:
					continue
			}
		}
	}

	//1st token cleanup
	//TODO: Try to eliminate this part
	for x := 0; x < len(tokenArray); x++ {
		if(tokenArray[x].Type == TOKEN_TYPE_SPACE || tokenArray[x].Type == TOKEN_TYPE_SINGLE_COMMENT || tokenArray[x].Type == TOKEN_TYPE_MULTI_COMMENT || tokenArray[x].Type == TOKEN_TYPE_CLOSE_MULTI_COMMENT || tokenArray[x].Type == TOKEN_TYPE_CLOSE_STRING || tokenArray[x].Type == TOKEN_TYPE_TAB) {
			//ignore space, tab, comments etc...
			continue
		}
		cleanTokenArray = append(cleanTokenArray, tokenArray[x])
	}

	//2nd token cleanup 
	var ignoreOpenP bool = false
	var f_count int = 0
	var op_count int = 0
	for x := 0; x < len(cleanTokenArray); x++ {
		if(ignoreOpenP) {
			//ignore open parenthesis if the last token is a function
			ignoreOpenP = false
			continue
		}
		if(cleanTokenArray[x].Type == TOKEN_TYPE_OPEN_PARENTHESIS) {
			op_count += 1
		}
		if(cleanTokenArray[x].Type == TOKEN_TYPE_CLOSE_PARENTHESIS) {
			if(op_count > 0) {
				op_count -= 1
			} else {
				if(f_count > 0) {
					f_count -= 1
					cleanTokenArray[x].Type = TOKEN_TYPE_INVOKE_FUNCTION
				}
			}
		}
		if(cleanTokenArray[x].Type == TOKEN_TYPE_IDENTIFIER) {
			if(IsReservedWord(cleanTokenArray[x].Value)) {
				//Convert identifier to keyword if existing in reserved words
				cleanTokenArray[x].Type = TOKEN_TYPE_KEYWORD
			} else {
				//Check if the next token is '(', if yes then it's a function call
				if((x + 1) <= len(cleanTokenArray) - 1 ) {
					if(cleanTokenArray[x+1].Type == TOKEN_TYPE_OPEN_PARENTHESIS) {
						cleanTokenArray[x].Type = TOKEN_TYPE_FUNCTION
						ignoreOpenP = true
						f_count += 1
					}
				}
			}
		} else if(cleanTokenArray[x].Type == TOKEN_TYPE_STRING) {
			if(!((x+1) < len(cleanTokenArray))) {
				return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x].Line, cleanTokenArray[x].Column, "Expected closing of string", lexer.FileName))
			}
		} else if(cleanTokenArray[x].Type == TOKEN_TYPE_MULTI_COMMENT) {
			if(!((x+1) < len(cleanTokenArray))) {
				return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x].Line, cleanTokenArray[x].Column, "Expected closing of multi line comment", lexer.FileName))
			}
		}

		if(x != 0 && (cleanTokenArray[x].Type == TOKEN_TYPE_FLOAT || cleanTokenArray[x].Type == TOKEN_TYPE_INTEGER)) {
			if(cleanTokenArray[x - 1].Type == TOKEN_TYPE_MINUS) {

				if((x - 2) >= 0) {
					if(cleanTokenArray[x-2].Type != TOKEN_TYPE_FLOAT && cleanTokenArray[x-2].Type != TOKEN_TYPE_INTEGER && cleanTokenArray[x-2].Type != TOKEN_TYPE_IDENTIFIER) {
						//set to negative number
						finalTokenArray[len(finalTokenArray) - 1].Value += cleanTokenArray[x].Value
						finalTokenArray[len(finalTokenArray) - 1].Type = cleanTokenArray[x].Type
					} else {
						finalTokenArray = append(finalTokenArray, cleanTokenArray[x])
					}
				} else {
					finalTokenArray = append(finalTokenArray, cleanTokenArray[x])
				}
				
			} else {
				finalTokenArray = append(finalTokenArray, cleanTokenArray[x])
			}
		} else {
			finalTokenArray = append(finalTokenArray, cleanTokenArray[x])
		}
	}

	return finalTokenArray, nil
}
