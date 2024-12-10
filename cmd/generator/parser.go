package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Parser struct {
	fset     *token.FileSet
	pkgPath  string
	packages map[string]*ast.Package
}

func NewParser(pkgPath string) (*Parser, error) {
	fset := token.NewFileSet()

	packages, err := parser.ParseDir(fset, pkgPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse directory: %w", err)
	}

	return &Parser{
		fset:     fset,
		pkgPath:  pkgPath,
		packages: packages,
	}, nil
}

func (p *Parser) CollectTypes(interfaces []Interface) ([]CustomType, []string, error) {
	types := make(map[string]CustomType)
	imports := make(map[string]bool)

	for _, iface := range interfaces {
		for _, method := range iface.Methods {
			// Parameter tiplerini kontrol et
			for _, param := range method.Params {
				p.collectType(param.Type, types, imports)
			}
			// Return tiplerini kontrol et
			for _, ret := range method.Returns {
				p.collectType(ret.Type, types, imports)
			}
		}
	}

	// Map to slice
	var customTypes []CustomType
	var importPaths []string

	for _, t := range types {
		customTypes = append(customTypes, t)
	}
	for imp := range imports {
		importPaths = append(importPaths, imp)
	}

	return customTypes, importPaths, nil
}

func (p *Parser) ParseInterface(names ...string) ([]Interface, error) {
	var interfaces []Interface

	for _, name := range names {
		iface, err := p.FindInterface(name)
		if err != nil {
			return nil, fmt.Errorf("finding interface %s: %w", name, err)
		}

		methods, err := p.GetInterfaceMethods(iface)
		if err != nil {
			return nil, fmt.Errorf("getting methods for interface %s: %w", name, err)
		}

		interfaces = append(interfaces, Interface{
			Name:    name,
			Methods: methods,
			Doc:     p.getInterfaceDoc(iface),
		})
	}

	return interfaces, nil
}

func (p *Parser) FindInterface(name string) (*ast.InterfaceType, error) {
	for _, pkg := range p.packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || typeSpec.Name.Name != name {
						continue
					}

					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						return interfaceType, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("interface %s not found", name)
}

func (p *Parser) GetInterfaceMethods(iface *ast.InterfaceType) ([]Method, error) {
	var methods []Method

	for _, m := range iface.Methods.List {
		funcType, ok := m.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		method := Method{
			Name:    m.Names[0].Name,
			Params:  p.parseFieldList(funcType.Params),
			Returns: p.parseFieldList(funcType.Results),
		}

		if m.Doc != nil {
			method.Doc = m.Doc.Text()
		}

		methods = append(methods, method)
	}

	return methods, nil
}

func (p *Parser) parseFieldList(fields *ast.FieldList) []Parameter {
	var params []Parameter
	if fields == nil {
		return params
	}

	for _, field := range fields.List {
		typ := p.typeToString(field.Type)

		if len(field.Names) == 0 {
			params = append(params, Parameter{
				Type:       typ,
				IsVariadic: p.isVariadic(field.Type),
			})
			continue
		}

		for _, name := range field.Names {
			params = append(params, Parameter{
				Name:       name.Name,
				Type:       typ,
				IsVariadic: p.isVariadic(field.Type),
				Doc:        field.Doc.Text(),
			})
		}
	}

	return params
}

