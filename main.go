package main

import (
	"bytes"
	"encoding/binary"
	"false-vm/bf"
	false2 "false-vm/false"
	"false-vm/input"
	vm2 "false-vm/vm"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var bcf string
	var src string
	var lang string
	var out string
	var run bool
	var verbose bool
	var memSize int
	var opStackSize int
	var callStackSize int
	flag.StringVar(&bcf, "b", "", "bytecode file (has more priority than source file parameter)")
	flag.StringVar(&src, "s", "", "source file (.bf and .false are supported)")
	flag.StringVar(&lang, "l", "auto", "force set language: auto (autodetect by file extension), false - FALSE, bf - Brainfuck")
	flag.StringVar(&out, "o", "", "output compiled bytecode to file")
	flag.BoolVar(&run, "r", true, "run compiled file")
	flag.BoolVar(&verbose, "v", false, "verbose log mode")
	flag.IntVar(&memSize, "m", 131072, "total memory size (32-bit integers)")
	flag.IntVar(&opStackSize, "os", 1280, "operation stack size (part of total memory; 32-bit integers)")
	flag.IntVar(&callStackSize, "cs", 640, "call stack size (part of total memory; 32-bit integers)")
	flag.Parse()

	var err error
	var bc []byte
	if bcf != "" {
		if bc, err = os.ReadFile(bcf); err != nil {
			log.Fatalln("unable to read bytecode file:", err.Error())
		}
	} else {
		if src == "" {
			log.Fatalln("source file is required")
		}

		if lang == "auto" {
			ext := strings.ToLower(filepath.Ext(src))
			switch ext {
			case ".bf":
				lang = "bf"
				break
			case ".f":
			case ".false":
				lang = "false"
				break
			default:
				log.Fatalln("unsupported file extension:", ext)
			}
		}

		var p input.Parser
		switch lang {
		case "bf":
			p = bf.NewParser()
			break
		case "false":
			p = false2.NewParser()
			break
		default:
			log.Fatalln("unsupported language:", lang)
		}

		r, err := os.Open(src)
		if err != nil {
			log.Fatalln("unable to open file:", err.Error())
		}
		w := new(bytes.Buffer)
		err = p.Parse(r, w)
		if err != nil {
			log.Fatalln("parsing failed")
		}
		err = r.Close()
		if err != nil {
			log.Fatalln("closing failed")
		}

		bc = w.Bytes()
	}

	if out != "" {
		if err := os.WriteFile(out, bc, 0644); err != nil {
			log.Fatalln("bytecode writing failed with error,", err.Error())
		}
		fmt.Printf("%d bytes written to file %s\n", len(bc), filepath.Base(out))
	}

	if run {
		unitSize := 4
		if len(bc)%unitSize != 0 {
			log.Fatalln("invalid byte alignment")
		}
		v := ""
		img := make([]int, len(bc)/unitSize)
		c := 0
		for i := 0; i < len(bc); i += unitSize {
			u := binary.LittleEndian.Uint32(bc[i : i+unitSize])
			img[c] = int(u)
			c++
			v += fmt.Sprintf("%d ", u)
		}
		logV(verbose, "image loaded: %s\n", v)

		vm := vm2.NewVM(memSize, opStackSize, callStackSize)
		err = vm.Load(img)
		if err != nil {
			log.Fatalln("image loading failed:", err)
		}

		before := time.Now().UnixMilli()
		ic, err := vm.Run()
		after := time.Now().UnixMilli()

		fmt.Printf("cpu instructions %d, time %d milliseconds\n", ic, after-before)

		if err != nil {
			log.Fatalln("vm fault:", err.Error())
		}
		if verbose {
			vm.Dump()
		}
	}
}

func logV(verbose bool, format string, a ...any) {
	if verbose {
		fmt.Printf(format, a...)
	}
}
