package service

import "testing"

func TestValidRegistrationTransition(t *testing.T) {
	tests := []struct {
		current string
		next    string
		valid   bool
	}{
		{"PENDING", "CONFIRMED", true},
		{"PENDING", "REJECTED", true},
		{"PENDING", "WITHDRAWN", true},
		{"CONFIRMED", "WITHDRAWN", true},
		{"CONFIRMED", "REJECTED", false},
		{"REJECTED", "CONFIRMED", false},
		{"WITHDRAWN", "CONFIRMED", false},
	}

	for _, test := range tests {
		if got := ValidRegistrationTransition(test.current, test.next); got != test.valid {
			t.Errorf("ValidRegistrationTransition(%q, %q) = %t, want %t", test.current, test.next, got, test.valid)
		}
	}
}