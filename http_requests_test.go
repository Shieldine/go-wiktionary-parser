package go_wiktionary_parser

import (
	"testing"
)

func TestSearchWordsForLanguage(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		lang        string
		wantErr     bool
		errMessage  string
		checkResult bool // only check result if true
		wantLen     int  // minimum expected length of results
	}{
		{
			name:        "valid english search",
			query:       "test",
			lang:        "en",
			wantErr:     false,
			checkResult: true,
			wantLen:     1,
		},
		{
			name:        "valid german search",
			query:       "haus",
			lang:        "de",
			wantErr:     false,
			checkResult: true,
			wantLen:     1,
		},
		{
			name:       "empty query",
			query:      "",
			lang:       "en",
			wantErr:    true,
			errMessage: "empty query",
		},
		{
			name:       "invalid language",
			query:      "test",
			lang:       "xx",
			wantErr:    true,
			errMessage: "invalid language",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchWordsForLanguage(tt.query, tt.lang)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SearchWordsForLanguage() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errMessage {
					t.Errorf("SearchWordsForLanguage() error = %v, wantErr %v", err, tt.errMessage)
				}
				return
			}

			if err != nil {
				t.Errorf("SearchWordsForLanguage() unexpected error = %v", err)
				return
			}

			// Check results if required
			if tt.checkResult {
				if len(got) < tt.wantLen {
					t.Errorf("SearchWordsForLanguage() got %d results, want at least %d", len(got), tt.wantLen)
				}
			}
		})
	}
}

func TestRetrieveArticleForLanguage(t *testing.T) {
	tests := []struct {
		name       string
		word       string
		lang       string
		wantErr    bool
		errMessage string
		want       *ArticleContent
	}{
		{
			name:    "valid english article",
			word:    "test",
			lang:    "en",
			wantErr: false,
		},
		{
			name:    "valid german article",
			word:    "Haus",
			lang:    "de",
			wantErr: false,
		},
		{
			name:       "invalid language",
			word:       "test",
			lang:       "xx",
			wantErr:    true,
			errMessage: "invalid language",
		},
		{
			name:       "nonexistent word",
			word:       "asdfasdfasdfasdf",
			lang:       "en",
			wantErr:    true,
			errMessage: "API error:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RetrieveArticleForLanguage(tt.word, tt.lang)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RetrieveArticleForLanguage() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMessage != "" && err.Error() != tt.errMessage {
					if tt.errMessage == "API error:" && !contains(err.Error(), tt.errMessage) {
						t.Errorf("RetrieveArticleForLanguage() error = %v, want error containing %v", err, tt.errMessage)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("RetrieveArticleForLanguage() unexpected error = %v", err)
				return
			}

			// Basic validation of returned article
			if got.Title == "" {
				t.Error("RetrieveArticleForLanguage() returned article with empty title")
			}
			if got.HTML == "" {
				t.Error("RetrieveArticleForLanguage() returned article with empty HTML")
			}
			if got.Language != tt.lang {
				t.Errorf("RetrieveArticleForLanguage() returned article with wrong language = %v, want %v", got.Language, tt.lang)
			}
		})
	}
}

func TestWrapper_Functions(t *testing.T) {
	// Test SearchWords (English wrapper)
	t.Run("SearchWords", func(t *testing.T) {
		got, err := SearchWords("test")
		if err != nil {
			t.Errorf("SearchWords() error = %v", err)
			return
		}
		if len(got) == 0 {
			t.Error("SearchWords() returned empty result")
		}
	})

	// Test RetrieveArticle (English wrapper)
	t.Run("RetrieveArticle", func(t *testing.T) {
		got, err := RetrieveArticle("test")
		if err != nil {
			t.Errorf("RetrieveArticle() error = %v", err)
			return
		}
		if got.Title == "" || got.HTML == "" {
			t.Error("RetrieveArticle() returned incomplete article")
		}
	})
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
