package debug

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewDebugDefaults(t *testing.T) {
	c := &cobra.Command{Use: "ip"}
	d := NewDebug(c)
	if d.Mode != RELEASE {
		t.Fatalf("NewDebug without debug flag: Mode = %v, want RELEASE", d.Mode)
	}
	if !d.CombinationKey {
		t.Fatal("NewDebug: CombinationKey = false, want true")
	}
	if !d.AutoAppendEntry {
		t.Fatal("NewDebug: AutoAppendEntry = false, want true")
	}
	if d.Command != c {
		t.Fatal("NewDebug: Command not wired")
	}
}

func TestNewDebugWithDebugFlag(t *testing.T) {
	c := &cobra.Command{Use: "ip"}
	c.Flags().Bool("debug", false, "debug mode")
	if err := c.Flags().Set("debug", "true"); err != nil {
		t.Fatalf("set debug flag: %v", err)
	}
	d := NewDebug(c)
	if d.Mode != DEVELOP {
		t.Fatalf("debug flag enabled: Mode = %v, want DEVELOP", d.Mode)
	}
}

func TestEnableDebugWithoutFlag(t *testing.T) {
	c := &cobra.Command{Use: "ip"}
	d := &Debug{Command: c}
	if d.EnableDebug() {
		t.Fatal("EnableDebug without flag = true, want false")
	}
}

func TestAppendEntry(t *testing.T) {
	cases := []struct {
		name   string
		entry  bool
		source string
		want   string
	}{
		{"appends newline", true, "msg", "msg\n"},
		{"keeps source", false, "msg", "msg"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := &Debug{AutoAppendEntry: tc.entry}
			if got := d.AppendEntry(tc.source); got != tc.want {
				t.Fatalf("AppendEntry(%q) = %q, want %q", tc.source, got, tc.want)
			}
		})
	}
}

func TestFillSpace(t *testing.T) {
	d := &Debug{}
	cases := []struct {
		name   string
		width  int
		middle string
		value  string
		want   string
	}{
		{"pads short value", 10, " : ", "admin", "     admin : "},
		{"pads single space", 10, " : ", "wmsjhappy", " wmsjhappy : "},
		{"no pad when exact", 5, "|", "abcde", "abcde|"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := d.FillSpace(tc.width, tc.middle, tc.value); got != tc.want {
				t.Fatalf("FillSpace(%d, %q, %q) = %q, want %q", tc.width, tc.middle, tc.value, got, tc.want)
			}
		})
	}
}
