package models

import (
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
	// Get the display name for this folder
	splitPath := strings.Split(folder, string(os.PathSeparator))
	displayName := splitPath[len(splitPath)-1] // Get the last part of the path
	displayName = strings.ToUpper(displayName[:1]) + displayName[1:]

	rootFolder := Folder{
		Name:       displayName,
		Files:      []File{},
		Subfolders: []Folder{},
	}

	// Read only the immediate contents of this directory
	entries, err := os.ReadDir(folder)
	if err != nil {
		return Folder{}, err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(folder, entry.Name())

		if entry.IsDir() {
			// Recursively build the subfolder tree
			subFolder, err := FileTree(entryPath)
			if err != nil {
				return Folder{}, err
			}

			rootFolder.Subfolders = append(rootFolder.Subfolders, subFolder)

		} else if strings.HasSuffix(entry.Name(), ".md") {
			// Only process .md files
			fileName := strings.TrimSuffix(entry.Name(), ".md")
			filePath := strings.TrimSuffix(entryPath, ".md")
			filePath = strings.TrimPrefix(filePath, "public/")

			file := File{
				Name: fileName,
				Path: filePath,
			}

			rootFolder.Files = append(rootFolder.Files, file)
		}
	}

	return rootFolder, nil
}
