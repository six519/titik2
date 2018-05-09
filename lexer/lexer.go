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
			lexer.fileContents = append(lexer.fileContents, scanner.Text())
		}

	} else {
		return errors.New("Error: no such file: '" + lexer.FileName + "'")
	}

	return nil
}

func (lexer Lexer) GenerateToken() ([]map[string]string, error) {
	var tokenArray []map[string]string
	tokenizerState := TOKENIZER_STATE_GET_WORD

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
