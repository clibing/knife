package utils

import "testing"

func TestBeautifyUnit(t *testing.T) {
	cases := []struct {
		name  string
		in    int64
		unit  string
		value string
	}{
		{"zero bytes", 0, "B", "0.00"},
		{"under 1KB stays B", KB - 1, "B", "1023.00"},
		{"exactly 1KB", KB, "KB", "1.00"},
		{"kilobytes", 2 * KB, "KB", "2.00"},
		{"megabytes", 3 * MB, "MB", "3.00"},
		{"gigabytes", 5 * GB, "GB", "5.00"},
		{"terabytes", 2 * TB, "TB", "2.00"},
		{"petabytes", 1 * PB, "PB", "1.00"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			unit, value := BeautifyUnit(tc.in)
			if unit != tc.unit || value != tc.value {
				t.Fatalf("BeautifyUnit(%d) = (%q, %q), want (%q, %q)", tc.in, unit, value, tc.unit, tc.value)
			}
		})
	}
}

func TestBeautifyValue(t *testing.T) {
	cases := []struct {
		name string
		in   int64
		want string
	}{
		{"bytes", 512, "512.00 B"},
		{"kilobytes", 2048, "2.00 KB"},
		{"gigabytes", 5 * GB, "5.00 GB"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := BeautifyValue(tc.in); got != tc.want {
				t.Fatalf("BeautifyValue(%d) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
