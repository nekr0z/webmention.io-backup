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
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var (
	update = flag.Bool("update", false, "update .golden files")
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

func TestReadFile(t *testing.T) {
	tt := map[string]struct {
		inFile     string
		goldenFile string
	}{
		"api/mentions":     {"page.json", "page_processed.json"},
		"api/mentions.jf2": {"jf2.json", "jf2_processed.json"},
		"simple list":      {"single_file.json", "single_file_processed.json"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			wantF := filepath.Join("testdata", tc.inFile)
			golden := filepath.Join("testdata", tc.goldenFile)
			got, err := readFile(wantF)
			if err != nil {
				t.Fatal(err)
			}
			writeAndCompare(t, got, golden)
		})
	}
}

func TestReadFileErr(t *testing.T) {
	_, err := readFile("testdata")
	if err == nil {
		t.Fatalf("want error, got nil")
	}
}

func writeAndCompare(t *testing.T, mm []interface{}, fn string) {
	t.Helper()
	wantF := filepath.Join(fn)
	gotF := filepath.Join("testdata", "test_output.json")
	err := writeFile(mm, gotF)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(gotF)
	got, err := ioutil.ReadFile(gotF)
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, got, wantF)
}

func TestGetNew(t *testing.T) {
	pageF := filepath.Join("testdata", "page.json")
	golden := filepath.Join("testdata", "page_processed.json")
	page, _ := ioutil.ReadFile(pageF)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("since_id") != "20" {
			t.Fatalf("expected since_id not in request")
		}
		if r.URL.Query().Get("page") == "" || r.URL.Query().Get("page") == "0" {
			fmt.Fprintf(w, "%s", page)
		} else {
			fmt.Fprintf(w, "%s", `{"links":[]}`)
		}
	}))
	defer ts.Close()

	got, err := getNew(ts.URL, 20)
	if err != nil {
		t.Fatal(err)
	}

	writeAndCompare(t, got, golden)
}

func assertGolden(t *testing.T, actual []byte, golden string) {
	t.Helper()

	if *update {
		if _, err := os.Stat(golden); os.IsNotExist(err) {
			if err := ioutil.WriteFile(golden, actual, 0644); err != nil {
				t.Fatal(err)
			}
		} else {
			t.Log("file", golden, "exists, remove it to record new golden result")
		}
	}
	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Error("no file:", golden)
	}

	if !bytes.Equal(actual, expected) {
		t.Fatalf("want:\n%s\ngot:\n%s\n", expected, actual)
	}
}
