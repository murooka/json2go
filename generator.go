package main

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strconv"
	"text/template"

	"github.com/iancoleman/strcase"
)

var outputTmpl = `
// Code generated by "{{.Command}}"; DO NOT EDIT.

package {{.Package}}

type {{.TypeName}} {{.TypeDef}}

var {{.VarName}} = {{.VarDef}}
`

func Generate(cmd string, pkg string, typeName, varName string, typ *JSONType, structurePaths []string, v interface{}) ([]byte, error) {
	tmpl := template.Must(template.New("output").Parse(outputTmpl))
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, map[string]interface{}{
		"Command":  cmd,
		"Package":  pkg,
		"TypeName": typeName,
		"TypeDef":  typ.ToGoType(),
		"VarName":  varName,
		"VarDef":   makeVarDef(typeName, typ, v, structurePaths),
	})
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

func makeVarDef(typeName string, typ *JSONType, v interface{}, structurePaths []string) string {
	buf := NewExtBuffer()
	makeVarBody(buf, typ, v, structurePaths)
	return makeVarType(typeName, structurePaths) + buf.String()
}

func makeVarType(typeName string, structurePaths []string) string {
	if len(structurePaths) == 0 {
		return typeName
	}

	switch structurePaths[0] {
	case "slice":
		return fmt.Sprintf("[]%s", makeVarType(typeName, structurePaths[1:]))
	case "map":
		return fmt.Sprintf("map[string]%s", makeVarType(typeName, structurePaths[1:]))
	}

	panic("assertion error")
}

func makeVarBody(b *ExtBuffer, typ *JSONType, v interface{}, structurePaths []string) {
	if len(structurePaths) == 0 {
		b.Printf(toLiteral(v, typ))
		return
	}

	switch structurePaths[0] {
	case "slice":
		a := v.([]interface{})

		b.Println("{")
		for _, e := range a {
			makeVarBody(b, typ, e, structurePaths[1:])
			b.Println(",")
		}
		b.Print("}")
	case "map":
		m := v.(map[string]interface{})

		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)

		b.Println("{")
		for _, k := range ks {
			e := m[k]
			b.Printf(`%s: `, strconv.Quote(k))
			makeVarBody(b, typ, e, structurePaths[1:])
			b.Println(",")
		}
		b.Print("}")
	}

}

func toLiteral(v interface{}, typ *JSONType) string {
	switch v := v.(type) {
	case map[string]interface{}:
		if typ.Object == nil {
			panic("assertion error")
		}

		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		buf := &bytes.Buffer{}
		buf.WriteString("{\n")
		for _, k := range keys {
			e := v[k]
			ctyp := typ.Object[k]
			buf.WriteString(fmt.Sprintf("%s: %s,\n", strcase.ToCamel(
				k), toLiteral(e, ctyp)))
		}
		buf.WriteString("}")
		return buf.String()
	case []interface{}:
		if typ.Array == nil {
			panic("assertion error")
		}

		buf := &bytes.Buffer{}
		buf.WriteString(fmt.Sprintf("[]%s{\n", typ.Array.ToGoType()))
		for _, e := range v {
			buf.WriteString(toLiteral(e, typ.Array))
			buf.WriteString(",\n")
		}
		buf.WriteString("}")
		return buf.String()
	case string:
		return strconv.Quote(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "nil"
	}

	panic(fmt.Sprintf("unknown type of value: %#v", v))
}
