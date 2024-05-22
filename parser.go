package main

import (
	"errors"
)

type Parser struct {
}

const (
	InstrPush = 1

	InstrDup  = 2
	InstrDrop = 3
	InstrSwap = 4
	InstrRot  = 5
	InstrPick = 6

	InstrPlus     = 7
	InstrMinus    = 8
	InstrMultiply = 9
	InstrDivide   = 10
	InstrNegative = 11
	InstrAnd      = 12
	InstrOr       = 13
	InstrNot      = 14

	InstrMore   = 15
	InstrEquals = 16

	InstrReadChar  = 17
	InstrWriteChar = 18
	InstrWriteInt  = 19
	InstrWriteStr  = 20
	InstrFlush     = 21

	InstrStore = 22
	InstrFetch = 23

	InstrCall   = 24
	InstrReturn = 25

	InstrGoto   = 26
	InstrGotoIf = 27

	InstrEnd = 28
)

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

func (p *Parser) Parse(str string) ([]int, error) {
	input := TokenInput{Input: StringInput{Str: str}}
	out := NewDynArray()

	vars := make(map[string]int)
	subs := NewStack(
		make([]int, 1000),
		0,
		1000,
	)
	for !input.Eof() {
		if input.IsInt() {
			if v, err := input.ReadInt(); err == nil {
				out.Add(InstrPush)
				out.Add(v)
			} else {
				return nil, err
			}
		} else if input.IsCharCode() {
			if v, err := input.ReadCharCode(); err == nil {
				out.Add(InstrPush)
				out.Add(v)
			} else {
				return nil, err
			}
		} else if input.IsVar() {
			if v, m, err := input.ReadVar(); err == nil {
				addr := 0
				ok := false
				if addr, ok = vars[v]; !ok {
					// New variable
					out.Add(InstrGoto)
					out.Add(out.Size + 2) // Allocate
					addr = out.Size
					vars[v] = addr
					out.Add(0) // Reserved variable value
				}

				switch m {
				case STORE_VAR:
					out.Add(InstrStore)
					out.Add(addr)
					break
				case FETCH_VAR:
					out.Add(InstrFetch)
					out.Add(addr)
					break
				}
			} else {
				return nil, err
			}
		} else if input.IsSubStart() {
			input.SkipSubStart()
			if err := subs.Push(out.Size + 1); err != nil {
				input.Input.Croak(err.Error())
				return nil, err
			}
			out.Add(InstrGoto)
			out.Add(0) // Goto address will be set after sub end
		} else if input.IsSubEnd() {
			input.SkipSubEnd()
			out.Add(InstrReturn)
			if v, err := subs.Pop(); err == nil {
				out.Data[v] = out.Size
				out.Add(InstrPush)
				out.Add(v + 1) // Sub start point
			} else {
				input.Input.Croak(err.Error())
				return nil, err
			}
		} else if input.IsSubCall() {
			input.SkipSubCall()
			out.Add(InstrCall)
		} else if input.IsIf() {
			input.SkipIf()
			out.Add(InstrGoto)
			out.Add(out.Size + 2)
			addr := out.Size
			out.Add(0) // Conditional lambda address reserve
			out.Add(InstrStore)
			out.Add(addr)
			// Add goto condition body address to stack
			out.Add(InstrPush)
			out.Add(out.Size + 4)

			out.Add(InstrGotoIf)

			// Skip condition goto body
			out.Add(InstrGoto)
			out.Add(out.Size + 4)
			// Call conditional lambda
			out.Add(InstrFetch)
			out.Add(addr)
			out.Add(InstrCall)
		} else if input.IsWhile() {
			input.SkipWhile()
			out.Add(InstrGoto)
			out.Add(out.Size + 3)
			ca := out.Size
			out.Add(0) // Reserved condition address
			ba := out.Size
			out.Add(0) // Reserved body address
			// Take body and condition addresses down from stack
			out.Add(InstrStore)
			out.Add(ba)
			out.Add(InstrStore)
			out.Add(ca)
			// Skip body
			out.Add(InstrGoto)
			out.Add(out.Size + 4)
			// Call body (by goto) // TODO: this is not sub!!!
			bca := out.Size
			out.Add(InstrFetch)
			out.Add(ba)
			out.Add(InstrCall)
			// Call condition
			out.Add(InstrFetch)
			out.Add(ca)
			out.Add(InstrCall)
			// Push to stack call body address
			out.Add(InstrPush)
			out.Add(bca)
			// Call condition
			out.Add(InstrGotoIf)
		} else if input.IsCommand() {
			if ic, err := input.ReadCommand(); err == nil {
				if cmd, ok := InstrMap[ic]; ok {
					out.Add(cmd)
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
				out.Add(InstrWriteStr)
				out.Add(len(s))
				for _, c := range s {
					out.Add(int(c))
				}
			} else {
				input.Input.Croak(err.Error())
				return nil, err
			}
		} else if input.IsWhitespace() {
			input.SkipWhitespace()
		}
	}
	out.Add(InstrEnd)
	return out.ToArray(), nil
}
