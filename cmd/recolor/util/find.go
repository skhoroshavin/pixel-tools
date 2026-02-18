package util

import (
	"os"
	"path/filepath"
	"strings"
)

func FindFiles(dir, suffix string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	res := make([]string, 0, len(files))
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		name := filepath.Join(dir, entry.Name())
		if strings.HasSuffix(name, suffix) {
			res = append(res, name)
		}
	}
	return res
}

func ReplaceSuffix(fname, srcSuffix, dstSuffix string) string {
	return strings.TrimSuffix(fname, srcSuffix) + dstSuffix
}

func FindPairs(dir, suffix1, suffix2 string) []string {
	tmp := FindFiles(dir, suffix1)
	res := make([]string, 0, len(tmp))
	for _, fname := range tmp {
		out := ReplaceSuffix(fname, suffix1, suffix2)
		finfo, err := os.Stat(out)
		if err != nil {
			continue
		}
		if finfo.IsDir() {
			continue
		}
		res = append(res, fname)
	}
	return res
}
