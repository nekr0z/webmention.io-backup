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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const endpoint = "https://webmention.io/api/mentions"

var (
	filename, token string
)

var version string = "custom"

type mention struct {
	Source       string          `json:"source"`
	Verified     bool            `json:"verified"`
	VerifiedDate time.Time       `json:"verified_date"`
	Id           int             `json:"id"`
	Private      bool            `json:"private"`
	Target       string          `json:"target"`
	Data         mentionData     `json:"data"`
	Activity     mentionActivity `json:"activity"`
}

type mentionData struct {
	Author      mentionAuthor `json:"author"`
	Url         string        `json:"url"`
	Name        string        `json:"name"`
	Content     string        `json:"content"`
	Published   time.Time     `json:"published"`
	PublishedTs int           `json:"published_ts"`
}

type mentionAuthor struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Photo string `json:"photo"`
}

type mentionActivity struct {
	TypeOf       string `json:"type"`
	Sentence     string `json:"sentence"`
	SentenceHtml string `json:"sentence_html"`
}

type mentionFile struct {
	Links *[]mention `json:"links"`
}

func main() {
	fmt.Printf("webmention.io-backup version %s\n", version)
	flag.StringVar(&filename, "f", "webmentions.json", "filename")
	flag.StringVar(&token, "t", "", "API token")
	flag.Parse()
	url := fmt.Sprintf("%s?token=%s", endpoint, token)

	mm, err := readFile(filename)
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
		err = writeFile(mm, filename)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("saved %d webmentions to %s\n", len(mm), filename)
		}
	}

	fmt.Println("all done!")
}

func readFile(fn string) (mm []mention, err error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}
	mm, err = parsePage(data)
	return
}

func findLatest(mm []mention) (latest int) {
	for _, m := range mm {
		if m.Id > latest {
			latest = m.Id
		}
	}
	return
}

func writeFile(mm []mention, fn string) error {
	file := mentionFile{&mm}
	var bb bytes.Buffer
	enc := json.NewEncoder(&bb)
	enc.SetEscapeHTML(false)
	err := enc.Encode(file)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fn, bb.Bytes(), 0644)
	return err
}

func getPage(url string) (mm []mention, err error) {
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

func parsePage(b []byte) (mm []mention, err error) {
	file := mentionFile{&mm}
	err = json.Unmarshal(b, &file)
	return
}

func getNew(uri string, latest int) (mm []mention, err error) {
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

func getNextPage(u *url.URL, page int) (mm []mention, err error) {
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	mm, err = getPage(u.String())
	return
}
