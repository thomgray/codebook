package model

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/thomgray/codebook/config"
	"github.com/thomgray/codebook/util"
	"golang.org/x/net/html"
)

type Location struct {
	BaseDir              string
	RelativePath         []string
	RelativePathWithName string
}

type File struct {
	Path      string
	Extension string
	Name      string
	Content   []byte
	Locations []Location
	Body      *html.Node
	Document  *Document
}

type FileManager struct {
	Files           []*File
	Config          *config.Config
	CurrentLocation *Location
}

func MakeFileManager(config *config.Config) *FileManager {
	fm := FileManager{
		Config: config,
	}

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

func (fm *FileManager) SetLocation(location *Location) {
	fm.CurrentLocation = location
}

func (fm *FileManager) SuggestPaths(fragment string) []string {
	res := make([]string, 0)
	upTo := filepath.Dir(fragment)
	pathPieces := strings.Split(fragment, string(os.PathSeparator))
	remainder := pathPieces[len(pathPieces)-1]
	for _, sp := range fm.Config.SearchPaths {
		maybeDirPath := filepath.Join(sp, upTo)
		if pathInfo, exists := util.PathExists(maybeDirPath); exists && pathInfo.IsDir() {
			filesInDir := util.ListFilesShort(maybeDirPath)
			for _, file := range filesInDir {
				if fileIsRelevant(file) {
					name := file.Name()
					if strings.HasPrefix(name, remainder) {
						var toAppend string = completionString(file, upTo)
						log.Println(toAppend)
						res = append(res, toAppend)
					}
				}
			}
		}
	}
	log.Printf("Autocompletion suggestions = %v", res)
	return res
}

func fileIsRelevant(info os.FileInfo) bool {
	if info.IsDir() {
		return true
	}
	switch filepath.Ext(info.Name()) {
	case ".md":
		return true
		// case ".html"
	}
	return false
}

func completionString(info os.FileInfo, base string) string {
	if info.IsDir() {
		return filepath.Join(base, info.Name()) + string(os.PathSeparator)
	}
	return filepath.Join(
		base,
		strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
	)
}

func (fm *FileManager) TraversePath(path string) []*File {
	files := make([]*File, 0)
	dir := filepath.Dir(path)
	fileName := filepath.Base(path)
	for _, sp := range fm.Config.SearchPaths {
		fullDirPath := filepath.Join(sp, dir)
		if _, exists := util.PathExists(fullDirPath); exists {
			filesInDir := util.ListFilesShort(fullDirPath)
			for _, fileInDir := range filesInDir {
				fileInDirName := fileInDir.Name()
				ext := filepath.Ext(fileInDirName)
				fileWithoutExt := strings.TrimSuffix(fileInDirName, ext)
				log.Println(fileWithoutExt)
				if strings.EqualFold(fileName, fileWithoutExt) {
					fullFilePath := filepath.Join(fullDirPath, fileInDirName)
					file := LoadCodeFile(fullFilePath)
					if file != nil {
						files = append(files, file)
						log.Println("Matched a file!")
					}
				}
			}
		}
	}

	return files
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
			// md := DocumentFromNode(node, filename)
			file.Body = node
			// file.Document = md
		}
	} else if extn == ".html" {
		node, err := util.HtmlToNode(fc)
		if err == nil {
			// md := DocumentFromNode(node, filename)
			// file.Document = md
			file.Body = node
		}
	}
	return &file
}
