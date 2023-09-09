package main

import (
	"bufio"
	"bytes"
	"embed"
	"flag"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

var (
	//go:embed templates/*
	templatesFS embed.FS
	tmpl        = template.Must(
		template.ParseFS(templatesFS, "templates/package.go.tmpl"))
	formatter = &printer.Config{
		Mode:     printer.TabIndent | printer.UseSpaces,
		Tabwidth: 8,
	}
)

type Query struct {
	// Name of the query
	Name string
	// GoName is a name used for variable
	GoName string
	// Filename is an orignal name of the file which contained query
	Filename string
	// SQL contains SQL query text
	SQL string
}

func main() {
	var export bool
	flag.BoolVar(&export, "export", false, "Generate queries as exproted variables")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		exitf("missing path argument")
	}

	dir := args[0]
	if !isDirectory(dir) {
		flag.Usage()
		exitf("file %q is not a directory or does not exist", dir)
	}

	dirQueries := make(map[string][]*Query)

	werr := filepath.WalkDir(dir, func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() && path != dir {
			return fs.SkipDir
		}

		if filepath.Ext(path) != ".sql" {
			return nil
		}

		q, err := readQuery(path, export)
		if err != nil {
			warnf("Failed to read query from file %q: %v", path, err)
		}

		dir := filepath.Dir(path)
		dirQueries[dir] = append(dirQueries[dir], q)

		return nil
	})
	if werr != nil {
		exitf("failed to scan directory for sql files: %v", werr)
	}

	for dir, queries := range dirQueries {
		if err := writeQueries(dir, queries); err != nil {
			exitf("failed to write queries: %v", err)
		}
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

func warnf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func isDirectory(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func readQuery(path string, export bool) (*Query, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	filename := filepath.Base(path)
	ext := filepath.Ext(filename)

	return &Query{
		Name:     filename[:len(filename)-len(ext)],
		GoName:   sanitize(filename, export),
		Filename: path,
		SQL:      string(b),
	}, nil
}

func sanitize(s string, export bool) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)

	if export {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return s
}

func writeQueries(dir string, queries []*Query) error {
	if len(queries) == 0 {
		return nil
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("absolute path: %w", err)
	}

	pkg := filepath.Base(abs)

	buf := bytes.NewBuffer(nil)
	terr := tmpl.ExecuteTemplate(buf, "package.go.tmpl", struct {
		PackageName string
		Queries     []*Query
	}{
		PackageName: pkg,
		Queries:     queries,
	})

	if terr != nil {
		return fmt.Errorf("template execute: %w", terr)
	}

	b := buf.Bytes()

	fset := token.NewFileSet()
	file, perr := parser.ParseFile(fset, "", b, parser.ParseComments)
	if perr != nil {
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(b))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		return fmt.Errorf("unparsable Go source: %v\n%v", perr, src.String())
	}

	out, oferr := os.OpenFile(filepath.Join(dir, "queries.gen.go"),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if oferr != nil {
		return fmt.Errorf("open file: %w", oferr)
	}
	defer out.Close()

	if ferr := formatter.Fprint(out, fset, file); ferr != nil {
		return fmt.Errorf("can not reformat Go source: %v", ferr)
	}

	return nil
}
