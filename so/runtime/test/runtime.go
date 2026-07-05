package main

import (
	"solod.dev/so/runtime"
	"solod.dev/so/testing"
)

func TestVersion(t *testing.T) {
	v := runtime.Version()
	if len(v) == 0 {
		t.Error("Empty version")
	}
}

func TestGOOS(t *testing.T) {
	os := runtime.GOOS
	if os != "bare" && os != "darwin" && os != "linux" && os != "windows" && os != "wasip1" {
		t.Error("Unexpected GOOS")
	}
}

func TestGOARCH(t *testing.T) {
	arch := runtime.GOARCH
	if arch != "amd64" && arch != "arm64" && arch != "386" && arch != "riscv64" && arch != "wasm" {
		t.Error("Unexpected GOARCH")
	}
}
