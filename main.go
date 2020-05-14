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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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
	fmt.Println("vim-go")
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
	b, err := json.Marshal(file)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fn, b, 0644)
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
