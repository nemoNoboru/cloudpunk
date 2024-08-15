package upload

import (
	"cloudpunk/cloud"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Function to generate a label from a file path
func generateLabel(path string) string {
	// Normalize the path to use the correct separators
	path = filepath.ToSlash(path)

	// Find the position of "/content/" in the path
	contentPos := strings.Index(path, "content/")

	// If "/content/" is not found, return the original path without modifications
	if contentPos == -1 {
		return ""
	}

	// Extract everything after "/content/"
	subPath := path[contentPos+len("content/"):]

	// Split the path into components
	components := strings.Split(subPath, "/")

	// Join the components with a dash, ignoring any empty strings
	var parts []string
	for _, component := range components {
		if component != "" {
			parts = append(parts, component)
		}
	}

	// Join parts to form the final label
	label := strings.Join(parts, "-")
	return strings.Split(label, ".")[0]
}

func UploadAllFilesInDir(dirPath string) error {
	// Walk through all files in the directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if it is a file
		if !info.IsDir() {
			// Read the file content
			label := generateLabel(path)
			fileData, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			log.Printf("uploading file %s, into %s", path, label)
			cloud.StorageLoad(label, fileData)
		}
		return nil
	})
	return err
}
