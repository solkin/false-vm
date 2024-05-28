package false

import (
	"errors"
	"false-vm/input"
	"false-vm/vm"
	"io"
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

func (p *Parser) Parse(r io.Reader, w io.Writer) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	ti := TokenInput{Input: &input.StringInput{Str: string(data)}}

	bc := vm.NewBytecodeWriter()

	vars := make(map[string]int)
	for !ti.Eof() {
		if ti.IsInt() {
			if v, err := ti.ReadInt(); err == nil {
				bc.WritePush(v)
			} else {
				return err
			}
		} else if ti.IsCharCode() {
			if v, err := ti.ReadCharCode(); err == nil {
				bc.WritePush(int(v))
			} else {
				return err
			}
		} else if ti.IsVar() {
			if v, m, err := ti.ReadVar(); err == nil {
				var addr int
				ok := false
				if addr, ok = vars[v]; !ok {
					// New variable
					addr = bc.WriteVar(0)
					vars[v] = addr
				}
				switch m {
				case STORE_VAR:
					bc.WriteStore(addr)
					break
				case FETCH_VAR:
					bc.WriteFetch(addr)
					break
				}
			} else {
				return err
			}
		} else if ti.IsSubStart() {
			ti.SkipSubStart()
			bc.SubCreate()
		} else if ti.IsSubEnd() {
			ti.SkipSubEnd()
			if err := bc.SubReturn(); err != nil {
				ti.Input.Croak(err.Error())
				return err
			}
		} else if ti.IsSubCall() {
			ti.SkipSubCall()
			bc.WriteCall()
		} else if ti.IsIf() {
			ti.SkipIf()
			bc.WriteCallIf()
		} else if ti.IsWhile() {
			ti.SkipWhile()
			// Reserved condition and body addresses
			ca := bc.WriteVar(0)
			ba := bc.WriteVar(0)
			// Take body and condition addresses down from stack
			bc.WriteStore(ba)
			bc.WriteStore(ca)
			// Call body (by goto)
			bc.BlockCreate()
			bc.WriteFetch(ba)
			bc.WriteCall()
			bca, err := bc.BlockSkip()
			if err != nil {
				return err
			}
			// Call condition
			bc.WriteFetch(ca)
			bc.WriteCall()
			// Push to stack call body address
			bc.WritePush(bca)
			// Call condition
			bc.WriteGotoIf()
		} else if ti.IsCommand() {
			if ic, err := ti.ReadCommand(); err == nil {
				if cmd, ok := InstrMap[ic]; ok {
					bc.WriteCommand(cmd)
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
				bc.WriteString(s)
			} else {
				ti.Input.Croak(err.Error())
				return err
			}
		} else if ti.IsCommentStart() {
			if _, err := ti.ReadComment(); err != nil {
				ti.Input.Croak(err.Error())
				return err
			}
		} else if ti.IsWhitespace() {
			ti.SkipWhitespace()
		}
	}
	bc.WriteEnd()
	_, err = bc.WriteTo(w)
	return err
}
