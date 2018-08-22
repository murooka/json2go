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

	g := NewGenerator(strings.Join(os.Args, " "), opts.Package, opts.TypeName, opts.VarName, typ, v, fields)

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
