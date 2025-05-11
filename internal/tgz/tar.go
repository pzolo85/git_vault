package tgz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func CreateTarGz(source, target string) error {
	// Create the output file
	outFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer outFile.Close()

	// Create a gzip writer
	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	// Create a tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Walk through the source directory/file and add it to the tar.gz archive
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		// Create a tar header for the current file/directory
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return fmt.Errorf("error creating tar header for %s: %w", path, err)
		}

		// Update the header name to maintain the correct relative path
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return fmt.Errorf("error updating header name for %s: %w", path, err)
		}

		// Write the header to the tar file
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("error writing header for %s: %w", path, err)
		}

		// If it's a directory, we don't need to write its content
		if info.IsDir() {
			return nil
		}

		// Open the file to copy its content
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file %s: %w", path, err)
		}
		defer file.Close()

		// Copy the file content to the tar writer
		if _, err := io.Copy(tarWriter, file); err != nil {
			return fmt.Errorf("error writing file content for %s: %w", path, err)
		}

		return nil
	})
}

// extractTarGz extracts a .tar.gz file to the specified destination directory
func ExtractTarGz(source, destination string) error {
	// Open the .tar.gz file
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// Create a gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzReader)

	// Iterate through the files in the tar archive
	for {
		// Read the next tar header
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of archive
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar header: %w", err)
		}

		// Determine the path for the extracted file
		targetPath := filepath.Join(destination, header.Name)

		// Handle directory creation
		switch header.Typeflag {
		case tar.TypeDir:
			// Create the directory
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
		case tar.TypeReg:
			// Ensure the directory for the file exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory for file %s: %w", targetPath, err)
			}

			// Create the file
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}
			defer outFile.Close()

			// Copy the file content from the tar reader
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to write file content to %s: %w", targetPath, err)
			}
		default:
			// Ignore other file types (e.g., symlinks, devices, etc.)
			fmt.Printf("Skipping unknown file type: %s\n", header.Name)
		}
	}

	return nil
}

// GetLastTgzEncFile lists all files in the current directory ending with ".tgz.enc"
// and returns the last one in alphabetical order.
func GetLastTgzEncFile() (string, error) {
	// Read the current directory
	files, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	var tgzEncFiles []string

	// Filter files with ".tgz.enc" extension
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".tgz.enc") {
			tgzEncFiles = append(tgzEncFiles, file.Name())
		}
	}

	// If no matching files are found, return an error
	if len(tgzEncFiles) == 0 {
		return "", fmt.Errorf("no files found with .tgz.enc extension")
	}

	// Sort files alphabetically
	sort.Strings(tgzEncFiles)

	// Return the last file in alphabetical order
	return tgzEncFiles[len(tgzEncFiles)-1], nil
}
