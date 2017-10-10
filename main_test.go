package main

import (
	"fmt"
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
		t.Errorf("%s => %q, want %q", cgfile, got[0], want)
	}
	want = "j2"
	if got[1] != want {
		t.Errorf("%s => %q, want %q", cgfile, got[1], want)
	}
}

func TestTwoDep(t *testing.T) {
	cgfile := "two-dep.yaml"
	got, err := loadNRun(cgfile)
	if err != nil {
		t.Errorf("Can't load call graph %s: %v", cgfile, err)
		t.FailNow()
	}
	want := "root"
	if got[0] != want {
		t.Errorf("%s => %q, want %q", cgfile, got[0], want)
	}
	want0 := "j2"
	want1 := "j3"
	if !((got[1] == want0 && got[2] == want1) ||
		(got[1] == want1 && got[2] == want0)) {
		t.Errorf("%s => %q, want %q or %q", cgfile, got[1], want0, want1)
	}
}

func TestDiamond(t *testing.T) {
	cgfile := "diamond.yaml"
	got, err := loadNRun(cgfile)
	if err != nil {
		t.Errorf("Can't load call graph %s: %v", cgfile, err)
		t.FailNow()
	}
	want := "root"
	if got[0] != want {
		t.Errorf("%s => %q, want %q", cgfile, got[0], want)
	}
	want0 := "j2"
	want1 := "j3"
	if !((got[1] == want0 && got[2] == want1) ||
		(got[1] == want1 && got[2] == want0)) {
		t.Errorf("%s => %q, want %q or %q", cgfile, got[1], want0, want1)
	}
	want = "j4"
	if got[3] != want {
		t.Errorf("%s => %q, want %q", cgfile, got[3], want)
	}
}

func TestDeep(t *testing.T) {
	cgfile := "deep.yaml"
	got, err := loadNRun(cgfile)
	if err != nil {
		t.Errorf("Can't load call graph %s: %v", cgfile, err)
		t.FailNow()
	}
	want := "root"
	if got[0] != want {
		t.Errorf("%s => %q, want %q", cgfile, got[0], want)
	}
	want0 := "j2"
	want1 := "j3"
	if !((got[1] == want0 && got[2] == want1) ||
		(got[1] == want1 && got[2] == want0)) {
		t.Errorf("%s => %q, want %q or %q", cgfile, got[1], want0, want1)
	}
	want = "j4"
	if got[3] != want {
		t.Errorf("%s => %q, want %q", cgfile, got[3], want)
	}
	want = "j5"
	if got[4] != want {
		t.Errorf("%s => %q, want %q", cgfile, got[4], want)
	}
}

func TestSeq(t *testing.T) {
	cgfile := "seq.yaml"
	got, err := loadNRun(cgfile)
	if err != nil {
		t.Errorf("Can't load call graph %s: %v", cgfile, err)
		t.FailNow()
	}
	want := "root"
	if got[0] != want {
		t.Errorf("%s => %q, want %q", cgfile, got[0], want)
	}
	for i, g := range got[1:] {
		want = fmt.Sprintf("j%d", i+2)
		if g != want {
			t.Errorf("%s => %q, want %q", cgfile, g, want)
		}
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
