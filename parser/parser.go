package parser

import (
	"github.com/six519/titik2/lexer"
)

//parser states
const PARSER_STATE_START int = 0

type Parser struct {
}

func (parser Parser) Parse(tokenArray []lexer.Token) error {
	var strippedTokenArray []lexer.Token
	parserState := PARSER_STATE_START

	//remove spaces, newlines, comments and tabs
	for x := 0; x < len(tokenArray); x++ {
		if(tokenArray[x].Type != lexer.TOKEN_TYPE_NEWLINE && tokenArray[x].Type != lexer.TOKEN_TYPE_SPACE && tokenArray[x].Type != lexer.TOKEN_TYPE_SINGLE_COMMENT && tokenArray[x].Type != lexer.TOKEN_TYPE_MULTI_COMMENT && tokenArray[x].Type != lexer.TOKEN_TYPE_CLOSE_MULTI_COMMENT && tokenArray[x].Type != lexer.TOKEN_TYPE_CLOSE_STRING && tokenArray[x].Type != lexer.TOKEN_TYPE_TAB) {
			strippedTokenArray = append(strippedTokenArray, tokenArray[x])
		}
	}

	//parse the stripped tokens below
	for x := 0; x < len(strippedTokenArray); x++ {
		switch(parserState) {
			case PARSER_STATE_START:
				//start
			default:
				continue
		}
	}

	return nil
}