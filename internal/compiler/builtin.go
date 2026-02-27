package compiler

import (
	"embed"
	"log/slog"
	"os"
	"path/filepath"
)

//go:embed builtin/so.h builtin/so.c
var builtinFS embed.FS

func writeBuiltin(outDir string) {
	for _, name := range []string{"so.h", "so.c"} {
		data, err := builtinFS.ReadFile("builtin/" + name)
		if err != nil {
			slog.Error("failed to read embedded builtin file", "name", name, "error", err)
			os.Exit(1)
		}
		if err := os.WriteFile(filepath.Join(outDir, name), data, 0o644); err != nil {
			slog.Error("failed to write builtin file", "error", err)
			os.Exit(1)
		}
	}
}
