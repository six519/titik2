package parser

import (
	"github.com/six519/titik2/lexer"
	"github.com/six519/titik2/variable"
)

type Parser struct {
}

func (parser Parser) Parse(tokenArray []lexer.Token, globalVariableArray *[]variable.Variable) error {

	for x := 0; x < len(tokenArray); x++ {
		//ignore newline, space, tab and comments
		if(tokenArray[x].Type != lexer.TOKEN_TYPE_NEWLINE && tokenArray[x].Type != lexer.TOKEN_TYPE_SPACE && tokenArray[x].Type != lexer.TOKEN_TYPE_SINGLE_COMMENT && tokenArray[x].Type != lexer.TOKEN_TYPE_MULTI_COMMENT && tokenArray[x].Type != lexer.TOKEN_TYPE_CLOSE_MULTI_COMMENT && tokenArray[x].Type != lexer.TOKEN_TYPE_CLOSE_STRING && tokenArray[x].Type != lexer.TOKEN_TYPE_TAB) {
			//TODO: TRY TO APPLY SHUNTING-YARD BELOW
		}
	}

	return nil
}