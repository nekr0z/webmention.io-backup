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
	"strings"
	"testing"
	"time"
)

var (
	update = flag.Bool("update", false, "update .golden files")
)

func TestFindLatest(t *testing.T) {
	tt := map[string]struct {
		file string
		want int
	}{
		"legacy": {"page.json", 792685},
		"jf2":    {"jf2.json", 1183052},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			mm, err := readFile(filepath.Join("testdata", tc.file))
			if err != nil {
				t.Fatal(err)
			}
			got := findLast(mm)
			if got != tc.want {
				t.Fatalf("want %v, got %v", got, tc.want)
			}
		})
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
			writeAndCompare(t, got, cfg{tlo: false}, golden)
		})
	}
}

func TestReadFileErr(t *testing.T) {
	_, err := readFile("testdata")
	if err == nil {
		t.Fatalf("want error, got nil")
	}
}

func TestSaveToDirs(t *testing.T) {
	cdir := filepath.Join("testdata", "site", "content")
	mib := filepath.Join(cdir, "posts", "2020", "microblog-is-bad")
	if err := os.MkdirAll(mib, 0777); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(filepath.Join("testdata", "site"))

	ex, err := readFile(filepath.Join("testdata", "existing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := writeFile(ex, cfg{filename: filepath.Join(mib, "webmentions.json"), useJF2: true, tlo: false}); err != nil {
		t.Fatal(err)
	}

	mm, err := readFile(filepath.Join("testdata", "page.json"))
	if err != nil {
		t.Fatal(err)
	}
	c := cfg{squashLeft: []string{"en"}, contentDir: cdir, filename: "webmentions.json", timestamp: true}

	if err := os.Chmod(cdir, 0555); err != nil {
		t.Fatal(err)
	}
	if err := saveToDirs(mm, c); err == nil {
		t.Fatalf("expected error on read-only directory")
	}
	if err := os.Chmod(cdir, 0777); err != nil {
		t.Fatal(err)
	}

	if err := saveToDirs(mm, c); err != nil {
		t.Fatal(err)
	}

	m, err := readFile(filepath.Join("testdata", "site", "content", "posts", "2020", "microblog-is-bad", "webmentions.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 5 {
		t.Fatalf("unexpected number of mentions saved, want 5, got %d", len(m))
	}

	c.timestamp = false
	mm, err = readFile(filepath.Join("testdata", "jf2.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := saveToDirs(mm, c); err != nil {
		t.Fatal(err)
	}
	mm, err = readFile(filepath.Join(c.contentDir, c.filename))
	if err != nil {
		t.Fatal(err)
	}
	want, _ := time.Parse(time.RFC3339, "2020-05-05T14:54:13+00:00")
	got := getTimestamp(mm)
	if !want.Equal(got) {
		t.Fatalf("wrong timestamp, want %s, got %s", want, got)
	}
	if len(mm) != 18 {
		t.Fatalf("unexpected number of mentions saved, want 18, got %d", len(mm))
	}
}

func TestSaveToDirsErr(t *testing.T) {
	tests := map[string]struct {
		config cfg
		fail   string
	}{
		"non-existent dir": {cfg{contentDir: "nosuchdir", filename: "wm.json"}, "expected error on non-existent directory"},
		"no filename":      {cfg{contentDir: "testdata"}, "expected error on no filename given"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := saveToDirs([]interface{}{"empty mention"}, tc.config); err == nil {
				t.Fatalf(tc.fail)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tt := map[string]struct {
		config  cfg
		goldenF string
	}{
		"legacy":   {cfg{useJF2: false, tlo: true}, "legacy.out"},
		"JF2 feed": {cfg{useJF2: true, tlo: true}, "jf2.out"},
		"array":    {cfg{useJF2: false, tlo: false}, "array.out"},
		"pretty":   {cfg{pretty: true, tlo: true}, "pretty.out"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			mm := []interface{}{}
			writeAndCompare(t, mm, tc.config, filepath.Join("testdata", tc.goldenF))
		})
	}
}

func writeAndCompare(t *testing.T, mm []interface{}, c cfg, fn string) {
	t.Helper()
	wantF := filepath.Join(fn)
	c.filename = filepath.Join("testdata", "test_output.json")
	err := writeFile(mm, c)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(c.filename)
	got, err := ioutil.ReadFile(c.filename)
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, got, wantF)
}

func TestGetNew(t *testing.T) {
	tests := map[string]struct {
		arg  interface{}
		qKey string
		qVal string
	}{
		"ID":        {20, "since_id", "20"},
		"Timestamp": {time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), "since", "2018-01-01T00:00:00Z"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pageF := filepath.Join("testdata", "page.json")
			golden := filepath.Join("testdata", "page_processed.json")
			page, _ := ioutil.ReadFile(pageF)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get(tc.qKey) != tc.qVal {
					t.Fatalf("expected since key not in request")
				}
				if r.URL.Query().Get("page") == "" || r.URL.Query().Get("page") == "0" {
					fmt.Fprintf(w, "%s", page)
				} else {
					fmt.Fprintf(w, "%s", `{"links":[]}`)
				}
			}))
			defer ts.Close()

			got, err := getNew(ts.URL, tc.arg)
			if err != nil {
				t.Fatal(err)
			}

			writeAndCompare(t, got, cfg{tlo: false}, golden)
		})
	}
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

func TestEndpointUrl(t *testing.T) {
	tt := map[string]struct {
		config cfg
		want   string
	}{
		"token":  {cfg{token: "t0K3n"}, "?token=t0K3n"},
		"jf2":    {cfg{useJF2: true}, ".jf2"},
		"domain": {cfg{domain: "example.org"}, "?domain=example.org"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := strings.TrimPrefix(endpointUrl(tc.config), endpoint)
			if got != tc.want {
				t.Fatalf("want:\n%s\ngot:\n%s\n", tc.want, got)
			}
		})
	}
}

func TestGetTimestamp(t *testing.T) {
	bb := []byte(`[{"timestamp":"2021-06-07T22:21:17Z"}]`)
	mm, err := parsePage(bb)
	if err != nil {
		t.Fatal(err)
	}

	got := getTimestamp(mm)
	want := time.Date(2021, 6, 7, 22, 21, 17, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("want %s, got %s", want, got)
	}
}
