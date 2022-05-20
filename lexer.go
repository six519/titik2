package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

//Tokenizer States
const (
	TOKENIZER_STATE_GET_WORD = iota
	TOKENIZER_STATE_GET_SINGLE_COMMENT
	TOKENIZER_STATE_GET_MULTI_COMMENT
	TOKENIZER_STATE_GET_STRING
	TOKENIZER_STATE_GET_FLOAT
)

//Token Types
const (
	TOKEN_TYPE_IDENTIFIER = iota
	TOKEN_TYPE_PERIOD
	TOKEN_TYPE_KEYWORD
	TOKEN_TYPE_SPACE
	TOKEN_TYPE_OPEN_PARENTHESIS
	TOKEN_TYPE_CLOSE_PARENTHESIS
	TOKEN_TYPE_SINGLE_COMMENT
	TOKEN_TYPE_MULTI_COMMENT
	TOKEN_TYPE_CLOSE_MULTI_COMMENT
	TOKEN_TYPE_NEWLINE
	TOKEN_TYPE_EQUALS
	TOKEN_TYPE_STRING
	TOKEN_TYPE_CLOSE_STRING
	TOKEN_TYPE_INTEGER
	TOKEN_TYPE_FLOAT
	TOKEN_TYPE_PLUS
	TOKEN_TYPE_MINUS
	TOKEN_TYPE_DIVIDE
	TOKEN_TYPE_MULTIPLY
	TOKEN_TYPE_COLON
	TOKEN_TYPE_SEMI_COLON
	TOKEN_TYPE_COMMA
	TOKEN_TYPE_OPEN_BRACKET
	TOKEN_TYPE_CLOSE_BRACKET
	TOKEN_TYPE_OPEN_BRACES
	TOKEN_TYPE_CLOSE_BRACES
	TOKEN_TYPE_AMPERSAND
	TOKEN_TYPE_GREATER_THAN
	TOKEN_TYPE_LESS_THAN
	TOKEN_TYPE_OR
	TOKEN_TYPE_TAB
	TOKEN_TYPE_CARRIAGE_RETURN
	TOKEN_TYPE_NONE
	TOKEN_TYPE_FUNCTION
	TOKEN_TYPE_INVOKE_FUNCTION
	TOKEN_TYPE_FUNCTION_DEF_START
	TOKEN_TYPE_FUNCTION_PARAM_END
	TOKEN_TYPE_FUNCTION_DEF_END
	TOKEN_TYPE_FUNCTION_RETURN
	TOKEN_TYPE_FOR_LOOP_START
	TOKEN_TYPE_FOR_LOOP_PARAM_END
	TOKEN_TYPE_FOR_LOOP_END
	TOKEN_TYPE_LOOP_BREAK
	TOKEN_TYPE_BOOLEAN
	TOKEN_TYPE_EQUALITY
	TOKEN_TYPE_INEQUALITY
	TOKEN_TYPE_LESS_THAN_OR_EQUALS
	TOKEN_TYPE_GREATER_THAN_OR_EQUALS
	TOKEN_TYPE_IF_START
	TOKEN_TYPE_IF_PARAM_END
	TOKEN_TYPE_IF_END
	TOKEN_TYPE_ELSE
	TOKEN_TYPE_ELIF_START
	TOKEN_TYPE_ELIF_PARAM_END
	TOKEN_TYPE_GET_ARRAY_START
	TOKEN_TYPE_GET_ARRAY_END
	TOKEN_TYPE_ARRAY
	TOKEN_TYPE_KEY_VALUE_PAIR
	TOKEN_TYPE_ASSOCIATIVE_ARRAY
	TOKEN_TYPE_WHILE_LOOP_START
	TOKEN_TYPE_WHILE_LOOP_PARAM_END
	TOKEN_TYPE_WHILE_LOOP_END
)

