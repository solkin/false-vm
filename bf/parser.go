package bf

import (
	"false-vm/input"
	"false-vm/vm"
	"io"
)

type Parser struct {
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

	bc.BlockCreate()
	for i := 0; i < 30720; i++ {
		bc.WriteInt(0)
	}
	mem, err := bc.BlockSkip()
	if err != nil {
		return err
	}
	mp := bc.WriteVar(mem)

	for !ti.Eof() {
		if !ti.IsCommand() {
			ti.Skip()
		}
		cmd := ti.Next()
		switch cmd {
		case NEXT:
			bc.WriteFetch(mp)
			bc.WritePush(1)
			bc.WriteCommand(vm.InstrPlus)
			bc.WriteStore(mp)
			break
		case PREV:
			bc.WriteFetch(mp)
			bc.WritePush(1)
			bc.WriteCommand(vm.InstrMinus)
			bc.WriteStore(mp)
			break
		case PLUS:
			bc.WriteFetch(mp)
			bc.WriteCommand(vm.InstrDup)
			bc.WriteStore(bc.Len() + 10)
			bc.WriteStore(bc.Len() + 3)
			bc.WriteFetch(0) // Stub address, will be written by command before
			bc.WritePush(1)
			bc.WriteCommand(vm.InstrPlus)
			bc.WriteStore(0)
			break
		case MINUS:
			bc.WriteFetch(mp)
			bc.WriteCommand(vm.InstrDup)
			bc.WriteStore(bc.Len() + 10)
			bc.WriteStore(bc.Len() + 3)
			bc.WriteFetch(0) // Stub address, will be written by command before
			bc.WritePush(1)
			bc.WriteCommand(vm.InstrMinus)
			bc.WriteStore(0)
			break
		case IN:
			bc.WriteCommand(vm.InstrReadChar)
			bc.WriteCommand(vm.InstrDup)
			bc.WriteCommand(vm.InstrWriteChar)
			bc.WriteFetch(mp)
			bc.WriteStore(bc.Len() + 3)
			bc.WriteStore(0) // Stub address, will be written by command before
			break
		case OUT:
			bc.WriteFetch(mp)
			bc.WriteStore(bc.Len() + 3)
			bc.WriteFetch(0) // Stub address, will be written by command before
			bc.WriteCommand(vm.InstrWriteChar)
			break
		case SUB:
			// Create conditional sub
			bc.SubCreate()
			bc.WriteFetch(mp)
			bc.WriteStore(bc.Len() + 3)
			bc.WriteFetch(0) // Stub address, will be written by command before
			if err = bc.SubReturn(); err != nil {
				ti.Input.Croak(err.Error())
				return err
			}
			bc.SubCreate()
			break
		case RETURN:
			if err = bc.SubReturn(); err != nil {
				ti.Input.Croak(err.Error())
				return err
			}
			// Reserved condition and body addresses
			ba := bc.WriteVar(0)
			ca := bc.WriteVar(0)
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
			break
		}
	}
	bc.WriteEnd()
	_, err = bc.WriteTo(w)
	return err
}
