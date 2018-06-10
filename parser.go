package main

import (
	//"errors"
	//"strconv"
	//"fmt"
)

type Parser struct {
}

func (parser Parser) Parse(tokenArray []Token, globalVariableArray *[]Variable, scopeName string) error {
	var currentTokens []Token
	for x := 0; x < len(tokenArray); x++ {
		//ignore space, tab and comments
		if(tokenArray[x].Type != TOKEN_TYPE_SPACE && tokenArray[x].Type != TOKEN_TYPE_SINGLE_COMMENT && tokenArray[x].Type != TOKEN_TYPE_MULTI_COMMENT && tokenArray[x].Type != TOKEN_TYPE_CLOSE_MULTI_COMMENT && tokenArray[x].Type != TOKEN_TYPE_CLOSE_STRING && tokenArray[x].Type != TOKEN_TYPE_TAB) {
			if(tokenArray[x].Type != TOKEN_TYPE_NEWLINE) {
				//execute shunting yard
			} else {
				//put the token to stack for shunting yard process later
				currentTokens = append(currentTokens, tokenArray[x])
			}
		}
	}

	return nil
}