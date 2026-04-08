package main

import "solod.dev/so/time"

func format() {
	buf := make([]byte, 64)
	t := time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC)
	{
		// RFC3339.
		s := t.Format(buf, time.RFC3339, time.UTC)
		if s != "2024-03-15T14:30:45Z" {
			panic("unexpected RFC3339 format")
		}
	}
	{
		// RFC3339Nano.
		t = time.Date(2024, time.March, 15, 14, 30, 45, 123456789, time.UTC)
		s := t.Format(buf, time.RFC3339Nano, time.UTC)
		if s != "2024-03-15T14:30:45.123456789Z" {
			panic("unexpected RFC3339Nano format")
		}
	}
	{
		// DateTime.
		s := t.Format(buf, time.DateTime, time.UTC)
		if s != "2024-03-15 14:30:45" {
			panic("unexpected DateTime format")
		}
	}
	{
		// DateOnly.
		s := t.Format(buf, time.DateOnly, time.UTC)
		if s != "2024-03-15" {
			panic("unexpected DateOnly format")
		}
	}
	{
		// TimeOnly.
		s := t.Format(buf, time.TimeOnly, time.UTC)
		if s != "14:30:45" {
			panic("unexpected TimeOnly format")
		}
	}
	{
		// Custom format.
		s := t.Format(buf, "%d.%m.%Y", time.UTC)
		if s != "15.03.2024" {
			panic("unexpected custom format")
		}
	}
	{
		// Time.String.
		s := t.String(buf)
		if s != "2024-03-15T14:30:45Z" {
			panic("unexpected String format")
		}
	}
}