func (p *Parser) collectType(typeName string, types map[string]CustomType, imports map[string]bool) {
	// bypass builtin types
	if isBuiltinType(typeName) {
		return
	}

	// cleanup pointer character
	typeName = strings.TrimPrefix(typeName, "*")

	// find type definition
	for _, pkg := range p.packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if typeSpec.Name.Name == typeName {
								// Tip tanımını kaydet
								types[typeName] = CustomType{
									Name:       typeName,
									Definition: p.exprToString(typeSpec.Type),
								}

								// İlgili importları topla
								if sel, ok := typeSpec.Type.(*ast.SelectorExpr); ok {
									if x, ok := sel.X.(*ast.Ident); ok {
										if imp := p.findImportPath(file, x.Name); imp != "" {
											imports[imp] = true
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (p *Parser) findImportPath(file *ast.File, pkgName string) string {
	for _, imp := range file.Imports {
		if imp.Name != nil && imp.Name.Name == pkgName {
			return strings.Trim(imp.Path.Value, "\"")
		}
		parts := strings.Split(strings.Trim(imp.Path.Value, "\""), "/")
		if parts[len(parts)-1] == pkgName {
			return strings.Trim(imp.Path.Value, "\"")
		}
	}
	return ""
}

func isBuiltinType(typeName string) bool {
	builtins := map[string]bool{
		"bool":       true,
		"byte":       true,
		"complex64":  true,
		"complex128": true,
		"error":      true,
		"float32":    true,
		"float64":    true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"rune":       true,
		"string":     true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"uintptr":    true,
	}
	return builtins[strings.TrimPrefix(typeName, "*")]
}

func (p *Parser) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", p.typeToString(t.X), t.Sel.Name)
	case *ast.StarExpr:
		return "*" + p.typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + p.typeToString(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", p.typeToString(t.Key), p.typeToString(t.Value))
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return "..." + p.typeToString(t.Elt)
	case *ast.FuncType:
		return "func" + p.funcTypeToString(t)
	case *ast.ChanType:
		switch t.Dir {
		case ast.SEND:
			return "chan<- " + p.typeToString(t.Value)
		case ast.RECV:
			return "<-chan " + p.typeToString(t.Value)
		default:
			return "chan " + p.typeToString(t.Value)
		}
	default:
		return fmt.Sprintf("unsupported type: %T", expr)
	}
}

func (p *Parser) funcTypeToString(ft *ast.FuncType) string {
	var result string
	result += "("

	if ft.Params != nil {
		params := p.parseFieldList(ft.Params)
		for i, param := range params {
			if i > 0 {
				result += ", "
			}
			result += param.Type
		}
	}

	result += ")"

	if ft.Results != nil {
		returns := p.parseFieldList(ft.Results)
		if len(returns) > 0 {
			result += " ("
			for i, ret := range returns {
				if i > 0 {
					result += ", "
				}
				result += ret.Type
			}
			result += ")"
		}
	}

	return result
}

func (p *Parser) exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + p.exprToString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + p.exprToString(t.Elt)
		}
		return fmt.Sprintf("[%s]%s", p.exprToString(t.Len), p.exprToString(t.Elt))
	case *ast.SelectorExpr:
		return p.exprToString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", p.exprToString(t.Key), p.exprToString(t.Value))
	case *ast.StructType:
		var fields []string
		for _, field := range t.Fields.List {
			var names []string
			for _, name := range field.Names {
				names = append(names, name.Name)
			}
			fieldStr := ""
			if len(names) > 0 {
				fieldStr = strings.Join(names, ", ") + " "
			}
			fieldStr += p.exprToString(field.Type)
			if field.Tag != nil {
				fieldStr += " " + field.Tag.Value
			}
			fields = append(fields, fieldStr)
		}
		return fmt.Sprintf("struct {\n\t%s\n}", strings.Join(fields, "\n\t"))
	case *ast.InterfaceType:
		var methods []string
		for _, method := range t.Methods.List {
			if len(method.Names) > 0 {
				methodType := p.exprToString(method.Type)
				methods = append(methods, fmt.Sprintf("%s %s", method.Names[0], methodType))
			}
		}
		return fmt.Sprintf("interface {\n\t%s\n}", strings.Join(methods, "\n\t"))
	case *ast.FuncType:
		var params, results []string

		if t.Params != nil {
			for _, param := range t.Params.List {
				paramType := p.exprToString(param.Type)
				if len(param.Names) > 0 {
					for _, name := range param.Names {
						params = append(params, fmt.Sprintf("%s %s", name.Name, paramType))
					}
				} else {
					params = append(params, paramType)
				}
			}
		}

		if t.Results != nil {
			for _, result := range t.Results.List {
				resultType := p.exprToString(result.Type)
				if len(result.Names) > 0 {
					for _, name := range result.Names {
						results = append(results, fmt.Sprintf("%s %s", name.Name, resultType))
					}
				} else {
					results = append(results, resultType)
				}
			}
		}

		paramsStr := strings.Join(params, ", ")
		if len(results) == 0 {
			return fmt.Sprintf("func(%s)", paramsStr)
		}
		resultsStr := strings.Join(results, ", ")
		if len(results) == 1 && !strings.Contains(results[0], " ") {
			return fmt.Sprintf("func(%s) %s", paramsStr, resultsStr)
		}
		return fmt.Sprintf("func(%s) (%s)", paramsStr, resultsStr)
	default:
		return fmt.Sprintf("/* unsupported type: %T */", expr)
	}
}

func (p *Parser) isVariadic(expr ast.Expr) bool {
	_, ok := expr.(*ast.Ellipsis)
	return ok
}

func (p *Parser) getInterfaceDoc(iface *ast.InterfaceType) string {
	// InterfaceType'dan doğrudan doc comment'e erişemeyiz
	// Bu yüzden parent node üzerinden bulmamız gerekiyor
	for _, pkg := range p.packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						if interfaceType == iface {
							// İlgili interface'i bulduk
							if typeSpec.Doc != nil {
								return typeSpec.Doc.Text()
							}
							if genDecl.Doc != nil {
								return genDecl.Doc.Text()
							}
							return ""
						}
					}
				}
			}
		}
	}
	return ""
}
