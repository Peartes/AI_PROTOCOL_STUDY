package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
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

	model, err := GetFineTunedModel(filename, openai.FineTuningJobNewParamsModelGPT3_5Turbo)
	if err != nil {
		return fmt.Errorf("error while checking if fine tuned model exists %w ", err)
	}
	fmt.Printf("Model status: %s \nModel id: %s \nModel name: %s\n", model.JobStatus, model.ModelId, model.ModelName)
	return nil
}

/*
This function creates a fine tuned model with a particular doc
Returns the (status, model id, model name, error if any)
*/
type FineTuneJob struct {
	JobStatus string `json:"job_status"`
	ModelId   string `json:"model_id"`
	ModelName string `json:"model_name"`
}

func GetFineTunedModel(trainingFile string, model openai.FineTuningJobNewParamsModel) (*FineTuneJob, error) {
	fd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error while getting working directory %w ", err)
	}
	file, err := os.OpenFile(path.Join(fd, "data", trainingFile), os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error while opening file %w ", err)
	}
	defer file.Close()

	checksum, err := createChecksum(file)
	if err != nil {
		return nil, fmt.Errorf("error while creating checksum %w ", err)
	}

	f, exist := db.Get(checksum, "fine_tuned_models")

	if exist {
		var res FineTuneJob
		bz, err := json.Marshal(f)
		if err != nil {
			return nil, fmt.Errorf("error while marshalling fine tuned model from db %w ", err)
		}
		err = json.Unmarshal(bz, &res)
		if err != nil {
			return nil, fmt.Errorf("error while un-marshalling fine tuned model from db %w ", err)
		}
		// check the status of the fine tuned model
		if res.JobStatus == "succeeded" {
			return &res, nil
		} else if res.JobStatus != "failed" {
			// model is not ready - check the status
			job, err := client.FineTuning.Jobs.Get(ctx, res.ModelId)
			if err != nil {
				return nil, fmt.Errorf("error while getting fine tuned model status %w ", err)
			}
			t := FineTuneJob{
				JobStatus: string(job.Status),
				ModelId:   job.ID,
				ModelName: job.FineTunedModel,
			}
			db.Set(checksum, &t, "fine_tuned_models")
			return &t, nil
		}
		return &res, nil
	} else {
		fd, err := GetFileUploadFD(trainingFile)
		if err != nil {
			return nil, fmt.Errorf("error while getting file upload descriptor %w ", err)
		}
		job, err := CreateFineTunedModel(openai.FineTuningJobNewParams{
			Model:        openai.Raw[openai.FineTuningJobNewParamsModel](openai.FineTuningJobNewParamsModelGPT3_5Turbo),
			TrainingFile: openai.String(fd),
			Suffix:       openai.String("motivational_speaker"),
		})
		if err != nil {
			return nil, fmt.Errorf("error while creating a new fine tuned model %w ", err)
		}
		db.Set(checksum, job, "fine_tuned_models")

		return job, nil
	}
}

/*
This function checks if there exists a finetuned model trained with a particular doc
Returns the (FineTuneJob, model name, error if any)
Returns the model name and id if the model is ready or the model id only if the model is not ready
*/
func CreateFineTunedModel(opts openai.FineTuningJobNewParams) (*FineTuneJob, error) {
	fineTune, err := client.FineTuning.Jobs.New(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("error while creating fine tuning job %w ", err)
	}
	if fineTune.Status == "failed" {
		return nil, fmt.Errorf("error while fine tuning model %w ", err)
	}
	return &FineTuneJob{
		JobStatus: string(fineTune.Status),
		ModelId:   fineTune.ID,
		ModelName: fineTune.FineTunedModel,
	}, nil
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