//for debugging purpose only
var TOKEN_TYPES_STRING = []string{
	"TOKEN_TYPE_IDENTIFIER",
	"TOKEN_TYPE_PERIOD",
	"TOKEN_TYPE_KEYWORD",
	"TOKEN_TYPE_SPACE",
	"TOKEN_TYPE_OPEN_PARENTHESIS",
	"TOKEN_TYPE_CLOSE_PARENTHESIS",
	"TOKEN_TYPE_SINGLE_COMMENT",
	"TOKEN_TYPE_MULTI_COMMENT",
	"TOKEN_TYPE_CLOSE_MULTI_COMMENT",
	"TOKEN_TYPE_NEWLINE",
	"TOKEN_TYPE_EQUALS",
	"TOKEN_TYPE_STRING",
	"TOKEN_TYPE_CLOSE_STRING",
	"TOKEN_TYPE_INTEGER",
	"TOKEN_TYPE_FLOAT",
	"TOKEN_TYPE_PLUS",
	"TOKEN_TYPE_MINUS",
	"TOKEN_TYPE_DIVIDE",
	"TOKEN_TYPE_MULTIPLY",
	"TOKEN_TYPE_COLON",
	"TOKEN_TYPE_SEMI_COLON",
	"TOKEN_TYPE_COMMA",
	"TOKEN_TYPE_OPEN_BRACKET",
	"TOKEN_TYPE_CLOSE_BRACKET",
	"TOKEN_TYPE_OPEN_BRACES",
	"TOKEN_TYPE_CLOSE_BRACES",
	"TOKEN_TYPE_AMPERSAND",
	"TOKEN_TYPE_GREATER_THAN",
	"TOKEN_TYPE_LESS_THAN",
	"TOKEN_TYPE_OR",
	"TOKEN_TYPE_TAB",
	"TOKEN_TYPE_CARRIAGE_RETURN",
	"TOKEN_TYPE_NONE",
	"TOKEN_TYPE_FUNCTION",
	"TOKEN_TYPE_INVOKE_FUNCTION",
	"TOKEN_TYPE_FUNCTION_DEF_START",
	"TOKEN_TYPE_FUNCTION_PARAM_END",
	"TOKEN_TYPE_FUNCTION_DEF_END",
	"TOKEN_TYPE_FUNCTION_RETURN",
	"TOKEN_TYPE_FOR_LOOP_START",
	"TOKEN_TYPE_FOR_LOOP_PARAM_END",
	"TOKEN_TYPE_FOR_LOOP_END",
	"TOKEN_TYPE_LOOP_BREAK",
	"TOKEN_TYPE_BOOLEAN",
	"TOKEN_TYPE_EQUALITY",
	"TOKEN_TYPE_INEQUALITY",
	"TOKEN_TYPE_LESS_THAN_OR_EQUALS",
	"TOKEN_TYPE_GREATER_THAN_OR_EQUALS",
	"TOKEN_TYPE_IF_START",
	"TOKEN_TYPE_IF_PARAM_END",
	"TOKEN_TYPE_IF_END",
	"TOKEN_TYPE_ELSE",
	"TOKEN_TYPE_ELIF_START",
	"TOKEN_TYPE_ELIF_PARAM_END",
	"TOKEN_TYPE_GET_ARRAY_START",
	"TOKEN_TYPE_GET_ARRAY_END",
	"TOKEN_TYPE_ARRAY",
	"TOKEN_TYPE_KEY_VALUE_PAIR",
	"TOKEN_TYPE_ASSOCIATIVE_ARRAY",
	"TOKEN_TYPE_WHILE_LOOP_START",
	"TOKEN_TYPE_WHILE_LOOP_PARAM_END",
	"TOKEN_TYPE_WHILE_LOOP_END",
}

//token object
type Token struct {
	Value               string
	FileName            string
	Context             string
	Array               []Token
	AssociativeArray    map[string]Token
	Type                int
	Line                int
	Column              int
	OtherInt            int
	Array_is_ref        bool
	Array_is_assoc      bool
	Array_ref_var_name  string
	Array_ref_index     int
	Array_ref_index_str string
}

