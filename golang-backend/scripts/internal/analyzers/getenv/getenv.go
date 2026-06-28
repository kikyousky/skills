package getenv

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/internalutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "golangbackend_getenv",
	Doc:  "forbids raw os.Getenv outside internal/config",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	if internalutil.IsConfigPath(pass.Pkg.Path()) {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			x, ok := sel.X.(*ast.Ident)
			if ok && x.Name == "os" && sel.Sel.Name == "Getenv" {
				pass.Reportf(call.Pos(), "raw os.Getenv usage is forbidden outside internal/config. Load env vars in internal/config, validate them, and pass typed config inward.")
			}
			return true
		})
	}

	return nil, nil
}
