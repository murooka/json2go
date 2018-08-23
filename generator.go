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

func (g *Generator) Generate(cmd string, pkg string, typeName, varName string, typ *JSONType, v interface{}) ([]byte, error) {
	g.Printlnf(`// Code generated by "%s"; DO NOT EDIT.`, cmd)
	g.Printlnf("")
	g.Printlnf("package %s", pkg)
	g.Printlnf("")
	g.Printlnf("type %s %s", typeName, typ.ToGoType())
	g.Printlnf("")

	g.Printlnf("var %s = []%s {", varName, typeName)

	a := v.([]interface{})
	for _, e := range a {
		g.Printlnf(g.toLiteral(e, typ) + ",")
	}
	g.Printlnf("}")

	return format.Source(g.buf.Bytes())
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
