package arithmetic

import (
	"errors"
	"false-vm/input"
	"false-vm/vm"
	"io"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

type operator struct {
	Weight int
	Value  rune
}

func (p *Parser) Parse(r io.Reader, w io.Writer) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	ti := TokenInput{Input: &input.StringInput{Str: string(data)}}

	bc := vm.NewBytecodeWriter()

	priority := make(map[rune]int)
	priority[Open] = 0
	priority[Plus] = 1
	priority[Minus] = 1
	priority[Multiply] = 2
	priority[Divide] = 2

	s := vm.NewStack()

	for !ti.Eof() {
		if ti.IsOperand() {
			v, err := ti.ReadOperand()
			if err != nil {
				return err
			}
			bc.WritePush(v)
		} else if ti.IsOperator() {
			o := ti.ReadOperator()
			rw, ok := priority[o]
			if !ok {
				return errors.New("unknown operator")
			}
			ro := operator{
				Weight: rw,
				Value:  o,
			}
			if s.Len() != 0 && s.Peek().(operator).Weight > rw {
				for s.Peek().(operator).Weight > rw && s.Peek().(operator).Value != Open {
					popOperatorStack(s, bc)
				}
			}
			s.Push(ro)
		} else if ti.IsCommaStat() {
			ti.Skip()
			s.Push(operator{
				Weight: priority[Open],
				Value:  Open,
			})
		} else if ti.IsCommaEnd() {
			ti.Skip()
			for s.Peek().(operator).Value != Open {
				popOperatorStack(s, bc)
			}
		} else {
			ti.Skip()
		}
	}
	for s.Len() > 0 {
		popOperatorStack(s, bc)
	}
	bc.WriteCommand(vm.InstrWriteInt)
	bc.WriteEnd()
	_, err = bc.WriteTo(w)
	return err
}

func popOperatorStack(s *vm.Stack, bc *vm.BytecodeWriter) {
	so := s.Pop().(operator)
	switch so.Value {
	case Plus:
		bc.WriteCommand(vm.InstrPlus)
		break
	case Minus:
		bc.WriteCommand(vm.InstrMinus)
		break
	case Multiply:
		bc.WriteCommand(vm.InstrMultiply)
		break
	case Divide:
		bc.WriteCommand(vm.InstrDivide)
		break
	}
}
