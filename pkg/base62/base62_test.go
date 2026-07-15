package base62

import (
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		id       uint64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{61, "Z"},
		{62, "10"},
		{100, "1C"},
		{999, "g7"},
		{123456789, "8m0Kx"},
	}

	for _, tc := range tests {
		actual := Encode(tc.id)
		if actual != tc.expected {
			t.Errorf("Encode(%d) = %s; expected %s", tc.id, actual, tc.expected)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		encoded  string
		expected uint64
		hasErr   bool
	}{
		{"0", 0, false},
		{"1", 1, false},
		{"Z", 61, false},
		{"10", 62, false},
		{"1C", 100, false},
		{"g7", 999, false},
		{"8m0Kx", 123456789, false},
		{"invalid-char!", 0, true},
	}

	for _, tc := range tests {
		actual, err := Decode(tc.encoded)
		if tc.hasErr {
			if err == nil {
				t.Errorf("Decode(%s) expected error but got none", tc.encoded)
			}
		} else {
			if err != nil {
				t.Errorf("Decode(%s) unexpected error: %v", tc.encoded, err)
			}
			if actual != tc.expected {
				t.Errorf("Decode(%s) = %d; expected %d", tc.encoded, actual, tc.expected)
			}
		}
	}
}