type TokenArray struct {
	Tokens []Token
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
		fmt.Printf("Token Type: %s\n", TOKEN_TYPES_STRING[tokenArray[x].Type])
		fmt.Printf("Line #: %d\n", tokenArray[x].Line)
		fmt.Printf("Column #: %d\n", tokenArray[x].Column)

		if tokenArray[x].Type == TOKEN_TYPE_ARRAY {
			strVal := ""

			for x2 := 0; x2 < len(tokenArray[x].Array); x2++ {
				strVal = strVal + tokenArray[x].Array[x2].Value
				if (x2 + 1) != len(tokenArray[x].Array) {
					strVal = strVal + " , "
				}
			}

			fmt.Printf("Value: [ %s ]\n", strVal)
		} else {
			fmt.Printf("Value: %s\n", tokenArray[x].Value)
		}

		fmt.Printf("Context: %s\n", tokenArray[x].Context)
		fmt.Printf("OtherInt #: %d\n", tokenArray[x].OtherInt)
		fmt.Printf("====================================\n")
	}
}

//lexer object
type Lexer struct {
	FileName     string
	fileContents []string
}

func (lexer *Lexer) ReadSourceFile() error {

	file, err := os.Open(lexer.FileName)
	defer file.Close()

	if err == nil {
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			lexer.fileContents = append(lexer.fileContents, scanner.Text()+"\n")
		}

	} else {
		return errors.New("Error: no such file: '" + lexer.FileName + "'")
	}

	return nil
}

func (lexer *Lexer) ReadString(inputString string) {
	lexer.fileContents = append(lexer.fileContents, inputString)
}

