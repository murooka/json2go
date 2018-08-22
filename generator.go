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
	buf      *bytes.Buffer
	cmd      string
	pkg      string
	typeName string
	varName  string
	typ      *JSONType
	v        interface{}
	fields   []string
}

func NewGenerator(cmd string, pkg string, typeName, varName string, typ *JSONType, v interface{}, fields []string) *Generator {
	return &Generator{
		buf:      &bytes.Buffer{},
		cmd:      cmd,
		pkg:      pkg,
		typeName: typeName,
		varName:  varName,
		typ:      typ,
		v:        v,
		fields:   fields,
	}
}

func (g *Generator) Generate() ([]byte, error) {
	g.Printlnf(`// Code generated by "%s"; DO NOT EDIT.`, g.cmd)
	g.Printlnf("")
	g.Printlnf("package %s", g.pkg)
	g.Printlnf("")
	g.Printlnf("type %s %s", g.typeName, g.typ.ToGoType())
	g.Printlnf("")

	g.Printlnf("var %s = []%s {", g.varName, g.typeName)

	a := g.v.([]interface{})
	for _, e := range a {
		g.Printlnf(g.toLiteral(e, g.typ) + ",")
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

		var keys []string
		if g.fields == nil {
			keys = make([]string, 0, len(v))
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)
		} else {
			keys = g.fields
		}

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
		return "null"
	}

	panic(fmt.Sprintf("unknown type of value: %#v", v))
}