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
var db *MemStore[JsonTypes]

func init() {
	client = *openai.NewClient(openoptions.WithAPIKey(config.GetOpenAIKey()), openoptions.WithOrganization(config.GetOpenAIOrganization()), openoptions.WithProject(config.GetOpenAIProject()))
	ctx = context.Background()
	db = NewMemStore[JsonTypes]("data/db.json")
}

// / RunApp runs the application
// / It uploads the fine tune data file to the openai server
// / fileName is the name of the file to upload. Must be a jsonl file and in the data directory
// / It returns an error if any
func RunApp(filename string) error {
	fd, err := GetFileUploadFD(filename)
	if err != nil {
		return fmt.Errorf("error while getting file upload descriptor %w ", err)
	}
	fmt.Println("File uploaded successfully with id: ", fd)

	modelId, err := GetFineTunedModel(filename, openai.FineTuningJobNewParamsModelGPT3_5Turbo)
	if err != nil {
		return fmt.Errorf("error while checking if fine tuned model exists %w ", err)
	}
	fmt.Println("Fine tuned model already exists with id: ", modelId)
	return nil
}

/*
This function creates a fine tuned model with a particular doc
Returns the (model id, error if any)
*/
func GetFineTunedModel(trainingFile string, model openai.FineTuningJobNewParamsModel) (string, error) {
	fd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error while getting working directory %w ", err)
	}
	file, err := os.OpenFile(path.Join(fd, "data", trainingFile), os.O_RDONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("error while opening file %w ", err)
	}
	defer file.Close()

	checksum, err := createChecksum(file)
	if err != nil {
		return "", fmt.Errorf("error while creating checksum %w ", err)
	}

	f, exist := db.Get(checksum, "fine_tuned_models")

	if exist {
		res, ok := (*f).(string)
		if ok {
			return res, nil
		}
		return "", fmt.Errorf("error while getting fine tuned model from db %w ", err)
	} else {
		fd, err := GetFileUploadFD(trainingFile)
		if err != nil {
			return "", fmt.Errorf("error while getting file upload descriptor %w ", err)
		}
		modelId, err := CreateFineTunedModel(openai.FineTuningJobNewParams{
			Model:        openai.Raw[openai.FineTuningJobNewParamsModel](openai.FineTuningJobNewParamsModelGPT3_5Turbo),
			TrainingFile: openai.String(fd),
			Suffix:       openai.String("motivational_speaker"),
		})
		if err != nil {
			return "", fmt.Errorf("error while creating a new fine tuned model %w ", err)
		}
		db.Set(checksum, modelId, "fine_tuned_models")
		return modelId, nil
	}
}

/*
This function checks if there exists a finetuned model trained with a particular doc
Returns the (model id, error if any)
*/
func CreateFineTunedModel(opts openai.FineTuningJobNewParams) (string, error) {
	fineTune, err := client.FineTuning.Jobs.New(ctx, opts)
	if err != nil {
		return "", fmt.Errorf("error while creating fine tuning job %w ", err)
	}
	return fineTune.ID, nil
}

/*
This function checks if a file has been uploaded to an openai organization before.
It returns a (isUploaded, file upload descriptor, error if any)
*/
func GetFileUploadFD(filename string) (string, error) {
	fd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error while getting working directory %w ", err)
	}
	file, err := os.OpenFile(path.Join(fd, "data", filename), os.O_RDONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("error while opening file %w ", err)
	}
	defer file.Close()

	checksum, err := createChecksum(file)
	if err != nil {
		return "", fmt.Errorf("error while creating checksum %w ", err)
	}

	f, exist := db.Get(checksum, "uploads")

	if exist {
		res, ok := (*f).(string)
		if ok {
			return res, nil
		}
		return "", fmt.Errorf("error while getting file upload descriptor from db %w ", err)
	} else {
		fd, err := uploadFileToOpenAI(file)
		if err != nil {
			return "", fmt.Errorf("error while uploading file %w ", err)
		}
		db.Set(checksum, fd, "uploads")
		return fd, nil
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

func createChecksum(file *os.File) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("error while hashing file %w ", err)
	}

	data := hasher.Sum([]byte("file"))

	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}
