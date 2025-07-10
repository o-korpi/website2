package models

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name string
	Path string
}

type Folder struct {
	Name       string
	Files      []File
	Subfolders []Folder
}

func FileTree(folder string) (Folder, error) {

	splitPath := strings.Split(folder, string(os.PathSeparator))
	displayName := splitPath[len(splitPath)-1] // Get the last part of the path
	displayName = strings.ToUpper(displayName[:1]) + displayName[1:]
	rootFolder := Folder{
		Name:       displayName,
		Files:      []File{},
		Subfolders: []Folder{},
	}

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path == folder {
				// Skip the root folder itself
				return nil
			}

			log.Println("Visiting folder:", path)
			subFolder, err := FileTree(path)
			if err != nil {
				return err
			}

			rootFolder.Subfolders = append(rootFolder.Subfolders, subFolder)
		} else if strings.HasSuffix(info.Name(), ".md") {
			fileName := strings.TrimSuffix(info.Name(), ".md")
			filePath := strings.TrimPrefix(path, "public/")

			file := File{
				Name: fileName,
				Path: filePath,
			}

			rootFolder.Files = append(rootFolder.Files, file)
		}

		return nil
	})

	if err != nil {
		return Folder{}, err
	}

	log.Printf("Folder %v\n", rootFolder)

	return rootFolder, nil
}
