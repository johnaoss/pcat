package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/chroma"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func getfile(str string) []byte {
	body, err := ioutil.ReadFile(str)
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
		os.Exit(2)
	}
	return body
}

func main() {
	flag.Parse()

	var body []byte
	var err error
	if len(os.Args) == 2 && os.Args[1] != "" {
		body = getfile(os.Args[1])
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// piped from stdin
			body, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Printf("failed to read from stdin: %v", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("no input on stdin or filename provided")
			os.Exit(1)
		}

	}

	content := string(body)

	lex := lexers.Analyse(content)
	if lex == nil {
		lex = lexers.Fallback
	}
	lex = chroma.Coalesce(lex)

	style := styles.Get("monokai")
	formatter := formatters.Get("terminal256")

	iterator, err := lex.Tokenise(nil, content)
	if err != nil {
		fmt.Printf("failed to tokenize: %v", err)
		os.Exit(1)
	}

	if err := formatter.Format(os.Stdout, style, iterator); err != nil {
		fmt.Printf("failed to format: %v", err)
		os.Exit(1)
	}
}
