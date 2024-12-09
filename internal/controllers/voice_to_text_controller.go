package controllers

import (
	"context"
	"fmt"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"google.golang.org/api/option"
)

type VoiceToTextController struct{}

func (v *VoiceToTextController) ConvertVoiceToText(audioFilePath string) (string, error) {
	// Open the audio file
	file, err := os.Open(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %v", err)
	}
	defer file.Close()

	// Decode the WAV file
	decoder := wav.NewDecoder(file)

	// Check if the file is valid
	if !decoder.IsValidFile() {
		return "", fmt.Errorf("invalid WAV file")
	}

	// Read the WAV file header and decode the audio data
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return "", fmt.Errorf("failed to decode WAV file: %v", err)
	}

	// Check if it's stereo (2 channels)
	if buf.Format.NumChannels != 2 {
		return v.ConvertVoiceToTextNormal(audioFilePath)
	}

	// Convert stereo to mono by averaging the left and right channels
	monoData := make([]int16, len(buf.Data)/2)
	for i := 0; i < len(monoData); i++ {
		// Average left and right channels
		monoData[i] = (int16(buf.Data[i*2]) + int16(buf.Data[i*2+1])) / 2
	}

	monoDataInt := make([]int, len(monoData))
	for i, v := range monoData {
		monoDataInt[i] = int(v)
	}

	// Create a new file for the mono WAV file
	monoFilePath := "mono_" + audioFilePath
	monoFile, err := os.Create(monoFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create mono file: %v", err)
	}
	defer monoFile.Close()

	// Create a new WAV encoder for mono file
	encoder := wav.NewEncoder(monoFile, int(decoder.SampleRate), 16, 1, 1) // Convert SampleRate to int

	monoBuffer := &audio.IntBuffer{Data: monoDataInt, Format: &audio.Format{SampleRate: int(decoder.SampleRate), NumChannels: 1}}

	// Write the mono data to the encoder
	if err := encoder.Write(monoBuffer); err != nil {
		return "", fmt.Errorf("failed to write mono WAV file: %v", err)
	}

	// Use the mono file for speech-to-text
	return v.transcribeMonoFile(monoFilePath)
}

func (v *VoiceToTextController) transcribeMonoFile(monoFilePath string) (string, error) {
	ctx := context.Background()

	// Read credentials from environment variable
	credsJSON := os.Getenv("CRED_JSON")
	if credsJSON == "" {
		return "", fmt.Errorf("credentials not configured")
	}

	// Save credentials to a temporary file
	credsByteSlice := []byte(credsJSON)
	tmpFile, err := os.CreateTemp("", "creds-*.json")
	if err != nil {
		return "", fmt.Errorf("temp file creation failed: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write credentials to temp file
	if _, err := tmpFile.Write(credsByteSlice); err != nil {
		return "", fmt.Errorf("credential write failed: %v", err)
	}
	tmpFile.Close()

	// Create a new Speech client
	client, err := speech.NewClient(ctx, option.WithCredentialsFile(tmpFile.Name()))
	if err != nil {
		return "", fmt.Errorf("speech client creation failed: %v", err)
	}
	defer client.Close()

	// Read the mono WAV file into memory
	data, err := os.ReadFile(monoFilePath)
	if err != nil {
		return "", fmt.Errorf("audio file read failed: %v", err)
	}

	// Configure the request to Google Cloud Speech API
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:          speechpb.RecognitionConfig_LINEAR16,
			LanguageCode:      "en-US",
			AudioChannelCount: 1, // mono
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	// Send the request to Google Cloud Speech API
	resp, err := client.Recognize(ctx, req)
	if err != nil {
		return "", fmt.Errorf("speech recognition failed: %v", err)
	}

	// Collect all transcriptions from the response
	var resultText string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			resultText += alt.Transcript + " "
		}
	}

	// Return the result text (transcribed speech)
	return resultText, nil
}

func (v *VoiceToTextController) ConvertVoiceToTextNormal(audioFilePath string) (string, error) {
	ctx := context.Background()

	credsJSON := os.Getenv("CRED_JSON")
	if credsJSON == "" {
		return "", fmt.Errorf("credentials not configured")
	}

	credsByteSlice := []byte(credsJSON)
	tmpFile, err := os.CreateTemp("", "creds-*.json")
	if err != nil {
		return "", fmt.Errorf("temp file creation failed: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(credsByteSlice); err != nil {
		return "", fmt.Errorf("credential write failed: %v", err)
	}
	tmpFile.Close()

	client, err := speech.NewClient(ctx, option.WithCredentialsFile(tmpFile.Name()))
	if err != nil {
		return "", fmt.Errorf("speech client creation failed: %v", err)
	}
	defer client.Close()

	data, err := os.ReadFile(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("audio file read failed: %v", err)
	}

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:          speechpb.RecognitionConfig_LINEAR16,
			LanguageCode:      "en-US",
			AudioChannelCount: 1,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	resp, err := client.Recognize(ctx, req)
	if err != nil {
		return "", fmt.Errorf("speech recognition failed: %v", err)
	}

	var resultText string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			resultText += alt.Transcript
		}
	}

	return resultText, nil
}
