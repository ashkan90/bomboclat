package generator

import (
	"fmt"
	"go/ast"
	"go/types"
)

type FunctionAnalyzer struct {
	signatures map[string]*FuncSignature
}
type CallAnalysis struct {
	FuncName string
	Args     []Argument
	Returns  []ReturnType
}

type Argument struct {
	Name    string
	Type    string
	IsConst bool
	Value   interface{}
}

type ReturnType struct {
	Type    string
	IsError bool
}

func (fa *FunctionAnalyzer) AnalyzeCall(call *ast.CallExpr) (*CallAnalysis, error) {
	return &CallAnalysis{
		FuncName: extractFuncName(call),
		Args:     extractArgs(call),
		Returns:  extractReturns(call),
	}, nil
}

func extractFuncName(call *ast.CallExpr) string {
	if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
		return selector.Sel.Name
	}
	if ident, ok := call.Fun.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

func extractArgs(call *ast.CallExpr) []Argument {
	args := make([]Argument, len(call.Args))
	for i, arg := range call.Args {
		args[i] = extractArgument(arg)
	}
	return args
}

func extractArgument(expr ast.Expr) Argument {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return Argument{
			Type:    e.Kind.String(),
			IsConst: true,
			Value:   e.Value,
		}
	case *ast.Ident:
		return Argument{
			Name: e.Name,
			Type: e.Obj.Kind.String(), // acaba ?
		}
	default:
		return Argument{Type: fmt.Sprintf("%T", expr)}
	}
}

func extractReturns(call *ast.CallExpr) []ReturnType {
	if t, ok := call.Fun.(*ast.Ident); ok {
		if t.Obj != nil {
			if fn, ok := t.Obj.Decl.(*ast.FuncDecl); ok {
				return extractFuncReturns(fn.Type.Results)
			}
		}
	}
	return nil
}

func extractFuncReturns(fields *ast.FieldList) []ReturnType {
	if fields == nil {
		return nil
	}

	returns := make([]ReturnType, 0, len(fields.List))
	for _, field := range fields.List {
		returns = append(returns, ReturnType{
			Type:    types.ExprString(field.Type),
			IsError: isErrorType(field.Type),
		})
	}
	return returns
}

func isErrorType(expr ast.Expr) bool {
	if id, ok := expr.(*ast.Ident); ok {
		return id.Name == "error"
	}
	return false
}
