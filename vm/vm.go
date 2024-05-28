package vm

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"strconv"
)

type VM struct {
	Memory    []int
	pmOffset  int
	pmSize    int
	ip        int
	OpStack   *IntStack
	CallStack *IntStack
}

func NewVM(size int, opStackSize int, callStackSize int) *VM {
	memory := make([]int, size)
	return &VM{
		Memory:    memory,
		pmOffset:  0,
		pmSize:    size - callStackSize - opStackSize,
		OpStack:   NewIntStack(memory, len(memory)-callStackSize-opStackSize, opStackSize),
		CallStack: NewIntStack(memory, len(memory)-callStackSize, callStackSize),
	}
}

func (vm *VM) Load(img []int) error {
	if len(img) > len(vm.Memory) {
		return errors.New("image size is larger than allocated memory")
	}
	copy(vm.Memory, img)
	vm.ip = 0
	vm.OpStack.Reset()
	vm.CallStack.Reset()
	return nil
}

func (vm *VM) Run() error {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)

	defer func(w *bufio.Writer) {
		_ = w.Flush()
	}(w)
	_, _ = w.WriteString("vm started\n\n")
	w.Flush()

	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		i, err := vm.next()
		if err != nil {
			return err
		}
		switch i {
		case InstrPush:
			v := 0
			if v, err = vm.next(); err == nil {
				if err = vm.OpStack.Push(v); err == nil {
					break
				}
			}
			return err
		case InstrDup:
			v := vm.OpStack.Peek()
			err = vm.OpStack.Push(v)
			if err != nil {
				return err
			}
			break
		case InstrDrop:
			_, err = vm.OpStack.Pop()
			if err != nil {
				return err
			}
			break
		case InstrSwap:
			if v1, err := vm.OpStack.Pop(); err == nil {
				if v2, err := vm.OpStack.Pop(); err == nil {
					err = vm.OpStack.Push(v1)
					if err != nil {
						return err
					}
					err = vm.OpStack.Push(v2)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				return err
			}
			break
		case InstrRot:
			if x2, err := vm.OpStack.Pop(); err == nil {
				if x1, err := vm.OpStack.Pop(); err == nil {
					if x, err := vm.OpStack.Pop(); err == nil {
						err = vm.OpStack.Push(x1)
						if err != nil {
							return err
						}
						err = vm.OpStack.Push(x2)
						if err != nil {
							return err
						}
						err = vm.OpStack.Push(x)
						if err != nil {
							return err
						}
					} else {
						return err
					}
				} else {
					return err
				}
			} else {
				return err
			}
			break
		case InstrPick:
			v, err := vm.OpStack.Pop()
			if err != nil {
				return err
			}
			v, err = vm.OpStack.Pick(v)
			if err != nil {
				return err
			}
			err = vm.OpStack.Push(v)
			if err != nil {
				return err
			}
			break
		case InstrPlus:
			v1, v2, err := vm.OpStack.PopPop()
			if err == nil {
				if err = vm.OpStack.Push(v1 + v2); err == nil {
					break
				}
			}
			return err
		case InstrMinus:
			v1, v2, err := vm.OpStack.PopPop()
			if err == nil {
				if err = vm.OpStack.Push(v2 - v1); err == nil {
					break
				}
			}
			return err
		case InstrMultiply:
			v1, v2, err := vm.OpStack.PopPop()
			if err == nil {
				if err = vm.OpStack.Push(v1 * v2); err == nil {
					break
				}
			}
			return err
		case InstrDivide:
			v1, v2, err := vm.OpStack.PopPop()
			if err == nil {
				if err = vm.OpStack.Push(v2 / v1); err == nil {
					break
				}
			}
			return err
		case InstrNegative:
			v, err := vm.OpStack.Pop()
			if err == nil {
				if err = vm.OpStack.Push(-v); err == nil {
					break
				}
			}
			return err
		case InstrAnd:
			v1, v2, err := vm.OpStack.PopPop()
			if err == nil {
				res := 0
				if v1 != 0 && v2 != 0 {
					res = 1
				}
				if err = vm.OpStack.Push(res); err == nil {
					break
				}
			}
			return err
		case InstrOr:
			v1, v2, err := vm.OpStack.PopPop()
			if err == nil {
				res := 0
				if v1 != 0 || v2 != 0 {
					res = 1
				}
				if err = vm.OpStack.Push(res); err == nil {
					break
				}
			}
			return err
		case InstrNot:
			v1, err := vm.OpStack.Pop()
			if err == nil {
				res := 0
				if v1 == 0 {
					res = 1
				}
				if err = vm.OpStack.Push(res); err == nil {
					break
				}
			}
			return err
		case InstrEquals:
			v1, v2, err := vm.OpStack.PopPop()
			if err == nil {
				eq := 0
				if v1 == v2 {
					eq = 1
				}
				if err = vm.OpStack.Push(eq); err == nil {
					break
				}
			}
			return err
		case InstrMore:
			v1, v2, err := vm.OpStack.PopPop()
			if err == nil {
				eq := 0
				if v2 > v1 {
					eq = 1
				}
				if err = vm.OpStack.Push(eq); err == nil {
					break
				}
			}
			return err
		case InstrWriteInt:
			v, err := vm.OpStack.Pop()
			if err != nil {
				return err
			}
			_, _ = w.WriteString(strconv.Itoa(v))
			break
		case InstrWriteChar:
			v, err := vm.OpStack.Pop()
			if err != nil {
				return err
			}
			_, _ = w.WriteString(string(rune(v)))
			break
		case InstrWriteStr:
			l, err := vm.next()
			if err != nil {
				return err
			}
			for i := 0; i < l; i++ {
				v, err := vm.next()
				if err != nil {
					return err
				}
				_, _ = w.WriteString(string(rune(v)))
			}
			break
		case InstrReadChar:
			r, _, _ := r.ReadRune()
			err = vm.OpStack.Push(int(r))
			if err != nil {
				return err
			}
			break
		case InstrFlush:
			_ = w.Flush()
			break
		case InstrStore:
			var addr, val int
			if addr, err = vm.next(); err == nil {
				if addr >= vm.pmOffset && addr < vm.pmSize {
					if val, err = vm.OpStack.Pop(); err == nil {
						vm.Memory[addr] = val
						break
					}
				} else {
					err = errors.New("out of memory")
				}
			}
			return err
		case InstrFetch:
			var addr int
			if addr, err = vm.next(); err == nil {
				if addr >= vm.pmOffset && addr < vm.pmSize {
					val := vm.Memory[addr]
					if err = vm.OpStack.Push(val); err == nil {
						break
					}
				} else {
					err = errors.New("out of memory")
				}
			}
			return err
		case InstrCopy:
			var addr1, addr2 int
			if addr1, err = vm.next(); err == nil {
				if addr2, err = vm.next(); err == nil {
					if addr1 >= vm.pmOffset && addr1 < vm.pmSize &&
						addr2 >= vm.pmOffset && addr2 < vm.pmSize {
						vm.Memory[addr2] = vm.Memory[addr1]
						break
					} else {
						err = errors.New("out of memory")
					}
				}
			}
			return err
		case InstrCall:
			addr, err := vm.OpStack.Pop()
			if err != nil {
				return err
			}
			err = vm.CallStack.Push(vm.ip)
			if err != nil {
				return err
			}
			vm.ip = addr
			break
		case InstrCallIf:
			addr, err := vm.OpStack.Pop()
			if err != nil {
				return err
			}
			cond, err := vm.OpStack.Pop()
			if err != nil {
				return err
			}
			if cond != 0 {
				err = vm.CallStack.Push(vm.ip)
				if err != nil {
					return err
				}
				vm.ip = addr
			}
			break
		case InstrReturn:
			addr, err := vm.CallStack.Pop()
			if err != nil {
				return err
			}
			vm.ip = addr
			break
		case InstrGoto:
			addr, err := vm.next()
			if err != nil {
				return err
			}
			vm.ip = addr
			break
		case InstrGotoIf:
			addr := 0
			if addr, err = vm.OpStack.Pop(); err != nil {
				return err
			}
			cond := 0
			if cond, err = vm.OpStack.Pop(); err != nil {
				return err
			}
			if cond != 0 {
				vm.ip = addr
			}
			break
		case InstrEnd:
			_, _ = w.WriteString("\n\nvm gracefully stopped\n")
			return nil
		default:
			return errors.New("invalid instruction " + strconv.Itoa(i))
		}
	}
}

