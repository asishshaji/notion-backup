package actions

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ExtractorAction struct{}

func (ExtractorAction) String() string {
	return "ExtractorAction"
}

// used to extract the downloaded zip
// the download path is available from the DownloaderAction
func (eA ExtractorAction) Act(s *SharedData) error {
	outerZipFile, err := zip.OpenReader(s.DownloadedFilePath)
	if err != nil {
		return err
	}
	defer outerZipFile.Close()

	// Create a buffer to hold the contents of the inner ZIP file
	var innerZipBuffer bytes.Buffer

	// Find and extract the inner ZIP file from the outer ZIP file
	var innerZipFile *zip.File
	for _, file := range outerZipFile.File {
		innerZipFile = file
	}

	if innerZipFile == nil {
		return fmt.Errorf("inner zip file not found")
	}

	innerZipReader, err := innerZipFile.Open()
	if err != nil {
		return err
	}
	defer innerZipReader.Close()
	_, err = io.Copy(&innerZipBuffer, innerZipReader)
	if err != nil {
		return err
	}

	// Unzip the contents of the inner ZIP file
	innerZipReaderAt := bytes.NewReader(innerZipBuffer.Bytes())
	innerZip, err := zip.NewReader(innerZipReaderAt, int64(innerZipBuffer.Len()))
	if err != nil {
		return err
	}

	extracted_dir := fmt.Sprintf("extracted_%s", s.ExportType)
	removePreviousAndCreate(extracted_dir)

	for _, file := range innerZip.File {
		innerFile, _ := file.Open()
		defer innerFile.Close()

		extractPath := filepath.Join(extracted_dir, file.Name)
		if err := os.MkdirAll(filepath.Dir(extractPath), 0770); err != nil {
			return err
		}
		extractFile, err := os.Create(extractPath)
		if err != nil {
			return err
		}
		fmt.Printf("creating file : %s\n", file.Name)
		defer extractFile.Close()

		// Copy the file contents
		if _, err := io.Copy(extractFile, innerFile); err != nil {
			return err
		}

	}

	fmt.Printf("extract completed %s\n", s.DownloadedFilePath)
	return nil
}

func removePreviousAndCreate(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat directory :%s", err)
	}

	if info.IsDir() {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove directory: %s", err)
		}
		fmt.Printf("removed directory :%s\n", path)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}
