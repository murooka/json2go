package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/xeipuuv/gojsonpointer"
)

func main() {
	var opts struct {
		Output    string `long:"output"    description:"output filepath" default:"-"`
		TypeName  string `long:"typename"  description:"type name to be generated"`
		VarName   string `long:"varname"   description:"variable name to be generated"`
		Package   string `long:"package"   description:"output package name"`
		Root      string `long:"root"      description:"JSON pointer specifing root object or array"`
		Structure string `long:"structure" description:"comma separated \"map\" or \"slice\" sequence" default:"slice"`
		Fields    string `long:"fields"    description:"comma separated target property names"`
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
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

	structurePaths := strings.Split(opts.Structure, ",")
	for _, path := range structurePaths {
		if path != "slice" && path != "map" {
			log.Fatalf(`--structure must be constructed either "slice" or "map"`)
		}
	}

	v, err := loadJSON(args, opts.Root, structurePaths)
	if err != nil {
		log.Fatalf("failed to parse JSON: %s", err)
	}

	if opts.Fields != "" {
		fields := strings.Split(opts.Fields, ",")
		sort.Strings(fields)
		filterFields(v, fields)
	}

	typ, err := detectTypeInStructure(v, structurePaths)
	if err != nil {
		log.Fatalf("failed to detect JSON type: %s", err)
	}

	src, err := Generate(strings.Join(os.Args, " "), opts.Package, opts.TypeName, opts.VarName, typ, structurePaths, v)
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

func loadJSON(filenames []string, root string, structurePaths []string) (interface{}, error) {
	vs := make([]interface{}, 0, len(filenames))

	pointer, err := gojsonpointer.NewJsonPointer(root)
	if err != nil {
		return nil, fmt.Errorf("invalid json pointer: %s", err)
	}

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

		v, _, err = pointer.Get(v)
		if err != nil {
			log.Fatalf("failed to get root value: %s", err)
		}

		vs = append(vs, v)
	}

	if err := validateStructures(vs, structurePaths); err != nil {
		return nil, err
	}

	return mergeJSONs(vs, structurePaths), nil
}

func validateStructures(vs []interface{}, structurePaths []string) error {
	for _, v := range vs {
		if err := validateStructure(v, structurePaths); err != nil {
			return err
		}
	}

	return nil
}

func validateStructure(v interface{}, structurePaths []string) error {
	if len(structurePaths) == 0 {
		return nil
	}

	// TODO: explicit error message
	path := structurePaths[0]
	switch path {
	case "slice":
		a, ok := v.([]interface{})
		if !ok {
			fmt.Errorf("structure path is slice, but got non-array value")
		}

		for _, e := range a {
			if err := validateStructure(e, structurePaths[1:]); err != nil {
				return err
			}
		}

		break
	case "map":
		m, ok := v.(map[string]interface{})
		if !ok {
			fmt.Errorf("structure path is map, but got non-object value")
		}

		for _, e := range m {
			if err := validateStructure(e, structurePaths[1:]); err != nil {
				return err
			}
		}

		break
	default:
		panic("assertion error")
	}

	return nil
}

func mergeJSONs(vs []interface{}, structurePaths []string) interface{} {
	if len(structurePaths) == 0 {
		return vs[len(vs)-1]
	}

	switch structurePaths[0] {
	case "slice":
		ret := make([]interface{}, 0)
		for _, v := range vs {
			a := v.([]interface{})

			for _, e := range a {
				ret = append(ret, e)
			}
		}

		return ret
	case "map":
		keyMap := map[string]struct{}{}
		for _, v := range vs {
			m := v.(map[string]interface{})
			for k := range m {
				keyMap[k] = struct{}{}
			}
		}

		ret := make(map[string]interface{})
		for key := range keyMap {
			cs := make([]interface{}, 0, len(vs))
			for _, v := range vs {
				m := v.(map[string]interface{})
				e, ok := m[key]
				if ok {
					cs = append(cs, e)
				}
			}
			ret[key] = mergeJSONs(cs, structurePaths[1:])
		}

		return ret
	}

	panic(fmt.Sprintf("assertion error: unexpected structure type: %s", structurePaths[0]))
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
