package main

type Execute func()

type Function struct {
	Name string
	IsSystem bool
	Tokens []Token
	Run Execute
}