package main

import "solod.dev/so/time"

func parse() {
	// All tests use variants of 2024-03-15T14:30:45Z as the input time.
	{
		// RFC3339.
		t, err := time.Parse(time.RFC3339, "2024-03-15T14:30:45Z", 0)
		if err != nil {
			panic("unexpected Parse RFC3339 error")
		}
		date := t.Date(time.UTC)
		if date.Year != 2024 || date.Month != time.March || date.Day != 15 {
			panic("unexpected Parse RFC3339 date")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 14 || clock.Minute != 30 || clock.Second != 45 {
			panic("unexpected Parse RFC3339 clock")
		}
	}
	{
		// RFC3339Nano.
		t, err := time.Parse(time.RFC3339Nano, "2024-03-15T14:30:45.123456789Z", 0)
		if err != nil {
			panic("unexpected Parse RFC3339Nano error")
		}
		date := t.Date(time.UTC)
		if date.Year != 2024 || date.Month != time.March || date.Day != 15 {
			panic("unexpected Parse RFC3339Nano date")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 14 || clock.Minute != 30 || clock.Second != 45 {
			panic("unexpected Parse RFC3339Nano clock")
		}
		if t.Nanosecond() != 123456789 {
			panic("unexpected Parse RFC3339Nano nanosecond")
		}
	}
	{
		// RFC3339 with positive offset.
		// 14:30:45+05:00 is 09:30:45 UTC.
		t, err := time.Parse(time.RFC3339, "2024-03-15T14:30:45+05:00", 0)
		if err != nil {
			panic("unexpected Parse RFC3339+offset error")
		}
		date := t.Date(time.UTC)
		if date.Year != 2024 || date.Month != time.March || date.Day != 15 {
			panic("unexpected Parse RFC3339+offset date")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 9 || clock.Minute != 30 || clock.Second != 45 {
			panic("unexpected Parse RFC3339+offset clock")
		}
	}
	{
		// RFC3339 with negative offset.
		// 14:30:45-03:00 is 17:30:45 UTC.
		t, err := time.Parse(time.RFC3339, "2024-03-15T14:30:45-03:00", 0)
		if err != nil {
			panic("unexpected Parse RFC3339-offset error")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 17 || clock.Minute != 30 || clock.Second != 45 {
			panic("unexpected Parse RFC3339-offset clock")
		}
	}
	{
		// RFC3339Nano with offset.
		// 14:30:45+05:30 is 09:00:45 UTC.
		t, err := time.Parse(time.RFC3339Nano, "2024-03-15T14:30:45.123456789+05:30", 0)
		if err != nil {
			panic("unexpected Parse RFC3339Nano+offset error")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 9 || clock.Minute != 0 || clock.Second != 45 {
			panic("unexpected Parse RFC3339Nano+offset clock")
		}
		if t.Nanosecond() != 123456789 {
			panic("unexpected Parse RFC3339Nano+offset nanosecond")
		}
	}
	{
		// DateTime.
		t, err := time.Parse(time.DateTime, "2024-03-15 14:30:45", time.UTC)
		if err != nil {
			panic("unexpected Parse DateTime error")
		}
		date := t.Date(time.UTC)
		if date.Year != 2024 || date.Month != time.March || date.Day != 15 {
			panic("unexpected Parse DateTime date")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 14 || clock.Minute != 30 || clock.Second != 45 {
			panic("unexpected Parse DateTime clock")
		}
	}
	{
		// DateTime with offset parameter.
		// 14:30:45+05:30 is 09:00:45 UTC.
		offset := time.Offset(5*3600 + 30*60) // UTC+5:30
		t, err := time.Parse(time.DateTime, "2024-03-15 14:30:45", offset)
		if err != nil {
			panic("unexpected Parse DateTime+offset error")
		}
		date := t.Date(time.UTC)
		if date.Year != 2024 || date.Month != time.March || date.Day != 15 {
			panic("unexpected Parse DateTime+offset date")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 9 || clock.Minute != 0 || clock.Second != 45 {
			panic("unexpected Parse DateTime+offset clock")
		}
	}
	{
		// DateOnly.
		t, err := time.Parse(time.DateOnly, "2024-03-15", time.UTC)
		if err != nil {
			panic("unexpected Parse DateOnly error")
		}
		date := t.Date(time.UTC)
		if date.Year != 2024 || date.Month != time.March || date.Day != 15 {
			panic("unexpected Parse DateOnly date")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 0 || clock.Minute != 0 || clock.Second != 0 {
			panic("unexpected Parse DateOnly clock")
		}
	}
	{
		// TimeOnly.
		t, err := time.Parse(time.TimeOnly, "14:30:45", time.UTC)
		if err != nil {
			panic("unexpected Parse TimeOnly error")
		}
		date := t.Date(time.UTC)
		if date.Year != 0 || date.Month != time.January || date.Day != 1 {
			panic("unexpected Parse TimeOnly date")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 14 || clock.Minute != 30 || clock.Second != 45 {
			panic("unexpected Parse TimeOnly clock")
		}
	}
	{
		// Custom format.
		t, err := time.Parse("%d.%m.%Y", "15.03.2024", time.UTC)
		if err != nil {
			panic("unexpected Parse custom error")
		}
		date := t.Date(time.UTC)
		if date.Year != 2024 || date.Month != time.March || date.Day != 15 {
			panic("unexpected Parse custom date")
		}
		clock := t.Clock(time.UTC)
		if clock.Hour != 0 || clock.Minute != 0 || clock.Second != 0 {
			panic("unexpected Parse custom clock")
		}
	}
	{
		// time.Parse error.
		_, err := time.Parse("%Y-%m-%d", "not-a-date", time.UTC)
		if err == nil {
			panic("expected Parse error")
		}
	}
}
