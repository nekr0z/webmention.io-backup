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
	"strconv"
)

const endpoint = "https://webmention.io/api/mentions"

type cfg struct {
	filename string
	token    string
	domain   string
	useJF2   bool
	tlo      bool
}

var version string = "custom"

func main() {
	fmt.Printf("webmention.io-backup version %s\n", version)

	config := cfg{}
	flag.StringVar(&config.filename, "f", "webmentions.json", "filename")
	flag.StringVar(&config.token, "t", "", "API token")
	flag.StringVar(&config.domain, "d", "", "domain to fetch webmentions for")
	flag.BoolVar(&config.useJF2, "jf2", false, "use JF2 endpoint instead of the classic one")
	flag.BoolVar(&config.tlo, "tlo", true, "wrap output in a top-level object (links list or feed)")
	flag.Parse()
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
		fmt.Printf("appending %d new webmentions\n", len(m))
		mm = append(mm, m...)
		err = writeFile(mm, config)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("saved %d webmentions to %s\n", len(mm), config.filename)
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
		id, ok := m["id"]
		if !ok {
			id = m["wm-id"]
		}
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
		mentions, ok := m["links"]
		if !ok {
			mentions = m["children"]
		}
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
