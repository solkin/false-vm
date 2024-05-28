package bf

import (
	"io"
	"sandbox-vm/input"
	"sandbox-vm/vm"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(str string, out io.Writer) error {
	ti := TokenInput{Input: &input.StringInput{Str: str}}

	w := vm.NewBytecodeWriter()

	w.BlockCreate()
	for i := 0; i < 1024; i++ {
		w.WriteInt(0)
	}
	mem, err := w.BlockSkip()
	if err != nil {
		return err
	}
	mp := w.WriteVar(mem)

	for !ti.Eof() {
		cmd, err := ti.Next()
		if err != nil {
			ti.Input.Croak(err.Error())
			return err
		}
		switch cmd {
		case NEXT:
			w.WriteFetch(mp)
			w.WritePush(1)
			w.WriteCommand(vm.InstrPlus)
			w.WriteStore(mp)
			break
		case PREV:
			w.WriteFetch(mp)
			w.WritePush(1)
			w.WriteCommand(vm.InstrMinus)
			w.WriteStore(mp)
			break
		case PLUS:
			w.WriteFetch(mp)
			w.WriteCommand(vm.InstrDup)
			w.WriteStore(w.Len() + 10)
			w.WriteStore(w.Len() + 3)
			w.WriteFetch(0) // Stub address, will be written by command before
			w.WritePush(1)
			w.WriteCommand(vm.InstrPlus)
			w.WriteStore(0)
			break
		case MINUS:
			w.WriteFetch(mp)
			w.WriteCommand(vm.InstrDup)
			w.WriteStore(w.Len() + 10)
			w.WriteStore(w.Len() + 3)
			w.WriteFetch(0) // Stub address, will be written by command before
			w.WritePush(1)
			w.WriteCommand(vm.InstrMinus)
			w.WriteStore(0)
			break
		case IN:
			w.WriteCommand(vm.InstrReadChar)
			w.WriteFetch(mp)
			w.WriteStore(w.Len() + 3)
			w.WriteStore(0) // Stub address, will be written by command before
			break
		case OUT:
			w.WriteFetch(mp)
			w.WriteStore(w.Len() + 3)
			w.WriteFetch(0) // Stub address, will be written by command before
			w.WriteCommand(vm.InstrWriteChar)
			break
		case SUB:
			// Create conditional sub
			w.SubCreate()
			w.WriteFetch(mp)
			w.WriteStore(w.Len() + 3)
			w.WriteFetch(0) // Stub address, will be written by command before
			if err = w.SubReturn(); err != nil {
				ti.Input.Croak(err.Error())
				return err
			}
			w.SubCreate()
			break
		case RETURN:
			if err = w.SubReturn(); err != nil {
				ti.Input.Croak(err.Error())
				return err
			}
			// Reserved condition and body addresses
			ba := w.WriteVar(0)
			ca := w.WriteVar(0)
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
			break
		}
	}
	w.WriteEnd()
	_, err = w.WriteTo(out)
	return err
}
