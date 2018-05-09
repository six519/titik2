package lexer

import (
	"os"
	"bufio"
	"errors"
	//"fmt"
)

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

	//read the file contents line by line
	for x := 0; x < len(lexer.fileContents); x++ {
		/*
		fmt.Println(lexer.fileContents[x])
		fmt.Printf("Line #%d \n", x + 1)
		*/
	}

	return tokenArray, nil
}
