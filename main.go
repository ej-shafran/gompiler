package main

import (
	"fmt"
	"os"

	"github.com/ej-shafran/gompiler/pkg/lexer"
	"github.com/ej-shafran/gompiler/pkg/token"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Missing required argument: filename")
		os.Exit(1)
		return
	}

	filename := args[0]

	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	l := lexer.NewLexer(filename, string(file))

	_, perr := l.ExpectIdentifier("int")
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	_, perr = l.ExpectIdentifier("main")
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	_, perr = l.ExpectSymbol("(")
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	_, perr = l.ExpectIdentifier("void")
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	_, perr = l.ExpectSymbol(")")
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	_, perr = l.ExpectSymbol("{")
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	_, perr = l.ExpectIdentifier("return")
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	t, perr := l.ExpectTokenKind(token.TOKEN_NUMBER_LITERAL)
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	_, perr = l.ExpectSymbol(";")
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	_, perr = l.ExpectSymbol("}")
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}
	_, perr = l.ExpectTokenKind(token.TOKEN_END_OF_FILE)
	if perr != nil {
		fmt.Fprintln(os.Stderr, perr)
		os.Exit(1)
		return
	}

	value := l.TokenValue(t)
	fmt.Println("\t.globl\t_main")
	fmt.Println("_main:")
	fmt.Println("\t.cfi_startproc")
	fmt.Printf("\tmov\tw0, #%s\n", value)
	fmt.Println("\tret")
	fmt.Println("\t.cfi_endproc")
}
