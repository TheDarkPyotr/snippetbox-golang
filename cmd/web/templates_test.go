package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name  string
		given time.Time
		want  string
	}{
		{
			name:  "UTC",
			given: time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want:  "17 Dec 2020 at 10:00",
		},
		{

			name:  "Empty",
			given: time.Time{},
			want:  "",
		},
		/*
			{
				name:  "CET",
				given: time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60)),
				want:  "17 Dec 2020 at 09:00",
			},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hd := humanDate(tt.given)

			if hd != tt.want {
				t.Errorf("want %q got %q", tt.want, hd)
			}

		})
	}

}
