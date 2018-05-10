package lexer

import (
	"os"
	"bufio"
	"errors"
	"unicode"
	"github.com/six519/titik2/info"
	//"fmt"
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

func (lexer Lexer) GenerateToken() ([]map[string]string, error) {
	var tokenArray []map[string]string
	tokenizerState := TOKENIZER_STATE_GET_WORD
	isTokenInit := false

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

						if(!isTokenInit) {

						}
					} else if(currentChar == "\n") {
						//new line
					} else if(currentChar == "\t") {
						//tab
					} else if(currentChar == " ") {
						//space
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
