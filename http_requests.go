package go_wiktionary_parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

const (
	defaultLimit = 10
	apiURLFormat = "https://%s.wiktionary.org/w/api.php"
)

func SearchWordsForLanguage(query string, lang string) ([]string, error) {
	if query == "" {
		return nil, errors.New("empty query")
	}

	language := Language(lang)

	if !language.IsValid() {
		return []string{}, errors.New("invalid language")
	}

	baseURL, _ := url.Parse(fmt.Sprintf(apiURLFormat, lang))

	params := url.Values{}
	params.Add("action", "opensearch")
	params.Add("format", "json")
	params.Add("search", query)
	params.Add("limit", strconv.Itoa(defaultLimit))

	baseURL.RawQuery = params.Encode()

	resp, err := client.Get(baseURL.String())
	if err != nil {
		return []string{}, fmt.Errorf("error while querying data: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("failed to close response body: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, fmt.Errorf("error while reading the response: %v", err)
	}

	var result []interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return []string{}, fmt.Errorf("error while parsing JSON: %v", err)
	}

	suggestions, ok := result[1].([]interface{})
	if !ok {
		return []string{}, fmt.Errorf("error while extracting the suggestions")
	}

	words := make([]string, 0, len(suggestions))
	for _, suggestion := range suggestions {
		word, ok := suggestion.(string)
		if !ok {
			continue
		}
		words = append(words, word)
	}

	if len(words) == 0 {
		return nil, errors.New("no valid suggestions found")
	}

	return words, nil
}

func SearchWords(query string) ([]string, error) {
	res, err := SearchWordsForLanguage(query, "en")

	return res, err
}

func RetrieveArticleForLanguage(word string, lang string) (*ArticleContent, error) {
	language := Language(lang)

	if !language.IsValid() {
		return nil, errors.New("invalid language")
	}

	baseURL, err := url.Parse(fmt.Sprintf("https://%s.wiktionary.org/w/api.php", lang))
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	params := url.Values{}
	params.Add("action", "parse")
	params.Add("format", "json")
	params.Add("page", word)
	params.Add("prop", "text")
	params.Add("formatversion", "2")

	baseURL.RawQuery = params.Encode()

	resp, err := client.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("failed to close response body: %v", err)
		}
	}(resp.Body)

	var parseResp ParseResponse
	if err := json.NewDecoder(resp.Body).Decode(&parseResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if parseResp.Error.Info != "" {
		return nil, fmt.Errorf("API error: %s", parseResp.Error.Info)
	}

	article := &ArticleContent{
		Title:    parseResp.Parse.Title,
		HTML:     parseResp.Parse.Text,
		Language: lang,
	}

	return article, nil
}

func RetrieveArticle(word string) (*ArticleContent, error) {
	res, err := RetrieveArticleForLanguage(word, "en")

	return res, err
}
