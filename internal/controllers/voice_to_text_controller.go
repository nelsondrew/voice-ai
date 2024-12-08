package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/api/option"
)

type VoiceToTextController struct{}

// ConvertVoiceToText converts an audio file to text using Google Cloud Speech-to-Text API
func (v *VoiceToTextController) ConvertVoiceToText(audioFilePath string) (string, error) {
	// Set up the client context
	ctx := context.Background()

	// Retrieve the credentials JSON string from the environment variable
	credsJSON := os.Getenv("CRED_JSON")
	if credsJSON == "" {
		return "", fmt.Errorf("CRED_JSON environment variable is not set")
	}

	// Convert the credentials JSON string to a byte slice
	credsByteSlice := []byte(credsJSON)

	// Create a temporary file to store the credentials
	tmpFile, err := ioutil.TempFile("", "creds-*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary credentials file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Ensure the temporary file is removed after use

	// Write the credentials to the temporary file
	if _, err := tmpFile.Write(credsByteSlice); err != nil {
		return "", fmt.Errorf("failed to write credentials to temporary file: %v", err)
	}

	// Create the Speech client using the credentials file
	client, err := speech.NewClient(ctx, option.WithCredentialsFile(tmpFile.Name()))
	if err != nil {
		return "", fmt.Errorf("failed to create speech client: %v", err)
	}
	defer client.Close()

	// Read the audio file from the disk
	data, err := ioutil.ReadFile(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read audio file: %v", err)
	}

	// Configure the recognition request
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:     speechpb.RecognitionConfig_LINEAR16, // or speechpb.RecognitionConfig_FLAC if the file is FLAC
			LanguageCode: "en-US",                             // Language of the speech in the audio
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	// Perform the speech-to-text recognition
	resp, err := client.Recognize(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to recognize speech: %v", err)
	}

	// Extract and combine all recognized text
	var resultText string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			resultText += alt.Transcript
		}
	}

	// Return the recognized text
	return resultText, nil
}
