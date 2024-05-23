package main

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
	ip        int
	OpStack   *Stack
	CallStack *Stack
}

func NewVM(size int, opStackSize int, callStackSize int) *VM {
	memory := make([]int, size)
	return &VM{
		Memory:    memory,
		OpStack:   NewStack(memory, len(memory)-callStackSize-opStackSize, opStackSize),
		CallStack: NewStack(memory, len(memory)-callStackSize, callStackSize),
	}
}

func (vm *VM) Load(img []int) error {
	if len(img) > len(vm.Memory) {
		return errors.New("image size is larger than allocated memory")
	}
	for c := 0; c < len(vm.Memory); c++ {
		if c < len(img) {
			vm.Memory[c] = img[c]
		} else {
			vm.Memory[c] = 0
		}
	}
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
			v, err := vm.next()
			if err != nil {
				return err
			}
			if err := vm.OpStack.Push(v); err != nil {
				vm.Fault()
				vm.Dump()
			}
			break
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
			if v1, err := vm.OpStack.Pop(); err == nil {
				if v2, err := vm.OpStack.Pop(); err == nil {
					err = vm.OpStack.Push(v1 + v2)
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
		case InstrMinus:
			if v1, err := vm.OpStack.Pop(); err == nil {
				if v2, err := vm.OpStack.Pop(); err == nil {
					err = vm.OpStack.Push(v2 - v1)
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
		case InstrMultiply:
			if v1, err := vm.OpStack.Pop(); err == nil {
				if v2, err := vm.OpStack.Pop(); err == nil {
					err = vm.OpStack.Push(v1 * v2)
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
		case InstrDivide:
			if v1, err := vm.OpStack.Pop(); err == nil {
				if v2, err := vm.OpStack.Pop(); err == nil {
					err = vm.OpStack.Push(v2 / v1)
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
		case InstrNegative:
			if v, err := vm.OpStack.Pop(); err == nil {
				err = vm.OpStack.Push(-v)
				if err != nil {
					return err
				}
			} else {
				return err
			}
			break
		case InstrAnd:
			if v1, err := vm.OpStack.Pop(); err == nil {
				if v2, err := vm.OpStack.Pop(); err == nil {
					res := 0
					if v1 != 0 && v2 != 0 {
						res = 1
					}
					err = vm.OpStack.Push(res)
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
		case InstrOr:
			if v1, err := vm.OpStack.Pop(); err == nil {
				if v2, err := vm.OpStack.Pop(); err == nil {
					res := 0
					if v1 != 0 || v2 != 0 {
						res = 1
					}
					err = vm.OpStack.Push(res)
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
		case InstrNot:
			if v1, err := vm.OpStack.Pop(); err == nil {
				res := 0
				if v1 == 0 {
					res = 1
				}
				err = vm.OpStack.Push(res)
				if err != nil {
					return err
				}
			} else {
				return err
			}
			break
		case InstrEquals:
			if v1, err := vm.OpStack.Pop(); err == nil {
				if v2, err := vm.OpStack.Pop(); err == nil {
					eq := 0
					if v1 == v2 {
						eq = 1
					}
					err = vm.OpStack.Push(eq)
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
		case InstrMore:
			if v1, err := vm.OpStack.Pop(); err == nil {
				if v2, err := vm.OpStack.Pop(); err == nil {
					eq := 0
					if v2 > v1 {
						eq = 1
					}
					err = vm.OpStack.Push(eq)
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
			addr, err := vm.next()
			if err != nil {
				return err
			}
			val, err := vm.OpStack.Pop()
			if err != nil {
				return err
			}
			vm.Memory[addr] = val
			break
		case InstrFetch:
			addr, err := vm.next()
			if err != nil {
				return err
			}
			val := vm.Memory[addr]
			err = vm.OpStack.Push(val)
			if err != nil {
				return err
			}
			break
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
	h := ""
	for i := 0; i < vm.ip; i++ {
		h += fmt.Sprintf("%d ", vm.Memory[i])
	}
	log.Println("     image dump: ", h)
	h = ""
	for i := vm.OpStack.Offset + vm.OpStack.Size - 1; i >= vm.OpStack.p; i-- {
		h += fmt.Sprintf("%d ", vm.OpStack.Array[i])
	}
	log.Println("  op stack dump: ", h)
	h = ""
	for i := vm.CallStack.Offset + vm.CallStack.Size - 1; i >= vm.CallStack.p; i-- {
		h += fmt.Sprintf("%d ", vm.CallStack.Array[i])
	}
	log.Println("call stack dump: ", h)
}
