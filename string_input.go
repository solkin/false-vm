package main

import (
	"log"
	"unicode/utf8"
)

type StringInput struct {
	Str  string
	pos  int
	line int
	col  int
}

func (s *StringInput) Peek() rune {
	c, _ := s.getChar(s.pos)
	return c
}

func (s *StringInput) Next() rune {
	c, w := s.getChar(s.pos)
	s.pos += w
	if c == '\n' || c == '\r' {
		s.line++
		s.col = 0
	} else {
		s.col++
	}
	return c
}

func (s *StringInput) Eof() bool {
	return s.Peek() == 0
}

func (s *StringInput) Croak(msg string) {
	log.Printf("%s (%d:%d)\n", msg, s.line, s.col)
}

func (s *StringInput) getChar(pos int) (rune, int) {
	if pos >= 0 && pos < len(s.Str) {
		c, w := utf8.DecodeRuneInString(s.Str[pos:])
		return c, w
	}
	return 0, 0
}
