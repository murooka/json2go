package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sort"
	"strings"

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

	v, err := loadJSON(args)
	if err != nil {
		log.Fatalf("failed to parse JSON: %s", err)
	}

	if opts.Fields != "" {
		fields := strings.Split(opts.Fields, ",")
		sort.Strings(fields)
		filterFields(v, fields)
	}

	typ, err := detectTypeOfItem(v)
	if err != nil {
		log.Fatalf("failed to detect JSON type: %s", err)
	}

	g := NewGenerator()

	src, err := g.Generate(strings.Join(os.Args, " "), opts.Package, opts.TypeName, opts.VarName, typ, v)
	if err != nil {
		log.Fatalf("failed to format output source code: %s", err)
	}

	var out io.Writer
	if opts.Output == "-" {
		out = os.Stdout
	} else {
		f, err := os.Create(opts.Output)
		if err != nil {
			log.Fatalf("failed to open file: %s", err)
		}
		defer f.Close()

		out = f
	}

	out.Write(src)
}

func loadJSON(filenames []string) (interface{}, error) {
	items := make([]interface{}, 0)

	for _, filename := range filenames {
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

		// TODO: consider map
		a := v.([]interface{})
		for _, e := range a {
			items = append(items, e)
		}
	}

	var v interface{}
	v = items

	return v, nil
}

func filterFields(v interface{}, fields []string) {
	fieldMap := map[string]bool{}
	for _, field := range fields {
		fieldMap[field] = true
	}

	a := v.([]interface{})
	for _, obj := range a {
		obj := obj.(map[string]interface{})
		for key := range obj {
			if !fieldMap[key] {
				delete(obj, key)
			}
		}
	}
}
