package go_wiktionary_parser

type Language string

const (
	English Language = "en"
	German  Language = "de"
)

func (l Language) IsValid() bool {
	switch l {
	case English, German:
		return true
	}
	return false
}
