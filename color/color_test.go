package color

import (
	"fmt"
	"testing"
)

func TestText(t *testing.T) {
	tests := []struct {
		code    int
		value   interface{}
		want    string
		hasAnsi bool
	}{
		{0, "hello", "\u001b[38;5;0mhello\u001b[0m", true},
		{1, "world", "\u001b[38;5;1mworld\u001b[0m", true},
		{255, "!", "\u001b[38;5;255m!\u001b[0m", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d-%v", test.code, test.value), func(t *testing.T) {
			got := Text(test.code, test.value)
			if got != test.want {
				t.Errorf("Text(%d, %v) = %q, want %q", test.code, test.value, got, test.want)
			}
			if (len(got) > 4) != test.hasAnsi {
				t.Errorf("Text(%d, %v) has ANSI = %t, want %t", test.code, test.value, len(got) > 4, test.hasAnsi)
			}
		})
	}
}

func TestBg(t *testing.T) {
	tests := []struct {
		code    int
		value   interface{}
		want    string
		hasAnsi bool
	}{
		{0, "hello", "\u001b[48;5;0mhello\u001b[0m", true},
		{1, "world", "\u001b[48;5;1mworld\u001b[0m", true},
		{255, "!", "\u001b[48;5;255m!\u001b[0m", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d-%v", test.code, test.value), func(t *testing.T) {
			got := Bg(test.code, test.value)
			if got != test.want {
				t.Errorf("Bg(%d, %v) = %q, want %q", test.code, test.value, got, test.want)
			}
			if (len(got) > 4) != test.hasAnsi {
				t.Errorf("Bg(%d, %v) has ANSI = %t, want %t", test.code, test.value, len(got) > 4, test.hasAnsi)
			}
		})
	}
}
