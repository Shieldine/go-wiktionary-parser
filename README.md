# Go Wiktionary Parser
A package of useful functions to deal with Wiktionary articles

> [!IMPORTANT]  
> Supported Wiktionaries: English (en.wiktionary.org) and German (de.wiktionary.org).
> For now, the package only supports parsing the first meaning in the article.
> This will likely change in the future. Check [the development section](#development-and-contribution) for current
> plans.

## Quickstart
There isn't much, so I'll be brief.

### Installation
Install with:

```bash
go get github.com/Shieldine/go-wiktionary-parser
```

### Structs
You will mainly deal with the following struct:

```Go
type ArticleContent struct {
	Title    string   `json:"title"`
	HTML     string   `json:"html"`
	Language string   `json:"language"`
	WordInfo WordInfo `json:"word_info,omitempty"`
}
```

WordInfo being most important and dependent on your language choice. For English, you'll find the following:

```Go
type EnglishWordInfo struct {
	Word        string   `json:"word"`        // The main word being described
	Plural      string   `json:"plural"`      // The plural form of the word
	Etymology   string   `json:"etymology"`   // Origin of the word and historical development
	Definitions []string `json:"definitions"` // List of definitions or meanings of the word.
}
```

German looks as follows:

```Go
type GermanWordInfo struct {
	Word                string   `json:"word"`                 // The word being defined (e.g., "Baum").
	GrammaticalCategory string   `json:"grammatical_category"` // The grammatical category (e.g., noun).
	Gender              string   `json:"gender"`               // The grammatical gender of the word (e.g., masculine).
	Singular            string   `json:"singular"`             // The singular form of the word.
	Plural              string   `json:"plural"`               // The plural form of the word.
	Definitions         []string `json:"definitions"`          // List of definitions or meanings of the word.
	Etymology           string   `json:"etymology"`            // Historical origin and linguistic evolution of the word.
	Examples            []string `json:"examples"`             // Usage examples
	Phrases             []string `json:"phrases"`              // Phrases used in everyday speech
}
```

### Usage example
The example can be found in [example.go](./example/example.go)

```Go
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
		// fmt.Println(parsed)

		// and this is the parsed info - you need to cast it as it is defined as an interface
		fmt.Println(parsed.WordInfo.(*parser.GermanWordInfo))
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
		// fmt.Println(articleRes)
		fmt.Println(articleRes.Title)
	}

	// you can default to English here as well
	englishArticle, err := parser.RetrieveArticle("tree")
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(englishArticle)
		fmt.Println(englishArticle.Title)
	}

	// and here's the parsed info
	parsedInfo, err := parser.ParseArticle(englishArticle, "en")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(parsedInfo)
	}

	// you can add it to your ArticleContent
	englishArticle.WordInfo = parsedInfo
}
```

## Supported languages
The supported languages are:

- English (en.wiktionary.org)
- German (de.wiktionary.org)

## Development and contribution
This project is in occasional development (I sadly don't have as much time as I'd like to).

Current plans:
- Parse all polysemes of a word instead of the first one
- Parse further sections (currently, only a few get parsed)
- Add support for other Wiktionaries (pl, it, fr)

If you want to contribute, feel free to create issues and merge requests.

## License
This project is licensed under Apache 2.0.
You can find the license [here](./LICENSE).
