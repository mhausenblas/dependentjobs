package main

import (
	"path/filepath"
	"strings"
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
	idroot, _, endroot := extract(got[0])
	if idroot != want {
		t.Errorf("%s => %q, want %q", cgfile, idroot, want)
	}
	idj2, startj2, _ := extract(got[1])
	if startj2 <= endroot {
		t.Errorf("%s => %q before %q", cgfile, idj2, idroot)
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
	idroot, _, endroot := extract(got[0])
	if idroot != want {
		t.Errorf("%s => %q, want %q", cgfile, idroot, want)
	}
	ida, starta, _ := extract(got[1])
	if starta <= endroot {
		t.Errorf("%s => %q before %q", cgfile, ida, idroot)
	}
	idb, startb, _ := extract(got[2])
	if startb <= endroot {
		t.Errorf("%s => %q before %q", cgfile, idb, idroot)
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
	idroot, _, endroot := extract(got[0])
	if idroot != want {
		t.Errorf("%s => %q, want %q", cgfile, idroot, want)
	}
	ida, starta, enda := extract(got[1])
	if starta <= endroot {
		t.Errorf("%s => %q before %q", cgfile, ida, idroot)
	}
	idb, startb, endb := extract(got[2])
	if startb <= endroot {
		t.Errorf("%s => %q before %q", cgfile, idb, idroot)
	}
	idlast, startlast, _ := extract(got[3])
	if startlast <= enda {
		t.Errorf("%s => %q before %q", cgfile, idlast, ida)
	}
	if startlast <= endb {
		t.Errorf("%s => %q before %q", cgfile, idlast, idb)
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
	idroot, _, endroot := extract(got[0])
	if idroot != want {
		t.Errorf("%s => %q, want %q", cgfile, idroot, want)
	}
	ida, starta, enda := extract(got[1])
	if starta <= endroot {
		t.Errorf("%s => %q before %q", cgfile, ida, idroot)
	}
	idb, startb, endb := extract(got[2])
	if startb <= endroot {
		t.Errorf("%s => %q before %q", cgfile, idb, idroot)
	}
	idc, startc, endc := extract(got[3])
	if startc <= enda {
		t.Errorf("%s => %q before %q", cgfile, idc, idroot)
	}
	if startc <= endb {
		t.Errorf("%s => %q before %q", cgfile, idc, idroot)
	}
	idlast, startlast, _ := extract(got[4])
	if startlast <= endroot {
		t.Errorf("%s => %q before %q", cgfile, idlast, idroot)
	}
	if startlast <= enda {
		t.Errorf("%s => %q before %q", cgfile, idlast, ida)
	}
	if startlast <= endb {
		t.Errorf("%s => %q before %q", cgfile, idlast, idb)
	}
	if startlast <= endc {
		t.Errorf("%s => %q before %q", cgfile, idlast, idc)
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
	idroot, startroot, endroot := extract(got[0])
	if idroot != want {
		t.Errorf("%s => %q, want %q", cgfile, idroot, want)
	}
	idprev, _, endprev := idroot, startroot, endroot
	for _, g := range got[1:] {
		idcurrent, startcurrent, endcurrent := extract(g)
		if startcurrent <= endprev {
			t.Errorf("%s => %q before %q", cgfile, idcurrent, idprev)
		}
		idprev, endprev = idcurrent, endcurrent
	}
}

func TestTree(t *testing.T) {
	cgfile := "tree.yaml"
	got, err := loadNRun(cgfile)
	if err != nil {
		t.Errorf("Can't load call graph %s: %v", cgfile, err)
		t.FailNow()
	}
	want := "root"
	idroot, _, endroot := extract(got[0])
	if idroot != want {
		t.Errorf("%s => %q, want %q", cgfile, idroot, want)
	}
	ida, starta, enda := extract(got[1])
	if starta <= endroot {
		t.Errorf("%s => %q before %q", cgfile, ida, idroot)
	}
	idb, startb, endb := extract(got[2])
	if startb <= endroot {
		t.Errorf("%s => %q before %q", cgfile, idb, idroot)
	}
	idleafa, startleafa, _ := extract(got[3])
	if idleafa == "j4" && ida == "j2" && startleafa <= enda {
		t.Errorf("%s => %q before %q", cgfile, idleafa, ida)
	}
	if idleafa == "j4" && idb == "j2" && startleafa <= endb {
		t.Errorf("%s => %q before %q", cgfile, idleafa, idb)
	}
	idleafb, startleafb, _ := extract(got[3])
	if idleafb == "j5" && ida == "j3" && startleafb <= enda {
		t.Errorf("%s => %q before %q", cgfile, idleafb, ida)
	}
	if idleafa == "j5" && idb == "j3" && startleafb <= endb {
		t.Errorf("%s => %q before %q", cgfile, idleafb, idb)
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
	dj.Complete()
	return dj.CallSeq(), nil
}

func extract(csentry string) (id, start, end string) {
	res := strings.Split(csentry, " ")
	id, start, end = res[0], res[1], res[2]
	return
}
