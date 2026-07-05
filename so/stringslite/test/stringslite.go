package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/stringslite"
	"solod.dev/so/testing"
)

func TestClone(t *testing.T) {
	s := "hello"
	c := stringslite.Clone(mem.System, s)
	defer mem.FreeString(mem.System, c)
	if c != s {
		t.Error("Clone failed")
	}
}

func TestCut(t *testing.T) {
	before, after := stringslite.Cut("hello world", " ")
	if before != "hello" || after != "world" {
		t.Error("Cut failed")
	}
}

func TestCutPrefix(t *testing.T) {
	after, found := stringslite.CutPrefix("hello world", "hello ")
	if after != "world" || !found {
		t.Error("CutPrefix failed")
	}
}

func TestCutSuffix(t *testing.T) {
	before, found := stringslite.CutSuffix("hello world", " world")
	if before != "hello" || !found {
		t.Error("CutSuffix failed")
	}
}

func TestHasPrefix(t *testing.T) {
	if !stringslite.HasPrefix("hello world", "hello") {
		t.Error("HasPrefix failed")
	}
	if stringslite.HasPrefix("hello world", "world") {
		t.Error("HasPrefix failed")
	}
}

func TestHasSuffix(t *testing.T) {
	if !stringslite.HasSuffix("hello world", "world") {
		t.Error("HasSuffix failed")
	}
	if stringslite.HasSuffix("hello world", "hello") {
		t.Error("HasSuffix failed")
	}
}

func TestIndex(t *testing.T) {
	idx := stringslite.Index("hello world", "world")
	if idx != 6 {
		t.Error("Index failed")
	}
}

func TestIndexByte(t *testing.T) {
	idx := stringslite.IndexByte("hello world", 'o')
	if idx != 4 {
		t.Error("IndexByte failed")
	}
}

func TestTrimPrefix(t *testing.T) {
	s := stringslite.TrimPrefix("hello world", "hello ")
	if s != "world" {
		t.Error("TrimPrefix failed")
	}
}

func TestTrimSuffix(t *testing.T) {
	s := stringslite.TrimSuffix("hello world", " world")
	if s != "hello" {
		t.Error("TrimSuffix failed")
	}
}
