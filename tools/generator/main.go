package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

// Param represents a function parameter name and type
type Param struct {
	Name string
	Type string
}

// MethodSpec represents a method specification to be generated
type MethodSpec struct {
	Name        string
	Params      []Param
	RetType     string
	RawRet      string
	IsSlice     bool
	SuccessType string
	Doc         []string
}

func main() {
	targetFile := "../client/oas_response_decoders_gen.go"
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		targetFile = "pkg/client/oas_response_decoders_gen.go"
	}

	fset := token.NewFileSet()

	// 1. Parse the response decoders file into an Abstract Syntax Tree (AST)
	node, err := parser.ParseFile(fset, targetFile, nil, 0)
	if err != nil {
		panic(err)
	}

	mappings := make(map[string]string)

	// 2. Walk the AST to find decoder functions
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if !strings.HasPrefix(fn.Name.Name, "decode") || !strings.HasSuffix(fn.Name.Name, "Response") {
			return true
		}

		if fn.Type.Results == nil || len(fn.Type.Results.List) == 0 {
			return true
		}
		retTypeExpr := fn.Type.Results.List[0].Type
		ident, ok := retTypeExpr.(*ast.Ident)
		if !ok {
			return true
		}
		interfaceName := ident.Name

		successType := findSuccessStructInDecoder(fn.Body)
		if successType != "" {
			mappings[interfaceName] = successType
		}

		return true
	})

	// 2.5 Parse oas_schemas_gen.go to find slice types and map them to their underlying slice types
	schemasFile := "../client/oas_schemas_gen.go"
	if _, err := os.Stat(schemasFile); os.IsNotExist(err) {
		schemasFile = "pkg/client/oas_schemas_gen.go"
	}

	schemasNode, err := parser.ParseFile(fset, schemasFile, nil, 0)
	if err != nil {
		panic(err)
	}

	sliceMappings := make(map[string]string)
	ast.Inspect(schemasNode, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		arrayType, ok := ts.Type.(*ast.ArrayType)
		if !ok {
			return true
		}
		if arrayType.Len == nil {
			elemType := renderType(arrayType.Elt)
			sliceMappings[ts.Name.Name] = "[]" + elemType
		}
		return true
	})

	// 3. Parse oas_client_gen.go to extract the Invoker interface methods
	clientFile := "../client/oas_client_gen.go"
	if _, err := os.Stat(clientFile); os.IsNotExist(err) {
		clientFile = "pkg/client/oas_client_gen.go"
	}

	clientNode, err := parser.ParseFile(fset, clientFile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var methods []MethodSpec

	ast.Inspect(clientNode, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if ts.Name.Name != "Invoker" {
			return true
		}
		itype, ok := ts.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		for _, field := range itype.Methods.List {
			if len(field.Names) == 0 {
				continue
			}
			methodName := field.Names[0].Name
			ft, ok := field.Type.(*ast.FuncType)
			if !ok {
				continue
			}

			var params []Param
			for i, pField := range ft.Params.List {
				typeStr := renderType(pField.Type)
				if len(pField.Names) == 0 {
					params = append(params, Param{
						Name: fmt.Sprintf("arg%d", i),
						Type: typeStr,
					})
				} else {
					for _, pName := range pField.Names {
						params = append(params, Param{
							Name: pName.Name,
							Type: typeStr,
						})
					}
				}
			}

			var retType string
			var rawRet string
			var isSlice bool
			var successType string

			if ft.Results != nil && len(ft.Results.List) > 0 {
				rawRet = renderType(ft.Results.List[0].Type)
				lookupName := strings.TrimPrefix(rawRet, "client.")
				if structName, ok := mappings[lookupName]; ok {
					if isBuiltin(structName) {
						retType = structName
						successType = structName
					} else {
						successType = "*client." + structName
						if sliceType, found := sliceMappings[structName]; found {
							isSlice = true
							retType = sliceType
						} else {
							retType = "*client." + structName
						}
					}
				} else {
					retType = rawRet
					successType = rawRet
				}
			} else {
				retType = "error"
				successType = "error"
			}

			var doc []string
			if field.Doc != nil {
				for _, comment := range field.Doc.List {
					doc = append(doc, comment.Text)
				}
			}

			methods = append(methods, MethodSpec{
				Name:        methodName,
				Params:      params,
				RetType:     retType,
				RawRet:      rawRet,
				IsSlice:     isSlice,
				SuccessType: successType,
				Doc:         doc,
			})
		}
		return false
	})

	// 4. Determine the correct output path for method_gen.go
	outputPath := "pkg/api/method_gen.go"
	if _, err := os.Stat("pkg/api"); os.IsNotExist(err) {
		outputPath = "method_gen.go"
	}

	// 5. Generate code using a template
	tmpl, err := template.New("methods").Funcs(template.FuncMap{
		"hasEllipsis": func(s string) bool {
			return strings.HasPrefix(s, "...")
		},
	}).Parse(tmplSource)
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, struct {
		Methods []MethodSpec
	}{
		Methods: methods,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully generated %d strongly-typed wrapper methods in %s\n", len(methods), outputPath)
}

