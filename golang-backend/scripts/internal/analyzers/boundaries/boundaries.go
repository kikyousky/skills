package boundaries

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/internalutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "golangbackend_boundaries",
	Doc:  "enforces backend package boundaries and naming suffixes",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	pkgPath := pass.Pkg.Path()

	for _, file := range pass.Files {
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, "\"")
			switch {
			case internalutil.IsDomainPath(pkgPath) && (internalutil.IsTransportPath(path) || internalutil.IsServicePath(path)):
				pass.Reportf(imp.Pos(), "domain package must not import %q. Keep domain isolated from transport and service.", path)
			case internalutil.IsServicePath(pkgPath) && internalutil.IsTransportPath(path):
				pass.Reportf(imp.Pos(), "service package must not import %q. Accept only validated service inputs.", path)
			case internalutil.IsServicePath(pkgPath) && internalutil.IsTransportFrameworkImport(path):
				pass.Reportf(imp.Pos(), "service package must not import transport framework %q. Keep delivery concerns in internal/transport.", path)
			case internalutil.IsRepoPath(pkgPath) && (internalutil.IsTransportPath(path) || internalutil.IsServicePath(path)):
				pass.Reportf(imp.Pos(), "repo package must not import %q. Persist domain concepts, not transport or service-layer types.", path)
			}
		}

		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if internalutil.HasForbiddenSuffix(ts.Name.Name) {
					pass.Reportf(ts.Name.Pos(), "type %q uses forbidden suffix. Use a precise domain-specific name instead of Manager, Helper, Util, or Data.", ts.Name.Name)
				}
			}
		}
	}

	return nil, nil
}
