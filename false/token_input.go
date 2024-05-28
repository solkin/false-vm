package false

import (
	"errors"
	"sandbox-vm/input"
	"strconv"
	"unicode"
)

type TokenInput struct {
	Input input.RuneInput
}

const (
	DUP  rune = '$'
	DROP rune = '%'
	SWAP rune = '\\'
	ROT  rune = '@'
	PICK rune = 'ø'

	PLUS     rune = '+'
	MINUS    rune = '-'
	MULTIPLY rune = '*'
	DIVIDE   rune = '/'
	NEGATIVE rune = '_'
	AND      rune = '&'
	OR       rune = '|'
	NOT      rune = '~'

	GREATER rune = '>'
	EQUALS  rune = '='

	READ_CHAR  rune = '^'
	WRITE_CHAR rune = ','
	WRITE_INT  rune = '.'
	FLUSH      rune = 'ß'

	STORE_VAR rune = ':'
	FETCH_VAR rune = ';'
)

func (ti *TokenInput) IsInt() bool {
	b := ti.Input.Peek()
	return unicode.IsDigit(b)
}

func (ti *TokenInput) ReadInt() (int, error) {
	b := make([]rune, 0)
	for ti.IsInt() && !ti.Input.Eof() {
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
	case DUP, DROP, SWAP, ROT, PICK,
		PLUS, MINUS, MULTIPLY, DIVIDE, NEGATIVE, AND, OR, NOT,
		GREATER, EQUALS,
		READ_CHAR, WRITE_CHAR, WRITE_INT, FLUSH:
		return true
	default:
		return false
	}
}

func (ti *TokenInput) ReadCommand() (rune, error) {
	if ti.IsCommand() {
		return ti.Input.Next(), nil
	}
	err := errors.New("not an arithmetic operation")
	ti.Input.Croak(err.Error())
	return 0, err
}

func (ti *TokenInput) IsString() bool {
	return ti.Input.Peek() == '"'
}

func (ti *TokenInput) ReadString() (string, error) {
	b := make([]rune, 0)
	if ti.IsString() {
		ti.Input.Next()
		escaped := false
		for !ti.Input.Eof() {
			c := ti.Input.Next()
			if escaped {
				escaped = false
			} else if c == '\\' {
				escaped = true
				continue
			}
			b = append(b, c)
			if ti.IsString() {
				ti.Input.Next()
				break
			}
		}
	} else {
		err := errors.New("not a string")
		ti.Input.Croak(err.Error())
		return "", err
	}
	return string(b), nil
}

func (ti *TokenInput) IsVar() bool {
	b := ti.Input.Peek()
	return isLetter(b)
}

func isLetter(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func (ti *TokenInput) ReadVar() (string, rune, error) {
	b := make([]rune, 1)
	var mode rune
	if ti.IsVar() {
		b[0] = ti.Input.Next()
		mode = ti.NextSkipWhitespaces()
		switch mode {
		case STORE_VAR, FETCH_VAR:
			break
		default:
			err := errors.New("invalid var mode")
			ti.Input.Croak(err.Error())
			return "", 0, err
		}
	} else {
		err := errors.New("not a variable")
		ti.Input.Croak(err.Error())
		return "", 0, err
	}
	return string(b), mode, nil
}

func (ti *TokenInput) IsSubStart() bool {
	b := ti.Input.Peek()
	return b == '['
}

func (ti *TokenInput) SkipSubStart() {
	if ti.IsSubStart() {
		ti.Input.Next()
	}
}

func (ti *TokenInput) IsSubEnd() bool {
	b := ti.Input.Peek()
	return b == ']'
}

func (ti *TokenInput) SkipSubEnd() {
	if ti.IsSubEnd() {
		ti.Input.Next()
	}
}

func (ti *TokenInput) IsSubCall() bool {
	b := ti.Input.Peek()
	return b == '!'
}

func (ti *TokenInput) SkipSubCall() {
	if ti.IsSubCall() {
		ti.Input.Next()
	}
}

func (ti *TokenInput) IsIf() bool {
	b := ti.Input.Peek()
	return b == '?'
}

func (ti *TokenInput) SkipIf() {
	if ti.IsIf() {
		ti.Input.Next()
	}
}

func (ti *TokenInput) IsWhile() bool {
	b := ti.Input.Peek()
	return b == '#'
}

func (ti *TokenInput) SkipWhile() {
	if ti.IsWhile() {
		ti.Input.Next()
	}
}

func (ti *TokenInput) IsCommentStart() bool {
	return ti.Input.Peek() == '{'
}

func (ti *TokenInput) IsCommentEnd() bool {
	return ti.Input.Peek() == '{'
}

func (ti *TokenInput) ReadComment() (string, error) {
	b := make([]rune, 0)
	if ti.IsCommentStart() {
		ti.Input.Next()
		escaped := false
		for !ti.Input.Eof() {
			c := ti.Input.Next()
			if escaped {
				escaped = false
			} else if c == '\\' {
				escaped = true
				continue
			}
			b = append(b, c)
			if ti.IsCommentEnd() {
				ti.Input.Next()
				break
			}
		}
	} else {
		err := errors.New("not a comment")
		ti.Input.Croak(err.Error())
		return "", err
	}
	return string(b), nil
}

func (ti *TokenInput) IsWhitespace() bool {
	c := ti.Input.Peek()
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func (ti *TokenInput) SkipWhitespace() {
	if ti.IsWhitespace() {
		ti.Input.Next()
	}
}

func (ti *TokenInput) NextSkipWhitespaces() rune {
	if ti.IsWhitespace() {
		ti.Input.Next()
	}
	return ti.Input.Next()
}

func (ti *TokenInput) Eof() bool {
	return ti.Input.Eof()
}
