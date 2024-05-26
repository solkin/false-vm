package main

import (
	"errors"
)

type Parser struct {
}

var InstrMap = map[rune]int{
	DUP:        InstrDup,
	DROP:       InstrDrop,
	SWAP:       InstrSwap,
	ROT:        InstrRot,
	PICK:       InstrPick,
	PLUS:       InstrPlus,
	MINUS:      InstrMinus,
	MULTIPLY:   InstrMultiply,
	DIVIDE:     InstrDivide,
	NEGATIVE:   InstrNegative,
	AND:        InstrAnd,
	OR:         InstrOr,
	NOT:        InstrNot,
	GREATER:    InstrMore,
	EQUALS:     InstrEquals,
	READ_CHAR:  InstrReadChar,
	WRITE_CHAR: InstrWriteChar,
	WRITE_INT:  InstrWriteInt,
	FLUSH:      InstrFlush,
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(str string) ([]byte, error) {
	input := TokenInput{Input: StringInput{Str: str}}

	w := NewBytecodeWriter()

	vars := make(map[string]int)
	for !input.Eof() {
		if input.IsInt() {
			if v, err := input.ReadInt(); err == nil {
				w.WritePush(v)
			} else {
				return nil, err
			}
		} else if input.IsCharCode() {
			if v, err := input.ReadCharCode(); err == nil {
				w.WritePush(int(v))
			} else {
				return nil, err
			}
		} else if input.IsVar() {
			if v, m, err := input.ReadVar(); err == nil {
				var addr int
				ok := false
				if addr, ok = vars[v]; !ok {
					// New variable
					addr = w.WriteVar(0)
					vars[v] = addr
				}
				switch m {
				case STORE_VAR:
					w.WriteStore(addr)
					break
				case FETCH_VAR:
					w.WriteFetch(addr)
					break
				}
			} else {
				return nil, err
			}
		} else if input.IsSubStart() {
			input.SkipSubStart()
			w.SubCreate()
		} else if input.IsSubEnd() {
			input.SkipSubEnd()
			if err := w.SubReturn(); err != nil {
				input.Input.Croak(err.Error())
				return nil, err
			}
		} else if input.IsSubCall() {
			input.SkipSubCall()
			w.WriteCall()
		} else if input.IsIf() {
			input.SkipIf()
			w.WriteCallIf()
		} else if input.IsWhile() {
			input.SkipWhile()
			// Reserved condition and body addresses
			ca := w.WriteVar(0)
			ba := w.WriteVar(0)
			// Take body and condition addresses down from stack
			w.WriteStore(ba)
			w.WriteStore(ca)
			// Call body (by goto)
			w.BlockCreate()
			w.WriteFetch(ba)
			w.WriteCall()
			bca, err := w.BlockSkip()
			if err != nil {
				return nil, err
			}
			// Call condition
			w.WriteFetch(ca)
			w.WriteCall()
			// Push to stack call body address
			w.WritePush(bca)
			// Call condition
			w.WriteGotoIf()
		} else if input.IsCommand() {
			if ic, err := input.ReadCommand(); err == nil {
				if cmd, ok := InstrMap[ic]; ok {
					w.WriteCommand(cmd)
				} else {
					err := errors.New("invalid command")
					input.Input.Croak(err.Error())
					return nil, err
				}
			} else {
				return nil, err
			}
		} else if input.IsString() {
			if s, err := input.ReadString(); err == nil {
				w.WriteString(s)
			} else {
				input.Input.Croak(err.Error())
				return nil, err
			}
		} else if input.IsWhitespace() {
			input.SkipWhitespace()
		}
	}
	w.WriteEnd()
	return w.Bytes(), nil
}
