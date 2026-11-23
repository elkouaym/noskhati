package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDemo(t *testing.T) {
	var out bytes.Buffer
	if err := run(&out, []string{"--demo"}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"AILLEURS", "DIFFÉRENT", "IDENTIQUE", "MANQUANT", "2/4"} {
		if !strings.Contains(out.String(), want) {
			t.Fatal(out.String())
		}
	}
}
func TestDisk(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	for _, dir := range []string{src, dst} {
		if err := os.WriteFile(filepath.Join(dir, "test.txt"), []byte("exemple"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	var out bytes.Buffer
	if err := run(&out, []string{src, dst}); err != nil || !strings.Contains(out.String(), "1/1") {
		t.Fatal(err, out.String())
	}
}
func TestNested(t *testing.T) {
	parent := t.TempDir()
	child := filepath.Join(parent, "copie")
	if err := os.Mkdir(child, 0700); err != nil {
		t.Fatal(err)
	}
	if run(&bytes.Buffer{}, []string{parent, child}) == nil {
		t.Fatal("dossiers imbriqués acceptés")
	}
}
