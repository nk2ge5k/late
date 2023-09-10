package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	var (
		inDir       string
		outFilename string
	)

	flag.StringVar(&inDir, "migrations-dir", "",
		"Path to the directory that contains database migrations")
	flag.StringVar(&outFilename, "out", "", "Path to the output file")
	flag.Parse()

	if outFilename == "" {
		flag.Usage()
		exitf("missing output file name")
	}

	if !isDirectory(inDir) {
		flag.Usage()
		exitf("path %q does not exists or not a directory", inDir)
	}

	filenames, err := sqlPaths(inDir)
	if err != nil {
		exitf("failed to collect sql files: %v", err)
	}

	sort.Strings(filenames)

	f, err := os.OpenFile(outFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		exitf("open file %q: %v", outFilename, err)
	}
	defer f.Close()

	for _, filename := range filenames {
		fmt.Fprintln(f, filepath.Base(filename))
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

func sqlPaths(dir string) ([]string, error) {
	var filenames []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() && path != dir {
			return fs.SkipDir
		}

		if filepath.Ext(path) != ".sql" {
			return nil
		}

		filenames = append(filenames, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk dir: %w", err)
	}

	return filenames, nil
}

func isDirectory(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}
