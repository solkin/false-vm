package main

type DynArray struct {
	Data []int
	Size int
}

func NewDynArray() *DynArray {
	return &DynArray{
		Data: make([]int, 1024),
		Size: 0,
	}
}

func (a *DynArray) Add(v int) {
	p := a.Size + 1
	if len(a.Data) <= p {
		a.Enlarge()
	}
	a.Data[p-1] = v
	a.Size = p
}

func (a *DynArray) ToArray() []int {
	exp := make([]int, a.Size)
	copy(exp, a.Data[:a.Size])
	return exp
}

func (a *DynArray) Enlarge() {
	newData := make([]int, len(a.Data)+1024)
	copy(newData, a.Data)
	a.Data = newData
}
