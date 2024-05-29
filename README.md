# False-VM
Virtual Machine + False and Brainfuck compilers

![media/term.png](media/term.png)

Build instructions
------------------

Requirements:

* Go 1.11+

Assuming you have it, type:

```
go build
```

This will produce `false-vm` executable in the project folder.

Usage
------------------

Run false-vm with `-help` to get all flags info:

```
./false-vm -help

Usage of ./false-vm:
  -b string
    	bytecode file (has more priority than source file parameter)
  -cs int
    	call stack size (part of total memory; 32-bit integers) (default 640)
  -l string
    	force set language: auto (autodetect by file extension), false - FALSE, bf - Brainfuck (default "auto")
  -m int
    	total memory size (32-bit integers) (default 131072)
  -o string
    	output compiled bytecode to file
  -os int
    	operation stack size (part of total memory; 32-bit integers) (default 1280)
  -r	run compiled file (default true)
  -s string
    	source file (.bf and .false are supported)
  -v	verbose log mode
```

To compile and run Fibonacci sample:

```
./false-vm -s false/samples/fibonacci.false
```

This will compile Fibonacci sample to bytecode file without running:

```
./false-vm -s false/samples/fibonacci.false -o fib.fbc -r=0
```

and when you need to run, just type:

```
./false-vm -b fib.fbc
```

VM bytecode specification
------------------

Bytecode consists of 32-bit ints in little-endian encoding

| Instruction | Code | Args | Stack Change | Description                                                                                |
|-------------|------|------|--------------|--------------------------------------------------------------------------------------------|
| Push        | 1    | 1    | +1           | Push argument integer to the stack                                                         |
| Dup         | 2    | 0    | +1           | Duplicate topmost stack item                                                               |
| Drop        | 3    | 0    | -1           | Delete topmost stack item                                                                  |
| Swap        | 4    | 0    | 0            | Swap to topmost stack-items                                                                |
| Rot         | 5    | 0    | 0            | Rotate 3rd stack item to top                                                               |
| Pick        | 6    | 0    | +1           | Copy n-th item to top                                                                      |
| Plus        | 7    | 0    | -1           | Sum two topmost stack-items and push result to stack                                       |
| Minus       | 8    | 0    | -1           | Minus two topmost stack-items and push result to stack                                     |
| Multiply    | 9    | 0    | -1           | Multiply two topmost stack-items and push result to stack                                  |
| Divide      | 10   | 0    | -1           | Divide two topmost stack-items and push result to stack                                    |
| Negative    | 11   | 0    | 0            | Change topmost stack item sign value                                                       |
| And         | 12   | 0    | -1           | Apply logical 'and' for two topmost stack-items and push result to stack                   |
| Or          | 13   | 0    | -1           | Apply logical 'or' for two topmost stack-items and push result to stack                    |
| Not         | 14   | 0    | 0            | Apply logical 'not' topmost stack item and push result to stack                            |
| More        | 15   | 0    | -1           | Check for one topmost stack item is more than another                                      |
| Equals      | 16   | 0    | -1           | Check for tow topmost stack items has the are same value                                   |
| ReadChar    | 17   | 0    | +1           | Reads char from user and put item to top                                                   |
| WriteChar   | 18   | 0    | -1           | Output top stack item as char to the term                                                  |
| WriteInt    | 19   | 0    | -1           | Output top stack item as int to the term                                                   |
| WriteString | 20   | 1+n  | 0            | Read argument int as string len + n string chars                                           |
| Flush       | 21   | 0    | 0            | Push argument integer to the stack                                                         |
| Store       | 22   | 1    | -1           | Store stack item to args address                                                           |
| Fetch       | 23   | 1    | +1           | Put value by args address as a stack item                                                  |
| Copy        | 24   | 2    | 0            | Copy value from one args address to another address                                        |
| Call        | 25   | 0    | -1           | Save current pc to call stack; takes sub pointer from stack; go to sub                     |
| CallIf      | 26   | 0    | -2           | Take condition and body sub pointers from stack; run condition, call body on 'true' result |
| Return      | 27   | 0    | 0            | Take pc from call stack and go to it                                                       |
| Goto        | 28   | 1    | 0            | Change pc to the argument pointer                                                          |
| GotoIf      | 29   | 0    | -2           | Same as CallIf, but goto instead of call                                                   |
| End         | 30   | 0    | 0            | Exit program                                                                               |

