package models

type Token struct {
	TokenType string
	Value     string
	UserId    int
	Exp       int
}
