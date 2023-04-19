package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tm := time.Date(2023, 4, 19, 12, 56, 0, 0, time.UTC)
	hd := humanDate(tm)

	if hd != "19 Apr 2023 at 12:56" {
		t.Errorf("got %q; want %q", hd, "19 Apr 2023 at 12:56")
	}
}
