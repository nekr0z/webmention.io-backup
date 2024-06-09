package path_test

import (
	"testing"

	"evgenykuznetsov.org/go/webmention.io-backup/internal/path"
)

func TestDirFromUrl(t *testing.T) {
	testcases := []struct {
		url      string
		prefixes []string
		want     string
	}{
		{
			url:      "https://evgenykuznetsov.org/posts/2024/elevator/",
			prefixes: []string{"en"},
			want:     "posts/2024/elevator",
		},
		{
			url:      "https://evgenykuznetsov.org/en/posts/2021/covid/",
			prefixes: []string{"en", "ru"},
			want:     "posts/2021/covid",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.url, func(t *testing.T) {
			got := path.DirFromUrl(tc.url, tc.prefixes)
			if got != tc.want {
				t.Errorf("\nwant: %s,\n got: %s", tc.want, got)
			}
		})
	}
}

func TestFilenameFromUrl(t *testing.T) {
	testcases := []struct {
		url      string
		prefixes []string
		filename string
		want     string
	}{
		{
			url:      "https://evgenykuznetsov.org/posts/2024/elevator/",
			prefixes: []string{"en"},
			filename: "webmentions.json",
			want:     "webmentions.json",
		},
		{
			url:      "https://evgenykuznetsov.org/en/posts/2021/covid/",
			prefixes: []string{"en", "ru"},
			filename: "webmentions.json",
			want:     "webmentions.en.json",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.url, func(t *testing.T) {
			got := path.FilenameFromUrl(tc.url, tc.prefixes, tc.filename)
			if got != tc.want {
				t.Errorf("\nwant: %s,\n got: %s", tc.want, got)
			}
		})
	}
}
