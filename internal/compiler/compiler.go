package compiler

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/nalgeon/soan/internal/clang"
)

// Translate loads all Go packages from dir (including So stdlib dependencies),
// translates them to C, and writes the output to outDir.
func Translate(dir string, outDir string) {
	pkgs := loadPackages(dir)
	if len(pkgs) == 0 {
		slog.Error("no packages found")
		os.Exit(1)
	}

	entry := pkgs[0]

	// Walk import graph and collect transpilable packages in topological order
	ordered := topoSort(entry, entry.Module)

	// Translate each package
	for _, pkg := range ordered {
		pkgOutDir := packageOutDir(pkg, entry, outDir)
		clang.Emit(clang.EmitOptions{
			Pkg:    pkg,
			OutDir: pkgOutDir,
		})
	}

	// Write embedded builtin files (so.h, so.c) into the output directory
	writeBuiltin(outDir)
}

// packageOutDir returns the output directory for a package.
// Entry package goes to outDir directly.
// Other packages strip their module prefix (e.g. github.com/nalgeon/soan/math -> math).
func packageOutDir(pkg, entry *packages.Package, outDir string) string {
	if pkg.PkgPath == entry.PkgPath {
		return outDir
	}
	relPath := strings.TrimPrefix(pkg.PkgPath, pkg.Module.Path+"/")
	return filepath.Join(outDir, relPath)
}

// loadPackages uses go/packages to load the entry package and all dependencies.
func loadPackages(dir string) []*packages.Package {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedImports | packages.NeedDeps |
			packages.NeedModule | packages.NeedTypesInfo,
		Dir: dir,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		slog.Error("failed to load packages", "error", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}
	return pkgs
}

// topoSort walks the import graph from entry and returns transpilable packages
// (module-internal + So stdlib) in topological order (dependencies before dependents).
func topoSort(entry *packages.Package, entryModule *packages.Module) []*packages.Package {
	var ordered []*packages.Package
	visited := make(map[string]bool)

	var walk func(pkg *packages.Package)
	walk = func(pkg *packages.Package) {
		if visited[pkg.PkgPath] {
			return
		}
		visited[pkg.PkgPath] = true

		// Visit dependencies first (post-order)
		for _, dep := range pkg.Imports {
			if shouldTranspile(dep, entryModule) {
				walk(dep)
			}
		}
		ordered = append(ordered, pkg)
	}
	walk(entry)
	return ordered
}

// shouldTranspile returns true if a package should be transpiled to C.
// This includes module-internal packages and So stdlib packages.
func shouldTranspile(pkg *packages.Package, entryModule *packages.Module) bool {
	if isModuleInternal(pkg, entryModule) {
		return true
	}
	return pkg.Module != nil && pkg.Module.Path == getOwnModulePath()
}

// isModuleInternal checks if a package belongs to the same module as the entry package.
func isModuleInternal(pkg *packages.Package, entryModule *packages.Module) bool {
	if entryModule == nil || pkg.Module == nil {
		return false
	}
	return pkg.Module.Path == entryModule.Path
}

// getOwnModulePath returns the module path of the running binary.
func getOwnModulePath() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Main.Path
}
