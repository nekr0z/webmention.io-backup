package path

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// DirFromUrl returns a (relative to content directory) dir to store mentions.
func DirFromUrl(t string, prefixes []string) string {
	p, err := trimmedPath(t)
	if err != nil {
		return ""
	}

	p = trimOne(p, prefixes)
	p = path.Dir(p)
	p = strings.TrimPrefix(p, "/")
	return p
}

// FilenameFromUrl returns a filename with language extension inserted.
func FilenameFromUrl(t string, prefixes []string, filename string) string {
	p, err := trimmedPath(t)
	if err != nil {
		return ""
	}

	for _, pref := range prefixes {
		if strings.HasPrefix(p, pref) {
			base := filepath.Base(filename)
			ext := filepath.Ext(filename)
			base = strings.TrimSuffix(base, ext)
			ext = strings.TrimPrefix(ext, ".")
			return strings.Join([]string{base, pref, ext}, ".")
		}
	}

	return filename
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

func trimmedPath(u string) (string, error) {
	ur, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	p := strings.TrimPrefix(ur.Path, "/")

	return p, nil
}
