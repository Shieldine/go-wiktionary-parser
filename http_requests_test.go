package go_wiktionary_parser

import (
	"strings"
	"testing"
)

func TestSearchWordsForLanguage(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		lang        string
		wantErr     bool
		errContains string
		checkResult func([]string) bool
	}{
		{
			name:        "empty query",
			query:       "",
			lang:        "en",
			wantErr:     true,
			errContains: "empty query",
		},
		{
			name:        "invalid language",
			query:       "test",
			lang:        "invalid",
			wantErr:     true,
			errContains: "invalid language",
		},
		{
			name:    "valid english search",
			query:   "test",
			lang:    "en",
			wantErr: false,
			checkResult: func(results []string) bool {
				return len(results) > 0 && len(results) <= defaultLimit
			},
		},
		{
			name:    "partial word search",
			query:   "prog",
			lang:    "en",
			wantErr: false,
			checkResult: func(results []string) bool {
				return len(results) > 0 && strings.Contains(strings.Join(results, " "), "prog")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchWordsForLanguage(tt.query, tt.lang)

			if (err != nil) != tt.wantErr {
				t.Errorf("searchWordsForLanguage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("searchWordsForLanguage() error = %v, should contain %v", err, tt.errContains)
				return
			}

			if !tt.wantErr && tt.checkResult != nil {
				if !tt.checkResult(got) {
					t.Errorf("searchWordsForLanguage() got = %v, failed result check", got)
				}
			}
		})
	}
}

func TestRetrieveArticleForLanguage(t *testing.T) {
	tests := []struct {
		name        string
		word        string
		lang        string
		wantErr     bool
		errContains string
		checkResult func(string) bool
	}{
		{
			name:        "invalid language",
			word:        "test",
			lang:        "invalid",
			wantErr:     true,
			errContains: "invalid language",
		},
		{
			name:    "valid english article",
			word:    "test",
			lang:    "en",
			wantErr: false,
			checkResult: func(result string) bool {
				return len(result) > 0 && strings.Contains(result, "test")
			},
		},
		{
			name:        "non-existent word",
			word:        "asdfqwerzxcv",
			lang:        "en",
			wantErr:     true,
			errContains: "API error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := retrieveArticleForLanguage(tt.word, tt.lang)

			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveArticleForLanguage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("retrieveArticleForLanguage() error = %v, should contain %v", err, tt.errContains)
				return
			}

			// Result checking
			if !tt.wantErr && tt.checkResult != nil {
				if !tt.checkResult(got) {
					t.Errorf("retrieveArticleForLanguage() got result that failed validation")
				}
			}
		})
	}
}

func TestSearchWords(t *testing.T) {
	got, err := searchWords("test")
	if err != nil {
		t.Errorf("searchWords() error = %v", err)
		return
	}
	if len(got) == 0 {
		t.Error("searchWords() returned empty result")
	}
}

func TestRetrieveArticle(t *testing.T) {
	got, err := retrieveArticle("test", "en")
	if err != nil {
		t.Errorf("retrieveArticle() error = %v", err)
		return
	}
	if len(got) == 0 {
		t.Error("retrieveArticle() returned empty result")
	}
}

// TestIntegration tests the workflow of searching and then retrieving an article
func TestIntegration(t *testing.T) {
	words, err := searchWords("test")
	if err != nil {
		t.Fatalf("Failed to search words: %v", err)
	}
	if len(words) == 0 {
		t.Fatal("No words found in search")
	}

	article, err := retrieveArticle(words[0], "en")
	if err != nil {
		t.Fatalf("Failed to retrieve article: %v", err)
	}
	if len(article) == 0 {
		t.Fatal("Retrieved empty article")
	}
}
