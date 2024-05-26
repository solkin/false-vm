package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sandbox-vm/bf"
	"sandbox-vm/false"
	vm2 "sandbox-vm/vm"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("file name argument required")
	}
	name := os.Args[1]
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatalln("unable to read file:", err.Error())
	}

	ext := filepath.Ext(name)
	var img []byte
	if ext == ".false" {
		p := false.NewParser()
		img, err = p.Parse(string(data))
		if err != nil {
			log.Fatalln("parsing failed")
		}
	} else if ext == ".bf" {
		p := bf.NewParser()
		img, err = p.Parse(string(data))
		if err != nil {
			log.Fatalln("parsing failed")
		}
	} else {
		log.Fatalln("unsupported file type")
	}

	h := ""
	img32 := make([]int, len(img)/4)
	c := 0
	for i := 0; i < len(img); i += 4 {
		i32 := binary.LittleEndian.Uint32(img[i : i+4])
		img32[c] = int(i32)
		c++
		h += fmt.Sprintf("%d ", i32)
	}
	log.Printf("   image loaded: %s\n", h)

	vm := vm2.NewVM(10240, 60, 60)
	err = vm.Load(img32)
	if err != nil {
		log.Fatalln("image loading failed:", err)
	}
	before := time.Now().UnixMilli()
	err = vm.Run()
	after := time.Now().UnixMilli()
	log.Println("cpu time: ", after-before, "milliseconds")
	vm.Dump()
	if err != nil {
		log.Fatalln("vm fault:", err.Error())
	}
}
