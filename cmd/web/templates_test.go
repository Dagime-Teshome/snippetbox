package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	tableTest := []struct {
		name     string
		data     time.Time
		expected string
	}{
		{
			name:     "test empty date",
			data:     time.Time{},
			expected: "",
		},
		{
			name:     "test UTC",
			data:     time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			expected: "17 Dec 2020 at 10:00",
		},
		{
			name:     "CET",
			data:     time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			expected: "17 Dec 2020 at 09:00",
		},
	}

	for _, test := range tableTest {
		t.Run(test.name, func(t *testing.T) {
			testRes := humanDate(test.data)
			if testRes != test.expected {
				t.Errorf("expected %q but got %q", test.expected, testRes)
			}
		})
	}
}
