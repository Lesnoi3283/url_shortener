package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShortenURL(t *testing.T) {
	type args struct {
		url []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				url: []byte("https://googlegooglegooglegoogle.com"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShortenURL(tt.args.url)
			for _, b := range result {
				isLetter := (b <= 'z') && (b >= 'A')
				require.Truef(t, isLetter, "found a strange byte (not a letter) '%c' in str '%s'", b, result)
			}
		})
	}
}
