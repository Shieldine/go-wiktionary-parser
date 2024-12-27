package go_wiktionary_parser

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func FetchAndParseArticleForWord(word string, lang string) (*ArticleContent, error) {
	article, err := RetrieveArticleForLanguage(word, lang)
	if err != nil {
		return nil, err
	}

	parsed, err := ParseArticle(article, lang)
	if err != nil {
		return nil, err
	}

	article.WordInfo = parsed

	return article, nil
}

func ParseArticle(article *ArticleContent, lang string) (WordInfo, error) {
	l := Language(lang)
	if !l.IsValid() {
		return nil, errors.New("invalid language")
	}

	switch l {
	case German:
		parsed, err := parseGerman(article)
		if err != nil {
			return nil, err
		}
		return &parsed, nil
	case English:
		parseEnglish(article)
		return nil, nil
	default:
		return nil, errors.New("invalid language")
	}
}

func parseGerman(content *ArticleContent) (*GermanWordInfo, error) {
	htmlContent := content.HTML

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var info GermanWordInfo

	info.Word = content.Title

	// singular + plural
	var foundNominative bool
	doc.Find("th").Each(func(i int, s *goquery.Selection) {
		if foundNominative {
			return
		}

		if strings.Contains(s.Text(), "Nominativ") {
			singularTD := s.Parent().Find("td").First()
			info.Singular = strings.TrimSpace(singularTD.Text())

			secondTD := singularTD.Next()
			info.Plural = strings.TrimSpace(secondTD.Text())

			foundNominative = true
		}
	})

	// gender + category
	var foundGenderCategory bool
	doc.Find("div.mw-heading.mw-heading3 h3").Each(func(i int, s *goquery.Selection) {
		if foundGenderCategory {
			return
		}

		grammarLink := s.Find("a[title='Hilfe:Wortart']")
		info.GrammaticalCategory = strings.TrimSpace(grammarLink.Text())

		genderEl := s.Find("em")
		info.Gender = strings.TrimSpace(genderEl.Text())

		foundGenderCategory = true
	})

	// definitions
	pSelection := doc.Find(`p:contains("Bedeutungen:")`).First()
	if pSelection.Length() != 0 {
		nextSiblings := pSelection.NextUntil("p")

		var definitions []string
		nextSiblings.Find("dd").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			definitions = append(definitions, text)
		})
		info.Definitions = definitions
	}

	// Etymology
	pSelection = doc.Find(`p:contains("Herkunft:")`).First()
	if pSelection.Length() != 0 {
		nextSiblings := pSelection.NextUntil("p")

		var etymology string
		nextSiblings.Find("dd").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			etymology += text
		})
		info.Etymology = etymology
	}

	// examples
	pSelection = doc.Find(`p:contains("Beispiele:")`).First()
	if pSelection.Length() != 0 {
		nextSiblings := pSelection.NextUntil("p")

		var examples []string
		nextSiblings.Find("dd").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			examples = append(examples, text)
		})
		info.Examples = examples
	}

	// phrases
	pSelection = doc.Find(`p:contains("Redewendungen:")`).First()
	if pSelection.Length() != 0 {
		nextSiblings := pSelection.NextUntil("p")

		var phrases []string
		nextSiblings.Find("dd").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			phrases = append(phrases, text)
		})
		info.Phrases = phrases
	}

	return &info, nil
}

func parseEnglish(content *ArticleContent) (*EnglishWordInfo, error) {
	htmlContent := content.HTML

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var info EnglishWordInfo
	info.Word = content.Title

	// plural form
	iPlural := doc.Find("i:contains('plural')").First()
	if iPlural.Length() > 0 {
		iPlural.NextAll().Filter("b").Each(func(_ int, bSel *goquery.Selection) {
			formText := strings.TrimSpace(bSel.Text())
			if formText != "" {
				info.Plural += formText + " "
			}
		})
	}

	// etymology
	etyHeading := doc.Find(`h3#Etymology`).First()
	if etyHeading.Length() > 0 {

		etyParagraph := etyHeading.
			Parent().
			NextAll().
			Filter("p").
			First()

		if etyParagraph.Length() > 0 {
			etymologyText := strings.TrimSpace(etyParagraph.Text())

			info.Etymology = etymologyText
		}

	}

	// definitions
	pWithHeadword := doc.Find(`p:has(span.headword-line)`).First()
	if pWithHeadword.Length() > 0 {

		defsOl := pWithHeadword.NextAll().Filter("ol").First()
		if defsOl.Length() > 0 {

			defsOl.Find("li").Each(func(i int, liSel *goquery.Selection) {
				if goquery.NodeName(liSel.Parent()) == "ol" {

					clonedLi := liSel.Clone()

					clonedLi.Find("ul").Remove()
					rawText := strings.TrimSpace(clonedLi.Text())
					rawText = strings.Join(strings.Fields(rawText), " ")

					info.Definitions = append(info.Definitions, rawText)
				} else {
					return
				}
			})
		}
	}

	return &info, nil
}
