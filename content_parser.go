package go_wiktionary_parser

func FetchAndParseArticleForWord(word string, lang string) {
	article, err := RetrieveArticleForLanguage(word, lang)

	if err != nil {

	}

	ParseArticleForWord(article, lang)
}

func ParseArticleForWord(article *ArticleContent, lang string) {
	l := Language(lang)
	if !l.IsValid() {

	}

	switch l {
	case German:
		parseGerman(article)
	case English:
		parseEnglish(article)
	}
}

func parseGerman(content *ArticleContent) {

}

func parseEnglish(content *ArticleContent) {

}
