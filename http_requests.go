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

func searchWordsForLanguage(query string, lang string) ([]string, error) {
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

func searchWords(query string) ([]string, error) {
	res, err := searchWordsForLanguage(query, "en")

	return res, err
}

func retrieveArticleForLanguage(word string, lang string) (string, error) {
	language := Language(lang)

	if !language.IsValid() {
		return "", errors.New("invalid language")
	}

	baseURL, err := url.Parse(fmt.Sprintf("https://%s.wiktionary.org/w/api.php", lang))
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
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
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("failed to close response body: %v", err)
		}
	}(resp.Body)

	var result struct {
		Parse struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		} `json:"parse"`
		Error struct {
			Info string `json:"info"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if result.Error.Info != "" {
		return "", fmt.Errorf("API error: %s", result.Error.Info)
	}

	return result.Parse.Text, nil
}

func retrieveArticle(word string, lang string) (string, error) {
	res, err := retrieveArticleForLanguage(word, "en")

	return res, err
}
