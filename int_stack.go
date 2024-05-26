package main

import "errors"

type IntStack struct {
	Array  []int
	Offset int
	Size   int
	p      int
}

func NewIntStack(arr []int, offset int, size int) *IntStack {
	return &IntStack{
		Array:  arr,
		Offset: offset,
		Size:   size,
		p:      offset + size,
	}
}

func (s *IntStack) Reset() {
	s.p = s.Offset + s.Size

}

func (s *IntStack) Pick(v int) (int, error) {
	if s.p+v >= s.Offset+s.Size {
		return 0, errors.New("stack out of range")

	}
	return s.Array[s.p+v], nil
}

func (s *IntStack) Push(v int) error {
	if s.p-1 < s.Offset {
		return errors.New("stack overflow")

	}
	s.p--
	s.Array[s.p] = v
	return nil
}

func (s *IntStack) Peek() int {
	return s.Array[s.p]
}

func (s *IntStack) Pop() (int, error) {
	if s.p >= s.Offset+s.Size {
		return 0, errors.New("stack underflow")
	}
	i := s.Array[s.p]
	s.p++
	return i, nil
}
