package main

import (
	"errors"
)

type Parser struct {
}

func (parser Parser) Parse(tokenArray []Token, globalVariableArray *[]Variable) error {

	var finalTokenArray []Token
	//operatorPrecedences := map[string] int{"+": 0, "-": 0, "/": 1, "*": 1} //operator order of precedences
	//var operatorStack []Token
	var outputQueue []Token

	for x := 0; x < len(tokenArray); x++ {
		//ignore newline, space, tab and comments
		if(tokenArray[x].Type != TOKEN_TYPE_NEWLINE && tokenArray[x].Type != TOKEN_TYPE_SPACE && tokenArray[x].Type != TOKEN_TYPE_SINGLE_COMMENT && tokenArray[x].Type != TOKEN_TYPE_MULTI_COMMENT && tokenArray[x].Type != TOKEN_TYPE_CLOSE_MULTI_COMMENT && tokenArray[x].Type != TOKEN_TYPE_CLOSE_STRING && tokenArray[x].Type != TOKEN_TYPE_TAB) {
			finalTokenArray = append(finalTokenArray, tokenArray[x])
		}
	}

	if(len(finalTokenArray) > 0) {
		if(finalTokenArray[0].Type == TOKEN_TYPE_PLUS || finalTokenArray[0].Type == TOKEN_TYPE_MINUS || finalTokenArray[0].Type == TOKEN_TYPE_DIVIDE || finalTokenArray[0].Type == TOKEN_TYPE_MULTIPLY) {
			//syntax error if the first token is an operator
			return errors.New(SyntaxErrorMessage(finalTokenArray[0].Line, finalTokenArray[0].Column, "Unexpected token '" + finalTokenArray[0].Value + "'", finalTokenArray[0].FileName))
		}

		if(finalTokenArray[len(finalTokenArray)-1].Type == TOKEN_TYPE_PLUS || finalTokenArray[len(finalTokenArray)-1].Type == TOKEN_TYPE_MINUS || finalTokenArray[len(finalTokenArray)-1].Type == TOKEN_TYPE_DIVIDE || finalTokenArray[len(finalTokenArray)-1].Type == TOKEN_TYPE_MULTIPLY) {
			//syntax error if the last token is an operator
			return errors.New(SyntaxErrorMessage(finalTokenArray[len(finalTokenArray)-1].Line, finalTokenArray[len(finalTokenArray)-1].Column, "Unfinished operation", finalTokenArray[len(finalTokenArray)-1].FileName))
		}

		//shunting-yard
		for len(finalTokenArray) > 0 {
			currentToken := finalTokenArray[0]
			finalTokenArray = append(finalTokenArray[:0], finalTokenArray[1:]...) //pop the first element

			if(currentToken.Type == TOKEN_TYPE_INTEGER || currentToken.Type == TOKEN_TYPE_FLOAT) {
				//If it's a number add it to queue
				outputQueue = append(outputQueue, currentToken)
			}

		}
	}

	return nil
}