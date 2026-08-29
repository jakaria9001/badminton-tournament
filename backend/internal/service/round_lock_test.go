package service

import "testing"

func TestCanLockRound(t *testing.T) {
	for _, test := range []struct {
		roundName string
		want      bool
	}{
		{roundName: "FINAL", want: true},
		{roundName: "SEMIFINAL", want: false},
		{roundName: "ROUND_1", want: false},
	} {
		if got := CanLockRound(test.roundName); got != test.want {
			t.Errorf("CanLockRound(%q) = %t, want %t", test.roundName, got, test.want)
		}
	}
}