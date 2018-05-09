package lexer

import (
	"os"
	"bufio"
	"errors"
)

type Lexer struct {
	FileName string
	fileContents []string
}

func (lexer Lexer) ReadSourceFile() error {

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

	return tokenArray, nil
}
