package shell

import "os"

// func loadFile(fpath string) ([]string, error) {
// 	fp, err := os.Open(fpath)
// 	if err != nil {
// 		return nil, errors.New("Error reading file: " + fpath)
// 	}
// 	defer fp.Close()

// 	lines := []string{}
// 	scanner := bufio.NewScanner(fp)
// 	for scanner.Scan() {
// 		lines = append(lines, scanner.Text())
// 	}
// 	return lines, nil
// }

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func GetFileMtime(filepath string) int64 {
	info, err := os.Stat(filepath)
	if err != nil {
		return 0
	}
	return info.ModTime().Unix()
}
