package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/thomgray/codebook/util"
)

// Config ...
type Config struct {
	SearchPaths []string
	NotePaths   []string
}

// MakeConfig ...
func MakeConfig() *Config {
	return (&Config{}).Init()
}

// Init ...
func (c *Config) Init() *Config {
	c.SearchPaths = loadSeachPaths()
	c.NotePaths = loadNotePaths(c.SearchPaths)
	return c
}

func loadSeachPaths() []string {
	bytes, _ := util.ReadFile(ConfigNotePathsPath())
	paths := util.ReadLines(bytes)
	return paths
}

func loadNotePaths(searchPaths []string) []string {
	files := make([]string, 0)

	for _, sp := range searchPaths {
		f, err := ioutil.ReadDir(sp)
		if err == nil {
			for _, file := range f {
				if file.Mode().IsRegular() && filepath.Ext(file.Name()) == ".md" {
					files = append(files, fmt.Sprintf("%s/%s", sp, file.Name()))
				}
				log.Printf("File %s\n", fmt.Sprintf("%s/%s", sp, file.Name()))
			}
		}
		// don't walk as it recurs directories.
		// filepath.Walk(sp, func(path string, info os.FileInfo, err error) error {
		// 	if !info.IsDir() && filepath.Ext(path) == ".md" {
		// 		files = append(files, path)
		// 	}
		// 	return nil
		// })
	}
	return files
}

var _homedir *string = nil

// ConfigDirectory ...
func ConfigDirectory() string {
	return fmt.Sprintf("%s/.codebook", GetAppConfig().HomeDir)
}

// ConfigNotePathsPath ...
func ConfigNotePathsPath() string {
	return fmt.Sprintf("%s/paths", ConfigDirectory())
}

func (c *Config) AddSearchPath(sp string) {
	c.SearchPaths = append(c.SearchPaths, sp)
	c.updateSearchPathConfig()
}

func (c *Config) updateSearchPathConfig() {
	serlaised := []byte(strings.Join(c.SearchPaths, "\n"))
	ioutil.WriteFile(ConfigNotePathsPath(), serlaised, 0644)
}

func (c *Config) RemoveSearchPath(sp string) {
	for i, p := range c.SearchPaths {
		if p == sp {
			newSp := append(c.SearchPaths[:i], c.SearchPaths[i+1:]...)
			c.SearchPaths = newSp
		}
	}
	c.updateSearchPathConfig()
}

func (c *Config) ReloadNotes() {
	c.NotePaths = loadNotePaths(c.SearchPaths)
}
