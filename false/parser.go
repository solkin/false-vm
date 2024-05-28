package false

import (
	"errors"
	"io"
	"sandbox-vm/input"
	"sandbox-vm/vm"
)

type Parser struct {
}

var InstrMap = map[rune]int{
	DUP:        vm.InstrDup,
	DROP:       vm.InstrDrop,
	SWAP:       vm.InstrSwap,
	ROT:        vm.InstrRot,
	PICK:       vm.InstrPick,
	PLUS:       vm.InstrPlus,
	MINUS:      vm.InstrMinus,
	MULTIPLY:   vm.InstrMultiply,
	DIVIDE:     vm.InstrDivide,
	NEGATIVE:   vm.InstrNegative,
	AND:        vm.InstrAnd,
	OR:         vm.InstrOr,
	NOT:        vm.InstrNot,
	GREATER:    vm.InstrMore,
	EQUALS:     vm.InstrEquals,
	READ_CHAR:  vm.InstrReadChar,
	WRITE_CHAR: vm.InstrWriteChar,
	WRITE_INT:  vm.InstrWriteInt,
	FLUSH:      vm.InstrFlush,
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(str string, out io.Writer) error {
	ti := TokenInput{Input: &input.StringInput{Str: str}}

	w := vm.NewBytecodeWriter()

	vars := make(map[string]int)
	for !ti.Eof() {
		if ti.IsInt() {
			if v, err := ti.ReadInt(); err == nil {
				w.WritePush(v)
			} else {
				return err
			}
		} else if ti.IsCharCode() {
			if v, err := ti.ReadCharCode(); err == nil {
				w.WritePush(int(v))
			} else {
				return err
			}
		} else if ti.IsVar() {
			if v, m, err := ti.ReadVar(); err == nil {
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
				return err
			}
		} else if ti.IsSubStart() {
			ti.SkipSubStart()
			w.SubCreate()
		} else if ti.IsSubEnd() {
			ti.SkipSubEnd()
			if err := w.SubReturn(); err != nil {
				ti.Input.Croak(err.Error())
				return err
			}
		} else if ti.IsSubCall() {
			ti.SkipSubCall()
			w.WriteCall()
		} else if ti.IsIf() {
			ti.SkipIf()
			w.WriteCallIf()
		} else if ti.IsWhile() {
			ti.SkipWhile()
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
				return err
			}
			// Call condition
			w.WriteFetch(ca)
			w.WriteCall()
			// Push to stack call body address
			w.WritePush(bca)
			// Call condition
			w.WriteGotoIf()
		} else if ti.IsCommand() {
			if ic, err := ti.ReadCommand(); err == nil {
				if cmd, ok := InstrMap[ic]; ok {
					w.WriteCommand(cmd)
				} else {
					err := errors.New("invalid command")
					ti.Input.Croak(err.Error())
					return err
				}
			} else {
				return err
			}
		} else if ti.IsString() {
			if s, err := ti.ReadString(); err == nil {
				w.WriteString(s)
			} else {
				ti.Input.Croak(err.Error())
				return err
			}
		} else if ti.IsWhitespace() {
			ti.SkipWhitespace()
		}
	}
	w.WriteEnd()
	_, err := w.WriteTo(out)
	return err
}
