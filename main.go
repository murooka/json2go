package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jessevdk/go-flags"
)

func main() {
	var opts struct {
		Output   string `long:"output" default:"-"`
		TypeName string `long:"typename"`
		VarName  string `long:"varname"`
		Package  string `long:"package"`
		Fields   string `long:"fields"`
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		if err := err.(*flags.Error); err.Type == flags.ErrHelp {
			os.Exit(1)
		} else {
			log.Fatalf("failed to parse arguments: %s", err)
		}
	}

	if len(args) == 0 {
		log.Fatal("argument is not found")
	} else if len(args) > 1 {
		log.Fatalf("too many arguments: %v", args)
	}

	if opts.Package == "" {
		log.Fatal("--package must be specifed")
	}

	if opts.TypeName == "" {
		log.Fatal("--typename must be specified")
	}

	if opts.VarName == "" {
		log.Fatal("--varname must be specified")
	}

	var fields []string
	if opts.Fields != "" {
		fields = strings.Split(opts.Fields, ",")
		sort.Strings(fields)
	}

	filename := args[0]
	v, err := loadJSON(filename)
	if err != nil {
		log.Fatalf("failed to parse JSON: %s", err)
	}

	typ, err := detectTypeOfItem(v, fields)
	if err != nil {
		log.Fatalf("failed to detect JSON type: %s", err)
	}

	g := NewGenerator(fields)
	g.Printlnf(`// Code generated by "%s"; DO NOT EDIT.`, strings.Join(os.Args, " "))
	g.Printlnf("")
	g.Printlnf("package %s", opts.Package)
	g.Printlnf("")
	g.Printlnf("type %s %s", opts.TypeName, typ.ToGoType())
	g.Printlnf("")

	g.Printlnf("var %s = []%s {", opts.VarName, opts.TypeName)
	a := v.([]interface{})
	for _, e := range a {
		g.Printlnf(g.toLiteral(e, typ) + ",")
	}
	g.Printlnf("}")

	src, err := g.Generate()
	if err != nil {
		log.Fatalf("failed to format output source code: %s", err)
	}

	var out io.Writer
	if opts.Output == "-" {
		out = os.Stdout
	} else {
		out, err = os.Create(opts.Output)
		if err != nil {
			log.Fatalf("failed to open file: %s", err)
		}
	}

	out.Write(src)
}

func loadJSON(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer f.Close()

	var v interface{}
	err = json.NewDecoder(f).Decode(&v)
	if err != nil {
		log.Fatalf("failed to decode JSON: %s", err)
	}

	return v, nil
}

type Generator struct {
	buf    *bytes.Buffer
	fields []string
}

func NewGenerator(fields []string) *Generator {
	return &Generator{
		buf:    &bytes.Buffer{},
		fields: fields,
	}
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

func (g *Generator) Printlnf(format string, args ...interface{}) {
	g.buf.WriteString(fmt.Sprintf(format, args...))
	g.buf.WriteByte('\n')
}

func (g *Generator) Generate() ([]byte, error) {
	return format.Source(g.buf.Bytes())
}
