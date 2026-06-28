package serviceparams

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/internalutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "golangbackend_serviceparams",
	Doc:  "enforces service function parameter restrictions",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	if !internalutil.IsServicePath(pass.Pkg.Path()) {
		return nil, nil
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Type == nil || fd.Type.Params == nil || isConstructor(fd) {
				continue
			}

			for _, field := range fd.Type.Params.List {
				paramType := pass.TypesInfo.TypeOf(field.Type)
				if paramType == nil || isContext(paramType) {
					continue
				}
				if isAllowedServiceParam(pass.Pkg, paramType) {
					continue
				}
				pass.Reportf(field.Pos(), "service function %q accepts forbidden parameter type %s. Accept only context.Context, service.*Input, or service-owned interfaces.", fd.Name.Name, types.TypeString(paramType, qualifier))
			}
		}
	}

	return nil, nil
}

func isConstructor(fd *ast.FuncDecl) bool {
	return strings.HasPrefix(fd.Name.Name, "New")
}

func isContext(t types.Type) bool {
	named, ok := deref(t).(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}
	return named.Obj().Pkg().Path() == "context" && named.Obj().Name() == "Context"
}

func isAllowedServiceParam(pkg *types.Package, t types.Type) bool {
	if iface, ok := deref(t).Underlying().(*types.Interface); ok {
		_ = iface
		named, ok := deref(t).(*types.Named)
		if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
			return false
		}
		return named.Obj().Pkg().Path() == pkg.Path()
	}

	named, ok := deref(t).(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}
	return named.Obj().Pkg().Path() == pkg.Path() && strings.HasSuffix(named.Obj().Name(), "Input")
}

func deref(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}

func qualifier(other *types.Package) string {
	if other == nil {
		return ""
	}
	return other.Name()
}
