package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

const (
	InstrPush int = 1

	InstrDup  int = 2
	InstrDrop int = 3
	InstrSwap int = 4
	InstrRot  int = 5
	InstrPick int = 6

	InstrPlus     int = 7
	InstrMinus    int = 8
	InstrMultiply int = 9
	InstrDivide   int = 10
	InstrNegative int = 11
	InstrAnd      int = 12
	InstrOr       int = 13
	InstrNot      int = 14

	InstrMore   int = 15
	InstrEquals int = 16

	InstrReadChar  int = 17
	InstrWriteChar int = 18
	InstrWriteInt  int = 19
	InstrWriteStr  int = 20
	InstrFlush     int = 21

	InstrStore int = 22
	InstrFetch int = 23

	InstrCall   int = 24
	InstrCallIf int = 25
	InstrReturn int = 26

	InstrGoto   int = 27
	InstrGotoIf int = 28

	InstrEnd int = 29
)

type BufferWrapper struct {
	buf          *bytes.Buffer
	prevLenBytes int
}

type BytecodeWriter struct {
	order binary.ByteOrder
	buffs *Stack
}

func NewBytecodeWriter() *BytecodeWriter {
	buffs := NewStack()
	buffs.Push(&BufferWrapper{buf: new(bytes.Buffer), prevLenBytes: 0})
	return &BytecodeWriter{
		order: binary.LittleEndian,
		buffs: buffs,
	}
}

func (w *BytecodeWriter) bufWrap() *BufferWrapper {
	if w.buffs.Len() > 0 {
		return w.buffs.Peek().(*BufferWrapper)
	}
	return nil
}

func (w *BytecodeWriter) buf() *bytes.Buffer {
	bw := w.bufWrap()
	if bw != nil {
		return bw.buf
	}
	return nil
}

func (w *BytecodeWriter) Len() int {
	return w.LenBytes() / 4
}

func (w *BytecodeWriter) LenBytes() int {
	bw := w.bufWrap()
	return bw.prevLenBytes + bw.buf.Len()
}

func (w *BytecodeWriter) WritePush(v int) {
	w.WriteCommand(InstrPush)
	w.WriteInt(v)
}

func (w *BytecodeWriter) WriteGoto(addr int) {
	w.WriteCommand(InstrGoto)
	w.WriteInt(addr)
}

func (w *BytecodeWriter) WriteGotoIf() {
	w.WriteCommand(InstrGotoIf)
}

func (w *BytecodeWriter) WriteGotoRel(diff int) int {
	w.WriteCommand(InstrGoto)
	addr := w.Len() + 1 + diff
	w.WriteInt(addr)
	return addr
}

func (w *BytecodeWriter) WriteStore(addr int) {
	w.WriteCommand(InstrStore)
	w.WriteInt(addr)
}

func (w *BytecodeWriter) WriteFetch(addr int) {
	w.WriteCommand(InstrFetch)
	w.WriteInt(addr)
}

func (w *BytecodeWriter) WriteCall() {
	w.WriteCommand(InstrCall)
}

func (w *BytecodeWriter) WriteCallIf() {
	w.WriteCommand(InstrCallIf)
}

func (w *BytecodeWriter) WriteString(s string) {
	w.WriteCommand(InstrWriteStr)
	b := []byte(s)
	w.WriteInt(len(b))
	for _, c := range b {
		w.WriteInt(int(c))
	}
}

func (w *BytecodeWriter) WriteEnd() {
	w.WriteCommand(InstrEnd)
}

func (w *BytecodeWriter) WriteVar(v int) int {
	w.WriteGotoRel(1)
	addr := w.Len()
	w.WriteInt(v)
	return addr
}

func (w *BytecodeWriter) BlockCreate() {
	w.buffs.Push(&BufferWrapper{buf: new(bytes.Buffer), prevLenBytes: w.LenBytes() + 8}) // TODO: check +8
}

func (w *BytecodeWriter) BlockSkip() (int, error) {
	bw := w.buffs.Pop().(*BufferWrapper)
	buf := bw.buf
	if buf == nil {
		return 0, fmt.Errorf("block end without creation")
	}
	w.WriteGotoRel(buf.Len() / 4) // Goto address to skip sub
	addr := w.Len()
	w.WriteBytes(buf.Bytes()) // Write sub content
	buf.Reset()
	return addr, nil
}

func (w *BytecodeWriter) SubCreate() {
	w.buffs.Push(&BufferWrapper{buf: new(bytes.Buffer), prevLenBytes: w.LenBytes() + 8}) // TODO: check +8
}

func (w *BytecodeWriter) SubReturn() error {
	w.WriteCommand(InstrReturn)
	bw := w.buffs.Pop().(*BufferWrapper)
	buf := bw.buf
	if buf == nil {
		return fmt.Errorf("sub return without creation")
	}
	w.WriteGotoRel(buf.Len() / 4) // Goto address to skip sub
	addr := w.Len()               // Sub start address
	w.WriteBytes(buf.Bytes())     // Write sub content
	w.WritePush(addr)             // Push sub start point to stack
	buf.Reset()
	return nil
}

func (w *BytecodeWriter) WriteCommand(c int) {
	w.WriteInt(c)
}

func (w *BytecodeWriter) WriteInt(v int) {
	err := binary.Write(w.buf(), w.order, int32(v))
	w.assertError(err)
}

func (w *BytecodeWriter) WriteBytes(p []byte) {
	_, err := w.buf().Write(p)
	w.assertError(err)
}

func (w *BytecodeWriter) Bytes() []byte {
	return w.buf().Bytes()
}

func (w *BytecodeWriter) assertError(err error) {
	if err != nil {
		log.Fatalln("bytecode write error:", err)
	}
}
