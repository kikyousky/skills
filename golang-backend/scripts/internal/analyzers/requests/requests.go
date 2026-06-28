package requests

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/internalutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "golangbackend_requests",
	Doc:  "enforces request type placement and ValidateAndSanitize signatures",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	pkgPath := pass.Pkg.Path()
	for _, file := range pass.Files {
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
				if _, ok := ts.Type.(*ast.StructType); !ok {
					continue
				}
				if strings.HasSuffix(ts.Name.Name, "Request") {
					if !internalutil.IsTransportPath(pkgPath) {
						pass.Reportf(ts.Name.Pos(), "request type %q is declared outside internal/transport. Move boundary request structs under internal/transport/...", ts.Name.Name)
						continue
					}
					checkValidateAndSanitize(pass, ts)
					continue
				}

				if internalutil.IsTransportPath(pkgPath) && looksLikeBoundaryInput(ts.Type) {
					pass.Reportf(ts.Name.Pos(), "transport input struct %q must be named with the Request suffix. All external input must use *Request types.", ts.Name.Name)
				}
			}
		}
	}

	return nil, nil
}

func looksLikeBoundaryInput(expr ast.Expr) bool {
	_, ok := expr.(*ast.StructType)
	return ok
}

func checkValidateAndSanitize(pass *analysis.Pass, ts *ast.TypeSpec) {
	obj := pass.TypesInfo.Defs[ts.Name]
	if obj == nil {
		return
	}
	named, ok := obj.Type().(*types.Named)
	if !ok {
		return
	}

	method := lookupMethod(named, "ValidateAndSanitize")
	if method == nil {
		pass.Reportf(ts.Name.Pos(), "request %q does not implement ValidateAndSanitize() (service.XInput, error). Boundary request structs must validate and normalize input before entering internal/service.", ts.Name.Name)
		return
	}

	sig, ok := method.Type().(*types.Signature)
	if !ok {
		pass.Reportf(method.Pos(), "ValidateAndSanitize on %q must be a method with signature (service.XInput, error).", ts.Name.Name)
		return
	}
	if sig.Params().Len() != 0 || sig.Results().Len() != 2 {
		pass.Reportf(method.Pos(), "ValidateAndSanitize on %q must return (service.XInput, error).", ts.Name.Name)
		return
	}

	resultType := sig.Results().At(0).Type()
	resultNamed, ok := deref(resultType).(*types.Named)
	if !ok || !strings.HasSuffix(resultNamed.Obj().Name(), "Input") || !internalutil.IsServicePath(resultNamed.Obj().Pkg().Path()) {
		pass.Reportf(method.Pos(), "ValidateAndSanitize on %q must return a type from internal/service whose name ends with Input.", ts.Name.Name)
	}

	errType := sig.Results().At(1).Type()
	if !types.Identical(errType, types.Universe.Lookup("error").Type()) {
		pass.Reportf(method.Pos(), "ValidateAndSanitize on %q must return error as its second result.", ts.Name.Name)
	}
}

func lookupMethod(named *types.Named, name string) *types.Func {
	for i := 0; i < named.NumMethods(); i++ {
		method := named.Method(i)
		if method.Name() == name {
			return method
		}
	}
	methodSet := types.NewMethodSet(types.NewPointer(named))
	for i := 0; i < methodSet.Len(); i++ {
		selection := methodSet.At(i)
		if selection.Obj().Name() == name {
			if fn, ok := selection.Obj().(*types.Func); ok {
				return fn
			}
		}
	}
	return nil
}

func deref(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}