func (lexer Lexer) GenerateToken() ([]Token, error) {
	var tokenArray []Token
	var cleanTokenArray []Token
	var finalTokenArray []Token
	tokenizerState := TOKENIZER_STATE_GET_WORD
	isTokenInit := false
	usePeriod := true
	withPeriod := false
	stringOpener := "\""
	var ignoreNext bool = false

	//read the file contents line by line
	for x := 0; x < len(lexer.fileContents); x++ {
		//read character by character
		for x2 := 0; x2 < len(lexer.fileContents[x]); x2++ {
			if ignoreNext {
				ignoreNext = false
				continue
			}
			currentChar := string(lexer.fileContents[x][x2])

		switch_label:
			switch tokenizerState {
			case TOKENIZER_STATE_GET_WORD:
				//get word
				if unicode.IsLetter([]rune(currentChar)[0]) || currentChar == "_" || currentChar == "!" {
					//alphabetic or underscore
					if len(tokenArray) > 0 {
						if tokenArray[len(tokenArray)-1].Type == TOKEN_TYPE_INTEGER {
							isTokenInit = false
						}
					}
					if !isTokenInit {
						setToken(true, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_IDENTIFIER, lexer.FileName, "") //init token
					}
					tokenArray[len(tokenArray)-1].Value += currentChar
				} else if currentChar == "\n" {
					//new line
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_NEWLINE, lexer.FileName, currentChar) //set token
				} else if currentChar == "\t" {
					//tab
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_TAB, lexer.FileName, currentChar) //set token
				} else if currentChar == " " {
					//space
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_SPACE, lexer.FileName, currentChar) //set token
				} else if currentChar == "=" {
					//equals
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_EQUALS, lexer.FileName, currentChar) //set token
					if (x2 + 1) < len(lexer.fileContents[x]) {
						if string(lexer.fileContents[x][x2+1]) == "=" {
							//equality operator instead of equals
							ignoreNext = true
							tokenArray[len(tokenArray)-1].Type = TOKEN_TYPE_EQUALITY
							tokenArray[len(tokenArray)-1].Value = "=="
						}
					}
				} else if currentChar == "+" {
					//plus
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_PLUS, lexer.FileName, currentChar) //set token
				} else if currentChar == "-" {
					//minus
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_MINUS, lexer.FileName, currentChar) //set token
				} else if currentChar == "/" {
					//divide
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_DIVIDE, lexer.FileName, currentChar) //set token
				} else if currentChar == "*" {
					//multiply
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_MULTIPLY, lexer.FileName, currentChar) //set token
				} else if currentChar == "." {
					//period
					usePeriod = true
					withPeriod = false
					if len(tokenArray) > 0 {
						if tokenArray[len(tokenArray)-1].Type == TOKEN_TYPE_INTEGER {
							tokenizerState = TOKENIZER_STATE_GET_FLOAT
							usePeriod = false
						}
					}
					if usePeriod {
						setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_PERIOD, lexer.FileName, currentChar) //set token
					}
				} else if currentChar == "(" {
					//open parenthesis
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_OPEN_PARENTHESIS, lexer.FileName, currentChar) //set token
				} else if currentChar == ")" {
					//close parenthesis
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_CLOSE_PARENTHESIS, lexer.FileName, currentChar) //set token
				} else if currentChar == "{" {
					//open bracket
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_OPEN_BRACKET, lexer.FileName, currentChar) //set token
				} else if currentChar == "}" {
					//close brackeT
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_CLOSE_BRACKET, lexer.FileName, currentChar) //set token
				} else if currentChar == "[" {
					//open braces
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_OPEN_BRACES, lexer.FileName, currentChar) //set token
				} else if currentChar == "]" {
					//close braces
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_CLOSE_BRACES, lexer.FileName, currentChar) //set token
				} else if currentChar == ":" {
					//colon
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_COLON, lexer.FileName, currentChar) //set token
				} else if currentChar == "," {
					//comma
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_COMMA, lexer.FileName, currentChar) //set token
				} else if currentChar == ";" {
					//semi colon
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_SEMI_COLON, lexer.FileName, currentChar) //set token
				} else if currentChar == "&" {
					//ampersand
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_AMPERSAND, lexer.FileName, currentChar) //set token
				} else if currentChar == ">" {
					//greater than
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_GREATER_THAN, lexer.FileName, currentChar) //set token
					if (x2 + 1) < len(lexer.fileContents[x]) {
						if string(lexer.fileContents[x][x2+1]) == "=" {
							//greater than or equals instead of greater than only
							ignoreNext = true
							tokenArray[len(tokenArray)-1].Type = TOKEN_TYPE_GREATER_THAN_OR_EQUALS
							tokenArray[len(tokenArray)-1].Value = ">="
						}
					}
				} else if currentChar == "<" {
					//less than
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_LESS_THAN, lexer.FileName, currentChar) //set token
					if (x2 + 1) < len(lexer.fileContents[x]) {
						if string(lexer.fileContents[x][x2+1]) == ">" {
							//inequality operator instead of less than
							ignoreNext = true
							tokenArray[len(tokenArray)-1].Type = TOKEN_TYPE_INEQUALITY
							tokenArray[len(tokenArray)-1].Value = "<>"
						} else if string(lexer.fileContents[x][x2+1]) == "=" {
							//less than or equals operator instead of less than
							ignoreNext = true
							tokenArray[len(tokenArray)-1].Type = TOKEN_TYPE_LESS_THAN_OR_EQUALS
							tokenArray[len(tokenArray)-1].Value = "<="
						}
					}
				} else if currentChar == "|" {
					//or
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_OR, lexer.FileName, currentChar) //set token
				} else if currentChar == "'" || currentChar == "\"" {
					//string
					stringOpener = currentChar
					setToken(true, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_STRING, lexer.FileName, "") //init token
					tokenizerState = TOKENIZER_STATE_GET_STRING
				} else if currentChar == "#" {
					//start of single line comment
					setToken(true, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_SINGLE_COMMENT, lexer.FileName, "") //init token
					tokenizerState = TOKENIZER_STATE_GET_SINGLE_COMMENT
				} else if currentChar == "^" {
					//start of multiline comment
					setToken(true, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_MULTI_COMMENT, lexer.FileName, "") //init token
					tokenizerState = TOKENIZER_STATE_GET_MULTI_COMMENT
				} else if unicode.IsDigit([]rune(currentChar)[0]) {
					//integer
					if !isTokenInit {
						setToken(true, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_INTEGER, lexer.FileName, "") //init token
					}
					tokenArray[len(tokenArray)-1].Value += currentChar
				} else if []rune(currentChar)[0] == 13 {
					//for windows
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_NEWLINE, lexer.FileName, "\n") //set token
				} else {
					return finalTokenArray, errors.New(SyntaxErrorMessage(x+1, x2+1, "Invalid token", lexer.FileName))
				}
			case TOKENIZER_STATE_GET_SINGLE_COMMENT:
				//get single comment
				if currentChar == "\n" || []rune(currentChar)[0] == 13 { //add char 13 for windows
					tokenizerState = TOKENIZER_STATE_GET_WORD
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_NEWLINE, lexer.FileName, "\n") //set token
				} else {
					tokenArray[len(tokenArray)-1].Value += currentChar
				}
			case TOKENIZER_STATE_GET_MULTI_COMMENT:
				//get multi comment
				if currentChar == "^" {
					tokenizerState = TOKENIZER_STATE_GET_WORD
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_CLOSE_MULTI_COMMENT, lexer.FileName, currentChar) //set token
				} else if []rune(currentChar)[0] == 13 {
					//for windows
					tokenArray[len(tokenArray)-1].Value += "\n"
				} else {
					tokenArray[len(tokenArray)-1].Value += currentChar
				}
			case TOKENIZER_STATE_GET_STRING:
				//get string
				if currentChar == stringOpener {
					tokenizerState = TOKENIZER_STATE_GET_WORD
					setToken(false, &tokenArray, &isTokenInit, x+1, x2+1, TOKEN_TYPE_CLOSE_STRING, lexer.FileName, currentChar) //set token
				} else if []rune(currentChar)[0] == 13 {
					//for windows
					tokenArray[len(tokenArray)-1].Value += "\n"
				} else {
					tokenArray[len(tokenArray)-1].Value += currentChar
				}
			case TOKENIZER_STATE_GET_FLOAT:
				//get float
				if unicode.IsDigit([]rune(currentChar)[0]) {
					if !withPeriod {
						withPeriod = true
						tokenArray[len(tokenArray)-1].Value += "."
						tokenArray[len(tokenArray)-1].Type = TOKEN_TYPE_FLOAT
					}

					tokenArray[len(tokenArray)-1].Value += currentChar
				} else {
					if !withPeriod {
						setToken(false, &tokenArray, &isTokenInit, x+1, x2, TOKEN_TYPE_PERIOD, lexer.FileName, ".") //set token
					}
					tokenizerState = TOKENIZER_STATE_GET_WORD
					isTokenInit = false
					goto switch_label
				}
			default:
				continue
			}
		}
	}

	//1st token cleanup
	//TODO: Try to eliminate this part
	for x := 0; x < len(tokenArray); x++ {
		if tokenArray[x].Type == TOKEN_TYPE_SPACE || tokenArray[x].Type == TOKEN_TYPE_SINGLE_COMMENT || tokenArray[x].Type == TOKEN_TYPE_MULTI_COMMENT || tokenArray[x].Type == TOKEN_TYPE_CLOSE_MULTI_COMMENT || tokenArray[x].Type == TOKEN_TYPE_CLOSE_STRING || tokenArray[x].Type == TOKEN_TYPE_TAB {
			//ignore space, tab, comments etc...
			continue
		}
		cleanTokenArray = append(cleanTokenArray, tokenArray[x])
	}

	//2nd token cleanup
	var ignoreOpen bool = false
	var isOpenP bool = false
	var isOpenB bool = false
	var isFunctionDef bool = false
	var isForLoop bool = false
	var isWhileLoop bool = false
	var isForIf bool = false
	var isForEf bool = false
	var openFunctionCount int = 0
	var f_count int = 0
	var op_count map[string]int
	var ob_count map[string]int
	op_count = make(map[string]int)
	ob_count = make(map[string]int)

	var contextName = []string{"main_context"}

	for x := 0; x < len(cleanTokenArray); x++ {
		if ignoreOpen {
			//ignore open parenthesis if the last token is a function
			//or ignore open braces if the last token is a variable that accessing its index
			ignoreOpen = false
			continue
		}
		if cleanTokenArray[x].Type == TOKEN_TYPE_OPEN_PARENTHESIS {
			op_count[contextName[len(contextName)-1]] += 1
		}
		if cleanTokenArray[x].Type == TOKEN_TYPE_OPEN_BRACES {
			ob_count[contextName[len(contextName)-1]] += 1
		}
		if cleanTokenArray[x].Type == TOKEN_TYPE_CLOSE_BRACES {
			if ob_count[contextName[len(contextName)-1]] > 0 {
				ob_count[contextName[len(contextName)-1]] -= 1
			} else {
				cleanTokenArray[x].Type = TOKEN_TYPE_GET_ARRAY_END
				cleanTokenArray[x].Context = contextName[len(contextName)-1]
				contextName = contextName[:len(contextName)-1]
			}
		}
		if cleanTokenArray[x].Type == TOKEN_TYPE_CLOSE_PARENTHESIS {
			if op_count[contextName[len(contextName)-1]] > 0 {
				op_count[contextName[len(contextName)-1]] -= 1
			} else {
				if f_count > 0 {
					f_count -= 1

					if isFunctionDef {
						isFunctionDef = false
						cleanTokenArray[x].Type = TOKEN_TYPE_FUNCTION_PARAM_END
					} else {
						cleanTokenArray[x].Type = TOKEN_TYPE_INVOKE_FUNCTION
						cleanTokenArray[x].Context = contextName[len(contextName)-1]
						contextName = contextName[:len(contextName)-1]
					}
				} else {
					if isForLoop {
						isForLoop = false
						cleanTokenArray[x].Type = TOKEN_TYPE_FOR_LOOP_PARAM_END
						cleanTokenArray[x].Context = contextName[len(contextName)-1]
						contextName = contextName[:len(contextName)-1]
					}
					if isWhileLoop {
						isWhileLoop = false
						cleanTokenArray[x].Type = TOKEN_TYPE_WHILE_LOOP_PARAM_END
						cleanTokenArray[x].Context = contextName[len(contextName)-1]
						contextName = contextName[:len(contextName)-1]
					}
					if isForIf {
						isForIf = false
						cleanTokenArray[x].Type = TOKEN_TYPE_IF_PARAM_END
						cleanTokenArray[x].Context = contextName[len(contextName)-1]
						contextName = contextName[:len(contextName)-1]
					}
					if isForEf {
						isForEf = false
						cleanTokenArray[x].Type = TOKEN_TYPE_ELIF_PARAM_END
						cleanTokenArray[x].Context = contextName[len(contextName)-1]
						contextName = contextName[:len(contextName)-1]
					}
				}
			}
		}
		if cleanTokenArray[x].Type == TOKEN_TYPE_IDENTIFIER {
			if IsReservedWord(cleanTokenArray[x].Value) {
				//Convert identifier to keyword if existing in reserved words
				cleanTokenArray[x].Type = TOKEN_TYPE_KEYWORD
				if cleanTokenArray[x].Value == "fd" {
					//function definition
					openFunctionCount += 1
					continue
				}
				if cleanTokenArray[x].Value == "df" {
					//end of function definition
					cleanTokenArray[x].Type = TOKEN_TYPE_FUNCTION_DEF_END
					openFunctionCount -= 1
				}
				if cleanTokenArray[x].Value == "rtn" {
					//function return
					cleanTokenArray[x].Type = TOKEN_TYPE_FUNCTION_RETURN
				}
				if cleanTokenArray[x].Value == "brk" {
					//loop break
					cleanTokenArray[x].Type = TOKEN_TYPE_LOOP_BREAK
				}
				if cleanTokenArray[x].Value == "fl" {
					//for loop
					isForLoop = true
					cleanTokenArray[x].Type = TOKEN_TYPE_FOR_LOOP_START
					ignoreOpen = true

					thisSuffix := strconv.Itoa(cleanTokenArray[x].Column)
					contextName = append(contextName, "fl_"+thisSuffix)

					if (x + 1) <= len(cleanTokenArray)-1 {
						if cleanTokenArray[x+1].Type != TOKEN_TYPE_OPEN_PARENTHESIS {
							return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x+1].Line, cleanTokenArray[x+1].Column, "Unexpected token '"+cleanTokenArray[x+1].Value+"'", cleanTokenArray[x+1].FileName))
						}
					} else {
						return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x].Line, cleanTokenArray[x].Column, "Unfinished statement", cleanTokenArray[x].FileName))
					}
					//continue
				}
				if cleanTokenArray[x].Value == "lf" {
					cleanTokenArray[x].Type = TOKEN_TYPE_FOR_LOOP_END
				}

				if cleanTokenArray[x].Value == "wl" {
					//while loop
					isWhileLoop = true
					cleanTokenArray[x].Type = TOKEN_TYPE_WHILE_LOOP_START
					ignoreOpen = true

					thisSuffix := strconv.Itoa(cleanTokenArray[x].Column)
					contextName = append(contextName, "wl_"+thisSuffix)

					if (x + 1) <= len(cleanTokenArray)-1 {
						if cleanTokenArray[x+1].Type != TOKEN_TYPE_OPEN_PARENTHESIS {
							return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x+1].Line, cleanTokenArray[x+1].Column, "Unexpected token '"+cleanTokenArray[x+1].Value+"'", cleanTokenArray[x+1].FileName))
						}
					} else {
						return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x].Line, cleanTokenArray[x].Column, "Unfinished statement", cleanTokenArray[x].FileName))
					}
				}
				if cleanTokenArray[x].Value == "lw" {
					cleanTokenArray[x].Type = TOKEN_TYPE_WHILE_LOOP_END
				}

				if cleanTokenArray[x].Value == "if" || cleanTokenArray[x].Value == "ef" {
					//if or ef statement
					thisSuffix := strconv.Itoa(cleanTokenArray[x].Column)

					if cleanTokenArray[x].Value == "if" {
						//if statement
						isForIf = true
						cleanTokenArray[x].Type = TOKEN_TYPE_IF_START
						contextName = append(contextName, "if_"+thisSuffix)
					} else {
						//ef statement
						isForEf = true
						cleanTokenArray[x].Type = TOKEN_TYPE_ELIF_START
						contextName = append(contextName, "ef_"+thisSuffix)
					}

					ignoreOpen = true
					if (x + 1) <= len(cleanTokenArray)-1 {
						if cleanTokenArray[x+1].Type != TOKEN_TYPE_OPEN_PARENTHESIS {
							return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x+1].Line, cleanTokenArray[x+1].Column, "Unexpected token '"+cleanTokenArray[x+1].Value+"'", cleanTokenArray[x+1].FileName))
						}
					} else {
						return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x].Line, cleanTokenArray[x].Column, "Unfinished statement", cleanTokenArray[x].FileName))
					}
					//continue
				}
				if cleanTokenArray[x].Value == "fi" {
					cleanTokenArray[x].Type = TOKEN_TYPE_IF_END
				}
				if cleanTokenArray[x].Value == "el" {
					cleanTokenArray[x].Type = TOKEN_TYPE_ELSE
				}
			} else {
				if (x + 1) <= len(cleanTokenArray)-1 {
					isOpenP = false
					//Check if the next token is '(', if yes then it's a function call
					if cleanTokenArray[x+1].Type == TOKEN_TYPE_OPEN_PARENTHESIS {
						isOpenP = true
					}
					if isOpenP {
						//set to function call
						cleanTokenArray[x].Type = TOKEN_TYPE_FUNCTION
						ignoreOpen = true
						f_count += 1
						if (x - 1) >= 0 {
							if cleanTokenArray[x-1].Value == "fd" {
								//set to function definition
								cleanTokenArray[x].Type = TOKEN_TYPE_FUNCTION_DEF_START
								isFunctionDef = true
								if openFunctionCount > 1 {
									//if it's already true then
									//you define a function inside a function
									//but it's prohibited so raise an error
									return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x].Line, cleanTokenArray[x].Column, "You cannot define a function inside a function", lexer.FileName))
								}
							}
						}
						if cleanTokenArray[x].Type == TOKEN_TYPE_FUNCTION {
							thisSuffix := strconv.Itoa(cleanTokenArray[x].Column)
							contextName = append(contextName, cleanTokenArray[x].Value+"_"+thisSuffix)
						}
					}
					isOpenB = false
					//Check if the next token is '[', if yes then it's an array getter
					if cleanTokenArray[x+1].Type == TOKEN_TYPE_OPEN_BRACES {
						isOpenB = true
					}
					if isOpenB {
						cleanTokenArray[x].Type = TOKEN_TYPE_GET_ARRAY_START
						ignoreOpen = true
						thisSuffix := strconv.Itoa(cleanTokenArray[x].Column)
						contextName = append(contextName, "array_get"+thisSuffix)
					}
				}
			}
		} else if cleanTokenArray[x].Type == TOKEN_TYPE_STRING {
			if !((x + 1) < len(cleanTokenArray)) {
				return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x].Line, cleanTokenArray[x].Column, "Expected closing of string", lexer.FileName))
			}
		} else if cleanTokenArray[x].Type == TOKEN_TYPE_MULTI_COMMENT {
			if !((x + 1) < len(cleanTokenArray)) {
				return finalTokenArray, errors.New(SyntaxErrorMessage(cleanTokenArray[x].Line, cleanTokenArray[x].Column, "Expected closing of multi line comment", lexer.FileName))
			}
		}

		//set context
		if cleanTokenArray[x].Context == "" {
			cleanTokenArray[x].Context = contextName[len(contextName)-1]
		}

		if x != 0 && (cleanTokenArray[x].Type == TOKEN_TYPE_FLOAT || cleanTokenArray[x].Type == TOKEN_TYPE_INTEGER) {
			if cleanTokenArray[x-1].Type == TOKEN_TYPE_MINUS {

				if (x - 2) >= 0 {
					if cleanTokenArray[x-2].Type != TOKEN_TYPE_FLOAT && cleanTokenArray[x-2].Type != TOKEN_TYPE_INTEGER && cleanTokenArray[x-2].Type != TOKEN_TYPE_IDENTIFIER && cleanTokenArray[x-2].Type != TOKEN_TYPE_CLOSE_PARENTHESIS && cleanTokenArray[x-2].Type != TOKEN_TYPE_INVOKE_FUNCTION { //added && cleanTokenArray[x-2].Type != TOKEN_TYPE_CLOSE_PARENTHESIS
						//set to negative number
						finalTokenArray[len(finalTokenArray)-1].Value += cleanTokenArray[x].Value
						finalTokenArray[len(finalTokenArray)-1].Type = cleanTokenArray[x].Type
					} else {
						finalTokenArray = append(finalTokenArray, cleanTokenArray[x])
					}
				} else {
					finalTokenArray = append(finalTokenArray, cleanTokenArray[x])
				}

			} else {
				finalTokenArray = append(finalTokenArray, cleanTokenArray[x])
			}
		} else {
			finalTokenArray = append(finalTokenArray, cleanTokenArray[x])
		}
	}

	//3rd token cleanup (for open and close braces)
	var contextNameBrace = []string{"main_context"}
	var contextToReplaceBrace string = ""
	var contextNameBracket = []string{"main_context"}
	var contextToReplaceBracket string = ""
	for x := 0; x < len(finalTokenArray); x++ {

		//braces
		if finalTokenArray[x].Type == TOKEN_TYPE_OPEN_BRACES {
			thisSuffix := strconv.Itoa(finalTokenArray[x].Column)
			contextNameBrace = append(contextNameBrace, "ob_"+thisSuffix)
			contextToReplaceBrace = finalTokenArray[x].Context
		}

		if finalTokenArray[x].Type == TOKEN_TYPE_CLOSE_BRACES {
			finalTokenArray[x].Context = contextNameBrace[len(contextNameBrace)-1]
			contextNameBrace = contextNameBrace[:len(contextNameBrace)-1]
		}

		if finalTokenArray[x].Context == contextToReplaceBrace {
			finalTokenArray[x].Context = contextNameBrace[len(contextNameBrace)-1]
		}

		//bracket
		if finalTokenArray[x].Type == TOKEN_TYPE_OPEN_BRACKET {
			thisSuffix := strconv.Itoa(finalTokenArray[x].Column)
			contextNameBracket = append(contextNameBracket, "obr_"+thisSuffix)
			contextToReplaceBracket = finalTokenArray[x].Context
		}

		if finalTokenArray[x].Type == TOKEN_TYPE_CLOSE_BRACKET {
			finalTokenArray[x].Context = contextNameBracket[len(contextNameBracket)-1]
			contextNameBracket = contextNameBracket[:len(contextNameBracket)-1]
		}

		if finalTokenArray[x].Context == contextToReplaceBracket {
			finalTokenArray[x].Context = contextNameBracket[len(contextNameBracket)-1]
		}

	}

	return finalTokenArray, nil
}
