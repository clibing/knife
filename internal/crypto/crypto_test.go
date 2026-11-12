package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestIsHex(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"0170341a68ba2d86a951338b02a539b3", true},
		{"ABCDEF", true},
		{"abcdef", true},
		{"", true},
		{"4JkQGN1ReWXHq4MRrwiFtw==", false},
		{"clibing", false},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			if got := IsHex([]byte(tc.in)); got != tc.want {
				t.Fatalf("IsHex(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestPKCS5RoundTrip(t *testing.T) {
	plains := [][]byte{[]byte("a"), []byte("clibing"), []byte("0123456789abcdef"), {}}
	for _, plain := range plains {
		padded := PKCS5Padding(plain, BlockSize)
		if len(padded) == 0 || len(padded)%BlockSize != 0 {
			t.Fatalf("PKCS5Padding(%q) length %d is not a positive multiple of %d", plain, len(padded), BlockSize)
		}
		got, err := PKCS5UnPadding(padded)
		if err != nil {
			t.Fatalf("PKCS5UnPadding(%q): %v", plain, err)
		}
		if !bytes.Equal(got, plain) {
			t.Fatalf("round trip = %q, want %q", got, plain)
		}
	}
}

func TestPKCS5UnPaddingErrors(t *testing.T) {
	if _, err := PKCS5UnPadding(nil); err == nil {
		t.Fatal("empty input: want error")
	}
	if _, err := PKCS5UnPadding([]byte{200}); err == nil {
		t.Fatal("padding byte larger than data length: want error")
	}
}

func TestGenerateKey(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []byte
	}{
		{"16 bytes kept as-is", "1234567890abcedf", []byte("1234567890abcedf")},
		{"short key zero padded", "abc", append([]byte("abc"), make([]byte, 13)...)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := GenerateKey([]byte(tc.in)); !bytes.Equal(got, tc.want) {
				t.Fatalf("GenerateKey(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}

	t.Run("long key folded with xor", func(t *testing.T) {
		key := []byte("1234567890abcdefABCDEFGH") // 24 bytes
		got := GenerateKey(key)
		for j := 0; j < 8; j++ {
			if got[j] != key[j]^key[16+j] {
				t.Fatalf("folded byte %d = %d, want %d", j, got[j], key[j]^key[16+j])
			}
		}
		if !bytes.Equal(got[8:16], key[8:16]) {
			t.Fatalf("bytes 8..16 = %q, want untouched %q", got[8:16], key[8:16])
		}
	})
}

func TestECBRoundTrip(t *testing.T) {
	key := []byte("1234567890abcedf")
	plains := []string{"a", "clibing", "0123456789abcdef", "一个中文的测试"}
	for _, plain := range plains {
		encrypted := ECBEncrypt([]byte(plain), key)
		if len(encrypted)%BlockSize != 0 {
			t.Fatalf("ECBEncrypt(%q) length %d not multiple of %d", plain, len(encrypted), BlockSize)
		}
		got := string(ECBDecrypt(encrypted, key))
		if got != plain {
			t.Fatalf("ECB round trip = %q, want %q", got, plain)
		}
	}
}

func TestECBKnownAnswer(t *testing.T) {
	// 来自 knife sign aes -k 1234567890abcedf -m ECB clibing 的历史输出
	got := hex.EncodeToString(ECBEncrypt([]byte("clibing"), []byte("1234567890abcedf")))
	want := "867a374e5caee0f4c8462a658cc431b7"
	if got != want {
		t.Fatalf("ECBEncrypt known answer = %s, want %s", got, want)
	}
}

func TestCBCRoundTrip(t *testing.T) {
	keys := []string{
		"1234567890abcedf",
		"1234567890abcdef12345678",
		"1234567890abcdef1234567890abcdef",
	}
	iv := "0123456789abcdef"
	plains := []string{"a", "clibing", "exactly-16-bytes", "中文内容也可以"}
	for _, key := range keys {
		for _, plain := range plains {
			encrypted, err := CBCEncrypt([]byte(plain), []byte(key), []byte(iv))
			if err != nil {
				t.Fatalf("CBCEncrypt(key=%q, plain=%q): %v", key, plain, err)
			}
			got, err := CBCDecrypt(encrypted, []byte(key), []byte(iv))
			if err != nil {
				t.Fatalf("CBCDecrypt(key=%q, plain=%q): %v", key, plain, err)
			}
			if string(got) != plain {
				t.Fatalf("CBC round trip = %q, want %q", got, plain)
			}
		}
	}
}

func TestCBCErrors(t *testing.T) {
	key := []byte("1234567890abcedf")
	iv := []byte("1234567890abcedf")

	cases := []struct {
		name string
		act  func() error
	}{
		{"encrypt rejects short key", func() error {
			_, err := CBCEncrypt([]byte("x"), []byte("short"), iv)
			return err
		}},
		{"encrypt rejects wrong iv length", func() error {
			_, err := CBCEncrypt([]byte("x"), key, []byte("shortiv"))
			return err
		}},
		{"decrypt rejects non-multiple length", func() error {
			_, err := CBCDecrypt([]byte("not-a-multiple-of-16"), key, iv)
			return err
		}},
		{"decrypt rejects empty ciphertext", func() error {
			_, err := CBCDecrypt(nil, key, iv)
			return err
		}},
		{"decrypt rejects invalid padding", func() error {
			_, err := CBCDecrypt(make([]byte, BlockSize), key, iv)
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.act(); err == nil {
				t.Fatal("want error")
			}
		})
	}
}

func TestCFBRoundTrip(t *testing.T) {
	key := []byte("1234567890abcedf")
	for _, plain := range []string{"a", "clibing", "中文"} {
		encrypted, err := CFBEncrypt([]byte(plain), key)
		if err != nil {
			t.Fatalf("CFBEncrypt(%q): %v", plain, err)
		}
		got, err := CFBDecrypt(encrypted, key)
		if err != nil {
			t.Fatalf("CFBDecrypt(%q): %v", plain, err)
		}
		if string(got) != plain {
			t.Fatalf("CFB round trip = %q, want %q", got, plain)
		}
	}

	e1, _ := CFBEncrypt([]byte("clibing"), key)
	e2, _ := CFBEncrypt([]byte("clibing"), key)
	if bytes.Equal(e1, e2) {
		t.Fatal("random IV should produce different ciphertexts")
	}
}

func TestCFBErrors(t *testing.T) {
	key := []byte("1234567890abcedf")
	if _, err := CFBDecrypt([]byte("short"), key); err == nil {
		t.Fatal("short ciphertext: want error")
	}
	if _, err := CFBEncrypt([]byte("x"), []byte("badkey")); err == nil {
		t.Fatal("invalid key: want error")
	}
}