func (s *IntStack) PopPop() (int, int, error) {
	if v1, err := s.Pop(); err == nil {
		if v2, err := s.Pop(); err == nil {
			return v1, v2, err
		} else {
			return 0, 0, err
		}
	} else {
		return 0, 0, err
	}
}

func (vm *VM) next() (int, error) {
	if vm.ip >= len(vm.Memory) {
		return 0, errors.New("out of memory")
	}
	i := vm.Memory[vm.ip]
	vm.ip++
	return i, nil
}

func (vm *VM) Fault() {
	log.Printf("fault on address %d\n", vm.ip)
}

func (vm *VM) Dump() {
	log.Println("instruction pointer: ", vm.ip)
	h := ""
	for i := 0; i < vm.ip; i++ {
		h += fmt.Sprintf("%d ", vm.Memory[i])
	}
	log.Println("image dump: ", h)
	h = ""
	for i := vm.OpStack.Offset + vm.OpStack.Size - 1; i >= vm.OpStack.p; i-- {
		h += fmt.Sprintf("%d ", vm.OpStack.Array[i])
	}
	log.Println("op stack dump: ", h)
	h = ""
	for i := vm.CallStack.Offset + vm.CallStack.Size - 1; i >= vm.CallStack.p; i-- {
		h += fmt.Sprintf("%d ", vm.CallStack.Array[i])
	}
	log.Println("call stack dump: ", h)
}
