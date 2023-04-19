package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 4, 19, 10, 01, 0, 0, time.UTC),
			want: "19 Apr 2023 at 10:01",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "IRST",
			tm:   time.Date(2023, 4, 19, 10, 01, 0, 0, time.FixedZone("IRST", -(3*60*60+30*60))),
			want: "19 Apr 2023 at 13:31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			if hd != tt.want {
				t.Errorf("got %q; want %q", hd, tt.want)
			}
		})
	}

	tm := time.Date(2023, 4, 19, 12, 56, 0, 0, time.UTC)
	hd := humanDate(tm)

	if hd != "19 Apr 2023 at 12:56" {
		t.Errorf("got %q; want %q", hd, "19 Apr 2023 at 12:56")
	}
}
