package bf

import (
	"errors"
	"false-vm/input"
)

type TokenInput struct {
	Input input.RuneInput
}

const (
	NEXT rune = '>'
	PREV rune = '<'

	PLUS  rune = '+'
	MINUS rune = '-'

	IN  rune = ','
	OUT rune = '.'

	SUB    rune = '['
	RETURN rune = ']'
)

func (ti *TokenInput) IsCharCode() bool {
	return ti.Input.Peek() == '\''
}

func (ti *TokenInput) ReadCharCode() (rune, error) {
	if ti.IsCharCode() {
		ti.Input.Next()
		return ti.Input.Next(), nil
	}
	err := errors.New("not a char")
	ti.Input.Croak(err.Error())
	return 0, err
}

func (ti *TokenInput) IsCommand() bool {
	c := ti.Input.Peek()
	switch c {
	case NEXT, PREV,
		PLUS, MINUS,
		IN, OUT,
		SUB, RETURN:
		return true
	default:
		return false
	}
}

func (ti *TokenInput) Skip() {
	ti.Input.Next()
}

func (ti *TokenInput) Next() rune {
	return ti.Input.Next()
}

func (ti *TokenInput) Eof() bool {
	return ti.Input.Eof()
}
