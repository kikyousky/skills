package jsondecode

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "golangbackend_jsondecode",
	Doc:  "requires HTTP JSON decoding to use json.Decoder with DisallowUnknownFields",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	if !strings.Contains(pass.Pkg.Path(), "/internal/transport/http") {
		return nil, nil
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}
			checkFunc(pass, fd)
		}
	}

	return nil, nil
}

func checkFunc(pass *analysis.Pass, fd *ast.FuncDecl) {
	configured := map[string]bool{}
	decoders := map[string]token.Pos{}

	for _, stmt := range fd.Body.List {
		ast.Inspect(stmt, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.AssignStmt:
				if len(n.Lhs) != 1 || len(n.Rhs) != 1 {
					return true
				}
				ident, ok := n.Lhs[0].(*ast.Ident)
				if !ok {
					return true
				}
				if isJSONNewDecoder(n.Rhs[0]) {
					decoders[ident.Name] = ident.Pos()
				}
			case *ast.CallExpr:
				if isJSONUnmarshal(n) {
					pass.Reportf(n.Pos(), "json.Unmarshal is forbidden for HTTP boundary decoding. Use json.NewDecoder plus DisallowUnknownFields instead.")
				}
				if recv, ok := selectorReceiver(n, "DisallowUnknownFields"); ok {
					configured[recv] = true
				}
				if recv, ok := selectorReceiver(n, "Decode"); ok {
					if _, exists := decoders[recv]; exists && !configured[recv] {
						pass.Reportf(n.Pos(), "decoder %q must call DisallowUnknownFields() before Decode(). Reject unknown HTTP JSON fields at the boundary.", recv)
					}
				}
				if isInlineDecode(n) {
					pass.Reportf(n.Pos(), "inline json.NewDecoder(...).Decode(...) is forbidden for HTTP boundary decoding. Create a decoder, call DisallowUnknownFields(), then Decode().")
				}
			}
			return true
		})
	}
}

func isJSONNewDecoder(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	x, ok := sel.X.(*ast.Ident)
	return ok && x.Name == "json" && sel.Sel.Name == "NewDecoder"
}

func isJSONUnmarshal(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	x, ok := sel.X.(*ast.Ident)
	return ok && x.Name == "json" && sel.Sel.Name == "Unmarshal"
}

func selectorReceiver(call *ast.CallExpr, method string) (string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != method {
		return "", false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", false
	}
	return ident.Name, true
}

func isInlineDecode(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Decode" {
		return false
	}
	return isJSONNewDecoder(sel.X)
}
