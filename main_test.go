// Copyright (C) 2020 Evgeny Kuznetsov (evgeny@kuznetsov.md)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/udhos/equalfile"
)

func TestFindLatest(t *testing.T) {
	mm, err := readFile(filepath.Join("testdata", "page.json"))
	if err != nil {
		t.Fatal(err)
	}
	got := findLatest(mm)
	want := 792685
	if got != want {
		t.Fatalf("want %v, got %v", got, want)
	}
}

func TestReadFileErr(t *testing.T) {
	_, err := readFile("testdata")
	if err == nil {
		t.Fatalf("want error, got nil")
	}
}

func TestWriteFile(t *testing.T) {
	wantF := filepath.Join("testdata", "page.json")
	gotF := filepath.Join("testdata", "page_processed.json")
	mm, err := readFile(wantF)
	if err != nil {
		t.Fatal(err)
	}
	err = writeFile(mm, gotF)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(gotF)
	if !filesEqual(t, wantF, gotF) {
		t.Fatalf("files don't match")
	}
}

func filesEqual(t *testing.T, newFile, oldFile string) bool {
	t.Helper()
	cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
	r1, err := os.Open(newFile)
	if err != nil {
		t.Fatalf("could not open %s", newFile)
	}
	defer r1.Close()
	r2, err := os.Open(oldFile)
	if err != nil {
		t.Fatalf("could not open %s", oldFile)
	}
	defer r2.Close()

	equal, err := cmp.CompareReader(r1, r2)
	if err != nil {
		t.Fatal(err)
	}
	return !equal
}
