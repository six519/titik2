package parser

import (
	"errors"
	"github.com/six519/titik2/lexer"
	"github.com/six519/titik2/info"
	"github.com/six519/titik2/variable"
)

//parser states
const PARSER_STATE_START int = 0

type Parser struct {
}

func (parser Parser) Parse(tokenArray []lexer.Token, globalVariableArray *[]variable.Variable) error {
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
				//only accepts identifier & keywords at start
				switch(strippedTokenArray[x].Type) {
					case lexer.TOKEN_TYPE_IDENTIFIER:
						//identifier type
					case lexer.TOKEN_TYPE_KEYWORD:
						//keyword type

					default:
						//token error
						return errors.New(info.ErrorMessage(false, strippedTokenArray[x].Line, strippedTokenArray[x].Column, "Unexpected token", strippedTokenArray[x].FileName))
				}
			default:
				continue
		}
	}

	return nil
}