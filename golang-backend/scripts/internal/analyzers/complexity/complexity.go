package complexity

import (
	"go/ast"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	maxFileLines     = 400
	maxFunctionLines = 60
	maxParams        = 4
)

var Analyzer = &analysis.Analyzer{
	Name: "golangbackend_complexity",
	Doc:  "enforces file and function size constraints",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if ast.IsGenerated(file) || isExemptFile(pass.Fset.Position(file.Pos()).Filename) {
			continue
		}
		start := pass.Fset.Position(file.Package).Line
		end := pass.Fset.Position(file.End()).Line
		if end-start+1 > maxFileLines {
			pass.Reportf(file.Package, "file exceeds %d lines. Split mixed responsibilities into smaller files.", maxFileLines)
		}

		for _, decl := range file.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Body == nil || isConstructor(fd) {
				continue
			}
			checkFunctionLines(pass, fd)
			checkParameterCount(pass, fd)
		}
	}

	return nil, nil
}

func isExemptFile(filename string) bool {
	base := filepath.ToSlash(filename)
	return strings.Contains(base, "/migrations/") || strings.Contains(base, "/fixtures/")
}

func isConstructor(fd *ast.FuncDecl) bool {
	return strings.HasPrefix(fd.Name.Name, "New")
}

func checkFunctionLines(pass *analysis.Pass, fd *ast.FuncDecl) {
	start := pass.Fset.Position(fd.Pos()).Line
	end := pass.Fset.Position(fd.End()).Line
	if end-start+1 > maxFunctionLines {
		pass.Reportf(fd.Name.Pos(), "function %q exceeds %d lines. Split the function into smaller units or move boundary orchestration outward.", fd.Name.Name, maxFunctionLines)
	}
}

func checkParameterCount(pass *analysis.Pass, fd *ast.FuncDecl) {
	count := 0
	for _, field := range fd.Type.Params.List {
		if isContextField(pass, field) {
			continue
		}
		names := len(field.Names)
		if names == 0 {
			names = 1
		}
		count += names
	}
	if count > maxParams {
		pass.Reportf(fd.Name.Pos(), "function %q has %d parameters excluding context.Context. Maximum allowed is %d.", fd.Name.Name, count, maxParams)
	}
}

func isContextField(pass *analysis.Pass, field *ast.Field) bool {
	t := pass.TypesInfo.TypeOf(field.Type)
	named, ok := deref(t).(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}
	return named.Obj().Pkg().Path() == "context" && named.Obj().Name() == "Context"
}

func deref(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}
