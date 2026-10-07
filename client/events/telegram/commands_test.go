package telegram_test

import (
	"errors"
	"strings"
	"testing"

	telegram "github.com/Zhenka07/TelegramBot/client/events/telegram"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr error
	}{
		{
			name:    "valid https url",
			url:     "https://github.com/obra/superpowers",
			wantErr: nil,
		},
		{
			name:    "valid http url with query params",
			url:     "http://example.com/path?foo=bar&num=42",
			wantErr: nil,
		},
		{
			name:    "invalid scheme ftp",
			url:     "ftp://example.com/files",
			wantErr: telegram.ErrInvalidScheme,
		},
		{
			name:    "invalid scheme javascript",
			url:     "javascript:alert(1)",
			wantErr: telegram.ErrInvalidScheme,
		},
		{
			name:    "url too long (over 2048 chars)",
			url:     "https://example.com/" + strings.Repeat("a", 2040),
			wantErr: telegram.ErrURLTooLong,
		},
		{
			name:    "forbidden localhost",
			url:     "http://localhost:8080/admin",
			wantErr: telegram.ErrForbiddenHost,
		},
		{
			name:    "forbidden 127.0.0.1",
			url:     "http://127.0.0.1:9090",
			wantErr: telegram.ErrForbiddenHost,
		},
		{
			name:    "missing host",
			url:     "http:///no-host-here",
			wantErr: telegram.ErrInvalidHost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := telegram.ValidateURL(tt.url)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
			} else {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got: %v", tt.wantErr, err)
				}
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	if !telegram.IsURL("https://example.com") {
		t.Errorf("expected 'https://example.com' to be a valid URL")
	}
	if telegram.IsURL("http://localhost:8080") {
		t.Errorf("expected localhost to be rejected as valid URL")
	}
	if telegram.IsURL("not-a-url") {
		t.Errorf("expected 'not-a-url' to be false")
	}
}
