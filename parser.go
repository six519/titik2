package main

import (
	"errors"
)

type Parser struct {
}

func (parser Parser) Parse(tokenArray []Token, globalVariableArray *[]Variable) error {

	operatorPrecedences := map[string] int{"+": 0, "-": 0, "/": 1, "*": 1} //operator order of precedences
	var operatorStack []Token
	var outputQueue []Token

	for x := 0; x < len(tokenArray); x++ {
		//ignore newline, space, tab and comments
		if(tokenArray[x].Type != TOKEN_TYPE_NEWLINE && tokenArray[x].Type != TOKEN_TYPE_SPACE && tokenArray[x].Type != TOKEN_TYPE_SINGLE_COMMENT && tokenArray[x].Type != TOKEN_TYPE_MULTI_COMMENT && tokenArray[x].Type != TOKEN_TYPE_CLOSE_MULTI_COMMENT && tokenArray[x].Type != TOKEN_TYPE_CLOSE_STRING && tokenArray[x].Type != TOKEN_TYPE_TAB) {

			if(tokenArray[x].Type == TOKEN_TYPE_PLUS || tokenArray[x].Type == TOKEN_TYPE_MINUS || tokenArray[x].Type == TOKEN_TYPE_DIVIDE || tokenArray[x].Type == TOKEN_TYPE_MULTIPLY) {
				//syntax error if the first token is an operator
				return errors.New(SyntaxErrorMessage(tokenArray[x].Line, tokenArray[x].Column, "Unexpected token '" + tokenArray[x].Value + "'", tokenArray[x].FileName))
			}

		}
	}

	return nil
}