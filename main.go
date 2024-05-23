package main

import (
	"fmt"
	"log"
	"os"
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

	p := NewParser()
	img, err := p.Parse(string(data))
	if err != nil {
		log.Fatalln("parsing failed")
	}
	h := ""
	for i := 0; i < len(img); i++ {
		h += fmt.Sprintf("%d ", img[i])
	}
	log.Printf("   image loaded: %s\n", h)

	vm := NewVM(1024, 60, 60)
	err = vm.Load(img)
	if err != nil {
		log.Fatalln("image loading failed")
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
