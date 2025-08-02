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

	for {
		t, parseErr := l.ConsumeToken()
		if parseErr != nil {
			fmt.Fprintln(os.Stderr, parseErr)
			os.Exit(1)
			return
		}

		fmt.Printf(
			"%s\n",
			t.Start.FileInfo.Contents[t.Start.Cursor:t.End.Cursor],
		)

		if t.Kind == token.TOKEN_END_OF_FILE {
			break
		}
	}
}
