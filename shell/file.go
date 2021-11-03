package shell

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func FileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil
}

func GetFileMtime(fpath string) int64 {
	info, err := os.Stat(fpath)
	if err != nil {
		return 0
	}
	return info.ModTime().Unix()
}

func LoadFile(fpath string) (string, error) {
	fp, err := os.Open(fpath)
	if err != nil {
		return "", errors.New("Error reading file: " + fpath)
	}
	defer fp.Close()

	lines := []string{}
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n"), nil
}

func PathRelativeTo(fpath string, relativeTo string) string {
	dirpath, err := filepath.Abs(path.Dir(relativeTo))
	if err != nil {
		return fpath
	}
	return path.Join(dirpath, fpath)
}
