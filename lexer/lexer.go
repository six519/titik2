package lexer

import (
	"os"
	"bufio"
	"errors"
	"unicode"
	"github.com/six519/titik2/info"
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
const TOKEN_TYPE_NONE = 33

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
		tokenType := ""

		switch(tokenArray[x].Type) {
			case TOKEN_TYPE_IDENTIFIER:
				tokenType = "TOKEN_TYPE_IDENTIFIER"
			case TOKEN_TYPE_PERIOD:
				tokenType = "TOKEN_TYPE_PERIOD"
			case TOKEN_TYPE_KEYWORD:
				tokenType = "TOKEN_TYPE_KEYWORD"
			case TOKEN_TYPE_SPACE:
				tokenType = "TOKEN_TYPE_SPACE"
			case TOKEN_TYPE_OPEN_PARENTHESIS:
				tokenType = "TOKEN_TYPE_OPEN_PARENTHESIS"
			case TOKEN_TYPE_CLOSE_PARENTHESIS:
				tokenType = "TOKEN_TYPE_CLOSE_PARENTHESIS"
			case TOKEN_TYPE_SINGLE_COMMENT:
				tokenType = "TOKEN_TYPE_SINGLE_COMMENT"
			case TOKEN_TYPE_MULTI_COMMENT:
				tokenType = "TOKEN_TYPE_MULTI_COMMENT"
			case TOKEN_TYPE_CLOSE_MULTI_COMMENT:
				tokenType = "TOKEN_TYPE_CLOSE_MULTI_COMMENT"
			case TOKEN_TYPE_NEWLINE:
				tokenType = "TOKEN_TYPE_NEWLINE"
			case TOKEN_TYPE_EQUALS:
				tokenType = "TOKEN_TYPE_EQUALS"
			case TOKEN_TYPE_STRING:
				tokenType = "TOKEN_TYPE_STRING"
			case TOKEN_TYPE_CLOSE_STRING:
				tokenType = "TOKEN_TYPE_CLOSE_STRING"
			case TOKEN_TYPE_INTEGER:
				tokenType = "TOKEN_TYPE_INTEGER"
			case TOKEN_TYPE_FLOAT:
				tokenType = "TOKEN_TYPE_FLOAT"
			case TOKEN_TYPE_PLUS:
				tokenType = "TOKEN_TYPE_PLUS"
			case TOKEN_TYPE_MINUS:
				tokenType = "TOKEN_TYPE_MINUS"
			case TOKEN_TYPE_DIVIDE:
				tokenType = "TOKEN_TYPE_DIVIDE"
			case TOKEN_TYPE_MULTIPLY:
				tokenType = "TOKEN_TYPE_MULTIPLY"
			case TOKEN_TYPE_COLON:
				tokenType = "TOKEN_TYPE_COLON"
			case TOKEN_TYPE_SEMI_COLON:
				tokenType = "TOKEN_TYPE_SEMI_COLON"
			case TOKEN_TYPE_COMMA:
				tokenType = "TOKEN_TYPE_COMMA"
			case TOKEN_TYPE_OPEN_BRACKET:
				tokenType = "TOKEN_TYPE_OPEN_BRACKET"
			case TOKEN_TYPE_CLOSE_BRACKET:
				tokenType = "TOKEN_TYPE_CLOSE_BRACKET"
			case TOKEN_TYPE_OPEN_BRACES:
				tokenType = "TOKEN_TYPE_OPEN_BRACES"
			case TOKEN_TYPE_CLOSE_BRACES:
				tokenType = "TOKEN_TYPE_CLOSE_BRACES"
			case TOKEN_TYPE_AMPERSAND:
				tokenType = "TOKEN_TYPE_AMPERSAND"
			case TOKEN_TYPE_GREATER_THAN:
				tokenType = "TOKEN_TYPE_GREATER_THAN"
			case TOKEN_TYPE_LESS_THAN:
				tokenType = "TOKEN_TYPE_LESS_THAN"
			case TOKEN_TYPE_OR:
				tokenType = "TOKEN_TYPE_OR"
			case TOKEN_TYPE_EXCLAMATION:
				tokenType = "TOKEN_TYPE_EXCLAMATION"
			case TOKEN_TYPE_TAB:
				tokenType = "TOKEN_TYPE_TAB"
			case TOKEN_TYPE_CARRIAGE_RETURN:
				tokenType = "TOKEN_TYPE_CARRIAGE_RETURN"
			default:
				tokenType = "TOKEN_TYPE_NONE"
		}

        fmt.Printf("Token Type: %s\n", tokenType)
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

func (lexer Lexer) GenerateToken() ([]Token, error) {
	var tokenArray []Token
	tokenizerState := TOKENIZER_STATE_GET_WORD
	isTokenInit := false
	usePeriod := true
	withPeriod := false

	//read the file contents line by line
	for x := 0; x < len(lexer.fileContents); x++ {
		//read character by character
		for x2 := 0; x2 < len(lexer.fileContents[x]); x2++ {
			currentChar := string(lexer.fileContents[x][x2])

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
						return tokenArray, errors.New(info.TokenErrorMessage(x + 1, x2 + 1, "Invalid token", lexer.FileName))
					}
				case TOKENIZER_STATE_GET_SINGLE_COMMENT:
					//get single comment
				case TOKENIZER_STATE_GET_MULTI_COMMENT:
					//get multi comment
				case TOKENIZER_STATE_GET_STRING:
					//get string
				case TOKENIZER_STATE_GET_FLOAT:
					//get float
				default:
					continue
			}
		}
	}

	return tokenArray, nil
}
