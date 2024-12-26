package go_wiktionary_parser

type ParseResponse struct {
	Parse struct {
		Title    string `json:"title"`
		PageID   int    `json:"pageid"`
		Text     string `json:"text"`
		Sections []struct {
			Title  string `json:"line"`
			Index  string `json:"index"`
			Level  string `json:"level"`
			Number string `json:"number"`
			Anchor string `json:"anchor"`
		} `json:"sections"`
	} `json:"parse"`
	Error struct {
		Code string `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
}

type ArticleContent struct {
	Title           string          `json:"title"`
	HTML            string          `json:"html"`
	Language        string          `json:"language"`
	Sections        []Section       `json:"sections"`
	GrammaticalInfo GrammaticalInfo `json:"grammatical_info,omitempty"`
}

type Section struct {
	Title  string `json:"title"`
	Level  string `json:"level"`
	Anchor string `json:"anchor"`
}

type GrammaticalInfo interface {
	GetWordClass() string
}

type GermanGrammaticalInfo struct {
	Genus     string   `json:"genus,omitempty"`
	Plural    string   `json:"plural,omitempty"`
	Article   string   `json:"article,omitempty"`
	Cases     []string `json:"cases,omitempty"`
	WordClass string   `json:"word_class,omitempty"`
}

type EnglishGrammaticalInfo struct {
	WordClass   string   `json:"word_class,omitempty"`
	Plural      string   `json:"plural,omitempty"`
	Participles []string `json:"participles,omitempty"`
	Comparative string   `json:"comparative,omitempty"`
	Superlative string   `json:"superlative,omitempty"`
	Countable   *bool    `json:"countable,omitempty"`
}
