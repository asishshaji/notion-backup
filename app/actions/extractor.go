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

	extracted_dir := "extracted"
	os.MkdirAll(extracted_dir, 0755)

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
