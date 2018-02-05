package main

import (
	"fmt"
	"os"
	"errors"
	"bufio"
)

var tape []int
var loopLocs []int
var currentByte int
var pointer int
var code []byte
var mode string

func printErr(err error) {
	fmt.Println()
	fmt.Println("ERROR @byte #", currentByte+1)
	fmt.Println(err)
	os.Exit(69)
}
func handleErr(err error) {
	if err != nil {
		printErr(err)
	}
}

func getCode() ([]byte, string) {
	if len(os.Args) >= 2 { //file input
		fileLoc := os.Args[1]
		f, err := os.Open(fileLoc)
		handleErr(err)
		stats, err := f.Stat()
		handleErr(err)
		buffer := make([]byte, stats.Size()) //gets size of file
		f.Read(buffer)
		f.Close()
		return buffer, "file"
	} else { //command line input
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("\n>>> ")
		scanner.Scan(); //get input
		return []byte(scanner.Text()), "loop"
	}
}

func movePointer(i int) {
	pointer += i
	if pointer < 0 {
		printErr(errors.New("Cannot access negitive memory location"))
	}
	if pointer > len(tape) - 1 {
		tape = append(tape, 0)
	}
}

func changeMem(i int) {
	tape[pointer] += i 
}

func output(n int) { //I know its inefficienct and messy to check for an extra arg every time it prints instead of having it check for the flag in getCode() but idgaf
	if len(os.Args) > 2 && os.Args[2] == "-s" && n >= 0 && n < 256 {
		fmt.Print(string(byte(n)))
	} else {
		fmt.Print(n)
	}
	fmt.Print(" ")
}

func getInput() int {
	var i int
	fmt.Print("\n>")
	_, err := fmt.Scan(&i)
	handleErr(err)
	return i
}

func startLoop() {
	if tape[pointer] == 0 { //skip the loop
		internalLoopStarts := 0 //TO HANDLE SHIT LIKE [[]] BECAUSE IF THE CURRENT CELL IS 0 WHEN THE PGRM REACHES THE FIRST "["" THEN IT WILL LOOK FOR THE NEXT "]" BUT IN THIS CASE THE NEXT "]" IS ACTUALLY BEFORE THE CORRECT "]" SO IT NEEDS TO COUNT HOW MANY INTERNAL LOOPS TO SKIP 
		currentByte ++ //goes to the next byte right away because the current byte is a "[" which will add to the internal loops if we don't skip it
		for currentByte < len(code) { //loops through bytes until it runs out of codes
			s := string(code[currentByte])
			if s == "]" { //if a loop end is found
				if internalLoopStarts == 0 { //breaks if there are no internal loops left to skip
					return
				} else { //if the loop end is part of an internal loop that must be skipped
					internalLoopStarts --
				}
			} else if s == "[" { //if the start of internal loop is found
				internalLoopStarts ++
			}
			currentByte ++
		}
		if currentByte == len(code)-1 && string(code[currentByte]) != "]" {
			printErr(errors.New("Missing ]"))
		}
	} else { //do the loop
		loopLocs = append([]int{currentByte}, loopLocs...)
	}
}

func endLoop() {
	if len(loopLocs) == 0 {
		printErr(errors.New("End of loop found with no start"))
	}
	if tape[pointer] == 0 {
		loopLocs = loopLocs[1:]
	} else {
		currentByte = loopLocs[0]
	}
}

func interpretSymbol(s string) {
	switch s {
	case ">":
		movePointer(1)
	case "<":
		movePointer(-1)
	case "+":
		changeMem(1)
	case "-":
		changeMem(-1)
	case ".":
		output(tape[pointer])
	case ",":
		tape[pointer] = getInput()
	case "[":
		startLoop()
	case "]":
		endLoop()
	}
}


func main() {
	tape = []int{0}
	pointer = 0
	mode = "loop"

	for mode == "loop" {
		//loopLocs = []int{} //reset loopLocs
		currentByte = 0 //reset currentByte
		code, mode = getCode()
		for currentByte < len(code) {
			chr := string(code[currentByte])
			interpretSymbol(chr)
			currentByte ++
		}
	}
}
