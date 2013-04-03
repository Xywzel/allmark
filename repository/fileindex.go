package repository

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func NewFileIndex(directory string) *FileIndex {

	files := getFiles(directory)

	return &FileIndex{
		Items: func() []*File {
			return files
		},

		path: directory,
	}
}

type FileIndex struct {
	Items func() []*File
	path  string
}

func (fileIndex *FileIndex) String() string {
	return fmt.Sprintf("File Index %s", fileIndex.path)
}

func (fileIndex *FileIndex) Path() string {
	return fileIndex.path
}

func getFiles(directory string) []*File {

	filesDirectoryEntries, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Printf("Cannot read files from directory %q. Error: %s", directory, err)
		return make([]*File, 0)
	}

	files := make([]*File, 0, 5)
	for _, directoryEntry := range filesDirectoryEntries {

		// recurse
		if directoryEntry.IsDir() {
			subDirectory := filepath.Join(directory, directoryEntry.Name())
			files = append(files, getFiles(subDirectory)...)
			continue
		}

		// append new file
		filePath := filepath.Join(directory, directoryEntry.Name())
		files = append(files, NewFile(filePath))
	}

	return files
}