func isBuiltin(s string) bool {
	switch s {
	case "bool", "string", "int", "int32", "int64", "float32", "float64", "error", "any":
		return true
	default:
		return false
	}
}

func renderType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		if isBuiltin(t.Name) {
			return t.Name
		}
		return "client." + t.Name
	case *ast.StarExpr:
		return "*" + renderType(t.X)
	case *ast.SelectorExpr:
		pkgIdent, ok := t.X.(*ast.Ident)
		if ok {
			return pkgIdent.Name + "." + t.Sel.Name
		}
		return renderType(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + renderType(t.Elt)
		}
		return fmt.Sprintf("[%s]%s", renderType(t.Len), renderType(t.Elt))
	case *ast.Ellipsis:
		return "..." + renderType(t.Elt)
	case *ast.BasicLit:
		return t.Value
	default:
		return ""
	}
}

// findSuccessStructInDecoder walks inside a decoder function to look for
// "case 200" or "case 201" and find the local variable instantiated there.
func findSuccessStructInDecoder(body *ast.BlockStmt) string {
	var foundStruct string

	ast.Inspect(body, func(n ast.Node) bool {
		switchStmt, ok := n.(*ast.SwitchStmt)
		if !ok {
			return true
		}

		for _, stmt := range switchStmt.Body.List {
			caseClause, ok := stmt.(*ast.CaseClause)
			if !ok {
				continue
			}

			isSuccessStatus := false
			for _, expr := range caseClause.List {
				if basicLit, ok := expr.(*ast.BasicLit); ok && (basicLit.Value == "200" || basicLit.Value == "201") {
					isSuccessStatus = true
					break
				}
			}

			if isSuccessStatus {
				ast.Inspect(caseClause, func(innerNode ast.Node) bool {
					declStmt, ok := innerNode.(*ast.DeclStmt)
					if !ok {
						return true
					}
					genDecl, ok := declStmt.Decl.(*ast.GenDecl)
					if !ok || genDecl.Tok != token.VAR {
						return true
					}
					for _, spec := range genDecl.Specs {
						valueSpec, ok := spec.(*ast.ValueSpec)
						if !ok || len(valueSpec.Names) == 0 {
							continue
						}
						if valueSpec.Names[0].Name == "response" {
							if typeIdent, ok := valueSpec.Type.(*ast.Ident); ok {
								foundStruct = typeIdent.Name
								return false
							}
						}
					}
					return true
				})
			}
		}
		return true
	})

	return foundStruct
}

const tmplSource = `package api

// Code generated by tools/generator; DO NOT EDIT.

import (
	"context"

	"opencode-go-kit/pkg/client"
	"github.com/go-faster/jx"
)
{{range .Methods}}
{{range .Doc}}{{.}}
{{end}}func (a *API) {{.Name}}(
{{range .Params}}	{{.Name}} {{.Type}},
{{end}}) ({{.RetType}}, error) {
{{if .IsSlice}}	res, err := Extract[{{.SuccessType}}](a.Client.{{.Name}}({{range $i, $p := .Params}}{{if $i}}, {{end}}{{$p.Name}}{{if hasEllipsis $p.Type}}...{{end}}{{end}}))
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return {{.RetType}}(*res), nil
{{else}}	return Extract[{{.RetType}}](a.Client.{{.Name}}({{range $i, $p := .Params}}{{if $i}}, {{end}}{{$p.Name}}{{if hasEllipsis $p.Type}}...{{end}}{{end}}))
{{end}}}
{{end}}
`
