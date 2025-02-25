package go_wiktionary_parser

type ParseResponse struct {
	Parse struct {
		Title  string `json:"title"`
		PageID int    `json:"pageid"`
		Text   string `json:"text"`
	} `json:"parse"`
	Error struct {
		Code string `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
}

type ArticleContent struct {
	Title    string   `json:"title"`
	HTML     string   `json:"html"`
	Language string   `json:"language"`
	WordInfo WordInfo `json:"word_info,omitempty"`
}

type WordInfo interface{}

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

type EnglishWordInfo struct {
	Word        string   `json:"word"`        // The main word being described
	Plural      string   `json:"plural"`      // The plural form of the word
	Etymology   string   `json:"etymology"`   // Origin of the word and historical development
	Definitions []string `json:"definitions"` // List of definitions or meanings of the word.
}
