package input

type RuneInput interface {
	Peek() rune
	Next() rune
	Eof() bool
	Croak(msg string)
}
