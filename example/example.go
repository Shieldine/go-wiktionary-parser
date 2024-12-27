package main

import (
	"fmt"
	parser "github.com/Shieldine/go-wiktionary-parser"
)

func main() {
	// as a quickstart: you'll probably just need this
	parsed, err := parser.FetchAndParseArticleForWord("Baum", "de")
	if err != nil {
		fmt.Println(err)
	} else {
		// this is the whole object of type ArticleContent
		fmt.Println(parsed)

		// and this is the parsed info
		fmt.Println(parsed.WordInfo)
	}

	// and now some fine-grained examples

	// search for words (define language)
	res, err := parser.SearchWordsForLanguage("Bau", "de")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	// if you want english anyway, you can use a function that defaults to english
	res, err = parser.SearchWords("tes")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	// retrieve raw articles for a given language - WordInfo will be empty
	articleRes, err := parser.RetrieveArticleForLanguage("Baum", "de")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(articleRes)
	}

	// you can default to English here as well
	englishArticle, err := parser.RetrieveArticle("Tree")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(englishArticle)
	}

	// and here's the parsed info
	parsedInfo, _ := parser.ParseArticle(englishArticle, "en")

	fmt.Println(parsedInfo)

	// you can add it to your ArticleContent
	englishArticle.WordInfo = parsedInfo
}
