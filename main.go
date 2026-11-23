package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing/fstest"
)

type fingerprint [sha256.Size]byte
type inventory map[string]fingerprint

func scan(folder fs.FS) (inventory, error) {
	files := inventory{}
	err := fs.WalkDir(folder, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		f, err := folder.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		hash := sha256.New()
		if _, err = io.Copy(hash, f); err != nil {
			return err
		}
		var key fingerprint
		copy(key[:], hash.Sum(nil))
		files[path] = key
		return nil
	})
	return files, err
}

func report(out io.Writer, source, backup inventory) {
	byContent := map[fingerprint]string{}
	for path, key := range backup {
		previous, ok := byContent[key]
		if !ok || path < previous {
			byContent[key] = path
		}
	}
	paths := make([]string, 0, len(source))
	for path := range source {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	covered := 0
	for _, path := range paths {
		key := source[path]
		saved, exists := backup[path]
		other, found := byContent[key]
		switch {
		case exists && saved == key:
			fmt.Fprintf(out, "IDENTIQUE  %q\n", path)
			covered++
		case found:
			fmt.Fprintf(out, "AILLEURS   %q -> %q\n", path, other)
			covered++
		case exists:
			fmt.Fprintf(out, "DIFFÉRENT  %q\n", path)
		default:
			fmt.Fprintf(out, "MANQUANT   %q\n", path)
		}
	}
	fmt.Fprintf(out, "Contenu retrouvé : %d/%d fichiers source\n", covered, len(paths))
}

func root(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%q doit être un dossier", path)
	}
	return resolved, nil
}

func contains(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func run(out io.Writer, args []string) error {
	var source, backup fs.FS
	switch {
	case len(args) == 1 && args[0] == "--demo":
		source, backup = demo()
	case len(args) == 2:
		a, err := root(args[0])
		if err != nil {
			return err
		}
		b, err := root(args[1])
		if err != nil {
			return err
		}
		if contains(a, b) || contains(b, a) {
			return fmt.Errorf("choisissez deux dossiers distincts et non imbriqués")
		}
		source, backup = os.DirFS(a), os.DirFS(b)
	default:
		return fmt.Errorf("usage : noskhati --demo | noskhati SOURCE SAUVEGARDE")
	}
	original, err := scan(source)
	if err != nil {
		return err
	}
	copie, err := scan(backup)
	if err != nil {
		return err
	}
	report(out, original, copie)
	return nil
}

func demo() (fs.FS, fs.FS) {
	original := fstest.MapFS{
		"cours.txt":    {Data: []byte("Algorithmes")},
		"notes.txt":    {Data: []byte("Version B")},
		"exercice.txt": {Data: []byte("Exercice 3")},
		"plan.txt":     {Data: []byte("Plan du TP")},
	}
	backup := fstest.MapFS{
		"archives/chapitre.txt": {Data: []byte("Algorithmes")},
		"notes.txt":             {Data: []byte("Version A")},
		"plan.txt":              {Data: []byte("Plan du TP")},
	}
	return original, backup
}

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
