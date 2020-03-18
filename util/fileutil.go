package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// ListFiles ...
func ListFiles(path string) []string {
	res := make([]string, 0)
	files, err := ioutil.ReadDir(path)
	// err := filepath.Walk(path, func(pth string, info os.FileInfo, err error) error {
	// 	if !info.IsDir() {
	// 		res = append(res, path)
	// 	}
	// 	return nil
	// })

	Check(err)

	for _, f := range files {
		log.Print(f.Name())
		if f.Mode().IsRegular() {
			res = append(res, fmt.Sprintf("%s/%s", path, f.Name()))
		}
	}

	return res
}

// ReadFile - reads a file contents. Empty bytes if file doesn't exist, second param is false is file doesn't exist.
// Panics on any other error
func ReadFile(path string) (data []byte, exists bool) {
	bytes, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return make([]byte, 0), false
	}
	Check(err)
	return bytes, true
}
