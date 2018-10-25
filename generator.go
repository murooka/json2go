package main

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strconv"

	"github.com/iancoleman/strcase"
)

type Generator struct {
	buf *bytes.Buffer
}

func NewGenerator() *Generator {
	return &Generator{
		buf: &bytes.Buffer{},
	}
}

func (g *Generator) Generate(cmd string, pkg string, typeName, varName string, typ *JSONType, structurePaths []string, v interface{}) ([]byte, error) {
	g.Printlnf(`// Code generated by "%s"; DO NOT EDIT.`, cmd)
	g.Printlnf("")
	g.Printlnf("package %s", pkg)
	g.Printlnf("")
	g.Printlnf("type %s %s", typeName, typ.ToGoType())
	g.Printlnf("")
	g.Printf("var %s = %s", varName, makeVarType(typeName, structurePaths))

	g.printStructure(structurePaths, typ, v)

	return format.Source(g.buf.Bytes())
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

func (g *Generator) printStructure(structurePaths []string, typ *JSONType, v interface{}) {
	if len(structurePaths) == 0 {
		g.Printf(g.toLiteral(v, typ))
		return
	}

	switch structurePaths[0] {
	case "slice":
		a := v.([]interface{})

		g.Printlnf("{")
		for _, e := range a {
			g.printStructure(structurePaths[1:], typ, e)
			g.Println(",")
		}
		g.Print("}")
	case "map":
		m := v.(map[string]interface{})

		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)

		g.Printlnf("{")
		for _, k := range ks {
			e := m[k]
			g.Printf(`%s: `, strconv.Quote(k))
			g.printStructure(structurePaths[1:], typ, e)
			g.Println(",")
		}
		g.Print("}")
	}

}

func (g *Generator) Print(s string) {
	g.buf.WriteString(s)
}

func (g *Generator) Println(s string) {
	g.buf.WriteString(s)
	g.buf.WriteByte('\n')
}

func (g *Generator) Printf(format string, args ...interface{}) {
	g.buf.WriteString(fmt.Sprintf(format, args...))
}

func (g *Generator) Printlnf(format string, args ...interface{}) {
	g.buf.WriteString(fmt.Sprintf(format, args...))
	g.buf.WriteByte('\n')
}

func (g *Generator) toLiteral(v interface{}, typ *JSONType) string {
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
				k), g.toLiteral(e, ctyp)))
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
			buf.WriteString(g.toLiteral(e, typ.Array))
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
