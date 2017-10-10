package main

import (
	"path/filepath"
	"testing"
)

func TestOneDep(t *testing.T) {
	cgfile := "one-dep.yaml"
	got, err := loadNRun(cgfile)
	if err != nil {
		t.Errorf("Can't load call graph %s: %v", cgfile, err)
		t.FailNow()
	}
	want := "root"
	if got[0] != want {
		t.Errorf("%s => %q, want %q", cgfile, got, want)
	}
}

func loadNRun(cg string) ([]string, error) {
	dj := New()
	cgfile, err := filepath.Abs(filepath.Join("examples", cg))
	if err != nil {
		return []string{}, err
	}
	err = dj.FromFile(cgfile)
	if err != nil {
		return []string{}, err
	}
	dj.Run()
	return dj.CallSeq(), nil
}
