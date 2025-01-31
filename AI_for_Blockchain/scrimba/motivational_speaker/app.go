package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"

	openai "github.com/openai/openai-go"
	openoptions "github.com/openai/openai-go/option"
	"github.com/peartes/scrimba/motivational_speaker/config"
)

var client openai.Client
var ctx context.Context
var db *MemStore

func init() {
	client = *openai.NewClient(openoptions.WithAPIKey(config.GetOpenAIKey()), openoptions.WithOrganization(config.GetOpenAIOrganization()), openoptions.WithProject(config.GetOpenAIProject()))
	ctx = context.Background()
	db = NewMemStore("data/db.json")
}

// / RunApp runs the application
// / It uploads the fine tune data file to the openai server
// / fileName is the name of the file to upload. Must be a jsonl file and in the data directory
// / It returns an error if any
func RunApp(filename string) error {
	isUploaded, fd, err := GetFileUploadFD(filename)
	if err != nil {
		return fmt.Errorf("error while getting file upload descriptor %w ", err)
	}
	if !isUploaded {
		return fmt.Errorf("error while uploading file %w ", err)
	}
	fmt.Println("File uploaded successfully with file descriptor: ", fd)
	return nil
}

/*
This function checks if a file has been uploaded to an openai organization before.
It returns a (isUploaded, file upload descriptor, error if any)
*/
func GetFileUploadFD(filename string) (bool, string, error) {
	fd, err := os.Getwd()
	if err != nil {
		return false, "", fmt.Errorf("error while getting working directory %w ", err)
	}
	file, err := os.OpenFile(path.Join(fd, "data", filename), os.O_RDONLY, 0644)
	if err != nil {
		return false, "", fmt.Errorf("error while opening file %w ", err)
	}
	defer file.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return false, "", fmt.Errorf("error while hashing file %w ", err)
	}

	data := hasher.Sum([]byte("file"))

	checksum := fmt.Sprintf("%x", sha256.Sum256(data))

	f, exist := db.Get(checksum)
	if exist {
		return true, f, nil
	} else {
		fd, err := uploadFileToOpenAI(file)
		if err != nil {
			return false, "", fmt.Errorf("error while uploading file %w ", err)
		}
		db.Set(checksum, fd)
		return true, fd, nil
	}
}

/*
Uploads a file to the openai server.
Filename is the name of the file to upload. Must be a jsonl file and in the data directory
It returns (file upload descriptor, error if any)
*/
func uploadFileToOpenAI(file *os.File) (string, error) {
	// Reset the read pointer to the beginning
	_, err := file.Seek(0, 0) // Move back to start
	if err != nil {
		return "", fmt.Errorf("error while seeking file %w ", err)
	}
	uploadFD, err := client.Files.New(ctx, openai.FileNewParams{
		File:    openai.Raw[io.Reader](file),
		Purpose: openai.Raw[openai.FilePurpose](openai.FilePurposeFineTune),
	})
	if err != nil {
		return "", fmt.Errorf("error while uploading file %w ", err)
	}
	return uploadFD.ID, nil
}
