package parser

import (
	"github.com/six519/titik2/lexer"
)

type Parser struct {
}

func (parser Parser) Parse(tokenArray []lexer.Token) error {
	var strippedTokenArray []lexer.Token

	//remove spaces, newlines, comments and tabs
	for x := 0; x < len(tokenArray); x++ {
		if(tokenArray[x].Type != lexer.TOKEN_TYPE_NEWLINE && tokenArray[x].Type != lexer.TOKEN_TYPE_SPACE && tokenArray[x].Type != lexer.TOKEN_TYPE_SINGLE_COMMENT && tokenArray[x].Type != lexer.TOKEN_TYPE_MULTI_COMMENT && tokenArray[x].Type != lexer.TOKEN_TYPE_CLOSE_MULTI_COMMENT && tokenArray[x].Type != lexer.TOKEN_TYPE_CLOSE_STRING && tokenArray[x].Type != lexer.TOKEN_TYPE_TAB) {
			strippedTokenArray = append(strippedTokenArray, tokenArray[x])
		}
	}

	//parse the stripped tokens below

	return nil
}