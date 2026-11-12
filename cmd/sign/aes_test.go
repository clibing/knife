package sign

import (
	"bytes"
	"testing"
)

func TestAesIV(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		key       string
		blockSize int
		want      []byte
		wantErr   bool
	}{
		{"defaults to key", "", "1234567890abcedf", 16, []byte("1234567890abcedf"), false},
		{"defaults to key prefix when key longer", "", "1234567890abcdef12345678", 16, []byte("1234567890abcdef"), false},
		{"uses input when exact block size", "0234567890abcedf", "1234567890abcedf", 16, []byte("0234567890abcedf"), false},
		{"truncates long input", "0123456789abcdefXXXX", "1234567890abcedf", 16, []byte("0123456789abcdef"), false},
		{"rejects short input", "short", "1234567890abcedf", 16, nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := aesIV(tc.input, tc.key, tc.blockSize)
			if (err != nil) != tc.wantErr {
				t.Fatalf("aesIV(%q, %q, %d) error = %v, wantErr %v", tc.input, tc.key, tc.blockSize, err, tc.wantErr)
			}
			if !tc.wantErr && !bytes.Equal(got, tc.want) {
				t.Fatalf("aesIV(%q, %q, %d) = %q, want %q", tc.input, tc.key, tc.blockSize, got, tc.want)
			}
		})
	}
}
