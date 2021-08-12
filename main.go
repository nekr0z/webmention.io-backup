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

//go:generate go run version_generate.go

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const endpoint = "https://webmention.io/api/mentions"

type cfg struct {
	filename   string
	token      string
	domain     string
	useJF2     bool
	tlo        bool
	contentDir string
	squashLeft []string
}

var version string = "custom"

func main() {
	fmt.Printf("webmention.io-backup version %s\n", version)

	config := cfg{}
	var sl string
	flag.StringVar(&config.filename, "f", "webmentions.json", "filename")
	flag.StringVar(&config.token, "t", "", "API token")
	flag.StringVar(&config.domain, "d", "", "domain to fetch webmentions for")
	flag.BoolVar(&config.useJF2, "jf2", false, "use JF2 endpoint instead of the classic one")
	flag.BoolVar(&config.tlo, "tlo", true, "wrap output in a top-level object (links list or feed)")
	flag.StringVar(&config.contentDir, "cd", "", "directory to look for structure in; if specified, attempts are made to save according to paths")
	flag.StringVar(&sl, "l", "", "list of top-level subdirs to drop while saving according to paths, comma-separated")
	flag.Parse()
	config.squashLeft = strings.Split(sl, ",")
	url := endpointUrl(config)

	mm, err := readFile(config.filename)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("found %d existing webmentions\n", len(mm))

	m, err := getNew(url, findLatest(mm))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(m) == 0 {
		fmt.Println("no new webmentions")
	} else {
		if config.contentDir != "" {
			if err := saveToDirs(m, config); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("appending %d new webmentions\n", len(m))
			mm = append(mm, m...)
			err = writeFile(mm, config)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("saved %d webmentions to %s\n", len(mm), config.filename)
			}
		}
	}

	fmt.Println("all done!")
}

func readFile(fn string) (mm []interface{}, err error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}
	mm, err = parsePage(data)
	return
}

func findLatest(mm []interface{}) (latest int) {
	for _, f := range mm {
		m, _ := f.(map[string]interface{})
		id := either(m, []string{"id", "wm-id"})
		if l, ok := id.(float64); ok {
			this := int(l)
			if this > latest {
				latest = this
			}
		}
	}
	return
}

func writeFile(mm []interface{}, c cfg) error {
	var bb bytes.Buffer
	var f interface{}
	if !c.tlo {
		f = mm
	} else {
		if c.useJF2 {
			f = struct {
				Type     string        `json:"type"`
				Name     string        `json:"name"`
				Children []interface{} `json:"children"`
			}{"feed", "Webmentions", mm}
		} else {
			f = struct {
				Links []interface{} `json:"links"`
			}{mm}
		}
	}
	enc := json.NewEncoder(&bb)
	enc.SetEscapeHTML(false)
	err := enc.Encode(f)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.filename, bb.Bytes(), 0644)
	return err
}

func saveToDirs(mm []interface{}, c cfg) (err error) {
	for _, m := range mm {
		if !saveToDir(m, c) {
			if err = saveToContentDir(m, c); err != nil {
				return
			}
		}
	}
	return
}

func saveToDir(m interface{}, c cfg) bool {
	if dir, ok := suggestDir(m, c); ok {
		c.contentDir = filepath.Join(c.contentDir, dir)
		if err := saveToContentDir(m, c); err == nil {
			return true
		}
	}
	return false
}

func saveToContentDir(m interface{}, c cfg) error {
	if c.filename == "" {
		return fmt.Errorf("no filename specified")
	}
	c.filename = filepath.Join(c.contentDir, c.filename)
	return saveToFile(m, c)
}

func saveToFile(m interface{}, c cfg) (err error) {
	mm, _ := readFile(c.filename)
	for _, exm := range mm {
		if sameMention(exm, m) {
			return
		}
	}
	mm = append(mm, m)
	fmt.Printf("Saving new mention to %s\n", c.filename)
	err = writeFile(mm, c)
	return
}

func sameMention(ma interface{}, mb interface{}) bool {
	mapa, oka := ma.(map[string]interface{})
	mapb, okb := mb.(map[string]interface{})
	if !oka || !okb {
		return false
	}

	sa := either(mapa, []string{"source", "wm-source"})
	sb := either(mapb, []string{"source", "wm-source"})
	oa, oka := sa.(string)
	ob, okb := sb.(string)
	if !oka || !okb || oa != ob {
		return false
	}

	sa = either(mapa, []string{"verified_date", "wm-received"})
	sb = either(mapb, []string{"verified_date", "wm-received"})
	oa, oka = sa.(string)
	ob, okb = sb.(string)
	ta, err := time.Parse(time.RFC3339, oa)
	if err != nil {
		return false
	}
	tb, err := time.Parse(time.RFC3339, ob)
	if err != nil {
		return false
	}
	if !oka || !okb || !ta.Equal(tb) {
		return false
	}

	return true
}

func suggestDir(m interface{}, c cfg) (dir string, ok bool) {
	mn, _ := m.(map[string]interface{})
	t := either(mn, []string{"target", "wm-target"})
	if tgt, ok := t.(string); ok {
		return dirFromUrl(tgt, c), true
	}
	return "", false
}

func dirFromUrl(t string, c cfg) string {
	u, err := url.Parse(t)
	if err != nil {
		return ""
	}

	p := strings.TrimPrefix(u.Path, "/")
	p = trimOne(p, c.squashLeft)
	p = path.Dir(p)
	p = strings.TrimPrefix(p, "/")
	return p
}

func trimOne(s string, vv []string) string {
	for _, l := range vv {
		r := strings.TrimPrefix(s, l)
		if s != r {
			return r
		}
	}
	return s
}

func getPage(url string) (mm []interface{}, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	mm, err = parsePage(b)
	return
}

func parsePage(b []byte) (mm []interface{}, err error) {
	var f interface{}
	err = json.Unmarshal(b, &f)

	// can be classic api/mentions with "links" array as a root object
	// or JF2 feed
	// or just an array of objects like we write it
	switch m := f.(type) {
	case map[string]interface{}:
		mentions := either(m, []string{"links", "children"})
		if mnts, ok := mentions.([]interface{}); ok {
			mm = mnts
		}
	case []interface{}:
		mm = m
	default:
		err = fmt.Errorf("could not parse JSON")
	}

	return
}

func getNew(uri string, latest int) (mm []interface{}, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}

	q := u.Query()
	q.Set("since_id", strconv.Itoa(latest))
	u.RawQuery = q.Encode()

	ok := true

	for i := 0; ok; i++ {
		m, err := getNextPage(u, i)
		if err != nil {
			return mm, err
		}
		mm = append(mm, m...)
		if len(m) == 0 {
			ok = false
		}
	}

	return
}

func getNextPage(u *url.URL, page int) (mm []interface{}, err error) {
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	mm, err = getPage(u.String())
	return
}

func endpointUrl(c cfg) string {
	q := url.Values{}
	vv := map[string]string{
		"token":  c.token,
		"domain": c.domain,
	}
	for k, v := range vv {
		if v != "" {
			q.Set(k, v)
		}
	}

	var e string
	if c.useJF2 {
		e = ".jf2"
	}
	ep := fmt.Sprintf("%s%s", endpoint, e)

	u, _ := url.Parse(ep)
	u.RawQuery = q.Encode()
	return u.String()
}

// either returns the first value it finds while iterating kk for key
func either(m map[string]interface{}, kk []string) interface{} {
	for _, k := range kk {
		if v, ok := m[k]; ok {
			return v
		}
	}
	return nil
}
