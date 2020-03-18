package model

import (
	"path/filepath"
	"strings"

	"github.com/thomgray/codebook/util"
)

type File struct {
	Path      string
	Extension string
	Name      string
	Content   []byte
	Document  *Document
}

type FileManager struct {
	Files []*File
}

func MakeFileManager() *FileManager {
	fm := FileManager{}

	return &fm
}

func (fm *FileManager) LoadFiles(filepaths []string) {
	files := make([]*File, 0)
	for _, path := range filepaths {
		f := LoadCodeFile(path)
		if f != nil {
			files = append(files, f)
		}
	}
	fm.Files = files
}

func LoadCodeFile(path string) *File {
	extn := filepath.Ext(path)
	_, n := filepath.Split(path)
	filename := strings.TrimSuffix(n, extn)

	file := File{
		Path:      path,
		Extension: extn,
		Name:      filename,
	}
	fc, _ := util.ReadFile(path)
	file.Content = fc

	if extn == ".md" {
		node, err := util.MarkdownToNode(fc)
		if err == nil {
			md := DocumentFromNode(node)
			file.Document = md
		}
	}
	return &file
}
