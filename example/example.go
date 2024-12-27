package main

import (
	"fmt"
	"github.com/Shieldine/go-wiktionary-parser"
)

func main() {
	fmt.Println("test")
	res, _ := go_wiktionary_parser.RetrieveArticle("test")

	fmt.Println(res)
}
