package main

import (
	"solod.dev/so/time"
	"solod.dev/so/testing"
)

func TestDate(t *testing.T) {
	tm := time.Date(2021, time.May, 10, 12, 33, 44, 777888999, time.UTC)
	if tm.Year() != 2021 {
		t.Error("unexpected Time.Year")
	}
	if tm.Month() != time.May {
		t.Error("unexpected Time.Month")
	}
	if tm.Day() != 10 {
		t.Error("unexpected Time.Day")
	}
	if tm.Hour() != 12 {
		t.Error("unexpected Time.Hour")
	}
	if tm.Minute() != 33 {
		t.Error("unexpected Time.Minute")
	}
	if tm.Second() != 44 {
		t.Error("unexpected Time.Second")
	}
	if tm.Nanosecond() != 777888999 {
		t.Error("unexpected Time.Nanosecond")
	}
}

func TestNow(t *testing.T) {
	tm := time.Now()
	if tm.IsZero() {
		t.Error("unexpected Time.IsZero")
	}
}
