package arithmetic

import (
	"false-vm/input"
	"strconv"
	"unicode"
)

type TokenInput struct {
	Input input.RuneInput
}

const (
	Plus     rune = '+'
	Minus    rune = '-'
	Multiply rune = '*'
	Divide   rune = '/'

	Open  rune = '('
	Close rune = ')'
)

func (ti *TokenInput) IsOperand() bool {
	b := ti.Input.Peek()
	return unicode.IsDigit(b)
}

func (ti *TokenInput) ReadOperand() (int, error) {
	b := make([]rune, 0)
	for ti.IsOperand() && !ti.Input.Eof() {
		b = append(b, ti.Input.Next())
	}
	s := string(b)
	v, e := strconv.Atoi(s)
	if e != nil {
		ti.Input.Croak("invalid integer format: " + s)
		return 0, e
	}
	return v, nil
}

func (ti *TokenInput) IsOperator() bool {
	c := ti.Input.Peek()
	switch c {
	case Plus, Minus,
		Multiply, Divide:
		return true
	default:
		return false
	}
}

func (ti *TokenInput) ReadOperator() rune {
	return ti.Input.Next()
}

func (ti *TokenInput) IsCommaStat() bool {
	return ti.Input.Peek() == Open
}

func (ti *TokenInput) IsCommaEnd() bool {
	return ti.Input.Peek() == Close
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
