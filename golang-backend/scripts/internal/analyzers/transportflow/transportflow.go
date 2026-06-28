package transportflow

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/xiabing/skills/golang-backend/scripts/internal/analyzers/internalutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "golangbackend_transportflow",
	Doc:  "prevents raw request types from crossing the transport-service boundary",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	if !internalutil.IsTransportPath(pass.Pkg.Path()) {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			calleeType := pass.TypesInfo.TypeOf(call.Fun)
			if calleeType == nil {
				return true
			}
			sig, ok := calleeType.(*types.Signature)
			if !ok {
				return true
			}
			if !isServiceLikeCall(sig, call.Fun) {
				return true
			}
			for _, arg := range call.Args {
				argType := pass.TypesInfo.TypeOf(arg)
				if isRequestType(argType) {
					pass.Reportf(arg.Pos(), "service call receives request %q directly. Call ValidateAndSanitize() first and pass the resulting service.*Input value instead.", requestTypeName(argType))
				}
			}
			return true
		})
	}

	return nil, nil
}

func isServiceLikeCall(sig *types.Signature, fun ast.Expr) bool {
	if _, ok := fun.(*ast.SelectorExpr); ok {
		if sig.Recv() != nil {
			if recvType, ok := deref(sig.Recv().Type()).(*types.Named); ok && recvType.Obj() != nil && recvType.Obj().Pkg() != nil && internalutil.IsServicePath(recvType.Obj().Pkg().Path()) {
				return true
			}
		}
	}
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		if named, ok := deref(params.At(i).Type()).(*types.Named); ok && named.Obj() != nil && named.Obj().Pkg() != nil && internalutil.IsServicePath(named.Obj().Pkg().Path()) {
			return true
		}
	}
	return false
}

func isRequestType(t types.Type) bool {
	named, ok := deref(t).(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}
	return strings.HasSuffix(named.Obj().Name(), "Request") && internalutil.IsTransportPath(named.Obj().Pkg().Path())
}

func requestTypeName(t types.Type) string {
	named, ok := deref(t).(*types.Named)
	if !ok || named.Obj() == nil {
		return types.TypeString(t, nil)
	}
	return named.Obj().Name()
}

func deref(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return ptr.Elem()
	}
	return t
}
