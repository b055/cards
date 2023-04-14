package models

import (
	"errors"
	"testing"
)

func Test_CodeToSuitValue_Invalid(t *testing.T) {
	for _, code := range []string{"Ak", "AK", "aa", "ab", "1", "1111"} {
		_, _, err := CodeToSuitValue(code)
		if err == nil {
			t.Fatalf(`CodeToSuitValue(%q) = _, _, %v, want _, _, %v`, code, err, errors.New("invalid suit "+code[:1]))
		}
	}
}

func Test_CodeToSuitValue_Valid(t *testing.T) {
	for _, code := range []string{"AC", "10H", "JD", "8S"} {
		_, _, err := CodeToSuitValue(code)
		if err != nil {
			t.Fatalf(`CodeToSuitValue(%q) = _, _, %v, want _, _, %v`, code, err, errors.New("invalid suit "+code[:1]))
		}
	}
}

func Test_ValueString_Valid(t *testing.T) {
	expected := []string{"2", "Jack"}
	for i, v := range []int{2, 11} {
		result := Value(v)
		if result.String() != expected[i] {
			t.Fatalf(`Value(%q) = %v, want %v`, v, expected, result)
		}

	}
}

func Test_ValueString_Invalid(t *testing.T) {
	expected := []string{"unknown"}
	for i, v := range []int{100} {
		result := Value(v)
		if result.String() != expected[i] {
			t.Fatalf(`Value(%q) = %v, want %v`, v, expected, result)
		}

	}
}
